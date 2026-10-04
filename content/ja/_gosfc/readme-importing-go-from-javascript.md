<!-- Synced by cmd/syncdocs. Edit the source instead. -->

Go module の中にある `.js`、`.ts`、`.astro` の module は、`go:` specifier で Go package を直接 import できます。gosfc は import した側のファイルが属する Go module で goesm を使ってその package をコンパイルし、Vite が他の module と同じように bundle します。API は goesm のものです：export された関数と型があり、Go の文字列や slice は各 module が `$runtime` として re-export する runtime で変換します。

```astro
---
// src/pages/[slug].astro
import { Slugs, $runtime as rt } from "go:example.com/app/content";

export function getStaticPaths() {
  return rt.toArray(Slugs()).map((s) => ({ params: { slug: rt.toJSString(s) } }));
}
---
```
