package exchange

import (
	"cryptocurrency/trade"
	"math"
)

type Position struct { // order
	Symbol    string
	Quoted    string
	Settled   string
	RiskLevel int
	OpenPosition
	ClosedPosition
}

type Record struct {
	price    float64
	Size     int64
	Quantity float64
}

type ClosedPosition struct {
	// dollar
	EntryPrice float64
	// crypto
	EntryValue float64
	// dollar
	ExitPrice float64
	// crypto
	ExitValue float64

	// crypto
	Proceeds float64 // value that is earned (gain/loss)
	// crypto
	Gross float64 // entryValue + proceeds

	// dollar
	proceedsAmount float64 // dollar value at exit price
	// dollar
	grossAmount float64 // dollar value at exit price

	// rate
	GrossYield float64
	// dollar
	priceDifference float64
	// rate
	pricePercentage float64
	// dollar
	Size int64
	// dollar
	PNL float64
	// dollar
	Total float64
	// rate
	TotalYield float64
}

type OpenPosition struct {
	OpenPrice        float64
	ContractQuantity int64
	OrderValue       float64
	PositionValue    float64
	Leverage         float64
	InitialMargin    float64
	MaintenaceMargin float64
	Margin           float64
}

type Price float64
type Value float64

func Long(entryPrice, exitPrice Price, size int64) ClosedPosition {
	const (
		PROCEEDS   = 0
		ENTRYVALUE = 1
		EXITVALUE  = 2
	)
	factor := math.Pow(10.0, 8.0)
	difference := exitPrice.subtract(entryPrice)
	percentage := difference.divide(entryPrice)
	// percentage := math.Floor(difference/entryPrice*100000) / 1000
	pos := trade.Long(entryPrice, exitPrice, size)
	gross := math.Floor(pos[PROCEEDS]*factor+pos[ENTRYVALUE]*factor) / factor
	yield := math.Floor(pos[PROCEEDS]/pos[ENTRYVALUE]*100000) / 100000
	pnl := trade.Amount(exitPrice, pos[PROCEEDS])

	total := trade.Amount(exitPrice, gross)
	totalYield := math.Floor((total-float64(size))/float64(size)*100000) / 100000
	// totalYield := math.Floor((total-float64(size))/float64(size)*100000) / 1000

	return ClosedPosition{
		entryPrice,
		pos[ENTRYVALUE],
		exitPrice,
		pos[EXITVALUE],
		pos[PROCEEDS],
		gross,
		yield,
		difference,
		percentage,
		size,
		pnl,
		total,
		totalYield}
}

func (a Price) subtract(b Price) Price {
	return subtract(float64(a), float64(b), 2)
}

func (a Price) divide(b Price) float64 {
	factor := math.Pow(10, 5)
	return math.Floor(float64(a)/float64(b)*factor) / factor
}

func (a Value) subtract(b Value) Value {
	return subtract(float64(a), float64(b), 8)
}

func subtract(a, b, pow float64) float64 {
	factor := math.Pow(10, pow)
	return math.Floor(a*factor-b*factor) / factor
}

func (p *Position) test() {}

func truncate(a, b float64) float64 {
	// truncate
	return math.Floor(a*b) / b
}

// long ->
//      -> openValue
//      -> closeValue
//      -> openQty
//      -> closeQty
//      -> returnValue
//      ->
//      -> openPrice
//      -> closePrice
//      -> priceDifference
//      -> pricePercentage
//      -> size
//      ->

// symbol -> Minimum Price Increment
// symbol -> RiskRate
// symbol -> inverse or linear

// BTCUSD

// (BTC)

// 100 150 200 250 300 350 400 450 500 550

// IM

// 1% 1.5% 2% 2.5% 3% 3.5% 4% 4.5% 5% 5.5%

// MM

// 0.5% 1% 1.5% 2% 2.5% 3% 3.5% 4% 4.5% 5%

// Exchange
//	-> Bettrix
// 	-> Phemex
//	-> Coinex
//	-> CoinBase

//	-> Binance
//	-> ByBit
//	-> Gemini
//	-> Kraken
//	-> Bitstamp

//	-> OKCoin
//	-> KuCoin
//	-> Bitfinex
//	-> BitMex
//	-> FTX
//	-> CEX.io
//	-> decentralized exchanges

