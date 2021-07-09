package tax

import "fmt"

func newCostBasisEntry(asset *AssetTrade, trade *trade, transaction *Transaction) *CostBasisEntry {
	fmt.Println("asset in newCostBasisEntry", asset)
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
		executedPrice}
	return build(asset, trade, transaction, cb)
}

// func build(t *Transaction, f func(){}) {

// }

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

type builder func(index int, entry *CostBasisEntry, asset *AssetTrade, trade *trade, transaction *Transaction)

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

func build(asset *AssetTrade, trade *trade, transaction *Transaction, cb []builder) *CostBasisEntry {
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
		next(index, &entry, asset, trade, transaction)
	}

	return &entry
}

func createID(index int, entry *CostBasisEntry, asset *AssetTrade, trade *trade, transaction *Transaction) {
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

func executedPrice(index int, entry *CostBasisEntry, asset *AssetTrade, trade *trade, transaction *Transaction) {

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

// build
// create
// raise
// construct
// define

// middleware
// closure
