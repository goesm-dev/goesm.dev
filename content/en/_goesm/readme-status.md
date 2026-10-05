<!-- Synced by cmd/syncdocs. Edit the source instead. -->

The proof of concept compiles most of the Go language and a good part of the standard library:

- **Language:** functions and closures, structs, arrays, slices, maps, pointers, interfaces, type switches, generics (including Go 1.27 generic methods), method values, `defer` / `panic` / `recover`, goroutines, channels, `select`, range over integers and functions, `goto`, labeled statements, package initialization order.
- **Standard library**, compiled from Go source: `strings`, `strconv`, `unicode`, `sort`, `slices`, `maps`, `errors`, `math`, `math/bits`, `fmt`, `reflect`, `encoding/json`, `sync`, `time` (on the host's timers), `os` standard streams, and more. A function goesm cannot lower yet becomes a stub that panics if called; `goesm build -v` lists them.
- **Compile-time instrumentation:** `goesm build -toolexec "otelc toolexec"` builds a program instrumented by OpenTelemetry's [otelc](https://github.com/open-telemetry/opentelemetry-go-compile-instrumentation), and it emits the spans a native build does, to the console or to an OTLP/HTTP collector ([docs/otelc.md](/reference/otelc/)). `//go:linkname` between packages and `//go:embed` work too.
- **Go's test suite:** 898 of the 961 runnable tests in `$GOROOT/test` produce the output native Go does ([docs/conformance.md](/reference/conformance/)).
- **Use cases:** command-line tools (including cobra), build tools, server-side rendering, HTTP and Connect servers with `http.ListenAndServe` under Node.js, Bun and Deno, the same `http.Handler` as a Cloudflare Workers fetch handler, Connect clients, DOM code, and React, Preact or Next.js apps calling Go. [docs/use-cases.md](/reference/use-cases/) says what is supported where, how each is checked, and what is not supported yet.

Not there yet (details in [ARCHITECTURE.md §11](/reference/architecture/#11-implemented--not-implemented--differences-from-native-go)):

- `int` and `uint` are JS numbers: exact below 2^53, but they do not wrap on 64-bit overflow. `int64` and `uint64` are exact (BigInt).
- Goroutine-local `recover` state, and deadlock detection while the host has pending work.
