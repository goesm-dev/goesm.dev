---
description: goesm.dev is built with goesm. Its VitePress-style docs engine is a Go package, run by Astro and rendered by Vue through gosfc.
---

# How this site is built

goesm.dev is built with goesm. The documentation engine, the Markdown renderer,
the code highlighter and the search are Go packages. goesm compiles them into ES
modules, [gosfc](./vue.md) connects them to Astro and Vue, and Astro writes static
HTML. Its source is [goesm-dev/goesm.dev](https://github.com/goesm-dev/goesm.dev).

```text
content/**/*.md ──► press (Go) ──► goesm ──► Astro pages + Vue theme ──► static HTML
                                                │
                                                └─► islands in the browser (Go too)
```

## press: a VitePress-style engine in Go

`goesm.dev/press` does what VitePress does for a docs site:

- Markdown through [goldmark](https://github.com/yuin/goldmark): CommonMark,
  GitHub tables, strikethrough and autolinks
- `::: tip` / `info` / `warning` / `danger` / `details` containers and GitHub
  alerts (`> [!NOTE]`)
- highlighted code blocks with a language label, a copy button and line
  highlights (` ```go {2,4-5}`); the highlighter is `press/highlight`
- heading anchors compatible with GitHub's, and the "On this page" outline
- `<!--@include: ./file.md-->`, and links to `.md` files rewritten to routes
- the nav bar, sidebars, previous / next links and "Edit this page" links
- locales: `/` is English and `/ja/` is Japanese, with translated UI, a
  language menu that keeps the page, and `hreflang` links
- a search index per locale, and `sitemap.xml`

`goesm.dev/site` is the configuration (what `.vitepress/config` holds in
VitePress): locales, nav, sidebars. The pages themselves are Markdown files in
`content/en` and `content/ja`, embedded with `go:embed`.

::: tip Tested with go test
press is ordinary Go, so it is tested natively with `go test`. The same code then
runs under Node.js at build time, compiled by goesm.
:::

## Astro and Vue

The Astro pages call Go directly. `src/pages/[...slug].astro` asks the engine
for the list of pages:

```astro
---
import { RoutesJSON, $runtime as rt } from "go:goesm.dev/site";

export function getStaticPaths() {
  const routes = JSON.parse(rt.toJSString(RoutesJSON()));
  return routes.map((r) => ({ params: { slug: r.slug || undefined }, props: { route: r.route } }));
}
---
```

The theme is a set of Vue components whose `<script setup>` is Go. The doc page
component, for example, calls the engine and hands the result to its template:

```vue
<script setup lang="go">
import "goesm.dev/site"

type Props struct {
	Route string
}

doc := site.Doc(props.Route)
</script>
```

Astro renders these components to HTML at build time. The pages ship no
JavaScript for them: rendering Markdown, highlighting and building the sidebar
all happen once, during the build.

## Go in the browser

A few components are hydrated as Astro islands, and their Go code runs in the
browser:

| Island | Go package | What it does | Size (gzip) |
| --- | --- | --- | --- |
| Search | `client/searchbox`, `press/search`, `client/webmcp` | loads the locale's index, ranks results, builds snippets; the WebMCP tools | 9 KB |
| Live demo (home) | `client/demo` | the cart: Go `int` arithmetic on every click | 3 KB |
| Page enhancements | `client/enhance` | copy buttons, the outline following the scroll | 1.5 KB |
| Theme toggle | (its own Go block) | dark mode, remembered in `localStorage` | 1 KB |

Besides Vue itself, they share the goesm runtime and `syscall/js`: about 11 KB
gzip, loaded once.

The search is the same idea as VitePress's local search: press writes a
plain-text index of every section at build time, and the browser loads it on
the first search. The query code avoids `strings`, `sort` and `strconv` on
purpose: with goesm today, importing them also brings in Unicode tables and
reflection, which would make the island many times larger.

## For AI agents

The site is written for programs as well as for people:

- Every page has a Markdown version next to it: `/guide/getting-started/` is
  also `/guide/getting-started.md`, with includes expanded, containers as block
  quotes and absolute links. Each page links to it ("View as Markdown", and a
  `<link rel="alternate" type="text/markdown">`).
- Each site and language has an [llms.txt](https://llmstxt.org/) listing its
  pages, and an `llms-full.txt` with all of them: `/llms.txt`, `/ja/llms.txt`,
  `/gosfc/llms.txt` and `/gosfc/ja/llms.txt`.
- The Worker in front of the static files (`worker/index.ts`) looks at the
  `Accept` header. A browser asks for HTML and gets the page; `curl`, `fetch`
  and other clients that do not ask for HTML get the Markdown from the same
  URL. `curl https://goesm.dev/guide/` prints Markdown.
- In a browser with [WebMCP](https://github.com/webmachinelearning/webmcp),
  the page registers three read-only tools with `document.modelContext`, so
  an agent working in the page can use them instead of reading the screen:
  `search_docs` (the search dialog's index), `get_page` (a page as Markdown)
  and `list_pages` (the llms.txt). They are Go too, in the search island.

The Markdown and the llms.txt files are written by press at build time, from
the same pages as the HTML.

## No WebAssembly needed

TinyGo and WebAssembly were an option for heavy work in the browser. Nothing on
this site turned out to need it: the heavy work (Markdown, highlighting, the
index) is done at build time, and searching the index (about 100 KB of text per language) is fast in
plain JavaScript. The islands stay small and call the DOM directly.

## Docs synced from the repositories

The reference pages and the README sections are not copied by hand.
`cmd/syncdocs` reads them from the goesm and gosfc repositories at the versions
in `go.mod`, rewrites their links and images, and writes them under `content/`.
CI runs it with `-check`, so the site cannot drift from the versions it
documents.
