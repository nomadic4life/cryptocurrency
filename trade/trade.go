package trade

import (
	"fmt"
	"math"
)

type Trade struct {
	price float64
	size  int64
}

type Position struct {
	Symbol               string
	Order                string // [long, short]
	Size                 int64
	Value                float64
	leverage             float64
	AverageEntryPrice    float64
	LiquidationPrice     float64
	BankruptPrice        float64
	Margin               float64
	InitialMargin        float64
	MaintenaceMargin     float64
	InitialMarginRate    float64
	MaintenaceMarginRate float64
	RiskLevel            int64
	RiskRate             float64
	ActiveOrders         struct {
		entry []Order
		exit  []Order
	}
	Fills struct {
		entry []Fill
		exit  []Fill
	} // [entry, exits]

}

type Fill struct {
	order string // [entry, exit]
	price float64
	size  int64
}

type Order struct {
	order string // [entry, exit]
	Price float64
	Size  int64
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

// type crypto float64

func Open(leverage, price float64, size int64) {
	position := new(OpenPosition)
	position.OpenPrice = price
	position.ContractQuantity = size
	position.OrderValue = Quantity(price, size)
}

func Quantity(price float64, size int64) float64 {
	// price PRICE, size int64 -> Value
	factor := math.Pow(10.0, 8.0)
	return truncate((float64(size) / price), factor)
}

func Amount(price, value float64) float64 {
	// price Price, value Value
	factor := math.Pow(10.0, 2.0)
	return truncate((value * price), factor)
}

func Size(price, value float64) int64 {
	// price Price, value Value
	return int64(Amount(price, value))
}

func Long(entryPrice, exitPrice float64, size int64) [4]float64 {
	entryValue := Quantity(entryPrice, size)
	exitValue := Quantity(exitPrice, size)
	proceeds := closePosition(entryValue, exitValue)
	gross := gross(exitValue, proceeds)

	return [4]float64{proceeds, entryValue, exitValue, gross}
}

func Short(entryPrice, exitPrice float64, size int64) [4]float64 {
	entryValue := Quantity(entryPrice, size)
	exitValue := Quantity(exitPrice, size)
	proceeds := closePosition(exitValue, entryValue)
	gross := gross(entryValue, proceeds)

	// total := Amount(exitPrice, gross)

	return [4]float64{proceeds, entryValue, exitValue, gross}
}

func short(entryValue, exitValue float64) float64 {
	factor := math.Pow(10.0, 8.0)
	return truncate((exitValue - entryValue), factor)
}

func long(entryValue, exitValue float64) float64 {
	factor := math.Pow(10.0, 8.0)
	return truncate((entryValue - exitValue), factor)
}

func closePosition(a, b float64) float64 {
	factor := math.Pow(10.0, 8.0)
	return truncate((a - b), factor)
}

func gross(base, proceeds float64) float64 {
	factor := math.Pow(10.0, 8.0)
	return truncate(base+proceeds, factor)
}

func truncate(a, b float64) float64 {
	// truncate
	return math.Floor(a*b) / b
}

func AveragePrice(orders []Order) {

	var size int64 = 0
	sum := 0.0

	for i := 0; i < len(orders); i++ {
		value := Quantity(orders[i].Price, orders[i].Size)
		sum = math.Floor((sum+value)*math.Pow(10, 8)) / math.Pow(10, 8)
		size += orders[i].Size
	}

	result := math.Floor(float64(size) / sum)
	fmt.Println(result)
}

// structs
//	-> Trade
//	-> Postion
//	-> Record
// 	-> Transaction
//	-> CostBasisRecord
//	-> CostBasisQueue
//	-> OpenPosition
//	-> ClosedPosition

// accrue
// total
// aggregate
// accrual
// return
// yeild
// profit
// gain
// gross

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

// Account
// Accounts
// MasterAccount
// SubAccount
// Order
// PlaceOrder
//	-> limit
//	-> market
//	-> conditional
//	-> Order Value
//	-> Cost (Order Value + Margin (Initial Margin + Maintence Margin))
//	-> Available Balance
//	-> Quantity
//	-> Limit Price
//	-> Buy/Long
//	-> Sell/Short
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
// Contract Details
// Open Position
//	-> Symbol
//	-> Size
//	-> Value
//	-> Entry Price
//	-> Mark Price
//	-> Liq. Price
//	-> Margin
//	-> Unrealized PNL [2][mark price, last traded price]
//	-> Realized PNL
// Closed Position
//	-> Symbol
//	-> Total Size
//	-> Closed PNL
//	-> Exchange Fee Paid
//	-> Funding Fee Paid
//	-> Realized PNL
// Active Orders
//	-> symbol
//	-> QTY.
//	-> Order Price
//	-> Filled/Remaining
//	-> Order Value
//	-> TP/SL
//	-> Fill Price
//	-> Type
//	-> Status
//	-> Time

// Account
//	-> TotalBalance
//	-> AvailableBalance
//	-> UnrealizedPNL
//	-> RiskLimit
//	-> Positions []Position

//	-> OrderHistory []Orders
//	-> fills
//	-> OpenPositions []OpenPosition
//	-> AcitveOrders []ActiveOrder
//	-> ClosedPostions []ClosedPosition
// 	-> PNL

// Position
//	-> size
//	-> value
//	-> price
//	-> leverage
//	-> margin
// 	-> intialMargin
//	-> maintanceMargin
//	-> positionMargin
//	-> riskLimit
//	-> activePostions []entry/open / []exit/close or map[]entry/open / map[]exit/close
//	-> entryTrade []entry
//	-> closeTrade []close
// OpenPostion
//	-> size
//	-> pr
// ClosePostion

// Position
// OpenedPostion
// ClosedPosition
