<template>
  <aside class="sidebar">
    <!-- On narrow screens the nav bar keeps only the brand, search and the
         theme toggle; its links, languages and GitHub move here. -->
    <nav class="sidebar-nav" aria-label="Main">
      <a v-for="l in bar.nav" :key="l.link" :href="l.link" :class="{ active: l.active }" :target="l.external ? '_blank' : undefined" :rel="l.external ? 'noreferrer' : undefined"
        >{{ l.text }}<svg v-if="l.external" class="external" viewBox="0 0 24 24" width="12" height="12" aria-hidden="true"><path fill="currentColor" d="M9 5v2h6.59L4 18.59 5.41 20 17 8.41V15h2V5z" /></svg></a
      >
      <a v-for="s in bar.social" :key="s.link" :href="s.link" target="_blank" rel="noreferrer">GitHub<svg class="external" viewBox="0 0 24 24" width="12" height="12" aria-hidden="true"><path fill="currentColor" d="M9 5v2h6.59L4 18.59 5.41 20 17 8.41V15h2V5z" /></svg></a>
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
