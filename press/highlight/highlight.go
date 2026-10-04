// Package highlight turns source code into HTML with token classes.
//
// It is a small hand-written lexer per language family, enough for the code
// blocks of goesm's documentation: Go, TypeScript / JavaScript, Vue / HTML,
// shell, JSON, TOML / YAML and diffs. Unknown languages are escaped only.
//
// Tokens become <span class="hl-…"> elements: kw (keyword), str (string),
// com (comment), num (number), fn (called function), type (type name or
// builtin), attr (shell flag, HTML attribute, mapping key), tag (HTML tag),
// add / del (diff lines), prompt (shell prompt).
package highlight

import (
	"strings"
)

// Code returns the HTML of src highlighted as lang. The result contains no
// enclosing <pre> or <code>.
func Code(lang, src string) string {
	lang = strings.ToLower(lang)
	if a, ok := aliases[lang]; ok {
		lang = a
	}
	var b strings.Builder
	switch lang {
	case "go":
		clike(&b, src, goLang)
	case "ts":
		clike(&b, src, tsLang)
	case "sh":
		shell(&b, src)
	case "json":
		clike(&b, src, jsonLang)
	case "toml":
		conf(&b, src, '=')
	case "yaml":
		conf(&b, src, ':')
	case "html":
		markup(&b, src)
	case "diff":
		diff(&b, src)
	default:
		escape(&b, src)
	}
	return b.String()
}

// Label returns the display name of a language, as shown on code blocks.
func Label(lang string) string {
	l := strings.ToLower(lang)
	if a, ok := aliases[l]; ok && a == "sh" {
		return "sh"
	}
	return l
}

var aliases = map[string]string{
	"golang": "go", "go": "go", "go.mod": "toml", "gomod": "toml",
	"ts": "ts", "typescript": "ts", "js": "ts", "javascript": "ts", "mjs": "ts", "tsx": "ts", "jsx": "ts",
	"sh": "sh", "bash": "sh", "shell": "sh", "console": "sh", "zsh": "sh",
	"json": "json", "jsonc": "json",
	"toml": "toml", "yaml": "yaml", "yml": "yaml",
	"html": "html", "vue": "html", "astro": "html", "xml": "html", "svg": "html",
	"diff": "diff",
}

type langDef struct {
	keywords map[string]bool
	types    map[string]bool
	consts   map[string]bool
	line     []string // line comment starts
	block    bool     // /* */ comments
	raw      byte     // raw / template string delimiter (` in Go and JS)
	single   bool     // '...' is a string (JS), not a rune (Go: also string-like)
}

func set(words string) map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(words) {
		m[w] = true
	}
	return m
}

var goLang = &langDef{
	keywords: set(`break case chan const continue default defer else fallthrough for func go goto if import
		interface map package range return select struct switch type var`),
	types: set(`any bool byte comparable complex64 complex128 error float32 float64 int int8 int16 int32 int64
		rune string uint uint8 uint16 uint32 uint64 uintptr append cap clear close complex copy delete imag len
		make max min new panic print println real recover`),
	consts: set(`true false nil iota`),
	line:   []string{"//"},
	block:  true,
	raw:    '`',
	single: true,
}

var tsLang = &langDef{
	keywords: set(`abstract as async await break case catch class const continue debugger declare default delete do
		else enum export extends finally for from function get if implements import in instanceof interface
		let new of package private protected public readonly return satisfies set static super switch this throw
		try type typeof var void while with yield`),
	types: set(`any boolean number string symbol bigint unknown never object Promise Array Map Set Record
		BigInt Number String Uint8Array console`),
	consts: set(`true false null undefined NaN Infinity`),
	line:   []string{"//"},
	block:  true,
	raw:    '`',
	single: true,
}

var jsonLang = &langDef{
	keywords: map[string]bool{},
	types:    map[string]bool{},
	consts:   set(`true false null`),
	line:     []string{"//"},
	block:    true,
}

func span(b *strings.Builder, class, text string) {
	b.WriteString(`<span class="hl-`)
	b.WriteString(class)
	b.WriteString(`">`)
	escape(b, text)
	b.WriteString(`</span>`)
}

func escape(b *strings.Builder, s string) {
	start := 0
	for i := 0; i < len(s); i++ {
		var rep string
		switch s[i] {
		case '<':
			rep = "&lt;"
		case '>':
			rep = "&gt;"
		case '&':
			rep = "&amp;"
		case '"':
			rep = "&quot;"
		default:
			continue
		}
		b.WriteString(s[start:i])
		b.WriteString(rep)
		start = i + 1
	}
	b.WriteString(s[start:])
}

func isIdentStart(c byte) bool {
	return c == '_' || c == '$' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80
}

