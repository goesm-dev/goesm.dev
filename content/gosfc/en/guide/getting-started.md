---
description: Add gosfc to an Astro or Vite project and write your first <script setup lang="go">.
---

# Getting started

You need Go 1.27 or later and Node.js. The app is a Go module (a `go.mod` at
its root) as well as a Node.js project: the Go packages your components import
live in it.

## Install and set up Astro

<!--@include: ../_gosfc/readme-usage-astro.md-->

## Next steps

- [Writing the Go block](./go-block.md): what the block can contain and how
  the template sees it
- [Importing Go from JavaScript](./importing-go.md): `go:` imports in `.astro`,
  `.ts` and `.js` files
