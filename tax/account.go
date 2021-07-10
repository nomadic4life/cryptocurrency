package tax

import (
	"fmt"
)

type tradeLog struct {
	// trade data
	// trade record
	// transcription
	// dictation
	// log
	// registry
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
	assetRecords   []CostBasisEntry
	queue          struct {
		quote []AssetTrade
		base  []AssetTrade
	}
}

func NewAccount(capital float64) *Account {
	account := Account{}
	account.Statement.TotalCapital = capital
	account.AssetsHoldings = make(map[string]float64)
	// account.Ledger.TransactionHistory = make([]Transaction, 0, 20)
	account.CostBasisAssetQueue = make(map[string][]AssetTrade)
	account.AssetsHoldings["USD"] = capital
	account.AssetsHoldings["BTC"] = 0.0
	account.AssetsHoldings["ETH"] = 0.0
	account.AssetsHoldings["LBC"] = 0.0
	account.AssetsHoldings["DOGE"] = 0.0
	account.AssetsHoldings["XRP"] = 0.0
	return &account
}

func (a *Account) CreateTransaction(input TradeInput) {
	// Concepts:
	//    debit     /   credit
	//    outflow   /   inflow
	//    deduct    /   append
	transaction := newTransaction(a, input)
	log := a.initLog(transaction)

	a.outFlow(log, transaction)
	//	-> deduct()
	//	-> newCostBasisRecord()
	a.inflow(log, transaction)
	//	-> append()
	//	-> newCostBasisRecord()
	a.updateAccount(log, transaction)
	// fmt.Println("create transaction: ", transaction, trade)
	a.log()
}

func (a *Account) initLog(t *Transaction) *tradeLog {
	log := tradeLog{}
	log.balance.quote = a.getAssetHoldings(t.quote())
	log.balance.base = a.getAssetHoldings(t.base())
	log.PNL = a.Statement.PNL
	// log.unrealizedPNL = 0

	if _, ok := a.CostBasisAssetQueue[t.base()]; ok {
		// copy(log.queue.base, a.CostBasisAssetQueue[t.base()])
		log.queue.base = a.CostBasisAssetQueue[t.base()][:]
	} else {
		log.queue.base = make([]AssetTrade, 0, 10)
	}

	if _, ok := a.CostBasisAssetQueue[t.quote()]; ok {
		// copy(log.queue.quote, a.CostBasisAssetQueue[t.quote()])
		log.queue.quote = a.CostBasisAssetQueue[t.quote()][:]
	} else {
		log.queue.quote = make([]AssetTrade, 0, 10)
	}
	return &log
}

func (a *Account) outFlow(log *tradeLog, t *Transaction) {
	fmt.Println("	\n::OUTFLOW::	")

	if t.OrderType == "BUY" {
		fmt.Println("\t::BUY::	")

		log.symbol.deduct = t.quote()
		log.balance.quote -= t.OrderAmount
		// a.deduct(o)

		if t.quote() != "USD" {

			// trade.amountDeducted = transaction[AMOUNT];
			log.amountDeducted = t.OrderAmount

			// const [deductions, record] = this.deduct(trade);
			// deductions, record := a.deduct(o)

			// trade.quoteQueue = [...record];
			// copy(log.queue.quote, record)

			// trade.records.push(...deductions.map(item => new CostBasisRecord(account, item, transaction, trade)));
			// for i := 0; i < len(deductions); i++ {
			// 		log.records.append(newCostBasisRecord(deduction[i], t, o))
			// }

		}
	}

	if t.OrderType == "SELL" {
		fmt.Println("\t::SELL::	")
		a.log()
		// trade.deductSymbol = transaction.base;
		log.symbol.deduct = t.base()

		// trade.baseBalance -= transaction[QUANTITY];
		log.balance.base -= t.OrderQuantity

		// trade.amountDeducted = transaction[QUANTITY];
		// log.balance.deduct = t.OrderQuantity
		log.amountDeducted = t.OrderQuantity

		// const [deductions, record] = this.deduct(trade);
		deductions, record := a.deduct(log)
		// test := *deductions
		// fmt.Println("deduct: ", test[0], "record", record)

		// trade.baseQueue = [...record];
		// log.queue.base = *record
		copy(log.queue.base, *record)

		// trade.records.push(...deductions.map(item => new CostBasisRecord(account, item, transaction, trade)));
		for i := 0; i < len(deductions); i++ {
			log.assetRecords = append(log.assetRecords, *newCostBasisEntry(&deductions[i], log, t))
		}
	}
}

