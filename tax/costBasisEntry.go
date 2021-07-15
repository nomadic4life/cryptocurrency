package tax

import (
	"fmt"
	"math"
)

func newCostBasisEntry(asset *AssetCostBasis, log *tradeLog) *CostBasisEntry {
	// middleware implementation
	// create ID
	// configue Excuted Price
	// update Change Amount
	// update Holdings
	// update PNL
	// display results?
	// build("Hello World!", func(message string) { fmt.Println(message) })
	cb := []builder{
		createID,
		executedPrice,
		updateAllocation,
		updateBalanceRemaining,
		updateHoldings,
		updatePNL}
	return build(asset, log, &log.ledger.transaction, cb)
}

type exchange struct {
	isBuy  func() currency
	isSell func() currency
}

type currency struct {
	isUSD    func() action
	isCrypto func() action
}

type action struct {
	isAppend func() bool
	isDeduct func() bool
}

type builder func(index int, entry *CostBasisEntry, asset *AssetCostBasis, log *tradeLog, transaction *TransactionEntry)

func checkCondition(orderType, orderPair string, asset *AssetCostBasis) exchange {
	result := false

	isAppend := func() bool {
		if result && asset == nil {
			result = true
		} else {
			result = false
		}

		return result
	}

	isDeduct := func() bool {
		if result && asset != nil {
			result = true
		} else {
			result = false
		}

		return result
	}

	isUSD := func() action {
		if result && orderPair == "USD" {
			result = true
		} else {
			result = false
		}

		return action{
			isAppend,
			isDeduct}
	}

	isCrypto := func() action {
		if result && orderPair != "USD" {
			result = true
		} else {
			result = false
		}

		return action{
			isAppend,
			isDeduct}
	}

	isBuy := func() currency {
		if orderType == "BUY" {
			result = true
		} else {
			result = false
		}

		return currency{
			isUSD,
			isCrypto}
	}

	isSell := func() currency {
		if orderType == "SELL" {
			result = true
		} else {
			result = false
		}

		return currency{
			isUSD,
			isCrypto}
	}

	return exchange{
		isBuy,
		isSell}
}

func build(asset *AssetCostBasis, log *tradeLog, transaction *TransactionEntry, cb []builder) *CostBasisEntry {
	entry := CostBasisEntry{}
	check := checkCondition(transaction.OrderType, transaction.quote(), asset)
	index := 0

	if check.isBuy().isUSD().isDeduct() {
		// OUTFLOW -> BUY -> C/USD -> [NO COST BASIS ENTRY]
		fmt.Println("Throw an Error if reached this condition")
	} else if check.isBuy().isUSD().isAppend() {
		// INFLOW -> BUY -> C/USD -> APPEND -> [COST BASIS ENTRY]
		index = 1
	} else if check.isSell().isUSD().isDeduct() {
		// OUTFLOW -> SELL -> C/USD -> DEDUCTION -> [COST BASIS ENTRY]
		index = 2
	} else if check.isSell().isUSD().isAppend() {
		// INFLOW -> SELL -> C/USD -> APPEND -> [NO COST BASIS ENTRY]
		fmt.Println("Throw an Error if reached this condition")
	} else if check.isBuy().isCrypto().isDeduct() {
		// OUTFLOW -> BUY -> C/C -> DEDUCTION -> [COST BASIS ENTRY]
		index = 4
	} else if check.isBuy().isCrypto().isAppend() {
		// INFLOW -> BUY -> C/C -> APPEND -> [COST BASIS ENTRY]
		index = 5
	} else if check.isSell().isCrypto().isDeduct() {
		// OUTFLOW -> SELL -> C/C DEDUCTION -> [COST BASIS ENTRY]
		index = 6
	} else if check.isSell().isCrypto().isAppend() {
		// INFLOW -> SELL -> C/C APPEND -> [COST BASIS ENTRY]
		index = 7
	}

	entry.meta.orderPair = transaction.OrderPair
	entry.meta.date = int(transaction.Date)

	for _, next := range cb {
		next(index, &entry, asset, log, transaction)
	}

	return &entry
}

