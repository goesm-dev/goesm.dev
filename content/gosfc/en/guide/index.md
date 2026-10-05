---
description: gosfc is a thin integration layer that lets the <script setup> of a Vue component be real Go, compiled by goesm.
---

# What is gosfc?

<!--@include: ../_gosfc/readme-intro.md-->

## Who does what

gosfc finds the Go block of a `.vue` file, turns it into an ordinary Go file,
and connects what goesm compiles to Vue. Everything else is done by the tools
that already do it:

| | Does | Does not |
| --- | --- | --- |
| gosfc | finds the Go block, builds the Go file, calls goesm, exposes the bindings to the template, keeps source positions; the Vite plugin and the Astro integration | parse or type-check Go, resolve modules, compile templates or styles, bundle, render |
| goesm | Go packages and modules, parsing, type checking, Go semantics, TypeScript output with source maps to `.go` | anything about Vue |
| Vue tooling | SFC parsing, template compilation, scoped CSS, HMR | Go |
| Vite | dev server, TypeScript to JavaScript, bundling | Go, SFCs |
| Astro | pages, SSR, static builds, islands | Go, the inside of SFCs |

[The architecture](../reference/architecture.md) describes each step.