// // type crypto float64

// func Open(leverage, price float64, size int64) {
// 	position := new(OpenPosition)
// 	position.OpenPrice = price
// 	position.ContractQuantity = size
// 	position.OrderValue = Value(price, size)
// }

// func Value(price float64, size int64) float64 {
// 	// price PRICE, size int64 -> Value
// 	factor := math.Pow(10.0, 8.0)
// 	return truncate((float64(size) / price), factor)
// }

// func Amount(price, value float64) float64 {
// 	// price Price, value Value
// 	factor := math.Pow(10.0, 2.0)
// 	return truncate((value * price), factor)
// }

// func CountractSize(price, value float64) int64 {
// 	// price Price, value Value
// 	return int64(Amount(price, value))
// }

// func Long(entryPrice, exitPrice float64, size int64) [4]float64 {
// 	entryValue := Value(entryPrice, size)
// 	exitValue := Value(exitPrice, size)
// 	proceeds := exit(entryValue, exitValue)
// 	gross := Total(exitValue, proceeds)

// 	return [4]float64{proceeds, entryValue, exitValue, gross}
// }

// func Short(entryPrice, exitPrice float64, size int64) [4]float64 {
// 	entryValue := Value(entryPrice, size)
// 	exitValue := Value(exitPrice, size)
// 	proceeds := exit(exitValue, entryValue)
// 	gross := Total(entryValue, proceeds)

// 	// total := Amount(exitPrice, gross)

// 	return [4]float64{proceeds, entryValue, exitValue, gross}
// }

// func short(entryValue, exitValue float64) float64 {
// 	factor := math.Pow(10.0, 8.0)
// 	return truncate((exitValue - entryValue), factor)
// }

// func long(entryValue, exitValue float64) float64 {
// 	factor := math.Pow(10.0, 8.0)
// 	return truncate((entryValue - exitValue), factor)
// }

// func exit(a, b float64) float64 {
// 	factor := math.Pow(10.0, 8.0)
// 	return truncate((a - b), factor)
// }

// func Total(base, proceeds float64) float64 {
// 	factor := math.Pow(10.0, 8.0)
// 	return truncate(base+proceeds, factor)
// }

// // func truncate(a, b float64) float64 {
// // 	// truncate
// // 	return math.Floor(a*b) / b
// // }

// func AveragePrice(orders []Order) float64 {

// 	size := 0.0
// 	sum := 0.0

// 	for i := 0; i < len(orders); i++ {
// 		value := Value(orders[i].Price, orders[i].Size)
// 		sum = math.Floor((sum+value)*math.Pow(10, 8)) / math.Pow(10, 8)
// 		size += float64(orders[i].Size)
// 	}

// 	return math.Floor(size/sum*10) / 10
// }

// func findNearestIntervial(price float64) float64 {
// 	cent := math.Floor((math.Ceil(price) - price) * 100)

// 	if 25 <= cent && cent < 75 {
// 		return ((math.Floor(price) * 100) + 50) / 100
// 	}

// 	return math.Round(price)
// }

// // structs
// //	-> Trade
// //	-> Postion
// //	-> Record
// // 	-> Transaction
// //	-> CostBasisRecord
// //	-> CostBasisQueue
// //	-> OpenPosition
// //	-> ClosedPosition

// // accrue
// // total
// // aggregate
// // accrual
// // return
// // yeild
// // profit
// // gain
// // gross

// Trade Account
//	-> BTC Account Balance
//		-> Unrealized PNL
//		-> Margin Balance
//		-> Position Margin
//		-> Order Margin
//		-> Available Balance
//	-> USD Account Balance
//		-> Unrealized PNL
//		-> Margin Balance
//		-> Position Margin
//		-> Order Margin
//		-> Available Balance

// type Fill struct {
// 	order string // [entry, exit]
// 	price float64
// 	size  int64
// }

// type Order struct {
// 	order string // [entry, exit]
// 	Price float64
// 	Size  int64
// }

// type OpenPosition struct {
// 	OpenPrice        float64
// 	ContractQuantity int64
// 	OrderValue       float64
// 	PositionValue    float64
// 	Leverage         float64
// 	InitialMargin    float64
// 	MaintenaceMargin float64
// 	Margin           float64
// }
