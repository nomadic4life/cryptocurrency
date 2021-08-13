package phemex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	Req        *http.Request
	Header     *http.Header
	Method     string
	URL        string // endpoint + path + query
	Path       string
	Query      string
	Body       []byte
	Expiry     string
	Signature  string // HEADER ->  x-phemex-request-signature
	Signed     string
	HMACSHA256 string // URL Path + QueryString + Expiry + body
}

func (r *Request) setPath(method, path string) {
	r.Path = path
	r.Method = method
	r.URL = client.HostHTTP + r.Path
}

func (r *Request) setQuery(query map[string]string) {
	if query != nil {
		list := make([]string, 0, len(query))
		for key, element := range query {
			list = append(list, key+"="+element)
		}
		r.Query = strings.Join(list, "&")
		r.URL += "?"
		r.URL += r.Query
	}
}

func (r *Request) setBody(body map[string]interface{}) {
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			panic("OH shit!")
		}
		r.Body = data
	}
}

func (r *Request) setRequest() {
	if len(r.Body) == 0 {
		req, err := http.NewRequest(r.Method, r.URL, nil)
		if err != nil {
			panic("Holy Shit")
		}
		r.Req = req
		return
	}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))

	if err != nil {
		panic("Holy Shit")
	}
	r.Req = req
}

func (r *Request) isPrivate() bool {
	if r.Path == "/exchange/public/nomics/trades" || r.Path == "/exchange/public/products" {
		return false
	}
	return true
}

func (r *Request) sign(a *Account) {

	if r.isPrivate() {
		minute := 60
		time := int(time.Now().Unix())

		r.Expiry = strconv.Itoa(time + minute)

		byteMessage := []byte(r.Path + r.Query + r.Expiry + string(r.Body))

		a.hmac.Write(byteMessage)
		r.Signature = fmt.Sprintf("%x", a.hmac.Sum(nil))

		a.hmac.Reset()

		r.Req.Header.Add("x-phemex-access-token", a.API_KEY)
		r.Req.Header.Add("x-phemex-request-expiry", r.Expiry)
		r.Req.Header.Add("x-phemex-request-signature", r.Signature)
	}
}

func (r *Request) send(res *Response) {
	r.Req.Header.Add("content-type", "application/json")
	response, err := client.conn.Do(r.Req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error: %s\n", err)
		return
	}

	res.data, _ = ioutil.ReadAll(response.Body)
	res.req = r
}

// :: /orders
//	:Request:
//	->	<ORDERINPUT>
//	->	symbol <- Required
//	->	clOrdID <- Required
//	-> 	side <- Required
//	->	orderQty <- Required
//	:Response:
//	-> 	<ORDER>

// :: /orders/replace
//	Q->	<ORDERINPUT>
//	Q-> symbol <- required
//	Q-> orderID <- required
//	Q-> origClOrdID
//	:Response:
//	-> 	<ORDER>

// :: /orders/cancel
//	Q-> symbol
//	Q-> orderID
//	:Response:
//	-> 	<ORDER>

// :: /orders
//	Q-> symbol
//	Q-> []orderID
//	:Response:
//	-> 	[]<ORDER>

// :: /orders/all
//	Q-> symbol
//	Q-> untriggered <- [true, false]
//	Q-> text
//	:Response:
//	-> 	subject to change

// :: /accounts/accountPositions
// :: /accounts/positions <- unrealizedPnL (high rate limit)
//	Q-> currency
//	:Response:
//	->	account: {}
// 	->	->	<ACCOUNT>
//	-> 	positions: [{},...]
// 	->	-> 	<POSITION>

// Inverse contract: unRealizedPnl = (posSize/contractSize) / avgEntryPrice - (posSize/contractSize) / markPrice)
// Linear contract:  unRealizedPnl = (posSize/contractSize) * markPrice - (posSize/contractSize) * avgEntryPrice
// posSize is a signed vaule. contractSize is a fixed value.

// :: /orders/activeList <- open orders :: <-O = open order status
//	Q->	symbol
//	:Response:
//	-> rows: [{}]
//	->	->	<ORDER>

// :: /orders/active
//	Q-> orderID
//	:Response: ???

// :: /positions/leverage
//	Q-> symbol
//	Q-> leverage
//	Q-> leverageEr
//	:Response:
//	msg: "OK"

// :: /positions/riskLimit
//	Q-> symbol
//	Q-> riskLimit
//	Q-> riskLimitEv
//	:Response: ???

