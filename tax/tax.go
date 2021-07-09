package tax

type Account struct {
	Statement struct {
		TotalCapital float64
		PNL          float64
	}
	AssetsHoldings map[string]float64
	Ledger         struct {
		TransactionHistory []Transaction
		CostBasesHistory   []CostBasisEntry
	}

	CostBasisAssetQueue map[string][]AssetTrade
}

type Transaction struct {
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

type Trade struct {
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
