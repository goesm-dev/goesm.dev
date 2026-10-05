import type { APIRoute } from "astro";
import { Sitemap } from "go:goesm.dev/site";

export const GET: APIRoute = () =>
  new Response(Sitemap(), { headers: { "content-type": "application/xml; charset=utf-8" } });
