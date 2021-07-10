package tax

import "fmt"

func newCostBasisEntry(asset *AssetTrade, trade *tradeLog, transaction *Transaction) *CostBasisEntry {
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
		updateChangeAmount,
		updateBalanceRemaining}
	return build(asset, trade, transaction, cb)
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

type builder func(index int, entry *CostBasisEntry, asset *AssetTrade, log *tradeLog, transaction *Transaction)

func checkCondition(orderType, orderPair string, asset *AssetTrade) exchange {
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

func build(asset *AssetTrade, log *tradeLog, transaction *Transaction, cb []builder) *CostBasisEntry {
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

	for _, next := range cb {
		next(index, &entry, asset, log, transaction)
	}

	return &entry
}

func createID(index int, entry *CostBasisEntry, asset *AssetTrade, log *tradeLog, transaction *Transaction) {
	// from -> debited
	// to -> credited
	if asset != nil {
		entry.TransactionID.From = asset.To
	} else {
		entry.TransactionID.From = transaction.TransactionID
	}

	entry.TransactionID.To = transaction.TransactionID

	fmt.Println("created ID", entry)
}

func executedPrice(index int, entry *CostBasisEntry, asset *AssetTrade, log *tradeLog, transaction *Transaction) {

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

	fmt.Println("Executed Prices", entry)
}

func updateChangeAmount(index int, entry *CostBasisEntry, asset *AssetTrade, log *tradeLog, transaction *Transaction) {

	value := [9]float64{
		transaction.OrderAmount,
		transaction.OrderAmount * entry.lastUSDPrice(),
		transaction.OrderQuantity,
		transaction.OrderQuantity * entry.lastQuotePrice(),
		transaction.OrderQuantity * entry.lastUSDPrice(),
		0.0,
		0.0,
		0.0,
		0.0}

	if asset != nil {
		value[5] = asset.ChangeAmount
		value[6] = asset.ChangeAmount * entry.lastQuotePrice()
		value[7] = asset.ChangeAmount * entry.lastUSDPrice()
		value[8] = entry.ChangeAmount.QuoteAmount * entry.lastUSDPrice()
	}

	var table map[string]map[int]int
	table = make(map[string]map[int]int)
	table["Base Quantity"] = map[int]int{7: 0, 1: 2, 5: 2, 2: 5, 4: 5, 6: 5}
	table["Quote Amount"] = map[int]int{7: 2, 1: 3, 5: 3, 2: 6, 4: 6, 6: 6}
	table["USDValue"] = map[int]int{7: 1, 1: 4, 2: 7, 4: 7, 5: 8, 6: 8}

	entry.ChangeAmount.BaseQuantity = value[table["Base Quantity"][index]]
	entry.ChangeAmount.QuoteAmount = value[table["Quote Amount"][index]]
	entry.ChangeAmount.USDValue = value[table["USDValue"][index]]

	fmt.Println("updated Change amount", entry)
}

func updateBalanceRemaining(index int, entry *CostBasisEntry, asset *AssetTrade, log *tradeLog, transaction *Transaction) {
	oldValue := 0
	newValue := 1

	if index == 1 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.BaseAmount[oldValue] = 0
		entry.BalanceRemaining.BaseAmount[newValue] = entry.ChangeAmount.BaseQuantity
		entry.BalanceRemaining.BaseValue = entry.QuotePriceEntry * entry.BalanceRemaining.BaseAmount[newValue]
		entry.BalanceRemaining.USDValue = entry.BalanceRemaining.BaseValue

	} else if index == 2 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.BaseAmount[oldValue] = asset.BaseAmount
		entry.BalanceRemaining.BaseAmount[newValue] = entry.BalanceRemaining.BaseAmount[oldValue] - entry.ChangeAmount.BaseQuantity
		entry.BalanceRemaining.BaseValue = entry.QuotePriceEntry * entry.BalanceRemaining.BaseAmount[newValue]
		entry.BalanceRemaining.USDValue = entry.BalanceRemaining.BaseValue

	} else if index == 4 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.BaseAmount[oldValue] = asset.BaseAmount
		entry.BalanceRemaining.BaseAmount[newValue] = entry.BalanceRemaining.BaseAmount[oldValue] - entry.ChangeAmount.BaseQuantity
		entry.BalanceRemaining.BaseValue = entry.QuotePriceEntry * entry.BalanceRemaining.BaseAmount[newValue]
		entry.BalanceRemaining.USDValue = entry.USDPriceEntry * entry.BalanceRemaining.BaseValue

	} else if index == 5 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.BaseAmount[oldValue] = 0
		entry.BalanceRemaining.BaseAmount[newValue] = entry.ChangeAmount.BaseQuantity
		entry.BalanceRemaining.BaseValue = entry.QuotePriceEntry * entry.BalanceRemaining.BaseAmount[newValue]
		entry.BalanceRemaining.USDValue = entry.USDPriceEntry * entry.BalanceRemaining.BaseValue

	} else if index == 6 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.BaseAmount[oldValue] = asset.BaseAmount
		entry.BalanceRemaining.BaseAmount[newValue] = entry.BalanceRemaining.BaseAmount[oldValue] - entry.ChangeAmount.BaseQuantity
		entry.BalanceRemaining.BaseValue = entry.QuotePriceEntry * entry.BalanceRemaining.BaseAmount[newValue]
		entry.BalanceRemaining.USDValue = entry.USDPriceEntry * entry.BalanceRemaining.BaseValue

	} else if index == 7 {

		//      -- BALANCE REMAINING --
		entry.BalanceRemaining.BaseAmount[oldValue] = 0
		entry.BalanceRemaining.BaseAmount[newValue] = entry.ChangeAmount.BaseQuantity
		entry.BalanceRemaining.BaseValue = entry.QuotePriceEntry * entry.BalanceRemaining.BaseAmount[newValue]
		entry.BalanceRemaining.USDValue = entry.BalanceRemaining.BaseValue

	}
	fmt.Println("updated Balance Remaining", entry, index)
}

func updateHoldings(index int, entry *CostBasisEntry, asset *AssetTrade, trade *tradeLog, transaction *Transaction) {

}

func updatePNL(index int, entry *CostBasisEntry, asset *AssetTrade, log *tradeLog, transaction *Transaction) {

	// isAppend [1, 5, 7]
	// isDeduct isUSD [2]
	// isDeduct isCrypto isBuy [4]
	// isDeduct isCrypto isSell [6]

	if index == 1 || index == 5 || index == 7 {
		entry.PNL.Amount = 0
		entry.PNL.Total = 0

	} else {
		// doesn't account for fee amount
		entry.PNL.Amount = entry.ChangeAmount.USDValue - (entry.USDPriceEntry * entry.ChangeAmount.BaseQuantity) // - (transaction.fee || 0);
		if index == 2 {
			entry.PNL.Total = entry.ChangeAmount.USDValue - (entry.USDPriceEntry * entry.ChangeAmount.BaseQuantity) // - (transaction.fee || 0)
		} else if index == 4 {
			entry.PNL.Total = entry.ChangeAmount.USDValue - (entry.QuotePriceEntry * entry.ChangeAmount.BaseQuantity * entry.USDPriceEntry) // - (transaction.fee || 0)
		} else if index == 6 {
			entry.PNL.Total = entry.ChangeAmount.USDValue - (entry.USDPriceEntry * entry.ChangeAmount.BaseQuantity) // - (transaction.fee || 0)
		}

	}

	// asset.PNL +=
	// this['P & L']['USD Total'] = trade.profitAndLoss;
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
	return e.QuotePriceEntry
}

// build
// create
// raise
// construct
// define

// middleware
// closure
