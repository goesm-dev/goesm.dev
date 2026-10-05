---
description: Add gosfc to an Astro or Vite project and write your first Go in a Vue component or an .astro file.
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
- [Go in .astro files](./astro.md): Go in the frontmatter and `<script>` of
  `.astro` files
- [Importing Go from JavaScript](./importing-go.md): `go:` imports in `.astro`,
  `.ts` and `.js` files
