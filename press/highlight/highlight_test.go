package highlight

import (
	"strings"
	"testing"
)

func TestCode(t *testing.T) {
	cases := []struct {
		lang, src string
		want      []string
	}{
		{"go", "func Total(items []Item) int {\n\treturn 0 // zero\n}", []string{`<span class="hl-kw">func</span> <span class="hl-fn">Total</span>`, `<span class="hl-type">int</span>`, `<span class="hl-com">// zero</span>`, `<span class="hl-num">0</span>`}},
		{"ts", `import { Total } from "./cart.ts";`, []string{`<span class="hl-kw">import</span>`, `<span class="hl-str">&quot;./cart.ts&quot;</span>`}},
		{"sh", "go tool goesm emit-ts -o out ./cart # build", []string{`<span class="hl-fn">go</span>`, `<span class="hl-attr">-o</span>`, `<span class="hl-com"># build</span>`}},
		{"vue", "<template>\n  <div>{{ total }}</div>\n</template>\n<script setup lang=\"go\">\nx := 1\n</script>", []string{`<span class="hl-tag">template</span>`, `<span class="hl-attr">lang</span>=<span class="hl-str">&quot;go&quot;</span>`, `<span class="hl-num">1</span>`}},
		{"toml", "[tools]\ngo = \"1.27.1\"", []string{`<span class="hl-type">[tools]</span>`, `<span class="hl-attr">go </span>`}},
		{"json", `{"type": "module", "n": 1}`, []string{`<span class="hl-attr">&quot;type&quot;</span>`, `<span class="hl-str">&quot;module&quot;</span>`}},
		{"text", "<a>", []string{"&lt;a&gt;"}},
	}
	for _, c := range cases {
		got := Code(c.lang, c.src)
		for _, w := range c.want {
			if !strings.Contains(got, w) {
				t.Errorf("%s: %q\nmissing %q", c.lang, got, w)
			}
		}
	}
}
