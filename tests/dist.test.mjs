// Checks of the built site (run `pnpm build` first): every page the engine
// lists is there, in both languages, with its head and the Go islands.
import assert from "node:assert/strict";
import { existsSync, readdirSync, readFileSync } from "node:fs";
import { test } from "node:test";

const dist = new URL("../dist/", import.meta.url);
const read = (p) => readFileSync(new URL(p, dist), "utf8");
const page = (route) => read(route.slice(1) + "index.html");

const sitemap = existsSync(new URL("sitemap.xml", dist)) ? read("sitemap.xml") : "";
const routes = [...sitemap.matchAll(/<loc>https:\/\/goesm\.dev([^<]*)<\/loc>/g)].map((m) => m[1]);

test("the site is built", () => {
  assert.ok(routes.length > 20, "run pnpm build first");
});

test("every page has a title, a canonical link and its alternates", () => {
  for (const r of routes) {
    const html = page(r);
    assert.match(html, /<title>[^<]+<\/title>/, r);
    assert.ok(html.includes(`<link rel="canonical" href="https://goesm.dev${r}">`), `${r}: canonical`);
    assert.match(html, new RegExp(`<html lang="${r.startsWith("/ja/") ? "ja" : "en"}"`), r);
    assert.ok(!html.includes("<!--@include"), `${r}: include left unexpanded`);
  }
});

test("pages exist in both languages", () => {
  const en = routes.filter((r) => !r.startsWith("/ja/"));
  for (const r of en) assert.ok(routes.includes("/ja" + r), `no Japanese page for ${r}`);
  assert.match(page("/ja/guide/"), /hreflang="en" href="https:\/\/goesm\.dev\/guide\/"/);
});

test("doc pages are rendered at build time", () => {
  const html = page("/guide/getting-started/");
  assert.match(html, /<h2 id="install" tabindex="-1">Install/);
  assert.match(html, /class="hl-kw"/); // highlighted code
  assert.match(html, /class="sidebar/);
  assert.match(html, /class="outline/);
});

test("search indexes cover both languages", () => {
  for (const l of ["en", "ja"]) {
    const lines = read(`search/${l}.txt`).trim().split("\n");
    assert.ok(lines.length > 50, `${l}: ${lines.length} records`);
    for (const line of lines) assert.equal(line.split("\x1f").length, 4, line);
  }
});

test("the first paint needs no request after the HTML", () => {
  // Inlined CSS: on a slow connection each extra request before the first
  // paint costs a full round trip.
  for (const r of ["/", "/ja/guide/", "/reference/architecture/"]) {
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
