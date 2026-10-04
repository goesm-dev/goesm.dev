package press

import (
	"fmt"
	"strconv"
	"strings"
)

// splitFrontmatter separates a leading "---" YAML block from the Markdown
// body. A file without frontmatter has an empty map.
func splitFrontmatter(src string) (map[string]any, string, error) {
	if !strings.HasPrefix(src, "---\n") && !strings.HasPrefix(src, "---\r\n") {
		return map[string]any{}, src, nil
	}
	rest := src[strings.IndexByte(src, '\n')+1:]
	end := -1
	for i := 0; i < len(rest); {
		j := strings.IndexByte(rest[i:], '\n')
		line := rest[i:]
		if j >= 0 {
			line = rest[i : i+j]
		}
		if strings.TrimRight(line, "\r") == "---" {
			end = i
			break
		}
		if j < 0 {
			break
		}
		i += j + 1
	}
	if end < 0 {
		return nil, "", fmt.Errorf("frontmatter: missing closing ---")
	}
	fm, err := parseYAML(rest[:end])
	if err != nil {
		return nil, "", err
	}
	body := ""
	if i := strings.IndexByte(rest[end:], '\n'); i >= 0 {
		body = rest[end+i+1:]
	}
	return fm, body, nil
}

// parseYAML parses the subset of YAML that frontmatter uses: nested block
// mappings and sequences, plain, single- and double-quoted scalars, and
// comments. Values are string, map[string]any or []any.
func parseYAML(src string) (map[string]any, error) {
	var lines []yamlLine
	for n, raw := range strings.Split(src, "\n") {
		raw = strings.TrimRight(raw, " \t\r")
		trimmed := strings.TrimLeft(raw, " ")
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		lines = append(lines, yamlLine{indent: len(raw) - len(trimmed), text: trimmed, num: n + 1})
	}
	p := &yamlParser{lines: lines}
	if len(lines) == 0 {
		return map[string]any{}, nil
	}
	v, err := p.block(lines[0].indent)
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.lines) {
		return nil, p.errorf("unexpected indentation")
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("frontmatter: expected a mapping at the top level")
	}
	return m, nil
}

type yamlLine struct {
	indent int
	text   string
	num    int
}

type yamlParser struct {
	lines []yamlLine
	pos   int
}

func (p *yamlParser) errorf(format string, args ...any) error {
	n := 0
	if p.pos < len(p.lines) {
		n = p.lines[p.pos].num
	}
	return fmt.Errorf("frontmatter line %d: %s", n, fmt.Sprintf(format, args...))
}

// block parses the mapping or sequence whose lines start at indent.
func (p *yamlParser) block(indent int) (any, error) {
	if p.pos >= len(p.lines) {
		return "", nil
	}
	if strings.HasPrefix(p.lines[p.pos].text, "- ") || p.lines[p.pos].text == "-" {
		return p.sequence(indent)
	}
	return p.mapping(indent, map[string]any{})
}

func (p *yamlParser) mapping(indent int, m map[string]any) (map[string]any, error) {
	for p.pos < len(p.lines) {
		l := p.lines[p.pos]
		if l.indent < indent {
			break
		}
		if l.indent > indent {
			return nil, p.errorf("unexpected indentation")
		}
		key, rest, ok := splitKey(l.text)
		if !ok {
			return nil, p.errorf("expected key: value, got %q", l.text)
		}
		p.pos++
		if rest != "" {
			v, err := scalar(rest)
			if err != nil {
				return nil, p.errorf("%v", err)
			}
			m[key] = v
			continue
		}
		// A nested block, or an empty value.
		if p.pos < len(p.lines) {
			next := p.lines[p.pos]
			isSeq := strings.HasPrefix(next.text, "- ") || next.text == "-"
			if next.indent > indent || (isSeq && next.indent == indent) {
				v, err := p.block(next.indent)
				if err != nil {
					return nil, err
				}
				m[key] = v
				continue
			}
		}
		m[key] = ""
	}
	return m, nil
}

func (p *yamlParser) sequence(indent int) ([]any, error) {
	var out []any
	for p.pos < len(p.lines) {
		l := p.lines[p.pos]
		if l.indent != indent || !(strings.HasPrefix(l.text, "- ") || l.text == "-") {
			if l.indent > indent {
				return nil, p.errorf("unexpected indentation")
			}
			break
		}
		item := strings.TrimLeft(strings.TrimPrefix(l.text, "-"), " ")
		if item == "" {
			p.pos++
			v, err := p.block(p.lines[min(p.pos, len(p.lines)-1)].indent)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
			continue
		}
		if _, _, ok := splitKey(item); ok {
			// "- key: value" starts a mapping whose other keys are indented
			// to the column of key.
			col := indent + len(l.text) - len(item)
			m := map[string]any{}
			p.lines[p.pos] = yamlLine{indent: col, text: item, num: l.num}
			if _, err := p.mapping(col, m); err != nil {
				return nil, err
			}
			out = append(out, m)
			continue
		}
		v, err := scalar(item)
		if err != nil {
			return nil, p.errorf("%v", err)
		}
		out = append(out, v)
		p.pos++
	}
	return out, nil
}

// splitKey splits "key: value" (value may be empty). Keys are plain words.
func splitKey(s string) (key, rest string, ok bool) {
	if strings.HasPrefix(s, `"`) || strings.HasPrefix(s, `'`) {
		return "", "", false
	}
	i := strings.Index(s, ":")
	if i <= 0 || (i+1 < len(s) && s[i+1] != ' ') {
		return "", "", false
	}
	key = strings.TrimSpace(s[:i])
	if strings.ContainsAny(key, " \t") {
		return "", "", false
	}
	return key, strings.TrimSpace(s[i+1:]), true
}

func scalar(s string) (string, error) {
	switch {
	case strings.HasPrefix(s, `"`):
		end := closingQuote(s)
		if end < 0 {
			return "", fmt.Errorf("unterminated string %s", s)
		}
		v, err := strconv.Unquote(s[:end+1])
		if err != nil {
			return "", fmt.Errorf("bad string %s", s)
		}
		return v, nil
	case strings.HasPrefix(s, `'`):
		var b strings.Builder
		for i := 1; i < len(s); i++ {
			if s[i] == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					b.WriteByte('\'')
					i++
					continue
				}
				return b.String(), nil
			}
			b.WriteByte(s[i])
		}
		return "", fmt.Errorf("unterminated string %s", s)
	}
	if i := strings.Index(s, " #"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	return s, nil
}

func closingQuote(s string) int {
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '"':
			return i
		}
	}
	return -1
}

// Helpers to read decoded frontmatter values.

func str(m map[string]any, key string) string {
	s, _ := m[key].(string)
	return s
}

func mapOf(m map[string]any, key string) map[string]any {
	v, _ := m[key].(map[string]any)
	if v == nil {
		return map[string]any{}
	}
	return v
}

func listOf(m map[string]any, key string) []map[string]any {
	l, _ := m[key].([]any)
	var out []map[string]any
	for _, x := range l {
		if mm, ok := x.(map[string]any); ok {
			out = append(out, mm)
		}
	}
	return out
}
