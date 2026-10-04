<!-- Synced by cmd/syncdocs. Edit the source instead. -->

[bench/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.0/bench) runs the same Go kernels compiled by goesm, [GopherJS](https://github.com/gopherjs/gopherjs), Go's own `GOOS=js GOARCH=wasm` port and [TinyGo](https://tinygo.org)'s wasm target under Node.js, Bun and Chromium, checks every result against native Go, and also measures startup time and output size. Native Go and the same kernels written by hand in JavaScript are references. What each kernel does, how the outputs are built and called, and the fairness caveats are in [bench/README.md](/reference/benchmark/); all numbers, per runtime, are in [bench/results/results.md](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.0/bench/results/results.md).

<!-- bench:start -->
- Intel(R) Xeon(R) Processor @ 2.80GHz (4 threads), linux 6.18.44-fc-v64
- goesm 7be619e, go version go1.27.1 linux/amd64
- GopherJS 1.21.0+go1.21.13
- tinygo version 0.42.0 linux/amd64 (using go version go1.27.1 and LLVM version 22.1.4)
- Node.js v26.10.0, Bun 1.4.2, Chromium 141.0.7390.37
- 2026-10-04; warm-up ≥ 300 ms, then the median of ≥ 10 calls and ≥ 1000 ms per kernel

![Slowdown vs native Go](/repo/goesm/bench/results/charts/slowdown.svg)

![Total time](/repo/goesm/bench/results/charts/total.svg)

Slowdown vs native Go (geometric mean of the per-kernel time ratios; lower is better):

| Runtime | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 1.0× | 3.0× | 11.4× | 3.0× | **1.7×** |
| Bun 1.4.2 | 1.2× | 3.8× | 9.5× | 2.7× | **2.0×** |
| Chromium 141.0.7390.37 | 0.9× | 2.6× | 9.1× | 3.4× | **1.8×** |

Total ms to run every kernel once (the sum of the medians, the calling kernels' whole loops included; * leaves out kernels the implementation lacks; lower is better):

| Runtime | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 329* | 1534 | 11622 | 1475 | **833** |
| Bun 1.4.2 | 961* | 2806 | 9998 | 1289 | **1036** |
| Chromium 141.0.7390.37 | 282* | 1288 | 9153 | 1684 | **947** |

Median ms per call under Node.js v26.10.0 (lower is better):

| Kernel | Exercises | Native Go | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Fib | recursive calls, int arithmetic | 4.66 | 10.1 | 9.33 | 9.60 | 14.4 | **3.98** |
| Sieve | []bool, tight loops | 8.21 | 11.1 | 45.4 | 77.4 | 14.3 | **10.9** |
| Mandelbrot | float64 loops | 15.4 | 16.1 | 16.1 | 16.0 | 15.7 | **14.2** |
| NBody | float64 struct fields via pointers | 8.81 | 10.5 | 16.2 | 1074 | 11.6 | **9.20** |
| FNV32 | uint32 multiply and xor | 10.9 | 11.1 | 13.3 | 26.8 | 25.3 | **11.1** |
| FNV64 | uint64 multiply and xor | 11.2 | 47.9 | 45.6 | 335 | 25.3 | **10.8** |
| BinaryTrees | allocation, GC | 66.7 | 55.1 | **54.4** | 72.8 | 280 | 108 |
| Interfaces | interface method calls | 21.0 | 13.2 | 35.7 | 45.2 | 105 | **30.7** |
| MapInt | map[int]int insert, lookup, delete | 30.5 | 24.5 | 49.4 | **46.1** | 107 | 145 |
| MapString | map[string]int counting | 9.60 | 24.1 | **26.6** | 135 | 32.4 | 80.3 |
| Strings | strings.Builder, strconv, Split, Join | 10.3 | 18.3 | 46.2 | 611 | 35.5 | **13.5** |
| Sort | sort.Ints, sort.Strings | 34.4 | 51.5 | 108 | 192 | 126 | **39.8** |
| JSON | encoding/json Marshal + Unmarshal | 9.63 | 2.92 | 101 | 656 | 29.0 | **28.1** |
| Sprintf | fmt.Sprintf | 16.9 | 10.4 | 160 | 1657 | 70.1 | **46.2** |
| Channels | goroutines, unbuffered channels | 91.0 | — | 113 | 394 | 284 | **31.1** |
| Add (ns/call) | calls from JS: two numbers in, one out | — | 0.59 | **0.59** | 3772 | 6.30 | 2.10 |
| Upper (ns/call) | calls from JS: strings.ToUpper, a string in and out | 149 | 58.2 | 1358 | 15573 | 1201 | **896** |
| Handle (ns/call) | calls from JS: a JSON request handler, a string in and out | 4051 | 1613 | 55864 | 434067 | 17849 | **16108** |
| **Total ms, every kernel once** | | 405* | 329* | 1534 | 11622 | 1475 | **833** |
| **Geometric mean vs native Go** | | 1× | 1.0× | 3.0× | 11.4× | 3.0× | **1.7×** |

Startup (ms from starting to load the output to the first callable function):

| Runtime | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Node.js | 132 | 90.6 | 55.0 | **17.9** |
| Bun | 252 | 297 | 63.2 | **23.1** |
| Chromium | 86.8 | 79.3 | 69.5 | **31.6** |

Output size (all kernels and the standard library they use):

| | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Files | `kernels.js` | `bench.js` | `bench.wasm` + `wasm_exec.js` | `bench.wasm` + `wasm_exec.js` |
| Raw | 1194 KiB | 1153 KiB | 4393 KiB | **1104 KiB** |
| gzip -9 | 309 KiB | **228 KiB** | 1211 KiB | 406 KiB |
| brotli -11 | 238 KiB | **169 KiB** | 894 KiB | 299 KiB |
<!-- bench:end -->

What these numbers say about goesm today:

- **Total time: TinyGo is fastest; goesm and Go wasm are close, and GopherJS is far behind.** Running every kernel once takes goesm 1.3–2.8 s, Go wasm 1.3–1.7 s, TinyGo 0.8–1.0 s and GopherJS 9.2–11.6 s. Under Chromium goesm (1.3 s) is ahead of Go wasm (1.7 s), under Node.js they are even (1.5 s), and under Bun goesm is behind (2.8 s, more than a fifth of it FNV64's BigInt arithmetic, which is slow in Bun's engine). In geometric mean against native Go that is 2.6–3.8× for goesm (4.2–4.6× before optimization started), 2.7–3.4× for Go wasm, 1.7–2.0× for TinyGo and 9.1–11.4× for GopherJS.
- **Crossing from JavaScript is free for numbers under goesm, but strings and `encoding/json` still cost more than in wasm.** A goesm function is a JS function, so `Add` costs what it costs in hand-written JS (0.6–0.7 ns); through plain wasm exports TinyGo takes 2–3 ns and Go wasm 5–8 ns (microseconds through `syscall/js`), and GopherJS 3–4 µs. With a string in and out (`Upper`), goesm takes 1.1–1.4 µs per call under Chromium and Node.js (3.0 µs under Bun) and wasm 0.9–1.7 µs (hand-written JS: 0.05–0.07 µs). The JSON request handler takes goesm 48–56 µs per call under Chromium and Node.js (93 µs under Bun) against 16–23 µs in wasm: goesm's `encoding/json`, which works through reflection, is its slowest part.
- **Where goesm stands out and where it lags.** It is the fastest at allocation-heavy BinaryTrees on every runtime, at Interfaces under Bun and Chromium and at the map kernels on some runtimes, and matches hand-written JS on Fib, Mandelbrot, FNV32 and FNV64. It is furthest from native Go on Sprintf and JSON (about 10× under Node.js), standard-library code on byte strings and reflection, and on Sieve, whose `[]bool` is a JS array of booleans. Its output starts in 87–252 ms against 18–70 ms for wasm, and is 238 KiB with brotli against 169 KiB for GopherJS, 299 KiB for TinyGo and 894 KiB for Go wasm.
- **The gap is in goesm's lowering, not in JavaScript.** Hand-written JS is within about 1× of native Go under Node.js and Chromium, so what separates goesm from it is code goesm generates and its runtime, which is where optimization continues.
