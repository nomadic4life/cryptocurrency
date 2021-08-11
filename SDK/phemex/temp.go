package phemex

type AccountData struct {
	Account struct {
		AccountId          int64
		Currency           string
		AccountBalanceEv   int64
		TotalUsedBalanceEv int64
	}
	Positions []Position
}

type ActiveList struct {
	Orders []Order
}

type Order struct {
	bizError       int64   // int64
	orderID        string  // UUID string
	clOrdID        string  // UUID string
	symbol         string  // string
	side           string  // string
	actionTimeNs   int64   // int64
	transactTimeNs int64   // int64
	orderType      string  // string?
	priceEp        int64   // int64
	price          int64   // int64
	orderQty       int64   // int64
	displayQty     int64   // int64
	timeInForce    string  // string?
	reduceOnly     bool    // bool
	stopPxEp       int64   // int64
	closedPnlEv    int64   // int64
	closedPnl      int64   // int64
	closedSize     int64   // int64
	cumQty         int64   // int64
	cumValueEv     int64   // int64
	cumValue       int64   // int64
	leavesQty      int64   // int64
	leavesValueEv  int64   // int64
	leavesValue    float64 // float64
	stopPx         int64   // int64
	stopDirection  string  // string
	ordStatus      string  // string

	Currency     string
	Action       string
	TradeType    string
	ExecQty      int64
	ExectPriceEP int64
	ExecValueEv  int64
	FeeRateEr    int64
	ExecFeeEv    int64
	Ordtype      string
	ExecStatus   string
}

type MarketData struct {
	Err     string
	Id      int64
	Results struct {
		Books []BookData
	}
	Depth     int64
	Sequence  int64
	Timestamp int64
	Symbol    string
	Type      string
}

type BookData struct {
	Asks [][]int64
	Bids [][]int64
}

type TradeData struct {
	Err     string
	Id      int64
	Results struct {
		Type     string
		Sequence int64
		Symbol   string
		Trades   [][]interface{}
	}
}

type TickerData struct {
	Err     string
	Id      int64
	Results struct {
		Open            int64
		High            int64
		Low             int64
		Close           int64
		IndexPrice      int64
		MarkPrice       int64
		OpenInterest    int64
		FundingRate     int64
		PredFundingRate int64
		Symbol          string
		Turnover        int64
		Volume          int64
		TimeStamp       int64
	}
}

type Nomics struct {
	Data []struct {
		Id           string
		Amount_Quote string
		Price        string
		Side         string
		Timestamp    string
		Type         string
	}
	MSG string
}

// ASSET API LIST
type Users struct {
	Code int64
	MSG  string
	Data struct {
		Total int64
		Rows  []struct {
			UserId        string
			Email         string
			NickName      string
			PasswordState int64
			ClientCnt     int64
			TOTP          int64
			Logon         int64
			ParentId      int64
			ParentEmail   string
			Status        int64
			Wallet        struct {
				TotalBalance    string
				TotalBalanceEv  int64
				AvailBalance    string
				AvailBalanceEv  int64
				FreezeBalance   string
				FreezeBalanceEv int64
				Currency        string
				CurrencyCode    int64
			}
			UserMarginVo []struct {
				Currency           string
				AccountBalance     string
				TotalUsedBalance   string
				AccountBalanceEv   int64
				TotalUsedBalanceEv int64
				BonusBalanceEv     int64
				BonusBalance       string
			}
		}
	}
}

type Wallet struct {
	Amount    int64
	AmountEv  int64
	ClientCnt int64
	Currency  string
}

type ResponseWallet struct {
	Code int64
	MSG  string
	Data string
}

type Margins struct {
	BTCAmount   float64
	BTCAmountEv int64
	LinkKey     string
	MoveOp      []int64
	USDAmount   float64
	USDAmountEv float64
}

type ResponseMargin struct {
	Code int64
	MSG  string
	Data struct {
		MoveOp           int64
		FromCurrencyName string
		ToCurrencyName   string
		FromAmount       string
		ToAmount         string
		LinkKey          string
		Status           int64
	}
}

type ResponseWalletHistory struct {
	Code int64
	MSG  string
	Data struct {
		Total int64
		Rows  []struct {
			MoveOp           int64
			FromCurrencyName string
			ToCurrencyName   string
			FromAmount       string
			ToAmount         string
			LinkKey          string
			Status           int64
			CreateTime       int64
		}
	}
}

type WithDraw struct {
	Address  string
	Amountev int64
	Currency string
}

type ResonseWithDraw struct {
	Code int64
	MSG  string
	Data struct {
		Id          int64
		Currency    string
		Status      string
		AmountEv    int64
		FeeEv       int64
		Address     string
		TxHash      string
		SubmitedAt  int64
		Expiredtime int64
	}
}

type ResponseConfirmWithDraw struct {
	Code int64
	MSG  string
}

type CancelWithdraw struct {
	Id int64
}

type WithdrawAddress struct {
	Address  string
	Currency string
	Remark   string
}

type ResponseConfirmAddress struct {
	Code int64
	MSG  string
	Data int64
}

type Query struct {
	currency, ordStatus, symbol, orderID, origClOrdID, clOrdID,
	price, priceep, orderQty, stopPx, stopPxEp, takeProfit, takeProfitEP,
	stopLoss, stopLossEp, pegOffsetValueEp, pegPriceType, untriggered,
	leverage, leverageEr, riskLimit, riskLimitEv, posBalance, posBalanceEv,
	start, end, offset, limit, tradeType, withCount, market, since, optCode, code string
}

type Body struct {
	symbol, clOrdID, side,
	priceEp, orderQty, ordType,
	reduceOnly, timeInForce, takeProfitEp,
	stopLossEp, actionBy, pegPriceType,
	pegOffsetValueEp, stopPxEp, closeOnTrigger,
	triggerType, address, amountEv, currency, remark string
}
