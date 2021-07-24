package phemex

// import (
// 	"math"
// )

// type Trade struct {
// 	price float64
// 	size  int64
// }

// type Position struct {
// 	Symbol               string
// 	Order                string // [long, short]
// 	Size                 int64
// 	Value                float64
// 	leverage             float64
// 	AverageEntryPrice    float64
// 	LiquidationPrice     float64
// 	BankruptPrice        float64
// 	Margin               float64
// 	InitialMargin        float64
// 	MaintenaceMargin     float64
// 	InitialMarginRate    float64
// 	MaintenaceMarginRate float64
// 	RiskLevel            int64
// 	RiskRate             float64
// 	ActiveOrders         struct {
// 		entry []Order
// 		exit  []Order
// 	}
// 	Fills struct {
// 		entry []Fill
// 		exit  []Fill
// 	} // [entry, exits]

// }

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

// // on program boot up
// // 	-> pull data from Database that is relevent to account info and active trades
// //	-> pull data from exchange that is relevent to acocunt info and active trades
// //	-> compare data from database to exchange
// //	-> update with current/correct account data information
// //	-> track and monitor relevent market symbol data
// //	-> post trades when an event is triggered

// // processes of program
// //	-> execute trades through SDK to Exchange API
// // 	-> set triggers and automate trades, configure events algorithms
// //	-> track and monitor exhange data of symbols
// //	-> data processing, data anylsis of historic data and monitored data
// //	-> store data -> account data, tracked data, events, proccessed data, anylized data, historic data
// //	-> front end to interface with backend/daemon
// //	->	-> display data in a meaningful way, layout, graph visualized
// //	-> 	-> interact -> set trades, triggers, deposit, withdraw, manage
// //	-> 	-> manage user account info
