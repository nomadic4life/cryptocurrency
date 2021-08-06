package phemex

import (
	"crypto"
	"crypto/hmac"
	_ "crypto/sha256"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// HEADERS:
// Every HTTP Rest Request must have the following Headers:

// x-phemex-access-token : This is API-KEY (id field) from Phemex site.
// x-phemex-request-expiry : This describes the Unix EPoch SECONDS to expire the request, normally it should be (Now() + 1 minute)
// x-phemex-request-signature : This is HMAC SHA256 signature of the http request. Secret is API Secret, its formula is : HMacSha256( URL Path + QueryString + Expiry + body )

// Optional Headers:
// x-phemex-request-tracing: a unique string to trace http-request, less than 40 bytes.

// X-RateLimit-Remaining-CONTRACT, # Remaining request permits in this minute
// X-RateLimit-Capacity-CONTRACT, # Request ratelimit capacity
// X-RateLimit-Retry-After-CONTRACT, # Reset timeout in seconds for current ratelimited user

var client = setupClient()

const (
	// Order Type Const
	LIMIT                      string = "Limit"
	MARKET                     string = "Market"
	STOP                       string = "Stop"
	STOP_LIMIT                 string = "StopLimit"
	MARKET_IF_TOUCHED          string = "MarketIfTouched"
	LIMIT_IF_TOUCHED           string = "LimitIfTouched"
	MARKET_AS_LIMIT            string = "MarketAsLimit"
	STOP_AS_LIMIT              string = "StopAsLimit"
	MARKET_IF_TOUCHED_AS_LIMIT string = "MarketIfTouchedAsLimit"

	// Order Status Const
	UNTRIGGERED      string = "Untriggered"     // Conditional order waiting to be triggered
	TRIGGERED        string = "Triggered"       // Conditional order being triggered
	REJECTED         string = "Rejected"        // Order rejected
	NEW              string = "New"             // Order placed in cross engine
	PARTIALLY_FILLED string = "PartiallyFilled" // Order partially filled
	FILLED           string = "Filled"          // Order fully filled
	CANCELED         string = "Canceled"        // Order canceled

	// timeInForce Const
	GOOD_TILL_CANCEL    string = "GoodTillCancel"
	POST_ONLY           string = "PostOnly"
	IMMEDIATE_OR_CANCEL string = "ImmediateOrCancel"
	FILL_OR_KILL        string = "FillOrKill"

	// Execution instruction
	REDUCE_ONLY      string = "ReduceOnly"     // reduce position size, never increase position size
	CLOSE_ON_TRIGGER string = "CloseOnTrigger" // close the position

	// Trigger source
	BY_MARK_PRICE string = "ByMarkPrice" // trigger by mark price
	BY_LAST_PRICE string = "ByLastPrice" // trigger by last price
)

// NET/HTTP package
// type Client, Header, Request, Response, File

// ID:   "2b6ba6d3-a14c-44ef-be5b-0e590ab35126"
// API Secret:   "cgJEMS_iIEnzg6DV9Ac-Eug3XcDy-fIugzeJhU0ohZ43YjJmMTI1Zi1lZDhiLTQ0YjItOGU4NS02YWEwYWUwOThhMjM"

type Client struct {
	conn   http.Client
	hmac   hash.Hash
	header *http.Header
	host   string
	apiID  string
	// socket websocket
}

type Request struct {
	Req        *http.Request
	Header     *http.Header
	Method     string
	URL        string // endpoint + path + query
	Path       string
	Query      string
	Body       string
	Expiry     string
	Signature  string // HEADER ->  x-phemex-request-signature
	Signed     string
	HMACSHA256 string // URL Path + QueryString + Expiry + body
	Private    bool
	// apiSecret = Base64::urlDecode(API Secret)
}

func setupClient() *Client {
	client := new(Client)
	client.conn = *http.DefaultClient
	client.hmac = hmac.New(crypto.SHA256.New, []byte("API SECRET"))
	client.host = "https://api.phemex.com"
	// read from config.json .env file
	//	-> api secret
	//	-> host -> production or test
	//	-> API-KEY -> headers -> x-phemex-access-token : (id field)

	return client
}
func (r *Request) Sign() {

	minute := 60
	time := int(time.Now().Unix())
	r.Expiry = strconv.Itoa(time + minute)

	byteMessage := []byte(r.Path + r.Query + r.Expiry + r.Body)

	client.hmac.Write(byteMessage)
	r.Signature = fmt.Sprintf("%x", client.hmac.Sum(nil))

	client.hmac.Reset()
}

// CreateRequest(method, url string, body io.ReadCloser)
func CreateRequest() {
	request := new(Request)
	// https://api.phemex.com/accounts/accountPositions?currency=BTC

	// req, err := http.NewRequest(method, url, body)
	// if err != nil {
	// 	panic("Holy Shit")
	// }
	// request.Req = req

	// examples
	// 	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	// resp, err := client.Do(req)

	request.Method = "GET"
	request.URL = "https://api.phemex.com/accounts/accountPositions?currency=BTC"
	request.Path = "/accounts/accountPositions"
	request.Query = "currency=BTC"
	request.Body = ""
	request.Private = false

	req, err := http.NewRequest("GET", "https://api.phemex.com/accounts/accountPositions?currency=BTC", nil)
	if err != nil {
		panic("Holy Shit")
	}

	// x-phemex-access-token : This is API-KEY (id field) from Phemex site.
	// x-phemex-request-expiry : This describes the Unix EPoch SECONDS to expire the request, normally it should be (Now() + 1 minute)
	// x-phemex-request-signature : This is HMAC SHA256 signature of the http request. Secret is API Secret, its formula is : HMacSha256( URL Path + QueryString + Expiry + body )

	request.Sign()
	req.Header.Add("x-phemex-access-token", client.apiID)
	req.Header.Add("x-phemex-request-expiry", request.Expiry)
	req.Header.Add("x-phemex-request-signature", request.Signature)

	response, err := client.conn.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error: %s\n", err)
	} else {
		var data []byte
		data, _ = ioutil.ReadAll(response.Body)

		// json.Unmarshal(data, &markets)
		fmt.Printf("RESPONSE %s\n", data)
		fmt.Println("\n", request.Signature)
	}
}

