package tax

import "fmt"

func newCostBasisEntry(t *Transaction) *CostBasisEntry {
	// middleware implementation
	// create ID
	// configue Excuted Price
	// update Change Amount
	// update Holdings
	// update PNL
	// display results?
	// build("Hello World!", func(message string) { fmt.Println(message) })
	return &CostBasisEntry{}
}

// func build(t *Transaction, f func(){}) {

// }

func build(greeting string, f []func(*string)) {
	fmt.Println("Build input: -> greeting", greeting)
	f[0](&greeting)
	test := func() {
		fmt.Println("Build closure: -> greeting", greeting)
	}

	test()

}

func createID(entry *CostBasisEntry, asset *AssetTrade, trade *trade, transaction *Transaction) {
	// from -> debited
	// to -> credited
	if asset != nil {
		entry.TransactionID.From = asset.To
	} else {
		entry.TransactionID.From = transaction.TransactionID
	}

	entry.TransactionID.To = transaction.TransactionID
}

func executedPrice(entry *CostBasisEntry, asset *AssetTrade, trade *trade, transaction *Transaction) {
	if transaction.OrderType == "BUY" {
		if transaction.quote() == "USD" {

		} else {

		}

	} else if transaction.OrderType == "BUY" {
		if transaction.quote() == "USD" {

		} else {

		}
	}
}

// build
// create
// raise
// construct
// define

// middleware
// closure
