<!-- Synced by cmd/syncdocs. Edit the source instead. -->

[bench/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.4/bench) では、同じ Go のカーネルを goesm、[GopherJS](https://github.com/gopherjs/gopherjs)、Go 公式の `GOOS=js GOARCH=wasm`、[TinyGo](https://tinygo.org) の wasm ターゲットでコンパイルし、Node.js、Bun、Chromium で実行します。すべての結果をネイティブ Go と照合し、起動時間と出力サイズも計測します。ネイティブ Go と、同じカーネルを JavaScript で手書きしたものを基準として載せています。各カーネルの内容、出力のビルド方法と呼び出し方、公平性についての注意は [bench/README.ja.md](/ja/reference/benchmark/) にあります。ランタイムごとの全数値は [bench/results/results.md](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.4/bench/results/results.md) にあります。

計測方法は次のとおりです。

- どの実装についても、事前にビルドした出力を読み込んでから計測します。ビルドと TypeScript から JavaScript への変換は事前に済ませるので、どの計測値にも含まれません。出力の読み込み、wasm のコンパイルとインスタンス化にかかる時間は、カーネルの時間には含めず、起動時間に含めます。
- ネイティブ Go については、`go build` で作ったバイナリの中で、同じ計測ループを Go で実行します。
- 手書き JS については、[bench/js/handwritten.mjs](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.4/bench/js/handwritten.mjs) を ES モジュールとしてそのまま読み込みます。
- goesm については、`goesm build -minify` が出力した TypeScript を、goesm に組み込まれた esbuild で 1 つの ES モジュール `kernels.js` にバンドルし、それを読み込みます。
- GopherJS については、`gopherjs build -m` が出力したスクリプトを読み込みます。Go wasm と TinyGo wasm については、`.wasm` ファイルと `wasm_exec.js` を読み込みます。
- 各カーネルは JavaScript から呼び出します。ウォームアップとして 300 ms 以上かつ 3 回以上実行したあと、10 回以上かつ 1000 ms 以上計測し、その中央値を結果とします。1 回の実行が遅いカーネルでは、10 回に届かなくても、3 回以上かつ合計 10 秒以上計測した時点で計測を終えます。
- Add、Upper、Handle では、JavaScript から関数を 1 回呼び出すのにかかる時間を測ります。goesm と wasm の時間には、JS の文字列と Go の文字列を相互に変換する時間が含まれます。この変換は、呼び出し側が実際に負担するコストだからです。

<!-- bench:start -->
- Intel(R) Xeon(R) Processor @ 2.10GHz (4 threads), linux 6.18.44-fc-v70
- goesm v0.0.1-beta.1-33-g9297b61, go version go1.27.1 linux/amd64
- GopherJS 1.21.0+go1.21.13
- tinygo version 0.42.0 linux/amd64 (using go version go1.27.1 and LLVM version 22.1.4)
- Node.js v26.10.0, Bun 1.4.2, Chromium 141.0.7390.37
- 2026-10-06 に計測しました。各カーネルについて、300 ms 以上ウォームアップしたあと 10 回以上かつ 1000 ms 以上計測し、その中央値を結果としています。1 回の実行が遅いカーネルは、3 回以上かつ合計 10 秒以上で計測を終えています。

![ネイティブ Go に対する遅さ](/repo/goesm/bench/results/charts/slowdown.ja.svg)

![合計時間](/repo/goesm/bench/results/charts/total.ja.svg)

次の表は、ネイティブ Go に対する遅さを、カーネルごとの時間の比の幾何平均で示します。値が小さいほど速いことを表します。

| ランタイム | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 1.01× | **1.05×** | 10.83× | 2.49× | 1.67× |
| Bun 1.4.2 | 1.15× | **1.01×** | 9.27× | 2.31× | 1.79× |
| Chromium 141.0.7390.37 | 0.91× | **0.99×** | 8.64× | 2.57× | 1.71× |

次の表は、全カーネルを 1 回ずつ実行した合計時間を ms で示します。合計は各カーネルの中央値の和で、呼び出し系のカーネルはループ全体の時間を足しています。* の付いた値は、その実装にないカーネルを除いた合計です。値が小さいほど速いことを表します。

| ランタイム | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 296* | **352** | 8831 | 1080 | 685 |
| Bun 1.4.2 | 742* | **307** | 7736 | 964 | 745 |
| Chromium 141.0.7390.37 | 238* | **302** | 7212 | 1165 | 738 |

次の表は、各実装がネイティブ Go と手書き JS に対して最も遅くなるカーネルと、その時間の比を示します。対象には次の段落の崖のカーネルも含みます。

| ランタイム | | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | ネイティブ Go 比で最も遅いカーネル | 66.2× RSASign | 374× RSASign | **6.55× Upper** | 8.76× RSASign |
| Node.js v26.10.0 | 手書き JS 比で最も遅いカーネル | 122× RSASign | 5120× Add | 20.2× Upper | **16.1× RSASign** |
| Bun 1.4.2 | ネイティブ Go 比で最も遅いカーネル | 100× Rand64 | 663× RSASign | **7.07× RSASign** | 9.91× RSASign |
| Bun 1.4.2 | 手書き JS 比で最も遅いカーネル | 107× RSASign | 3324× Add | **18.1× Pull** | 23.2× Pull |
| Chromium 141.0.7390.37 | ネイティブ Go 比で最も遅いカーネル | 67.9× RSASign | 308× RSASign | 10.1× Upper | **8.95× Upper** |
| Chromium 141.0.7390.37 | 手書き JS 比で最も遅いカーネル | 106× RSASign | 3887× Add | 33.6× Upper | **29.8× Upper** |

次の表は、Node.js v26.10.0 での 1 回あたりの時間の中央値を ms で示します。ns/回 と書いた行だけは、JS から関数を 1 回呼び出すのにかかる時間を ns で示します。値が小さいほど速いことを表します。

| カーネル | 対象 | ネイティブ Go | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Fib | 再帰呼び出し、int 演算 | 4.07 | 7.80 | 7.53 | 7.42 | 11.6 | **3.88** |
| Sieve | []bool、密なループ | 6.93 | 10.2 | 10.7 | 68.2 | 11.4 | **9.37** |
| Mandelbrot | float64 のループ | 10.1 | 11.3 | 11.3 | 11.4 | 11.2 | **10.3** |
| NBody | ポインタ経由の構造体の float64 フィールド | 8.14 | 7.38 | 7.81 | 717 | 8.49 | **6.04** |
| FNV32 | uint32 の乗算と xor | 11.2 | 11.2 | 11.5 | 19.3 | 11.6 | **11.0** |
| FNV64 | uint64 の乗算と xor | 11.5 | 47.0 | 16.0 | 302 | 13.3 | **10.9** |
| BinaryTrees | メモリ確保、GC | 65.1 | 50.1 | **59.4** | 59.7 | 231 | 110 |
| Interfaces | インターフェースのメソッド呼び出し | 17.4 | 13.5 | **15.9** | 41.2 | 76.4 | 24.2 |
| MapInt | map[int]int の挿入・検索・削除 | 26.3 | 27.0 | **27.9** | 46.4 | 50.6 | 107 |
| MapString | map[string]int での集計 | 7.40 | 23.9 | **17.8** | 124 | 20.6 | 64.4 |
| Strings | strings.Builder、strconv、Split、Join | 8.16 | 14.2 | 13.9 | 487 | 27.0 | **12.7** |
| Sort | sort.Ints、sort.Strings | 29.2 | 47.3 | 46.2 | 148 | 98.6 | **38.4** |
| JSON | encoding/json の Marshal + Unmarshal | 7.00 | 2.42 | **3.67** | 514 | 21.3 | 18.2 |
| Sprintf | fmt.Sprintf | 13.8 | 5.48 | **5.62** | 1308 | 49.6 | 31.8 |
| Channels | goroutine、バッファなしチャネル | 75.1 | — | 67.5 | 264 | 208 | **20.5** |
| Add, ns/回 | JS からの呼び出し: 数値 2 つを渡して 1 つ受け取る | — | 0.57 | **0.61** | 2931 | 5.08 | 1.67 |
| Upper, ns/回 | JS からの呼び出し: strings.ToUpper、文字列を渡して受け取る | 135 | 43.9 | **110** | 13395 | 885 | 682 |
| Handle, ns/回 | JS からの呼び出し: JSON のリクエストハンドラ、文字列を渡して受け取る | 3046 | 1287 | **1838** | 308221 | 14026 | 13876 |
| **全カーネルを 1 回ずつ実行した合計 ms** | | 345* | 296* | **352** | 8831 | 1080 | 685 |
| **ネイティブ Go 比の幾何平均** | | 1.00× | 1.01× | **1.05×** | 10.83× | 2.49× | 1.67× |

次の表は、崖のカーネルの Node.js v26.10.0 での時間を ms で示します。崖のカーネルとは、JS へのコンパイラがネイティブ Go や手書き JS より極端に遅くしやすい Go の書き方を測るカーネルです。上の幾何平均と合計には含めていません。

| カーネル | 対象 | ネイティブ Go | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Parallel | 4 つの goroutine に分けた CPU 処理 | 7.21 | 34.6 | 32.8 | 59.3 | 38.6 | **16.7** |
| Rand64 | 構造体フィールドでの uint64 演算 | 2.12 | 18.5 | 19.6 | 114 | 3.54 | **2.09** |
| MaybeBlocking | 別の実装がブロックするインターフェースの呼び出し | 2.12 | 0.71 | 1.15 | 1.43 | 1.92 | **0.40** |
| Pull | iter.Pull | 4.86 | 0.99 | **2.12** | — | 8.90 | 9.21 |
| RSASign | crypto/rsa の 2048 ビット署名、math/big | 3.76 | 2.05 | 249 | 1405 | **22.6** | 32.9 |

次の表は起動時間を ms で示します。起動時間は、出力を読み込み始めてから関数を呼べるようになるまでの時間です。

| ランタイム | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Node.js | 112 | 144 | 68.0 | **45.7** |
| Bun | **215** | 948 | 1380 | 559 |
| Chromium | 111 | 257 | 347 | **48.6** |

次の表は出力サイズを示します。サイズには、カーネル一式と、カーネルが使う標準ライブラリが含まれます。

| | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| ファイル | `kernels.js` | `bench.js` | `bench.wasm` + `wasm_exec.js` | `bench.wasm` + `wasm_exec.js` |
| 非圧縮 | **760 KiB** | 1633 KiB | 5070 KiB | 1338 KiB |
| gzip -9 | **221 KiB** | 331 KiB | 1399 KiB | 504 KiB |
| brotli -11 | **176 KiB** | 242 KiB | 1029 KiB | 363 KiB |
<!-- bench:end -->

この数値から読み取れる goesm の現状は次のとおりです。

- **合計時間は、どのランタイムでも 4 つの変換方式の中で goesm が最も短く、Bun では手書き JS よりも短くなっています。** 全カーネルを 1 回ずつ実行すると、goesm は 0.30〜0.35 秒、TinyGo は 0.69〜0.75 秒、Go wasm は 0.96〜1.17 秒、GopherJS は 7.2〜8.8 秒かかります。手書き JS は、Channels を除いて 0.24〜0.74 秒です。ネイティブ Go に対する遅さの幾何平均は、goesm が 0.99〜1.05 倍、手書き JS が 0.91〜1.15 倍、TinyGo が 1.67〜1.79 倍、Go wasm が 2.31〜2.57 倍、GopherJS が 8.64〜10.83 倍です。Bun では goesm の幾何平均が手書き JS を下回ります。goesm の値は、最適化を始める前は 4.2〜4.6 倍でした。
- **幾何平均は、ネイティブ Go より遅いカーネルと速いカーネルが打ち消し合った値です。** Node.js では、手書き JS は Fib で 1.9 倍、FNV64 で 4.1 倍、MapString で 3.2 倍の時間がかかる一方、JSON では 0.35 倍、Sprintf では 0.40 倍、Upper では 0.32 倍の時間で済みます。これは、JSON と文字列の処理では JS エンジンの組み込み関数が Go の標準ライブラリより速いためです。goesm も同じ傾向を示します。goesm は、数値計算のカーネルではネイティブ Go の 0.96〜1.90 倍の時間がかかります。Sprintf、JSON、Upper、JSON ハンドラ、BinaryTrees、Interfaces、Channels では、どのランタイムでもネイティブ Go より短い時間で終わります。
- **JavaScript から呼ぶコストは wasm より小さくなっています。** goesm の関数は JS の関数そのものなので、`Add` の 1 回の呼び出しは手書き JS と同じく 0.60〜0.62 ns で終わります。wasm の export を直接呼ぶ場合、TinyGo は 1.7〜2.3 ns、Go wasm は 5.1〜6.1 ns かかります。Go wasm を `syscall/js` 経由で呼ぶと、1 回にマイクロ秒単位の時間がかかります。GopherJS は 1 回に 2.3〜2.9 µs かかります。文字列を渡して受け取る `Upper` の 1 回の呼び出しは、goesm が 103〜110 ns、wasm が 0.68〜1.4 µs、手書き JS が 41〜55 ns です。JSON のリクエストハンドラの 1 回の呼び出しは、goesm が 1.5〜1.8 µs、TinyGo が 13〜16 µs、Go wasm が 13〜15 µs、手書き JS が 0.89〜1.4 µs です。
- **ほとんどのカーネルは手書き JS に並んでいます。** goesm は、カーネルごとに手書き JS とネイティブ Go の速い方に近づくコードを出力します。たとえば、ループで使う uint64 のローカル変数は 32 ビットの整数 2 つで持ちます。構造体の `json.Marshal` はその型のために生成したエンコーダを通し、構造体へのポインタはそれ自体がインターフェース値になります。FNV64 はネイティブ Go の 1.23〜1.40 倍で、BigInt を使う手書き JS の 2.3〜45 倍より速くなっています。Node.js では、17 のカーネルのうち 10 で、goesm は手書き JS との差が 5% 以内か、それより速くなっています。差が大きいのは、Upper の 1.9〜2.5 倍、JSON の 1.4〜1.6 倍、JSON ハンドラの 1.1〜1.9 倍、Interfaces の 1.1〜1.7 倍です。BinaryTrees は今回の計測では 1.2 倍ですが、単独で計ると手書き JS と同じ速さです。
- **崖のカーネルは、Go と JavaScript の違いが表れる箇所を示します。** 崖のカーネルは、どれも Go と JavaScript の一方が速く処理し、もう一方では遅くなりうる処理です。goesm が最も遅いのは RSASign で、ネイティブ Go の 66〜87 倍、WebCrypto の 106〜122 倍の時間がかかります。math/big は 32 ビットのワードの乗算を float64 の演算で行います。BigInt は演算が定数時間で終わらないため使っていません。Rand64 は uint64 を構造体のフィールドに持ち、その値は BigInt です。goesm は、そのフィールドへの代入が続く間はフィールドをローカル変数に写して計算するので、Rand64 は BigInt を使う手書き JS の 1.04〜1.11 倍の時間で終わります。ネイティブ Go と比べると、Node.js と Chromium では 5.6〜9.3 倍、BigInt の演算が遅い Bun では 100 倍です。Pull は、関数リテラルで書いたシーケンスを goroutine ではなく JS のジェネレータとして進め、ネイティブ Go の 0.37〜0.48 倍、手書き JS のジェネレータの 2.1〜5.2 倍の時間がかかります。Parallel は JavaScript では 1 つのスレッドしか使えず、処理を順に実行する手書き JS とほぼ同じ時間がかかります。ブロックする実装を持つインターフェースを呼ぶ MaybeBlocking は、ネイティブ Go の 0.26〜0.54 倍です。RSA と Parallel の崖は、JavaScript が提供する機能から生じています。
- **起動時間は TinyGo より長く、出力サイズは 4 つの中で最も小さくなっています。** RSASign のため、カーネルには crypto/rsa と math/big が含まれます。Node.js と Chromium では、goesm の出力は起動に 111 ms かかります。GopherJS は 144〜257 ms、Go wasm は 68〜347 ms、TinyGo は 46〜49 ms です。今回の計測では、Bun でどの実装も起動が遅くなりました。goesm は 215 ms、ほかの実装は 0.56〜1.4 秒です。brotli で圧縮したサイズは、goesm が 176 KiB、GopherJS が 242 KiB、TinyGo が 363 KiB、Go wasm が 1029 KiB です。
- **残る差の原因は、値の表現にあります。** JS の文字列を Go のバイト文字列として受け取るときに、ASCII かどうかを調べる走査を呼び出しごとに行います。Upper の差の大部分はこの走査です。JSON は、`JSON.parse` では表せない Go のデコード規則に従い、Go の値を組み立てます。手書き JS より遅いカーネルごとの理由と割り当てるメモリは、[bench/README.ja.md](/ja/reference/benchmark/#goesm-が手書き-js-より遅いところ) にまとめています。
