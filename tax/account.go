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
		deduct float64
		quote  float64
		base   float64
	}
	amountDeducted float64
	PNL            float64
	unrealizedPNL  float64
	assetRecords   []int
	queue          struct {
		quote []AssetTrade
		base  []AssetTrade
	}
}

func NewAccount(capital float64) *Account {
	account := Account{}
	account.Statement.TotalCapital = capital
	account.AssetsHoldings = make(map[string]float64)
	account.Ledger.TransactionHistory = make([]Transaction, 0, 20)
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
	trade := a.initTrade(transaction)
	a.outFlow(trade, transaction)
	//	-> deduct()
	//	-> newCostBasisRecord()
	// a.inflow(trade, transaction)
	//	-> append()
	//	-> newCostBasisRecord()
	a.updateAccount(trade, transaction)
	fmt.Println("create transaction: ", transaction, trade)
}

func (a *Account) initTrade(t *Transaction) *trade {
	ret := trade{}
	ret.balance.quote = a.getAssetHoldings(t.quote())
	ret.balance.base = a.getAssetHoldings(t.base())
	ret.PNL = a.Statement.PNL
	// ret.unrealizedPNL = 0
	copy(ret.queue.base, a.CostBasisAssetQueue[t.base()])
	copy(ret.queue.quote, a.CostBasisAssetQueue[t.base()])
	return &ret
}

func (a *Account) outFlow(o *trade, t *Transaction) { // -> trade, transaction
	fmt.Println("	\n::OUTFLOW::	")

	if t.OrderType == "BUY" {

		o.symbol.deduct = t.quote()
		o.balance.quote -= t.OrderAmount

		if t.quote() != "USD" {
			// trade.amountDeducted = transaction[AMOUNT];
			o.amountDeducted = t.OrderAmount

			// const [deductions, record] = this.deduct(trade);
			// deductions, record := a.deduct(o)

			// trade.quoteQueue = [...record];
			// copy(o.queue.quote, record)

			// trade.records.push(...deductions.map(item => new CostBasisRecord(account, item, transaction, trade)));
			// for i := 0; i < len(deductions); i++ {
			// 		o.records.append(newCostBasisRecord(deduction[i], t, o))
			// }

		}
	}

	if t.OrderType == "SELL" {
		// trade.deductSymbol = transaction.base;
		o.symbol.deduct = t.base()

		// trade.baseBalance -= transaction[QUANTITY];
		o.balance.base -= t.OrderAmount

		// trade.amountDeducted = transaction[QUANTITY];
		o.amountDeducted = t.OrderQuantity

		// const [deductions, record] = this.deduct(trade);
		// deductions, record := a.deduct(o)

		// trade.baseQueue = [...record];
		// copy(o.queue.base, record)

		// trade.records.push(...deductions.map(item => new CostBasisRecord(account, item, transaction, trade)));
		// for i := 0; i < len(deductions); i++ {
		// 		o.records.append(newCostBasisRecord(deduction[i], t, o))
		// }

	}

}

func (a *Account) inflow(o trade, t Transaction) { // -> trade, transaction
	fmt.Println(t)
}

func (a *Account) deduct(t Transaction) { // -> trade
	fmt.Println(t)
}

func (a *Account) append(o trade, t Transaction) { // -> queue, trade, transaction
	fmt.Println(t)

}

func (a *Account) updateAccount(o *trade, t *Transaction) { // -> trade, transaction
	if _, ok := a.CostBasisAssetQueue[t.base()]; ok {
		a.CostBasisAssetQueue[t.base()] = make([]AssetTrade, 0, 10)
	}
	if _, ok := a.CostBasisAssetQueue[t.quote()]; ok {
		a.CostBasisAssetQueue[t.quote()] = make([]AssetTrade, 0, 10)
	}

	a.Ledger.TransactionHistory = append(a.Ledger.TransactionHistory, *t)
	a.AssetsHoldings[t.base()] = o.balance.base
	a.AssetsHoldings[t.quote()] = o.balance.quote
}

func (a *Account) getID() int64 {
	return int64(len(a.Ledger.TransactionHistory))
}

func (a *Account) getAssetHoldings(symbol string) float64 {
	if val, ok := a.AssetsHoldings[symbol]; ok {
		return val
	}

	return 0.0
}
