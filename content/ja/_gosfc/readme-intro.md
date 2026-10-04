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

`.go` は普通の Go package、import は普通の Go import です。Go のコンパイルは [goesm](https://github.com/goesm-dev/goesm)、SFC と template は Vue tooling、build は Vite、ページと SSR は Astro が担当します。設計は [ARCHITECTURE.md](https://github.com/goesm-dev/gosfc/blob/60ac235c2d5339028c79d300a40265e0c3d14ec6/ARCHITECTURE.md) を見てください。

**状態：PoC。** 未実装の項目は ARCHITECTURE.md の「未実装・未決事項」にあります。
