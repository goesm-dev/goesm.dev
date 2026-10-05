import { rename, rmdir } from "node:fs/promises";
import { defineConfig } from "astro/config";
import gosfc from "@gosfc/astro";

// Cloudflare's asset server answers a missing page with the nearest 404.html
// above it (cloudflare.config.ts), so each 404 page Astro writes as
// <dir>/404/index.html (all but the root one) moves to <dir>/404.html.
const localized404 = {
  name: "localized-404",
  hooks: {
    "astro:build:done": async ({ dir }) => {
      for (const d of ["ja/", "gosfc/", "gosfc/ja/"]) {
        await rename(new URL(d + "404/index.html", dir), new URL(d + "404.html", dir));
        await rmdir(new URL(d + "404/", dir));
      }
    },
  },
};

export default defineConfig({
  site: "https://goesm.dev",
  trailingSlash: "always",
  // The stylesheet is small (about 5 KB gzip): inlining it saves a round trip
  // before the first paint, which is most of the wait on slow connections.
  build: { inlineStylesheets: "always" },
  integrations: [gosfc(), localized404],
});