func (a *Account) inflow(log *tradeLog, t *Transaction) {
	fmt.Println("   \n:INFLOW:  ")

	if t.OrderType == "BUY" {
		log.balance.base += t.OrderQuantity
		log.symbol.append = t.base()

		// log.symbol.base = t.base()
		log.assetRecords = append(log.assetRecords, a.append(&log.queue.base, log, t))
	}

	if t.OrderType == "SELL" {
		log.balance.quote += t.OrderAmount

		if t.quote() != "USD" {
			// hasn't been tested could be bugged
			log.symbol.append = t.quote()
			// trade.quoteSymbol = t.quote();
			// trade.records.push(this.append(trade.quoteQueue, transaction, trade));
		}
	}
}

func (a *Account) deduct(log *tradeLog) ([]AssetTrade, *[]AssetTrade) { // -> trade

	deductions := make([]AssetTrade, 0, 40)
	records := append([]AssetTrade(nil), a.CostBasisAssetQueue[log.symbol.deduct]...)

	for log.amountDeducted > 0 {
		deductions = append(deductions, records[0])

		if records[0].BaseAmount < log.amountDeducted {
			log.amountDeducted -= records[0].BaseAmount
			deductions[len(deductions)-1].ChangeAmount = records[0].BaseAmount
			records = records[1:]
		} else {
			records[0].BaseAmount -= log.amountDeducted
			deductions[len(deductions)-1].ChangeAmount = log.amountDeducted
			log.amountDeducted = 0.0
		}

		if records[0].BaseAmount == 0 {
			records = records[1:]
		}

	}
	return deductions, &records
}

func (a *Account) append(queue *[]AssetTrade, log *tradeLog, transaction *Transaction) CostBasisEntry {

	entry := newCostBasisEntry(nil, log, transaction)

	asset := AssetTrade{
		TransactionID: entry.TransactionID,
		QuotePrice:    entry.QuotePriceEntry,
		USDPriceValue: entry.USDPriceEntry,
		BaseAmount:    entry.BalanceRemaining.BaseAmount[1]}

	*queue = append(*queue, asset)

	return *entry
}

func (a *Account) updateAccount(log *tradeLog, transaction *Transaction) {
	fmt.Println("update account: Log", log)
	fmt.Println("update account: Transaction ", transaction)

	// if _, ok := a.CostBasisAssetQueue[t.base()]; ok == false {
	// 	a.CostBasisAssetQueue[t.base()] = make([]AssetTrade, 0, 10)
	// }
	// if _, ok := a.CostBasisAssetQueue[t.quote()]; ok == false && t.quote() != "USD" {
	// 	a.CostBasisAssetQueue[t.quote()] = make([]AssetTrade, 0, 10)
	// }

	a.Ledger.TransactionHistory = append(a.Ledger.TransactionHistory, *transaction)
	a.Ledger.CostBasesHistory = append(a.Ledger.CostBasesHistory, log.assetRecords...)

	a.AssetsHoldings[transaction.base()] = log.balance.base
	a.AssetsHoldings[transaction.quote()] = log.balance.quote

	a.CostBasisAssetQueue[transaction.base()] = log.queue.base
	if transaction.quote() != "USD" {
		a.CostBasisAssetQueue[transaction.quote()] = log.queue.quote
	}
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

func (a *Account) log() {
	fmt.Print("\n")
	fmt.Println("STATEMENT:", a.Statement)
	fmt.Println("ASSETHOLDINGS:", a.AssetsHoldings)
	fmt.Println("LEDGER -> TRANSACTION:", a.Ledger.TransactionHistory)
	fmt.Println("LEDGER -> COST BASIS:", a.Ledger.CostBasesHistory)
	fmt.Println("COST BASIS ASSET QUEUE:", a.CostBasisAssetQueue)
	fmt.Print("\n")
}

func (l *tradeLog) log() {
	fmt.Println("\t symbol: ->  deduct:", l.symbol.deduct)
	fmt.Println("\t symbol: ->  append:", l.symbol.append)
	fmt.Println("\t balance: ->  quote:", l.balance.quote)
	fmt.Println("\t balance: ->  base:", l.balance.base)
	fmt.Println("\t amount deducted:", l.amountDeducted)
	fmt.Println("\t PNL:", l.PNL)
	fmt.Println("\t unrealizedPNL:", l.unrealizedPNL)
	fmt.Println("\t assetRecords:", l.assetRecords)
	fmt.Println("\t queue: -> quote", l.queue.quote)
	fmt.Println("\t queue: -> base", l.queue.base)
}
