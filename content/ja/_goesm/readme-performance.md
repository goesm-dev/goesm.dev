<!-- Synced by cmd/syncdocs. Edit the source instead. -->

[bench/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.3/bench) では、同じ Go のカーネルを goesm、[GopherJS](https://github.com/gopherjs/gopherjs)、Go 公式の `GOOS=js GOARCH=wasm`、[TinyGo](https://tinygo.org) の wasm ターゲットでコンパイルし、Node.js、Bun、Chromium で実行します。すべての結果をネイティブ Go と照合し、起動時間と出力サイズも計測します。ネイティブ Go と、同じカーネルを JavaScript で手書きしたものを基準として載せています。各カーネルの内容、出力のビルド方法と呼び出し方、公平性についての注意は [bench/README.ja.md](/ja/reference/benchmark/) にあります。ランタイムごとの全数値は [bench/results/results.md](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.3/bench/results/results.md) にあります。

計測方法は次のとおりです。

- どの実装についても、事前にビルドした出力を読み込んでから計測します。ビルドと TypeScript から JavaScript への変換は事前に済ませるので、どの計測値にも含まれません。出力の読み込み、wasm のコンパイルとインスタンス化にかかる時間は、カーネルの時間には含めず、起動時間に含めます。
- ネイティブ Go については、`go build` で作ったバイナリの中で、同じ計測ループを Go で実行します。
- 手書き JS については、[bench/js/handwritten.mjs](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.3/bench/js/handwritten.mjs) を ES モジュールとしてそのまま読み込みます。
- goesm については、`goesm build -minify` が出力した TypeScript を、goesm に組み込まれた esbuild で 1 つの ES モジュール `kernels.js` にバンドルし、それを読み込みます。
- GopherJS については、`gopherjs build -m` が出力したスクリプトを読み込みます。Go wasm と TinyGo wasm については、`.wasm` ファイルと `wasm_exec.js` を読み込みます。
- 各カーネルは JavaScript から呼び出します。ウォームアップとして 300 ms 以上かつ 3 回以上実行したあと、10 回以上かつ 1000 ms 以上計測し、その中央値を結果とします。1 回の実行が遅いカーネルでは、10 回に届かなくても、3 回以上かつ合計 10 秒以上計測した時点で計測を終えます。
- Add、Upper、Handle では、JavaScript から関数を 1 回呼び出すのにかかる時間を測ります。goesm と wasm の時間には、JS の文字列と Go の文字列を相互に変換する時間が含まれます。この変換は、呼び出し側が実際に負担するコストだからです。

<!-- bench:start -->
- Intel(R) Xeon(R) Processor @ 2.80GHz (4 threads), linux 6.18.44-fc-v70
- goesm f2910eb, go version go1.27.1 linux/amd64
- GopherJS 1.21.0+go1.21.13
- tinygo version 0.42.0 linux/amd64 (using go version go1.27.1 and LLVM version 22.1.4)
- Node.js v26.10.0, Bun 1.4.2, Chromium 141.0.7390.37
- 2026-10-05 に計測しました。各カーネルについて、300 ms 以上ウォームアップしたあと 10 回以上かつ 1000 ms 以上計測し、その中央値を結果としています。1 回の実行が遅いカーネルは、3 回以上かつ合計 10 秒以上で計測を終えています。

![ネイティブ Go に対する遅さ](/repo/goesm/bench/results/charts/slowdown.ja.svg)

![合計時間](/repo/goesm/bench/results/charts/total.ja.svg)

次の表は、ネイティブ Go に対する遅さを、カーネルごとの時間の比の幾何平均で示します。値が小さいほど速いことを表します。

| ランタイム | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 0.97× | **1.14×** | 11.37× | 2.78× | 1.73× |
| Bun 1.4.2 | 1.38× | **1.15×** | 10.59× | 3.05× | 1.90× |
| Chromium 141.0.7390.37 | 0.92× | **1.06×** | 9.61× | 3.05× | 1.79× |

次の表は、全カーネルを 1 回ずつ実行した合計時間を ms で示します。合計は各カーネルの中央値の和で、呼び出し系のカーネルはループ全体の時間を足しています。* の付いた値は、その実装にないカーネルを除いた合計です。値が小さいほど速いことを表します。

| ランタイム | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 407* | **569** | 15305 | 1884 | 1186 |
| Bun 1.4.2 | 1335* | **538** | 15170 | 1878 | 1338 |
| Chromium 141.0.7390.37 | 367* | **503** | 12336 | 2127 | 1269 |

次の表は、各実装がネイティブ Go と手書き JS に対して最も遅くなるカーネルと、その時間の比を示します。対象には次の段落の崖のカーネルも含みます。

| ランタイム | | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | ネイティブ Go 比で最も遅いカーネル | 55.0× RSASign | 327× RSASign | **8.12× Upper** | 8.61× Upper |
| Node.js v26.10.0 | 手書き JS 比で最も遅いカーネル | 88.3× RSASign | 10672× Add | **21.0× Upper** | 22.3× Upper |
| Bun 1.4.2 | ネイティブ Go 比で最も遅いカーネル | 141× Rand64 | 631× RSASign | **6.13× RSASign** | 12.0× RSASign |
| Bun 1.4.2 | 手書き JS 比で最も遅いカーネル | 98.2× RSASign | 3598× Add | **14.0× Upper** | 27.9× Pull |
| Chromium 141.0.7390.37 | ネイティブ Go 比で最も遅いカーネル | 57.5× RSASign | 277× RSASign | **11.6× Upper** | 11.8× Upper |
| Chromium 141.0.7390.37 | 手書き JS 比で最も遅いカーネル | 129× RSASign | 8009× Add | **31.4× Upper** | 31.9× Upper |

次の表は、Node.js v26.10.0 での 1 回あたりの時間の中央値を ms で示します。ns/回 と書いた行だけは、JS から関数を 1 回呼び出すのにかかる時間を ns で示します。値が小さいほど速いことを表します。

| カーネル | 対象 | ネイティブ Go | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Fib | 再帰呼び出し、int 演算 | 4.80 | 12.5 | 13.1 | 13.1 | 18.3 | **5.43** |
| Sieve | []bool、密なループ | 10.8 | 13.7 | 17.5 | 115 | 18.5 | **12.6** |
| Mandelbrot | float64 のループ | 16.5 | 16.1 | 16.0 | 16.0 | 16.6 | **15.0** |
| NBody | ポインタ経由の構造体の float64 フィールド | 11.1 | 15.1 | 18.6 | 1210 | 15.1 | **11.8** |
| FNV32 | uint32 の乗算と xor | 12.0 | 11.9 | 12.0 | 28.8 | 22.5 | **11.4** |
| FNV64 | uint64 の乗算と xor | 11.5 | 57.6 | 21.4 | 457 | 23.6 | **11.2** |
| BinaryTrees | メモリ確保、GC | 90.3 | 58.3 | **76.1** | 91.9 | 407 | 210 |
| Interfaces | インターフェースのメソッド呼び出し | 43.1 | 15.1 | **31.8** | 56.8 | 121 | 42.3 |
| MapInt | map[int]int の挿入・検索・削除 | 63.9 | 43.6 | 80.7 | **79.2** | 131 | 184 |
| MapString | map[string]int での集計 | 12.5 | 28.7 | **24.4** | 201 | 37.3 | 84.9 |
| Strings | strings.Builder、strconv、Split、Join | 12.2 | 19.9 | 21.8 | 803 | 42.8 | **17.8** |
| Sort | sort.Ints、sort.Strings | 40.8 | 70.0 | 70.8 | 225 | 142 | **53.4** |
| JSON | encoding/json の Marshal + Unmarshal | 13.3 | 5.11 | **5.98** | 900 | 35.8 | 33.4 |
| Sprintf | fmt.Sprintf | 24.4 | 11.6 | **16.5** | 2020 | 77.3 | 58.1 |
| Channels | goroutine、バッファなしチャネル | 125 | — | 91.0 | 512 | 330 | **64.5** |
| Add, ns/回 | JS からの呼び出し: 数値 2 つを渡して 1 つ受け取る | — | 0.59 | **0.59** | 6314 | 7.02 | 2.14 |
| Upper, ns/回 | JS からの呼び出し: strings.ToUpper、文字列を渡して受け取る | 186 | 72.0 | **198** | 21159 | 1513 | 1604 |
| Handle, ns/回 | JS からの呼び出し: JSON のリクエストハンドラ、文字列を渡して受け取る | 4852 | 2010 | **3095** | 582772 | 29347 | 20886 |
| **全カーネルを 1 回ずつ実行した合計 ms** | | 560* | 407* | **569** | 15305 | 1884 | 1186 |
| **ネイティブ Go 比の幾何平均** | | 1.00× | 0.97× | **1.14×** | 11.37× | 2.78× | 1.73× |

次の表は、崖のカーネルの Node.js v26.10.0 での時間を ms で示します。崖のカーネルとは、JS へのコンパイラがネイティブ Go や手書き JS より極端に遅くしやすい Go の書き方を測るカーネルです。上の幾何平均と合計には含めていません。

| カーネル | 対象 | ネイティブ Go | 手書き JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Parallel | 4 つの goroutine に分けた CPU 処理 | 8.72 | 42.5 | 33.1 | 60.3 | 47.6 | **12.2** |
| Rand64 | 構造体フィールドでの uint64 演算 | 2.80 | 24.0 | 134 | 177 | 3.63 | **2.25** |
| MaybeBlocking | 別の実装がブロックするインターフェースの呼び出し | 1.63 | 0.97 | 4.90 | 2.60 | 2.62 | **0.47** |
| Pull | iter.Pull | 4.42 | 1.17 | 70.7 | — | **12.9** | 23.6 |
| RSASign | crypto/rsa の 2048 ビット署名、math/big | 6.09 | 3.80 | 335 | 1992 | **32.7** | 51.3 |

次の表は起動時間を ms で示します。起動時間は、出力を読み込み始めてから関数を呼べるようになるまでの時間です。

| ランタイム | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Node.js | 220 | 217 | 82.0 | **64.8** |
| Bun | 496 | 474 | 115 | **50.1** |
| Chromium | 221 | 196 | 145 | **60.7** |

次の表は出力サイズを示します。サイズには、カーネル一式と、カーネルが使う標準ライブラリが含まれます。

| | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| ファイル | `kernels.js` | `bench.js` | `bench.wasm` + `wasm_exec.js` | `bench.wasm` + `wasm_exec.js` |
| 非圧縮 | **1221 KiB** | 1633 KiB | 5070 KiB | 1338 KiB |
| gzip -9 | 349 KiB | **331 KiB** | 1399 KiB | 504 KiB |
| brotli -11 | 277 KiB | **242 KiB** | 1028 KiB | 363 KiB |
<!-- bench:end -->

この数値から読み取れる goesm の現状は次のとおりです。

- **合計時間は、どのランタイムでも 4 つの変換方式の中で goesm が最も短くなっています。** 全カーネルを 1 回ずつ実行すると、goesm は 0.50〜0.57 秒、TinyGo は 1.19〜1.34 秒、Go wasm は 1.88〜2.13 秒、GopherJS は 12.3〜15.3 秒かかります。手書き JS は、Channels を除いて 0.37〜1.34 秒です。ネイティブ Go に対する遅さの幾何平均は、goesm が 1.06〜1.15 倍、手書き JS が 0.92〜1.38 倍、TinyGo が 1.73〜1.90 倍、Go wasm が 2.78〜3.05 倍、GopherJS が 9.61〜11.37 倍です。Bun では goesm の幾何平均が手書き JS を下回ります。goesm の値は、最適化を始める前は 4.2〜4.6 倍でした。今回の計測は前回と異なるマシンで行ったため、時間を以前の版の表と比べることはできません。
- **幾何平均は、ネイティブ Go より遅いカーネルと速いカーネルが打ち消し合った値です。** Node.js では、手書き JS は Fib で 2.6 倍、FNV64 で 5.0 倍、MapString で 2.3 倍の時間がかかる一方、JSON では 0.38 倍、Sprintf では 0.48 倍、Upper では 0.39 倍の時間で済みます。これは、JSON と文字列の処理では JS エンジンの組み込み関数が Go の標準ライブラリより速いためです。goesm も同じ傾向を示します。goesm は、数値計算のカーネルではネイティブ Go の 0.93〜2.7 倍の時間がかかり、Sprintf、JSON、BinaryTrees、Interfaces、JSON ハンドラではネイティブ Go より短い時間で終わります。
- **JavaScript から呼ぶコストは wasm より小さくなっています。** goesm の関数は JS の関数そのものなので、`Add` の 1 回の呼び出しは手書き JS と同じく 0.59〜0.88 ns で終わります。wasm の export を直接呼ぶ場合、TinyGo は 2.1〜4.9 ns、Go wasm は 6.8〜9.7 ns かかります。Go wasm を `syscall/js` 経由で呼ぶと、1 回にマイクロ秒単位の時間がかかります。GopherJS は 1 回に 3.7〜6.3 µs かかります。文字列を渡して受け取る `Upper` の 1 回の呼び出しは、goesm が 198 ns、wasm が 1.1〜2.2 µs、手書き JS が 69〜81 ns です。JSON のリクエストハンドラの 1 回の呼び出しは、goesm が 2.6〜3.1 µs、TinyGo が 21〜36 µs、Go wasm が 27〜29 µs、手書き JS が 1.5〜2.2 µs です。
- **手書き JS に並んだカーネルと、まだ差が大きいカーネルがあります。** goesm は、カーネルごとに手書き JS とネイティブ Go の速い方に近づくコードを出力します。たとえば、ループで使う uint64 のローカル変数は 32 ビットの整数 2 つで持ち、構造体の `json.Marshal` はその型のために生成したエンコーダを通します。FNV64 はネイティブ Go の 1.8〜2.3 倍で、BigInt を使う手書き JS の 3.4〜78 倍より速くなっています。Node.js では、Fib、Mandelbrot、FNV32、FNV64、MapString、Sort で、goesm は手書き JS との差が 5% 以内か、それより速くなっています。差が大きいのは、Upper の 2.5〜2.9 倍、Interfaces の 1.3〜2.1 倍、JSON ハンドラの 1.2〜1.9 倍、JSON の 1.2〜1.6 倍、Node.js での MapInt の最大 1.9 倍です。
- **崖のカーネルは、Go と JavaScript の違いが表れる箇所を示します。** 崖のカーネルは、どれも Go と JavaScript の一方が速く処理し、もう一方では遅くなりうる処理です。goesm が最も遅いのは RSASign で、ネイティブ Go の 55〜68 倍、WebCrypto の 88〜129 倍の時間がかかります。math/big は 32 ビットのワードの乗算を float64 の演算で行います。BigInt は演算が定数時間で終わらないため使っていません。Rand64 は uint64 を構造体のフィールドに持ち、その値が BigInt のまま残るため、ネイティブ Go の 45〜141 倍の時間がかかります。Bun では BigInt を使う手書き JS と同等です。Pull は iter.Pull がイテレータを goroutine で動かすため、ネイティブ Go の 13〜17 倍、JS のジェネレータの 40〜95 倍の時間がかかります。Parallel は JavaScript では 1 つのスレッドしか使えず、処理を順に実行する手書き JS とほぼ同じ時間がかかります。ブロックする実装を持つインターフェースを呼ぶ MaybeBlocking は、ネイティブ Go の 2.0〜3.6 倍です。このリリースで解析を直す前は、その約 13 倍の時間がかかっていました。uint64 のフィールドと iter.Pull の崖は今後対応する予定です。RSA と Parallel の崖は、JavaScript が提供する機能から生じています。
- **起動時間は wasm より長く、出力サイズは GopherJS に近くなっています。** RSASign のため、カーネルには crypto/rsa と math/big が含まれるようになりました。goesm の出力は、起動に 220〜496 ms かかります。GopherJS は 196〜474 ms、Go wasm は 82〜145 ms、TinyGo は 50〜65 ms です。brotli で圧縮したサイズは、goesm が 277 KiB、GopherJS が 242 KiB、TinyGo が 363 KiB、Go wasm が 1028 KiB です。
- **残る差の原因は、値の表現にあります。** Interfaces では、インターフェースに入れた構造体ごとに箱を 1 つ余分に確保します。Upper と JSON ハンドラでは、JS の文字列を Go のバイト文字列として受け取るときに、ASCII かどうかを調べる走査を呼び出しごとに行います。JSON では、Go の構造体へのデコードが、`JSON.parse` で素のオブジェクトを作る手書き JS より重くなります。
