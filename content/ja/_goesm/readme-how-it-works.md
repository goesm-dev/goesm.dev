<!-- Synced by cmd/syncdocs. Edit the source instead. -->

```text
.go ──► go/packages + go/types ──► goesm: Go の意味論 → TypeScript ──► バンドラー / ランタイム ──► ES モジュール
         （Go のツールチェーン）                                         （Vite、Rolldown、esbuild、Bun、Node.js）
```

goesm は、ソース言語については Go のツールチェーンを、出力については現代の JS ツールを権威として扱います。Go のパーサー、型システム、モジュールを置き換えることも、バンドラーを実装することもしません。型検査済みの Go を TypeScript と、TypeScript で書いた小さなランタイムに変換します。主な設計判断は次のとおりで、すべて [ARCHITECTURE.ja.md](/ja/reference/architecture/) で説明しています。

- コンパイルの単位は Go のプログラムではなく Go のパッケージです。goesm のプロジェクトに `package main` は要らず、どこから実行を始めるかは JS のアプリケーションが決めます。
- goroutine は協調的です。プログラム全体の解析でブロックしうる関数を見つけ、それらはブロック箇所に `await` を置いた `async` 関数になります。それ以外はすべて素の同期 JavaScript のままです。
- 標準ライブラリのパッケージは Go 自身のソースからコンパイルします。gc ランタイムに結びついた少数のパッケージ（`runtime`、`reflect`、`internal/reflectlite`、`sync`、`syscall/js`）は goesm 独自の Go ソースに置き換え、Go の本体を持たない関数はランタイムで実装します。
- Go の各型は実行時の型記述子を持ち、インターフェース、ジェネリクス（型辞書）、構造体をキーにしたマップ、リフレクションがそれを使います。

設計を GopherJS と比べたものは [docs/gopherjs-comparison.ja.md](/ja/reference/gopherjs-comparison/) にあります。

### 目標と目標外

goesm は、Go の開発体験（`go.mod`、`go.work`、`go fmt`、`gopls`、`go test`、`go vet`、`golangci-lint`、`govulncheck`）をそのまま保ち、出力側では JavaScript エコシステムのツールを再利用することを目指します。

Go 風の言語、WebAssembly ランタイム、パッケージマネージャー、Vite や Rolldown の代替、フレームワークではありません。フレームワークとの統合は別のアダプターの役割です。たとえば `<script setup lang="go">` を持つ Vue SFC は、Go の部分を goesm に渡し、テンプレート層は Vue が受け持つ、という形にできます。
