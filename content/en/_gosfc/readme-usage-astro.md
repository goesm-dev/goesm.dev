<!-- Synced by cmd/syncdocs. Edit the source instead. -->

1. Add goesm and gosfc to your Go module as tools. Their versions are pinned by go.mod / go.sum.

   ```sh
   go get -tool github.com/goesm-dev/goesm/cmd/goesm@<version>
   go get -tool github.com/goesm-dev/gosfc/cmd/gosfc@<version>
   ```

2. Add the integration to Astro. `@astrojs/vue` is added if it is not already there.

   ```js
   // astro.config.mjs
   import { defineConfig } from "astro/config";
   import gosfc from "@gosfc/astro";

   export default defineConfig({
     integrations: [gosfc()],
   });
   ```

3. Use `<script setup lang="go">` in your `.vue` files.

   ```astro
   ---
   import Summary from "../features/cart/Summary.vue";
   ---

   <Summary />
   ```

To use gosfc with Vite alone, put `@gosfc/vite` before `@vitejs/plugin-vue`.

```js
import vue from "@vitejs/plugin-vue";
import gosfc from "@gosfc/vite";

export default { plugins: [gosfc(), vue()] };
```

Components using `lang="ts"` or a plain `<script setup>` keep working alongside them. gosfc only touches components with `lang="go"`.
