package main

import (
	// "cryptocurrency/exchange/phemex"

	"cryptocurrency/exchange/number"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

func test() {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	type Secret struct {
		Key string `json:"SECRET"`
	}

	client := make(map[string]interface{})

	json.Unmarshal(byteValue, &client)

	fmt.Println(client)

	data, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		panic("yike")
	}
	fmt.Println(string(data))

	// client.hmac = hmac.New(crypto.SHA256.New, []byte(client.Secret))
}

// "cryptocurrency/tax"

type calc struct {
	factor float64
	number int64 // 3 -> 6 + 1 -> 9, 223,372,036, 854,775,807 -> 9.2 * 10^18 -> 0x7FFFFFFFFFFFFFFF
}

// 21,000,000. 000 000 00 -> 2,100,000,000,000,000  -> 21.0 * 10*14
// 2,100,000,000,000,000  * 1,000,000 = 2,100,000,000,000,000,000,000

//  2,100,000,000,000,000 * 1,000,000.50 [100,000,050]
//  1,000,000.50 * 10^2 = 			100,000,050
//  2,100,000,000,000,000 / 10^8 =  21,000,000
//  21,000,000 * 100,000,050 = 		2,100,001,050,000,000  [21,000,000,000,000]
//  2,100,001,050,000,000 / 10^2 = 	21,000,010,500,000.00

type Number struct {
	numberType string
	scale      int64
	number     int64
	value      string
}

func (Number) Crypto(num int64) *Number {
	number := new(Number)

	scale := int64(math.Pow(10, 8))
	value := fmt.Sprintf("%.8f", float64(num)/float64(scale))

	number.scale = scale
	number.number = num
	number.value = value

	return number
}

func (Number) USD(num int64) *Number {
	number := new(Number)

	scale := int64(math.Pow(10, 2))
	value := fmt.Sprintf("%.2f", float64(num)/float64(scale))

	number.scale = scale
	number.number = num
	number.value = value

	return number
}

func (Number) Ratio(num int64) *Number {
	number := new(Number)

	scale := int64(math.Pow(10, 5))
	value := fmt.Sprintf("%.5f", float64(num)/float64(scale))

	number.scale = scale
	number.number = num
	number.value = value

	return number
}

func (a *Number) Mul(b *Number) *Number {
	// 0.000005 * 45000.50
	// 500 * 450050
	c := new(Number)

	var num int64 = 0

	if a.number < a.scale {
		num = a.number * b.number / a.scale
	} else {

	}
	value := fmt.Sprintf("%.2f", float64(num)/float64(b.scale))

	c.scale = b.scale
	c.number = num
	c.value = value

	return c
}

func big(factor float64) *calc {

	return &calc{factor: factor, number: 0}
}

func (c *calc) add(number int64) *calc {
	c.number += number
	return c
}

func (c *calc) print() {
	fmt.Println(c)
}

