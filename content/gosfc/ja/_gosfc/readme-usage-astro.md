<!-- Synced by cmd/syncdocs. Edit the source instead. -->

1. Go モジュールに goesm と gosfc をツールとして追加します。バージョンは go.mod / go.sum で固定されます。gosfc には goesm v0.0.1-beta.3 以降が必要です。

   ```sh
   go get -tool github.com/goesm-dev/goesm/cmd/goesm@<version>
   go get -tool github.com/goesm-dev/gosfc/cmd/gosfc@<version>
   ```

2. Astro にインテグレーションを追加します。`@astrojs/vue` が無ければ追加されます。

   ```js
   // astro.config.mjs
   import { defineConfig } from "astro/config";
   import gosfc from "@gosfc/astro";

   export default defineConfig({
     integrations: [gosfc()],
   });
   ```

3. `.vue` で `<script setup lang="go">` を使います。

   ```astro
   ---
   import Summary from "../features/cart/Summary.vue";
   ---

   <Summary />
   ```

   `.astro` ファイルも Go で書けます。書き方は「[.astro ファイルで Go を使う](/gosfc/ja/guide/astro/)」で説明します。

Vite だけで使う場合は `@vitejs/plugin-vue` の前に `@gosfc/vite` を置きます。

```js
import vue from "@vitejs/plugin-vue";
import gosfc from "@gosfc/vite";

export default { plugins: [gosfc(), vue()] };
```

`lang="ts"` や素の `<script setup>` のコンポーネントはそのまま共存できます。gosfc が触るのは `lang="go"` のコンポーネントだけです。
