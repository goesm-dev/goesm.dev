<!-- Synced by cmd/syncdocs. Edit the source instead. -->

* トップレベルは Vue の `<script setup>` と同じくコンポーネントインスタンスごとに上から 1 回実行されます。`x := ...`、`var`、`const`、`type`、`func F() {...}` が書けます。
* トップレベルの変数・定数・関数はテンプレートから参照できます。Go 関数をテンプレートから呼ぶ（`@click="Increment"` など）と、表示が Go の値に追従します。値や、呼び出しの引数と戻り値は goesm の[エクスポートされた関数](https://github.com/goesm-dev/goesm/blob/main/docs/js-exports.ja.md)と同じく変換されます。文字列は文字列、スライスは配列、struct は普通のオブジェクトになります。
* `import` で取り込めるのは Go のパッケージだけです。`.astro` ファイルの Go のフロントマターでは、後述のとおり JavaScript のモジュールも import できます。Vue コンポーネントの Go ブロックには、Vue コンポーネント、TypeScript、JavaScript を goesm の [`//goesm:import`](https://github.com/goesm-dev/goesm/blob/main/docs/js-imports.ja.md) ディレクティブで取り込みます。コンポーネントやクラスなどの値は `var` の前に、関数は本体のない関数宣言の前にディレクティブを書きます。関数の引数と戻り値は自動で変換されます。どちらもブロックのほかの名前と同じくテンプレートから使えます。これらはパッケージレベルで宣言されるので、型にはインポートした型と組み込みの型を使え、ブロック内で宣言した型は使えません。

  ```vue
  <template>
    <p>{{ formatPrice(12800) }} <Badge label="新着" /></p>
  </template>

  <script setup lang="go">
  import "syscall/js"

  //goesm:import "./Badge.vue"
  var Badge js.Value

  //goesm:import "./format.ts" formatPrice
  func formatPrice(yen int) string
  </script>
  ```
* メソッドとジェネリック関数は Go パッケージに置いてください。
* props を受け取るには `type Props struct {...}` を宣言します。ブロックの中でその型の `props` 変数が使えます。フィールド `Route` は属性 `route`（json タグがあればその名前、または kebab-case の形）から読み、フィールドの型（string、bool、整数、浮動小数点数）に変換します。props はインスタンスの setup 時に 1 回だけ読み、ルート要素にはフォールスルーしません。

  ```vue
  <script setup lang="go">
  import "strings"

  type Props struct {
  	Title string
  }

  heading := strings.ToUpper(props.Title)
  </script>
  ```
