package press

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"goesm.dev/press/highlight"
)

// The Markdown dialect is CommonMark + GitHub Flavored Markdown (tables,
// strikethrough, autolinks) through goldmark, plus what VitePress adds:
//
//   - heading ids compatible with GitHub's anchors, and permalink anchors
//   - ::: tip / info / warning / danger / details containers
//   - GitHub alerts (> [!NOTE]) rendered as the same containers
//   - highlighted code blocks with a language label and a copy button,
//     and {1,3-4} line highlights in the info string
//   - links to other .md files rewritten to page routes
//   - <!--@include: ./file.md--> (expanded before parsing, see site.go)

// renderEnv carries what the AST transformers need about the page being
// rendered, through goldmark's parser context.
type renderEnv struct {
	site     *Site
	page     *Page
	slugs    map[string]int
	headings []Heading
}

var envKey = parser.NewContextKey()

func newMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.Linkify,
			// A line break between Japanese characters is not a space.
			extension.NewCJK(extension.WithEastAsianLineBreaks(extension.EastAsianLineBreaksCSS3Draft))),
		goldmark.WithParserOptions(
			parser.WithBlockParsers(util.Prioritized(&containerParser{}, 50)),
			parser.WithASTTransformers(
				util.Prioritized(alertTransformer{}, 100),
				util.Prioritized(headingTransformer{}, 200),
				util.Prioritized(linkTransformer{}, 300),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			renderer.WithNodeRenderers(util.Prioritized(&nodeRenderer{}, 100)),
		),
	)
}

// render converts the Markdown body of p to HTML.
func (s *Site) render(p *Page, body string) (string, []Heading, []section) {
	env := &renderEnv{site: s, page: p, slugs: map[string]int{}}
	ctx := parser.NewContext()
	ctx.Set(envKey, env)
	src := []byte(body)
	doc := s.md.Parser().Parse(text.NewReader(src), parser.WithContext(ctx))
	var buf bytes.Buffer
	if err := s.md.Renderer().Render(&buf, src, doc); err != nil {
		return "<pre>" + escapeHTML(err.Error()) + "</pre>", nil, nil
	}
	return buf.String(), env.headings, sections(doc, src)
}

func envOf(pc parser.Context) *renderEnv {
	e, _ := pc.Get(envKey).(*renderEnv)
	return e
}

// Heading is an entry of a page's outline.
type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
	ID    string `json:"id"`
}

// plainText returns the text of an inline tree, without markup.
func plainText(n ast.Node, src []byte) string {
	var b strings.Builder
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			switch c := c.(type) {
			case *ast.Text:
				b.Write(c.Segment.Value(src))
				if c.SoftLineBreak() || c.HardLineBreak() {
					b.WriteByte(' ')
				}
			case *ast.String:
				b.Write(c.Value)
			case *ast.RawHTML:
			case *ast.AutoLink:
				b.Write(c.Label(src))
			default:
				walk(c)
			}
		}
	}
	walk(n)
	return b.String()
}

// Slug returns GitHub's anchor for a heading text: lower case, punctuation
// removed, spaces turned into hyphens. Letters of every script are kept.
func Slug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) || r == '_' || r == '-':
			b.WriteRune(r)
		case r == ' ':
			b.WriteByte('-')
		}
	}
	return b.String()
}

type headingTransformer struct{}

func (headingTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	env := envOf(pc)
	src := reader.Source()
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		h, ok := n.(*ast.Heading)
		if !ok || !entering {
			return ast.WalkContinue, nil
		}
		txt := strings.TrimSpace(plainText(h, src))
		var id string
		if v, ok := h.AttributeString("id"); ok {
			id = string(v.([]byte))
		} else {
			base := Slug(txt)
			id = base
			if k := env.slugs[base]; k > 0 {
				id = base + "-" + itoa(k)
			}
			env.slugs[base]++
			h.SetAttributeString("id", []byte(id))
		}
		if h.Level == 1 && env.page.Title == "" {
			env.page.Title = txt
		}
		if h.Level >= 2 && h.Level <= 3 {
			env.headings = append(env.headings, Heading{Level: h.Level, Text: txt, ID: id})
		}
		return ast.WalkSkipChildren, nil
	})
}

