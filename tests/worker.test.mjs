// The Worker's content negotiation, against a fake asset server: browsers
// get HTML, curl and other clients that do not ask for HTML get Markdown.
import assert from "node:assert/strict";
import { test } from "node:test";
import { handle, markdownPath, wantsMarkdown } from "../worker/index.ts";

const files = {
  "/guide/": ["text/html; charset=utf-8", "<h1>Guide</h1>"],
  "/guide.md": ["text/markdown; charset=utf-8", "# Guide\n"],
  "/index.md": ["text/markdown; charset=utf-8", "# goesm\n"],
  "/llms.txt": ["text/plain; charset=utf-8", "# goesm\n"],
};
const assets = {
  async fetch(request) {
    const file = files[new URL(request.url).pathname];
    if (!file) return new Response("<h1>404</h1>", { status: 404, headers: { "content-type": "text/html; charset=utf-8" } });
    return new Response(file[1], { headers: { "content-type": file[0] } });
  },
};
const get = (path, accept) => handle(new Request("https://goesm.dev" + path, { headers: accept ? { accept } : {} }), assets);
const browser = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8";

test("wantsMarkdown", () => {
  for (const accept of [null, "", "*/*", "text/plain", "text/markdown", "text/markdown, text/html;q=0.9", "text/markdown;q=1, text/html;q=1"]) {
    assert.equal(wantsMarkdown(accept), true, String(accept));
  }
  for (const accept of [browser, "text/html", "text/html, text/markdown;q=0.5", "application/xhtml+xml", "text/markdown;q=0"]) {
    assert.equal(wantsMarkdown(accept), false, accept);
  }
});

test("markdownPath", () => {
  assert.equal(markdownPath("/"), "/index.md");
  assert.equal(markdownPath("/ja/"), "/ja.md");
  assert.equal(markdownPath("/guide/getting-started/"), "/guide/getting-started.md");
  assert.equal(markdownPath("/gosfc/guide"), "/gosfc/guide.md");
});

test("curl gets Markdown", async () => {
  const res = await get("/guide/");
  assert.equal(res.status, 200);
  assert.equal(res.headers.get("content-type"), "text/markdown; charset=utf-8");
  assert.equal(res.headers.get("vary"), "Accept");
  assert.equal(await res.text(), "# Guide\n");
  assert.equal(await (await get("/")).text(), "# goesm\n");
});

test("browsers get HTML with a link to the Markdown", async () => {
  const res = await get("/guide/", browser);
  assert.equal(res.headers.get("content-type"), "text/html; charset=utf-8");
  assert.equal(res.headers.get("vary"), "Accept");
  assert.equal(res.headers.get("link"), '</guide.md>; rel="alternate"; type="text/markdown"');
  assert.equal(await res.text(), "<h1>Guide</h1>");
});

test("files are served as they are", async () => {
  for (const accept of [undefined, browser]) {
    const res = await get("/llms.txt", accept);
    assert.equal(res.headers.get("content-type"), "text/plain; charset=utf-8");
    assert.equal(res.headers.get("vary"), null);
  }
});

test("a missing page is a Markdown 404 pointing to llms.txt", async () => {
  const res = await get("/gosfc/ja/nope/");
  assert.equal(res.status, 404);
  assert.equal(res.headers.get("content-type"), "text/markdown; charset=utf-8");
  assert.match(await res.text(), /https:\/\/goesm\.dev\/gosfc\/ja\/llms\.txt/);
  const html = await get("/nope/", browser);
  assert.equal(html.status, 404);
  assert.equal(html.headers.get("link"), null);
});
