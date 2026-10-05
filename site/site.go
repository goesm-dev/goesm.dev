// Package site is goesm.dev: the press engine loaded with the configuration
// and content of its two sites, goesm at / and gosfc at /gosfc/. Its
// functions are what the Astro pages and the Vue theme call; each takes the
// route of the page being rendered, which also says which site it is on.
//
// Functions returning structs are called from Vue templates (gosfc converts
// the result into plain objects, with the json tags as keys). Functions
// called from Astro or TypeScript return JSON or text.
package site

import (
	"encoding/json"
	"io/fs"
	"strings"

	"goesm.dev/content"
	"goesm.dev/press"
)

var (
	goesmSite = load(content.FS, config)
	gosfcSite = load(subFS(content.FS, "gosfc"), gosfcConfig)
	// Sites with a longer Base first, so that siteOf finds the most
	// specific one.
	sites = []*press.Site{gosfcSite, goesmSite}
)

func load(fsys fs.FS, cfg *press.Config) *press.Site {
	site, err := press.Load(fsys, cfg)
	if err != nil {
		panic(err)
	}
	return site
}

func subFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

// siteOf returns the site a route belongs to.
func siteOf(r string) *press.Site {
	for _, s := range sites {
		if s.Owns(r) {
			return s
		}
	}
	return goesmSite
}

type route struct {
	Slug  string `json:"slug"`
	Route string `json:"route"`
}

// RoutesJSON lists every page as {slug, route}, where slug is the route
// without its slashes, for Astro's getStaticPaths.
func RoutesJSON() string {
	var out []route
	for _, s := range sites {
		for _, p := range s.Pages() {
			out = append(out, route{Slug: strings.Trim(p.Route, "/"), Route: p.Route})
		}
	}
	return mustJSON(out)
}

// Head is what <head> needs.
type Head struct {
	Site        string         `json:"site"` // "goesm" or "gosfc"
	Found       bool           `json:"found"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Lang        string         `json:"lang"`
	Locale      string         `json:"locale"`
	Layout      string         `json:"layout"`
	Canonical   string         `json:"canonical"`
	Alternates  []altLink      `json:"alternates"`
	UI          press.UI       `json:"ui"`
	Search      string         `json:"search"` // URL of the locale's search index
	Home        string         `json:"home"`
	Social      []press.Social `json:"social"`
}

type altLink struct {
	Lang string `json:"lang"`
	Href string `json:"href"`
}

// HeadJSON returns the Head of route as JSON. For a route without a page
// (the 404 page), found is false and the rest describes the route's locale.
func HeadJSON(r string) string {
	s := siteOf(r)
	loc := s.LocaleOf(r)
	h := Head{
		Site:   loc.Title,
		Lang:   loc.Lang,
		Locale: loc.Code,
		UI:     loc.UI,
		Search: s.Config.Base + "/search/" + loc.Code + ".txt",
		Home:   s.Home(loc),
		Social: s.Config.Social,
		Title:  loc.UI.NotFound + " | " + loc.Title,
		Layout: "page",
	}
	p := s.Page(r)
	if p == nil {
		return mustJSON(h)
	}
	h.Found = true
	h.Layout = p.Layout
	h.Description = p.Description
	if h.Description == "" {
		h.Description = loc.Description
	}
	if p.Layout == "home" {
		h.Title = loc.Title
		if p.Hero.Text != "" {
			h.Title += ": " + p.Hero.Text
		}
	} else {
		h.Title = p.Title + " | " + loc.Title
	}
	h.Canonical = s.Config.SiteURL + p.Route
	for _, a := range s.Alternates(p.Route) {
		if a.Exists {
			h.Alternates = append(h.Alternates, altLink{Lang: a.Lang, Href: s.Config.SiteURL + a.Link})
		}
	}
	return mustJSON(h)
}

// SearchIndex returns the search index of a locale ("en", "ja") of the site
// at base ("" or "/gosfc").
func SearchIndex(base, locale string) string {
	for _, s := range sites {
		if s.Config.Base == base {
			return s.SearchIndex(locale)
		}
	}
	panic("no site at " + base)
}

// Sitemap returns sitemap.xml, for both sites.
func Sitemap() string { return press.Sitemap(goesmSite, gosfcSite) }

// Bar is the navigation bar.
type Bar struct {
	Title   string            `json:"title"`
	Home    string            `json:"home"`
	Nav     []press.NavLink   `json:"nav"`
	Locales []press.Alternate `json:"locales"`
	Social  []press.Social    `json:"social"`
	UI      press.UI          `json:"ui"`
}

func NavBar(r string) Bar {
	s := siteOf(r)
	loc := s.LocaleOf(r)
	return Bar{
		Title:   loc.Title,
		Home:    s.Home(loc),
		Nav:     s.Nav(r),
		Locales: s.Alternates(r),
		Social:  s.Config.Social,
		UI:      loc.UI,
	}
}

func Sidebar(r string) []press.SidebarView { return siteOf(r).Sidebar(r) }

// DocView is a documentation page.
type DocView struct {
	Title   string          `json:"title"`
	HTML    string          `json:"html"`
	Outline []press.Heading `json:"outline"`
	Edit    string          `json:"edit"`
	// Synced is the source file for a page synced from another repository
	// ("goesm-dev/goesm: docs/otelc.md"); the edit link then points there.
	Synced string        `json:"synced"`
	Prev   press.NavLink `json:"prev"`
	Next   press.NavLink `json:"next"`
	UI     press.UI      `json:"ui"`
}

func Doc(r string) DocView {
	s := siteOf(r)
	p := s.Page(r)
	if p == nil {
		return DocView{}
	}
	html, outline := s.Render(p)
	edit, synced := s.EditLink(p)
	v := DocView{Title: p.Title, HTML: html, Outline: outline, Edit: edit, UI: p.Locale.UI}
	if synced {
		v.Synced = strings.TrimPrefix(s.Config.Sources[p.Source].Repo, "https://github.com/") + ": " + p.SourcePath
	}
	v.Prev, v.Next = s.PrevNext(p.Route)
	return v
}

// HomeView is the landing page.
type HomeView struct {
	Hero     press.Hero      `json:"hero"`
	Features []press.Feature `json:"features"`
	HTML     string          `json:"html"`
}

func Home(r string) HomeView {
	s := siteOf(r)
	p := s.Page(r)
	if p == nil {
		return HomeView{}
	}
	html, _ := s.Render(p)
	return HomeView{Hero: p.Hero, Features: p.Features, HTML: html}
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
