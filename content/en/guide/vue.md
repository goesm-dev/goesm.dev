---
description: gosfc runs real Go in the <script setup> of Vue components and lets Astro pages import Go packages. Its documentation is at goesm.dev/gosfc.
---

# Vue and Astro (gosfc)

[gosfc](/gosfc/) connects goesm to Vue and Astro. A Vue component can have a
`<script setup lang="go">`, and an Astro page can import a Go package with a
`go:` specifier. goesm compiles the Go; Vue, Vite and Astro do the rest as
usual.

```vue
<template>
  <div>
    Total: {{ total }}
  </div>
</template>

<script setup lang="go">
import cart "example.com/app/src/features/cart/pkg"

items := []cart.Item{
	{
		Price:    100,
		Quantity: 2,
	},
}

total := cart.Total(items)
</script>
```

This site is built that way: its theme is Vue components whose script is Go
(see [How this site is built](./this-site.md)).

::: tip gosfc documentation
Setup, the rules of the Go block, importing Go from JavaScript and the
benchmark are in the [gosfc documentation](/gosfc/guide/).
:::
