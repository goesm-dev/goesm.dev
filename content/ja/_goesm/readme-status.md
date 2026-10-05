<!-- Synced by cmd/syncdocs. Edit the source instead. -->

PoC は Go 言語の大部分と、標準ライブラリのかなりの部分をコンパイルできます。

- **言語**: 関数とクロージャ、構造体、配列、スライス、マップ、ポインタ、インターフェース、型 switch、ジェネリクスと Go 1.27 のジェネリックメソッド、メソッド値、`defer` / `panic` / `recover`、goroutine、チャネル、`select`、整数と関数に対する range、`goto`、ラベル付き文、パッケージの初期化順序。
- **標準ライブラリ**: Go のソースからコンパイルします。`strings`、`strconv`、`unicode`、`sort`、`slices`、`maps`、`errors`、`math`、`math/bits`、`fmt`、`reflect`、`encoding/json`、`sync`、`time`、`os` の標準入出力などが動きます。`time` はホストのタイマーの上で動きます。goesm がまだ変換できない関数は、呼ぶと panic するスタブになります。スタブの一覧は `goesm build -v` で表示できます。
- **compile-time instrumentation**: OpenTelemetry の [otelc](https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation) が計装したプログラムを、`goesm build -toolexec "otelc toolexec"` でビルドできます。ビルドしたプログラムは、ネイティブのビルドと同じ span を console や OTLP/HTTP の collector に出力します。package 間の `//go:linkname` と `//go:embed` も動きます。詳しくは [docs/otelc.ja.md](/ja/reference/otelc/) を参照してください。
- **Go のテストスイート**: `$GOROOT/test` の実行可能なテスト 961 件のうち 898 件が、ネイティブ Go と同じ出力になります。詳しくは [docs/conformance.ja.md](/ja/reference/conformance/) を参照してください。
- **ユースケース**: cobra を含む CLI、ビルドツール、サーバーサイドレンダリング、Node.js・Bun・Deno 上の `http.ListenAndServe` による HTTP と Connect のサーバー、Cloudflare Workers の fetch ハンドラとして動く同じ `http.Handler`、Connect のクライアント、DOM 操作、Go を呼ぶ React・Preact・Next.js のアプリに対応しています。どこで何に対応し、どう確認しているか、何が未対応かは [docs/use-cases.ja.md](/ja/reference/use-cases/) にまとめています。

まだできないことは次のとおりです。詳細は [ARCHITECTURE.ja.md §11](/ja/reference/architecture/) にあります。

- `int` と `uint` は JS の number です。2^53 未満では正確ですが、64 ビットのオーバーフローで折り返しません。`int64` と `uint64` は BigInt で表すので正確です。
- JavaScript から Go を呼ぶときの ABI はまだありません。Go の文字列とスライスはランタイムのオブジェクトなので、呼び出し側が変換する必要があります。変換には、各モジュールが再 export しているランタイムの `rt.fromJSString`、`rt.sliceLit`、`rt.toArray` などを使います。ブロックしうる関数は Promise を返します。
- goroutine ごとの `recover` の状態、ホストに保留中の処理があるときのデッドロック検出、DOM バインディングは、まだ実装していません。
