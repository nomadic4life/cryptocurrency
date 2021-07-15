package tax

type Account struct {
	Statement struct {
		TotalCapital float64
		PNL          float64
	}
	Assets struct {
		Holdings            map[string]float64
		AssetCostBasesQueue map[string][]AssetCostBasis
	}
	Ledger struct {
		// Transactions []TransactionEntry
		// CostBases []CostBasisEntry
		Transactions transactionList
		CostBases    costBasisList
	}
}

type TransactionEntry struct {
	TransactionID string
	Date          int64
	OrderPair     string
	OrderType     string
	OrderPrice    float64
	OrderQuantity float64
	OrderAmount   float64
	USDPriceValue float64
	FeeAmount     float64
}

type CostBasisEntry struct {
	meta struct {
		orderPair string
		orderType string
		date      int
	}
	TransactionID
	QuotePriceEntry float64
	QuotePriceExit  float64
	USDPriceEntry   float64
	USDPriceExit    float64
	Allocation      struct {
		Quantity float64
		Amount   float64 // not relevent
		Value    float64
	}
	BalanceRemaining struct {
		// Balance
		// Remaining <- subtracted from
		Quantity [2]float64
		Amount   float64 // not relevent
		Value    float64
	}
	Holdings struct {
		// NOT RELEVENT
		Base  float64
		Quote float64
	}
	PNL struct {
		Amount float64
		Total  float64 // not relevent
		// Unrealized float64
	}
}

type AssetCostBasis struct {
	TransactionID
	QuotePrice    float64
	USDPriceValue float64
	Quantity      float64 // debit
	Credit        float64 // credit
}

type TransactionID struct {
	Credit string
	Debit  string
}

type TradeInput struct {
	Date     int64
	Pair     string
	Type     string
	Price    float64
	Quantity float64
	Amount   float64
	Value    float64
	Fee      float64
}

type table map[string]string

// type prop []string

// type queue []interface{}

type transactionList []TransactionEntry

type costBasisList []CostBasisEntry

type assetCostBasisList []AssetCostBasis

func (t table) filter(properties []string) []string {

	results := make([]string, 0, 20)

	for i := 0; i < len(properties); i++ {
		key := properties[i]
		if _, ok := t[key]; ok {
			results = append(results, t[key])
		}
	}

	return results
}

// commits -> changed type names to reflect their purpose more accuretly

// need TODO:
// -> need to test on gains and losses
// -> need to fix bugs
// -> Read / Write JSON Data
// -> clean up of names, prints, comments, code layout
// -> refactor code, enqueue / dequeue, round calculations, inport / exports, other?
// -> alternative cost basis implementation ex.> specific identification

// -> API access:
// -> -> BTC/USD price
// -> -> Date

// -> market holdings
// -> -> note: how to implement this?
