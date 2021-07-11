package tax

type Account struct {
	Statement struct {
		TotalCapital float64
		PNL          float64
	}
	Assets struct {
		Holdings            map[string]float64
		CostBasisAssetQueue map[string][]AssetTrade
	}
	Ledger struct {
		Transactions []TransactionEntry
		CostBases    []CostBasisEntry
	}
}

type TransactionEntry struct {
	TransactionID int64
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
	TransactionID
	QuotePriceEntry float64
	QuotePriceExit  float64
	USDPriceEntry   float64
	USDPriceExit    float64
	ChangeAmount    struct {
		BaseQuantity float64
		QuoteAmount  float64
		USDValue     float64
	}
	BalanceRemaining struct {
		BaseAmount [2]float64
		BaseValue  float64
		USDValue   float64
	}
	Holdings struct {
		TotalBaseBalance float64
		UnrealizedPNL    float64
	}
	PNL struct {
		Amount float64
		Total  float64
	}
}

type AssetTrade struct {
	TransactionID
	QuotePrice    float64
	USDPriceValue float64
	BaseAmount    float64 // debit
	ChangeAmount  float64 // credit
}

type TransactionID struct {
	From int64
	To   int64
}

type TradeInput struct { // InitTrade, TradeInput
	Date     int64
	Pair     string
	Type     string
	Price    float64
	Quantity float64
	Amount   float64
	Value    float64
	Fee      float64
}

// added comments that reflect the changes of amount going in and out of asset trade

// type trade struct {
// 	symbol struct {
// 		deduct string
// 		append string
// 	}
// 	balance struct {
// 		deduct float64
// 		quote  float64
// 		base   float64
// 	}
// 	amountDeducted float64
// 	PNL            float64
// 	unrealizedPNL  float64
// 	assetRecords   []CostBasisEntry
// 	queue          struct {
// 		quote []AssetTrade
// 		base  []AssetTrade
// 	}
// }
