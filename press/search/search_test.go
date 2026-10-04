package search

import (
	"strings"
	"testing"
)

const data = "/guide/#install\x1fGuide\x1fGuide\x1eInstall\x1fRun go get -tool to add goesm as a tool.\n" +
	"/guide/#usage\x1fGuide\x1fGuide\x1eUsage\x1fgoesm emit-ts writes TypeScript; install nothing else.\n" +
	"/ja/guide/#インストール\x1fガイド\x1fガイド\x1eインストール\x1fgoesm を tool として追加します。\n"

func TestSearch(t *testing.T) {
	idx := Parse(data)
	if idx.Len() != 3 {
		t.Fatalf("len %d", idx.Len())
	}
	r := idx.Search("install", 10)
	if len(r) != 2 || r[0].Link != "/guide/#install" {
		t.Fatalf("%+v", r)
	}
	r = idx.Search("goesm typescript", 10)
	if len(r) != 1 || r[0].Link != "/guide/#usage" {
		t.Fatalf("%+v", r)
	}
	var marked []string
	for _, p := range r[0].Snippet {
		if p.Mark {
			marked = append(marked, p.Text)
		}
	}
	if strings.Join(marked, ",") != "goesm,TypeScript" {
		t.Errorf("marks %+v", r[0].Snippet)
	}
	r = idx.Search("追加", 10)
	if len(r) != 1 || r[0].Titles[1] != "インストール" {
		t.Fatalf("%+v", r)
	}
	if idx.Search("  ", 10) != nil {
		t.Error("empty query")
	}
}
