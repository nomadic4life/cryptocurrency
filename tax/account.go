package tax

import (
	"fmt"
)

type trade struct {
	symbol struct {
		deduct string
		append string
	}
	balance struct {
		quote float64
		base  float64
	}
	PNL           float64
	unrealizedPNL float64
	assetRecords  []int
	queue         struct {
		quote []AssetTrade
		base  []AssetTrade
	}
}

func NewAccount(capital float64) *Account {
	account := Account{}
	account.Statement.TotalCapital = capital
	account.AssetsHoldings = make(map[string]float64)
	account.AssetsHoldings["USD"] = capital
	account.AssetsHoldings["BTC"] = 0.0
	account.AssetsHoldings["ETH"] = 0.0
	account.AssetsHoldings["LBC"] = 0.0
	account.AssetsHoldings["DOGE"] = 0.0
	account.AssetsHoldings["XRP"] = 0.0
	return &account
}

func (a *Account) CreateTransaction(t Trade) {
	transaction := newTransaction(a, t)
	// trade := a.initTrade(transaction)
	// a.outFlow(trade, transaction)
	// a.inflow(trade, transaction)
	// a.updateAccount(trade, transaction)
	fmt.Println(transaction)
}

func (a *Account) initTrade(t Transaction) (ret trade) {
	// ret.balance.quote = 0
	// ret.balance.base = 0
	// ret.PNL = a.Statement.PNL
	// ret.unrealizedPNL = 0
	// ret.queue.base = []
	// ret.queue.quote = []
	return ret
}

func (a *Account) outFlow(o trade, t Transaction) { // -> trade, transaction
	fmt.Println(t)
}

func (a *Account) inflow(o trade, t Transaction) { // -> trade, transaction
	fmt.Println(t)
}

func (a *Account) deduct(t Transaction) { // -> trade
	fmt.Println(t)
}

func (a *Account) append(t Transaction) { // -> queue, trade, transaction
	fmt.Println(t)
}

func (a *Account) updateAccount(o trade, t Transaction) { // -> trade, transaction
	fmt.Println(t)
}

func (a *Account) getID() int64 {
	// might be working, need to add a transaction to history to see if it works
	fmt.Println("ID: ", a.Ledger.TransactionHistory, len(a.Ledger.TransactionHistory))
	return int64(len(a.Ledger.TransactionHistory))
}
