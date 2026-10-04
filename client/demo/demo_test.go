package demo

import "testing"

func TestCart(t *testing.T) {
	c := NewCart("en")
	if got := Total(c.Items); got != 360+250+996 {
		t.Errorf("total %d", got)
	}
	if got := Discount(2408, 15); got != 2046 {
		t.Errorf("discount %d", got)
	}
	c.Add(0, -10)
	c.SetPercent(150)
	if c.Items[0].Quantity != 0 || c.Percent != 100 {
		t.Errorf("%+v", c)
	}
}

func TestFormula(t *testing.T) {
	g, j := Formula(2408, 15)
	if g != "2408 * 85 / 100 = 2046" || j != "2046.8" {
		t.Errorf("%q %q", g, j)
	}
	if _, j := Formula(1000, 15); j != "850" {
		t.Errorf("%q", j)
	}
	if _, j := Formula(1001, 15); j != "850.85" {
		t.Errorf("%q", j)
	}
	if Atoi("42") != 42 || Atoi("x") != 0 {
		t.Error("Atoi")
	}
}
