<template>
  <section class="demo">
    <div class="demo-text">
      <h2>{{ t.title }}</h2>
      <p>{{ t.lead }}</p>
      <p class="demo-note">{{ t.note }} <code>{{ exact }}</code> · JavaScript: <code>{{ float }}</code></p>
    </div>
    <div class="demo-card">
      <table>
        <tbody>
          <tr v-for="(item, i) in cart.items" :key="i">
            <th scope="row">{{ item.name }}</th>
            <td class="num">{{ item.price }}</td>
            <td class="qty">
              <button type="button" :aria-label="'− ' + item.name" @click="Add(i, -1)">−</button>
              <span>{{ item.quantity }}</span>
              <button type="button" :aria-label="'+ ' + item.name" @click="Add(i, 1)">+</button>
            </td>
          </tr>
        </tbody>
      </table>
      <label class="demo-discount">
        <span>{{ t.discount }} {{ cart.percent }}%</span>
        <input type="range" min="0" max="50" :value="cart.percent" @input="SetPercent($event.target.value)" />
      </label>
      <dl class="demo-totals">
        <div><dt>Total(items)</dt><dd>{{ total }}</dd></div>
        <div><dt>Discount(total, {{ cart.percent }})</dt><dd>{{ discounted }}</dd></div>
      </dl>
    </div>
  </section>
</template>

<script setup lang="go">
// Go running in the browser: this component's state and arithmetic are the
// Go package goesm.dev/client/demo, compiled by goesm. Every click calls a Go
// function, and the template shows Go's values.
import "goesm.dev/client/demo"

type Props struct {
	Lang string
}

type text struct {
	Title    string `json:"title"`
	Lead     string `json:"lead"`
	Note     string `json:"note"`
	Discount string `json:"discount"`
}

t := text{
	Title:    "Go, running in your browser",
	Lead:     "This cart is a Go package compiled by goesm. The buttons call Go functions; the totals are Go ints.",
	Note:     "Integer division stays Go's:",
	Discount: "Discount",
}
if props.Lang == "ja" {
	t = text{
		Title:    "ブラウザで動いている Go",
		Lead:     "このカートは goesm でコンパイルした Go のパッケージです。ボタンは Go の関数を呼び、合計は Go の int で計算しています。",
		Note:     "整数の割り算は Go のまま:",
		Discount: "割引",
	}
}

cart := demo.NewCart(props.Lang)
total := 0
discounted := 0
exact := ""
float := ""

func update() {
	total = demo.Total(cart.Items)
	discounted = demo.Discount(total, cart.Percent)
	exact, float = demo.Formula(total, cart.Percent)
}
update()

func Add(i, delta int) {
	cart.Add(i, delta)
	update()
}

func SetPercent(p string) {
	cart.SetPercent(demo.Atoi(p))
	update()
}
</script>
