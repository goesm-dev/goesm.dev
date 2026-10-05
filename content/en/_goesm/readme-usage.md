<!-- Synced by cmd/syncdocs. Edit the source instead. -->

### `goesm emit-ts`: a TypeScript tree for your bundler

Inputs are ordinary Go package patterns in an ordinary Go module:

```sh
cd testdata/example
go run ../../cmd/goesm emit-ts -o goesm-ts ./main   # -o defaults to goesm-ts
```

```text
goesm-ts/
├── example.com/app/main.ts    import * as mathx from "./mathx.ts"
├── example.com/app/mathx.ts   import * as $rt from "../../@goesm/runtime/index.ts"
└── @goesm/runtime/
    ├── index.ts
    └── ...                    the other runtime files
```

The module of Go package `p` is `<dir>/<p>.ts`, standard library packages included (`strings.ts`, `internal/bytealg.ts`), and the runtime is `<dir>/@goesm/runtime/` (a Go import path cannot start with `@`). Modules import each other with relative specifiers ending in `.ts`, so no resolver, plugin or bundler configuration is needed. Import the entry package from your own code:

```ts
// index.ts in a Vite project, or run directly: bun index.ts / node index.ts (Node.js 22.18+)
import { Result } from "./goesm-ts/example.com/app/main.ts";
console.log(Result()); // 3
```

To type-check code that imports the tree with `tsc`, enable `allowImportingTsExtensions`. [docs/example-output.md](/reference/example-output/) shows the generated TypeScript and JavaScript.

### `goesm build`: a bundled ES module

`goesm build` writes the same tree to a temporary directory and bundles it with esbuild's Go API, for when there is no bundler (and for goesm's own tests):

```sh
go run ../../cmd/goesm build ./main            # dist/main.js (+ .js.map pointing at the .go files)
go run ../../cmd/goesm build -minify ./main
go run ../../cmd/goesm build -split ./main     # one ES module per Go package: dist/example.com/app/main.js, ...
```

```js
import { Result } from "./dist/main.js";
Result(); // 3
```

The output is ES modules with a `.js` extension. Under Node.js, load them from a package whose `package.json` says `"type": "module"`: otherwise Node.js parses each module a second time to detect its format, which for a large bundle adds tens of milliseconds to startup.

### Rebuilds

goesm keeps the TypeScript module of every package it lowers in a cache, `goesm/modules` in the user cache directory. `GOESMCACHE` moves the cache, and `GOESMCACHE=off` disables it. A rebuild lowers only the packages whose module can have changed: the edited packages, the packages that depend on them, and the packages whose whole-program analysis results the edit changed, such as a function of a dependency that now has to await a callback. The output is the same with and without the cache. For the `site` package of goesm.dev, which has 77 packages, `emit-ts` takes 1.2 s without the cache and 0.7 s after editing one package; the Go frontend and the whole-program analysis take the rest.

### Calling Go from JavaScript

- Exported functions and types are exports of the package's module, with their Go types in TypeScript (`Total(items: $rt.S<Item>): number`).
- Numbers and booleans are JS numbers and booleans; `int64` / `uint64` are BigInts. Structs are classes with positional constructors (`new Item(name, price, quantity)`).
- Go strings are byte strings: `rt.fromJSString(s)` in, `rt.toJSString(s)` out. Slices: `rt.sliceLit([...])` in, `rt.toArray(s)` out. Multiple results come back as an array, and an `error` as a Go interface value.
- A function that may block (channel operations, `time.Sleep`, waiting on a mutex) is an `async function` and returns a Promise; the others are synchronous.

`rt` is the runtime, which every module re-exports as `$runtime`. [examples/](https://github.com/goesm-dev/goesm/tree/v0.0.1-beta.2/examples) has runnable examples (the cart, standard library use, goroutines) with the JavaScript that calls them.

### Calling JavaScript from Go

A function declared without a body imports a function of an ES module, and its Go types say how the values are converted:

```go
//goesm:import "./format.ts" formatPrice
func formatPrice(yen int, currency string) string

//goesm:import "./api.ts" fetchUser await
func fetchUser(id string) (User, error) // a Promise; an exception or rejection is the error
```

Strings, slices, maps, structs (as `encoding/json` names their fields), functions, `js.Value` and `any` cross the boundary, and a call costs a few nanoseconds more than the same call from JavaScript. [docs/js-imports.md](/reference/js-imports/) has the details; Vue components are used through [gosfc](https://github.com/goesm-dev/gosfc).

### Serving HTTP

An `http.Handler` (a `ServeMux`, a Connect service, middleware) serves requests on every host, through the host's own server:

```go
// Node.js, Bun and Deno: the host's HTTP server (node:http, Bun.serve, Deno.serve) runs it.
log.Fatal(http.ListenAndServe(":8080", api.Handler()))
```

```ts
// Cloudflare Workers (and Deno.serve, Bun.serve, service workers): a fetch handler.
import { Handler, $runtime as rt } from "./goesm-ts/example.com/app/api.ts";
export default { fetch: rt.fetchHandler(Handler()) };
```

Each request runs in its own goroutine, its body read in full first; the response is sent when the handler returns, or streams from its first flush (Server-Sent Events, Connect's server streaming). Under Workers, `os.Getenv` reads the Worker's text bindings and secrets with the `nodejs_compat` flag. The HTTP client is `fetch`. [docs/use-cases.md](/reference/use-cases/) lists what is supported where.
