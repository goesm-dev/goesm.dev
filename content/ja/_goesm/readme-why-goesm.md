<!-- Synced by cmd/syncdocs. Edit the source instead. -->

- **Go の関数がそのまま JavaScript の関数になります。** WebAssembly のインスタンスも、`wasm_exec.js` も、非同期のインスタンス化も、`syscall/js` 経由の値の変換もありません。`Discount` の呼び出しコストは普通の JS 関数の呼び出しと同じです。数値、真偽値、構造体は、変換されずにそのまま境界を越えます。
- **Go パッケージ 1 つが ES モジュール 1 つになります。** 出力は、相対 `.ts` 指定子で互いを import する TypeScript モジュールのツリーです。tree shaking、コード分割、minify、`.go` ファイルを指すソースマップの生成は、ホストのバンドラーが受け持ちます。TypeScript からは Go の API の型が見えます。
- **Go は Go のままです。** goesm は、Go のツールチェーンの go/packages と go/types をそのままフロントエンドとして使います。同じコードに対して `go.mod`、`go.work`、`gopls`、`go vet`、`go test` がそのまま使えます。goesm 独自の構文はありません。独自のディレクティブは `//goesm:import` だけで、JavaScript を呼ぶコードだけが使います。標準ライブラリも Go 自身のソースからコンパイルします。
- **ネイティブ Go と突き合わせて検証しています。** すべてのフィクスチャの結果を `go run` の結果と比較しています。Go 自身のテストスイートである `$GOROOT/test` も goesm で実行しています。
