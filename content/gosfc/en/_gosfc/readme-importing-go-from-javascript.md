<!-- Synced by cmd/syncdocs. Edit the source instead. -->

A `.js`, `.ts` or `.astro` module inside a Go module can import a Go package directly with a `go:` specifier. gosfc compiles the package with goesm in the Go module of the importing file, and Vite bundles it like any other module. The API is goesm's [JS calling ABI](https://github.com/goesm-dev/goesm/blob/main/docs/js-exports.md): the exported functions take and return plain JavaScript values (strings, arrays, objects).

```astro
---
// src/pages/[slug].astro
import { Slugs } from "go:example.com/app/content";

export function getStaticPaths() {
  return Slugs().map((slug) => ({ params: { slug } }));
}
---
```
