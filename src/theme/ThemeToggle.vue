<template>
  <button type="button" class="theme-toggle" :title="label" :aria-label="label" @click="Toggle">
    <svg class="sun" viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"><path fill="currentColor" d="M12 7a5 5 0 1 0 0 10 5 5 0 0 0 0-10zm0-5 1 3h-2zm0 20-1-3h2zM2 12l3-1v2zm20 0-3 1v-2zM4.9 4.9l2.8 1.4-1.4 1.4zm14.2 14.2-2.8-1.4 1.4-1.4zM4.9 19.1l1.4-2.8 1.4 1.4zM19.1 4.9l-1.4 2.8-1.4-1.4z" /></svg>
    <svg class="moon" viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"><path fill="currentColor" d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z" /></svg>
  </button>
</template>

<script setup lang="go">
// Dark mode: the class on <html> is set before the first paint by a small
// inline script in the layout; this island only flips it and remembers the
// choice. Both icons are rendered and CSS shows the right one, so the
// server-rendered HTML never disagrees with the client.
import "syscall/js"

type Props struct {
	Label string
}

label := props.Label

func Toggle() {
	root := js.Global().Get("document").Get("documentElement")
	dark := root.Get("classList").Call("toggle", "dark").Bool()
	mode := "light"
	if dark {
		mode = "dark"
	}
	js.Global().Get("localStorage").Call("setItem", "goesm-theme", mode)
}
</script>
