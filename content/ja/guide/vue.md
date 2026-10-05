---
description: gosfc は Vue コンポーネントの <script setup> で本物の Go を動かし、Astro のページから Go のパッケージを import できるようにします。ドキュメントは goesm.dev/gosfc にあります。
---

# Vue と Astro (gosfc)

[gosfc](/gosfc/ja/) は goesm を Vue と Astro につなぎます。Vue コンポーネントに
`<script setup lang="go">` を書けるようになり、Astro のページから `go:` で Go の
パッケージを import できます。Go のコンパイルは goesm が行い、それ以外はいつもどおり
Vue、Vite、Astro が担当します。

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

このサイトもこの方法で作られています。テーマは、スクリプトが Go の Vue コンポーネント
です ([このサイトの作り方](./this-site.md) を参照)。

::: tip gosfc のドキュメント
導入方法、Go ブロックの決まり、JavaScript からの Go の import、ベンチマークは
[gosfc のドキュメント](/gosfc/ja/guide/) にあります。
:::