// :: /positions/assign
//	Q-> symbol
//	Q-> PosBalance
//	Q-> PosBalanceEv
//	:Response: ???

// :: /exchange/order/list <- Closed Order Status
//	Q->	symbol
//	Q->	start <- Epoch Millis
//	Q->	end <- Epoch Millis
//	Q->	offset
//	Q->	limit
//	Q->	ordStatus <- [New, PartiallyFilled, Untriggered, Filled, Canceled]
//	:Response:
//	-> rows: [{}]
//	->	->	<ORDER>

// :: /exchange/order <- Order by orderID's or clOrdID's
//	Q->	symbol
//	Q->	orderID
//	Q->	clOrdID
//	:Response:
//	-> rows: [{}]
//	->	->	<ORDER>

// :: /exchange/public/nomics/trades
//	Q-> market: symbol
//	Q-> since: "0-0-0"
//	:Response:
//	-> data: []
//	->	-> <NOMICS>

// :: /phemex-user/users/children
//	Q-> offset
//	Q-> limit
//	Q-> withCount
//	:Response:
//	-> rows: [{}]
// 	->	->	"userId": 6XXX12,
//	->	->	"email": "x**@**.com",
//	->	->	"nickName": "nickName",
//	->	->	"passwordState": 1,
//	->	->	"clientCnt": 0,
//	->	->	"totp": 1,
//	->	->	"logon": 0,
//	->	->	"parentId": 0,
//	->	->	"parentEmail": null,
//	->	->	"status": 1,
//	->	->	"wallet": {}
//	->	->	->	"totalBalance": "989.25471319",
//	->	->	->	"totalBalanceEv": 98925471319,
//	->	->	->	"availBalance": "989.05471319",
//	->	->	->	"availBalanceEv": 98905471319,
//	->	->	->	"freezeBalance": "0.20000000",
//	->	->	->	"freezeBalanceEv": 20000000,
//	->	->	->	"currency": "BTC",
//	->	->	->	"currencyCode": 1
//	->	->	"userMarginVo": [{}]
//	->	-> 	->	"currency": "BTC",
//	->	-> 	->	"accountBalance": "3.90032508",
//	->	-> 	->	"totalUsedBalance": "0.00015666",
//	->	-> 	->	"accountBalanceEv": 390032508,
//	->	-> 	->	"totalUsedBalanceEv": 15666,
//	->	-> 	->	"bonusBalanceEv": 0,
//	->	-> 	->	"bonusBalance": "0"

// :: /exchange/wallets/transferOut
//	:Request:
//	->	"amount": 0, // unscaled amount
//	->	"amountEv": 0, // scaled amount, when both amount and amountEv are provided, amountEv wins.
//	->	"clientCnt": 0, // client number, this is from API in children list; when sub-client issues this API, client must be 0.
//	->	"currency": "string"
//	:Response:
//	->"OK"

// :: /exchange/wallets/transferIn
//	:Request:
//	->	"amount": 0, // unscaled amount
//	->	"amountEv": 0, // scaled amount, when both amount and amountEv are provided, amountEv wins.
//	->	"clientCnt": 0, // client number, this is from API in children list; when sub-client issues this API, client must be 0.
//	->	"currency": "string"
//	:Response:
//	->"OK"

// :: /exchange/margins
//	:Request:
// 	->	"btcAmount": 0.00,
// 	->	"btcAmountEv": 0,
// 	->	"linkKey": "unique-str-for-this-request",
// 	->	"moveOp": [1,2,3,4],
// 	->	"usdAmount": 0.00,
// 	->	"usdAmountEv": 0
//	:Response:
// 	->	"moveOp": 1,
// 	->	"fromCurrencyName": "BTC",
// 	->	"toCurrencyName": "BTC",
// 	->	"fromAmount": "0.10000000",
// 	->	"toAmount": "0.10000000",
// 	->	"linkKey": "2431ca9b-2dd4-44b8-91f3-2539bb62db2d",
// 	->	"status": 10,

// :: /exchange/margins/transfer
//	Q->start
//	Q->end
//	Q->offset
//	Q->limit
//	Q->withCount [true, false]
//	:Response:
//	-> rows: [{}]
// 	->	->	"moveOp": 1,
// 	->	->	"fromCurrencyName": "BTC",
// 	->	->	"toCurrencyName": "BTC",
// 	->	->	"fromAmount": "0.10000000",
// 	->	->	"toAmount": "0.10000000",
// 	->	->	"linkKey": "2431ca9b-2dd4-44b8-91f3-2539bb62db2d",
// 	->	->	"status": 10,
// 	->	->	"createTime": 1580201427000

