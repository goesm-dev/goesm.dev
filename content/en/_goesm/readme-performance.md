<!-- Synced by cmd/syncdocs. Edit the source instead. -->

[bench/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.1/bench) runs the same Go kernels compiled by goesm, [GopherJS](https://github.com/gopherjs/gopherjs), Go's own `GOOS=js GOARCH=wasm` port and [TinyGo](https://tinygo.org)'s wasm target under Node.js, Bun and Chromium, checks every result against native Go, and also measures startup time and output size. Native Go and the same kernels written by hand in JavaScript are references. What each kernel does, how the outputs are built and called, and the fairness caveats are in [bench/README.md](/reference/benchmark/). All numbers, per runtime, are in [bench/results/results.md](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.1/bench/results/results.md).

How the numbers are measured:

- Every implementation is measured on its prebuilt output. Building and the TypeScript-to-JavaScript step happen beforehand and are in no measured time. Loading the output, including compiling and instantiating wasm, is left out of the kernel times and counted in startup time.
- Native Go runs the same measuring loop in Go, inside a binary built with `go build`.
- Hand-written JS is [bench/js/handwritten.mjs](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.1/bench/js/handwritten.mjs), imported as an ES module as it is.
- For goesm, the TypeScript that `goesm build -minify` emits is bundled into one ES module, `kernels.js`, by the esbuild built into goesm, and that module is imported.
- For GopherJS, the script `gopherjs build -m` emits is loaded. For Go wasm and TinyGo wasm, the `.wasm` file and `wasm_exec.js` are loaded.
- Each kernel is called from JavaScript. After a warm-up of at least 300 ms and at least 3 calls, it is timed for at least 10 calls and at least 1000 ms, and the median is the result. A kernel whose calls are slow stops after at least 3 calls once 10 s have passed, even if it has fewer than 10 calls.
- Add, Upper and Handle measure one call of a function from JavaScript. For goesm and wasm, that time includes converting strings between JS and Go, because the caller pays for that conversion.

<!-- bench:start -->
- Intel(R) Xeon(R) Processor @ 2.80GHz (4 threads), linux 6.18.44-fc-v70
- goesm f2910eb, go version go1.27.1 linux/amd64
- GopherJS 1.21.0+go1.21.13
- tinygo version 0.42.0 linux/amd64 (using go version go1.27.1 and LLVM version 22.1.4)
- Node.js v26.10.0, Bun 1.4.2, Chromium 141.0.7390.37
- 2026-10-05; warm-up ≥ 300 ms, then the median of ≥ 10 calls and ≥ 1000 ms per kernel, or of ≥ 3 calls once 10 s have passed for slow kernels

![Slowdown vs native Go](/repo/goesm/bench/results/charts/slowdown.svg)

![Total time](/repo/goesm/bench/results/charts/total.svg)

Slowdown against native Go, as the geometric mean of the per-kernel time ratios. Lower is better.

| Runtime | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 0.97× | **1.14×** | 11.37× | 2.78× | 1.73× |
| Bun 1.4.2 | 1.38× | **1.15×** | 10.59× | 3.05× | 1.90× |
| Chromium 141.0.7390.37 | 0.92× | **1.06×** | 9.61× | 3.05× | 1.79× |

Total ms to run every kernel once: the sum of the medians, with the calling kernels' whole loops included. A value marked * leaves out the kernels that implementation lacks. Lower is better.

| Runtime | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 407* | **569** | 15305 | 1884 | 1186 |
| Bun 1.4.2 | 1335* | **538** | 15170 | 1878 | 1338 |
| Chromium 141.0.7390.37 | 367* | **503** | 12336 | 2127 | 1269 |

The kernel on which each implementation is the slowest relative to native Go and to hand-written JS, with the time ratio. The cliff kernels below are included.

| Runtime | | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | Worst kernel vs native Go | 55.0× RSASign | 327× RSASign | **8.12× Upper** | 8.61× Upper |
| Node.js v26.10.0 | Worst kernel vs hand-written JS | 88.3× RSASign | 10672× Add | **21.0× Upper** | 22.3× Upper |
| Bun 1.4.2 | Worst kernel vs native Go | 141× Rand64 | 631× RSASign | **6.13× RSASign** | 12.0× RSASign |
| Bun 1.4.2 | Worst kernel vs hand-written JS | 98.2× RSASign | 3598× Add | **14.0× Upper** | 27.9× Pull |
| Chromium 141.0.7390.37 | Worst kernel vs native Go | 57.5× RSASign | 277× RSASign | **11.6× Upper** | 11.8× Upper |
| Chromium 141.0.7390.37 | Worst kernel vs hand-written JS | 129× RSASign | 8009× Add | **31.4× Upper** | 31.9× Upper |

Median ms per call under Node.js v26.10.0. The rows marked ns/call give the time of one call from JS in ns. Lower is better.

| Kernel | Exercises | Native Go | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Fib | recursive calls, int arithmetic | 4.80 | 12.5 | 13.1 | 13.1 | 18.3 | **5.43** |
| Sieve | []bool, tight loops | 10.8 | 13.7 | 17.5 | 115 | 18.5 | **12.6** |
| Mandelbrot | float64 loops | 16.5 | 16.1 | 16.0 | 16.0 | 16.6 | **15.0** |
| NBody | float64 struct fields via pointers | 11.1 | 15.1 | 18.6 | 1210 | 15.1 | **11.8** |
| FNV32 | uint32 multiply and xor | 12.0 | 11.9 | 12.0 | 28.8 | 22.5 | **11.4** |
| FNV64 | uint64 multiply and xor | 11.5 | 57.6 | 21.4 | 457 | 23.6 | **11.2** |
| BinaryTrees | allocation, GC | 90.3 | 58.3 | **76.1** | 91.9 | 407 | 210 |
| Interfaces | interface method calls | 43.1 | 15.1 | **31.8** | 56.8 | 121 | 42.3 |
| MapInt | map[int]int insert, lookup, delete | 63.9 | 43.6 | 80.7 | **79.2** | 131 | 184 |
| MapString | map[string]int counting | 12.5 | 28.7 | **24.4** | 201 | 37.3 | 84.9 |
| Strings | strings.Builder, strconv, Split, Join | 12.2 | 19.9 | 21.8 | 803 | 42.8 | **17.8** |
| Sort | sort.Ints, sort.Strings | 40.8 | 70.0 | 70.8 | 225 | 142 | **53.4** |
| JSON | encoding/json Marshal + Unmarshal | 13.3 | 5.11 | **5.98** | 900 | 35.8 | 33.4 |
| Sprintf | fmt.Sprintf | 24.4 | 11.6 | **16.5** | 2020 | 77.3 | 58.1 |
| Channels | goroutines, unbuffered channels | 125 | — | 91.0 | 512 | 330 | **64.5** |
| Add, ns/call | calls from JS: two numbers in, one out | — | 0.59 | **0.59** | 6314 | 7.02 | 2.14 |
| Upper, ns/call | calls from JS: strings.ToUpper, a string in and out | 186 | 72.0 | **198** | 21159 | 1513 | 1604 |
| Handle, ns/call | calls from JS: a JSON request handler, a string in and out | 4852 | 2010 | **3095** | 582772 | 29347 | 20886 |
| **Total ms, every kernel once** | | 560* | 407* | **569** | 15305 | 1884 | 1186 |
| **Geometric mean vs native Go** | | 1.00× | 0.97× | **1.14×** | 11.37× | 2.78× | 1.73× |

The cliff kernels, in ms under Node.js v26.10.0: ways of writing Go that a compiler to JS can make far slower than native Go or than the JS one would write. The geometric means and totals above leave them out.

| Kernel | Exercises | Native Go | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Parallel | CPU work split over 4 goroutines | 8.72 | 42.5 | 33.1 | 60.3 | 47.6 | **12.2** |
| Rand64 | uint64 arithmetic on a struct field | 2.80 | 24.0 | 134 | 177 | 3.63 | **2.25** |
| MaybeBlocking | interface calls of which another implementation blocks | 1.63 | 0.97 | 4.90 | 2.60 | 2.62 | **0.47** |
| Pull | iter.Pull | 4.42 | 1.17 | 70.7 | — | **12.9** | 23.6 |
| RSASign | crypto/rsa 2048-bit signatures, math/big | 6.09 | 3.80 | 335 | 1992 | **32.7** | 51.3 |

Startup in ms, from starting to load the output to the first callable function.

| Runtime | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Node.js | 220 | 217 | 82.0 | **64.8** |
| Bun | 496 | 474 | 115 | **50.1** |
| Chromium | 221 | 196 | 145 | **60.7** |

Output size, covering all kernels and the standard library they use.

| | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Files | `kernels.js` | `bench.js` | `bench.wasm` + `wasm_exec.js` | `bench.wasm` + `wasm_exec.js` |
| Raw | **1221 KiB** | 1633 KiB | 5070 KiB | 1338 KiB |
| gzip -9 | 349 KiB | **331 KiB** | 1399 KiB | 504 KiB |
| brotli -11 | 277 KiB | **242 KiB** | 1028 KiB | 363 KiB |
<!-- bench:end -->

What these numbers say about goesm today:

- **Total time: goesm is the fastest of the four on every runtime.** Running every kernel once takes goesm 0.50–0.57 s, TinyGo 1.19–1.34 s, Go wasm 1.88–2.13 s and GopherJS 12.3–15.3 s. Hand-written JS takes 0.37–1.34 s without Channels. In geometric mean against native Go, goesm takes 1.06–1.15×, hand-written JS 0.92–1.38×, TinyGo 1.73–1.90×, Go wasm 2.78–3.05× and GopherJS 9.61–11.37×. Under Bun, goesm's geometric mean is below hand-written JS's. Before optimization started, goesm's was 4.2–4.6×. This run used a different machine from the previous one, so its times are not comparable with earlier versions of this table.
- **The geometric mean is what is left after slower and faster kernels cancel out.** Under Node.js, hand-written JS takes 2.6× native Go's time on Fib, 5.0× on FNV64 and 2.3× on MapString. It takes only 0.38× on JSON, 0.48× on Sprintf and 0.39× on Upper, because the JS engine's built-ins for JSON and strings are faster than Go's standard library. goesm shows the same pattern. It takes 0.93–2.7× native Go's time on the numeric kernels, and less time than native Go on Sprintf, JSON, BinaryTrees, Interfaces and the JSON handler.
- **Crossing from JavaScript costs less than in wasm.** A goesm function is a JS function, so a call of `Add` takes 0.59–0.88 ns, as in hand-written JS. Through plain wasm exports, TinyGo takes 2.1–4.9 ns and Go wasm 6.8–9.7 ns. Through `syscall/js`, a Go wasm call takes microseconds. GopherJS takes 3.7–6.3 µs per call. With a string in and out, a call of `Upper` takes goesm 198 ns, wasm 1.1–2.2 µs and hand-written JS 69–81 ns. A call of the JSON request handler takes goesm 2.6–3.1 µs, TinyGo 21–36 µs, Go wasm 27–29 µs and hand-written JS 1.5–2.2 µs.
- **Some kernels have caught up with hand-written JS, and some still lag.** For each kernel, goesm emits code that moves toward the faster of hand-written JS and native Go. For example, a uint64 local used in a loop is held as two 32-bit integers, and `json.Marshal` of a struct goes through an encoder generated for its type. FNV64 takes 1.8–2.3× native Go's time, faster than the 3.4–78× of hand-written JS, which uses BigInt. Under Node.js, goesm matches hand-written JS within 5% or beats it on Fib, Mandelbrot, FNV32, FNV64, MapString and Sort. The largest gaps are Upper at 2.5–2.9×, Interfaces at 1.3–2.1×, the JSON handler at 1.2–1.9×, JSON at 1.2–1.6× and MapInt at up to 1.9× under Node.js.
- **The cliff kernels show where Go and JavaScript differ.** Each cliff kernel is something one of the two does fast and the other may not. goesm is slowest on RSASign, where it takes 55–68× native Go's time and 88–129× WebCrypto's. math/big multiplies 32-bit words through float64 arithmetic, and BigInt is not used because its arithmetic does not run in constant time. Rand64 keeps a uint64 in a struct field, which stays a BigInt, so it takes 45–141× native Go's time; under Bun this matches hand-written JS with BigInt. Pull takes 13–17× native Go's time and 40–95× a JS generator's, because iter.Pull runs its iterator on a goroutine. Parallel cannot use more than one thread in JavaScript and takes about as long as hand-written JS, which runs the work sequentially. MaybeBlocking, where an interface has an implementation that blocks, takes 2.0–3.6× native Go's time; until the analysis was fixed in this release, it took about 13× longer. The uint64 field and iter.Pull cliffs are planned work; RSA and Parallel follow from what JavaScript provides.
- **Startup is slower than wasm, and size is close to GopherJS.** The kernels now include crypto/rsa and math/big for RSASign. goesm's output starts in 220–496 ms, against 196–474 ms for GopherJS, 82–145 ms for Go wasm and 50–65 ms for TinyGo. With brotli, goesm's output is 277 KiB, GopherJS's 242 KiB, TinyGo's 363 KiB and Go wasm's 1028 KiB.
- **What remains is value representation.** In Interfaces, each struct stored in an interface takes one extra box. In Upper and the JSON handler, receiving a JS string as a Go byte string takes a scan for non-ASCII characters on each call. In JSON, decoding into Go structs costs more than hand-written JS's `JSON.parse` into plain objects.