func main() {
	num := number.BTC(5000)
	fmt.Println(num)
	// phemex.GetAccounts()
	// phemex.MainSub()
	// phemex.Run()
	// phemex.Handler(phemex.Handle)
	// phemex.Subscribe("trade.subscribe", []interface{}{"BTCUSD"})
	// time.Sleep(4500 * time.Second)
	// body := map[string]interface{}{
	// 	// "actionBy":         "FromOrderPlacement",
	// 	"symbol":         "BTCUSD",
	// 	"side":           "Buy",
	// 	"clOrdID":        "2b6ba6d3-a14c-44ef-be5b-0e590ab35126",
	// 	"ordType":        "Limit",
	// 	"reduceOnly":     false,
	// 	"closeOnTrigger": false,
	// 	"timeInForce":    "GoodTillCancel",
	// 	"priceEp":        93185000,
	// 	// "triggerType":      "UNSPECIFIED",
	// 	// "pegPriceType":     "UNSPECIFIED",
	// 	// "takeProfitEp":     0,
	// 	// "stopLossEp":       0,
	// 	// "pegOffsetValueEp": 0,
	// 	"orderQty": 1}

	// query := map[string]string{
	// "currency": "BTC"}

	// var query map[string]interface{}
	// query = make(map[string]interface{})
	// query["priceEp"] = "53500000"
	// query["side"] = "Sell"
	// id := "5bfd51f1-1f65-4031-ba45-cf2bc17c49f3"
	// id := ""

	// phemex.Send("GET", "/exchange/public/products", nil, nil).HandleResponse(phemex.JSON).Display()
	// phemex.Send("GET", "/phemex-user/users/children", query, nil).HandleResponse(phemex.JSON).Display()
	// phemex.CreateOrder("Buy", 1, 45500000, query).Display()
	// phemex.AmendOrder(id, query).Display()
	// phemex.GetOrders(id, "").Display()
	// phemex.CancelOrders([]string{id}, "").Display()

	// res.Display()

	// value := Number{}.USD(85000000)
	// fmt.Println(value)

	// a := Number{}.Crypto(500)
	// b := Number{}.USD(4500050)
	// c := a.Mul(b)
	// fmt.Println(c)

	// big(10.0).add(20).add(20).print()

	// account := phemex.CreateTradeAccount("BTCUSD")
	// account.SetLeverage(1)
	// account.SetBalance(1)
	// fmt.Println(account)
	// account.GetAccount()
	// account.Entry("Short", 100, 500.0)
	// account.GetAccount()
	// account.Exit(100, 400.0)
	// account.GetAccount()

	// a := phemex.CreateTradeAccount("BTCUSD")
	// a.SetLeverage(2)
	// a.SetBalance(20)
	// fmt.Println(a)
	// a.GetAccount()
	// account.CalcMaxMargin()

	// account := tax.NewAccount(1000)

	// fmt.Println("Account:", account)

	// account.CreateTransaction(tax.TradeInput{
	// 	Pair:     "BTC/USD",
	// 	Type:     "BUY",
	// 	Price:    20.0,
	// 	Quantity: 10.0})

	// account.CreateTransaction(tax.TradeInput{
	// 	Pair:     "BTC/USD",
	// 	Type:     "BUY",
	// 	Price:    10.0,
	// 	Quantity: 10.0})

	// account.CreateTransaction(tax.TradeInput{
	// 	Pair:     "BTC/USD",
	// 	Type:     "BUY",
	// 	Price:    5.00,
	// 	Quantity: 20.00})

	// account.CreateTransaction(tax.TradeInput{
	// 	Pair:  "LBC/BTC",
	// 	Type:  "BUY",
	// 	Price: 0.10,
	// 	// Quantity: 10, // give correct results
	// 	Amount: 5.0, // seems not bugged anymore
	// 	Value:  10})

	// account.CreateTransaction(tax.TradeInput{
	// 	Pair:  "LBC/BTC",
	// 	Type:  "BUY",
	// 	Price: 0.05,
	// 	// Quantity: 10, // give correct results
	// 	Amount: 10.0, // seems not bugged anymore
	// 	Value:  40})

	// account.CreateTransaction(tax.TradeInput{
	// 	Pair:     "BTC/USD",
	// 	Type:     "SELL",
	// 	Price:    40.0,
	// 	Quantity: 10.0})

	// account.CreateTransaction(tax.TradeInput{
	// 	Pair:     "LBC/BTC",
	// 	Type:     "SELL",
	// 	Price:    0.20,
	// 	Quantity: 5, // give correct results
	// 	// Amount: 1.0, // seems not bugged anymore
	// 	Value: 80})

	// account.CreateTransaction(tax.TradeInput{
	// 	Pair:     "BTC/USD",
	// 	Type:     "SELL",
	// 	Price:    160.00,
	// 	Quantity: 14.0})

	// phemex.Run(31929.4, 34303.5, 3.0, 0.5, 18000)
	// phemex.Run(31929.4, 64000, 5.0, 1000.0, 18000)
	// phemex.Run(8719, 32000, 5.0, 50.0, 375)
	// phemex.Run(5354, 15000, 5.0, 5.0, 170)

	// phemex.Run(31929.4, 32000, 3.0, 0.5, 18000) [1.0635 8767] 1.1026 0312 -> 0.1102 6031 -> 0.2205 2062

	//1,338.52 -> 390
}

// field size
// padding
// minPadding
// maxPadding
// margin
// cell
// field
// padLeft
// padRight

// column
// row
// header
// field
// cell
