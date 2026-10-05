package site

import (
	"strings"

	"goesm.dev/press"
)

// Every page as Markdown, and llms.txt / llms-full.txt for each site and
// locale: /llms.txt, /ja/llms.txt, /gosfc/llms.txt, /gosfc/ja/llms.txt.

type mdPage struct {
	Slug  string `json:"slug"` // the Markdown path without "/" and ".md"
	Route string `json:"route"`
}

// MarkdownPagesJSON lists the Markdown version of every page, for Astro's
// getStaticPaths.
func MarkdownPagesJSON() string {
	var out []mdPage
	for _, s := range sites {
		for _, p := range s.Pages() {
			slug := strings.TrimSuffix(strings.TrimPrefix(press.MarkdownPath(p.Route), "/"), ".md")
			out = append(out, mdPage{Slug: slug, Route: p.Route})
		}
	}
	return mustJSON(out)
}

// PageMarkdown returns the Markdown version of the page at route.
func PageMarkdown(r string) string {
	s := siteOf(r)
	p := s.Page(r)
	if p == nil {
		return ""
	}
	return s.Markdown(p)
}

// LLMsPathsJSON lists the paths of every llms.txt and llms-full.txt,
// without "/" and ".txt" ("llms", "ja/llms-full", "gosfc/llms").
func LLMsPathsJSON() string {
	var out []string
	for _, s := range sites {
		for _, loc := range s.Config.Locales {
			for _, full := range []bool{false, true} {
				out = append(out, strings.TrimSuffix(strings.TrimPrefix(s.LLMsPath(loc, full), "/"), ".txt"))
			}
		}
	}
	return mustJSON(out)
}

// LLMs returns the llms.txt or llms-full.txt at path ("/ja/llms.txt").
func LLMs(path string) string {
	for _, s := range sites {
		for _, loc := range s.Config.Locales {
			switch path {
			case s.LLMsPath(loc, false):
				return s.LLMsTxt(loc, llmsMore(s, loc))
			case s.LLMsPath(loc, true):
				return s.LLMsFullTxt(loc)
			}
		}
	}
	return ""
}

// llmsMore is the last section of an llms.txt: the full text, the other
// language and the other site.
func llmsMore(s *press.Site, loc *press.Locale) press.LLMsSection {
	ja := loc.Code == "ja"
	url := func(p string) string { return s.Config.SiteURL + p }
	sec := press.LLMsSection{Title: pick(ja, "その他", "Optional")}
	sec.Links = append(sec.Links, press.LLMsLink{
		Text:  "llms-full.txt",
		URL:   url(s.LLMsPath(loc, true)),
		Notes: pick(ja, "全ページの Markdown をひとつにしたもの", "every page above in one Markdown file"),
	})
	for _, l := range s.Config.Locales {
		if l != loc {
			sec.Links = append(sec.Links, press.LLMsLink{Text: l.Label, URL: url(s.LLMsPath(l, false)), Notes: pick(ja, "英語版のドキュメント", "the documentation in Japanese")})
		}
	}
	for _, o := range sites {
		if o == s {
			continue
		}
		ol := o.Config.Locales[0]
		if ja {
			ol = o.Config.Locales[len(o.Config.Locales)-1]
		}
		sec.Links = append(sec.Links, press.LLMsLink{Text: ol.Title, URL: url(o.LLMsPath(ol, false)), Notes: ol.Description})
	}
	sec.Links = append(sec.Links, press.LLMsLink{Text: "GitHub", URL: s.Config.Social[0].Link, Notes: pick(ja, "ソースコード", "the source code")})
	return sec
}

func pick(ja bool, jaText, enText string) string {
	if ja {
		return jaText
	}
	return enText
}