// :: /exchange/wallets/createWithdraw
//	Q->optCode
//	:Request:
// 	->	"address": <address>,// address must set before withdraw
// 	->	"amountEv": <amountEv>, // scaled btc value
// 	->	"currency": <currency> // fixed to BTC
//	:Response:
// 	->	id
// 	->	currency
// 	->	status
// 	->	amountEv
// 	->	feeEv
// 	->	address
// 	->	txhash
// 	->	submitedAt
// 	->	expiredTime

// :: /exchange/wallets/confirm/withdraw
//	Q->code=<withdrawConfirmCode>
//	:Response:
//	->ok

// :: /exchange/wallets/cancelWithdraw
//	:Request:
// id: <withdrawRequestId>
//	:Response:???

// :: /exchange/wallets/withdrawList
//	Q->currency
//	Q->limit
//	Q->offset
//	Q->withCount [true, false]
//	:Response: ???

// :: /exchange/wallets/createWithdrawAddress
//	Q->optCode
//	:Request:
// 	->	"address": <address>,
// 	->	"currency": <currency>
// 	->	"remark": <name>
// :Response: sumject to change

// <ORDER>
// 	->	"bizError": 0,
// 	->	"orderID": "9cb95282-7840-42d6-9768-ab8901385a67",
// 	->	"clOrdID": "7eaa9987-928c-652e-cc6a-82fc35641706",
// 	->	"symbol": "BTCUSD",
// 	->	"side": "Buy",
// 	->	"actionTimeNs": 1580533011677666800,
// 	->	"transactTimeNs": 1580533011677666800,
// 	->	"orderType": null,
// 	->	"priceEp": 84000000,
// 	->	"price": 8400,
// 	->	"orderQty": 1,
// 	->	"displayQty": 1,
// 	->	"timeInForce": null,
// 	->	"reduceOnly": false,
// 	->	"stopPxEp": 0,
// 	->	"closedPnlEv": 0,
// 	->	"closedPnl": 0,
// 	->	"closedSize": 0,
// 	->	"cumQty": 0,
// 	->	"cumValueEv": 0,
// 	->	"cumValue": 0,
// 	->	"leavesQty": 0,
// 	->	"leavesValueEv": 0,
// 	->	"leavesValue": 0,
// 	->	"stopPx": 0,
// 	->	"stopDirection": "Falling",
// 	->	"ordStatus": "Untriggered" <- [New <-O, PartiallyFilled <-O, Filled, Canceld, Rejected, Triggered, Untriggered <-O]

// <USER TRADE>
// 	->	"transactTimeNs": 1578026629824704800,
// 	->	"symbol": "BTCUSD",
// 	->	"currency": "BTC",
// 	->	"action": "Replace",
// 	->	"side": "Sell",
// 	->	"tradeType": "Trade", <- [Trade, Funding, AdlTrade, LiqTrade]
// 	->	"execQty": 700,
// 	->	"execPriceEp": 71500000,
// 	->	"orderQty": 700,
// 	->	"priceEp": 71500000,
// 	->	"execValueEv": 9790209,
// 	->	"feeRateEr": -25000,
// 	->	"execFeeEv": -2447,
// 	->	"ordType": "Limit",
// 	->	"execID": "b01671a1-5ddc-5def-b80a-5311522fd4bf",
// 	->	"orderID": "b63bc982-be3a-45e0-8974-43d6375fb626",
// 	->	"clOrdID": "uuid-1577463487504",
// 	->	"execStatus": "MakerFill" [Init, MakerFill, TakerFill]

// <ORDERBOOK>
//	book: {}
//	-> 	asks: [[]]
//	->	->	priceEP
//	->	->	size
//	-> 	bids: [[]]
//	->	->	priceEP
//	->	->	size
// 	->	"depth": 30,
// 	->	"sequence": <sequence>,
// 	->	"timestamp": <timestamp>,
// 	->	"symbol": "<symbol>",
// 	->	"type": "snapshot" <- [snapshot, interval]

// <RECENTTRADES>
// 	->	"type": "snapshot", <- [snapshot, interval]
// 	->	"sequence": <sequence>,
// 	->	"symbol": "<symbol>",
//	-> 	"trades": [[]]
// 	->	->	<timestamp>,
// 	->	->	"<side>",
// 	->	->	<priceEp>,
// 	->	->	<size>

