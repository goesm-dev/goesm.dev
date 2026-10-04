package site

import (
	"encoding/json"
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
		if p := s.Page(r.Route); p.Layout == "home" {
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
}
