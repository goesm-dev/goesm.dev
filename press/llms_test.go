package press

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestMarkdownPath(t *testing.T) {
	for route, want := range map[string]string{
		"/":                     "/index.md",
		"/ja/":                  "/ja.md",
		"/guide/":               "/guide.md",
		"/gosfc/guide/x/#intro": "/gosfc/guide/x.md",
	} {
		if got := MarkdownPath(route); got != want {
			t.Errorf("MarkdownPath(%q) = %q, want %q", route, got, want)
		}
	}
}

func llmsSite(t *testing.T) *Site {
	t.Helper()
	fsys := fstest.MapFS{
		"en/index.md": {Data: []byte("---\nlayout: home\nhero:\n  name: goesm\n  text: Go as ESM\n  tagline: No wasm.\n  actions:\n    - text: Start\n      link: /guide/\nfeatures:\n  - title: Fast\n    details: Plain JS.\n    link: /guide/\n---\n")},
		"en/guide/index.md": {Data: []byte("---\ndescription: What it is.\n---\n\n# Guide\n\n<!--@include: ./_part.md-->\n\n" +
			"::: tip Note\nSee [start](./start.md#install) and ![logo](/logo.png).\n:::\n\n" +
			"```go {2}\n// [not a link](./x.md)\n```\n\n<img src=\"/a.svg\">\n")},
		"en/guide/_part.md": {Data: []byte("<!-- Synced. -->\n\nIncluded text.\n")},
		"en/guide/start.md": {Data: []byte("---\ntitle: Start\n---\nText without a heading.\n")},
		"en/extra.md":       {Data: []byte("# Extra\n")},
	}
	cfg := &Config{
		SiteURL: "https://goesm.dev",
		Locales: []*Locale{{
			Code: "en", Lang: "en", Label: "English", Title: "goesm", Description: "Go to ESM.",
			UI:      UI{OtherPages: "Other pages", Containers: map[string]string{"tip": "TIP"}},
			Sidebar: map[string][]SidebarGroup{"/guide/": {{Text: "Guide", Items: []SidebarItem{{Link: "/guide/"}, {Text: "Start", Link: "/guide/start"}}}}},
		}},
	}
	s, err := Load(fsys, cfg)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestMarkdown(t *testing.T) {
	s := llmsSite(t)
	got := s.Markdown(s.Page("/guide/"))
	want := "# Guide\n\n> What it is.\n\nIncluded text.\n\n" +
		"> **Note**\n>\n> See [start](https://goesm.dev/guide/start.md#install) and ![logo](https://goesm.dev/logo.png).\n\n" +
		"```go\n// [not a link](./x.md)\n```\n\n<img src=\"https://goesm.dev/a.svg\">\n"
	if got != want {
		t.Errorf("Markdown:\n%s\nwant:\n%s", got, want)
	}
	if got := s.Markdown(s.Page("/guide/start/")); !strings.HasPrefix(got, "# Start\n\nText without") {
		t.Errorf("page without a heading:\n%s", got)
	}
	home := s.Markdown(s.Page("/"))
	for _, want := range []string{"# goesm: Go as ESM\n\n> No wasm.\n", "[Start](https://goesm.dev/guide.md)", "- [**Fast**](https://goesm.dev/guide.md): Plain JS.\n"} {
		if !strings.Contains(home, want) {
			t.Errorf("home has no %q:\n%s", want, home)
		}
	}
}

func TestLLMsTxt(t *testing.T) {
	s := llmsSite(t)
	loc := s.Config.Locales[0]
	got := s.LLMsTxt(loc, LLMsSection{Title: "Optional", Links: []LLMsLink{{Text: "GitHub", URL: "https://github.com/goesm-dev/goesm"}}})
	want := "# goesm\n\n> Go to ESM.\n\n[goesm](https://goesm.dev/index.md): No wasm.\n\n" +
		"## Guide\n\n- [Guide](https://goesm.dev/guide.md): What it is.\n- [Start](https://goesm.dev/guide/start.md)\n\n" +
		"## Other pages\n\n- [Extra](https://goesm.dev/extra.md)\n\n" +
		"## Optional\n\n- [GitHub](https://github.com/goesm-dev/goesm)\n\n"
	if got != want {
		t.Errorf("LLMsTxt:\n%s\nwant:\n%s", got, want)
	}
	full := s.LLMsFullTxt(loc)
	if !strings.Contains(full, "Source: https://goesm.dev/guide/start/\n\n# Start") || strings.Count(full, "\n---\n") != 4 {
		t.Errorf("LLMsFullTxt:\n%s", full)
	}
	if got := s.LLMsPath(loc, true); got != "/llms-full.txt" {
		t.Errorf("LLMsPath %s", got)
	}
}