func createID(index int, entry *CostBasisEntry, asset *AssetCostBasis, log *tradeLog, transaction *TransactionEntry) {
	// from -> debited
	// to -> credited
	if asset != nil {
		entry.TransactionID.Credit = asset.Debit
	} else {
		entry.TransactionID.Credit = "-"
	}

	entry.TransactionID.Debit = transaction.TransactionID
}

func executedPrice(index int, entry *CostBasisEntry, asset *AssetCostBasis, log *tradeLog, transaction *TransactionEntry) {

	// return 0 -> [1, 5, 7], [], [1, 5, 7], []
	// return transaction.OrderPrice -> [1, 5], [2, 6], [], []
	// return transaction.USDPriceValue -> [7], [4], [1, 5, 7], [2, 4, 6]
	// return asset.QuotePrice -> [2, 6], [], [], []
	// return asset.USDPriceValue -> [4], [], [2, 4, 6], []

	price := [5]float64{
		0.0,
		transaction.OrderPrice,
		transaction.USDPriceValue,
		0.0,
		0.0}

	if asset != nil {
		price[3] = asset.QuotePrice
		price[4] = asset.USDPriceValue
	}

	var table map[string]map[int]int
	table = make(map[string]map[int]int)
	table["Quote Price Entry"] = map[int]int{1: 1, 5: 1, 7: 2, 2: 3, 6: 3, 4: 4}
	table["Quote Price Exit"] = map[int]int{1: 0, 5: 0, 7: 0, 2: 1, 6: 1, 4: 2}
	table["USD Price Entry"] = map[int]int{1: 2, 5: 2, 7: 2, 2: 4, 4: 4, 6: 4}
	table["USD Price Exit"] = map[int]int{1: 0, 5: 0, 7: 0, 2: 2, 4: 2, 6: 2}

	entry.QuotePriceEntry = price[table["Quote Price Entry"][index]]
	entry.QuotePriceExit = price[table["Quote Price Exit"][index]]
	entry.USDPriceEntry = price[table["USD Price Entry"][index]]
	entry.USDPriceExit = price[table["USD Price Exit"][index]]
}

func updateAllocation(index int, entry *CostBasisEntry, asset *AssetCostBasis, log *tradeLog, transaction *TransactionEntry) {

	value := [9]float64{
		transaction.OrderAmount,
		transaction.OrderAmount * entry.lastUSDPrice(),
		transaction.OrderQuantity,
		transaction.OrderQuantity * entry.lastQuotePrice(),
		transaction.OrderQuantity * entry.lastUSDPrice(),
		0.0,
		0.0,
		0.0,
		transaction.OrderQuantity * entry.lastQuotePrice() * entry.lastUSDPrice()}

	if asset != nil {
		value[5] = asset.Quantity
		value[6] = asset.Quantity * entry.lastQuotePrice()
		value[7] = asset.Quantity * entry.lastUSDPrice()
		value[8] = transaction.OrderQuantity * entry.lastQuotePrice() * entry.lastUSDPrice()
		// value[8] = entry.Allocate.Amount * entry.lastUSDPrice()
	}

	var table map[string]map[int]int
	table = make(map[string]map[int]int)
	table["Base Quantity"] = map[int]int{7: 0, 1: 2, 5: 2, 2: 5, 4: 5, 6: 5}
	table["Quote Amount"] = map[int]int{7: 2, 1: 3, 5: 3, 2: 6, 4: 6, 6: 6}
	table["USDValue"] = map[int]int{7: 1, 1: 4, 2: 7, 4: 7, 5: 8, 6: 8}

	entry.Allocation.Quantity = math.Floor(value[table["Base Quantity"][index]]*math.Pow(10, 8)) / math.Pow(10, 8)
	entry.Allocation.Amount = value[table["Quote Amount"][index]]
	entry.Allocation.Value = math.Floor(value[table["USDValue"][index]]*100) / 100
}

