package tax

import (
	"fmt"
	"strconv"
	"strings"
)

type tradeLog struct {
	symbol struct {
		deduct string
		append string
	}
	balance struct {
		deduct float64
		quote  float64
		base   float64
	}
	statement struct {
		PNL           float64
		unrealizedPNL float64
	}
	ledger struct {
		transaction TransactionEntry
		costBases   []CostBasisEntry
	}
	queue struct {
		quote []AssetCostBasis
		base  []AssetCostBasis
	}
}

func NewAccount(capital float64) *Account {
	account := new(Account)
	account.Statement.TotalCapital = capital

	account.Ledger.Transactions = make([]TransactionEntry, 0, 20)
	account.Ledger.CostBases = make([]CostBasisEntry, 0, 20)

	account.Assets.AssetCostBasesQueue = make(map[string][]AssetCostBasis)

	account.Assets.Holdings = make(map[string]float64)
	account.Assets.Holdings["USD"] = capital
	account.Assets.Holdings["BTC"] = 0.0
	account.Assets.Holdings["ETH"] = 0.0
	account.Assets.Holdings["LBC"] = 0.0
	account.Assets.Holdings["DOGE"] = 0.0
	account.Assets.Holdings["XRP"] = 0.0
	return account
}

func (a *Account) CreateTransaction(input TradeInput) {
	// Concepts:
	//    debit     /   credit
	//    outflow   /   inflow
	//    deduct    /   append

	log := a.initLog(input)

	a.outFlow(log)
	a.inflow(log)
	a.updateAccount(log)
	a.log()
}

func (a *Account) initLog(input TradeInput) *tradeLog {
	transaction := newTransaction(a, input)
	log := new(tradeLog)

	log.ledger.transaction = *transaction
	log.balance.quote = a.getAssetHoldings(transaction.quote())
	log.balance.base = a.getAssetHoldings(transaction.base())
	log.statement.PNL = a.Statement.PNL
	// log.statement.unrealizedPNL = a.Statement.unlrealizedPNL

	if _, ok := a.Assets.AssetCostBasesQueue[transaction.base()]; ok {
		log.queue.base = a.Assets.AssetCostBasesQueue[transaction.base()][:]
	} else {
		log.queue.base = make([]AssetCostBasis, 0, 10)
	}

	if _, ok := a.Assets.AssetCostBasesQueue[transaction.quote()]; ok {
		log.queue.quote = a.Assets.AssetCostBasesQueue[transaction.quote()][:]
	} else {
		log.queue.quote = make([]AssetCostBasis, 0, 10)
	}

	return log
}

func (a *Account) outFlow(log *tradeLog) {

	if log.ledger.transaction.OrderType == "BUY" {

		log.symbol.deduct = log.ledger.transaction.quote()
		log.balance.quote -= log.ledger.transaction.OrderAmount

		if log.ledger.transaction.quote() != "USD" {
			log.balance.deduct = log.ledger.transaction.OrderAmount
			deductions, record := a.deduct(log)
			copy(log.queue.quote, *record)

			for i := 0; i < len(deductions); i++ {
				log.ledger.costBases = append(log.ledger.costBases, *newCostBasisEntry(&deductions[i], log))
			}
		}
	}

	if log.ledger.transaction.OrderType == "SELL" {

		log.symbol.deduct = log.ledger.transaction.base()
		log.balance.base -= log.ledger.transaction.OrderQuantity
		log.balance.deduct = log.ledger.transaction.OrderQuantity
		deductions, record := a.deduct(log)
		copy(log.queue.base, *record)

		for i := 0; i < len(deductions); i++ {
			log.ledger.costBases = append(log.ledger.costBases, *newCostBasisEntry(&deductions[i], log))
		}
	}
}

func (a *Account) inflow(log *tradeLog) {

	if log.ledger.transaction.OrderType == "BUY" {
		log.balance.base += log.ledger.transaction.OrderQuantity
		log.symbol.append = log.ledger.transaction.base()
		log.ledger.costBases = append(log.ledger.costBases, a.append(&log.queue.base, log))

	}

	if log.ledger.transaction.OrderType == "SELL" {
		log.balance.quote += log.ledger.transaction.OrderAmount

		if log.ledger.transaction.quote() != "USD" {
			log.symbol.append = log.ledger.transaction.quote()
			log.ledger.costBases = append(log.ledger.costBases, a.append(&log.queue.quote, log))

		}
	}
}

