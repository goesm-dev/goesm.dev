// The Worker in front of goesm.dev's static files. Browsers get the HTML
// page as before; clients that do not ask for HTML (curl, fetch, AI agents)
// get the page's Markdown version from the same URL, so
// `curl https://goesm.dev/guide/` prints Markdown. Every page also has its
// Markdown at its own path (/guide.md), and each site and language has an
// llms.txt (/llms.txt, /ja/llms.txt, /gosfc/llms.txt, /gosfc/ja/llms.txt).
//
// The Worker runs only for page URLs (runWorkerFirst in cloudflare.config.ts);
// /_astro/ and the other files are served without it.

interface Assets {
  fetch(request: Request): Promise<Response>;
}

export interface Env {
  ASSETS: Assets;
}

export default {
  fetch(request: Request, env: Env): Promise<Response> {
    return handle(request, env.ASSETS);
  },
};

// The home page of each site and language, longest first: a missing page's
// Markdown 404 points to the llms.txt of the one it is under.
const homes = ["/gosfc/ja/", "/gosfc/", "/ja/", "/"];

export async function handle(request: Request, assets: Assets): Promise<Response> {
  const url = new URL(request.url);
  if ((request.method !== "GET" && request.method !== "HEAD") || isFile(url.pathname)) {
    return assets.fetch(request);
  }
  if (!wantsMarkdown(request.headers.get("accept"), request.headers.get("sec-fetch-site"))) {
    const res = await assets.fetch(request);
    return withHeaders(res, (h) => {
      h.append("vary", "Accept, Sec-Fetch-Site");
      if (res.ok && h.get("content-type")?.startsWith("text/html")) {
        h.set("link", `<${markdownPath(url.pathname)}>; rel="alternate"; type="text/markdown"`);
      }
    });
  }
  const md = markdownPath(url.pathname);
  const res = await assets.fetch(new Request(new URL(md, url), { method: request.method }));
  if (res.ok) {
    return withHeaders(res, (h) => {
      h.set("content-type", "text/markdown; charset=utf-8");
      h.set("content-location", md);
      h.append("vary", "Accept, Sec-Fetch-Site");
    });
  }
  const home = homes.find((p) => url.pathname.startsWith(p)) ?? "/";
  const body = `# Not found\n\nThere is no page at ${url.pathname}. The pages are listed in ${url.origin}${home}llms.txt.\n`;
  return new Response(request.method === "HEAD" ? null : body, {
    status: 404,
    headers: { "content-type": "text/markdown; charset=utf-8", vary: "Accept, Sec-Fetch-Site" },
  });
}

// wantsMarkdown reports whether a request with this Accept header should get
// Markdown: when it names text/markdown, or does not name HTML at all
// (curl's */*, a missing header). Browsers name text/html when they load a
// page. The site's own pages fetch pages with */* too (Astro's client router
// and its prefetching). The browser marks those requests with
// Sec-Fetch-Site: same-origin, which curl and other clients do not send, so
// they get HTML unless they name text/markdown.
export function wantsMarkdown(accept: string | null, fetchSite: string | null = null): boolean {
  const types = new Map<string, number>();
  for (const part of (accept ?? "").toLowerCase().split(",")) {
    const [type, ...params] = part.split(";").map((s) => s.trim());
    if (!type) continue;
    const q = params.find((p) => p.startsWith("q="));
    types.set(type, q ? Number(q.slice(2)) || 0 : 1);
  }
  const html = Math.max(types.get("text/html") ?? 0, types.get("application/xhtml+xml") ?? 0);
  const md = types.get("text/markdown");
  if (md === undefined) return html === 0 && fetchSite !== "same-origin";
  return md > 0 && md >= html;
}

// markdownPath is the Markdown version of a page path (press.MarkdownPath):
// /guide/x/ -> /guide/x.md, / -> /index.md, /ja/ -> /ja.md.
export function markdownPath(pathname: string): string {
  const p = pathname.replace(/\/+$/, "");
  return (p === "" ? "/index" : p) + ".md";
}

// isFile reports whether a path names a file (/llms.txt, /guide.md) rather
// than a page; files are served as they are.
function isFile(pathname: string): boolean {
  const last = pathname.slice(pathname.lastIndexOf("/") + 1);
  return last.includes(".");
}

function withHeaders(res: Response, edit: (h: Headers) => void): Response {
  const headers = new Headers(res.headers);
  edit(headers);
  return new Response(res.body, { status: res.status, statusText: res.statusText, headers });
}
