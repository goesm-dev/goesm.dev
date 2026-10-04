<!-- Synced by cmd/syncdocs. Edit the source instead. -->

A `.js`, `.ts` or `.astro` module inside a Go module can import a Go package directly with a `go:` specifier. gosfc compiles the package with goesm in the Go module of the importing file, and Vite bundles it like any other module. The API is goesm's: exported functions and types, with Go strings and slices converted through the runtime each module re-exports as `$runtime`.

```astro
---
// src/pages/[slug].astro
import { Slugs, $runtime as rt } from "go:example.com/app/content";

export function getStaticPaths() {
  return rt.toArray(Slugs()).map((s) => ({ params: { slug: rt.toJSString(s) } }));
}
---
```