func (a *Account) deduct(log *tradeLog) (assetCostBasisList, *assetCostBasisList) {

	deductions := make(assetCostBasisList, 0, 40)
	records := append(assetCostBasisList(nil), a.Assets.AssetCostBasesQueue[log.symbol.deduct]...)
	termination := 0.0

	for log.balance.deduct > termination {
		deductions = append(deductions, records[0])

		if records[0].BaseAmount < log.balance.deduct {
			log.balance.deduct -= records[0].BaseAmount
			deductions[len(deductions)-1].ChangeAmount = records[0].BaseAmount
			records.dequeue()
			// records = records[1:]

		} else {
			records[0].BaseAmount -= log.balance.deduct
			deductions[len(deductions)-1].ChangeAmount = log.balance.deduct
			log.balance.deduct = 0.0

		}

		if records[0].BaseAmount == 0 {
			records.dequeue()
			// records = records[1:]
		}

	}
	return deductions, &records
}

func (a *Account) append(queue *[]AssetCostBasis, log *tradeLog) CostBasisEntry {

	entry := newCostBasisEntry(nil, log)

	asset := AssetCostBasis{
		TransactionID: entry.TransactionID,
		QuotePrice:    entry.QuotePriceEntry,
		USDPriceValue: entry.USDPriceEntry,
		BaseAmount:    entry.BalanceRemaining.BaseAmount[1]}

	*queue = append(*queue, asset)

	return *entry
}

func (a *Account) updateAccount(log *tradeLog) {

	// if _, ok := a.AssetCostBasesQueue[t.base()]; ok == false {
	// 	a.AssetCostBasesQueue[t.base()] = make([]AssetCostBasis, 0, 10)
	// }

	// if _, ok := a.AssetCostBasesQueue[t.quote()]; ok == false && t.quote() != "USD" {
	// 	a.AssetCostBasesQueue[t.quote()] = make([]AssetCostBasis, 0, 10)
	// }

	a.Statement.PNL += log.statement.PNL

	// a.Ledger.Transactions = append(a.Ledger.Transactions, *&log.ledger.transaction)
	// a.Ledger.CostBases = append(a.Ledger.CostBases, log.ledger.costBases...)
	a.Ledger.Transactions.enqueue(&log.ledger.transaction)
	a.Ledger.CostBases.enqueue(&log.ledger.costBases)

	a.Assets.Holdings[log.ledger.transaction.base()] = log.balance.base
	a.Assets.Holdings[log.ledger.transaction.quote()] = log.balance.quote

	a.Assets.AssetCostBasesQueue[log.ledger.transaction.base()] = log.queue.base

	if log.ledger.transaction.quote() != "USD" {
		a.Assets.AssetCostBasesQueue[log.ledger.transaction.quote()] = log.queue.quote
	}
}

func (a *Account) getID() int64 {
	return int64(len(a.Ledger.Transactions))
}

func (a *Account) getAssetHoldings(symbol string) float64 {
	if val, ok := a.Assets.Holdings[symbol]; ok {
		return val
	}

	return 0.0
}

