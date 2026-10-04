import type { APIRoute } from "astro";
import { Sitemap, $runtime as rt } from "go:goesm.dev/site";

export const GET: APIRoute = () =>
  new Response(rt.toJSString(Sitemap()), { headers: { "content-type": "application/xml; charset=utf-8" } });
