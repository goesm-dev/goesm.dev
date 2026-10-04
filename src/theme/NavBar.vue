<template>
  <header class="navbar">
    <div class="navbar-inner">
      <label v-if="menu" for="nav-toggle" class="menu-button" :title="bar.ui.menu" :aria-label="bar.ui.menu">
        <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true"><path fill="currentColor" d="M3 6h18v2H3zm0 5h18v2H3zm0 5h18v2H3z" /></svg>
      </label>
      <a class="brand" :href="bar.home">
        <span class="brand-mark" aria-hidden="true"></span>
        <span class="brand-name">{{ bar.title }}</span>
      </a>
      <div class="navbar-search"><slot name="search" /></div>
      <nav class="navbar-links" aria-label="Main">
        <a
          v-for="l in bar.nav"
          :key="l.link"
          :href="l.link"
          :class="{ active: l.active }"
          :target="l.external ? '_blank' : undefined"
          :rel="l.external ? 'noreferrer' : undefined"
          >{{ l.text }}<span v-if="l.external" class="external" aria-hidden="true">↗</span></a
        >
      </nav>
      <details class="locale-menu">
        <summary :title="bar.ui.language" :aria-label="bar.ui.language">
          <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true"><path fill="currentColor" d="M12.87 15.07l-2.54-2.51.03-.03A17.52 17.52 0 0 0 14.07 6H17V4h-7V2H8v2H1v2h11.17C11.5 7.92 10.44 9.75 9 11.35 8.07 10.32 7.3 9.19 6.69 8h-2c.73 1.63 1.73 3.17 2.98 4.56l-5.09 5.02L4 19l5-5 3.11 3.11.76-2.04zM18.5 10h-2L12 22h2l1.12-3h4.75L21 22h2l-4.5-12zm-2.62 7l1.62-4.33L19.12 17h-3.24z" /></svg>
          <span>{{ current }}</span>
        </summary>
        <ul>
          <li v-for="a in bar.locales" :key="a.lang">
            <a :href="a.link" :lang="a.lang" :hreflang="a.lang" :class="{ active: a.current }">{{ a.label }}</a>
          </li>
        </ul>
      </details>
      <slot name="theme" />
      <a v-for="s in bar.social" :key="s.link" class="social" :href="s.link" target="_blank" rel="noreferrer" aria-label="GitHub">
        <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true"><path fill="currentColor" d="M12 .3a12 12 0 0 0-3.8 23.4c.6.1.8-.3.8-.6v-2c-3.3.7-4-1.6-4-1.6-.6-1.4-1.4-1.8-1.4-1.8-1-.7.1-.7.1-.7 1.2 0 1.9 1.2 1.9 1.2 1 1.8 2.8 1.3 3.5 1 0-.8.4-1.3.7-1.6-2.7-.3-5.5-1.3-5.5-6 0-1.2.5-2.3 1.3-3.1-.2-.4-.6-1.6 0-3.2 0 0 1-.3 3.4 1.2a11.5 11.5 0 0 1 6 0c2.3-1.5 3.3-1.2 3.3-1.2.6 1.6.2 2.8 0 3.2.9.8 1.3 1.9 1.3 3.2 0 4.6-2.8 5.6-5.5 5.9.5.4.9 1 .9 2.2v3.3c0 .3.1.7.8.6A12 12 0 0 0 12 .3" /></svg>
      </a>
    </div>
  </header>
</template>

<script setup lang="go">
import "goesm.dev/site"

type Props struct {
	Route string
	Menu  bool
}

bar := site.NavBar(props.Route)
menu := props.Menu

current := ""
for _, l := range bar.Locales {
	if l.Current {
		current = l.Label
	}
}
</script>
