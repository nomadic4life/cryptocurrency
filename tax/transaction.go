package tax

import (
	"math"
)

func newTransaction(account *Account, trade Trade) *Transaction {
	t := Transaction{}
	t.OrderPair = trade.Pair
	t.OrderType = trade.Type
	t.OrderPrice = trade.Price
	t.TransactionID = account.getID()
	t.Date = getDate(trade.Date)
	t.OrderQuantity = calcQuantity(trade.Price, trade.Quantity, trade.Amount)
	t.OrderAmount = calcAmount(trade.Price, trade.Quantity, trade.Amount)
	t.USDPriceValue = getUSDPrice(trade.Price, trade.Value, trade.Pair)
	t.FeeAmount = calcFee(trade.Price, trade.Quantity, trade.Amount, trade.Fee)
	return &t
}

func calcQuantity(price, quantity, amount float64) float64 {
	if quantity == 0 {
		return math.Floor(amount/price*math.Pow(10, 8)) / math.Pow(10, 8)
	} else {
		return quantity
	}
}

func calcAmount(price, quantity, amount float64) float64 {
	if amount == 0 {
		return quantity * price
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

func getUSDPrice(price, value float64, pair string) float64 {
	if pair == "BTC/USD" {
		return price
	} else if value == 0.0 {
		// get from api
		// t.USDPriceValue =
		return value
	} else {
		return value
	}
}