// Container is a ::: block or a GitHub alert.
type Container struct {
	ast.BaseBlock
	Variant string // tip, info, warning, danger, details, note, important, caution
	Title   string
	fence   int
}

var KindContainer = ast.NewNodeKind("Container")

func (n *Container) Kind() ast.NodeKind { return KindContainer }
func (n *Container) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{"Variant": n.Variant, "Title": n.Title}, nil)
}

var containerKinds = map[string]bool{
	"tip": true, "info": true, "warning": true, "danger": true, "details": true,
	"note": true, "important": true, "caution": true,
}

type containerParser struct{}

func (*containerParser) Trigger() []byte { return []byte{':'} }

func (*containerParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	s := strings.TrimRight(string(line), "\r\n")
	w, pos := util.IndentWidth(line, reader.LineOffset())
	if w > 3 {
		return nil, parser.NoChildren
	}
	s = s[pos:]
	fence := 0
	for fence < len(s) && s[fence] == ':' {
		fence++
	}
	if fence < 3 {
		return nil, parser.NoChildren
	}
	rest := strings.TrimSpace(s[fence:])
	kind, title, _ := strings.Cut(rest, " ")
	if !containerKinds[kind] {
		return nil, parser.NoChildren
	}
	title = strings.TrimSpace(title)
	if title == "" {
		if env := envOf(pc); env != nil {
			title = env.page.Locale.UI.Containers[kind]
		}
	}
	reader.AdvanceToEOL()
	return &Container{Variant: kind, Title: title, fence: fence}, parser.HasChildren
}

func (*containerParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	s := strings.TrimSpace(string(line))
	if s == strings.Repeat(":", node.(*Container).fence) {
		reader.AdvanceToEOL()
		return parser.Close
	}
	return parser.Continue | parser.HasChildren
}

func (*containerParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}
func (*containerParser) CanInterruptParagraph() bool                                { return true }
func (*containerParser) CanAcceptIndentedLine() bool                                { return false }

// alertTransformer turns GitHub alerts (a blockquote starting with [!NOTE],
// [!TIP], [!IMPORTANT], [!WARNING] or [!CAUTION]) into containers.
type alertTransformer struct{}

func (alertTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	env := envOf(pc)
	src := reader.Source()
	var quotes []*ast.Blockquote
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if q, ok := n.(*ast.Blockquote); ok && entering {
			quotes = append(quotes, q)
		}
		return ast.WalkContinue, nil
	})
	for _, q := range quotes {
		para, ok := q.FirstChild().(*ast.Paragraph)
		if !ok || para.Lines().Len() == 0 {
			continue
		}
		first := para.Lines().At(0)
		marker := strings.TrimSpace(string(first.Value(src)))
		if !strings.HasPrefix(marker, "[!") || !strings.HasSuffix(marker, "]") {
			continue
		}
		kind := strings.ToLower(marker[2 : len(marker)-1])
		if kind != "note" && kind != "tip" && kind != "important" && kind != "warning" && kind != "caution" {
			continue
		}
		// Drop the inline nodes of the marker line.
		for c := para.FirstChild(); c != nil; {
			next := c.NextSibling()
			start, end := inlineSpan(c)
			if start >= first.Stop || start < 0 && end < 0 {
				break
			}
			para.RemoveChild(para, c)
			if t, ok := c.(*ast.Text); ok && (t.SoftLineBreak() || t.HardLineBreak()) {
				break
			}
			c = next
		}
		c := &Container{Variant: kind, Title: env.page.Locale.UI.Containers[kind]}
		for ch := q.FirstChild(); ch != nil; {
			next := ch.NextSibling()
			if ch == para && para.ChildCount() == 0 {
				q.RemoveChild(q, ch)
			} else {
				c.AppendChild(c, ch)
			}
			ch = next
		}
		q.Parent().ReplaceChild(q.Parent(), q, c)
	}
}

// inlineSpan returns the source range of an inline node's first text.
func inlineSpan(n ast.Node) (int, int) {
	if t, ok := n.(*ast.Text); ok {
		return t.Segment.Start, t.Segment.Stop
	}
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if s, e := inlineSpan(c); s >= 0 {
			return s, e
		}
	}
	return -1, -1
}

// linkTransformer rewrites link and image destinations: .md files to page
// routes, repository paths of synced documents to GitHub or to the copied
// assets, and marks external links to open in a new tab.
type linkTransformer struct{}

