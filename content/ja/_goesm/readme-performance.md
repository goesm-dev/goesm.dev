<!-- Synced by cmd/syncdocs. Edit the source instead. -->

[bench/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.0/bench) では、同じ Go のカーネルを goesm、[GopherJS](https://github.com/gopherjs/gopherjs)、Go 公式の `GOOS=js GOARCH=wasm`、[TinyGo](https://tinygo.org) の wasm ターゲットでコンパイルし、Node.js、Bun、Chromium で実行します。すべての結果をネイティブ Go と照合し、起動時間と出力サイズも計測します。ネイティブ Go と、同じカーネルを JavaScript で手書きしたものを基準として載せています。各カーネルの内容、出力のビルド方法と呼び出し方、公平性についての注意は [bench/README.ja.md](/ja/reference/benchmark/) に、ランタイムごとの全数値は [bench/results/results.md](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.0/bench/results/results.md) にあります。

<!-- bench:start -->
- Intel(R) Xeon(R) Processor @ 2.80GHz (4 threads), linux 6.18.44-fc-v64
- goesm 7be619e, go version go1.27.1 linux/amd64
- GopherJS 1.21.0+go1.21.13
- tinygo version 0.42.0 linux/amd64 (using go version go1.27.1 and LLVM version 22.1.4)
- Node.js v26.10.0, Bun 1.4.2, Chromium 141.0.7390.37
- 2026-10-04 計測。300 ms 以上ウォームアップしたあと、カーネルごとに 10 回以上かつ 1000 ms 以上計測した中央値

![ネイティブ Go に対する遅さ](/repo/goesm/bench/results/charts/slowdown.ja.svg)

![合計時間](/repo/goesm/bench/results/charts/total.ja.svg)

ネイティブ Go に対する遅さ（各カーネルの時間の比の幾何平均。小さいほど速い）:

| ランタイム | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 1.0× | 3.0× | 11.4× | 3.0× | **1.7×** |
| Bun 1.4.2 | 1.2× | 3.8× | 9.5× | 2.7× | **2.0×** |
| Chromium 141.0.7390.37 | 0.9× | 2.6× | 9.1× | 3.4× | **1.8×** |

全カーネルを 1 回ずつ実行した合計時間（ms、中央値の和。呼び出し系カーネルはループ全体。* はないカーネルを除いた値。小さいほど速い）:

| ランタイム | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 329* | 1534 | 11622 | 1475 | **833** |
| Bun 1.4.2 | 961* | 2806 | 9998 | 1289 | **1036** |
| Chromium 141.0.7390.37 | 282* | 1288 | 9153 | 1684 | **947** |

Node.js v26.10.0 での 1 回あたりの時間（ms、中央値。小さいほど速い）:

| カーネル | 対象 | ネイティブ Go | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Fib | 再帰呼び出し、int 演算 | 4.66 | 10.1 | 9.33 | 9.60 | 14.4 | **3.98** |
| Sieve | []bool、密なループ | 8.21 | 11.1 | 45.4 | 77.4 | 14.3 | **10.9** |
| Mandelbrot | float64 のループ | 15.4 | 16.1 | 16.1 | 16.0 | 15.7 | **14.2** |
| NBody | ポインタ経由の構造体の float64 フィールド | 8.81 | 10.5 | 16.2 | 1074 | 11.6 | **9.20** |
| FNV32 | uint32 の乗算と xor | 10.9 | 11.1 | 13.3 | 26.8 | 25.3 | **11.1** |
| FNV64 | uint64 の乗算と xor | 11.2 | 47.9 | 45.6 | 335 | 25.3 | **10.8** |
| BinaryTrees | メモリ確保、GC | 66.7 | 55.1 | **54.4** | 72.8 | 280 | 108 |
| Interfaces | インターフェースのメソッド呼び出し | 21.0 | 13.2 | 35.7 | 45.2 | 105 | **30.7** |
| MapInt | map[int]int の挿入・検索・削除 | 30.5 | 24.5 | 49.4 | **46.1** | 107 | 145 |
| MapString | map[string]int での集計 | 9.60 | 24.1 | **26.6** | 135 | 32.4 | 80.3 |
| Strings | strings.Builder、strconv、Split、Join | 10.3 | 18.3 | 46.2 | 611 | 35.5 | **13.5** |
| Sort | sort.Ints、sort.Strings | 34.4 | 51.5 | 108 | 192 | 126 | **39.8** |
| JSON | encoding/json の Marshal + Unmarshal | 9.63 | 2.92 | 101 | 656 | 29.0 | **28.1** |
| Sprintf | fmt.Sprintf | 16.9 | 10.4 | 160 | 1657 | 70.1 | **46.2** |
| Channels | goroutine、バッファなしチャネル | 91.0 | — | 113 | 394 | 284 | **31.1** |
| Add (ns/回) | JS からの呼び出し: 数値 2 つを渡して 1 つ受け取る | — | 0.59 | **0.59** | 3772 | 6.30 | 2.10 |
| Upper (ns/回) | JS からの呼び出し: strings.ToUpper、文字列を渡して受け取る | 149 | 58.2 | 1358 | 15573 | 1201 | **896** |
| Handle (ns/回) | JS からの呼び出し: JSON のリクエストハンドラ、文字列を渡して受け取る | 4051 | 1613 | 55864 | 434067 | 17849 | **16108** |
| **合計 ms（全カーネルを 1 回ずつ）** | | 405* | 329* | 1534 | 11622 | 1475 | **833** |
| **ネイティブ Go 比の幾何平均** | | 1× | 1.0× | 3.0× | 11.4× | 3.0× | **1.7×** |

起動時間（出力を読み込み始めてから関数を呼べるようになるまで、ms）:

| ランタイム | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Node.js | 132 | 90.6 | 55.0 | **17.9** |
| Bun | 252 | 297 | 63.2 | **23.1** |
| Chromium | 86.8 | 79.3 | 69.5 | **31.6** |

出力サイズ（カーネル一式と、使っている標準ライブラリ）:

| | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| ファイル | `kernels.js` | `bench.js` | `bench.wasm` + `wasm_exec.js` | `bench.wasm` + `wasm_exec.js` |
| 非圧縮 | 1194 KiB | 1153 KiB | 4393 KiB | **1104 KiB** |
| gzip -9 | 309 KiB | **228 KiB** | 1211 KiB | 406 KiB |
| brotli -11 | 238 KiB | **169 KiB** | 894 KiB | 299 KiB |
<!-- bench:end -->

この数値から読み取れる goesm の現状:

- **合計時間は TinyGo が最速で、goesm と Go wasm がほぼ並び、GopherJS は大きく離れています。** 全カーネルを 1 回ずつ実行すると、goesm は 1.3〜2.8 秒、Go wasm は 1.3〜1.7 秒、TinyGo は 0.8〜1.0 秒、GopherJS は 9.2〜11.6 秒かかります。Chromium では goesm（1.3 秒）が Go wasm（1.7 秒）より速く、Node.js では同等（1.5 秒）、Bun では goesm が遅れています（2.8 秒。その 5 分の 1 以上は、Bun のエンジンで遅い BigInt を使う FNV64 です）。ネイティブ Go に対する幾何平均では、goesm が 2.6〜3.8 倍（最適化を始める前は 4.2〜4.6 倍）、Go wasm が 2.7〜3.4 倍、TinyGo が 1.7〜2.0 倍、GopherJS が 9.1〜11.4 倍です。
- **JavaScript から呼ぶコストは、数値なら goesm ではゼロですが、文字列と `encoding/json` ではまだ wasm より高くつきます。** goesm の関数は JS の関数そのものなので、`Add` のコストは手書き JS と同じ 0.6〜0.7 ns です。素の wasm export 経由では TinyGo が 2〜3 ns、Go wasm が 5〜8 ns（`syscall/js` 経由ならマイクロ秒単位）、GopherJS は 3〜4 µs です。文字列を渡して受け取る `Upper` は、goesm が Chromium と Node.js で 1 回 1.1〜1.4 µs（Bun では 3.0 µs）、wasm が 0.9〜1.7 µs です（手書き JS は 0.05〜0.07 µs）。JSON のリクエストハンドラは、goesm が Chromium と Node.js で 1 回 48〜56 µs（Bun では 93 µs）、wasm が 16〜23 µs です。リフレクションで動く goesm の `encoding/json` が最も遅い部分です。
- **goesm が強いところと弱いところ。** メモリ確保の多い BinaryTrees ではどのランタイムでも最速で、Interfaces でも Bun と Chromium で最速、map のカーネルでも一部のランタイムで最速です。Fib、Mandelbrot、FNV32、FNV64 では手書き JS と同等です。ネイティブ Go から最も離れているのは Sprintf と JSON（Node.js で約 10 倍）で、どちらもバイト文字列とリフレクションを使う標準ライブラリのコードです。`[]bool` が JS の真偽値の配列になる Sieve も離れています。起動は 87〜252 ms で wasm の 18〜70 ms より遅く、サイズは brotli で 238 KiB です（GopherJS 169 KiB、TinyGo 299 KiB、Go wasm 894 KiB）。
- **差は JavaScript ではなく goesm の変換にあります。** 手書き JS は Node.js と Chromium でネイティブ Go とほぼ同じ速さなので、goesm との差は goesm が生成するコードとランタイムにあり、引き続きそこを最適化していきます。
