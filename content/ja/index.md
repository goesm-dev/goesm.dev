---
layout: home
description: goesm は Go のパッケージを、Vite・Bun・Node.js・ブラウザ向けのネイティブな ES モジュールにコンパイルします。TypeScript の型付きの素の JavaScript で、WebAssembly は使いません。
hero:
  name: goesm
  text: Go のパッケージを、ネイティブな ES モジュールに
  tagline: TypeScript の型が付いた素の JavaScript。WebAssembly もグルーコードもなく、バンドラーは Go を他のモジュールと同じように import します。
  image: /goesm-hero.webp
  actions:
    - theme: brand
      text: はじめる
      link: /guide/getting-started
    - theme: alt
      text: goesm とは
      link: /guide/
    - theme: alt
      text: GitHub
      link: https://github.com/goesm-dev/goesm
features:
  - icon: ⚡
    title: Go の関数が JS の関数になる
    details: wasm のインスタンスも wasm_exec.js も値の変換もありません。export した Go の関数の呼び出しコストは、普通の JS 関数と同じです。
  - icon: 📦
    title: パッケージごとに 1 つの ES モジュール
    details: 相対 .ts import でつながる TypeScript モジュールのツリーを出力します。Vite・Rolldown・esbuild・Bun が tree shaking、分割、minify を行い、source map は .go まで戻ります。
  - icon: 🧪
    title: ネイティブ Go と突き合わせて検証
    details: すべての fixture をネイティブ Go と ESM の両方で実行し、結果の一致を確認しています。Go 自身のテストスイートは 961 件中 898 件が通ります。
    link: /reference/conformance
  - icon: 🛠️
    title: Go は Go のまま
    details: フロントエンドは go/packages と go/types です。go.mod、gopls、go vet、go test は同じコードでそのまま動きます。
  - icon: 🧩
    title: Vue と Astro
    details: gosfc を使うと、Vue コンポーネントの <script setup lang="go"> で本物の Go が動き、Astro のページからも Go のパッケージを直接 import できます。
    link: /guide/vue
  - icon: 🔭
    title: ゼロコード計装
    details: goesm build -toolexec "otelc toolexec" で、ネイティブのビルドと同じ OpenTelemetry の span を出力します。
    link: /reference/otelc
---

## Go を書いて、ESM として import する

```go
package cart

type Item struct {
	Name     string
	Price    int
	Quantity int
}

func Total(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Price * item.Quantity
	}
	return total
}
```

```sh
go tool goesm emit-ts ./cart   # goesm-ts/example.com/app/cart.ts とランタイム
```

```ts
import { Item, Total, $runtime as rt } from "./goesm-ts/example.com/app/cart.ts";

Total(rt.sliceLit([new Item(rt.fromJSString("apple"), 120, 3)])); // 360
```
