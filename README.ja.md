# goesm.dev

[English](README.md)

[goesm](https://github.com/goesm-dev/goesm) のウェブサイトです。goesm 自身で作っています。
Go で書いた VitePress 風のドキュメントエンジンを goesm でコンパイルし、
[gosfc](https://github.com/goesm-dev/gosfc) を通して Astro と Vue で描画します。
全体の仕組みはサイトの [このサイトの作り方](content/ja/guide/this-site.md) で説明しています。

## 構成

| パス | 内容 |
| --- | --- |
| `press/` | エンジン: Markdown (goldmark)、コンテナ、ハイライト、サイドバー、ロケール、検索インデックス、サイトマップ |
| `site/` | サイトの設定と、ページから呼ぶ関数 |
| `content/en`、`content/ja` | ページ (Markdown) |
| `client/` | ブラウザで動く Go (検索、ライブデモ、ページの補助) |
| `src/` | Astro のページと Vue のテーマ。コンポーネントの `<script setup>` は Go |
| `cmd/syncdocs` | goesm と gosfc のドキュメントを `go.mod` のバージョンでコピー |
| `third_party/gosfc` | gosfc の git submodule (Vite と Astro のパッケージ) |

## 開発

ツールのバージョンは `mise.toml` で固定しています。

```sh
git clone --recurse-submodules https://github.com/goesm-dev/goesm.dev
cd goesm.dev
mise install
pnpm install
pnpm dev          # http://localhost:4321
```

```sh
go test ./...     # エンジンとサイトをネイティブでテスト
pnpm build        # dist/ に静的サイト
node --test tests/  # dist/ の検査
```

`content/*/reference/` と `content/*/_goesm`、`_gosfc` 以下は生成されたページです。
更新するには `go.mod` の goesm か gosfc (と submodule) を上げてから、goesm を
このリポジトリの隣に clone した状態で次を実行します。

```sh
go run ./cmd/syncdocs -goesm ../goesm -gosfc third_party/gosfc
```

CI は同じものを `-check` 付きで実行します。

## デプロイ

サイトは Cloudflare Workers で、静的アセットだけの Worker として動きます。
`cloudflare.config.ts` が Worker を、`wrangler.config.ts` が `dist/` を指定します。

```sh
pnpm build:cf     # dist/ と .cloudflare/output の Build Output を作る
pnpm preview:cf   # Workers と同じ動きでローカルに配信
pnpm run deploy   # ビルドしてアップロード (`cf auth login` が必要)
```

CI は main へのマージごとに `mise run build` と `mise run deploy` でデプロイし、
そのコミットに `vYYYY.M.N` のタグを付けます (N は月ごとに 0 から、UTC)。
認証にはリポジトリのシークレット `CF_ID` と `CF_TOKEN` を使います。
