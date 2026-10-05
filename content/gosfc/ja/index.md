---
layout: home
description: gosfc は Vue コンポーネントの <script setup> で本物の Go を動かし、Astro のページから Go のパッケージを import できるようにします。goesm が素の JavaScript にコンパイルし、WebAssembly は使いません。
hero:
  name: gosfc
  text: Vue の SFC に Go を
  tagline: .vue ファイルに <script setup lang="go"> を書きます。中身は本物の Go で、goesm が素の JavaScript にコンパイルし、Vue が描画し、Vite と Astro がビルドします。
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
    title: Astro と SSR
    details: Astro にはインテグレーション 1 つで組み込めます。コンポーネントはサーバーで描画されてアイランドとしてハイドレートし、.astro ファイルからは go: で Go を import できます。
    link: /guide/importing-go
  - icon: ⚡
    title: WebAssembly なし
    details: goesm が Go を ES モジュールにコンパイルします。ベンチマークのページのクライアント JS は、同じコンポーネントを TypeScript で書いた場合より 13% 大きいだけです (gzip)。
    link: /guide/benchmark
  - icon: 🎯
    title: エラーは .vue ファイルを指す
    details: Go の診断メッセージ、ソースマップ、panic のスタックトレースは、すべて .vue ファイルの行を指します。
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

## 役割分担

gosfc は既存のツールをつなぐだけです。Go のコンパイルは goesm、SFC とテンプレートの
コンパイルは Vue のツール、ビルドは Vite、ページの描画は Astro が担当します。
詳しくは [アーキテクチャ](./reference/architecture.md) を見てください。
