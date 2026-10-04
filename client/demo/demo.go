// Package demo is the cart shown on the home page. It runs in the browser,
// compiled by goesm, so the numbers on the page are computed by Go: integer
// arithmetic truncates like Go's, not like JavaScript's floats.
package demo

// Item is a line of the cart.
type Item struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

// Total is the sum of price × quantity.
func Total(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Price * item.Quantity
	}
	return total
}

// Discount takes percent off total, with Go's integer division.
func Discount(total, percent int) int {
	return total * (100 - percent) / 100
}

// Cart is the demo's state.
type Cart struct {
	Items   []Item `json:"items"`
	Percent int    `json:"percent"`
}

// NewCart returns the starting cart with names in the given language.
func NewCart(lang string) *Cart {
	names := [3]string{"Apple", "Bread", "Coffee"}
	if lang == "ja" {
		names = [3]string{"りんご", "パン", "コーヒー"}
	}
	return &Cart{
		Items:   []Item{{names[0], 120, 3}, {names[1], 250, 1}, {names[2], 498, 2}},
		Percent: 15,
	}
}

// Add changes the quantity of item i by delta, never below zero.
func (c *Cart) Add(i, delta int) {
	if i < 0 || i >= len(c.Items) {
		return
	}
	c.Items[i].Quantity = max(0, c.Items[i].Quantity+delta)
}

// SetPercent sets the discount, clamped to 0–100.
func (c *Cart) SetPercent(p int) { c.Percent = min(100, max(0, p)) }

// Formula describes the discount computation with Go's integer division,
// and what JavaScript's floating-point division would give instead.
func Formula(total, percent int) (goExpr, jsValue string) {
	x := total * (100 - percent)
	goExpr = itoa(total) + " * " + itoa(100-percent) + " / 100 = " + itoa(x/100)
	jsValue = itoa(x / 100)
	if frac := x % 100; frac != 0 {
		d := itoa(frac + 100)[1:] // two digits
		if d[1] == '0' {
			d = d[:1]
		}
		jsValue += "." + d
	}
	return goExpr, jsValue
}

// Atoi parses a non-negative decimal number, or returns 0.
func Atoi(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0
		}
		n = n*10 + int(s[i]-'0')
	}
	return n
}

// itoa is strconv.Itoa for non-negative numbers. The demo does not import
// strconv, whose tables would make the island larger.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