func (a *Account) log() {

	fmt.Print("\n")
	fmt.Println("STATEMENT:", a.Statement)
	fmt.Println("ASSETHOLDINGS:", a.Assets.Holdings)
	fmt.Println("LEDGER -> TRANSACTION:")

	for i := 0; i < len(a.Ledger.Transactions); i++ {

		fmt.Print("\t", "ID: ", a.Ledger.Transactions[i].TransactionID)
		fmt.Print("\t", "DATE: ", a.Ledger.Transactions[i].Date)
		fmt.Print("\t\t", "PAIR: ", a.Ledger.Transactions[i].OrderPair)
		fmt.Print("\t\t", "TYPE: ", a.Ledger.Transactions[i].OrderType)
		fmt.Print("\t", "PRICE: ", a.Ledger.Transactions[i].OrderPrice)
		fmt.Print("\t", "QUANTITY: ", a.Ledger.Transactions[i].OrderQuantity)
		fmt.Print("\t", "AMOUNT: ", a.Ledger.Transactions[i].OrderAmount)
		fmt.Print("\t", "VALUE: ", a.Ledger.Transactions[i].USDPriceValue)
		fmt.Print("\n")

	}

	fmt.Print("\n")
	fmt.Println("LEDGER -> COST BASIS:")
	for i := 0; i < len(a.Ledger.CostBases); i++ {

		fmt.Print("\t", "ID ->: ", a.Ledger.CostBases[i].TransactionID.From)
		fmt.Print("", " - ", a.Ledger.CostBases[i].TransactionID.To)
		fmt.Print("\t", "QPEntry: ", a.Ledger.CostBases[i].QuotePriceEntry)
		fmt.Print("\t", "QPExit: ", a.Ledger.CostBases[i].QuotePriceExit)
		fmt.Print("\t", "USDPEntry: ", a.Ledger.CostBases[i].USDPriceEntry)
		fmt.Print("\t", "USDPExit: ", a.Ledger.CostBases[i].USDPriceExit)
		fmt.Print("\t", "C -> Q: ", a.Ledger.CostBases[i].ChangeAmount.BaseQuantity)
		fmt.Print("\t", "C -> A: ", a.Ledger.CostBases[i].ChangeAmount.QuoteAmount)
		fmt.Print("\t", "C -> V: ", a.Ledger.CostBases[i].ChangeAmount.USDValue)
		fmt.Print("\t", "R -> BA: ", a.Ledger.CostBases[i].BalanceRemaining.BaseAmount)
		fmt.Print("\t", "R -> BV: ", a.Ledger.CostBases[i].BalanceRemaining.BaseValue)
		fmt.Print("\t", "R -> V: ", a.Ledger.CostBases[i].BalanceRemaining.USDValue)
		fmt.Print("\n")
	}

	fmt.Println("COST BASIS ASSET QUEUE:", a.Assets.AssetCostBasesQueue)
	fmt.Print("\n")
}

func (l *tradeLog) log() {
	fmt.Print("\n")
	fmt.Println("\t symbol: ->  deduct:", l.symbol.deduct)
	fmt.Println("\t symbol: ->  append:", l.symbol.append)
	fmt.Println("\t balance: ->  quote:", l.balance.quote)
	fmt.Println("\t balance: ->  base:", l.balance.base)
	fmt.Println("\t amount deducted:", l.balance.deduct)
	fmt.Println("\t PNL:", l.statement.PNL)
	fmt.Println("\t unrealizedPNL:", l.statement.unrealizedPNL)
	fmt.Println("\t assetRecords:", l.ledger.costBases)
	fmt.Println("\t queue: -> quote", l.queue.quote)
	fmt.Println("\t queue: -> base", l.queue.base)
}

func dollarFormat(value float64) string {
	dollar := fmt.Sprintf("%.2f", value)
	split := strings.Split(dollar, ".")
	return "$" + commaSep(split)
}

func cryptoFormat(value float64) string {
	a, err := strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	if err != nil {
		fmt.Println("error")
	}

	b, err := strconv.ParseFloat(fmt.Sprintf("%.8f", value), 64)
	if err != nil {
		fmt.Println("error")
	}

	num := strings.Split(fmt.Sprintf("%f", value), ".")[0]

	var crypto string

	if len(num) <= 5 && a != b {
		crypto = fmt.Sprintf("%.8f", value)
	} else if len(num) <= 7 || a == b {
		crypto = fmt.Sprintf("%.4f", value)
	} else {
		crypto = fmt.Sprintf("%.1f", value)
		fmt.Println(crypto)
		split := strings.Split(crypto, ".")
		split[1] = "0"
		crypto = strings.Join(split, ".")
	}

	return commaSep(strings.Split(crypto, "."))
}

func commaSep(value []string) string {
	offset := len(value[0]) / 3
	size := len(value[0]) + offset
	results := make([]byte, size, size)
	counter := 0

	// adding comma seperation
	for i := len(value[0]) - 1; i >= 0; i-- {
		if (len(value[0])-1-i)%3 == 0 && len(value[0])-1-i != 0 {
			results[i+offset-counter] = ','
			counter++
		}

		results[i-counter+offset] = value[0][i]
	}
	value[0] = string(results)
	return strings.Join(value, ".")
}

// width of a field
//  - 4 min for USD
//  - $0.00         ->  6
//  - $00,000.00    ->  10
//  - $000,000      ->  8
//  - 6 min for crypto
//  - crypto
//  - 0.0000            ->  size 6
//  - 0.0000_0000       ->  size 10
//  - 0,000.0000        ->  size 10
//  - 0,000.0000_0000   ->  size 13
//  - 000,000.0000      ->  size 12
//  - 000,000,000.0000  ->  size 16
//  - 0,000,000         ->  size 9

// padding
