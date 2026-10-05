<!-- Synced by cmd/syncdocs. Edit the source instead. -->

* Like Vue's `<script setup>`, the top level runs once, top to bottom, per component instance. You can write `x := ...`, `var`, `const`, `type`, and `func F() {...}`.
* Top-level variables, constants, and functions are available in the template. When the template calls a Go function (for example `@click="Increment"`), the rendered output follows the Go values.
* Only Go imports are allowed. You cannot import `.vue`, `.ts`, or `.go` files.
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