// rest endpoint
// 	 -> production 	 https://api.phemex.com
// 	 -> test/dev	 https://testnet-api.phemex.com

// wss endpoint
//	 -> production 	wss://phemex.com/ws
//	 -> test/dev		wss://testnet.phemex.com/ws

// http return codes
//	 -> [401] -> unauthenticated
//	 -> [403] -> lack of priviledge
//	 -> [429] -> breaking a request rate limit
//	 -> [5xx] -> Phemex internal errors (UNKNOWN, could be succeed)

// response schema
// {
//     "code": <code>, 	-> 0 == success, non-zero == error
//     "msg": <msg>,	-> error message on non-zero code
//     "data": <data>	-> operation dependant
// }

// HTTP Rest Requst Headers:
//	-> x-phemex-access-token:		-> This is API-KEY (id field) from Phemex site
//	-> x-phemex-request-expiry:		-> EPoch SECONDS -> (Now() + 1 minute)
//	-> x-phemex-request-signature: 	-> HMacSha256( URL Path + QueryString + Expiry + body )

// Optional Headers:
//	-> x-phemex-request-tracing: 	-> a unique string to trace http-request, less than 40 bytes

// https://api.phemex.com/exchange/public/products			<- Query Prducts
// https://api.phemex.com/md/orderbook?symbol=<symbol>		<- Query Order Book
// https://api.phemex.com/md/trade?symbol=<symbol>			<- Query Recent Trades
// https://api.phemex.com/md/ticker/24hr?symbol=<symbol>	<- Query 24 Hours Ticker
// https://api.phemex.com/exchange/public/nomics/trades?market=<symbol>&since=<since>	<- Query History Trades By symbol

// https://api.phemex.com/orders
// https://api.phemex.com/accounts/accountPositions?currency=<currency>
// https://api.phemex.com/accounts/positions?currency=<currency>
// https://api.phemex.com/positions/leverage?symbol=<symbol>&leverage=<leverage>&leverageEr=<leverageEr>
// https://api.phemex.com/positions/riskLimit?symbol=<symbol>&riskLimit=<riskLimit>&riskLimitEv=<riskLimitEv>

