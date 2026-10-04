<template>
  <aside class="sidebar">
    <nav aria-label="Sidebar">
      <section v-for="g in groups" :key="g.text" class="sidebar-group">
        <p class="sidebar-title">{{ g.text }}</p>
        <a v-for="item in g.items" :key="item.link" :href="item.link" :class="{ active: item.active }" :aria-current="item.active ? 'page' : undefined">{{
          item.text
        }}</a>
      </section>
    </nav>
  </aside>
</template>

<script setup lang="go">
import "goesm.dev/site"

type Props struct {
	Route string
}

groups := site.Sidebar(props.Route)
</script>
