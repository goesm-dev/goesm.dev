<template>
  <div class="doc-container">
    <main id="main" class="doc-main">
      <p v-if="doc.synced" class="synced">
        {{ doc.ui.syncedFrom }} <a :href="doc.edit" target="_blank" rel="noreferrer">{{ doc.synced }}</a>
      </p>
      <article class="vp-doc" v-html="doc.html"></article>
      <footer class="doc-footer">
        <a class="edit-link" :href="doc.edit" target="_blank" rel="noreferrer">{{ doc.ui.editPage }}</a>
        <nav class="prev-next" aria-label="Pager">
          <a v-if="doc.prev.link" class="pager prev" :href="doc.prev.link">
            <span class="pager-label">{{ doc.ui.prev }}</span>
            <span class="pager-title">{{ doc.prev.text }}</span>
          </a>
          <a v-if="doc.next.link" class="pager next" :href="doc.next.link">
            <span class="pager-label">{{ doc.ui.next }}</span>
            <span class="pager-title">{{ doc.next.text }}</span>
          </a>
        </nav>
      </footer>
    </main>
    <aside v-if="hasOutline" class="outline" aria-labelledby="outline-title">
      <div class="outline-inner">
        <p id="outline-title" class="outline-title">{{ doc.ui.onThisPage }}</p>
        <ul>
          <li v-for="h in doc.outline" :key="h.id" :class="'level-' + h.level">
            <a :href="'#' + h.id">{{ h.text }}</a>
          </li>
        </ul>
      </div>
    </aside>
  </div>
</template>

<script setup lang="go">
import "goesm.dev/site"

type Props struct {
	Route string
}

doc := site.Doc(props.Route)
hasOutline := len(doc.Outline) > 0
</script>
