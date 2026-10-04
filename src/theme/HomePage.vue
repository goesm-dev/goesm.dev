<template>
  <main id="main" class="home">
    <section class="hero">
      <div class="hero-text">
        <h1 class="hero-name"><span class="clip">{{ home.hero.name }}</span></h1>
        <p class="hero-title">{{ home.hero.text }}</p>
        <p class="hero-tagline">{{ home.hero.tagline }}</p>
        <div class="hero-actions">
          <a
            v-for="a in home.hero.actions"
            :key="a.link"
            :class="['button', a.theme]"
            :href="a.link"
            :target="a.link.includes('://') ? '_blank' : undefined"
            :rel="a.link.includes('://') ? 'noreferrer' : undefined"
            >{{ a.text }}</a
          >
        </div>
      </div>
      <div v-if="home.hero.image" class="hero-image">
        <div class="hero-glow" aria-hidden="true"></div>
        <img :src="home.hero.image" :alt="home.hero.name" width="1000" height="250" />
      </div>
    </section>
    <section class="features">
      <component :is="f.link ? 'a' : 'div'" v-for="f in home.features" :key="f.title" class="feature" :href="f.link || undefined">
        <div class="feature-icon" aria-hidden="true">{{ f.icon }}</div>
        <h2 class="feature-title">{{ f.title }}</h2>
        <p class="feature-details">{{ f.details }}</p>
      </component>
    </section>
    <slot name="demo" />
    <section class="vp-doc home-content" v-html="home.html"></section>
  </main>
</template>

<script setup lang="go">
import "goesm.dev/site"

type Props struct {
	Route string
}

home := site.Home(props.Route)
</script>
