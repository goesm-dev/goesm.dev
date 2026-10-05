// Checks of the built site (run `pnpm build` first): every page the engine
// lists is there, on both sites (goesm and gosfc) and in both languages,
// with its head and the Go islands.
import assert from "node:assert/strict";
import { existsSync, readdirSync, readFileSync } from "node:fs";
import { test } from "node:test";

const dist = new URL("../dist/", import.meta.url);
const read = (p) => readFileSync(new URL(p, dist), "utf8");
const page = (route) => read(route.slice(1) + "index.html");

const sitemap = existsSync(new URL("sitemap.xml", dist)) ? read("sitemap.xml") : "";
const routes = [...sitemap.matchAll(/<loc>https:\/\/goesm\.dev([^<]*)<\/loc>/g)].map((m) => m[1]);

// The Japanese route of an English one: /ja/x/ on goesm, /gosfc/ja/x/ on gosfc.
const japanese = (r) => (r.startsWith("/gosfc/") ? "/gosfc/ja" + r.slice("/gosfc".length) : "/ja" + r);
const isJapanese = (r) => r.startsWith("/ja/") || r.startsWith("/gosfc/ja/");

test("the site is built", () => {
  assert.ok(routes.length > 20, "run pnpm build first");
});

test("every page has a title, a canonical link and its alternates", () => {
  for (const r of routes) {
    const html = page(r);
    assert.match(html, /<title>[^<]+<\/title>/, r);
    assert.ok(html.includes(`<link rel="canonical" href="https://goesm.dev${r}">`), `${r}: canonical`);
    assert.match(html, new RegExp(`<html lang="${isJapanese(r) ? "ja" : "en"}"`), r);
    assert.ok(!html.includes("<!--@include"), `${r}: include left unexpanded`);
  }
});

test("pages exist in both languages", () => {
  const en = routes.filter((r) => !isJapanese(r));
  for (const r of en) assert.ok(routes.includes(japanese(r)), `no Japanese page for ${r}`);
  assert.match(page("/ja/guide/"), /hreflang="en" href="https:\/\/goesm\.dev\/guide\/"/);
  assert.match(page("/gosfc/ja/guide/"), /hreflang="en" href="https:\/\/goesm\.dev\/gosfc\/guide\/"/);
});

test("the gosfc site is its own site", () => {
  assert.ok(routes.filter((r) => r.startsWith("/gosfc/")).length >= 16, "gosfc pages missing from the sitemap");
  const html = page("/gosfc/guide/go-block/");
  assert.match(html, /<title>Writing the Go block \| gosfc<\/title>/);
  assert.match(html, /<span class="brand-name" translate="no">gosfc<\/span>/);
  assert.match(html, /href="https:\/\/github\.com\/goesm-dev\/gosfc"/);
  assert.match(page("/guide/"), /<a href="\/gosfc\/"[^>]*>gosfc</);
  // The navbar's GitHub icon is the tool's repository; the footer links the site's own.
  assert.match(page("/guide/"), /class="social"[^>]*href="https:\/\/github\.com\/goesm-dev\/goesm"|href="https:\/\/github\.com\/goesm-dev\/goesm"[^>]*class="social"/);
  assert.match(html, /class="social"[^>]*href="https:\/\/github\.com\/goesm-dev\/gosfc"|href="https:\/\/github\.com\/goesm-dev\/gosfc"[^>]*class="social"/);
  assert.match(html, /<footer class="site-footer">[\s\S]*<a href="https:\/\/github\.com\/goesm-dev\/goesm\.dev">Source of this site<\/a>/);
  assert.match(page("/gosfc/ja/"), /<a href="https:\/\/github\.com\/goesm-dev\/goesm\.dev">このサイトのソース<\/a>/);
});

test("each locale has a 404.html where Cloudflare looks for it", () => {
  assert.match(read("404.html"), /<html lang="en"/);
  assert.match(read("ja/404.html"), /<html lang="ja"/);
  assert.match(read("gosfc/404.html"), /<html lang="en"[\s\S]*href="\/gosfc\/"/);
  assert.match(read("gosfc/ja/404.html"), /<html lang="ja"[\s\S]*href="\/gosfc\/ja\/"/);
});

