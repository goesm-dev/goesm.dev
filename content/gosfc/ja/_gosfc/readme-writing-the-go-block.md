<!-- Synced by cmd/syncdocs. Edit the source instead. -->

* トップレベルは Vue の `<script setup>` と同じくコンポーネントインスタンスごとに上から 1 回実行されます。`x := ...`、`var`、`const`、`type`、`func F() {...}` が書けます。
* トップレベルの変数・定数・関数はテンプレートから参照できます。Go 関数をテンプレートから呼ぶ（`@click="Increment"` など）と、表示が Go の値に追従します。
* インポートは Go のインポートだけです。`.vue`、`.ts`、`.go` ファイルのインポートはできません。
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
