// Package press is a small VitePress-like documentation engine: Markdown
// pages with frontmatter in a content tree, one directory per locale, turned
// into routes, rendered HTML, outlines, navigation, sidebars, prev / next
// links, a search index and a sitemap.
//
// It is plain Go. goesm.dev compiles it with goesm and runs it inside Astro
// and Vue (through gosfc): the theme is Vue, the engine is this package.
package press

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
)

// Site is a loaded content tree.
type Site struct {
	Config  *Config
	fsys    fs.FS
	md      goldmark.Markdown
	pages   map[string]*Page // by route
	byFile  map[string]*Page // by content path ("en/guide/x.md")
	bySrc   map[string]*Page // by locale code + ":" + source name + ":" + path
	ordered []*Page
}

// Page is one Markdown file.
type Page struct {
	File        string // content path: "en/guide/getting-started.md"
	Route       string // "/guide/getting-started/"
	Locale      *Locale
	Title       string
	Description string
	Layout      string // "doc" (default), "home" or "page"
	Hero        Hero
	Features    []Feature
	// Source is the synced document this page is a copy of, if any.
	Source     string // source name in Config.Sources ("goesm")
	SourcePath string // path in that repository ("docs/otelc.md")

	body     string
	rendered bool
	html     string
	headings []Heading
	sections []section
}

type Hero struct {
	Name    string   `json:"name"`
	Text    string   `json:"text"`
	Tagline string   `json:"tagline"`
	Image   string   `json:"image"`
	Actions []Action `json:"actions"`
}

type Action struct {
	Theme string `json:"theme"` // brand, alt
	Text  string `json:"text"`
	Link  string `json:"link"`
}

type Feature struct {
	Icon    string `json:"icon"`
	Title   string `json:"title"`
	Details string `json:"details"`
	Link    string `json:"link"`
}

// Load reads every .md file below the locale directories of fsys. Files and
// directories whose names start with "_" are not pages (they can be
// included).
func Load(fsys fs.FS, cfg *Config) (*Site, error) {
	s := &Site{
		Config: cfg,
		fsys:   fsys,
		md:     newMarkdown(),
		pages:  map[string]*Page{},
		byFile: map[string]*Page{},
		bySrc:  map[string]*Page{},
	}
	for _, loc := range cfg.Locales {
		err := fs.WalkDir(fsys, loc.Code, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if strings.HasPrefix(d.Name(), "_") {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
			if d.IsDir() || !strings.HasSuffix(p, ".md") {
				return nil
			}
			return s.add(loc, p)
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(s.ordered, func(i, j int) bool { return s.ordered[i].Route < s.ordered[j].Route })
	return s, nil
}

func (s *Site) add(loc *Locale, file string) error {
	raw, err := fs.ReadFile(s.fsys, file)
	if err != nil {
		return err
	}
	fm, body, err := splitFrontmatter(string(raw))
	if err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}
	body, err = s.expandIncludes(file, body, 0)
	if err != nil {
		return err
	}
	p := &Page{
		File:        file,
		Route:       routeOf(loc, file),
		Locale:      loc,
		Title:       str(fm, "title"),
		Description: str(fm, "description"),
		Layout:      str(fm, "layout"),
		body:        body,
	}
	if p.Layout == "" {
		p.Layout = "doc"
	}
	if src := str(fm, "source"); src != "" {
		name, sp, ok := strings.Cut(src, ":")
		if !ok || s.Config.Sources[name].Repo == "" {
			return fmt.Errorf("%s: unknown source %q", file, src)
		}
		p.Source, p.SourcePath = name, sp
		s.bySrc[loc.Code+":"+name+":"+sp] = p
	}
	hero := mapOf(fm, "hero")
	p.Hero = Hero{Name: str(hero, "name"), Text: str(hero, "text"), Tagline: str(hero, "tagline"), Image: str(hero, "image")}
	for _, a := range listOf(hero, "actions") {
		p.Hero.Actions = append(p.Hero.Actions, Action{Theme: str(a, "theme"), Text: str(a, "text"), Link: s.localLink(loc, str(a, "link"))})
	}
	for _, f := range listOf(fm, "features") {
		p.Features = append(p.Features, Feature{Icon: str(f, "icon"), Title: str(f, "title"), Details: str(f, "details"), Link: s.localLink(loc, str(f, "link"))})
	}
	if p.Title == "" && p.Layout == "home" {
		p.Title = p.Hero.Name
	}
	if p.Title == "" {
		// The first # heading; known after rendering, but the sidebar and
		// <title> need it before.
		p.Title = firstHeading(body)
	}
	if prev, dup := s.pages[p.Route]; dup {
		return fmt.Errorf("%s and %s both map to %s", prev.File, file, p.Route)
	}
	s.pages[p.Route] = p
	s.byFile[file] = p
	s.ordered = append(s.ordered, p)
	return nil
}

// expandIncludes replaces <!--@include: path--> with the file's content
// (path relative to the including file), as VitePress does. The directive
// must be a line of its own outside code fences, so that pages can show it
// in code.
func (s *Site) expandIncludes(file, body string, depth int) (string, error) {
	const open, end = "<!--@include:", "-->"
	if !strings.Contains(body, open) {
		return body, nil
	}
	if depth > 4 {
		return "", fmt.Errorf("%s: includes nested too deeply", file)
	}
	var b strings.Builder
	inFence := false
	for _, line := range strings.SplitAfter(body, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "```") {
			inFence = !inFence
		}
		if inFence || !strings.HasPrefix(t, open) || !strings.HasSuffix(t, end) {
			b.WriteString(line)
			continue
		}
		target := path.Join(path.Dir(file), strings.TrimSpace(t[len(open):len(t)-len(end)]))
		data, err := fs.ReadFile(s.fsys, target)
		if err != nil {
			return "", fmt.Errorf("%s: include: %w", file, err)
		}
		_, inc, err := splitFrontmatter(string(data))
		if err != nil {
			return "", fmt.Errorf("%s: %w", target, err)
		}
		inc, err = s.expandIncludes(target, inc, depth+1)
		if err != nil {
			return "", err
		}
		b.WriteString(strings.TrimRight(inc, "\n"))
		if strings.HasSuffix(line, "\n") {
			b.WriteByte('\n')
		}
	}
	return b.String(), nil
}

func firstHeading(body string) string {
	inFence := false
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
		}
		if !inFence && strings.HasPrefix(line, "# ") {
			return strings.Trim(strings.TrimSpace(line[2:]), "`")
		}
	}
	return ""
}

