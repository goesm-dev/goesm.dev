import { defineConfig } from "astro/config";
import gosfc from "@gosfc/astro";

export default defineConfig({
  site: "https://goesm.dev",
  trailingSlash: "always",
  integrations: [gosfc()],
});