// - Get Products
// https://api.phemex.com/exchange/public/products
// GET /exchange/public/products
// {
// 		symbol: "BTCUSD",
// 		underlyingSymbol: ".BTC",
// 		quoteCurrency: "USD",
// 		baseCurrency: "BTC",
// 		settlementCurrency: "BTC",
// 		maxOrderQty: 1000000,
// 		maxPriceEp: 100000000000000,
// 		lotSize: 1,
// 		tickSize: "0.5",
// 		contractSize: "1 USD",
// 		priceScale: 4,
// 		ratioScale: 8,
// 		valueScale: 8,
// 		defaultLeverage: 0,
// 		maxLeverage: 100,
// 		initMarginEr: "1000000",
// 		maintMarginEr: "500000",
// 		defaultRiskLimitEv: 10000000000,
// 		deleverage: true,
// 		makerFeeRateEr: -250000,
// 		takerFeeRateEr: 750000,
// 		fundingInterval: 8,
// 		marketUrl: "https://phemex.com/trade/BTCUSD",
// 		description: "BTCUSD is a BTC/USD perpetual contract priced on the .BTC Index. Each contract is worth 1 USD of Bitcoin. Funding is paid and received every 8 hours. At UTC time: 00:00, 08:00, 16:00.",
// 		type: "Perpetual"
// 	}

// - Place Order
// https://api.phemex.com/orders
// POST /orders
// {
//   "actionBy": "FromOrderPlacement",
//   "symbol": "BTCUSD",
//   "clOrdID": "uuid-1573058952273",
//   "side": "Sell",
//   "priceEp": 93185000,
//   "orderQty": 7,
//   "ordType": "Limit",
//   "reduceOnly": false,
//   "triggerType": "UNSPECIFIED",
//   "pegPriceType": "UNSPECIFIED",
//   "timeInForce": "GoodTillCancel",
//   "takeProfitEp": 0,
//   "stopLossEp": 0,
//   "pegOffsetValueEp": 0,
//   "pegPriceType": "UNSPECIFIED"
// }

// HTTP Response:
// {
//     "code": 0,
//         "msg": "",
//         "data": {
//             "bizError": 0,
//             "orderID": "ab90a08c-b728-4b6b-97c4-36fa497335bf",
//             "clOrdID": "137e1928-5d25-fecd-dbd1-705ded659a4f",
//             "symbol": "BTCUSD",
//             "side": "Sell",
//             "actionTimeNs": 1580547265848034600,
//             "transactTimeNs": 0,
//             "orderType": null,
//             "priceEp": 98970000,
//             "price": 9897,
//             "orderQty": 1,
//             "displayQty": 1,
//             "timeInForce": null,
//             "reduceOnly": false,
//             "stopPxEp": 0,
//             "closedPnlEv": 0,
//             "closedPnl": 0,
//             "closedSize": 0,
//             "cumQty": 0,
//             "cumValueEv": 0,
//             "cumValue": 0,
//             "leavesQty": 1,
//             "leavesValueEv": 10104,
//             "leavesValue": 0.00010104,
//             "stopPx": 0,
//             "stopDirection": "UNSPECIFIED",
//             "ordStatus": "Created"
//         }
// }

// Query trading account and positions
// https://api.phemex.com/accounts/accountPositions?currency=<currency>
// GET /accounts/accountPositions?currency=<currency>

