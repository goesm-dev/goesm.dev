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

TypeScript sees the types of these imports. With goesm v0.0.1-beta.8 or later, gosfc declares each package JavaScript imports with `go:` as a module, from the declaration file of its exports that goesm writes (`goesm emit-ts -dts`):

```ts
declare module "go:example.com/app/content" {
  export function Slugs(): Array<string>;
}
```

In Astro, `@gosfc/astro` writes the declarations to `.astro/integrations/_gosfc_astro/go-imports.d.ts`, which `.astro/types.d.ts` references, so `astro check` and editors see them. `astro sync`, and so `astro check`, compiles the Go to write them; `astro dev` and `astro build` rewrite them when they compile. With Vite alone, `@gosfc/vite` writes them to `gosfc.d.ts` in Vite's root, which the `include` of the project's `tsconfig.json` must cover; `gosfc({ dts: "src/gosfc.d.ts" })` picks another file and `gosfc({ dts: false })` none. Until the Go has been compiled once, the file declares every `go:` import as `any`.
