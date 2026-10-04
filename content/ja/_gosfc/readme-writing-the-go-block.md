<!-- Synced by cmd/syncdocs. Edit the source instead. -->

* トップレベルは Vue の `<script setup>` と同じく component instance ごとに上から 1 回実行されます。`x := ...`、`var`、`const`、`type`、`func F() {...}` が書けます。
* トップレベルの変数・定数・関数は template から参照できます。Go 関数を template から呼ぶ（`@click="Increment"` など）と、表示が Go の値に追従します。
* import は Go の import だけです。`.vue`、`.ts`、`.go` ファイルの import はできません。
* メソッドと generic 関数は Go package に置いてください。
* props を受け取るには `type Props struct {...}` を宣言します。block の中でその型の `props` 変数が使えます。フィールド `Route` は属性 `route`（json tag があればその名前、または kebab-case の形）から読み、フィールドの型（string、bool、整数、浮動小数点数）に変換します。props は instance の setup 時に 1 回だけ読み、root 要素には fallthrough しません。

  ```vue
  <script setup lang="go">
  import "strings"

  type Props struct {
  	Title string
  }

  heading := strings.ToUpper(props.Title)
  </script>
  ```
