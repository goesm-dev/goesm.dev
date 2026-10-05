---
description: gosfc は、Vue コンポーネントの <script setup> を本物の Go にする薄い統合レイヤーです。Go は goesm がコンパイルします。
---

# gosfc とは

<!--@include: ../_gosfc/readme-intro.md-->

## 役割分担

gosfc は `.vue` ファイルの Go ブロックを見つけて普通の Go ファイルに組み直し、goesm が
コンパイルしたものを Vue につなぎます。それ以外は、それを本来担当するツールが行います。

| | 担当すること | 担当しないこと |
| --- | --- | --- |
| gosfc | Go ブロックの検出、Go ファイルの構築、goesm の呼び出し、テンプレートへのバインディングの公開、ソース位置の維持、Vite プラグインと Astro インテグレーション | Go の構文解析・型検査、モジュールの解決、テンプレートやスタイルのコンパイル、バンドル、描画 |
| goesm | Go のパッケージとモジュール、構文解析、型検査、Go の意味論、`.go` へのソースマップ付きの TypeScript 出力 | Vue に関すること |
| Vue のツール | SFC の解析、テンプレートのコンパイル、scoped CSS、HMR | Go |
| Vite | 開発サーバー、TypeScript から JavaScript への変換、バンドル | Go、SFC |
| Astro | ページ、SSR、静的ビルド、アイランド | Go、SFC の中身 |

各ステップは [アーキテクチャ](../reference/architecture.md) で説明しています。
