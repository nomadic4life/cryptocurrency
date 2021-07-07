package tax

import "fmt"

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

func (a *Account) CreateTrade(t Trade) {
	fmt.Println(t)
}

func (a *Account) init(t Transaction) { // -> transaction
	fmt.Println(t)
}

func (a *Account) outFlow(t Transaction) { // -> trade, transaction
	fmt.Println(t)
}

func (a *Account) inflow(t Transaction) { // -> trade, transaction
	fmt.Println(t)
}

func (a *Account) deduct(t Transaction) { // -> trade
	fmt.Println(t)
}

func (a *Account) append(t Transaction) { // -> queue, trade, transaction
	fmt.Println(t)
}

func (a *Account) updateAccount(t Transaction) { // -> trade, transaction
	fmt.Println(t)
}

func (a *Account) getID() int64 {
	fmt.Println(len(a.Ledger.TransactionHistory))
	return int64(len(a.Ledger.TransactionHistory))
}
