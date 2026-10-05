<!-- Synced by cmd/syncdocs. Edit the source instead. -->

A frontmatter fence that opens with `---go` instead of `---` holds Go. It follows the rules of a Go block of a `.vue` file: the top level runs once per render (at build time, or per request with SSR), and its top-level variables, constants and functions are available in the template.

```astro
---go
import (
	"strconv"

	cart "example.com/app/src/features/cart/pkg"
	Line "../features/cart/Line.astro"
	Summary "../features/cart/Summary.vue"
)

items := []cart.Item{
	{Price: 120, Quantity: 3},
	{Price: 80, Quantity: 1},
}
total := cart.Total(items)

func Yen(n int) string {
	return "¥" + strconv.Itoa(n)
}
---

<Line label="Apples" price={120} quantity={3} />
<p>Total: {Yen(total)} ({items.length} items)</p>
<Summary />
```

* JavaScript modules (components, styles, npm packages) are imported with Go import syntax. `import Card "./Card.vue"` means `import Card from "./Card.vue"`, and `import _ "./global.css"` is a side-effect import. An import path is JavaScript when it starts with `./`, `../`, `/` or `@`, or contains `:` (such as `astro:assets`); any other path is a Go package.
* `type Props struct {...}` receives `Astro.props`, with the same field names as the props of a `.vue` block (`Label` is read from `label`).
* Values reach the template converted for JavaScript (strings, slices as arrays, structs as objects), and the template can call Go functions.
* If the Go code blocks (channels, `time.Sleep`, ...), the frontmatter waits for it.

A `<script lang="go">` in the template is Go that runs in the browser once per page, like Astro's processed `<script>`: Astro bundles it, and a component used several times includes it once. The template does not see its top-level names; use `syscall/js` to work with the DOM.

```astro
<button id="counter">0</button>

<script lang="go">
import (
	"strconv"
	"syscall/js"
)

count := 0
button := js.Global().Get("document").Call("getElementById", "counter")
button.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
	count++
	button.Set("textContent", strconv.Itoa(count))
	return nil
}))
</script>
```

Limitations:

* A Go frontmatter cannot export anything, so `getStaticPaths` and other exports need a TypeScript frontmatter. That frontmatter can still import Go with a `go:` specifier (next section).
* Go sees nothing of the `Astro` global but the props. For `Astro.url`, `Astro.cookies`, redirects and so on, use a TypeScript frontmatter.
* Go import syntax can write a default import or a side-effect import only. To use named exports of a JavaScript module, re-export them as the default export of a small module of your own and import that.
* `<script lang="go">` takes no other attributes (`is:inline`, `define:vars`, ...).
