package press

import (
	"path"
	"regexp"
	"sort"
	"strings"
)

// Pages for AI agents and other programs: every page as Markdown, and the
// llms.txt index of a locale (https://llmstxt.org/).

// MarkdownPath returns the path of the Markdown version of a route:
// /guide/x/ -> /guide/x.md, / -> /index.md, /ja/ -> /ja.md.
func MarkdownPath(route string) string {
	r, _, _ := strings.Cut(route, "#")
	r = strings.TrimSuffix(r, "/")
	if r == "" {
		r = "/index"
	}
	return r + ".md"
}

// LLMsPath returns the path of a locale's llms.txt ("/ja/llms.txt"), or of
// llms-full.txt with full.
func (s *Site) LLMsPath(loc *Locale, full bool) string {
	name := "llms.txt"
	if full {
		name = "llms-full.txt"
	}
	return s.Home(loc) + name
}

var (
	mdLinkRE   = regexp.MustCompile(`(!?)\[([^\]\n]*(?:\[[^\]\n]*\][^\]\n]*)*)\]\(([^)\s]+)\)`)
	htmlAttrRE = regexp.MustCompile(`(src|srcset|href)="(/[^"]*)"`)
)

// Markdown returns the page as self-contained Markdown: includes expanded,
// containers turned into block quotes, links absolute (pages link to their
// Markdown version) and a title with the page's description.
func (s *Site) Markdown(p *Page) string {
	var b strings.Builder
	if p.Layout == "home" {
		s.homeMarkdown(&b, p)
	}
	lines := strings.Split(strings.TrimRight(p.body, "\n"), "\n")
	titled := p.Layout == "home"
	inFence := ""
	quote := 0 // depth of open ::: containers
	for _, line := range lines {
		t := strings.TrimSpace(line)
		code := true
		switch {
		case inFence != "":
			if strings.HasPrefix(t, inFence) && strings.Trim(t, inFence[:1]) == "" {
				inFence = ""
			}
		case strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~"):
			n := 0
			for n < len(t) && t[n] == t[0] {
				n++
			}
			inFence = t[:n]
			// Drop the line highlights: ```go {2,4-5} -> ```go
			if i := strings.Index(line, "{"); i > 0 {
				line = strings.TrimRight(line[:i], " ")
			}
		case strings.HasPrefix(t, "<!--") && strings.HasSuffix(t, "-->"):
			continue
		case strings.HasPrefix(t, ":::") && strings.Trim(t, ":") == "":
			if quote > 0 {
				quote--
			}
			continue
		case strings.HasPrefix(t, ":::"):
			kind, title, _ := strings.Cut(strings.TrimSpace(strings.TrimLeft(t, ":")), " ")
			if !containerKinds[kind] {
				code = false
				break
			}
			if title = strings.TrimSpace(title); title == "" {
				title = p.Locale.UI.Containers[kind]
			}
			b.WriteString(strings.Repeat("> ", quote) + "> **" + title + "**\n" + strings.TrimRight(strings.Repeat("> ", quote+1), " ") + "\n")
			quote++
			continue
		case !titled && strings.HasPrefix(t, "# "):
			titled = true
			b.WriteString(line + "\n")
			if p.Description != "" {
				b.WriteString("\n> " + p.Description + "\n")
			}
			continue
		default:
			code = false
		}
		if !titled && t != "" && p.Title != "" {
			// A page without a # heading: its title comes first.
			titled = true
			b.WriteString("# " + p.Title + "\n\n")
			if p.Description != "" {
				b.WriteString("> " + p.Description + "\n\n")
			}
		}
		if !code {
			line = s.absoluteLinks(p, line)
		}
		if quote > 0 {
			line = strings.TrimRight(strings.Repeat("> ", quote)+line, " ")
		}
		// Dropped lines (comments, container fences) leave runs of blank
		// lines behind; keep one.
		if !code && inFence == "" && strings.TrimSpace(line) == "" && strings.HasSuffix(b.String(), "\n\n") {
			continue
		}
		b.WriteString(line + "\n")
	}
	return strings.Trim(b.String(), "\n") + "\n"
}

func (s *Site) homeMarkdown(b *strings.Builder, p *Page) {
	h := p.Hero
	title := h.Name
	if h.Text != "" {
		title += ": " + h.Text
	}
	b.WriteString("# " + title + "\n\n")
	if h.Tagline != "" {
		b.WriteString("> " + h.Tagline + "\n\n")
	}
	var actions []string
	for _, a := range h.Actions {
		actions = append(actions, "["+a.Text+"]("+s.absoluteURL(a.Link)+")")
	}
	if len(actions) > 0 {
		b.WriteString(strings.Join(actions, " · ") + "\n\n")
	}
	for _, f := range p.Features {
		name := "**" + f.Title + "**"
		if f.Link != "" {
			name = "[" + name + "](" + s.absoluteURL(f.Link) + ")"
		}
		b.WriteString("- " + name + ": " + f.Details + "\n")
	}
	if len(p.Features) > 0 {
		b.WriteString("\n")
	}
}

