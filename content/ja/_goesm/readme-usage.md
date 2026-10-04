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

Go パッケージ `p` のモジュールは `<dir>/<p>.ts` です（標準ライブラリも同じで、`strings.ts`、`internal/bytealg.ts` など）。ランタイムは `<dir>/@goesm/runtime/` です（Go の import パスは `@` で始まれません）。モジュールどうしは `.ts` で終わる相対指定子で import し合うので、リゾルバーもプラグインもバンドラーの設定も要りません。エントリのパッケージを自分のコードから import します。

```ts
// Vite プロジェクトの index.ts。または直接実行: bun index.ts / node index.ts（Node.js 22.18 以降）
import { Result } from "./goesm-ts/example.com/app/main.ts";
console.log(Result()); // 3
```

ツリーを import するコードを `tsc` で型検査するには `allowImportingTsExtensions` を有効にします。生成される TypeScript と JavaScript は [docs/example-output.ja.md](/ja/reference/example-output/) にあります。

### `goesm build`: バンドル済みの ES モジュール

`goesm build` は同じツリーを一時ディレクトリに書き出し、esbuild の Go API でバンドルします。バンドラーがない場合（と goesm 自身のテスト）向けです。

```sh
go run ../../cmd/goesm build ./main            # dist/main.js（+ .go を指す .js.map）
go run ../../cmd/goesm build -minify ./main
go run ../../cmd/goesm build -split ./main     # Go パッケージごとに 1 つの ES モジュール: dist/example.com/app/main.js など
```

```js
import { Result } from "./dist/main.js";
Result(); // 3
```

出力は拡張子 `.js` の ES モジュールです。Node.js では、`package.json` に `"type": "module"` があるパッケージから読み込んでください。そうでないと Node.js はモジュール形式を判定するために各モジュールをもう一度パースし、大きなバンドルでは起動に数十ミリ秒が加わります。

### JavaScript から Go を呼ぶ

- export された関数と型は、パッケージのモジュールの export になり、TypeScript 上も Go の型を持ちます（`Total(items: $rt.S<Item>): number`）。
- 数値と真偽値は JS の number と boolean、`int64` / `uint64` は BigInt です。構造体は、フィールドを順に受け取るコンストラクタを持つクラスです（`new Item(name, price, quantity)`）。
- Go の文字列はバイト列です。渡すときは `rt.fromJSString(s)`、受け取るときは `rt.toJSString(s)` を使います。スライスは `rt.sliceLit([...])` で渡し、`rt.toArray(s)` で受け取ります。複数の戻り値は配列で、`error` は Go のインターフェース値で返ります。
- ブロックしうる関数（チャネル操作、`time.Sleep`、ミューテックスの待ち）は `async function` で Promise を返します。それ以外は同期関数です。

`rt` はランタイムで、すべてのモジュールが `$runtime` として再 export しています。[examples/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.0/examples) に、呼び出し側の JavaScript 付きで実行できる例（cart、標準ライブラリの利用、goroutine）があります。