// routeOf maps a content path to its route: en/guide/x.md -> /guide/x/,
// en/index.md -> /, ja/guide/index.md -> /ja/guide/.
func routeOf(loc *Locale, file string) string {
	rel := strings.TrimSuffix(strings.TrimPrefix(file, loc.Code+"/"), ".md")
	rel = strings.TrimSuffix(rel, "index")
	r := loc.Prefix + "/" + rel
	if !strings.HasSuffix(r, "/") {
		r += "/"
	}
	return r
}

// localLink prefixes a site-relative link with the locale prefix.
func (s *Site) localLink(loc *Locale, link string) string {
	if link == "" || strings.Contains(link, "://") || !strings.HasPrefix(link, "/") {
		return link
	}
	return normalizeRoute(loc.Prefix + link)
}

func normalizeRoute(r string) string {
	p, frag, _ := strings.Cut(r, "#")
	if path.Ext(p) == "" && !strings.HasSuffix(p, "/") {
		p += "/"
	}
	if frag != "" {
		return p + "#" + frag
	}
	return p
}

// Page returns the page at route, or nil.
func (s *Site) Page(route string) *Page {
	return s.pages[normalizeRoute(route)]
}

// Pages returns every page, ordered by route.
func (s *Site) Pages() []*Page { return s.ordered }

// LocaleOf returns the locale a route belongs to.
func (s *Site) LocaleOf(route string) *Locale {
	root := s.Config.Locales[0]
	for _, l := range s.Config.Locales[1:] {
		if route == l.Prefix || strings.HasPrefix(route, l.Prefix+"/") {
			return l
		}
	}
	return root
}

// Render returns the page's HTML and outline, rendering it once.
func (s *Site) Render(p *Page) (string, []Heading) {
	if !p.rendered {
		p.html, p.headings, p.sections = s.render(p, p.body)
		p.rendered = true
	}
	return p.html, p.headings
}

// resolveLink rewrites a link of page p. It reports whether the result
// leaves the site.
func (s *Site) resolveLink(p *Page, dest string, image bool) (string, bool) {
	if dest == "" || strings.HasPrefix(dest, "#") {
		return dest, false
	}
	if strings.Contains(dest, "://") || strings.HasPrefix(dest, "mailto:") {
		if s.Config.SiteURL != "" && strings.HasPrefix(dest, s.Config.SiteURL+"/") {
			return strings.TrimPrefix(dest, s.Config.SiteURL), false
		}
		return dest, true
	}
	target, frag, _ := strings.Cut(dest, "#")
	if frag != "" {
		frag = "#" + frag
	}
	if strings.HasPrefix(target, "/") {
		if image || path.Ext(target) != "" && path.Ext(target) != ".md" {
			return dest, false
		}
		return normalizeRoute(strings.TrimSuffix(target, ".md")) + frag, false
	}
	if p.Source != "" {
		src := s.Config.Sources[p.Source]
		rel := path.Clean(path.Join(path.Dir(p.SourcePath), target))
		if !image {
			if q := s.sourcePage(p.Locale, p.Source, rel); q != nil {
				return q.Route + frag, false
			}
		}
		if image {
			return src.Assets + rel, false
		}
		return src.Repo + "/blob/" + src.Ref + "/" + rel + frag, true
	}
	rel := path.Clean(path.Join(path.Dir(p.File), target))
	if q := s.byFile[rel]; q != nil {
		return q.Route + frag, false
	}
	if q := s.byFile[path.Join(rel, "index.md")]; q != nil {
		return q.Route + frag, false
	}
	return dest, false
}

