//go:build js

// Package enhance adds the page behaviour that needs the browser: copy
// buttons on code blocks and the active entry of the page outline.
package enhance

import "syscall/js"

var installed bool

// Install registers the event listeners once. It does nothing outside a
// browser (during server-side rendering).
func Install(copied string) {
	doc := js.Global().Get("document")
	if installed || !doc.Truthy() {
		return
	}
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
		btn.Call("setAttribute", "aria-label", copied)
		js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) any {
			btn.Get("classList").Call("remove", "copied")
			return nil
		}), 2000)
		return nil
	}))
	trackOutline(doc)
}

// trackOutline marks the outline link of the heading at the top of the
// viewport.
func trackOutline(doc js.Value) {
	links := doc.Call("querySelectorAll", ".outline a")
	n := links.Get("length").Int()
	if n == 0 {
		return
	}
	type entry struct {
		link    js.Value
		heading js.Value
	}
	var entries []entry
	for i := 0; i < n; i++ {
		l := links.Call("item", i)
		id := l.Call("getAttribute", "href").String()[1:]
		if h := doc.Call("getElementById", id); h.Truthy() {
			entries = append(entries, entry{l, h})
		}
	}
	active := -1
	update := js.FuncOf(func(this js.Value, args []js.Value) any {
		cur := -1
		for i, e := range entries {
			if e.heading.Call("getBoundingClientRect").Get("top").Float() < 120 {
				cur = i
			}
		}
		if cur == active {
			return nil
		}
		if active >= 0 {
			entries[active].link.Get("classList").Call("remove", "active")
		}
		if cur >= 0 {
			entries[cur].link.Get("classList").Call("add", "active")
		}
		active = cur
		return nil
	})
	opts := js.Global().Get("Object").New()
	opts.Set("passive", true)
	js.Global().Call("addEventListener", "scroll", update, opts)
	update.Invoke()
}