func isIdent(c byte) bool {
	return isIdentStart(c) || c >= '0' && c <= '9'
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// clike lexes C-family languages (Go, TypeScript, JSON).
func clike(b *strings.Builder, s string, l *langDef) {
	prevWord := ""
	for i := 0; i < len(s); {
		c := s[i]
		// Comments.
		if l.block && strings.HasPrefix(s[i:], "/*") {
			end := strings.Index(s[i+2:], "*/")
			j := len(s)
			if end >= 0 {
				j = i + 2 + end + 2
			}
			span(b, "com", s[i:j])
			i = j
			continue
		}
		if isLineComment(s[i:], l.line) {
			j := lineEnd(s, i)
			span(b, "com", s[i:j])
			i = j
			continue
		}
		// Strings.
		if c == '"' || (c == '\'' && l.single) || (l.raw != 0 && c == l.raw) {
			j := stringEnd(s, i, c, c != l.raw)
			class := "str"
			if l == jsonLang && c == '"' && nextNonSpace(s, j) == ':' {
				class = "attr"
			}
			span(b, class, s[i:j])
			i = j
			prevWord = ""
			continue
		}
		// Numbers.
		if isDigit(c) || (c == '.' && i+1 < len(s) && isDigit(s[i+1])) {
			j := i + 1
			for j < len(s) && (isIdent(s[j]) || s[j] == '.' || ((s[j] == '+' || s[j] == '-') && (s[j-1] == 'e' || s[j-1] == 'E' || s[j-1] == 'p' || s[j-1] == 'P'))) {
				j++
			}
			span(b, "num", s[i:j])
			i = j
			continue
		}
		// Identifiers.
		if isIdentStart(c) {
			j := i + 1
			for j < len(s) && isIdent(s[j]) {
				j++
			}
			w := s[i:j]
			switch {
			case l.keywords[w]:
				span(b, "kw", w)
			case l.consts[w]:
				span(b, "num", w)
			case prevWord == "type" || prevWord == "interface" || prevWord == "class" || prevWord == "new" || prevWord == "extends" || prevWord == "implements":
				span(b, "type", w)
			case nextNonSpace(s, j) == '(' || prevWord == "func" || prevWord == "function":
				if l.types[w] {
					span(b, "type", w)
				} else {
					span(b, "fn", w)
				}
			case l.types[w]:
				span(b, "type", w)
			default:
				escape(b, w)
			}
			prevWord = w
			i = j
			continue
		}
		if c != ' ' && c != '\t' {
			prevWord = ""
		}
		escape(b, string(c))
		i++
	}
}

func isLineComment(s string, starts []string) bool {
	for _, p := range starts {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func lineEnd(s string, i int) int {
	j := strings.IndexByte(s[i:], '\n')
	if j < 0 {
		return len(s)
	}
	return i + j
}

// stringEnd returns the index after the string literal starting at i.
// Escaped strings stop at a newline.
func stringEnd(s string, i int, q byte, escapes bool) int {
	for j := i + 1; j < len(s); j++ {
		switch {
		case escapes && s[j] == '\\':
			j++
		case s[j] == q:
			return j + 1
		case escapes && s[j] == '\n':
			return j
		}
	}
	return len(s)
}

func nextNonSpace(s string, i int) byte {
	for ; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' {
			return s[i]
		}
	}
	return 0
}

// shell highlights command lines: prompts, the command word, flags, strings
// and comments.
func shell(b *strings.Builder, s string) {
	for n, line := range strings.Split(s, "\n") {
		if n > 0 {
			b.WriteByte('\n')
		}
		shellLine(b, line)
	}
}

func shellLine(b *strings.Builder, line string) {
	i := 0
	if strings.HasPrefix(line, "$ ") {
		span(b, "prompt", "$ ")
		i = 2
	}
	cmd := true
	for i < len(line) {
		c := line[i]
		switch {
		case c == '#' && (i == 0 || line[i-1] == ' '):
			span(b, "com", line[i:])
			return
		case c == '"' || c == '\'':
			j := stringEnd(line, i, c, c == '"')
			span(b, "str", line[i:j])
			i = j
			cmd = false
		case c == ' ' || c == '\t':
			b.WriteByte(c)
			i++
		case c == '|' || c == ';' || c == '&':
			escape(b, string(c))
			i++
			cmd = true
		default:
			j := i
			for j < len(line) && line[j] != ' ' && line[j] != '\t' && line[j] != ';' && line[j] != '|' {
				j++
			}
			w := line[i:j]
			switch {
			case cmd && strings.Contains(w, "=") && !strings.HasPrefix(w, "-"):
				// VAR=value prefix
				span(b, "attr", w)
			case cmd:
				span(b, "fn", w)
				cmd = false
			case strings.HasPrefix(w, "-"):
				span(b, "attr", w)
			default:
				escape(b, w)
			}
			i = j
		}
	}
}

// conf highlights TOML and YAML: keys, strings, numbers, comments and
// [section] headers.
func conf(b *strings.Builder, s string, sep byte) {
	for n, line := range strings.Split(s, "\n") {
		if n > 0 {
			b.WriteByte('\n')
		}
		t := strings.TrimLeft(line, " \t")
		indent := line[:len(line)-len(t)]
		b.WriteString(indent)
		switch {
		case strings.HasPrefix(t, "#") || strings.HasPrefix(t, "//"):
			span(b, "com", t)
			continue
		case strings.HasPrefix(t, "[") && sep == '=':
			span(b, "type", t)
			continue
		}
		if strings.HasPrefix(t, "- ") {
			b.WriteString("- ")
			t = t[2:]
		}
		k := strings.IndexByte(t, sep)
		if k > 0 && !strings.ContainsAny(t[:k], `"'`) {
			span(b, "attr", t[:k])
			b.WriteByte(sep)
			t = t[k+1:]
		}
		value(b, t)
	}
}

func value(b *strings.Builder, t string) {
	for i := 0; i < len(t); {
		c := t[i]
		switch {
		case c == '"' || c == '\'':
			j := stringEnd(t, i, c, c == '"')
			span(b, "str", t[i:j])
			i = j
		case c == '#' && (i == 0 || t[i-1] == ' '):
			span(b, "com", t[i:])
			return
		case isDigit(c) && (i == 0 || !isIdent(t[i-1])):
			j := i
			for j < len(t) && (isIdent(t[j]) || t[j] == '.') {
				j++
			}
			span(b, "num", t[i:j])
			i = j
		default:
			j := i + 1
			for j < len(t) && t[j] != '"' && t[j] != '\'' && t[j] != '#' && !(isDigit(t[j]) && !isIdent(t[j-1])) {
				j++
			}
			w := t[i:j]
			if tw := strings.TrimSpace(w); tw == "true" || tw == "false" || tw == "null" {
				escape(b, w[:strings.Index(w, tw)])
				span(b, "num", tw)
				escape(b, w[strings.Index(w, tw)+len(tw):])
			} else {
				escape(b, w)
			}
			i = j
		}
	}
}

// markup highlights HTML-like markup (HTML, Vue SFCs, Astro): tags,
// attributes, attribute values and comments. The content of <script> blocks
// is highlighted as Go when lang="go", else as TypeScript.
func markup(b *strings.Builder, s string) {
	i := 0
	// Astro frontmatter.
	if strings.HasPrefix(s, "---\n") {
		if end := strings.Index(s[4:], "\n---"); end >= 0 {
			span(b, "com", "---")
			b.WriteByte('\n')
			clike(b, s[4:4+end+1], tsLang)
			span(b, "com", "---")
			i = 4 + end + 4
		}
	}
	for i < len(s) {
		if strings.HasPrefix(s[i:], "<!--") {
			end := strings.Index(s[i:], "-->")
			j := len(s)
			if end >= 0 {
				j = i + end + 3
			}
			span(b, "com", s[i:j])
			i = j
			continue
		}
		if s[i] == '<' && i+1 < len(s) && (isIdentStart(s[i+1]) || s[i+1] == '/') {
			j := i + 1
			if s[j] == '/' {
				j++
			}
			k := j
			for k < len(s) && (isIdent(s[k]) || s[k] == '-' || s[k] == '.' || s[k] == ':') {
				k++
			}
			escape(b, s[i:j])
			name := s[j:k]
			span(b, "tag", name)
			// Attributes up to '>'.
			tagStart := k
			for k < len(s) && s[k] != '>' {
				c := s[k]
				switch {
				case c == '"' || c == '\'':
					e := stringEnd(s, k, c, false)
					span(b, "str", s[k:e])
					k = e
				case isIdentStart(c) || c == '@' || c == ':' || c == '#' || c == 'v':
					e := k + 1
					for e < len(s) && s[e] != '=' && s[e] != ' ' && s[e] != '>' && s[e] != '\n' && s[e] != '/' {
						e++
					}
					span(b, "attr", s[k:e])
					k = e
				default:
					escape(b, string(c))
					k++
				}
			}
			if k < len(s) {
				escape(b, ">")
				k++
			}
			attrs := s[tagStart:k]
			i = k
			if (name == "script" || name == "style") && s[j-1] != '/' {
				end := strings.Index(s[i:], "</"+name)
				if end < 0 {
					end = len(s) - i
				}
				body := s[i : i+end]
				switch {
				case name == "style":
					escape(b, body)
				case strings.Contains(attrs, `lang="go"`):
					clike(b, body, goLang)
				default:
					clike(b, body, tsLang)
				}
				i += end
			}
			continue
		}
		if s[i] == '{' && i+1 < len(s) && s[i+1] == '{' {
			end := strings.Index(s[i:], "}}")
			if end >= 0 {
				escape(b, "{{")
				clike(b, s[i+2:i+end], tsLang)
				escape(b, "}}")
				i += end + 2
				continue
			}
		}
		j := i + 1
		for j < len(s) && s[j] != '<' && s[j] != '{' {
			j++
		}
		escape(b, s[i:j])
		i = j
	}
}

func diff(b *strings.Builder, s string) {
	for n, line := range strings.Split(s, "\n") {
		if n > 0 {
			b.WriteByte('\n')
		}
		switch {
		case strings.HasPrefix(line, "+"):
			span(b, "add", line)
		case strings.HasPrefix(line, "-"):
			span(b, "del", line)
		case strings.HasPrefix(line, "@@"):
			span(b, "com", line)
		default:
			escape(b, line)
		}
	}
}
