<!-- Synced by cmd/syncdocs. Edit the source instead. -->

[bench/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.4/bench) runs the same Go kernels compiled by goesm, [GopherJS](https://github.com/gopherjs/gopherjs), Go's own `GOOS=js GOARCH=wasm` port and [TinyGo](https://tinygo.org)'s wasm target under Node.js, Bun and Chromium, checks every result against native Go, and also measures startup time and output size. Native Go and the same kernels written by hand in JavaScript are references. What each kernel does, how the outputs are built and called, and the fairness caveats are in [bench/README.md](/reference/benchmark/). All numbers, per runtime, are in [bench/results/results.md](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.4/bench/results/results.md).

How the numbers are measured:

- Every implementation is measured on its prebuilt output. Building and the TypeScript-to-JavaScript step happen beforehand and are in no measured time. Loading the output, including compiling and instantiating wasm, is left out of the kernel times and counted in startup time.
- Native Go runs the same measuring loop in Go, inside a binary built with `go build`.
- Hand-written JS is [bench/js/handwritten.mjs](https://github.com/goesm-dev/goesm/blob/v0.0.1-beta.4/bench/js/handwritten.mjs), imported as an ES module as it is.
- For goesm, the TypeScript that `goesm build -minify` emits is bundled into one ES module, `kernels.js`, by the esbuild built into goesm, and that module is imported.
- For GopherJS, the script `gopherjs build -m` emits is loaded. For Go wasm and TinyGo wasm, the `.wasm` file and `wasm_exec.js` are loaded.
- Each kernel is called from JavaScript. After a warm-up of at least 300 ms and at least 3 calls, it is timed for at least 10 calls and at least 1000 ms, and the median is the result. A kernel whose calls are slow stops after at least 3 calls once 10 s have passed, even if it has fewer than 10 calls.
- Add, Upper and Handle measure one call of a function from JavaScript. For goesm and wasm, that time includes converting strings between JS and Go, because the caller pays for that conversion.

<!-- bench:start -->
- Intel(R) Xeon(R) Processor @ 2.10GHz (4 threads), linux 6.18.44-fc-v70
- goesm v0.0.1-beta.1-33-g9297b61, go version go1.27.1 linux/amd64
- GopherJS 1.21.0+go1.21.13
- tinygo version 0.42.0 linux/amd64 (using go version go1.27.1 and LLVM version 22.1.4)
- Node.js v26.10.0, Bun 1.4.2, Chromium 141.0.7390.37
- 2026-10-06; warm-up ≥ 300 ms, then the median of ≥ 10 calls and ≥ 1000 ms per kernel, or of ≥ 3 calls once 10 s have passed for slow kernels

![Slowdown vs native Go](/repo/goesm/bench/results/charts/slowdown.svg)

![Total time](/repo/goesm/bench/results/charts/total.svg)

Slowdown against native Go, as the geometric mean of the per-kernel time ratios. Lower is better.

| Runtime | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 1.01× | **1.05×** | 10.83× | 2.49× | 1.67× |
| Bun 1.4.2 | 1.15× | **1.01×** | 9.27× | 2.31× | 1.79× |
| Chromium 141.0.7390.37 | 0.91× | **0.99×** | 8.64× | 2.57× | 1.71× |

Total ms to run every kernel once: the sum of the medians, with the calling kernels' whole loops included. A value marked * leaves out the kernels that implementation lacks. Lower is better.

| Runtime | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | 296* | **352** | 8831 | 1080 | 685 |
| Bun 1.4.2 | 742* | **307** | 7736 | 964 | 745 |
| Chromium 141.0.7390.37 | 238* | **302** | 7212 | 1165 | 738 |

The kernel on which each implementation is the slowest relative to native Go and to hand-written JS, with the time ratio. The cliff kernels below are included.

| Runtime | | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: |
| Node.js v26.10.0 | Worst kernel vs native Go | 66.2× RSASign | 374× RSASign | **6.55× Upper** | 8.76× RSASign |
| Node.js v26.10.0 | Worst kernel vs hand-written JS | 122× RSASign | 5120× Add | 20.2× Upper | **16.1× RSASign** |
| Bun 1.4.2 | Worst kernel vs native Go | 100× Rand64 | 663× RSASign | **7.07× RSASign** | 9.91× RSASign |
| Bun 1.4.2 | Worst kernel vs hand-written JS | 107× RSASign | 3324× Add | **18.1× Pull** | 23.2× Pull |
| Chromium 141.0.7390.37 | Worst kernel vs native Go | 67.9× RSASign | 308× RSASign | 10.1× Upper | **8.95× Upper** |
| Chromium 141.0.7390.37 | Worst kernel vs hand-written JS | 106× RSASign | 3887× Add | 33.6× Upper | **29.8× Upper** |

Median ms per call under Node.js v26.10.0. The rows marked ns/call give the time of one call from JS in ns. Lower is better.

| Kernel | Exercises | Native Go | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Fib | recursive calls, int arithmetic | 4.07 | 7.80 | 7.53 | 7.42 | 11.6 | **3.88** |
| Sieve | []bool, tight loops | 6.93 | 10.2 | 10.7 | 68.2 | 11.4 | **9.37** |
| Mandelbrot | float64 loops | 10.1 | 11.3 | 11.3 | 11.4 | 11.2 | **10.3** |
| NBody | float64 struct fields via pointers | 8.14 | 7.38 | 7.81 | 717 | 8.49 | **6.04** |
| FNV32 | uint32 multiply and xor | 11.2 | 11.2 | 11.5 | 19.3 | 11.6 | **11.0** |
| FNV64 | uint64 multiply and xor | 11.5 | 47.0 | 16.0 | 302 | 13.3 | **10.9** |
| BinaryTrees | allocation, GC | 65.1 | 50.1 | **59.4** | 59.7 | 231 | 110 |
| Interfaces | interface method calls | 17.4 | 13.5 | **15.9** | 41.2 | 76.4 | 24.2 |
| MapInt | map[int]int insert, lookup, delete | 26.3 | 27.0 | **27.9** | 46.4 | 50.6 | 107 |
| MapString | map[string]int counting | 7.40 | 23.9 | **17.8** | 124 | 20.6 | 64.4 |
| Strings | strings.Builder, strconv, Split, Join | 8.16 | 14.2 | 13.9 | 487 | 27.0 | **12.7** |
| Sort | sort.Ints, sort.Strings | 29.2 | 47.3 | 46.2 | 148 | 98.6 | **38.4** |
| JSON | encoding/json Marshal + Unmarshal | 7.00 | 2.42 | **3.67** | 514 | 21.3 | 18.2 |
| Sprintf | fmt.Sprintf | 13.8 | 5.48 | **5.62** | 1308 | 49.6 | 31.8 |
| Channels | goroutines, unbuffered channels | 75.1 | — | 67.5 | 264 | 208 | **20.5** |
| Add, ns/call | calls from JS: two numbers in, one out | — | 0.57 | **0.61** | 2931 | 5.08 | 1.67 |
| Upper, ns/call | calls from JS: strings.ToUpper, a string in and out | 135 | 43.9 | **110** | 13395 | 885 | 682 |
| Handle, ns/call | calls from JS: a JSON request handler, a string in and out | 3046 | 1287 | **1838** | 308221 | 14026 | 13876 |
| **Total ms, every kernel once** | | 345* | 296* | **352** | 8831 | 1080 | 685 |
| **Geometric mean vs native Go** | | 1.00× | 1.01× | **1.05×** | 10.83× | 2.49× | 1.67× |

The cliff kernels, in ms under Node.js v26.10.0: ways of writing Go that a compiler to JS can make far slower than native Go or than the JS one would write. The geometric means and totals above leave them out.

| Kernel | Exercises | Native Go | Hand-written JS | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Parallel | CPU work split over 4 goroutines | 7.21 | 34.6 | 32.8 | 59.3 | 38.6 | **16.7** |
| Rand64 | uint64 arithmetic on a struct field | 2.12 | 18.5 | 19.6 | 114 | 3.54 | **2.09** |
| MaybeBlocking | interface calls of which another implementation blocks | 2.12 | 0.71 | 1.15 | 1.43 | 1.92 | **0.40** |
| Pull | iter.Pull | 4.86 | 0.99 | **2.12** | — | 8.90 | 9.21 |
| RSASign | crypto/rsa 2048-bit signatures, math/big | 3.76 | 2.05 | 249 | 1405 | **22.6** | 32.9 |

Startup in ms, from starting to load the output to the first callable function.

| Runtime | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Node.js | 112 | 144 | 68.0 | **45.7** |
| Bun | **215** | 948 | 1380 | 559 |
| Chromium | 111 | 257 | 347 | **48.6** |

Output size, covering all kernels and the standard library they use.

| | goesm | GopherJS | Go wasm | TinyGo wasm |
| --- | ---: | ---: | ---: | ---: |
| Files | `kernels.js` | `bench.js` | `bench.wasm` + `wasm_exec.js` | `bench.wasm` + `wasm_exec.js` |
| Raw | **760 KiB** | 1633 KiB | 5070 KiB | 1338 KiB |
| gzip -9 | **221 KiB** | 331 KiB | 1399 KiB | 504 KiB |
| brotli -11 | **176 KiB** | 242 KiB | 1029 KiB | 363 KiB |
<!-- bench:end -->

What these numbers say about goesm today:

- **Total time: goesm is the fastest of the four on every runtime, and under Bun faster than hand-written JS.** Running every kernel once takes goesm 0.30–0.35 s, TinyGo 0.69–0.75 s, Go wasm 0.96–1.17 s and GopherJS 7.2–8.8 s. Hand-written JS takes 0.24–0.74 s without Channels. In geometric mean against native Go, goesm takes 0.99–1.05×, hand-written JS 0.91–1.15×, TinyGo 1.67–1.79×, Go wasm 2.31–2.57× and GopherJS 8.64–10.83×. Under Bun, goesm's geometric mean is below hand-written JS's. Before optimization started, goesm's was 4.2–4.6×.
- **The geometric mean is what is left after slower and faster kernels cancel out.** Under Node.js, hand-written JS takes 1.9× native Go's time on Fib, 4.1× on FNV64 and 3.2× on MapString. It takes only 0.35× on JSON, 0.40× on Sprintf and 0.32× on Upper, because the JS engine's built-ins for JSON and strings are faster than Go's standard library. goesm shows the same pattern. It takes 0.96–1.90× native Go's time on the numeric kernels. It takes less time than native Go on Sprintf, JSON, Upper, the JSON handler, BinaryTrees, Interfaces and Channels on every runtime.
- **Crossing from JavaScript costs less than in wasm.** A goesm function is a JS function, so a call of `Add` takes 0.60–0.62 ns, as in hand-written JS. Through plain wasm exports, TinyGo takes 1.7–2.3 ns and Go wasm 5.1–6.1 ns. Through `syscall/js`, a Go wasm call takes microseconds. GopherJS takes 2.3–2.9 µs per call. With a string in and out, a call of `Upper` takes goesm 103–110 ns, wasm 0.68–1.4 µs and hand-written JS 41–55 ns. A call of the JSON request handler takes goesm 1.5–1.8 µs, TinyGo 13–16 µs, Go wasm 13–15 µs and hand-written JS 0.89–1.4 µs.
- **Most kernels have caught up with hand-written JS.** For each kernel, goesm emits code that moves toward the faster of hand-written JS and native Go. For example, a uint64 local used in a loop is held as two 32-bit integers, `json.Marshal` of a struct goes through an encoder generated for its type, and a pointer to a struct is its own interface value. FNV64 takes 1.23–1.40× native Go's time, faster than the 2.3–45× of hand-written JS, which uses BigInt. Under Node.js, goesm matches hand-written JS within 5% or beats it on 10 of the 17 kernels. The largest gaps are Upper at 1.9–2.5×, JSON at 1.4–1.6×, the JSON handler at 1.1–1.9×, and Interfaces at 1.1–1.7×. BinaryTrees reads 1.2× in this run, but timed by itself it runs at hand-written JS's speed.
- **The cliff kernels show where Go and JavaScript differ.** Each cliff kernel is something one of the two does fast and the other may not. goesm is slowest on RSASign, where it takes 66–87× native Go's time and 106–122× WebCrypto's. math/big multiplies 32-bit words through float64 arithmetic, and BigInt is not used because its arithmetic does not run in constant time. Rand64 keeps a uint64 in a struct field, which is a BigInt. goesm works on a local copy of the field through a run of assignments, so Rand64 takes 1.04–1.11× hand-written JS's time with BigInt; that is 5.6–9.3× native Go's time under Node.js and Chromium, and 100× under Bun, whose BigInt arithmetic is slow. Pull steps a sequence written as a function literal as a JS generator instead of a goroutine, and takes 0.37–0.48× native Go's time and 2.1–5.2× a hand-written JS generator's. Parallel cannot use more than one thread in JavaScript and takes about as long as hand-written JS, which runs the work sequentially. MaybeBlocking, where an interface has an implementation that blocks, takes 0.26–0.54× native Go's time. RSA and Parallel follow from what JavaScript provides.
- **Startup is slower than TinyGo, and the output is the smallest of the four.** The kernels include crypto/rsa and math/big for RSASign. Under Node.js and Chromium, goesm's output starts in 111 ms, against 144–257 ms for GopherJS, 68–347 ms for Go wasm and 46–49 ms for TinyGo. Under Bun every implementation started more slowly in this run: goesm in 215 ms and the others in 0.56–1.4 s. With brotli, goesm's output is 176 KiB, GopherJS's 242 KiB, TinyGo's 363 KiB and Go wasm's 1029 KiB.
- **What remains is value representation.** Receiving a JS string as a Go byte string takes a scan for non-ASCII characters on each call, which is most of Upper's gap. JSON follows Go's decoding rules, which `JSON.parse` cannot express, and builds Go values. [bench/README.md](/reference/benchmark/#where-goesm-is-slower-than-hand-written-js) lists each kernel still slower than hand-written JS, with the reason and the memory each allocates.
