<!-- Synced by cmd/syncdocs. Edit the source instead. -->

```text
.go ──► go/packages + go/types ──► goesm: Go の意味論 → TypeScript ──► バンドラー / ランタイム ──► ES モジュール
```

go/packages と go/types は Go のツールチェーンの一部です。バンドラーやランタイムには、Vite、Rolldown、esbuild、Bun、Node.js などを使えます。

goesm は、ソース言語については Go のツールチェーンを、出力については現代の JS ツールを権威として扱います。Go のパーサー、型システム、モジュールを置き換えることも、バンドラーを実装することもしません。goesm は、型検査済みの Go を TypeScript に変換し、TypeScript で書いた小さなランタイムと組み合わせます。主な設計判断は次のとおりで、すべて [ARCHITECTURE.ja.md](/ja/reference/architecture/) で説明しています。

- コンパイルの単位は Go のプログラムではなく Go のパッケージです。goesm のプロジェクトに `package main` は要らず、どこから実行を始めるかは JS のアプリケーションが決めます。
- goroutine は協調的に動きます。goesm はプログラム全体を解析してブロックしうる関数を見つけ、それらをブロック箇所に `await` を置いた `async` 関数に変換します。それ以外の関数は、すべて普通の同期関数のままです。
- 標準ライブラリのパッケージは Go 自身のソースからコンパイルします。`runtime`、`reflect`、`internal/reflectlite`、`sync`、`syscall/js` など、gc ランタイムに結びついた少数のパッケージは、goesm 独自の Go ソースに置き換えます。Go の本体を持たない関数は、ランタイムで実装します。
- Go の各型は、実行時の型記述子を持ちます。インターフェース、型辞書を使うジェネリクス、構造体をキーにしたマップ、リフレクションは、この型記述子を使います。

設計を GopherJS と比べたものは [docs/gopherjs-comparison.ja.md](/ja/reference/gopherjs-comparison/) にあります。

### 目標と目標外

goesm は、`go.mod`、`go.work`、`go fmt`、`gopls`、`go test`、`go vet`、`golangci-lint`、`govulncheck` といった Go の開発体験をそのまま保ち、出力側では JavaScript エコシステムのツールを再利用することを目指します。

goesm は、Go 風の言語、WebAssembly ランタイム、パッケージマネージャー、Vite や Rolldown の代替、フレームワークのいずれでもありません。フレームワークとの統合は別のアダプターの役割です。たとえば `<script setup lang="go">` を持つ Vue SFC は、Go の部分を goesm に渡し、テンプレート層は Vue が受け持つ、という形にできます。
