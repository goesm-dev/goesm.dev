import { defineConfig } from "astro/config";
import gosfc from "@gosfc/astro";

export default defineConfig({
  site: "https://goesm.dev",
  trailingSlash: "always",
  // The stylesheet is small (about 5 KB gzip): inlining it saves a round trip
  // before the first paint, which is most of the wait on slow connections.
  build: { inlineStylesheets: "always" },
  integrations: [gosfc()],
});
