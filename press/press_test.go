package press

import (
	"strings"
	"testing"
	"testing/fstest"
)

func testSite(t *testing.T) *Site {
	t.Helper()
	fsys := fstest.MapFS{
		"en/index.md": {Data: []byte(`---
layout: home
hero:
  name: goesm
  text: "Go as ES modules"
  actions:
    - theme: brand
      text: Get started
      link: /guide/start
features:
  - icon: ⚡
    title: Fast
    details: 'No WebAssembly, it''s JS'
---
`)},
		"en/guide/start.md": {Data: []byte(`# Getting started

Read [the reference](../reference/arch.md#value-representation) or [go.dev](https://go.dev).

## Install ` + "`goesm`" + `

> [!NOTE]
> goesm is *experimental*.
> Second line.

::: tip Pin it
Use mise.
:::

## Install goesm

` + "```go {2}\npackage cart\nfunc Total() int { return 0 }\n```" + `

<!--@include: ./_part.md-->

Write ` + "`<!--@include: ./missing.md-->`" + ` on a line of its own.
`)},
		"en/guide/_part.md":    {Data: []byte("Included text.\n")},
		"en/reference/arch.md": {Data: []byte("---\nsource: goesm:ARCHITECTURE.md\n---\n# Architecture\n\nSee [conformance](docs/conformance.md), [testdata](testdata/x.go) and ![logo](docs/assets/goesm.png).\n\n## Value representation\n\n## 11. 実装済み / 未実装\n")},
		"en/reference/conf.md": {Data: []byte("---\nsource: goesm:docs/conformance.md\n---\n# Conformance\n\nBack to [the design](../ARCHITECTURE.md#value-representation).\n")},
		"ja/guide/start.md":    {Data: []byte("# はじめに\n\n## インストール\n\ngoesm を\nインストールします。\n")},
		"ja/reference/conf.md": {Data: []byte("---\nsource: goesm:docs/conformance.ja.md\n---\n# 適合性\n\n[設計](../ARCHITECTURE.ja.md)\n")},
	}
	ui := UI{Containers: map[string]string{"note": "Note", "tip": "Tip"}}
	cfg := &Config{
		SiteURL:  "https://goesm.dev",
		EditBase: "https://github.com/goesm-dev/goesm.dev/edit/main/content/",
		Sources:  map[string]Source{"goesm": {Repo: "https://github.com/goesm-dev/goesm", Ref: "v0.0.1-beta.0", Assets: "/repo/goesm/"}},
		Locales: []*Locale{
			{Code: "en", Prefix: "", Lang: "en", Label: "English", UI: ui,
				Nav:     []NavItem{{Text: "Guide", Link: "/guide/start", Match: "/guide/"}},
				Sidebar: map[string][]SidebarGroup{"/guide/": {{Text: "Guide", Items: []SidebarItem{{Link: "/guide/start"}, {Text: "Arch", Link: "/reference/arch"}}}}}},
			{Code: "ja", Prefix: "/ja", Lang: "ja", Label: "日本語", UI: ui},
		},
	}
	s, err := Load(fsys, cfg)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestRoutes(t *testing.T) {
	s := testSite(t)
	var routes []string
	for _, p := range s.Pages() {
		routes = append(routes, p.Route)
	}
	want := "/ /guide/start/ /ja/guide/start/ /ja/reference/conf/ /reference/arch/ /reference/conf/"
	if got := strings.Join(routes, " "); got != want {
		t.Errorf("routes = %s, want %s", got, want)
	}
	if s.Page("/guide/start") == nil || s.Page("/guide/start/") == nil {
		t.Error("route lookup without / with trailing slash")
	}
}

func TestHome(t *testing.T) {
	p := testSite(t).Page("/")
	if p.Layout != "home" || p.Title != "goesm" || p.Hero.Text != "Go as ES modules" {
		t.Errorf("%+v", p)
	}
	if len(p.Hero.Actions) != 1 || p.Hero.Actions[0].Link != "/guide/start/" {
		t.Errorf("actions %+v", p.Hero.Actions)
	}
	if len(p.Features) != 1 || p.Features[0].Details != "No WebAssembly, it's JS" || p.Features[0].Icon != "⚡" {
		t.Errorf("features %+v", p.Features)
	}
}

func TestRender(t *testing.T) {
	s := testSite(t)
	p := s.Page("/guide/start/")
	html, heads := s.Render(p)
	for _, want := range []string{
		`<h1 id="getting-started" tabindex="-1">Getting started <a class="header-anchor" href="#getting-started"`,
		`<a href="/reference/arch/#value-representation">the reference</a>`,
		`<a href="https://go.dev" target="_blank" rel="noreferrer">go.dev</a>`,
		`<h2 id="install-goesm" tabindex="-1">Install <code>goesm</code>`,
		`<h2 id="install-goesm-1"`,
		`<div class="custom-block note"><p class="custom-block-title">Note</p>`,
		`goesm is <em>experimental</em>.`,
		`<div class="custom-block tip"><p class="custom-block-title">Pin it</p>`,
		`<p>Use mise.</p>`,
		`<span class="line highlighted"><span class="hl-kw">func</span>`,
		`<span class="lang">go</span>`,
		`<p>Included text.</p>`,
		`<code>&lt;!--@include: ./missing.md--&gt;</code>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in\n%s", want, html)
		}
	}
	if strings.Contains(html, "[!NOTE]") {
		t.Error("alert marker left in output")
	}
	if len(heads) != 2 || heads[0].ID != "install-goesm" || heads[0].Text != "Install goesm" {
		t.Errorf("headings %+v", heads)
	}
}

func TestJapaneseLineBreaks(t *testing.T) {
	s := testSite(t)
	html, _ := s.Render(s.Page("/ja/guide/start/"))
	if !strings.Contains(html, "<p>goesm をインストールします。</p>") {
		t.Errorf("a line break between Japanese characters became a space:\n%s", html)
	}
}

func TestSyncedLinks(t *testing.T) {
	s := testSite(t)
	html, heads := s.Render(s.Page("/reference/arch/"))
	for _, want := range []string{
		`<a href="/reference/conf/">conformance</a>`,
		`<a href="https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.0/testdata/x.go" target="_blank" rel="noreferrer">testdata</a>`,
		`<img src="/repo/goesm/docs/assets/goesm.png" alt="logo">`,
		`id="11-実装済み--未実装"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in\n%s", want, html)
		}
	}
	if len(heads) != 2 {
		t.Errorf("headings %+v", heads)
	}
	html, _ = s.Render(s.Page("/reference/conf/"))
	if !strings.Contains(html, `<a href="/reference/arch/#value-representation">`) {
		t.Errorf("link back to the design: %s", html)
	}
	// The ja page links to the ja translation, which is missing here: the
	// link falls back to GitHub.
	html, _ = s.Render(s.Page("/ja/reference/conf/"))
	if !strings.Contains(html, `href="https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.0/ARCHITECTURE.ja.md"`) {
		t.Errorf("ja link: %s", html)
	}
	link, synced := s.EditLink(s.Page("/reference/arch/"))
	if !synced || link != "https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.0/ARCHITECTURE.md" {
		t.Errorf("edit link %s %v", link, synced)
	}
}

