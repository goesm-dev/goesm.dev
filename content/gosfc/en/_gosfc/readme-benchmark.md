<!-- Synced by cmd/syncdocs. Edit the source instead. -->

The same components written with `<script setup lang="go">` and with `<script setup lang="ts">` ([bench/](https://github.com/goesm-dev/gosfc/tree/a2bb163992154f4511b6988e2cbae4a372512717/bench)), built with Vite 8 and `@vitejs/plugin-vue`. The Go side adds `@gosfc/vite` in front; nothing else differs.

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/repo/gosfc/bench/results-dark.svg">
  <img alt="Bar chart of the benchmark: gosfc vs. Vue for client build time, client JS size and SSR render time (numbers in the table below)" src="/repo/gosfc/bench/results-light.svg">
</picture>

| | gosfc (`lang="go"`) | Vue (`lang="ts"`) | ratio |
|---|---:|---:|---:|
| Client build time | 196 ms | 130 ms | 1.51x |
| Client JS (minified) | 68.2 KiB | 59.8 KiB | 1.14x |
| Client JS (gzip) | 26.4 KiB | 23.3 KiB | 1.13x |
| SSR render, small component | 13.8 µs | 12.9 µs | 1.07x |
| SSR render, 1,000,000 items | 104.9 ms | 95.8 ms | 1.09x |

* **Client build time**: `vite build` of a page with a counter and a cart summary. Median of 9 builds, alternating which side builds first, after one warm-up build per side, so the goesm binary and the go command's build cache are warm, as in an edit-and-rebuild loop.
* **Client JS**: every JS file of that build, Vue runtime included. The difference (+8.4 KiB, +3.1 KiB gzip) is goesm's runtime and the code that keeps Go semantics.
* **SSR render**: `renderToString` with the production SSR build, both sides rendering the same HTML. Each side and component is measured in its own fresh Node process, 5 times with the order alternating, and the table shows the median. "small component" is a cart summary of 3 items; "1,000,000 items" builds and sums 1,000,000 items in the component's setup, which on both sides is mostly allocation. The small gap there is presumably the code goesm generates to keep Go semantics (for example, `range` copies each struct value); it has not been profiled.

Measured with `pnpm bench` on 2026-10-04 with the `mise.toml` versions (Node.js 26.10.0, Go 1.27.1) and goesm 20dbf1d, on a 4 vCPU Intel Xeon 2.80GHz cloud VM. Sizes are exact; timings move between runs on that machine (over 6 runs the ratios ranged from 1.51x to 1.73x for the build, 1.02x to 1.26x for the small component and 1.09x to 1.20x for 1,000,000 items).