func updateBalanceRemaining(index int, entry *CostBasisEntry, asset *AssetCostBasis, log *tradeLog, transaction *TransactionEntry) {
	oldValue := 0
	newValue := 1

	if index == 1 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.Quantity[oldValue] = 0
		entry.BalanceRemaining.Quantity[newValue] = math.Floor(entry.Allocation.Quantity*math.Pow(10, 8)) / math.Pow(10, 8)
		entry.BalanceRemaining.Amount = entry.QuotePriceEntry * entry.BalanceRemaining.Quantity[newValue]
		entry.BalanceRemaining.Value = math.Floor(entry.BalanceRemaining.Amount*100) / 100

	} else if index == 2 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.Quantity[oldValue] = asset.Quantity
		entry.BalanceRemaining.Quantity[newValue] = math.Floor((entry.BalanceRemaining.Quantity[oldValue]-entry.Allocation.Quantity)*math.Pow(10, 8)) / math.Pow(10, 8)
		entry.BalanceRemaining.Amount = entry.QuotePriceEntry * entry.BalanceRemaining.Quantity[newValue]
		entry.BalanceRemaining.Value = math.Floor(entry.BalanceRemaining.Amount*100) / 100

	} else if index == 4 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.Quantity[oldValue] = asset.Quantity
		entry.BalanceRemaining.Quantity[newValue] = math.Floor((entry.BalanceRemaining.Quantity[oldValue]-entry.Allocation.Quantity)*math.Pow(10, 8)) / math.Pow(10, 8)
		entry.BalanceRemaining.Amount = entry.QuotePriceEntry * entry.BalanceRemaining.Quantity[newValue]
		entry.BalanceRemaining.Value = math.Floor(entry.USDPriceEntry*entry.BalanceRemaining.Quantity[newValue]*100) / 100
		// entry.BalanceRemaining.Value = entry.USDPriceEntry * entry.BalanceRemaining.Amount

	} else if index == 5 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.Quantity[oldValue] = 0
		entry.BalanceRemaining.Quantity[newValue] = math.Floor(entry.Allocation.Quantity*math.Pow(10, 8)) / math.Pow(10, 8)
		entry.BalanceRemaining.Amount = entry.QuotePriceEntry * entry.BalanceRemaining.Quantity[newValue]
		entry.BalanceRemaining.Value = math.Floor(entry.USDPriceEntry*entry.BalanceRemaining.Amount*100) / 100

	} else if index == 6 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.Quantity[oldValue] = asset.Quantity
		entry.BalanceRemaining.Quantity[newValue] = math.Floor((entry.BalanceRemaining.Quantity[oldValue]-entry.Allocation.Quantity)*math.Pow(10, 8)) / math.Pow(10, 8)
		entry.BalanceRemaining.Amount = entry.QuotePriceEntry * entry.BalanceRemaining.Quantity[newValue]
		entry.BalanceRemaining.Value = math.Floor(entry.USDPriceEntry*entry.BalanceRemaining.Amount*100) / 100

	} else if index == 7 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.Quantity[oldValue] = 0
		entry.BalanceRemaining.Quantity[newValue] = math.Floor(entry.Allocation.Quantity*math.Pow(10, 8)) / math.Pow(10, 8)
		entry.BalanceRemaining.Amount = entry.QuotePriceEntry * entry.BalanceRemaining.Quantity[newValue]
		entry.BalanceRemaining.Value = math.Floor(entry.BalanceRemaining.Amount*100) / 100

	}
}

func updateHoldings(index int, entry *CostBasisEntry, asset *AssetCostBasis, log *tradeLog, transaction *TransactionEntry) {
	// what should it do or consider?
	entry.Holdings.Base = log.balance.base
	entry.Holdings.Quote = log.balance.quote
}