// <24HRTICKER>
// 	->	"open": <open priceEp>,
// 	->	"high": <high priceEp>,
// 	->	"low": <low priceEp>,
// 	->	"close": <close priceEp>,
// 	->	"indexPrice": <index priceEp>,
// 	->	"markPrice": <mark priceEp>,
// 	->	"openInterest": <open interest>,
// 	->	"fundingRate": <funding rateEr>,
// 	->	"predFundingRate": <predicated funding rateEr>,
// 	->	"symbol": "<symbol>",
// 	->	"turnover": <turnoverEv>,
// 	->	"volume": <volume>,
// 	->	"timestamp": <timestamp>

// <ORDERINPUT>
// 	->	"actionBy": "FromOrderPlacement",
// 	->	"symbol": "BTCUSD",
//	->	"origClOrdID": uuid-1573058952273",		<- Query input <AMEND>
// 	->	"clOrdID": "uuid-1573058952273",
// 	->	"side": "Sell",
// 	->	"priceEp": 93185000,
// 	->	"orderQty": 7,
// 	->	"ordType": "Limit",
// 	->	"reduceOnly": false,
// 	->	"triggerType": "UNSPECIFIED",
// 	->	"pegPriceType": "UNSPECIFIED",
// 	->	"timeInForce": "GoodTillCancel",
// 	->	"takeProfitEp": 0,
// 	->	"stopLossEp": 0,
// 	->	"pegOffsetValueEp": 0,
// 	->	"pegPriceType": "UNSPECIFIED"

// <POSITION>
// 	->	"accountID": 0,
// 	->	"symbol": "BTCUSD",
// 	->	"currency": "BTC",
// 	->	"side": "None",
// 	->	"positionStatus": "Normal",
// 	->	"crossMargin": false,
// 	->	"leverageEr": 0,
// 	->	"leverage": 0,
// 	->	"initMarginReqEr": 0,
// 	->	"initMarginReq": 0.01,
// 	->	"maintMarginReqEr": 500000,
// 	->	"maintMarginReq": 0.005,
// 	->	"riskLimitEv": 10000000000,
// 	->	"riskLimit": 100,
// 	->	"size": 0,
// 	->	"value": 0,
// 	->	"valueEv": 0,
// 	->	"avgEntryPriceEp": 0,
// 	->	"avgEntryPrice": 0,
// 	->	"posCostEv": 0,
// 	->	"posCost": 0,
// 	->	"assignedPosBalanceEv": 0,
// 	->	"assignedPosBalance": 0,
// 	->	"bankruptCommEv": 0,
// 	->	"bankruptComm": 0,
// 	->	"bankruptPriceEp": 0,
// 	->	"bankruptPrice": 0,
// 	->	"positionMarginEv": 0,
// 	->	"positionMargin": 0,
// 	->	"liquidationPriceEp": 0,
// 	->	"liquidationPrice": 0,
// 	->	"deleveragePercentileEr": 0,
// 	->	"deleveragePercentile": 0,
// 	->	"buyValueToCostEr": 1150750,
// 	->	"buyValueToCost": 0.0115075,
// 	->	"sellValueToCostEr": 1149250,
// 	->	"sellValueToCost": 0.0114925,
// 	->	"markPriceEp": 93169002,
// 	->	"markPrice": 9316.9002,
// 	->	"markValueEv": 0,
// 	->	"markValue": null,
// 	->	"estimatedOrdLossEv": 0,
// 	->	"estimatedOrdLoss": 0,
// 	->	"usedBalanceEv": 0,
// 	->	"usedBalance": 0,
// 	->	"takeProfitEp": 0,
// 	->	"takeProfit": null,
// 	->	"stopLossEp": 0,
// 	->	"stopLoss": null,
// 	->	"realisedPnlEv": 0,
// 	->	"realisedPnl": null,
//	->	"unRealisedPnlEv": 0, <-- calculated on client side :: or from /accounts/positions
//	->	"unRealisedPnl": null, <-- calculated on client side
// 	->	"cumRealisedPnlEv": 0,
// 	->	"cumRealisedPnl": null

// <ACCOUNT>
// 	->	"accountId": 0,
// 	->	"currency": "BTC",
// 	->	"accountBalanceEv": 0,
// 	->	"totalUsedBalanceEv": 0

// <NOMICS>
//	->	"id": "string",
//	->	"amount_quote": "string",
//	->	"price": "string",
//	->	"side": "string",
//	->	"timestamp": "string",
//	->	"type": "string"
