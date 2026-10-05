<!-- Synced by cmd/syncdocs. Edit the source instead. -->

An integration layer for using real Go in the `<script setup>` of Vue Single File Components.

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

`.go` files are ordinary Go packages, and imports are ordinary Go imports. Go compilation is handled by [goesm](https://github.com/goesm-dev/goesm), SFCs and templates by the Vue tooling, builds by Vite, and pages and SSR by Astro. See [ARCHITECTURE.md](/gosfc/reference/architecture/) for the design.

**Status: PoC.** What is not implemented yet is listed in ARCHITECTURE.md, section 11 ("Not yet implemented / open questions").
