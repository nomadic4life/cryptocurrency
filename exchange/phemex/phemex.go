package phemex

import (
	"fmt"
	"math"
)

const TAKER_FEE_RATE = 0.00075
const MAKER_FEE_RATE = 0.00025

type Account struct {
	Accounts map[string]TradeAccount
}

type TradeAccount struct {
	meta              *accountMeta
	activeOrders      *map[int]ActiveOrder // map[id]ActiveOrder
	conditionalOrders *[]ConditionalOrder
	openPosition      *OpenPosition
	closedPostion     *[]ClosedPosition
	configuration     *Configuration
}

type accountMeta struct {
	marketSymbol      string
	totalBalance      float64
	avialableBalance  float64
	leverage          float64
	riskLimit         float64
	initialMarginRate float64
	maintenanceRate   float64
	maxLeverage       float64
	maxCost           struct {
		long  float64
		short float64
	}
}

type Configuration struct {
	Limit struct {
		OrderCondition []string
		TimeInforce    string
		PostOnly       struct {
			TimeInforce string
		}
		ReduceOnly struct {
			TimeInforce string
		}
		BracketOrder struct {
			OrderType   string
			Trigger     string
			TimeInforce string
		}
	}
	Conditional struct {
		Limit struct {
			TriggerBy      string
			OrderCondition []string
			TimeInforce    string
			PostOnly       struct {
				TimeInforce string
			}
			CloseOnTrigger struct {
				TimeInforce string
			}
		}
		Market struct {
			TriggerBy      string
			OrderCondition string
		}
	}
}

type ActiveOrder struct {
	Side       string
	Symbol     string
	Quantity   int64
	OrderPrice float64
	OrderValue float64
	Filled     int64
	Remaining  int64
	// TakeProfit
	// StopLoss
	FillPrice float64
	Type      string
	Status    string
	Time      string
}

type ConditionalOrder struct {
	Symbol     string
	Quantity   int64
	OrderPrice float64
	// trigger
	// ActivationPrice
	// TriggeringPrice // (Distance)
	Type   string
	Status string
	Time   string
}

type TradeInput struct {
	Symbol      string
	Side        string
	Quantity    int64
	Price       float64
	Leverage    float64
	OrderType   string
	OrderOption string
	OrderLife   string
}

type OpenPosition struct { // TradeAccount, TradePosition
	Symbol            string
	Side              string
	Size              int64
	Value             float64
	AvgEntryPrice     float64
	MarkPrice         float64
	LiquidationPrice  float64
	BankruptPrice     float64
	Margin            float64
	InitialMargin     float64
	MaintenanceMargin float64
	// PositionMargin    float64
	// OrderMargin       float64
	// MarginBalance     float64
	PNL struct {
		UnrealisedPNL float64
		RealisedPNL   float64
	}
	TradeOrders struct {
		OpenTrade  []Position
		CloseTrade []ClosedPosition
	}
	OrderFees struct {
		OpenFees    []float64
		CloseFees   []float64
		FundingFees []float64
	}
	// TakeProfit float64
	// StopLoss   float64
}

type Position struct {
	Symbol            string
	Side              string
	Quantity          int64
	OrderValue        float64
	leverage          float64
	EntryPrice        float64
	OrderMargin       float64
	InitialMargin     float64
	MaintenanceMargin float64
	OpenFee           float64
}

// on program boot up
// 	-> pull data from Database that is relevent to account info and active trades
//	-> pull data from exchange that is relevent to acocunt info and active trades
//	-> compare data from database to exchange
//	-> update with current/correct account data information
//	-> track and monitor relevent market symbol data
//	-> post trades when an event is triggered

// processes of program
//	-> execute trades through SDK to Exchange API
// 	-> set triggers and automate trades, configure events algorithms
//	-> track and monitor exhange data of symbols
//	-> data processing, data anylsis of historic data and monitored data
//	-> store data -> account data, tracked data, events, proccessed data, anylized data, historic data
//	-> front end to interface with backend/daemon
//	->	-> display data in a meaningful way, layout, graph visualized
//	-> 	-> interact -> set trades, triggers, deposit, withdraw, manage
//	-> 	-> manage user account info

func ExchangeFeeAmount(value float64, marketOrder string) float64 {

	rate := 0.0

	if marketOrder == "taker" {
		rate = -0.0075
	} else if marketOrder == "maker" {
		rate = 0.0025
	}

	return truncate((value * rate), math.Pow(10, 8))
}

func FundingFeeAmount(value, rate float64) float64 {
	// pull rate and last traded mark price or pull funded fee amount from api every 8 hours
	return truncate((value * rate), math.Pow(10, 8))
}

type Trade struct {
	*Account
}

func CreateTradeAccount(symbol string) *TradeAccount {

	account := new(TradeAccount)
	account.setUp()
	account.SetDefaultConfiguration()
	account.SetMarketSymbol(symbol)
	account.SetRiskLimit(0.0)
	account.SetLeverage(0.0)
	account.SetBalance(0.0)

	return account
}

