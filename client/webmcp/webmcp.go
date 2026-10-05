//go:build js

// Package webmcp registers tools with the browser's WebMCP API
// (document.modelContext.registerTool), so that an AI agent working in the
// page can call them instead of reading the screen. Browsers without WebMCP
// (and server-side rendering) get nothing.
//
// Like package search, it imports no other Go package than syscall/js, to
// keep the island that uses it small.
package webmcp

import "syscall/js"

// Tool is one tool. Run gets the arguments the agent passed (an object
// matching Schema) and returns the text the agent reads; a failure is a
// sentence for the agent, too. Run may block: it runs in its own goroutine.
type Tool struct {
	Name        string
	Title       string
	Description string
	Schema      string // the JSON Schema of the input
	Run         func(input js.Value) string
}

// Context returns the page's WebMCP context, or a falsy value without one.
func Context() js.Value {
	doc := js.Global().Get("document")
	if !doc.Truthy() {
		return js.Undefined()
	}
	if mc := doc.Get("modelContext"); mc.Truthy() {
		return mc
	}
	// The older name, still in browsers that shipped before document.modelContext.
	return js.Global().Get("navigator").Get("modelContext")
}

// Register registers t with the page's WebMCP context, if there is one.
func Register(t Tool) {
	mc := Context()
	if !mc.Truthy() {
		return
	}
	object := js.Global().Get("Object")
	tool := object.New()
	tool.Set("name", t.Name)
	tool.Set("title", t.Title)
	tool.Set("description", t.Description)
	tool.Set("inputSchema", js.Global().Get("JSON").Call("parse", t.Schema))
	annotations := object.New()
	annotations.Set("readOnlyHint", true)
	tool.Set("annotations", annotations)
	tool.Set("execute", js.FuncOf(func(this js.Value, args []js.Value) any {
		input := js.Undefined()
		if len(args) > 0 {
			input = args[0]
		}
		return promise(func() string { return t.Run(input) })
	}))
	mc.Call("registerTool", tool)
}

// promise runs f in a goroutine and returns a Promise of its text.
func promise(f func() string) js.Value {
	var executor js.Func
	executor = js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		go func() {
			resolve.Invoke(f())
			executor.Release()
		}()
		return nil
	})
	return js.Global().Get("Promise").New(executor)
}

// String returns the string property name of v, or "" when it is missing or
// not a string.
func String(v js.Value, name string) string {
	if v.Type() != js.TypeObject {
		return ""
	}
	if p := v.Get(name); p.Type() == js.TypeString {
		return p.String()
	}
	return ""
}

// Int returns the number property name of v, or def when it is missing or
// not a number.
func Int(v js.Value, name string, def int) int {
	if v.Type() != js.TypeObject {
		return def
	}
	if p := v.Get(name); p.Type() == js.TypeNumber {
		return p.Int()
	}
	return def
}