func TestNavigation(t *testing.T) {
	s := testSite(t)
	nav := s.Nav("/guide/start/")
	if len(nav) != 1 || !nav[0].Active || nav[0].Link != "/guide/start/" {
		t.Errorf("nav %+v", nav)
	}
	sb := s.Sidebar("/guide/start/")
	if len(sb) != 1 || sb[0].Items[0].Text != "Getting started" || !sb[0].Items[0].Active || sb[0].Items[1].Text != "Arch" {
		t.Errorf("sidebar %+v", sb)
	}
	prev, next := s.PrevNext("/guide/start/")
	if prev.Link != "" || next.Link != "/reference/arch/" {
		t.Errorf("prev %+v next %+v", prev, next)
	}
	alts := s.Alternates("/guide/start/")
	if len(alts) != 2 || !alts[0].Current || alts[1].Link != "/ja/guide/start/" || !alts[1].Exists {
		t.Errorf("alternates %+v", alts)
	}
	alts = s.Alternates("/reference/arch/")
	if alts[1].Exists || alts[1].Link != "/ja/" {
		t.Errorf("missing translation %+v", alts)
	}
	if s.LocaleOf("/ja/guide/start/").Code != "ja" || s.LocaleOf("/japan/").Code != "en" {
		t.Error("LocaleOf")
	}
}

func TestSearchIndexAndSitemap(t *testing.T) {
	s := testSite(t)
	idx := s.SearchIndex("en")
	if !strings.Contains(idx, "/guide/start/#install-goesm\x1fGetting started\x1fGetting started\x1eInstall goesm\x1f") {
		t.Errorf("index:\n%q", idx)
	}
	if strings.Contains(idx, "はじめに") {
		t.Error("ja page in en index")
	}
	sm := s.Sitemap()
	if !strings.Contains(sm, `<loc>https://goesm.dev/guide/start/</loc>`) || !strings.Contains(sm, `hreflang="ja" href="https://goesm.dev/ja/guide/start/"`) {
		t.Errorf("sitemap:\n%s", sm)
	}
}

