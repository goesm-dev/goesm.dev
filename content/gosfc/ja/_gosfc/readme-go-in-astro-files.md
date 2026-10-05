<!-- Synced by cmd/syncdocs. Edit the source instead. -->

フロントマターを `---` ではなく `---go` で始めると、その中身は Go になります。規則は `.vue` の Go ブロックと同じです。トップレベルはレンダリングごとに 1 回実行されます。静的ビルドではビルド時、SSR ではリクエストごとの実行です。トップレベルの変数・定数・関数はテンプレートから参照できます。

```astro
---go
import (
	"strconv"

	cart "example.com/app/src/features/cart/pkg"
	Line "../features/cart/Line.astro"
	Summary "../features/cart/Summary.vue"
)

items := []cart.Item{
	{Price: 120, Quantity: 3},
	{Price: 80, Quantity: 1},
}
total := cart.Total(items)

func Yen(n int) string {
	return "¥" + strconv.Itoa(n)
}
---

<Line label="りんご" price={120} quantity={3} />
<p>合計: {Yen(total)}、{items.length} 品目</p>
<Summary />
```

* コンポーネントやスタイル、npm パッケージなどの JavaScript のモジュールは Go のインポート構文でインポートします。`import Card "./Card.vue"` は `import Card from "./Card.vue"` と同じ意味になり、`import _ "./global.css"` は副作用のためのインポートになります。インポートパスが `./`、`../`、`/`、`@` で始まるか `:` を含む場合は JavaScript のモジュールとして扱います。`astro:assets` も JavaScript のモジュールです。それ以外のパスは Go のパッケージです。
* `type Props struct {...}` を宣言すると `Astro.props` を受け取れます。フィールド名と props の名前の対応は `.vue` の Go ブロックと同じで、フィールド `Label` は `label` から読みます。
* テンプレートに渡る値は JavaScript 向けに変換されます。文字列は JavaScript の文字列、スライスは配列、struct はオブジェクトになります。テンプレートから Go の関数を呼ぶこともできます。
* チャネルや `time.Sleep` などで Go のコードがブロックする場合、フロントマターはその完了を待ちます。

テンプレートの中の `<script lang="go">` は、ブラウザでページごとに 1 回実行される Go です。Astro が処理する通常の `<script>` と同じように Astro がバンドルし、コンポーネントを何度使ってもページに含まれるのは 1 回だけです。トップレベルの名前はテンプレートから参照できません。DOM の操作には `syscall/js` を使います。

```astro
<button id="counter">0</button>

<script lang="go">
import (
	"strconv"
	"syscall/js"
)

count := 0
button := js.Global().Get("document").Call("getElementById", "counter")
button.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
	count++
	button.Set("textContent", strconv.Itoa(count))
	return nil
}))
</script>
```

制限は次のとおりです。

* Go のフロントマターは何もエクスポートできません。`getStaticPaths` などのエクスポートには TypeScript のフロントマターが必要です。TypeScript のフロントマターからも、次の節で説明する `go:` specifier で Go をインポートできます。
* Go のコードから使える `Astro` グローバルの情報は props だけです。`Astro.url`、`Astro.cookies`、リダイレクトなどを使うには TypeScript のフロントマターを書きます。
* Go のインポート構文で書けるのは default インポートと副作用のためのインポートだけです。JavaScript のモジュールの名前付きエクスポートを使うには、それを default エクスポートとして再エクスポートする小さなモジュールを作り、そのモジュールをインポートします。
* `<script lang="go">` には他の属性を付けられません。`is:inline` や `define:vars` も使えません。
