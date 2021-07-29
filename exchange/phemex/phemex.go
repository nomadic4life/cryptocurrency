package phemex

import (
	"fmt"
	"math"
)

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

func CreateTrade() {}

func (trade *Trade) Long(price float64, quantity int64, order string) {
	// Limit
	// -> Limit Price
	// -> Quanitity
	// -> Order Value
	// -> Available Balance
	// -> Cost
	// Options
	// -> Post-Only
	// 	-> GoodTillCancel
	// -> Reduce-Only
	// 	-> GoodTillCancel
	// 	-> ImmediateOrCancel
	// 	-> FillOrKill
	// -> BracketOrder
	// 	-> Order type: [Limit]
	// 	-> TP [Ticks]
	// 	-> Order type: [Limit]
	// 	-> SL [Ticks]
	// 	-> Trigger: Last Price
	// 	-> GoodTillCancel

	// Market
	// -> Quanitity
	// -> Order Value
	// -> Available Balance
	// -> Cost

	// Conditional
	// -> Market
	//	-> Trigger Price
	//	-> Trigger By [Last Price, Mark Price]
	//	-> Quantity
	// 	-> Cost
	//	-> Trigger
	// -> Limit
	//	-> Trigger Price
	//	-> Trigger By [Last Price, Mark Price]
	//	-> Limit Price
	//	-> Quantity
	// 	-> Cost
	//	-> Trigger
	// Options
	// -> Post-Only
	// 	-> GoodTillCancel
	// -> Close on Trigger
	// 	-> GoodTillCancel
	// 	-> ImmediateOrCancel
	// 	-> FillOrKill

}

// func (trade *Trade) Short(price float64, size int64) {

// }

// func calcCost() {
// 	// [short, long], leverage, quantity, price, -/+takerFee, initialMargin, maintenanceMargin

// 	// ((quantity / price) * takerFeeRate / leverage) + (((quantity / price) * initialMarginRate) + ((quantity / price) * maintenanceRate)) * 0.1
// }

// func (a *TradeAccount) LimitTrade(price float64, quantity int64, options []string) {}

// func (a *TradeAccount) MarketTrade(price float64, quantity int64, side string) {}

// func (a *TradeAccount) Trade(price float64, quantity int64, order []string) {
// 	const ORDER = 0
// 	const SIDE = 1
// 	feeRate := 0.0

// 	if order[ORDER] == "LIMIT" {
// 		feeRate = -0.00025
// 	} else if order[SIDE] == "LONG" && order[ORDER] == "MARKET" {
// 		feeRate = +0.00075
// 	} else if order[SIDE] == "SHORT" && order[ORDER] == "MARKET" {
// 		feeRate = -0.00075
// 	}
// 	// order [["Long", "Short"], ["Limit", "Market"]]
// 	// order[1]["Market"] -> FeeRate = +0.00075
// 	// order[1]["Limit"] -> FeeRate = -0.00025

// }

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

		if a.meta.riskLimit <= 100 {
			a.meta.initialMarginRate = 0.01
			a.meta.maintenanceRate = 0.005
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 150 {
			a.meta.initialMarginRate = 0.015
			a.meta.maintenanceRate = 0.01
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 200 {
			a.meta.initialMarginRate = 0.02
			a.meta.maintenanceRate = 0.015
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 250 {
			a.meta.initialMarginRate = 0.025
			a.meta.maintenanceRate = 0.02
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 300 {
			a.meta.initialMarginRate = 0.03
			a.meta.maintenanceRate = 0.025
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 350 {
			a.meta.initialMarginRate = 0.035
			a.meta.maintenanceRate = 0.03
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 400 {
			a.meta.initialMarginRate = 0.04
			a.meta.maintenanceRate = 0.035
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 450 {
			a.meta.initialMarginRate = 0.045
			a.meta.maintenanceRate = 0.04
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 500 {
			a.meta.initialMarginRate = 0.05
			a.meta.maintenanceRate = 0.0045
			// a.meta.MaxLeverage = 100
		} else if a.meta.riskLimit <= 550 {
			a.meta.initialMarginRate = 0.055
			a.meta.maintenanceRate = 0.05
			// a.meta.MaxLeverage = 100
		}
	}
}

func (a *TradeAccount) GetAccount() {
	fmt.Println("Status: \t", a.meta)
	fmt.Println("Active Orders: \t", a.activeOrders)
	fmt.Println("Conditional Orders: \t", a.conditionalOrders)
	fmt.Println("Open Postion: \t", a.openPosition)
	fmt.Println("Closed Positions: \t", a.closedPostion)
	// fmt.Println("Config: \t", a.configuration)
}

// Create TradeAccount with Default values
// Config TradeAccount with Default Configuration
