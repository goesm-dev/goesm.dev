package site

import "goesm.dev/press"

// The site's configuration: what .vitepress/config would hold in VitePress.

var guideSidebar = func(guide, intro, start, vue, status, perf, this, ref string) []press.SidebarGroup {
	return []press.SidebarGroup{
		{Text: guide, Items: []press.SidebarItem{
			{Text: intro, Link: "/guide/"},
			{Text: start, Link: "/guide/getting-started"},
			{Text: vue, Link: "/guide/vue"},
			{Text: status, Link: "/guide/status"},
			{Text: perf, Link: "/guide/performance"},
			{Text: this, Link: "/guide/this-site"},
		}},
		{Text: ref, Items: []press.SidebarItem{
			{Link: "/reference/architecture"},
			{Link: "/reference/example-output"},
			{Link: "/reference/use-cases"},
			{Link: "/reference/js-exports"},
			{Link: "/reference/js-imports"},
			{Link: "/reference/dom"},
			{Link: "/reference/concurrency"},
			{Link: "/reference/conformance"},
			{Link: "/reference/otelc"},
			{Link: "/reference/gopherjs-comparison"},
			{Link: "/reference/benchmark"},
			{Link: "/reference/compare"},
			{Link: "/reference/contributing"},
			{Link: "/reference/releasing"},
		}},
	}
}

var en = &press.Locale{
	Code:        "en",
	Prefix:      "",
	Lang:        "en",
	Label:       "English",
	Title:       "goesm",
	Description: "Compile Go packages into native ES modules: plain JavaScript, no WebAssembly.",
	Nav: []press.NavItem{
		{Text: "Guide", Link: "/guide/", Match: "/guide/"},
		{Text: "Reference", Link: "/reference/architecture", Match: "/reference/"},
		{Text: "Performance", Link: "/guide/performance"},
		{Text: "gosfc", Link: "/gosfc/", Root: true},
		{Text: "Releases", Link: "https://github.com/goesm-dev/goesm/releases"},
	},
	UI: press.UI{
		OnThisPage:    "On this page",
		Prev:          "Previous page",
		Next:          "Next page",
		EditPage:      "Edit this page on GitHub",
		SyncedFrom:    "This page is synced from",
		Search:        "Search",
		SearchHint:    "Search the docs",
		NoResults:     "No results for",
		Menu:          "Menu",
		ToggleTheme:   "Toggle dark mode",
		SkipToContent: "Skip to content",
		Language:      "Language",
		NotFound:      "Page not found",
		NotFoundText:  "There is no page at this address.",
		BackHome:      "Take me home",
		CopyCode:      "Copy code",
		Copied:        "Copied",
		Footer:        "Released under the BSD 3-Clause License. This site is built with goesm, gosfc, Astro and Vue.",
		SiteSource:    "Source of this site",
		ViewMarkdown:  "View as Markdown",
		OtherPages:    "Other pages",
		Containers: map[string]string{
			"tip": "TIP", "info": "INFO", "warning": "WARNING", "danger": "DANGER", "details": "Details",
			"note": "Note", "important": "Important", "caution": "Caution",
		},
	},
}

var ja = &press.Locale{
	Code:        "ja",
	Prefix:      "/ja",
	Lang:        "ja",
	Label:       "日本語",
	Title:       "goesm",
	Description: "Go のパッケージをネイティブな ES モジュールにコンパイルします。素の JavaScript で、WebAssembly は使いません。",
	Nav: []press.NavItem{
		{Text: "ガイド", Link: "/guide/", Match: "/guide/"},
		{Text: "リファレンス", Link: "/reference/architecture", Match: "/reference/"},
		{Text: "性能", Link: "/guide/performance"},
		{Text: "gosfc", Link: "/gosfc/ja/", Root: true},
		{Text: "リリース", Link: "https://github.com/goesm-dev/goesm/releases"},
	},
	UI: press.UI{
		OnThisPage:    "このページの内容",
		Prev:          "前のページ",
		Next:          "次のページ",
		EditPage:      "GitHub でこのページを編集",
		SyncedFrom:    "このページの元のファイル:",
		Search:        "検索",
		SearchHint:    "ドキュメントを検索",
		NoResults:     "見つかりませんでした:",
		Menu:          "メニュー",
		ToggleTheme:   "ダークモードの切り替え",
		SkipToContent: "本文へ移動",
		Language:      "言語",
		NotFound:      "ページが見つかりません",
		NotFoundText:  "このアドレスにページはありません。",
		BackHome:      "ホームへ戻る",
		CopyCode:      "コードをコピー",
		Copied:        "コピーしました",
		Footer:        "BSD 3-Clause License で公開しています。このサイトは goesm、gosfc、Astro、Vue で作られています。",
		SiteSource:    "このサイトのソース",
		ViewMarkdown:  "Markdown で表示",
		OtherPages:    "その他のページ",
		Containers: map[string]string{
			"tip": "ヒント", "info": "情報", "warning": "警告", "danger": "危険", "details": "詳細",
			"note": "注記", "important": "重要", "caution": "注意",
		},
	},
}

func init() {
	en.Sidebar = map[string][]press.SidebarGroup{}
	ja.Sidebar = map[string][]press.SidebarGroup{}
	enSide := guideSidebar("Guide", "What is goesm?", "Getting started", "Vue and Astro (gosfc)", "Status", "Performance", "How this site is built", "Reference")
	jaSide := guideSidebar("ガイド", "goesm とは", "はじめに", "Vue と Astro (gosfc)", "現状", "性能", "このサイトの作り方", "リファレンス")
	for _, prefix := range []string{"/guide/", "/reference/"} {
		en.Sidebar[prefix] = enSide
		ja.Sidebar[prefix] = jaSide
	}
}

var config = &press.Config{
	SiteURL:  "https://goesm.dev",
	EditBase: "https://github.com/goesm-dev/goesm.dev/edit/main/content/",
	Sources: map[string]press.Source{
		"goesm": {Repo: "https://github.com/goesm-dev/goesm", Ref: goesmRef, Assets: "/repo/goesm/"},
	},
	Social:  []press.Social{{Icon: "github", Link: "https://github.com/goesm-dev/goesm"}},
	Locales: []*press.Locale{en, ja},
}
