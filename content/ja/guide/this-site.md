---
description: goesm.dev は goesm で作られています。VitePress 風のドキュメントエンジンは Go のパッケージで、Astro が実行し、gosfc を通して Vue が描画します。
---

# このサイトの作り方

goesm.dev は goesm で作られています。ドキュメントエンジン、Markdown のレンダラー、
コードのハイライト、検索はすべて Go のパッケージです。goesm がそれを ES モジュールに
コンパイルし、[gosfc](./vue.md) が Astro と Vue につなぎ、Astro が静的な HTML を
書き出します。ソースは [goesm-dev/goesm.dev](https://github.com/goesm-dev/goesm.dev) です。

```text
content/**/*.md ──► press (Go) ──► goesm ──► Astro のページ + Vue のテーマ ──► 静的 HTML
                                                │
                                                └─► ブラウザで動くアイランド (これも Go)
```

## press: Go で書いた VitePress 風エンジン

`goesm.dev/press` は、ドキュメントサイトに対して VitePress がすることをします。

- [goldmark](https://github.com/yuin/goldmark) による Markdown: CommonMark、
  GitHub のテーブル、取り消し線、自動リンク
- `::: tip` / `info` / `warning` / `danger` / `details` コンテナと GitHub の
  アラート (`> [!NOTE]`)
- 言語ラベル、コピーボタン、行ハイライト (` ```go {2,4-5}`) 付きのコードブロック。
  ハイライターは `press/highlight`
- GitHub 互換の見出しアンカーと「このページの内容」
- `<!--@include: ./file.md-->`、`.md` へのリンクをページのルートに書き換え
- ナビゲーションバー、サイドバー、前後のページ、「GitHub でこのページを編集」
- ロケール: `/` が英語、`/ja/` が日本語。UI の翻訳、同じページのまま切り替える
  言語メニュー、`hreflang`
- ロケールごとの検索インデックスと `sitemap.xml`

`goesm.dev/site` は設定 (VitePress なら `.vitepress/config` に書くもの) です。
ロケール、ナビゲーション、サイドバーを持ちます。ページそのものは `content/en` と
`content/ja` の Markdown で、`go:embed` で埋め込んでいます。

::: tip go test でテスト
press は普通の Go なので、`go test` でネイティブにテストします。同じコードが
ビルド時には goesm でコンパイルされて Node.js 上で動きます。
:::

## Astro と Vue

Astro のファイルも Go で書いています。ページの外枠である `src/theme/Layout.astro` は、
gosfc の `---go` による Go のフロントマターを持ちます。そこでテーマのコンポーネントを
import し、`<head>` に必要な情報をエンジンから受け取ります。

```astro
---go
import (
	"goesm.dev/site"

	DocPage "./DocPage.vue"
	NavBar "./NavBar.vue"
	_ "../styles/theme.css"
)

type Props struct {
	Route string
}

route := props.Route
head := site.PageHead(route)
isDoc := head.Found && head.Layout == "doc"
---

<html lang={head.lang}>
  <head>
    <title>{head.title}</title>
```

次の部分は、それぞれの理由で TypeScript のままです。

- `src/pages/[...slug].astro` は `getStaticPaths` でページを列挙します。これは export であり、Go のフロントマターからは export できません。このフロントマターも `go:` の import で Go を呼び、`Routes().map(...)` とするだけです。
- Markdown 版、llms.txt、検索インデックス、サイトマップは Astro のエンドポイント（`.ts` ファイル）です。どれも Go の関数が書いた内容を返す数行のコードです。
- 最初の描画の前にダークテーマを設定するスクリプトは、インラインの `<script>` です。`<script lang="go">` はバンドルされ、ページの解析後に実行されるため間に合いません。
- `src/theme/ClientRouter.ts` は、Astro の `ClientRouter` をデフォルトエクスポートとして再エクスポートします。Go の import 構文はデフォルトエクスポートしか import できないためです。

テーマは `<script setup>` が Go の Vue コンポーネントです。たとえばドキュメントの
ページのコンポーネントは、エンジンを呼んで結果をテンプレートに渡します。

```vue
<script setup lang="go">
import "goesm.dev/site"

type Props struct {
	Route string
}

doc := site.Doc(props.Route)
</script>
```

Astro はこれらのコンポーネントをビルド時に HTML にします。これらのために
ページが JavaScript を配ることはありません。Markdown のレンダリング、ハイライト、
サイドバーの組み立ては、ビルド時に一度だけ行われます。

## ブラウザで動く Go

いくつかのコンポーネントは Astro のアイランドとしてハイドレートされ、その Go のコードが
ブラウザで動きます。

| アイランド | Go のパッケージ | すること | サイズ (gzip) |
| --- | --- | --- | --- |
| 検索 | `client/searchbox`、`press/search`、`client/webmcp` | ロケールのインデックスを読み込み、順位付けとスニペット作成。WebMCP のツール | 9 KB |
| ライブデモ (ホーム) | `client/demo` | カート。クリックのたびに Go の `int` で計算 | 3 KB |
| ページの補助 | `client/enhance` | コピーボタン、スクロールに追従する目次 | 1.5 KB |
| テーマ切り替え | (コンポーネント内の Go) | ダークモード。`localStorage` に記憶 | 1 KB |

Vue 本体のほかに、goesm のランタイムと `syscall/js` を共有します。gzip で約 11 KB で、
読み込みは一度だけです。

検索の仕組みは VitePress のローカル検索と同じです。press がビルド時に全セクションの
プレーンテキストのインデックスを書き出し、ブラウザは最初の検索のときにそれを
読み込みます。検索のコードは、あえて `strings`、`sort`、`strconv` を使っていません。
今の goesm では、これらを import すると Unicode のテーブルやリフレクションも一緒に
入り、アイランドが何倍にも大きくなるためです。

## AI エージェント向け

このサイトは、人だけでなくプログラムからも読めるように作っています。

- 各ページには Markdown 版があります。`/guide/getting-started/` の Markdown 版は
  `/guide/getting-started.md` です。インクルードを展開し、コンテナを引用に変え、
  リンクを絶対 URL にしています。各ページからは「Markdown で表示」のリンクと
  `<link rel="alternate" type="text/markdown">` でたどれます。
- サイトと言語ごとに、ページの一覧である [llms.txt](https://llmstxt.org/) と、
  全ページをまとめた `llms-full.txt` があります。場所は `/llms.txt`、`/ja/llms.txt`、
  `/gosfc/llms.txt`、`/gosfc/ja/llms.txt` です。
- 静的ファイルの前に置いた Worker (`worker/index.ts`) は `Accept` ヘッダーを見ます。
  HTML を求めるブラウザにはページを返し、`curl` や `fetch` など HTML を求めない
  クライアントには、同じ URL で Markdown を返します。`curl https://goesm.dev/ja/guide/`
  の結果は Markdown です。
- [WebMCP](https://github.com/webmachinelearning/webmcp) に対応したブラウザでは、
  ページが読み取り専用のツールを 3 つ `document.modelContext` に登録します。ページで
  作業する AI エージェントは、画面を読む代わりにこれらを使えます。検索ダイアログの
  インデックスを引く `search_docs`、ページを Markdown で返す `get_page`、llms.txt を
  返す `list_pages` です。これらも Go で、検索のアイランドに入っています。

Markdown 版と llms.txt は、HTML と同じページから press がビルド時に書き出します。

## WebAssembly は不要でした

ブラウザでの重い処理には TinyGo と WebAssembly を使う選択肢もありました。このサイトでは
必要になりませんでした。重い処理 (Markdown、ハイライト、インデックス作成) はビルド時に
済んでいて、インデックス (1 言語あたり約 100 KB のテキスト) の検索は素の JavaScript で十分に速いからです。アイランドは
小さいままで、DOM を直接呼べます。

## リポジトリから同期するドキュメント

リファレンスのページと README の各セクションは、手でコピーしていません。
`cmd/syncdocs` が goesm と gosfc のリポジトリから `go.mod` のバージョンのものを読み、
リンクと画像を書き換えて `content/` に書き出します。CI は `-check` 付きで実行するので、
サイトが説明しているバージョンとずれることはありません。