func (linkTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	env := envOf(pc)
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.Link:
			dest, external := env.site.resolveLink(env.page, string(n.Destination), false)
			n.Destination = []byte(dest)
			if external {
				n.SetAttributeString("target", []byte("_blank"))
				n.SetAttributeString("rel", []byte("noreferrer"))
			}
		case *ast.Image:
			dest, _ := env.site.resolveLink(env.page, string(n.Destination), true)
			n.Destination = []byte(dest)
		}
		return ast.WalkContinue, nil
	})
}

// nodeRenderer renders headings with permalinks, containers and code blocks.
type nodeRenderer struct{}

func (r *nodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHeading, r.heading)
	reg.Register(KindContainer, r.container)
	reg.Register(ast.KindFencedCodeBlock, r.code)
	reg.Register(ast.KindCodeBlock, r.code)
}

func (r *nodeRenderer) heading(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	tag := "h" + itoa(n.Level)
	if entering {
		w.WriteString("<" + tag)
		html.RenderAttributes(w, n, html.HeadingAttributeFilter)
		w.WriteString(` tabindex="-1">`)
		return ast.WalkContinue, nil
	}
	if id, ok := n.AttributeString("id"); ok {
		w.WriteString(` <a class="header-anchor" href="#`)
		w.Write(id.([]byte))
		w.WriteString(`" aria-hidden="true">#</a>`)
	}
	w.WriteString("</" + tag + ">\n")
	return ast.WalkContinue, nil
}

func (r *nodeRenderer) container(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*Container)
	if n.Variant == "details" {
		if entering {
			w.WriteString(`<details class="custom-block details"><summary>` + escapeHTML(n.Title) + "</summary>\n")
		} else {
			w.WriteString("</details>\n")
		}
		return ast.WalkContinue, nil
	}
	if entering {
		w.WriteString(`<div class="custom-block ` + n.Variant + `"><p class="custom-block-title">` + escapeHTML(n.Title) + "</p>\n")
	} else {
		w.WriteString("</div>\n")
	}
	return ast.WalkContinue, nil
}

func (r *nodeRenderer) code(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	lang, info := "", ""
	if f, ok := node.(*ast.FencedCodeBlock); ok && f.Info != nil {
		info = strings.TrimSpace(string(f.Info.Segment.Value(src)))
		lang, _, _ = strings.Cut(info, " ")
		if i := strings.IndexByte(lang, '{'); i >= 0 {
			lang = lang[:i]
		}
	}
	var code strings.Builder
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		code.Write(seg.Value(src))
	}
	body := strings.TrimSuffix(code.String(), "\n")
	hl := highlight.Code(lang, body)
	if marks := lineMarks(info); len(marks) > 0 {
		out := strings.Split(hl, "\n")
		for i := range out {
			if marks[i+1] {
				out[i] = `<span class="line highlighted">` + out[i] + "</span>"
			}
		}
		hl = strings.Join(out, "\n")
	}
	label := highlight.Label(lang)
	if label == "" {
		label = "text"
	}
	w.WriteString(`<div class="language-` + escapeHTML(label) + ` vp-code"><button title="Copy" class="copy"></button><span class="lang">` + escapeHTML(label) + "</span><pre><code>")
	w.WriteString(hl)
	w.WriteString("</code></pre></div>\n")
	return ast.WalkSkipChildren, nil
}

// lineMarks parses the {1,3-5} line ranges of a code block's info string.
func lineMarks(info string) map[int]bool {
	i := strings.IndexByte(info, '{')
	j := strings.IndexByte(info, '}')
	if i < 0 || j < i {
		return nil
	}
	marks := map[int]bool{}
	for _, part := range strings.Split(info[i+1:j], ",") {
		a, b, isRange := strings.Cut(strings.TrimSpace(part), "-")
		from, to := atoi(a), atoi(a)
		if isRange {
			to = atoi(b)
		}
		for k := from; k > 0 && k <= to; k++ {
			marks[k] = true
		}
	}
	return marks
}

func escapeHTML(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '&':
			b.WriteString("&amp;")
		case '"':
			b.WriteString("&quot;")
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func itoa(n int) string { return strconv.Itoa(n) }

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
