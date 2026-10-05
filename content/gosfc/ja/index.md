---
layout: home
description: gosfc は Vue コンポーネントと .astro ファイルで本物の Go を動かします。.vue ファイルの <script setup> と、.astro ファイルのフロントマターと <script> に Go を書けます。goesm が素の JavaScript にコンパイルし、WebAssembly は使いません。
hero:
  name: gosfc
  text: Vue と Astro に Go を
  tagline: .vue ファイルには <script setup lang="go"> を、.astro ファイルには ---go のフロントマターを書きます。中身は本物の Go で、goesm が素の JavaScript にコンパイルし、Vue と Astro が描画し、Vite がビルドします。
  actions:
    - theme: brand
      text: はじめに
      link: /guide/getting-started
    - theme: alt
      text: gosfc とは
      link: /guide/
    - theme: alt
      text: GitHub
      link: https://github.com/goesm-dev/gosfc
features:
  - icon: 🐹
    title: 本物の Go
    details: Go ブロックは go.mod にある普通の Go のパッケージを import します。gopls、go vet、go test もそのまま使えます。
    link: /guide/go-block
  - icon: 🧩
    title: Vue の SFC のまま
    details: テンプレート、スタイル、ファイルの残りの部分は Vue に任せます。lang="ts" のコンポーネントも並べて使えます。
  - icon: 🚀
    title: .astro ファイルにも Go
    details: Astro にはインテグレーション 1 つで組み込めます。.astro ファイルのフロントマターを Go で書け、<script lang="go"> はブラウザで Go を実行します。Vue コンポーネントはサーバーで描画されてアイランドとしてハイドレートします。
    link: /guide/astro
  - icon: ⚡
    title: WebAssembly なし
    details: goesm が Go を ES モジュールにコンパイルします。ベンチマークのページのクライアント JS は、同じコンポーネントを TypeScript で書いた場合より 13% 大きいだけです (gzip)。
    link: /guide/benchmark
  - icon: 🎯
    title: エラーは書いたファイルを指す
    details: Go の診断メッセージ、ソースマップ、panic のスタックトレースは、すべて .vue や .astro ファイルの行を指します。
    link: /reference/architecture
  - icon: 🔁
    title: ホットリロード
    details: テンプレートの編集は Go の状態を保ったまま再描画し、Go の編集はコンポーネントを読み直します。使っている .go ファイルの編集も同じです。
    link: /reference/architecture
---

## Go で書いたコンポーネント

```vue
<template>
  <button type="button" @click="Increment">{{ count }} 回クリックしました</button>
</template>

<script setup lang="go">
count := 0

func Increment() {
	count++
}
</script>
```

Go ブロックは、Vue の `<script setup>` と同じくコンポーネントのインスタンスごとに 1 回
実行されます。トップレベルの変数と関数がテンプレートから見え、Go の関数がそれを変えると
ページが更新されます。

## Go で書いたページ

```astro
---go
import (
	"strconv"

	cart "example.com/app/src/features/cart/pkg"
	Summary "../features/cart/Summary.vue"
)

items := []cart.Item{{Price: 120, Quantity: 3}, {Price: 80, Quantity: 1}}
total := cart.Total(items)

func Yen(n int) string {
	return "¥" + strconv.Itoa(n)
}
---

<p>合計: {Yen(total)} ({items.length} 品目)</p>
<Summary />
```

フロントマターを `---go` で始めた `.astro` ファイルも Go で書けます。フロントマターは
ビルド時またはリクエストごとに、描画のたびに 1 回実行され、その値と関数をテンプレートから
使えます。テンプレートの中の `<script lang="go">` は、ブラウザで Go を実行します
([.astro ファイルで Go を使う](./guide/astro.md))。

## 役割分担

gosfc は既存のツールをつなぐだけです。Go のコンパイルは goesm、SFC とテンプレートの
コンパイルは Vue のツール、`.astro` ファイルのコンパイルとページの描画は Astro、ビルドは
Vite が担当します。
詳しくは [アーキテクチャ](./reference/architecture.md) を見てください。
