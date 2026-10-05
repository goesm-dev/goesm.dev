<!-- Synced by cmd/syncdocs. Edit the source instead. -->

Vue Single File Component の `<script setup>` で本物の Go を使うための統合レイヤーです。

```vue
<template>
  <div>
    合計: {{ total }}
  </div>
</template>

<script setup lang="go">
import cart "example.com/app/src/features/cart/pkg"

items := []cart.Item{
	{
		Price:    100,
		Quantity: 2,
	},
}

total := cart.Total(items)
</script>
```

`.go` は普通の Go パッケージ、インポートは普通の Go のインポートです。Go のコンパイルは [goesm](https://github.com/goesm-dev/goesm)、SFC とテンプレートは Vue tooling、ビルドは Vite、ページと SSR は Astro が担当します。設計は [ARCHITECTURE.ja.md](/gosfc/ja/reference/architecture/) を見てください。

**状態：PoC。** 未実装の項目は ARCHITECTURE.ja.md の「未実装・未決事項」にあります。
