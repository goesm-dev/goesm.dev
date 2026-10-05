package site

import "goesm.dev/press"

// The gosfc site, at /gosfc/: the documentation of gosfc, which runs Go in
// Vue components and Astro pages. Its content is content/gosfc/.

var gosfcSidebar = func(guide, intro, start, block, astro, importing, bench, ref, arch, contrib string) []press.SidebarGroup {
	return []press.SidebarGroup{
		{Text: guide, Items: []press.SidebarItem{
			{Text: intro, Link: "/guide/"},
			{Text: start, Link: "/guide/getting-started"},
			{Text: block, Link: "/guide/go-block"},
			{Text: astro, Link: "/guide/astro"},
			{Text: importing, Link: "/guide/importing-go"},
			{Text: bench, Link: "/guide/benchmark"},
		}},
		{Text: ref, Items: []press.SidebarItem{
			{Text: arch, Link: "/reference/architecture"},
			{Text: contrib, Link: "/reference/contributing"},
		}},
	}
}

var gosfcEn = &press.Locale{
	Code:        "en",
	Prefix:      "",
	Lang:        "en",
	Label:       "English",
	Title:       "gosfc",
	Description: "Real Go in Vue components and .astro files: the <script setup> of a .vue file, and the frontmatter and <script> of an .astro file.",
	Nav: []press.NavItem{
		{Text: "Guide", Link: "/guide/", Match: "/guide/"},
		{Text: "Reference", Link: "/reference/architecture", Match: "/reference/"},
		{Text: "goesm", Link: "/", Root: true},
	},
	UI: withFooter(en.UI, "Released under the MIT License. This site is built with goesm, gosfc, Astro and Vue."),
}

var gosfcJa = &press.Locale{
	Code:        "ja",
	Prefix:      "/ja",
	Lang:        "ja",
	Label:       "日本語",
	Title:       "gosfc",
	Description: "Vue コンポーネントと .astro ファイルに本物の Go を書きます。.vue ファイルの <script setup>、.astro ファイルのフロントマターと <script> が Go になります。",
	Nav: []press.NavItem{
		{Text: "ガイド", Link: "/guide/", Match: "/guide/"},
		{Text: "リファレンス", Link: "/reference/architecture", Match: "/reference/"},
		{Text: "goesm", Link: "/ja/", Root: true},
	},
	UI: withFooter(ja.UI, "MIT License で公開しています。このサイトは goesm、gosfc、Astro、Vue で作られています。"),
}

func withFooter(ui press.UI, footer string) press.UI {
	ui.Footer = footer
	return ui
}

func init() {
	enSide := gosfcSidebar("Guide", "What is gosfc?", "Getting started", "Writing the Go block", "Go in .astro files", "Importing Go from JavaScript", "Benchmark", "Reference", "Architecture", "Contributing")
	jaSide := gosfcSidebar("ガイド", "gosfc とは", "はじめに", "Go ブロックの書き方", ".astro ファイルで Go を使う", "JavaScript から Go を import する", "ベンチマーク", "リファレンス", "アーキテクチャ", "コントリビューション")
	gosfcEn.Sidebar = map[string][]press.SidebarGroup{"/guide/": enSide, "/reference/": enSide}
	gosfcJa.Sidebar = map[string][]press.SidebarGroup{"/guide/": jaSide, "/reference/": jaSide}
}

var gosfcConfig = &press.Config{
	SiteURL:  "https://goesm.dev",
	Base:     "/gosfc",
	EditBase: "https://github.com/goesm-dev/goesm.dev/edit/main/content/gosfc/",
	Sources: map[string]press.Source{
		"gosfc": {Repo: "https://github.com/goesm-dev/gosfc", Ref: gosfcRef, Assets: "/repo/gosfc/"},
	},
	Social:  []press.Social{{Icon: "github", Link: "https://github.com/goesm-dev/gosfc"}},
	Locales: []*press.Locale{gosfcEn, gosfcJa},
}
