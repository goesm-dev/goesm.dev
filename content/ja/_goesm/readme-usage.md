<!-- Synced by cmd/syncdocs. Edit the source instead. -->

### `goesm emit-ts`: バンドラー向けの TypeScript ツリー

入力は、普通の Go モジュールでの普通の Go パッケージパターンです。

```sh
cd testdata/example
go run ../../cmd/goesm emit-ts -o goesm-ts ./main   # -o の既定は goesm-ts
```

```text
goesm-ts/
├── example.com/app/main.ts    import * as mathx from "./mathx.ts"
├── example.com/app/mathx.ts   import * as $rt from "../../@goesm/runtime/index.ts"
└── @goesm/runtime/
    ├── index.ts
    └── ...                    その他のランタイムのファイル
```

Go パッケージ `p` のモジュールは `<dir>/<p>.ts` に出力されます。標準ライブラリも同じで、`strings.ts` や `internal/bytealg.ts` のようになります。ランタイムは `<dir>/@goesm/runtime/` に出力されます。Go の import パスは `@` で始められないので、ランタイムの名前が Go のパッケージとぶつかることはありません。モジュールどうしは `.ts` で終わる相対指定子で import し合うので、リゾルバーもプラグインもバンドラーの設定も要りません。エントリのパッケージを自分のコードから import します。

```ts
// Vite プロジェクトの index.ts。bun index.ts や node index.ts で直接実行することもできる。Node.js は 22.18 以降が必要
import { Result } from "./goesm-ts/example.com/app/main.ts";
console.log(Result()); // 3
```

ツリーを import するコードを `tsc` で型検査するには `allowImportingTsExtensions` を有効にします。生成される TypeScript と JavaScript は [docs/example-output.ja.md](/ja/reference/example-output/) にあります。

### `goesm build`: バンドル済みの ES モジュール

`goesm build` は同じツリーを一時ディレクトリに書き出し、esbuild の Go API でバンドルします。バンドラーを使わない場合と、goesm 自身のテストのための機能です。

```sh
go run ../../cmd/goesm build ./main            # dist/main.js と、.go を指す dist/main.js.map を出力
go run ../../cmd/goesm build -minify ./main
go run ../../cmd/goesm build -split ./main     # Go パッケージごとに 1 つの ES モジュール: dist/example.com/app/main.js など
```

```js
import { Result } from "./dist/main.js";
Result(); // 3
```

出力は拡張子 `.js` の ES モジュールです。Node.js では、`package.json` に `"type": "module"` があるパッケージから読み込んでください。そうでないと Node.js はモジュール形式を判定するために各モジュールをもう一度パースし、大きなバンドルでは起動に数十ミリ秒が加わります。

### 再ビルド

goesm は、変換した各パッケージの TypeScript モジュールをキャッシュに保存します。キャッシュの場所は、ユーザーキャッシュディレクトリの `goesm/modules` です。`GOESMCACHE` で場所を変更でき、`GOESMCACHE=off` でキャッシュを無効にできます。再ビルドでは、モジュールが変わりうるパッケージだけを変換し直します。対象は、編集したパッケージ、それに依存するパッケージ、編集によってプログラム全体の解析結果が変わったパッケージです。たとえば、依存先の関数がコールバックを await する必要が新たに生じた場合は、その関数のパッケージが対象になります。出力は、キャッシュの有無にかかわらず同じです。goesm.dev の `site` パッケージは 77 個のパッケージを含みます。このパッケージの `emit-ts` は、キャッシュがない場合は 1.2 秒、1 つのパッケージを編集した後は 0.7 秒かかります。残りの時間は、Go のフロントエンドとプログラム全体の解析にかかります。

### JavaScript から Go を呼ぶ

- ビルドしたパッケージの export された関数は、そのモジュールの export になります。引数と戻り値は JavaScript の値で、Go の型に従って変換されます。文字列は JS の文字列、スライスは配列、構造体は `encoding/json` と同じ名前のフィールドを持つプレーンオブジェクト、`int64` / `uint64` は BigInt です。TypeScript にもこの型が見えます。たとえば `Total` の型は `Total(items: Array<{ Name?: string; Price?: number; Quantity?: number }> | null): number` です。
- メソッドを持つ構造体型へのポインタは Go のオブジェクトそのもので、JavaScript からそのメソッドを呼べます (`cart.Add(item)`)。
- 複数の戻り値は配列として返ります。最後の戻り値の `error` は `GoError` として投げられ、Go に渡し直すと元の Go のエラーに戻ります。
- チャネル操作、`time.Sleep`、ミューテックスの待ちのようにブロックしうる関数は、Promise を返す `async function` になります。それ以外の関数は同期関数です。

詳しくは [docs/js-exports.ja.md](/ja/reference/js-exports/) を参照してください。[examples/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.3/examples) には、呼び出し側の JavaScript と組み合わせて実行できる例があります。

### Go から JavaScript を呼ぶ

本体のない関数を宣言すると、ES モジュールの関数を取り込めます。値の変換は Go の型に従います。

```go
//goesm:import "./format.ts" formatPrice
func formatPrice(yen int, currency string) string

//goesm:import "./api.ts" fetchUser await
func fetchUser(id string) (User, error) // Promise を待ち、例外と reject は error になる
```

文字列、スライス、マップ、構造体、関数、`js.Value`、`any` が境界を越えられます。構造体のプロパティ名は `encoding/json` と同じ規則で決まります。1 回の呼び出しは、JavaScript から同じ関数を呼ぶ場合より数ナノ秒多くかかるだけです。詳しくは [docs/js-imports.ja.md](/ja/reference/js-imports/) にまとめています。Vue コンポーネントは [gosfc](https://github.com/goesm-dev/gosfc) から使います。

### HTTP を処理する

`ServeMux`、Connect のサービス、ミドルウェアなどの `http.Handler` は、どのホストでも、ホスト自身のサーバーを通じてリクエストを処理します。

```go
// Node.js、Bun、Deno では、node:http、Bun.serve、Deno.serve が処理する。
log.Fatal(http.ListenAndServe(":8080", api.Handler()))
```

```ts
// Cloudflare Workers、Deno.serve、Bun.serve、Service Worker では、fetch ハンドラとして公開する。
import { Handler } from "./goesm-ts/example.com/app/api.ts";
export default { fetch: Handler() }; // 戻り値の http.Handler は fetch ハンドラになる
```

各リクエストは専用の goroutine で処理され、ボディは最初に全部読み込まれます。レスポンスはハンドラが戻った時点で送られます。ハンドラが flush した場合は、その時点からレスポンスがストリームになります。Server-Sent Events や Connect のサーバーストリーミングはこの仕組みで動きます。Workers では、`nodejs_compat` フラグを有効にすると、Worker のテキストバインディングとシークレットを `os.Getenv` で読めます。HTTP クライアントは `fetch` を使います。どこで何に対応しているかは [docs/use-cases.ja.md](/ja/reference/use-cases/) にまとめています。