func (a *TradeAccount) setUp() {
	a.meta = new(accountMeta)
	a.activeOrders = new(map[int]ActiveOrder)
	a.closedPostion = new([]ClosedPosition)
	a.conditionalOrders = new([]ConditionalOrder)
	a.openPosition = new(OpenPosition)
}

func (a *TradeAccount) SetBalance(balance float64) {
	a.meta.totalBalance = balance
	a.meta.avialableBalance = a.meta.totalBalance
	if balance > 0.0 {
		a.CalcMaxMargin()
	}
}

func (a *TradeAccount) SetDefaultConfiguration() {
	config := new(Configuration)
	a.configuration = config

	config.Limit.OrderCondition = []string{"Post-Only"}
	config.Limit.TimeInforce = "GoodTillCancel"
	config.Limit.PostOnly.TimeInforce = "GoodTillCancel"
	config.Limit.ReduceOnly.TimeInforce = "GoodTillCancel"
	config.Limit.BracketOrder.TimeInforce = "GoodTillCancel"
	config.Limit.BracketOrder.OrderType = "Limit"

	config.Limit.BracketOrder.Trigger = "Last Price"
	config.Conditional.Limit.TriggerBy = "Last Price"
	config.Conditional.Market.TriggerBy = "Last Price"

	config.Conditional.Limit.OrderCondition = []string{"Post-Only"}
	config.Conditional.Limit.TimeInforce = "GoodTillCancel"
	config.Conditional.Limit.PostOnly.TimeInforce = "GoodTillCancel"
	config.Conditional.Limit.CloseOnTrigger.TimeInforce = "GoodTillCancel"

	config.Conditional.Market.OrderCondition = "Close On Trigger"
}

func (a *TradeAccount) SetLeverage(leverage float64) {
	a.meta.leverage = math.Floor(leverage*100) / 100

	if a.meta.leverage > 100.00 {
		a.meta.leverage = 100.00
	} else if a.meta.leverage < 1.00 {
		a.meta.leverage = 1.00
	}
}

func (a *TradeAccount) SetMarketSymbol(symbol string) {
	if symbol != "" {
		a.meta.marketSymbol = symbol
	} else {
		a.meta.marketSymbol = "BTCUSD"
	}
}

func (a *TradeAccount) SetRiskLimit(risk float64) {
	a.meta.riskLimit = risk

	if a.meta.marketSymbol == "BTCUSD" {

		if a.meta.riskLimit <= 100.0 {
			a.meta.initialMarginRate = 0.01
			a.meta.maintenanceRate = 0.005
			a.meta.maxLeverage = 100.00
		} else if a.meta.riskLimit <= 150.0 {
			a.meta.initialMarginRate = 0.015
			a.meta.maintenanceRate = 0.01
			a.meta.maxLeverage = 66.66
		} else if a.meta.riskLimit <= 200.0 {
			a.meta.initialMarginRate = 0.02
			a.meta.maintenanceRate = 0.015
			a.meta.maxLeverage = 50.0
		} else if a.meta.riskLimit <= 250.0 {
			a.meta.initialMarginRate = 0.025
			a.meta.maintenanceRate = 0.02
			a.meta.maxLeverage = 40.0
		} else if a.meta.riskLimit <= 300.0 {
			a.meta.initialMarginRate = 0.03
			a.meta.maintenanceRate = 0.025
			a.meta.maxLeverage = 33.33
		} else if a.meta.riskLimit <= 350.0 {
			a.meta.initialMarginRate = 0.035
			a.meta.maintenanceRate = 0.03
			a.meta.maxLeverage = 28.57
		} else if a.meta.riskLimit <= 400.0 {
			a.meta.initialMarginRate = 0.04
			a.meta.maintenanceRate = 0.035
			a.meta.maxLeverage = 25.0
		} else if a.meta.riskLimit <= 450.0 {
			a.meta.initialMarginRate = 0.045
			a.meta.maintenanceRate = 0.04
			a.meta.maxLeverage = 22.22
		} else if a.meta.riskLimit <= 500.0 {
			a.meta.initialMarginRate = 0.05
			a.meta.maintenanceRate = 0.0045
			a.meta.maxLeverage = 20.0
		} else if a.meta.riskLimit <= 550.0 {
			a.meta.initialMarginRate = 0.055
			a.meta.maintenanceRate = 0.05
			a.meta.maxLeverage = 18.18
		}
	}

	if a.meta.leverage > a.meta.maxLeverage {
		a.meta.leverage = a.meta.maxLeverage
	}
}

func (a *TradeAccount) GetAccount() {

	fmt.Println("Status: \t", a.meta)
	fmt.Println("Active Orders: \t", a.activeOrders)
	fmt.Println("Conditional Orders: \t", a.conditionalOrders)
	fmt.Println("Open Postion: \t", a.openPosition)
	fmt.Println("Closed Positions: \t", a.closedPostion)

	// fmt.Println("Config: \t", a.configuration)
	fmt.Println()
}

