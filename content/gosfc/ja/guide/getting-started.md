---
description: Astro か Vite のプロジェクトに gosfc を追加し、最初の <script setup lang="go"> を書きます。
---

# はじめに

Go 1.27 以降と Node.js が必要です。アプリは Node.js のプロジェクトであると同時に
Go のモジュールでもあります (ルートに `go.mod`)。コンポーネントが import する Go の
パッケージはその中に置きます。

## インストールと Astro の設定

<!--@include: ../_gosfc/readme-usage-astro.md-->

## 次に読むもの

- [Go ブロックの書き方](./go-block.md): ブロックに書けるものと、テンプレートからの見え方
- [JavaScript から Go を import する](./importing-go.md): `.astro`、`.ts`、`.js`
  ファイルでの `go:` の import
