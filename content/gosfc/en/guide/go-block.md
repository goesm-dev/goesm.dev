---
description: What a <script setup lang="go"> block can contain, how the template sees its variables and functions, and how it receives props.
---

# Writing the Go block

<!--@include: ../_gosfc/readme-writing-the-go-block.md-->

## State and events

Variables the block declares are the component's state. A template event that
calls a Go function runs it, and the template then shows the new values:

```vue
<template>
  <button type="button" @click="Increment">Clicked {{ count }} times</button>
  <input :value="name" @input="Rename($event.target.value)" />
  <p>Hello, {{ name }}</p>
</template>

<script setup lang="go">
count := 0
name := "Gopher"

func Increment() {
	count++
}

func Rename(s string) {
	name = s
}
</script>
```

The template gets a copy of each value, converted to JavaScript (a slice
becomes an array, a struct an object). Changing that copy in JavaScript does
not change the Go variable: change Go state by calling a Go function.
