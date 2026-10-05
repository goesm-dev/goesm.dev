<!-- Synced by cmd/syncdocs. Edit the source instead. -->

* Like Vue's `<script setup>`, the top level runs once, top to bottom, per component instance. You can write `x := ...`, `var`, `const`, `type`, and `func F() {...}`.
* Top-level variables, constants, and functions are available in the template. When the template calls a Go function (for example `@click="Increment"`), the rendered output follows the Go values.
* `import` takes Go packages only; the Go frontmatter of an `.astro` file can also import JavaScript modules, as described below. In a Vue component, Vue components, TypeScript and JavaScript come in through goesm's [`//goesm:import`](https://github.com/goesm-dev/goesm/blob/main/docs/js-imports.md) directive, before a `var` (a component, a class, any value) or a function without a body (a function, called with automatic conversion of its arguments and result). Both are template bindings like the block's other names. Their types may use imported and built-in types, not types declared in the block, since they are declared at package level. The directive needs goesm v0.0.1-beta.2 or later, the version in your go.mod; an older goesm does not recognize it.

  ```vue
  <template>
    <p>{{ formatPrice(12800) }} <Badge label="new" /></p>
  </template>

  <script setup lang="go">
  import "syscall/js"

  //goesm:import "./Badge.vue"
  var Badge js.Value

  //goesm:import "./format.ts" formatPrice
  func formatPrice(yen int) string
  </script>
  ```
* Put methods and generic functions in a Go package.
* To receive props, declare `type Props struct {...}`; the block then has a `props` variable of that type. Field `Route` is read from the attribute `route` (or the field's json tag name, or its kebab-case form), converted to the field's type (string, bool, integer or float kinds). Props are read once, when the instance is set up, and do not fall through to the root element.

  ```vue
  <script setup lang="go">
  import "strings"

  type Props struct {
  	Title string
  }

  heading := strings.ToUpper(props.Title)
  </script>
  ```
