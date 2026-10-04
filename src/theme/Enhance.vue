<template>
  <span hidden></span>
</template>

<script setup lang="go">
// Page behaviour that needs the DOM: copy buttons on code blocks and the
// outline that follows the scroll position. Written in Go with syscall/js.
import "goesm.dev/client/enhance"

type Props struct {
	Copied string
}

enhance.Install(props.Copied)
</script>
