<!-- Synced by cmd/syncdocs. Edit the source instead. -->

Go モジュールの中にある `.js`、`.ts`、`.astro` のモジュールは、`go:` specifier で Go パッケージを直接インポートできます。gosfc はインポートした側のファイルが属する Go モジュールで goesm を使ってそのパッケージをコンパイルし、Vite が他のモジュールと同じようにバンドルします。API は goesm の [JS 呼び出し ABI](https://github.com/goesm-dev/goesm/blob/main/docs/js-exports.ja.md) に従います。エクスポートされた関数は、文字列、配列、オブジェクトといった普通の JavaScript の値を受け取り、返します。

```astro
---
// src/pages/[slug].astro
import { Slugs } from "go:example.com/app/content";

export function getStaticPaths() {
  return Slugs().map((slug) => ({ params: { slug } }));
}
---
```

これらの import には TypeScript の型が付きます。goesm v0.0.1-beta.8 以降では、goesm がエクスポートの宣言ファイルを書き出します (`goesm emit-ts -dts`)。gosfc はそれをもとに、JavaScript が `go:` で import する各パッケージをモジュールとして宣言します。

```ts
declare module "go:example.com/app/content" {
  export function Slugs(): Array<string>;
}
```

Astro では、`@gosfc/astro` がこの宣言を `.astro/integrations/_gosfc_astro/go-imports.d.ts` に書き出します。このファイルは `.astro/types.d.ts` から参照されるので、`astro check` とエディターから型が見えます。`astro sync` と、それを実行する `astro check` は、Go をコンパイルして宣言を書き出します。`astro dev` と `astro build` は、コンパイルするたびに宣言を書き直します。Vite だけを使う場合、`@gosfc/vite` は Vite のルートの `gosfc.d.ts` に宣言を書き出します。このファイルは、プロジェクトの `tsconfig.json` の `include` に含めます。`gosfc({ dts: "src/gosfc.d.ts" })` で別のファイルを指定でき、`gosfc({ dts: false })` で書き出さないようにできます。Go を一度もコンパイルしていない間は、すべての `go:` の import を `any` として宣言します。
