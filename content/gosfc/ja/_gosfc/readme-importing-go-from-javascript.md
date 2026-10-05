<!-- Synced by cmd/syncdocs. Edit the source instead. -->

Go モジュールの中にある `.js`、`.ts`、`.astro` のモジュールは、`go:` specifier で Go パッケージを直接インポートできます。gosfc はインポートした側のファイルが属する Go モジュールで goesm を使ってそのパッケージをコンパイルし、Vite が他のモジュールと同じようにバンドルします。API は goesm のものです：エクスポートされた関数と型があり、Go の文字列やスライスは各モジュールが `$runtime` として再エクスポートするランタイムで変換します。

```astro
---
// src/pages/[slug].astro
import { Slugs, $runtime as rt } from "go:example.com/app/content";

export function getStaticPaths() {
  return rt.toArray(Slugs()).map((s) => ({ params: { slug: rt.toJSString(s) } }));
}
---
```