// Response
// {
//     "code": 0,
//         "msg": "",
//         "data": {
//             "account": {
//                 "accountId": 0,
//                 "currency": "BTC",
//                 "accountBalanceEv": 0,
//                 "totalUsedBalanceEv": 0
//             },
//             "positions": [
//             {
//                 "accountID": 0,
//                 "symbol": "BTCUSD",
//                 "currency": "BTC",
//                 "side": "None",
//                 "positionStatus": "Normal",
//                 "crossMargin": false,
//                 "leverageEr": 0,
//                 "leverage": 0,
//                 "initMarginReqEr": 0,
//                 "initMarginReq": 0.01,
//                 "maintMarginReqEr": 500000,
//                 "maintMarginReq": 0.005,
//                 "riskLimitEv": 10000000000,
//                 "riskLimit": 100,
//                 "size": 0,
//                 "value": 0,
//                 "valueEv": 0,
//                 "avgEntryPriceEp": 0,
//                 "avgEntryPrice": 0,
//                 "posCostEv": 0,
//                 "posCost": 0,
//                 "assignedPosBalanceEv": 0,
//                 "assignedPosBalance": 0,
//                 "bankruptCommEv": 0,
//                 "bankruptComm": 0,
//                 "bankruptPriceEp": 0,
//                 "bankruptPrice": 0,
//                 "positionMarginEv": 0,
//                 "positionMargin": 0,
//                 "liquidationPriceEp": 0,
//                 "liquidationPrice": 0,
//                 "deleveragePercentileEr": 0,
//                 "deleveragePercentile": 0,
//                 "buyValueToCostEr": 1150750,
//                 "buyValueToCost": 0.0115075,
//                 "sellValueToCostEr": 1149250,
//                 "sellValueToCost": 0.0114925,
//                 "markPriceEp": 93169002,
//                 "markPrice": 9316.9002,
//                 "markValueEv": 0,
//                 "markValue": null,
//                 "estimatedOrdLossEv": 0,
//                 "estimatedOrdLoss": 0,
//                 "usedBalanceEv": 0,
//                 "usedBalance": 0,
//                 "takeProfitEp": 0,
//                 "takeProfit": null,
//                 "stopLossEp": 0,
//                 "stopLoss": null,
//                 "realisedPnlEv": 0,
//                 "realisedPnl": null,
//                 "cumRealisedPnlEv": 0,
//                 "cumRealisedPnl": null
//             }
//             ]
//         }
// }

// https://api.phemex.com/accounts/positions?currency=<currency>
// GET /accounts/positions?currency=<currency>
// {
// 	"code": 0,
// 	"msg": "",
// 	"data": {
// 	  "account": {
// 		"accountId": 111100001,
// 		"currency": "BTC",
// 		"accountBalanceEv": 879599942377,
// 		"totalUsedBalanceEv": 285,
// 		"bonusBalanceEv": 0
// 	  },
// 	  "positions": [
// 		{
// 		  "accountID": 111100001,
// 		  "symbol": "BTCUSD",
// 		  "currency": "BTC",
// 		  "side": "Buy",
// 		  "positionStatus": "Normal",
// 		  "crossMargin": false,
// 		  "leverageEr": 0,
// 		  "initMarginReqEr": 1000000,
// 		  "maintMarginReqEr": 500000,
// 		  "riskLimitEv": 10000000000,
// 		  "size": 5,
// 		  "valueEv": 26435,
// 		  "avgEntryPriceEp": 189143181,
// 		  "posCostEv": 285,
// 		  "assignedPosBalanceEv": 285,
// 		  "bankruptCommEv": 750000,
// 		  "bankruptPriceEp": 5000,
// 		  "positionMarginEv": 879599192377,
// 		  "liquidationPriceEp": 5000,
// 		  "deleveragePercentileEr": 0,
// 		  "buyValueToCostEr": 1150750,
// 		  "sellValueToCostEr": 1149250,
// 		  "markPriceEp": 238287555,
// 		  "markValueEv": 0,
// 		  "unRealisedPosLossEv": 0,
// 		  "estimatedOrdLossEv": 0,
// 		  "usedBalanceEv": 285,
// 		  "takeProfitEp": 0,
// 		  "stopLossEp": 0,
// 		  "cumClosedPnlEv": -8913353,
// 		  "cumFundingFeeEv": 123996,
// 		  "cumTransactFeeEv": 940245,
// 		  "realisedPnlEv": 0,
// 		  "unRealisedPnlEv": 5452,
// 		  "cumRealisedPnlEv": 0
// 		}
// 	  ]
// 	}
//   }

// https://api.phemex.com/positions/leverage?symbol=<symbol>&leverage=<leverage>&leverageEr=<leverageEr>
// PUT /positions/leverage?symbol=<symbol>&leverage=<leverage>&leverageEr=<leverageEr>
// {
//     "code": 0,
//     "msg": "OK"
// }

