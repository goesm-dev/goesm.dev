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

// tools are the tools registered so far, by name, and controllers the
// AbortControllers whose signals unregister them.
var (
	tools       = map[string]Tool{}
	controllers = map[string]js.Value{}
)

// Register registers t with the page's WebMCP context, if there is one.
//
// The islands that register tools start again on every page the client
// router swaps in, so a tool may be registered again: the earlier one is
// unregistered first (by aborting its signal, or with unregisterTool where
// the browser has that instead). Where neither works, the tool registered
// first stays, and runs the Run of the latest registration.
func Register(t Tool) {
	mc := Context()
	if !mc.Truthy() {
		return
	}
	_, again := tools[t.Name]
	tools[t.Name] = t
	if again {
		unregister(mc, t.Name)
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
		run := tools[t.Name].Run
		return promise(func() string { return run(input) })
	}))
	controller := js.Global().Get("AbortController").New()
	controllers[t.Name] = controller
	options := object.New()
	options.Set("signal", controller.Get("signal"))
	try(func() { mc.Call("registerTool", tool, options) })
}

func unregister(mc js.Value, name string) {
	if c, ok := controllers[name]; ok {
		c.Call("abort")
		delete(controllers, name)
	}
	if mc.Get("unregisterTool").Type() == js.TypeFunction {
		try(func() { mc.Call("unregisterTool", name) })
	}
}

// try calls f, and drops the exception it throws: a browser that has
// unregistered the tool already may throw for unregisterTool, and one that
// cannot unregister may throw for a second registerTool.
func try(f func()) {
	defer func() { recover() }()
	f()
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
