<!-- Synced by cmd/syncdocs. Edit the source instead. -->

同じコンポーネントを `<script setup lang="go">` と `<script setup lang="ts">` で書いて比べています（[bench/](https://github.com/goesm-dev/gosfc/tree/3eb7cf0032859b09b29ec08802b29a1053e15399/bench)）。どちらも Vite 8 と `@vitejs/plugin-vue` でビルドし、Go 側は前に `@gosfc/vite` を置く以外は同じ設定です。

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/repo/gosfc/bench/results-dark.svg">
  <img alt="ベンチマークの棒グラフ：client build time、client JS のサイズ、SSR render time を gosfc と Vue で比較（数値は下の表）" src="/repo/gosfc/bench/results-light.svg">
</picture>

| | gosfc (`lang="go"`) | Vue (`lang="ts"`) | ratio |
|---|---:|---:|---:|
| Client build time | 196 ms | 130 ms | 1.51x |
| Client JS (minified) | 68.2 KiB | 59.8 KiB | 1.14x |
| Client JS (gzip) | 26.4 KiB | 23.3 KiB | 1.13x |
| SSR render, small component | 13.8 µs | 12.9 µs | 1.07x |
| SSR render, 1,000,000 items | 104.9 ms | 95.8 ms | 1.09x |

* **Client build time**：カウンターとカートのサマリーを置いたページの `vite build`。それぞれウォームアップのビルドを 1 回してから、先にビルドする側を入れ替えながら 9 回測った中央値です。goesm のバイナリと go コマンドのビルドキャッシュは温まった状態で、編集してビルドし直すときと同じ条件です。
* **Client JS**：そのビルドの JS ファイルすべて（Vue のランタイムを含む）。差（+8.4 KiB、gzip で +3.1 KiB）は goesm のランタイムと Go の意味論を守るためのコードです。
* **SSR render**：プロダクションの SSR ビルドで `renderToString` した時間で、両者は同じ HTML を出力します。それぞれの側とコンポーネントごとに新しい Node プロセスで、順番を入れ替えながら 5 回測った中央値です。「small component」は 3 品のカートのサマリー、「1,000,000 items」はコンポーネントの setup で 1,000,000 品を作って合計するもので、どちらも時間の大半はメモリ確保です。そこでの小さな差は、goesm が Go の意味論を守るために生成するコード（たとえば `range` での struct の値コピー）によるものと思われますが、プロファイルはまだ取っていません。

2026-10-04 に `pnpm bench` で測りました。`mise.toml` のバージョン（Node.js 26.10.0、Go 1.27.1）と goesm 20dbf1d を使い、4 vCPU の Intel Xeon 2.80GHz のクラウド VM で実行しています。サイズは決定的ですが、時間はこの環境では実行ごとに揺れます（6 回の実行で、比はビルドが 1.51x〜1.73x、small component が 1.02x〜1.26x、1,000,000 items が 1.09x〜1.20x でした）。
