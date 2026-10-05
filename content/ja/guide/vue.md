---
description: gosfc は Vue コンポーネントと .astro ファイルで本物の Go を動かし、JavaScript から Go のパッケージを import できるようにします。ドキュメントは goesm.dev/gosfc にあります。
---

# Vue と Astro (gosfc)

[gosfc](/gosfc/ja/) は goesm を Vue と Astro につなぎます。Vue コンポーネントには
`<script setup lang="go">` を、`.astro` ファイルには Go のフロントマター (`---go`) と
`<script lang="go">` を書けます。JavaScript のモジュールや Astro のページからは `go:` で
Go のパッケージを import できます。Go のコンパイルは goesm が行い、それ以外はいつもどおり
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
導入方法、Go ブロックの決まり、`.astro` ファイルでの Go、JavaScript からの Go の import、ベンチマークは
[gosfc のドキュメント](/gosfc/ja/guide/) にあります。
:::
