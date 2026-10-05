<!-- Synced by cmd/syncdocs. Edit the source instead. -->

Compile Go packages into native ES modules, emitted as TypeScript. No WebAssembly.

goesm takes ordinary Go packages from an ordinary Go module and turns each one into an ES module that Vite, Rolldown, esbuild, Bun, Node.js or a browser can import directly. Exported Go functions become JavaScript functions, exported types become classes with TypeScript types, and Go's semantics (integer arithmetic, slices, maps, interfaces, goroutines, `defer` / `panic` / `recover`, generics, reflection) are preserved and checked against native Go.

> [!NOTE]
> goesm is experimental. See [Status](/guide/status/) for what works today and what does not.

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
Discount(2408, 15);    // 2046: Go's integer division, not 2046.8
```