// absoluteLinks rewrites the links and HTML src/href attributes of a line
// of page p to absolute URLs.
func (s *Site) absoluteLinks(p *Page, line string) string {
	line = mdLinkRE.ReplaceAllStringFunc(line, func(m string) string {
		sm := mdLinkRE.FindStringSubmatch(m)
		dest, ext := s.resolveLink(p, sm[3], sm[1] == "!")
		if !ext && !strings.HasPrefix(dest, "#") {
			dest = s.absoluteURL(dest)
		}
		return sm[1] + "[" + sm[2] + "](" + dest + ")"
	})
	return htmlAttrRE.ReplaceAllStringFunc(line, func(m string) string {
		sm := htmlAttrRE.FindStringSubmatch(m)
		return sm[1] + `="` + s.Config.SiteURL + sm[2] + `"`
	})
}

// absoluteURL turns a path of the site into a URL: the Markdown version for
// a page route, the file itself for anything with an extension.
func (s *Site) absoluteURL(link string) string {
	if link == "" || strings.Contains(link, "://") || !strings.HasPrefix(link, "/") {
		return link
	}
	p, frag, hasFrag := strings.Cut(link, "#")
	if path.Ext(p) == "" {
		if hasFrag {
			frag = "#" + frag
		}
		return s.Config.SiteURL + MarkdownPath(p) + frag
	}
	return s.Config.SiteURL + link
}

// LLMsSection is a list of links in llms.txt.
type LLMsSection struct {
	Title string
	Links []LLMsLink
}

type LLMsLink struct {
	Text, URL, Notes string
}

// LLMsTxt returns the llms.txt of a locale: the site's title and
// description, then its pages grouped as in the sidebar, then extra.
func (s *Site) LLMsTxt(loc *Locale, extra ...LLMsSection) string {
	var b strings.Builder
	b.WriteString("# " + loc.Title + "\n\n> " + loc.Description + "\n\n")
	if home := s.Page(s.Home(loc)); home != nil {
		b.WriteString("[" + home.Title + "](" + s.Config.SiteURL + MarkdownPath(home.Route) + "): " + firstNonEmpty(home.Hero.Tagline, home.Description) + "\n\n")
	}
	for _, g := range s.llmsGroups(loc) {
		writeSection(&b, g)
	}
	for _, g := range extra {
		writeSection(&b, g)
	}
	return b.String()
}

func writeSection(b *strings.Builder, g LLMsSection) {
	b.WriteString("## " + g.Title + "\n\n")
	for _, l := range g.Links {
		b.WriteString("- [" + l.Text + "](" + l.URL + ")")
		if l.Notes != "" {
			b.WriteString(": " + l.Notes)
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

// LLMsFullTxt returns every page of a locale as Markdown, in llms.txt order.
func (s *Site) LLMsFullTxt(loc *Locale) string {
	var b strings.Builder
	b.WriteString("# " + loc.Title + "\n\n> " + loc.Description + "\n")
	for _, p := range s.llmsPages(loc) {
		b.WriteString("\n---\n\nSource: " + s.Config.SiteURL + p.Route + "\n\n")
		b.WriteString(s.Markdown(p))
	}
	return b.String()
}

// llmsGroups returns the sidebar groups of a locale with their pages, and
// the locale's other pages (except the home page) in a last group.
func (s *Site) llmsGroups(loc *Locale) []LLMsSection {
	seen := map[string]bool{s.Home(loc): true}
	var out []LLMsSection
	add := func(g *LLMsSection, text string, p *Page) {
		seen[p.Route] = true
		g.Links = append(g.Links, LLMsLink{Text: text, URL: s.Config.SiteURL + MarkdownPath(p.Route), Notes: p.Description})
	}
	byTitle := map[string]int{}
	for _, prefix := range sortedKeys(loc.Sidebar) {
		for _, g := range loc.Sidebar[prefix] {
			i, ok := byTitle[g.Text]
			if !ok {
				i = len(out)
				byTitle[g.Text] = i
				out = append(out, LLMsSection{Title: g.Text})
			}
			for _, it := range g.Items {
				if p := s.Page(s.localLink(loc, it.Link)); p != nil && !seen[p.Route] {
					add(&out[i], firstNonEmpty(it.Text, p.Title), p)
				}
			}
		}
	}
	rest := LLMsSection{Title: loc.UI.OtherPages}
	for _, p := range s.ordered {
		if p.Locale == loc && !seen[p.Route] {
			add(&rest, p.Title, p)
		}
	}
	if len(rest.Links) > 0 {
		out = append(out, rest)
	}
	return out
}

func (s *Site) llmsPages(loc *Locale) []*Page {
	var out []*Page
	if home := s.Page(s.Home(loc)); home != nil {
		out = append(out, home)
	}
	for _, g := range s.llmsGroups(loc) {
		for _, l := range g.Links {
			out = append(out, s.pages[routeOfMarkdownURL(s.Config.SiteURL, l.URL)])
		}
	}
	return out
}

// routeOfMarkdownURL inverts SiteURL + MarkdownPath(route) for page routes.
func routeOfMarkdownURL(site, u string) string {
	r := strings.TrimSuffix(strings.TrimPrefix(u, site), ".md")
	if r == "/index" {
		return "/"
	}
	return r + "/"
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func firstNonEmpty(s ...string) string {
	for _, x := range s {
		if x != "" {
			return x
		}
	}
	return ""
}
