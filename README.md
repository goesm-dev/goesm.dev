# goesm.dev

[日本語](README.ja.md)

The website of [goesm](https://github.com/goesm-dev/goesm) and, at
`/gosfc/`, of [gosfc](https://github.com/goesm-dev/gosfc), built with goesm
itself: a VitePress-style documentation engine written in Go, compiled by goesm,
and rendered by Astro and Vue through [gosfc](https://github.com/goesm-dev/gosfc).
How it fits together is explained on the site, in
[How this site is built](content/en/guide/this-site.md).

## Layout

| Path | What it is |
| --- | --- |
| `press/` | the engine: Markdown (goldmark), containers, highlighting, sidebars, locales, search index, sitemap |
| `site/` | the configuration of both sites and the functions the pages call |
| `content/en`, `content/ja` | the goesm pages, in Markdown |
| `content/gosfc/en`, `content/gosfc/ja` | the gosfc pages, served under `/gosfc/` |
| `client/` | Go that runs in the browser (search, WebMCP tools, live demo, page enhancements) |
| `worker/` | the Worker in front of the static files: Markdown for clients that do not ask for HTML |
| `src/` | Astro pages and the Vue theme; the components' `<script setup>` is Go |
| `cmd/syncdocs` | copies the goesm and gosfc docs at the versions in `go.mod` |
| `third_party/gosfc` | gosfc as a git submodule, for its Vite and Astro packages |

## Develop

Tool versions are pinned in `mise.toml`.

```sh
git clone --recurse-submodules https://github.com/goesm-dev/goesm.dev
cd goesm.dev
mise install
pnpm install
pnpm dev          # http://localhost:4321
```

```sh
go test ./...     # the engine and the site, natively
pnpm build        # static site in dist/
node --test tests/  # checks of dist/
```

Pages under `content/*/reference/` and `content/*/_goesm`, `_gosfc` are
generated. To update them, bump goesm or gosfc in `go.mod` (and the submodule),
then run, with a clone of goesm next to this repository:

```sh
go run ./cmd/syncdocs -goesm ../goesm -gosfc third_party/gosfc
```

CI runs the same with `-check`.

## Deploy

The site runs on Cloudflare Workers: `cloudflare.config.ts` describes the
Worker and `wrangler.config.ts` points it at `dist/`. Cloudflare serves the
static files; the script in `worker/` runs first for page URLs and answers
clients that do not ask for HTML (curl, AI agents) with the page's Markdown.

```sh
pnpm build:cf     # dist/, then the Build Output in .cloudflare/output
pnpm preview:cf   # serve it locally as Workers would
pnpm run deploy   # build and upload (needs `cf auth login`)
```

CI deploys every merge to main with `mise run build` and `mise run deploy`,
then tags the commit `vYYYY.M.N` (N counts from 0 each month, in UTC). It
authenticates with the `CF_ID` and `CF_TOKEN` repository secrets.
