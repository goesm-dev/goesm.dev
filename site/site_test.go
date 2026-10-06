package site

import (
	"encoding/json"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
)

var internalLink = regexp.MustCompile(`(?:href|src|srcset)="(/[^"#]*)(#[^"]*)?"`)

// Every page renders, and every link to the site itself points at a page, a
// heading of that page or a file in public/.
func TestLinks(t *testing.T) {
	var routes []route
	if err := json.Unmarshal([]byte(RoutesJSON()), &routes); err != nil {
		t.Fatal(err)
	}
	if len(routes) < 20 {
		t.Fatalf("only %d routes", len(routes))
	}
	ids := map[string]map[string]bool{}
	pages := map[string]string{}
	for _, r := range routes {
		var html string
		if p := siteOf(r.Route).Page(r.Route); p.Layout == "home" {
			html = Home(r.Route).HTML
		} else {
			d := Doc(r.Route)
			html = d.HTML
			if d.Title == "" {
				t.Errorf("%s: no title", r.Route)
			}
		}
		pages[r.Route] = html
		ids[r.Route] = map[string]bool{}
		for _, m := range regexp.MustCompile(` id="([^"]+)"`).FindAllStringSubmatch(html, -1) {
			ids[r.Route][m[1]] = true
		}
	}
	for route, html := range pages {
		for _, m := range internalLink.FindAllStringSubmatch(html, -1) {
			target, frag := m[1], strings.TrimPrefix(m[2], "#")
			// Markdown percent-encodes non-ASCII fragments; browsers decode
			// them before matching an id.
			if f, err := url.PathUnescape(frag); err == nil {
				frag = f
			}
			if pageIDs, ok := ids[target]; ok {
				if frag != "" && !pageIDs[frag] {
					t.Errorf("%s: link to missing heading %s#%s", route, target, frag)
				}
				continue
			}
			if _, err := os.Stat("../public" + target); err != nil {
				t.Errorf("%s: broken link %s", route, target)
			}
		}
		if strings.Contains(html, "<!--@include") {
			t.Errorf("%s: unexpanded include", route)
		}
	}
}

func TestHead(t *testing.T) {
	var h Head
	json.Unmarshal([]byte(HeadJSON("/ja/guide/")), &h)
	if !h.Found || h.Lang != "ja" || h.Title != "goesm とは | goesm" || len(h.Alternates) != 2 || h.Search != "/search/ja.txt" {
		t.Errorf("%+v", h)
	}
	json.Unmarshal([]byte(HeadJSON("/nope/")), &h)
	if h.Found {
		t.Error("found /nope/")
	}
	h = Head{}
	json.Unmarshal([]byte(HeadJSON("/gosfc/ja/guide/")), &h)
	if !h.Found || h.Site != "gosfc" || h.Title != "gosfc とは | gosfc" || h.Search != "/gosfc/search/ja.txt" || h.Home != "/gosfc/ja/" ||
		len(h.Alternates) != 2 || h.Alternates[0].Href != "https://goesm.dev/gosfc/guide/" {
		t.Errorf("%+v", h)
	}
	h = Head{}
	json.Unmarshal([]byte(HeadJSON("/gosfc/ja/nope/")), &h)
	if h.Found || h.Lang != "ja" || h.Home != "/gosfc/ja/" {
		t.Errorf("gosfc 404: %+v", h)
	}
}

// The two sites link to each other, and the sitemap lists both.
func TestSites(t *testing.T) {
	var gosfcLink bool
	for _, l := range NavBar("/ja/guide/").Nav {
		gosfcLink = gosfcLink || l.Link == "/gosfc/ja/"
	}
	if !gosfcLink {
		t.Error("no link to /gosfc/ja/ in the goesm nav")
	}
	if bar := NavBar("/gosfc/guide/"); bar.Title != "gosfc" || bar.Home != "/gosfc/" || bar.Social[0].Link != "https://github.com/goesm-dev/gosfc" {
		t.Errorf("gosfc bar %+v", bar)
	}
	sm := Sitemap()
	for _, u := range []string{"https://goesm.dev/guide/", "https://goesm.dev/gosfc/guide/go-block/", "https://goesm.dev/gosfc/ja/reference/architecture/"} {
		if !strings.Contains(sm, "<loc>"+u+"</loc>") {
			t.Errorf("sitemap has no %s", u)
		}
	}
	if SearchIndex("/gosfc", "en") == SearchIndex("", "en") {
		t.Error("one search index for both sites")
	}
}

func TestLLMs(t *testing.T) {
	var paths []string
	if err := json.Unmarshal([]byte(LLMsPathsJSON()), &paths); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(paths, " "); got != "gosfc/llms gosfc/llms-full gosfc/ja/llms gosfc/ja/llms-full llms llms-full ja/llms ja/llms-full" {
		t.Errorf("llms paths: %s", got)
	}
	for _, p := range paths {
		if txt := LLMs("/" + p + ".txt"); !strings.HasPrefix(txt, "# ") {
			t.Errorf("/%s.txt: %.80q", p, txt)
		}
	}
	ja := LLMs("/gosfc/ja/llms.txt")
	for _, want := range []string{"# gosfc\n", "(https://goesm.dev/gosfc/ja/guide/go-block.md)", "[English](https://goesm.dev/gosfc/llms.txt)", "[goesm](https://goesm.dev/ja/llms.txt)"} {
		if !strings.Contains(ja, want) {
			t.Errorf("/gosfc/ja/llms.txt has no %q", want)
		}
	}
	var pages []struct{ Slug, Route string }
	if err := json.Unmarshal([]byte(MarkdownPagesJSON()), &pages); err != nil {
		t.Fatal(err)
	}
	for _, p := range pages {
		if md := PageMarkdown(p.Route); !strings.HasPrefix(md, "# ") {
			t.Errorf("%s: %.80q", p.Route, md)
		}
	}
}

func TestVersion(t *testing.T) {
	if v := NavBar("/ja/guide/").Version; v.Label != goesmVersion || v.Links[0].Link != "https://github.com/goesm-dev/goesm/releases/tag/"+goesmVersion || v.Heading != "バージョン" {
		t.Errorf("goesm version %+v", v)
	}
	if v := NavBar("/gosfc/guide/").Version; v.Label != gosfcVersion || v.Links[1].Link != "https://pkg.go.dev/github.com/goesm-dev/gosfc@"+gosfcVersion {
		t.Errorf("gosfc version %+v", v)
	}
	for v, want := range map[string]string{
		"v0.0.0-20261005171715-4f0689be9e96":       "4f0689b 2026-10-05",
		"v1.2.4-pre.0.20261005171715-4f0689be9e96": "4f0689b 2026-10-05",
		"v0.0.1-beta.3": "v0.0.1-beta.3 ",
		"v0.0.1-beta.3.0.20261005171715-4f0689be9e96": "4f0689b 2026-10-05",
	} {
		if c, d := pseudo(v); c+" "+d != want {
			t.Errorf("pseudo(%s) = %s %s, want %s", v, c, d, want)
		}
	}
}