// https://api.phemex.com/positions/riskLimit?symbol=<symbol>&riskLimit=<riskLimit>&riskLimitEv=<riskLimitEv>
// PUT /positions/riskLimit?symbol=<symbol>&riskLimit=<riskLimit>&riskLimitEv=<riskLimitEv>

// Query open orders by symbol
// https://api.phemex.com/orders/activeList?symbol=<symbol>
// GET /orders/activeList?symbol=<symbol>
// {
//     "code": 0,
//         "msg": "",
//         "data": {
//             "rows": [
//             {
//                 "bizError": 0,
//                 "orderID": "9cb95282-7840-42d6-9768-ab8901385a67",
//                 "clOrdID": "7eaa9987-928c-652e-cc6a-82fc35641706",
//                 "symbol": "BTCUSD",
//                 "side": "Buy",
//                 "actionTimeNs": 1580533011677666800,
//                 "transactTimeNs": 1580533011677666800,
//                 "orderType": null,
//                 "priceEp": 84000000,
//                 "price": 8400,
//                 "orderQty": 1,
//                 "displayQty": 1,
//                 "timeInForce": null,
//                 "reduceOnly": false,
//                 "stopPxEp": 0,
//                 "closedPnlEv": 0,
//                 "closedPnl": 0,
//                 "closedSize": 0,
//                 "cumQty": 0,
//                 "cumValueEv": 0,
//                 "cumValue": 0,
//                 "leavesQty": 0,
//                 "leavesValueEv": 0,
//                 "leavesValue": 0,
//                 "stopPx": 0,
//                 "stopDirection": "Falling",
//                 "ordStatus": "Untriggered"
//             },
//             {
//                 "bizError": 0,
//                 "orderID": "93397a06-e76d-4e3b-babc-dff2696786aa",
//                 "clOrdID": "71c2ab5d-eb6f-0d5c-a7c4-50fd5d40cc50",
//                 "symbol": "BTCUSD",
//                 "side": "Sell",
//                 "actionTimeNs": 1580532983785506600,
//                 "transactTimeNs": 1580532983786370300,
//                 "orderType": null,
//                 "priceEp": 99040000,
//                 "price": 9904,
//                 "orderQty": 1,
//                 "displayQty": 1,
//                 "timeInForce": null,
//                 "reduceOnly": false,
//                 "stopPxEp": 0,
//                 "closedPnlEv": 0,
//                 "closedPnl": 0,
//                 "closedSize": 0,
//                 "cumQty": 0,
//                 "cumValueEv": 0,
//                 "cumValue": 0,
//                 "leavesQty": 1,
//                 "leavesValueEv": 10096,
//                 "leavesValue": 0.00010096,
//                 "stopPx": 0,
//                 "stopDirection": "UNSPECIFIED",
//                 "ordStatus": "New"
//             },
//             {
//                 "bizError": 0,
//                 "orderID": "2585817b-85df-4dea-8507-5db1920b9954",
//                 "clOrdID": "4b19fd1e-a1a7-2986-d02a-0288ad5137d4",
//                 "symbol": "BTCUSD",
//                 "side": "Buy",
//                 "actionTimeNs": 1580532966629408500,
//                 "transactTimeNs": 1580532966633276200,
//                 "orderType": null,
//                 "priceEp": 80040000,
//                 "price": 8004,
//                 "orderQty": 1,
//                 "displayQty": 1,
//                 "timeInForce": null,
//                 "reduceOnly": false,
//                 "stopPxEp": 0,
//                 "closedPnlEv": 0,
//                 "closedPnl": 0,
//                 "closedSize": 0,
//                 "cumQty": 0,
//                 "cumValueEv": 0,
//                 "cumValue": 0,
//                 "leavesQty": 1,
//                 "leavesValueEv": 12493,
//                 "leavesValue": 0.00012493,
//                 "stopPx": 0,
//                 "stopDirection": "UNSPECIFIED",
//                 "ordStatus": "New"
//             }
//             ]
//         }
// }