test("every page has a Markdown version", () => {
  const mdPath = (r) => (r === "/" ? "index.md" : r.slice(1, -1) + ".md");
  for (const r of routes) {
    const md = read(mdPath(r));
    assert.match(md, /^# \S/, `${r}: no title`);
    assert.ok(page(r).includes(`<link rel="alternate" type="text/markdown" href="/${mdPath(r)}">`), `${r}: alternate link`);
    const prose = md.replace(/^(```|~~~)[^]*?^\1/gm, "");
    assert.doesNotMatch(prose, /^\s*<!--|^:::|\]\(\.{0,2}\/(?!\/)/m, `${r}: comment, container or relative link left`);
  }
  assert.match(read("guide/getting-started.md"), /\]\(https:\/\/goesm\.dev\/guide\/[^)]*\.md[)#]/);
});

test("each site and language has an llms.txt", () => {
  for (const [dir, title] of [["", "goesm"], ["ja/", "goesm"], ["gosfc/", "gosfc"], ["gosfc/ja/", "gosfc"]]) {
    const llms = read(dir + "llms.txt");
    assert.match(llms, new RegExp(`^# ${title}\n\n> \\S`), dir);
    const links = [...llms.matchAll(/\]\((https:\/\/goesm\.dev\/[^)]*\.md)\)/g)].map((m) => m[1]);
    assert.ok(links.length >= 8, `${dir}llms.txt: ${links.length} pages`);
    for (const l of links) assert.ok(existsSync(new URL(l.slice("https://goesm.dev/".length), dist)), l);
    assert.ok(read(dir + "llms-full.txt").length > 20_000, `${dir}llms-full.txt`);
  }
});

test("every page has the favicon, touch icon and manifest", () => {
  for (const r of ["/", "/ja/guide/", "/gosfc/", "/gosfc/ja/guide/go-block/"]) {
    const html = page(r);
    assert.ok(html.includes('<link rel="icon" href="/favicon.ico" sizes="16x16 32x32 48x48">'), r);
    assert.ok(html.includes('<link rel="apple-touch-icon" href="/apple-touch-icon.png">'), r);
    assert.ok(html.includes('<link rel="manifest" href="/site.webmanifest">'), r);
  }
  const manifest = JSON.parse(read("site.webmanifest"));
  for (const icon of manifest.icons) assert.ok(existsSync(new URL(icon.src.slice(1), dist)), icon.src);
  for (const f of ["favicon.ico", "apple-touch-icon.png"]) assert.ok(existsSync(new URL(f, dist)), f);
});

test("doc pages are rendered at build time", () => {
  const html = page("/guide/getting-started/");
  assert.match(html, /<h2 id="install" tabindex="-1">Install/);
  assert.match(html, /class="hl-kw"/); // highlighted code
  assert.match(html, /class="sidebar/);
  assert.match(html, /class="outline/);
});

test("search indexes cover both languages", () => {
  for (const [index, min] of [["search/en", 50], ["search/ja", 50], ["gosfc/search/en", 20], ["gosfc/search/ja", 20]]) {
    const lines = read(`${index}.txt`).trim().split("\n");
    assert.ok(lines.length > min, `${index}: ${lines.length} records`);
    for (const line of lines) assert.equal(line.split("\x1f").length, 4, line);
  }
});

test("the first paint needs no request after the HTML", () => {
  // Inlined CSS: on a slow connection each extra request before the first
  // paint costs a full round trip.
  for (const r of ["/", "/ja/guide/", "/reference/architecture/", "/gosfc/", "/gosfc/ja/guide/go-block/"]) {
    const html = page(r);
    assert.doesNotMatch(html, /<link rel="stylesheet"/, r);
    assert.match(html, /<style>/, r);
    assert.doesNotMatch(html.slice(0, html.indexOf("</head>")), /<script[^>]+src=/, `${r}: blocking script`);
  }
});

test("the Go islands are small", () => {
  const assets = readdirSync(new URL("_astro/", dist));
  for (const name of ["Search", "LiveDemo", "Enhance", "ThemeToggle"]) {
    const file = assets.find((f) => f.startsWith(name + ".") && f.endsWith(".js"));
    assert.ok(file, `${name} island missing`);
    const size = readFileSync(new URL("_astro/" + file, dist)).length;
    assert.ok(size < 40_000, `${name}: ${size} bytes`);
  }
});
