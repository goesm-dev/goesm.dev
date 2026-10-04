<!-- Synced by cmd/syncdocs. Edit the source instead. -->

PoC は Go 言語の大部分と、標準ライブラリのかなりの部分をコンパイルできます。

- **言語**: 関数とクロージャ、構造体、配列、スライス、マップ、ポインタ、インターフェース、型 switch、ジェネリクス（Go 1.27 のジェネリックメソッドを含む）、メソッド値、`defer` / `panic` / `recover`、goroutine、チャネル、`select`、整数と関数に対する range、`goto`、ラベル付き文、パッケージの初期化順序。
- **標準ライブラリ**（Go のソースからコンパイル）: `strings`、`strconv`、`unicode`、`sort`、`slices`、`maps`、`errors`、`math`、`math/bits`、`fmt`、`reflect`、`encoding/json`、`sync`、`time`（ホストのタイマー上で動作）、`os` の標準入出力など。goesm がまだ変換できない関数は、呼ぶと panic するスタブになります。`goesm build -v` で一覧できます。
- **compile-time instrumentation**: `goesm build -toolexec "otelc toolexec"` で OpenTelemetry の [otelc](https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation) が計装したプログラムをビルドでき、ネイティブのビルドと同じ span を console や OTLP/HTTP の collector に出力します（[docs/otelc.ja.md](/ja/reference/otelc/)）。package 間の `//go:linkname` と `//go:embed` も動きます。
- **Go のテストスイート**: `$GOROOT/test` の実行可能なテスト 961 件のうち 898 件で、ネイティブ Go と同じ出力になります（[docs/conformance.ja.md](/ja/reference/conformance/)）。

まだできないこと（詳細は [ARCHITECTURE.ja.md §11](/ja/reference/architecture/)）:

- `int` と `uint` は JS の number です。2^53 未満では正確ですが、64 ビットのオーバーフローで折り返しません。`int64` と `uint64` は正確です（BigInt）。
- JS 呼び出し ABI はまだありません。Go の文字列とスライスはランタイムのオブジェクトなので、各モジュールが再 export しているランタイム（`rt.fromJSString`、`rt.sliceLit`、`rt.toArray` など）で手作業で変換します。ブロックしうる関数は Promise を返します。
- goroutine ごとの `recover` の状態、ホストに保留中の処理があるときのデッドロック検出、DOM バインディング。
