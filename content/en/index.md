---
layout: home
description: goesm compiles Go packages into native ES modules for Vite, Bun, Node.js and browsers. Plain JavaScript with TypeScript types, no WebAssembly.
hero:
  name: goesm
  text: Go packages as native ES modules
  tagline: Plain JavaScript with TypeScript types. No WebAssembly, no glue code. Your bundler imports Go like any other module.
  image: /goesm-hero.webp
  actions:
    - theme: brand
      text: Get started
      link: /guide/getting-started
    - theme: alt
      text: What is goesm?
      link: /guide/
    - theme: alt
      text: GitHub
      link: https://github.com/goesm-dev/goesm
features:
  - icon: ⚡
    title: Go functions are JS functions
    details: No wasm instance, no wasm_exec.js, no marshalling. Calling an exported Go function costs what calling any JS function costs.
  - icon: 📦
    title: One ES module per package
    details: A tree of TypeScript modules with relative .ts imports. Vite, Rolldown, esbuild and Bun tree-shake, split and minify it, with source maps back to .go.
  - icon: 🧪
    title: Checked against native Go
    details: Every fixture runs under native Go and as ESM, and must agree. 898 of the 961 runnable tests of Go's own test suite pass.
    link: /reference/conformance
  - icon: 🛠️
    title: Go stays Go
    details: go/packages and go/types are the frontend. go.mod, gopls, go vet and go test keep working on the same code.
  - icon: 🧩
    title: Vue and Astro
    details: With gosfc, real Go runs in the <script setup lang="go"> of Vue components and in the ---go frontmatter of .astro files.
    link: /guide/vue
  - icon: 🔭
    title: Zero-code instrumentation
    details: goesm build -toolexec "otelc toolexec" emits the same OpenTelemetry spans a native build does.
    link: /reference/otelc
---

## Write Go, import ESM

```go
package cart

type Item struct {
	Name     string
	Price    int
	Quantity int
}

func Total(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Price * item.Quantity
	}
	return total
}
```

```sh
go tool goesm emit-ts ./cart   # goesm-ts/example.com/app/cart.ts + the runtime
```

```ts
import { Item, Total, $runtime as rt } from "./goesm-ts/example.com/app/cart.ts";

Total(rt.sliceLit([new Item(rt.fromJSString("apple"), 120, 3)])); // 360
```