func TestFrontmatter(t *testing.T) {
	fm, body, err := splitFrontmatter("---\na: 1\nb:\n  c: \"x: y\" # comment\nl:\n- one\n- 'two'\n---\nbody\n")
	if err != nil {
		t.Fatal(err)
	}
	if body != "body\n" || str(fm, "a") != "1" || str(mapOf(fm, "b"), "c") != "x: y" {
		t.Errorf("%v %q", fm, body)
	}
	if l, _ := fm["l"].([]any); len(l) != 2 || l[1] != "two" {
		t.Errorf("list %v", fm["l"])
	}
}

func TestBase(t *testing.T) {
	fsys := fstest.MapFS{
		"en/index.md":       {Data: []byte("---\nlayout: home\nhero:\n  name: gosfc\n  actions:\n    - text: Start\n      link: /guide/\n---\n")},
		"en/guide/index.md": {Data: []byte("# What is gosfc?\n\nSee [the goesm guide](/guide/).\n")},
		"ja/guide/index.md": {Data: []byte("# gosfc とは\n")},
	}
	cfg := &Config{
		SiteURL: "https://goesm.dev",
		Base:    "/gosfc",
		Locales: []*Locale{
			{Code: "en", Prefix: "", Lang: "en", Label: "English",
				Nav:     []NavItem{{Text: "Guide", Link: "/guide/"}, {Text: "goesm", Link: "/", Root: true}},
				Sidebar: map[string][]SidebarGroup{"/guide/": {{Text: "Guide", Items: []SidebarItem{{Link: "/guide/"}}}}}},
			{Code: "ja", Prefix: "/ja", Lang: "ja", Label: "日本語"},
		},
	}
	s, err := Load(fsys, cfg)
	if err != nil {
		t.Fatal(err)
	}
	var routes []string
	for _, p := range s.Pages() {
		routes = append(routes, p.Route)
	}
	if got, want := strings.Join(routes, " "), "/gosfc/ /gosfc/guide/ /gosfc/ja/guide/"; got != want {
		t.Errorf("routes %q, want %q", got, want)
	}
	if !s.Owns("/gosfc/ja/") || s.Owns("/gosfcx/") || s.Owns("/guide/") {
		t.Error("Owns")
	}
	if l := s.LocaleOf("/gosfc/ja/guide/"); l.Code != "ja" {
		t.Errorf("LocaleOf: %s", l.Code)
	}
	if got := s.Page("/gosfc/").Hero.Actions[0].Link; got != "/gosfc/guide/" {
		t.Errorf("hero link %s", got)
	}
	nav := s.Nav("/gosfc/guide/")
	if nav[0].Link != "/gosfc/guide/" || !nav[0].Active || nav[1].Link != "/" || nav[1].Active {
		t.Errorf("nav %+v", nav)
	}
	if sb := s.Sidebar("/gosfc/guide/"); len(sb) != 1 || sb[0].Items[0].Text != "What is gosfc?" || !sb[0].Items[0].Active {
		t.Errorf("sidebar %+v", sb)
	}
	alt := s.Alternates("/gosfc/guide/")
	if alt[1].Link != "/gosfc/ja/guide/" || !alt[1].Exists {
		t.Errorf("alternates %+v", alt)
	}
	if alt := s.Alternates("/gosfc/"); alt[1].Link != "/gosfc/ja/" || alt[1].Exists {
		t.Errorf("alternates of home %+v", alt)
	}
	// Markdown links starting with / are paths from the origin.
	if html, _ := s.Render(s.Page("/gosfc/guide/")); !strings.Contains(html, `href="/guide/"`) {
		t.Errorf("origin link rewritten: %s", html)
	}
	if !strings.Contains(Sitemap(s), "<loc>https://goesm.dev/gosfc/guide/</loc>") {
		t.Error("sitemap")
	}
}

// A line break is a space in English, also before a link, and between
// Japanese and Latin text; it is nothing between Japanese characters.
func TestLineBreaks(t *testing.T) {
	fsys := fstest.MapFS{
		"en/a.md": {Data: []byte("# A\n\nIt renders pages.\n[The architecture](./a.md) has `x`\n`y`.\n")},
		"ja/a.md": {Data: []byte("# A\n\ngoesm を\nインストールします。Node.js と同時に\nGo も使います。\n")},
	}
	s, err := Load(fsys, &Config{Locales: []*Locale{{Code: "en", Lang: "en"}, {Code: "ja", Prefix: "/ja", Lang: "ja"}}})
	if err != nil {
		t.Fatal(err)
	}
	if html, _ := s.Render(s.Page("/a/")); !strings.Contains(html, "pages.\n<a") || !strings.Contains(html, "</code>\n<code>") {
		t.Errorf("en: %s", html)
	}
	if html, _ := s.Render(s.Page("/ja/a/")); !strings.Contains(html, "goesm をインストール") || !strings.Contains(html, "同時に\nGo") {
		t.Errorf("ja: %s", html)
	}
}