func updatePNL(index int, entry *CostBasisEntry, asset *AssetCostBasis, log *tradeLog, transaction *TransactionEntry) {

	// several ways to implement tracking of unrealized PNL
	// record of last traded price of order pair and last USD Price value of order pair
	// for every trade iterate over asset balance and pull from api the price and value of traded date
	// real time -> ticker rate -> pull price data -> iterate over asset balance and compute unreal holdings

	// for every trade need to iterate over asset balance
	// 	-> pull price data from api for each asset
	//	-> or record last traded price of that asset and use that info if no api access (not as accurate)
	//	-> compute order price with asset balance  (asset balance / order price)  order pair balance in quote
	//	-> computer usd price with asset balance in quote value -> (asset balance in quote value * usd price)

	// iterate on a ticker rate (most up to date and accurate)
	// 	-> pull price data from api for each asset
	//	-> compute order price with asset balance  (asset balance / order price)  order pair balance in quote
	//	-> computer usd price with asset balance in quote value -> (asset balance in quote value * usd price)

	// in this function we will not update the unrealized PNL and will not be on cost basis entry

	// isAppend [1, 5, 7]
	// isDeduct isUSD [2]
	// isDeduct isCrypto isBuy [4]
	// isDeduct isCrypto isSell [6]

	if index == 1 || index == 5 || index == 7 {
		entry.PNL.Amount = 0
		entry.PNL.Total = 0

	} else {
		// have unnessary conditions doing unnessary calcultions. need figure out want to keep and clean up.
		// doesn't account for fee amount
		entry.PNL.Amount = entry.Allocation.Value - (entry.USDPriceEntry * entry.Allocation.Quantity) // - (transaction.fee || 0);

		if index == 2 {
		} else if index == 4 {
		} else if index == 6 {
			entry.PNL.Amount = entry.Allocation.Value - (entry.USDPriceEntry * entry.Allocation.Amount) // - (transaction.fee || 0)
		}

		entry.PNL.Total = math.Floor(log.statement.PNL*100+entry.PNL.Amount*100) / 100
		log.statement.PNL = entry.PNL.Total
	}
}

func (e *CostBasisEntry) lastQuotePrice() float64 {
	if e.QuotePriceExit != 0 {
		return e.QuotePriceExit
	}
	return e.QuotePriceEntry
}

func (e *CostBasisEntry) lastUSDPrice() float64 {
	if e.USDPriceExit != 0 {
		return e.USDPriceExit
	}
	return e.USDPriceEntry
}

func (e *CostBasisEntry) filter(properties []string) []string {
	var t table = make(map[string]string)
	t = map[string]string{
		"Transaction ID -> Credit": fmt.Sprint(e.TransactionID.Credit),                   // Transaction ID format
		"Transaction ID -> Debit":  fmt.Sprint(e.TransactionID.Debit),                    // Transaction ID format
		"Quote Price -> Entry":     formatCurrency(e.QuotePriceEntry, e.quote()),         // formatCurrency
		"Quote Price -> Exit":      formatCurrency(e.QuotePriceExit, e.quote()),          // formatCurrency
		"USD Price -> Entry":       dollarFormat(e.USDPriceEntry),                        // dollarFormat
		"USD Price -> Exit":        dollarFormat(e.USDPriceExit),                         // dollarFormat
		"Allocation -> Quantity":   cryptoFormat(e.Allocation.Quantity),                  // cryptoFormat
		"Allocation -> Amount":     formatCurrency(e.Allocation.Amount, e.quote()),       // formatCurrency
		"Allocation -> Value":      dollarFormat(e.Allocation.Value),                     // dollarFormat
		"Balance -> Quantity":      fmt.Sprint(e.BalanceRemaining.Quantity),              // cryptoFormat
		"Balance -> Amount":        formatCurrency(e.BalanceRemaining.Amount, e.quote()), // formatCurrency
		"Balance -> Value":         dollarFormat(e.BalanceRemaining.Value),               // dollarFormat
		"Holdings -> Balance":      cryptoFormat(e.Holdings.Base),                        // cryptoFormat
		"PNL -> Amount":            dollarFormat(e.PNL.Amount),                           // dollarFormat
		"PNL -> Total":             dollarFormat(e.PNL.Total)}                            // dollarFormat

	if e.QuotePriceExit == 0.0 {
		t["Quote Price -> Exit"] = "-"
		t["USD Price -> Exit"] = "-"
		t["PNL -> Amount"] = "-"
	}

	return t.filter(properties)
}

// build
// create
// raise
// construct
// define

// middleware
// closure
