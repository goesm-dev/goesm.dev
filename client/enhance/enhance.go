//go:build js

// Package enhance adds the page behaviour that needs the browser: copy
// buttons on code blocks and the active entry of the page outline.
package enhance

import "syscall/js"

var installed bool

// Install registers the event listeners once, and finds the page's outline
// on every call: the island calling it starts again on every page the client
// router swaps in. It does nothing outside a browser (during server-side
// rendering).
func Install(copied string) {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return
	}
	if !installed {
		installed = true
		doc.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
			btn := args[0].Get("target").Call("closest", ".vp-code .copy")
			if !btn.Truthy() {
				return nil
			}
			code := btn.Get("parentElement").Call("querySelector", "pre")
			text := code.Get("textContent").String()
			js.Global().Get("navigator").Get("clipboard").Call("writeText", text)
			btn.Get("classList").Call("add", "copied")
			btn.Call("setAttribute", "aria-label", copiedLabel)
			js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) any {
				btn.Get("classList").Call("remove", "copied")
				return nil
			}), 2000)
			return nil
		}))
		opts := js.Global().Get("Object").New()
		opts.Set("passive", true)
		js.Global().Call("addEventListener", "scroll", js.FuncOf(func(this js.Value, args []js.Value) any {
			updateOutline()
			return nil
		}), opts)
	}
	copiedLabel = copied
	trackOutline(doc)
}

// copiedLabel is the copy button's label once it has copied, in the page's
// language.
var copiedLabel string

// The outline links of the page and their headings, and the index of the
// one marked active.
type entry struct {
	link    js.Value
	heading js.Value
}

var (
	entries []entry
	active  = -1
)

// trackOutline finds the outline links of the page and marks the one of the
// heading at the top of the viewport.
func trackOutline(doc js.Value) {
	entries = nil
	active = -1
	links := doc.Call("querySelectorAll", ".outline a")
	n := links.Get("length").Int()
	for i := 0; i < n; i++ {
		l := links.Call("item", i)
		id := l.Call("getAttribute", "href").String()[1:]
		if h := doc.Call("getElementById", id); h.Truthy() {
			entries = append(entries, entry{l, h})
		}
	}
	updateOutline()
}

func updateOutline() {
	cur := -1
	for i, e := range entries {
		if e.heading.Call("getBoundingClientRect").Get("top").Float() < 120 {
			cur = i
		}
	}
	if cur == active {
		return
	}
	if active >= 0 {
		entries[active].link.Get("classList").Call("remove", "active")
	}
	if cur >= 0 {
		entries[cur].link.Get("classList").Call("add", "active")
	}
	active = cur
}
