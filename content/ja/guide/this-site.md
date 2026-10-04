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

Astro のページは Go を直接呼びます。`src/pages/[...slug].astro` はページの一覧を
エンジンに尋ねます。

```astro
---
import { RoutesJSON, $runtime as rt } from "go:goesm.dev/site";

export function getStaticPaths() {
  const routes = JSON.parse(rt.toJSString(RoutesJSON()));
  return routes.map((r) => ({ params: { slug: r.slug || undefined }, props: { route: r.route } }));
}
---
```

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
| 検索 | `client/searchbox`、`press/search` | ロケールのインデックスを読み込み、順位付けとスニペット作成 | 8 KB |
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
