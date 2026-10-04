<template>
  <aside class="sidebar">
    <!-- On narrow screens the nav bar keeps only the brand, search and the
         theme toggle; its links, languages and GitHub move here. -->
    <nav class="sidebar-nav" aria-label="Main">
      <a v-for="l in bar.nav" :key="l.link" :href="l.link" :class="{ active: l.active }" :target="l.external ? '_blank' : undefined" :rel="l.external ? 'noreferrer' : undefined"
        >{{ l.text }}<span v-if="l.external" class="external" aria-hidden="true">↗</span></a
      >
      <a v-for="s in bar.social" :key="s.link" :href="s.link" target="_blank" rel="noreferrer">GitHub<span class="external" aria-hidden="true">↗</span></a>
      <p class="sidebar-title">{{ bar.ui.language }}</p>
      <a v-for="a in bar.locales" :key="a.lang" :href="a.link" :lang="a.lang" :hreflang="a.lang" :class="{ active: a.current }">{{ a.label }}</a>
    </nav>
    <nav v-if="groups && groups.length" aria-label="Sidebar">
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
bar := site.NavBar(props.Route)
</script>