func (a *TradeAccount) calcCost(side string, value float64) float64 {
	// [short, long], leverage, quantity, price, -/+takerFee, initialMargin, maintenanceMargin
	// ((qty / price * initMarginRate) + (qty / price * maintRate)) * 0.1) + (qty / price * -/+takerFeeRate / leverage)

	factor := math.Pow(10, 8)

	initialMargin := truncate((value / a.meta.leverage), factor)

	orderMargin := truncate((value * a.meta.initialMarginRate), factor)

	maintenanceMargin := truncate((value * a.meta.maintenanceRate), factor)

	margin := truncate(((orderMargin + maintenanceMargin) * 0.1), factor)

	takerFee := truncate((value * TAKER_FEE_RATE / a.meta.leverage), factor)

	// fmt.Printf("%.8f\t", value)
	// fmt.Printf("%.8f\t", initialMargin)
	// fmt.Printf("%.8f\t", orderMargin)
	// fmt.Printf("%.8f\t", maintenanceMargin)
	// fmt.Printf("%.8f\t", margin)
	// fmt.Printf("%.8f\t", takerFee)
	// fmt.Printf("%.8f\t", truncate(truncate(margin+takerFee, factor)+initialMargin, factor))
	// fmt.Printf("%.8f\t", (margin-takerFee)+initialMargin)
	// fmt.Println()

	if side == "Long" {
		return truncate(truncate(margin+takerFee, factor)+initialMargin, factor)
	} else if side == "Short" {
		return truncate(truncate(margin-takerFee, factor)+initialMargin, factor)
	}
	return 0.0
}

func (a *TradeAccount) CalcMaxMargin() {

	prev := a.meta.leverage

	factor := math.Pow(10, 8)

	findmaxLeverage := func() {

		list := [10]float64{100.0, 66.66, 50.0, 40.0, 33.33, 28.57, 25.0, 22.22, 20.0, 18.18}
		values := [10]float64{100.0, 150.0, 200.0, 250.0, 300.0, 350.0, 400.0, 450.0, 500.0, 550.0}

		tier := 0
		for i := tier; i < len(list)-1; i++ {
			value := truncate((a.meta.totalBalance * list[i]), factor)

			if value <= values[i] {
				tier = i
				break
			}
		}

		target := tier
		spread := 0
		for i := tier; i > 0; i-- {
			value := truncate((a.meta.totalBalance * list[i]), factor)
			if value <= values[target] {
				spread = tier - i
			}
		}

		if spread == 0 && a.meta.leverage > list[tier] {
			a.SetLeverage(list[tier])
			return
		} else if spread == 0 && a.meta.leverage <= list[tier] {
			return
		}

		min := list[tier]
		max := list[tier-1]
		mid := truncate(min+truncate((max-min), 100)/2, 100)
		value := truncate((a.meta.totalBalance * mid), factor)

		for value != values[tier-spread] {

			if value > values[tier-spread] {
				max = mid
			} else if value < values[tier-spread] {
				min = mid
			}

			mid = truncate(min+truncate((max-min), 100)/2, 100)
			value = truncate((a.meta.totalBalance * mid), factor)
		}

		if a.meta.leverage > mid {
			a.SetLeverage(mid)
		}
	}

	find := func(side string) float64 {
		findmaxLeverage()

		value := truncate((a.meta.totalBalance * a.meta.leverage), factor)

		a.SetRiskLimit(value)

		max := a.meta.totalBalance
		base := truncate((a.calcCost(side, value) - a.meta.totalBalance), factor)
		min := truncate((a.meta.totalBalance - base), factor)
		value = truncate((min * a.meta.leverage), factor)
		result := a.calcCost(side, value)

		for result > a.meta.totalBalance {
			base = truncate((base * 2), factor)
			min = truncate((a.meta.totalBalance - base), factor)
			value = truncate((min * a.meta.leverage), factor)
			result = a.calcCost(side, value)
		}

		mid := truncate((min + (base / 2)), factor)
		value = truncate((mid * a.meta.leverage), factor)
		result = a.calcCost(side, value)

		for result != a.meta.totalBalance {

			if base <= 0.00000002 {
				a.SetRiskLimit(0.0)
				return min

			} else if result > a.meta.totalBalance {
				max = mid
				value = truncate((max * a.meta.leverage), factor)
				base = math.Ceil((a.calcCost(side, value)-a.meta.totalBalance)*factor) / factor
				mid = truncate((min + (base / 2)), factor)

			} else if result < a.meta.totalBalance {
				min = mid
				value = truncate((min * a.meta.leverage), factor)
				base = math.Ceil((a.meta.totalBalance-a.calcCost(side, value))*factor) / factor
				mid = truncate((min + (base / 2)), factor)

			}

			value = truncate((mid * a.meta.leverage), factor)
			result = a.calcCost(side, value)

		}

		a.SetRiskLimit(0.0)
		return mid
	}

	a.meta.maxCost.long = find("Long")
	a.meta.maxCost.short = find("Short")
	a.SetRiskLimit(0.0)
	a.SetLeverage(prev)
}

// Create TradeAccount with Default values
// Config TradeAccount with Default Configuration
// removed variadic, converted to private data

// 8, 12, 18
