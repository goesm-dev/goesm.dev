// Package search queries the index that press.Site.SearchIndex writes.
//
// It runs in the browser, compiled by goesm into the search dialog, so it
// imports no other package: with goesm, importing strings or sort also brings
// in the unicode tables and reflection, which would multiply the size of the
// island. The few string helpers it needs are below; they work on bytes and
// lower-case ASCII only, which is what search needs (Japanese has no case).
package search

// Index is a parsed search index.
type Index struct {
	records []record
}

type record struct {
	link   string
	page   string
	titles []string
	text   string
	// Lower-cased copies for matching.
	lTitles string
	lText   string
	lLast   string // last heading
}

// Result is one match.
type Result struct {
	Link   string   `json:"link"`
	Page   string   `json:"page"`
	Titles []string `json:"titles"`
	// Snippet is the text around the first match, split into parts so the
	// UI can highlight the matched terms without building HTML.
	Snippet []Part `json:"snippet"`
}

// Part is a piece of a snippet; Mark is true for a matched term.
type Part struct {
	Text string `json:"text"`
	Mark bool   `json:"mark"`
}

// Parse reads an index: one record per line, fields separated by U+001F:
// link, page title, heading path (joined by U+001E), text.
func Parse(data string) *Index {
	idx := &Index{}
	for _, line := range split(data, '\n') {
		f := split(line, 0x1f)
		if len(f) != 4 {
			continue
		}
		titles := split(f[2], 0x1e)
		idx.records = append(idx.records, record{
			link:    f[0],
			page:    f[1],
			titles:  titles,
			text:    f[3],
			lTitles: lower(join(titles, ' ')),
			lText:   lower(f[3]),
			lLast:   lower(titles[len(titles)-1]),
		})
	}
	return idx
}

// Len returns the number of records.
func (idx *Index) Len() int { return len(idx.records) }

// Terms splits a query into lower-case terms at spaces and punctuation
// (ASCII, and the ideographic space and comma / full stop). Runs of CJK
// characters are kept whole and matched as substrings like any other term.
func Terms(q string) []string {
	q = lower(q)
	var out []string
	start := 0
	for i := 0; i < len(q); {
		n := sepLen(q, i)
		if n == 0 {
			i++
			continue
		}
		if i > start {
			out = append(out, q[start:i])
		}
		i += n
		start = i
	}
	if start < len(q) {
		out = append(out, q[start:])
	}
	return out
}

// sepLen returns the length of the separator at q[i], or 0.
func sepLen(q string, i int) int {
	c := q[i]
	if c < 0x80 {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			return 0
		}
		switch c {
		case '.', '_', '-', '/', '$', '@', '#':
			return 0
		}
		return 1
	}
	// U+3000 ideographic space, U+3001 、, U+3002 。 are E3 80 80..82.
	if c == 0xe3 && i+2 < len(q) && q[i+1] == 0x80 && q[i+2] <= 0x82 {
		return 3
	}
	return 0
}

type hit struct {
	i     int
	score int
}

// Search returns up to limit records containing every term, best first:
// matches in headings rank above matches in text, and earlier, word-start
// matches above later ones.
func (idx *Index) Search(q string, limit int) []Result {
	terms := Terms(q)
	if len(terms) == 0 {
		return nil
	}
	var top []hit // sorted by score, at most limit long
	for i := range idx.records {
		r := &idx.records[i]
		score := 0
		ok := true
		for _, t := range terms {
			s := 0
			if k := index(r.lTitles, t); k >= 0 {
				s += 100
				if k == 0 || !isWord(r.lTitles[k-1]) {
					s += 50
				}
				if r.lLast == t {
					s += 100
				}
			}
			if k := index(r.lText, t); k >= 0 {
				s += 10 + max(0, 20-k/50)
				if k == 0 || !isWord(r.lText[k-1]) {
					s += 10
				}
				s += count(r.lText, t, 10)
			}
			if s == 0 {
				ok = false
				break
			}
			score += s
		}
		if !ok {
			continue
		}
		// Shallower sections (page tops) first among equals.
		score -= len(r.titles)
		top = insert(top, hit{i, score}, limit)
	}
	out := make([]Result, 0, len(top))
	for _, h := range top {
		r := &idx.records[h.i]
		out = append(out, Result{Link: r.link, Page: r.page, Titles: r.titles, Snippet: snip(r.text, r.lText, terms)})
	}
	return out
}