// sourcePage finds the page synced from a source path, preferring the
// locale's own translation (docs/x.ja.md for ja) over the original.
func (s *Site) sourcePage(loc *Locale, name, rel string) *Page {
	base := strings.TrimSuffix(rel, ".md")
	for _, l := range s.Config.Locales {
		base = strings.TrimSuffix(base, "."+l.Code)
	}
	candidates := []string{base + "." + loc.Code + ".md", base + ".md", rel}
	for _, c := range candidates {
		if q := s.bySrc[loc.Code+":"+name+":"+c]; q != nil {
			return q
		}
	}
	return nil
}

// NavLink is a resolved navigation entry.
type NavLink struct {
	Text     string `json:"text"`
	Link     string `json:"link"`
	Active   bool   `json:"active"`
	External bool   `json:"external"`
}

// Nav returns the navigation bar entries for route.
func (s *Site) Nav(route string) []NavLink {
	loc := s.LocaleOf(route)
	var out []NavLink
	for _, n := range loc.Nav {
		ext := strings.Contains(n.Link, "://")
		link := n.Link
		if !ext {
			link = s.localLink(loc, n.Link)
		}
		match := n.Match
		if match == "" {
			match = n.Link
		}
		active := !ext && strings.HasPrefix(route, s.localLink(loc, match))
		out = append(out, NavLink{Text: n.Text, Link: link, Active: active, External: ext})
	}
	return out
}

// SidebarView is a sidebar group with resolved links.
type SidebarView struct {
	Text  string    `json:"text"`
	Items []NavLink `json:"items"`
}

// Sidebar returns the sidebar for route: the groups configured for the
// longest matching prefix in the route's locale.
func (s *Site) Sidebar(route string) []SidebarView {
	loc := s.LocaleOf(route)
	rel := strings.TrimPrefix(route, loc.Prefix)
	best := ""
	for prefix := range loc.Sidebar {
		if strings.HasPrefix(rel, prefix) && len(prefix) > len(best) {
			best = prefix
		}
	}
	if best == "" {
		return nil
	}
	var out []SidebarView
	for _, g := range loc.Sidebar[best] {
		v := SidebarView{Text: g.Text}
		for _, it := range g.Items {
			link := s.localLink(loc, it.Link)
			text := it.Text
			if text == "" {
				if p := s.Page(link); p != nil {
					text = p.Title
				}
			}
			v.Items = append(v.Items, NavLink{Text: text, Link: link, Active: link == route})
		}
		out = append(out, v)
	}
	return out
}

// PrevNext returns the sidebar neighbours of route (zero values if none).
func (s *Site) PrevNext(route string) (prev, next NavLink) {
	var flat []NavLink
	for _, g := range s.Sidebar(route) {
		flat = append(flat, g.Items...)
	}
	for i, l := range flat {
		if l.Link != route {
			continue
		}
		if i > 0 {
			prev = flat[i-1]
		}
		if i+1 < len(flat) {
			next = flat[i+1]
		}
		break
	}
	return prev, next
}

// Alternate is the same page in another locale.
type Alternate struct {
	Label   string `json:"label"`
	Lang    string `json:"lang"`
	Link    string `json:"link"`
	Current bool   `json:"current"`
	// Exists is false when the page has no translation and Link is the
	// locale's home page instead.
	Exists bool `json:"exists"`
}

// Alternates returns route in every locale, in the configured order.
func (s *Site) Alternates(route string) []Alternate {
	loc := s.LocaleOf(route)
	rel := strings.TrimPrefix(route, loc.Prefix)
	var out []Alternate
	for _, l := range s.Config.Locales {
		link := normalizeRoute(l.Prefix + rel)
		exists := s.Page(link) != nil
		if !exists {
			link = normalizeRoute(l.Prefix + "/")
		}
		out = append(out, Alternate{Label: l.Label, Lang: l.Lang, Link: link, Current: l == loc, Exists: exists})
	}
	return out
}

// EditLink returns the "edit this page" URL of a page: the source file in
// its own repository for a synced page, else the content file.
func (s *Site) EditLink(p *Page) (link string, synced bool) {
	if p.Source != "" {
		src := s.Config.Sources[p.Source]
		return src.Repo + "/blob/" + src.Ref + "/" + p.SourcePath, true
	}
	if s.Config.EditBase == "" {
		return "", false
	}
	return s.Config.EditBase + p.File, false
}

// Sitemap returns sitemap.xml with hreflang alternates.
func (s *Site) Sitemap() string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n")
	for _, p := range s.ordered {
		b.WriteString("  <url>\n    <loc>" + escapeHTML(s.Config.SiteURL+p.Route) + "</loc>\n")
		for _, a := range s.Alternates(p.Route) {
			if a.Exists {
				b.WriteString(`    <xhtml:link rel="alternate" hreflang="` + a.Lang + `" href="` + escapeHTML(s.Config.SiteURL+a.Link) + `"/>` + "\n")
			}
		}
		b.WriteString("  </url>\n")
	}
	b.WriteString("</urlset>\n")
	return b.String()
}
