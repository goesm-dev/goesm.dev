// The search index of a locale of the gosfc site (see ../../search/).
import type { APIRoute } from "astro";
import { SearchIndex, $runtime as rt } from "go:goesm.dev/site";

export function getStaticPaths() {
  return [{ params: { locale: "en" } }, { params: { locale: "ja" } }];
}

export const GET: APIRoute = ({ params }) =>
  new Response(rt.toJSString(SearchIndex(rt.fromJSString("/gosfc"), rt.fromJSString(params.locale!))), {
    headers: { "content-type": "text/plain; charset=utf-8" },
  });
