// The Markdown version of every page (/guide/getting-started.md), written by
// the Go engine at build time for AI agents and other programs. The Worker
// also serves it at the page's own URL to clients that do not ask for HTML.
import type { APIRoute } from "astro";
import { MarkdownPagesJSON, PageMarkdown } from "go:goesm.dev/site";

export function getStaticPaths() {
  const pages: { slug: string; route: string }[] = JSON.parse(MarkdownPagesJSON());
  return pages.map((p) => ({ params: { md: p.slug }, props: { route: p.route } }));
}

export const GET: APIRoute = ({ props }) =>
  new Response(PageMarkdown(props.route), {
    headers: { "content-type": "text/markdown; charset=utf-8" },
  });
