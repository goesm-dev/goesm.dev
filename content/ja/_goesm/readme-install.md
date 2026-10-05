<!-- Synced by cmd/syncdocs. Edit the source instead. -->

goesm には Go 1.27 以降のツールチェーンが必要です。それより古い `go` は、`GOTOOLCHAIN` の仕組みで Go 1.27 を自動でダウンロードします。Node.js や npm は不要です。ランタイムの `@goesm/runtime` は goesm のバイナリに埋め込まれていて、出力に書き出されます。

自分のモジュールのツールとして追加すると、自分のコードと同じツールチェーンでビルドされます。

```sh
go get -tool github.com/goesm-dev/goesm/cmd/goesm@latest
go tool goesm emit-ts ./cart     # goesm-ts/<cart の import パス>.ts + goesm-ts/@goesm/runtime/
```

`PATH` にインストールすることもできます。

```sh
go install github.com/goesm-dev/goesm/cmd/goesm@latest
```

ビルド済みのバイナリは配布していません。goesm は動作中に必ず `go` を実行するので、Go のツールチェーンはいずれにしても必要です。また、自分のツールチェーンで goesm をビルドすると、goesm の go/types がモジュールの使う Go と同じバージョンになります。実験段階の間、リリースは `v0.0.1-beta.N` という名前のプレリリースとして [GitHub Releases](https://github.com/goesm-dev/goesm/releases) で公開します。`@latest` は最新のリリースを指します。`goesm version` は、goesm のバージョンと、goesm のビルドに使った Go のバージョンを表示します。
