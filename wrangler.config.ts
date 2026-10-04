// Build settings for `wrangler build`, which turns Astro's output into the
// Build Output that `cf deploy --prebuilt` uploads. The Worker itself is
// described in cloudflare.config.ts. (`cf build` would do both steps, but its
// project detection stops here: it finds three Astro projects in this pnpm
// workspace, the site and the two gosfc packages.)
import { defineWranglerConfig } from "wrangler/experimental-config";

export default defineWranglerConfig({
  assetsDirectory: "./dist",
});
