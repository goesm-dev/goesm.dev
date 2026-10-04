// The search index of a locale, written by the Go engine at build time and
// fetched by the search dialog.
import type { APIRoute } from "astro";
import { SearchIndex, $runtime as rt } from "go:goesm.dev/site";

export function getStaticPaths() {
  return [{ params: { locale: "en" } }, { params: { locale: "ja" } }];
}

export const GET: APIRoute = ({ params }) =>
  new Response(rt.toJSString(SearchIndex(rt.fromJSString(params.locale!))), {
    headers: { "content-type": "text/plain; charset=utf-8" },
  });
