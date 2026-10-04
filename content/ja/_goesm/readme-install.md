<!-- Synced by cmd/syncdocs. Edit the source instead. -->

goesm には Go のツールチェーン（Go 1.27 以降。古い `go` は `GOTOOLCHAIN` で 1.27 を自動でダウンロードします）が必要です。Node.js や npm は不要です。ランタイム（`@goesm/runtime`）はバイナリに埋め込まれていて、出力に書き出されます。

自分のモジュールのツールとして追加すると、自分のコードと同じツールチェーンでビルドされます。

```sh
go get -tool github.com/goesm-dev/goesm/cmd/goesm@latest
go tool goesm emit-ts ./cart     # goesm-ts/<cart の import パス>.ts + goesm-ts/@goesm/runtime/
```

`PATH` にインストールすることもできます。

```sh
go install github.com/goesm-dev/goesm/cmd/goesm@latest
```

ビルド済みバイナリはありません。goesm はどのみち `go` を実行しますし、自分のツールチェーンでビルドすれば goesm の go/types がモジュールの使う Go と揃います。実験段階の間、リリースは `v0.0.1-beta.N` という名前のプレリリースです（[GitHub Releases](https://github.com/goesm-dev/goesm/releases)）。`@latest` は最新のものを指します。`goesm version` で goesm のバージョンと、ビルドに使った Go を表示します。
