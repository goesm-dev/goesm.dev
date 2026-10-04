// The Worker that serves goesm.dev on Cloudflare Workers. The site is static,
// so the Worker has no script: Cloudflare serves the files Astro writes to
// dist/ (wrangler.config.ts points there). Deploy with `pnpm run deploy`.
import { defineConfig } from "cf/config";

export default defineConfig({
  worker: {
    name: "goesm-dev",
    compatibilityDate: "2026-10-04",
    assets: {
      // Astro writes /guide/ as guide/index.html (trailingSlash: "always");
      // /guide and /guide/index.html redirect to /guide/.
      htmlHandling: "auto-trailing-slash",
      // A missing page gets the nearest 404.html: ja/404.html under /ja/.
      notFoundHandling: "404-page",
    },
  },
});
