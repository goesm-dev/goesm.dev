//go:build js

// Package searchbox is the client side of the search dialog: it loads a
// locale's index (written at build time by press) on first use and queries
// it with package search.
package searchbox

import (
	"syscall/js"

	"goesm.dev/press/search"
)

// Box is one search dialog.
type Box struct {
	url   string
	index *search.Index
}

func New(indexURL string) *Box { return &Box{url: indexURL} }

// Load fetches the index once. It blocks until the response arrives (goesm
// turns it, and its callers, into async functions).
func (b *Box) Load() {
	if b.index != nil {
		return
	}
	b.index = search.Parse(fetchText(b.url))
}

func (b *Box) Search(q string) []search.Result {
	if b.index == nil {
		return nil
	}
	return b.index.Search(q, 20)
}

// fetchText GETs url with the browser's fetch and waits for the body.
func fetchText(url string) string {
	ch := make(chan string, 1)
	var onText, onResp, onErr js.Func
	onText = js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- args[0].String()
		return nil
	})
	onErr = js.FuncOf(func(this js.Value, args []js.Value) any {
		ch <- ""
		return nil
	})
	onResp = js.FuncOf(func(this js.Value, args []js.Value) any {
		args[0].Call("text").Call("then", onText, onErr)
		return nil
	})
	js.Global().Call("fetch", url).Call("then", onResp, onErr)
	s := <-ch
	onText.Release()
	onResp.Release()
	onErr.Release()
	return s
}

// OnShortcut clicks the element matching selector when "/" or Ctrl-K /
// Cmd-K is pressed outside a text field. It does nothing outside a browser.
func OnShortcut(selector string) {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return
	}
	doc.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) any {
		e := args[0]
		key := e.Get("key").String()
		mod := e.Get("ctrlKey").Bool() || e.Get("metaKey").Bool()
		tag := e.Get("target").Get("tagName").String()
		editing := tag == "INPUT" || tag == "TEXTAREA" || e.Get("target").Get("isContentEditable").Bool()
		if (key == "k" && mod) || (key == "/" && !editing && !mod) {
			e.Call("preventDefault")
			Click(selector)
		}
		return nil
	}))
}

// Click clicks the first element matching selector.
func Click(selector string) {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Call("click")
	}
}

// ShowModal opens the <dialog> matching selector.
func ShowModal(selector string) {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() && !el.Get("open").Bool() {
		el.Call("showModal")
	}
}

// Close closes the <dialog> matching selector.
func Close(selector string) {
	if el := js.Global().Get("document").Call("querySelector", selector); el.Truthy() {
		el.Call("close")
	}
}

// Navigate goes to url.
func Navigate(url string) {
	js.Global().Get("location").Set("href", url)
}
