package number

import (
	"fmt"
	"math"
)

type Number struct { // Currency?
	scale    int64
	value    int64
	currency string // Symbol?
	// category string // ["PRICE", "VALUE", "QUANTITY"]
	// market string // ["BTCUSD"]
}

func NewNumber(value, scale int64, currency string) *Number {
	num := new(Number)
	num.value = value
	num.scale = scale
	num.currency = currency
	return num
}

func Crypto(num int64, currency string) *Number {
	return NewNumber(num, 8, currency)
}

func BTC(num int64) *Number {
	return Crypto(num, "BTC")
}

func Fiat(num int64, currency string) *Number {
	return NewNumber(num, 4, currency)
}

func USD(num int64) *Number {
	return Fiat(num, "USD")
}

func (n *Number) String() string {
	factor := int64(math.Pow(10, float64(n.scale)))
	a := n.value / factor
	b := n.value - (a * factor)
	symbol := ""
	fmt.Println(a, b, n.value, a*factor, n.value-a*factor)
	if n.currency == "USD" {
		symbol = "$"
		return fmt.Sprintf("%v %d.%-02d", symbol, a, (b / 100))
	} else if n.currency == "BTC" {
		symbol = "BTC"
		return fmt.Sprintf("%v %d.%0*d", symbol, a, n.scale, b)

	} else {
		return fmt.Sprintf("%v %d.%0*d", symbol, a, n.scale, b)
	}
}

func (n *Number) Add(m *Number) *Number {
	return n
}

func (n *Number) Sub(m *Number) *Number {
	return n
}

func (n *Number) Mul(m *Number) *Number {
	return n
}

func (n *Number) Div(m *Number) *Number {
	return n
}

func (n *Number) Mod(m *Number) *Number {
	return n
}