// Query user order by orderID or Query user order by client order ID
// https://api.phemex.com/exchange/order? symbol=<symbol> & orderID=<orderID1, orderID2>
// GET /exchange/order? symbol=<symbol> & orderID=<orderID1, orderID2>

// https://api.phemex.com/exchange/order? symbol=<symbol> & clOrdID=<clOrdID1, clOrdID2>
// GET /exchange/order? symbol=<symbol> & clOrdID=<clOrdID1, clOrdID2>
// {
//     "code": 0,
//         "msg": "OK",
//         "data": [
//         {
//             "orderID": "7d5a39d6-ff14-4428-b9e1-1fcf1800d6ac",
//             "clOrdID": "e422be37-074c-403d-aac8-ad94827f60c1",
//             "symbol": "BTCUSD",
//             "side": "Sell",
//             "orderType": "Limit",
//             "actionTimeNs": 1577523473419470300,
//             "priceEp": 75720000,
//             "price": null,
//             "orderQty": 12,
//             "displayQty": 0,
//             "timeInForce": "GoodTillCancel",
//             "reduceOnly": false,
//             "takeProfitEp": 0,
//             "takeProfit": null,
//             "stopLossEp": 0,
//             "closedPnlEv": 0,
//             "closedPnl": null,
//             "closedSize": 0,
//             "cumQty": 0,
//             "cumValueEv": 0,
//             "cumValue": null,
//             "leavesQty": 0,
//             "leavesValueEv": 0,
//             "leavesValue": null,
//             "stopLoss": null,
//             "stopDirection": "UNSPECIFIED",
//             "ordStatus": "Canceled",
//             "transactTimeNs": 1577523473425416400
//         },
//         {
//             "orderID": "b63bc982-be3a-45e0-8974-43d6375fb626",
//             "clOrdID": "uuid-1577463487504",
//             "symbol": "BTCUSD",
//             "side": "Sell",
//             "orderType": "Limit",
//             "actionTimeNs": 1577963507348468200,
//             "priceEp": 71500000,
//             "price": null,
//             "orderQty": 700,
//             "displayQty": 700,
//             "timeInForce": "GoodTillCancel",
//             "reduceOnly": false,
//             "takeProfitEp": 0,
//             "takeProfit": null,
//             "stopLossEp": 0,
//             "closedPnlEv": 0,
//             "closedPnl": null,
//             "closedSize": 0,
//             "cumQty": 700,
//             "cumValueEv": 9790209,
//             "cumValue": null,
//             "leavesQty": 0,
//             "leavesValueEv": 0,
//             "leavesValue": null,
//             "stopLoss": null,
//             "stopDirection": "UNSPECIFIED",
//             "ordStatus": "Filled",
//             "transactTimeNs": 1578026629824704800
//         }
//     ]
// }

// WSS
// PUBLIC:
//	-> trade 		-> sub to
//	-> orderbook

// PRIVATE:
//	-> account
//	-> postion
//	-> order

// {
// 	"method": "user.auth",
// 	"params": [
// 	  "API",
// 	  "<token>",
// 	  "<signature>",
// 	  <expiry>
// 	],
// 	"id": 1234
//   }

// Field		Type		Description	Possible values
// type			String		Token type	API
// token		String		API Key
// signature	String		Signature generated by a funtion as HMacSha256(API Key + expiry) with API Secret
// expiry		Integer		A future time after which request will be rejected, in epoch second.
//					 		Maximum expiry is request time plus 2 minutes

// sample:
// > {
// 	"method": "user.auth",
// 	"params": [
// 	  "API",
// 	  "806066b0-f02b-4d3e-b444-76ec718e1023",
// 	  "8c939f7a6e6716ab7c4240384e07c81840dacd371cdcf5051bb6b7084897470e",
// 	  1570091232
// 	],
// 	"id": 1234
//   }

//   < {
// 	"error": null,
// 	"id": 1234,
// 	"result": {
// 	  "status": "success"
// 	}
//   }

// Request
// {
//   "id": <id>,
//   "method": "trade.subscribe",
//   "params": [
//     "<symbol>"
//   ]
// }
// Response:
// {
//   "error": null,
//   "id": <id>,
//   "result": {
//     "status": "success"
//   }
// }
