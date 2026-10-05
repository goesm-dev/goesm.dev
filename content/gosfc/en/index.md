---
layout: home
description: gosfc runs real Go in the <script setup> of Vue components, and lets Astro pages import Go packages. Compiled to plain JavaScript by goesm, no WebAssembly.
hero:
  name: gosfc
  text: Go in Vue components
  tagline: Write <script setup lang="go"> in a .vue file. It is real Go, compiled by goesm to plain JavaScript; Vue renders it, and Vite and Astro build it.
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
    title: Astro and SSR
    details: One integration for Astro. Components render on the server and hydrate as islands, and .astro files import Go with go:.
    link: /guide/importing-go
  - icon: ⚡
    title: No WebAssembly
    details: goesm compiles the Go to ES modules. The client JS of the benchmark page is 13% larger (gzip) than the same components in TypeScript.
    link: /guide/benchmark
  - icon: 🎯
    title: Errors point at the .vue file
    details: Go's diagnostics, the source maps and the stack trace of a panic all point at lines of the .vue file.
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

## Who does what

gosfc only connects tools that already exist: goesm compiles Go, the Vue
tooling compiles SFCs and templates, Vite builds, and Astro renders pages.
[The architecture](./reference/architecture.md) has the details.
