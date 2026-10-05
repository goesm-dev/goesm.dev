//go:build js

package searchbox

import (
	"syscall/js"

	"goesm.dev/client/webmcp"
)

// RegisterTools registers the documentation's WebMCP tools, for AI agents
// working in the page: search_docs (the search dialog's index), get_page
// (a page as Markdown) and list_pages (the locale's llms.txt). site is the
// site's name ("goesm"), llms the path of the locale's llms.txt.
func (b *Box) RegisterTools(site, llms string) {
	if !webmcp.Context().Truthy() {
		return
	}
	location := js.Global().Get("location")
	origin := location.Get("origin").String()
	webmcp.Register(webmcp.Tool{
		Name:  "search_docs",
		Title: "Search the " + site + " documentation",
		Description: "Searches the " + site + " documentation, in the language of this page, and returns the matching sections: " +
			"their headings, URL and a snippet. Use get_page to read a whole page as Markdown.",
		Schema: `{"type":"object","properties":{` +
			`"query":{"type":"string","description":"The words to search for."},` +
			`"limit":{"type":"integer","minimum":1,"maximum":20,"description":"The maximum number of results (default 5)."}},` +
			`"required":["query"]}`,
		Run: func(in js.Value) string {
			q := webmcp.String(in, "query")
			if q == "" {
				return "Pass the words to search for in query."
			}
			b.Load()
			if b.index == nil {
				return "The search index could not be loaded. Try again, or use list_pages."
			}
			limit := webmcp.Int(in, "limit", 5)
			if limit < 1 || limit > 20 {
				limit = 5
			}
			results := b.index.Search(q, limit)
			if len(results) == 0 {
				return "No results for “" + q + "”. list_pages lists every page."
			}
			out := ""
			for _, r := range results {
				title := ""
				for i, t := range r.Titles {
					if i > 0 {
						title += " › "
					}
					title += t
				}
				snippet := ""
				for _, p := range r.Snippet {
					snippet += p.Text
				}
				out += "## " + title + "\n\n" + origin + r.Link + "\n\n" + snippet + "\n\n"
			}
			return out
		},
	})
	webmcp.Register(webmcp.Tool{
		Name:  "get_page",
		Title: "Read a page of the " + site + " documentation",
		Description: "Returns a page of the " + site + " documentation as Markdown, with absolute links. " +
			"Without a path, returns the page that is open.",
		Schema: `{"type":"object","properties":{` +
			`"path":{"type":"string","description":"The page's path or URL on this site, such as /guide/getting-started/."}}}`,
		Run: func(in js.Value) string {
			p := webmcp.String(in, "path")
			if p == "" {
				p = location.Get("pathname").String()
			}
			if hasPrefix(p, origin) {
				p = p[len(origin):]
			}
			if p == "" || p[0] != '/' || hasPrefix(p, "//") {
				return "Pass a path on this site, such as /guide/, in path."
			}
			text := fetchText(markdownPath(p))
			if text == "" {
				return "There is no page at " + p + ". list_pages lists every page."
			}
			return text
		},
	})
	webmcp.Register(webmcp.Tool{
		Name:  "list_pages",
		Title: "List the pages of the " + site + " documentation",
		Description: "Returns the llms.txt of the " + site + " documentation in the language of this page: " +
			"every page with its URL and summary, grouped as in the sidebar.",
		Schema: `{"type":"object","properties":{}}`,
		Run: func(in js.Value) string {
			if text := fetchText(llms); text != "" {
				return text
			}
			return "llms.txt could not be loaded. Try again."
		},
	})
}

// markdownPath is the path of a page's Markdown version, like
// press.MarkdownPath: /guide/x/ -> /guide/x.md, / -> /index.md. A query, a
// fragment and a path that already ends in .md are kept apart.
func markdownPath(p string) string {
	for i := 0; i < len(p); i++ {
		if p[i] == '?' || p[i] == '#' {
			p = p[:i]
			break
		}
	}
	if len(p) > 3 && p[len(p)-3:] == ".md" {
		return p
	}
	for len(p) > 0 && p[len(p)-1] == '/' {
		p = p[:len(p)-1]
	}
	if p == "" {
		p = "/index"
	}
	return p + ".md"
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