// insert adds h to the sorted list top, keeping at most limit entries.
// Among equal scores the earlier record stays first.
func insert(top []hit, h hit, limit int) []hit {
	pos := len(top)
	for pos > 0 && top[pos-1].score < h.score {
		pos--
	}
	if pos >= limit {
		return top
	}
	if len(top) < limit {
		top = append(top, hit{})
	}
	copy(top[pos+1:], top[pos:])
	top[pos] = h
	return top
}

func isWord(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c >= 0x80
}

const snippetLen = 160

// snip cuts about snippetLen bytes of text around the first term match and
// splits it at the term matches inside it.
func snip(text, low string, terms []string) []Part {
	first := -1
	for _, t := range terms {
		if k := index(low, t); k >= 0 && (first < 0 || k < first) {
			first = k
		}
	}
	start := 0
	if first > 40 {
		start = first - 40
	}
	end := min(len(text), start+snippetLen)
	// Do not cut UTF-8 sequences.
	for start > 0 && text[start]&0xc0 == 0x80 {
		start--
	}
	for end < len(text) && text[end]&0xc0 == 0x80 {
		end++
	}
	// Matched ranges in text[start:end], merged where they overlap.
	var spans []span
	lo := low[start:end]
	for _, t := range terms {
		for k := 0; ; {
			i := index(lo[k:], t)
			if i < 0 {
				break
			}
			spans = append(spans, span{start + k + i, start + k + i + len(t)})
			k += i + len(t)
		}
	}
	marks := merge(spans)
	var parts []Part
	if start > 0 {
		parts = append(parts, Part{Text: "… "})
	}
	at := start
	for i := 0; i+1 < len(marks); i += 2 {
		if marks[i] > at {
			parts = append(parts, Part{Text: text[at:marks[i]]})
		}
		parts = append(parts, Part{Text: text[marks[i]:marks[i+1]], Mark: true})
		at = marks[i+1]
	}
	if at < end {
		parts = append(parts, Part{Text: text[at:end]})
	}
	if end < len(text) {
		parts = append(parts, Part{Text: " …"})
	}
	return parts
}

type span struct{ s, e int }

// merge sorts spans by start and returns the merged ranges as s, e pairs.
func merge(spans []span) []int {
	for i := 1; i < len(spans); i++ {
		for j := i; j > 0 && spans[j].s < spans[j-1].s; j-- {
			spans[j], spans[j-1] = spans[j-1], spans[j]
		}
	}
	var out []int
	for _, sp := range spans {
		if n := len(out); n > 0 && sp.s <= out[n-1] {
			out[n-1] = max(out[n-1], sp.e)
			continue
		}
		out = append(out, sp.s, sp.e)
	}
	return out
}

// Byte-string helpers.

func lower(s string) string {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c >= 'A' && c <= 'Z' {
			b := []byte(s)
			for j := i; j < len(b); j++ {
				if c := b[j]; c >= 'A' && c <= 'Z' {
					b[j] = c + 'a' - 'A'
				}
			}
			return string(b)
		}
	}
	return s
}

func index(s, sub string) int {
	n := len(sub)
	if n == 0 {
		return 0
	}
	c := sub[0]
	for i := 0; i+n <= len(s); i++ {
		if s[i] == c && s[i:i+n] == sub {
			return i
		}
	}
	return -1
}

func count(s, sub string, limit int) int {
	n := 0
	for k := 0; n < limit; n++ {
		i := index(s[k:], sub)
		if i < 0 {
			break
		}
		k += i + len(sub)
	}
	return n
}

func split(s string, sep byte) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

func join(parts []string, sep byte) string {
	n := 0
	for _, p := range parts {
		n += len(p) + 1
	}
	b := make([]byte, 0, n)
	for i, p := range parts {
		if i > 0 {
			b = append(b, sep)
		}
		b = append(b, p...)
	}
	return string(b)
}
