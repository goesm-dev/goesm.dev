// The Worker that serves goesm.dev on Cloudflare Workers. Cloudflare serves
// the files Astro writes to dist/ (wrangler.config.ts points there); the
// script in worker/ runs first for page URLs only, to answer clients that do
// not ask for HTML with the page's Markdown. Deploy with `pnpm run deploy`.
import { bindings, defineConfig } from "cf/config";

export default defineConfig({
  worker: {
    name: "goesm-dev",
    compatibilityDate: "2026-10-04",
    // Custom domain: Cloudflare creates its DNS record and certificate on
    // deploy (the goesm.dev zone must be on the account).
    domains: ["goesm.dev"],
    entrypoint: "./worker/index.ts",
    env: {
      ASSETS: bindings.assets(),
    },
    assets: {
      // Astro writes /guide/ as guide/index.html (trailingSlash: "always");
      // /guide and /guide/index.html redirect to /guide/.
      htmlHandling: "auto-trailing-slash",
      // A missing page gets the nearest 404.html: ja/404.html under /ja/.
      notFoundHandling: "404-page",
      // Built files go straight to the static files, without the Worker.
      runWorkerFirst: ["/*", "!/_astro/*", "!/repo/*", "!/*.*"],
    },
  },
});
