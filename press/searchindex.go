package press

import (
	"strings"

	"github.com/yuin/goldmark/ast"
)

// A section is the part of a page under one heading, for search.
type section struct {
	anchor string   // heading id ("" for the text before the first heading)
	titles []string // heading path: h1, h2, h3
	text   string
}

// sections splits a rendered document at its headings (levels 1 to 3) and
// collects the plain text under each.
func sections(doc ast.Node, src []byte) []section {
	var out []section
	cur := section{}
	var path [4]string
	var text strings.Builder
	flush := func() {
		cur.text = strings.Join(strings.Fields(text.String()), " ")
		if cur.text != "" || cur.anchor != "" {
			out = append(out, cur)
		}
		text.Reset()
	}
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			switch c := c.(type) {
			case *ast.Heading:
				if c.Level > 3 {
					text.WriteString(plainText(c, src))
					text.WriteByte(' ')
					continue
				}
				flush()
				path[c.Level] = strings.TrimSpace(plainText(c, src))
				for l := c.Level + 1; l < len(path); l++ {
					path[l] = ""
				}
				var titles []string
				for _, t := range path[1:] {
					if t != "" {
						titles = append(titles, t)
					}
				}
				id, _ := c.AttributeString("id")
				b, _ := id.([]byte)
				cur = section{anchor: string(b), titles: titles}
			case *ast.FencedCodeBlock, *ast.CodeBlock:
				lines := c.Lines()
				for i := 0; i < lines.Len(); i++ {
					seg := lines.At(i)
					text.Write(seg.Value(src))
				}
				text.WriteByte(' ')
			case *ast.HTMLBlock:
			case *ast.Paragraph, *ast.TextBlock:
				text.WriteString(plainText(c, src))
				text.WriteByte(' ')
			default:
				if c.Type() == ast.TypeInline {
					text.WriteString(plainText(c, src))
				} else {
					walk(c)
					text.WriteByte(' ')
				}
			}
		}
	}
	walk(doc)
	flush()
	return out
}

// SearchIndex returns the search records of a locale's pages in a compact
// text format read by package search: one record per line, fields separated
// by U+001F: link, page title, heading path (joined by U+001E), text.
func (s *Site) SearchIndex(code string) string {
	var b strings.Builder
	for _, p := range s.ordered {
		if p.Locale.Code != code || p.Layout == "home" {
			continue
		}
		s.Render(p)
		for _, sec := range p.sections {
			link := p.Route
			if sec.anchor != "" {
				link += "#" + sec.anchor
			}
			titles := sec.titles
			if len(titles) == 0 {
				titles = []string{p.Title}
			}
			b.WriteString(link)
			b.WriteByte(0x1f)
			b.WriteString(clean(p.Title))
			b.WriteByte(0x1f)
			for i, t := range titles {
				if i > 0 {
					b.WriteByte(0x1e)
				}
				b.WriteString(clean(t))
			}
			b.WriteByte(0x1f)
			b.WriteString(clean(sec.text))
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func clean(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == 0x1e || r == 0x1f {
			return ' '
		}
		return r
	}, s)
}
