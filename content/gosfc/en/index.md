---
layout: home
description: gosfc runs real Go in Vue components and .astro files, the <script setup> of a .vue file and the frontmatter and <script> of an .astro file. Compiled to plain JavaScript by goesm, no WebAssembly.
hero:
  name: gosfc
  text: Go in Vue and Astro
  tagline: Write <script setup lang="go"> in a .vue file, or ---go frontmatter in an .astro file. It is real Go, compiled by goesm to plain JavaScript; Vue and Astro render it, and Vite builds it.
  actions:
    - theme: brand
      text: Get started
      link: /guide/getting-started
    - theme: alt
      text: What is gosfc?
      link: /guide/
    - theme: alt
      text: GitHub
      link: https://github.com/goesm-dev/gosfc
features:
  - icon: 🐹
    title: Real Go
    details: The Go block imports ordinary Go packages from your go.mod. gopls, go vet and go test keep working on them.
    link: /guide/go-block
  - icon: 🧩
    title: Still a Vue SFC
    details: The template, styles and the rest of the file are left to Vue. Components with lang="ts" keep working alongside.
  - icon: 🚀
    title: Go in .astro files
    details: One integration for Astro. The frontmatter of an .astro file can be Go, <script lang="go"> runs Go in the browser, and Vue components render on the server and hydrate as islands.
    link: /guide/astro
  - icon: ⚡
    title: No WebAssembly
    details: goesm compiles the Go to ES modules. The client JS of the benchmark page is 13% larger (gzip) than the same components in TypeScript.
    link: /guide/benchmark
  - icon: 🎯
    title: Errors point at your file
    details: Go's diagnostics, the source maps and the stack trace of a panic all point at lines of the .vue or .astro file.
    link: /reference/architecture
  - icon: 🔁
    title: Hot reload
    details: A template edit re-renders and keeps the Go state; a Go edit reloads the component, and so does editing a .go file it uses.
    link: /reference/architecture
---

## A component, in Go

```vue
<template>
  <button type="button" @click="Increment">Clicked {{ count }} times</button>
</template>

<script setup lang="go">
count := 0

func Increment() {
	count++
}
</script>
```

The block runs once per component instance, like Vue's `<script setup>`. Its
top-level variables and functions are what the template sees, and the page
updates when a Go function changes them.

## A page, in Go

```astro
---go
import (
	"strconv"

	cart "example.com/app/src/features/cart/pkg"
	Summary "../features/cart/Summary.vue"
)

items := []cart.Item{{Price: 120, Quantity: 3}, {Price: 80, Quantity: 1}}
total := cart.Total(items)

func Yen(n int) string {
	return "¥" + strconv.Itoa(n)
}
---

<p>Total: {Yen(total)} ({items.length} items)</p>
<Summary />
```

An `.astro` file whose frontmatter opens with `---go` is Go, too. It runs once
per render, at build time or per request, and the template sees its values and
functions. A `<script lang="go">` in the template runs Go in the browser
([Go in .astro files](./guide/astro.md)).

## Who does what

gosfc only connects tools that already exist: goesm compiles Go, the Vue
tooling compiles SFCs and templates, Astro compiles `.astro` files and renders
pages, and Vite builds.
[The architecture](./reference/architecture.md) has the details.
