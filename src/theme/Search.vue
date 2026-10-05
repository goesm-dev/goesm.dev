<template>
  <div class="search">
    <button type="button" class="search-button" :aria-label="label" @click="Open" @pointerenter="Preload" @focus="Preload">
      <svg viewBox="0 0 24 24" width="16" height="16" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><circle cx="10.5" cy="10.5" r="6.5" /><path d="m15.5 15.5 5 5" /></svg>
      <span class="search-label">{{ label }}</span>
      <kbd>/</kbd>
    </button>
    <dialog class="search-dialog" @click.self="Close" @close="Reset">
      <form method="dialog" class="search-form" @submit.prevent="First">
        <input
          class="search-input"
          type="search"
          :placeholder="hint"
          :aria-label="hint"
          autofocus
          autocomplete="off"
          spellcheck="false"
          :value="query"
          @input="Query($event.target.value)"
        />
      </form>
      <p v-if="loading" class="search-status">…</p>
      <ul v-else-if="results && results.length" class="search-results">
        <li v-for="r in results" :key="r.link">
          <a :href="r.link" @click="Close">
            <span class="search-titles">{{ r.titles.join(" › ") }}</span>
            <span class="search-snippet"
              ><template v-for="(p, i) in r.snippet" :key="i"><mark v-if="p.mark">{{ p.text }}</mark><template v-else>{{ p.text }}</template></template></span
            >
          </a>
        </li>
      </ul>
      <p v-else-if="query" class="search-status">{{ noResults }} “{{ query }}”</p>
    </dialog>
  </div>
</template>

<script setup lang="go">
// The search dialog. Everything here is Go compiled by goesm: loading the
// index with fetch, ranking and snippets (package press/search), and the
// keyboard shortcut. The index is the plain-text file press writes at build
// time for the page's locale. The same index answers AI agents: with WebMCP,
// the page registers search_docs, get_page and list_pages tools.
import (
	"goesm.dev/client/searchbox"
	"goesm.dev/press/search"
)

type Props struct {
	Index     string
	Label     string
	Hint      string
	NoResults string
	Site      string
	Llms      string
}

label := props.Label
hint := props.Hint
noResults := props.NoResults

box := searchbox.New(props.Index)
query := ""
loading := false
var results []search.Result

func Open() {
	searchbox.ShowModal(".search-dialog")
	loading = true
	box.Load()
	loading = false
	results = box.Search(query)
}

// Preload starts fetching the index when the pointer or focus reaches the
// button, so that on a slow connection it is (partly) there by the click.
func Preload() {
	box.Load()
}

func Query(q string) {
	query = q
	results = box.Search(q)
}

func First() {
	if len(results) > 0 {
		searchbox.Navigate(results[0].Link)
	}
}

func Close() {
	searchbox.Close(".search-dialog")
}

func Reset() {
	query = ""
	results = nil
}

searchbox.OnShortcut(".search-button")
box.RegisterTools(props.Site, props.Llms)
</script>
