// llms.txt and llms-full.txt (https://llmstxt.org/) of each site and locale:
// /llms.txt, /ja/llms.txt, /gosfc/llms.txt, /gosfc/ja/llms-full.txt, ...
import type { APIRoute } from "astro";
import { LLMs, LLMsPaths } from "go:goesm.dev/site";

export function getStaticPaths() {
  return LLMsPaths().map((llms) => ({ params: { llms } }));
}

export const GET: APIRoute = ({ params }) =>
  new Response(LLMs(`/${params.llms}.txt`), {
    headers: { "content-type": "text/plain; charset=utf-8" },
  });
