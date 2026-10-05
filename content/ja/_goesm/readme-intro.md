<!-- Synced by cmd/syncdocs. Edit the source instead. -->

Go のパッケージを、ネイティブな ES モジュールにコンパイルします。出力は TypeScript で、WebAssembly は使いません。

goesm は、普通の Go モジュールにある普通の Go パッケージを、パッケージごとに 1 つの ES モジュールに変換します。Vite、Rolldown、esbuild、Bun、Node.js、ブラウザからそのまま import できます。export された Go の関数は JavaScript の関数に、export された型は TypeScript の型を持つクラスになります。整数演算、スライス、マップ、インターフェース、goroutine、`defer` / `panic` / `recover`、ジェネリクス、リフレクションは Go と同じ意味で動き、その結果はネイティブ Go と突き合わせて検証しています。

> [!NOTE]
> goesm は実験段階です。現時点で動くものと動かないものは、[現状](/ja/guide/status/)を参照してください。

```go
package cart

type Item struct {
	Name     string
	Price    int
	Quantity int
}

func Total(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Price * item.Quantity
	}
	return total
}

func Discount(total, percent int) int {
	return total * (100 - percent) / 100
}
```

```ts
import { Discount, Total } from "./goesm-ts/example.com/app/cart.ts";

const items = [{ Name: "りんご", Price: 120, Quantity: 3 }, { Name: "bread", Price: 250, Quantity: 1 }];
Total(items);          // 610
Discount(2408, 15);    // 2046: Go の整数除算。2046.8 ではない
```
