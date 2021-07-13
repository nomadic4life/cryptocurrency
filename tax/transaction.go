package tax

import (
	"fmt"
	"math"
	"strings"
)

func newTransaction(account *Account, input TradeInput) *TransactionEntry {
	transaction := new(TransactionEntry)

	transaction.OrderPair = input.Pair
	transaction.OrderType = input.Type
	transaction.OrderPrice = input.Price

	transaction.TransactionID = account.getID()
	transaction.Date = getDate(input.Date)

	transaction.OrderQuantity = calcQuantity(input.Price, input.Quantity, input.Amount)
	transaction.OrderAmount = calcAmount(input.Price, input.Quantity, input.Amount, transaction.quote())
	transaction.USDPriceValue = getUSDPrice(input.Price, input.Value, transaction.quote())

	transaction.FeeAmount = calcFee(input.Price, input.Quantity, input.Amount, input.Fee)

	return transaction
}

func calcQuantity(price, quantity, amount float64) float64 {
	if quantity == 0 {
		return math.Floor(amount/price*math.Pow(10, 8)) / math.Pow(10, 8)
	} else {
		return quantity
	}
}

func calcAmount(price, quantity, amount float64, quoteSymbol string) float64 {
	if amount == 0 && quoteSymbol == "USD" {
		return math.Floor(quantity*price*100) / 100

	} else if amount == 0 && quoteSymbol != "USD" {
		return math.Floor(quantity*price*math.Pow(10, 8)) / math.Pow(10, 8)

	} else {
		return amount
	}
}

func calcFee(price, quantity, amount, fee float64) float64 {
	if fee == 0 {
		// calculate fee amount
		return 0
	} else {
		return fee
	}
}

func getDate(date int64) int64 {
	if date == 0 {
		// get from api
		// return from api
		return 0

	} else {
		return date

	}
}

func getUSDPrice(price, value float64, quote string) float64 {

	if quote == "USD" {
		return price

	} else if value == 0.0 {
		// get from api
		// t.USDPriceValue =
		return value

	} else {
		return value

	}
}

func (t *TransactionEntry) quote() string {
	return strings.Split(t.OrderPair, "/")[1]
}

func (t *TransactionEntry) base() string {
	return strings.Split(t.OrderPair, "/")[0]
}

func (e *TransactionEntry) filter(properties []string) []string {
	var t table = make(map[string]string)
	t = map[string]string{
		"Transaction ID":  fmt.Sprint(e.TransactionID),
		"Order Date":      fmt.Sprint(e.Date),
		"Order Pair":      fmt.Sprint(e.OrderPair),
		"Order Type":      fmt.Sprint(e.OrderType),
		"Order Price":     fmt.Sprint(e.OrderPrice),
		"Order Quantity":  fmt.Sprint(e.OrderQuantity),
		"Order Amount":    fmt.Sprint(e.OrderAmount),
		"USD Price Value": fmt.Sprint(e.USDPriceValue),
		"Fee Amount":      fmt.Sprint(e.FeeAmount)}

	return t.filter(properties)
}

// func (t *TransactionEntry) enqueue(transaction *TransactionEntry) {
//   t = append(t, *transaction)
// 	// a.Ledger.CostBases = append(a.Ledger.CostBases, log.ledger.costBases...)
// }

// TODO::
// api to get USD price
// api to get date
// calculate fee
// implement calcQuantity, calcAmount, calcFee as a method?
// change file name to transactionEntry.go
