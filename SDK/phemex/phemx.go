package phemex

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	_ "crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
)

var client, paths = setupClient()

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

type Paths []string

type Account struct {
	Type     string `json:"TYPE"` // ["main", "sub"]
	ID       string `json:"ID"`
	hmac     hash.Hash
	Socket   *websocket.Conn
	Accounts map[string]*Account
}

type Client struct {
	conn       http.Client
	HostHTTP   string `json:"HOSTHTTP"`
	HostWSS    string `json:"HOSTWSS"`
	SocketConn int
	Account    *Account
}

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

type Response struct {
	data   []byte                 // rawData, blob
	output map[string]interface{} // data
	req    *Request
}

type Position struct {
	AccountID              int64   // int64
	Symbol                 string  // string
	Currency               string  // string
	Side                   string  // string
	PositionStatus         string  // string
	CrossMargin            bool    // bool
	LeverageEr             int64   // int64
	Leverage               float64 // int64
	InitMarginReqEr        int64   // int64
	InitMarginReq          float64 // float64
	MaintMarginReqEr       int64   // int64
	MaintMarginReq         float64 // float64
	RiskLimitEv            int64   // int64
	RiskLimit              float64 // int64
	Size                   int64   // int64
	Value                  float64 // int64
	ValueEv                int64   // int64
	AvgEntryPriceEp        int64   // int64
	AvgEntryPrice          float64 // int64
	PosCostEv              int64   // int64
	PosCost                float64 // int64
	AssignedPosBalanceEv   int64   // int64
	AssignedPosBalance     float64 // int64
	BankruptCommEv         int64   // int64
	BankruptComm           float64 // int64
	BankruptPriceEp        int64   // int64
	BankruptPrice          float64 // int64
	PositionMarginEv       int64   // int64
	PositionMargin         float64 // int64
	LiquidationPriceEp     int64   // int64
	LiquidationPrice       float64 // int64
	DeleveragePercentileEr int64   // int64
	DeleveragePercentile   int64   // int64
	BuyValueToCostEr       int64   // int64
	BuyValueToCost         float64 // float64
	SellValueToCostEr      int64   // int64
	SellValueToCost        float64 // float64
	MarkPriceEp            int64   // int64
	MarkPrice              float64 // float64
	MarkValueEv            int64   // int64
	MarkValue              float64
	EstimatedOrdLossEv     int64 // int64
	EstimatedOrdLoss       int64 // int64
	UsedBalanceEv          int64 // int64
	UsedBalance            int64 // int64
	TakeProfitEp           int64 // int64
	TakeProfit             float64
	StopLossEp             int64 // int64
	StopLoss               float64
	RealisedPnlEv          int64 // int64
	RealisedPnl            float64
	CumRealisedPnlEv       int64 // int64
	CumRealisedPnl         float64
}

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

func (r *Request) sign() {

	if r.isPrivate() {
		minute := 60
		time := int(time.Now().Unix())

		r.Expiry = strconv.Itoa(time + minute)

		byteMessage := []byte(r.Path + r.Query + r.Expiry + string(r.Body))

		client.Account.hmac.Write(byteMessage)
		r.Signature = fmt.Sprintf("%x", client.Account.hmac.Sum(nil))

		client.Account.hmac.Reset()

		r.Req.Header.Add("x-phemex-access-token", client.Account.ID)
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

func (r *Response) Display() {
	if r.output == nil {
		r.HandleResponse(JSON)
	}

	data, err := json.MarshalIndent(r.output, "", "  ")
	if err != nil {
		panic("yike")
	}
	fmt.Println(string(data))
}

func (r *Response) HandleResponse(handler func(res *Response)) *Response {
	handler(r)
	return r
}

func JSON(res *Response) {
	var data map[string]interface{}
	err := json.Unmarshal(res.data, &data)
	if err != nil {
		panic("NOO")
	}

	res.output = data
}

func Send(method, path string, query map[string]string, body map[string]interface{}) *Response {
	request := new(Request)
	response := new(Response)
	request.setPath(method, path)
	request.setQuery(query)
	request.setBody(body)
	request.setRequest()
	request.sign()
	request.send(response)

	return response
}

func setupClient() (*Client, *Paths) {
	// setup client
	// set up paths
	// set up websockets

	paths := new(Paths)
	*paths = append(*paths, "/orders") // POST 	 -> Body {symbol, clOrdID, side, priceEp, ordrQty, actionBy, pegPriceType, pegOffsetValueEp, pegPriceType
	// 				  									  , reduceOnly, timeInforce, takeProfitEp, StopLossEp, stopPxEp, closeOnTrigger, triggertype}
	*paths = append(*paths, "/orders/replace")    // PUT 	 -> query
	*paths = append(*paths, "/orders")            // DELETE -> query {symbol, orderID=[]}
	*paths = append(*paths, "/orders/cancel")     // DELETE -> query {symbol, orderID}
	*paths = append(*paths, "/orders/all")        // DELETE -> query {symbol, untriggered, text}
	*paths = append(*paths, "/orders/activeList") // GET 	 -> query {symbol, ordStatus}

	*paths = append(*paths, "/positions/leverage")  // PUT    -> query {symbol, leverage, leverageEr}
	*paths = append(*paths, "/positions/riskLimit") // PUT    -> query {symbol, riskLimit, riskLimitEv}
	*paths = append(*paths, "/positions/assign")    // POST   -> query {symbol, posBalance, posBalanceEv}

	*paths = append(*paths, "/accounts/accountPositions")  // GET    -> query {currency}
	*paths = append(*paths, "/accounts/positions")         // GET    -> query {currency}
	*paths = append(*paths, "/phemex-user/users/children") // GET?	-> query {offset, limit, withCount}

	*paths = append(*paths, "/md/orderbook")   // GET	-> query {symbol}
	*paths = append(*paths, "/md/trade")       // GET	-> query {symbol}
	*paths = append(*paths, "/md/ticker/24hr") // GET	-> query {symbol}

	*paths = append(*paths, "/exchange/order")       // GET	-> query {symbol, orderID=[]}
	*paths = append(*paths, "/exchange/order")       // GET	-> query {symbol, clOrd=[]}
	*paths = append(*paths, "/exchange/order/list")  // GET	-> query {symbol, start, end, offset, limit, ordStatus, withcount}
	*paths = append(*paths, "/exchange/order/trade") // GET	-> query {symbol, tradeType, start, end, offset, limit, withcount}

	*paths = append(*paths, "/exchange/margins")          // POST	-> Body {btcAmount, btcAmountEv, linkKey, moveOp, usdAmount, usdAmountEv}
	*paths = append(*paths, "/exchange/margins/transfer") // GET	-> query {start, end, offset, limit, withCount}

	*paths = append(*paths, "/exchange/wallets/transferOut")           // POST	-> Body {amount, amountEv, clientCnt, currency}
	*paths = append(*paths, "/exchange/wallets/transferIn")            // POST	-> Body {amount, amountEv, clientCnt, currency}
	*paths = append(*paths, "/exchange/wallets/createWithdraw")        // POST	-> query {optCode} -> Body {address, amountEv, currency}
	*paths = append(*paths, "/exchange/wallets/confirm/withdraw")      // GET	-> query {code}
	*paths = append(*paths, "/exchange/wallets/cancelWithdraw")        // POST	-> Body {id}
	*paths = append(*paths, "/exchange/wallets/withdrawList")          // GET	-> query {currency, limit, offset, withCount}
	*paths = append(*paths, "/exchange/wallets/createWithdrawAddress") // POST	-> query {optCode} -> Body {address, currency, remark}

	*paths = append(*paths, "/exchange/public/nomics/trades") // GET	-> query {market, since}
	*paths = append(*paths, "/exchange/public/products")      // GET

	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	data := make(map[string]interface{})

	json.Unmarshal(byteValue, &data)

	client := new(Client)
	client.conn = *http.DefaultClient
	client.HostHTTP = data["HOSTHTTP"].(string)
	client.HostWSS = data["HOSTWSS"].(string)

	var addr = flag.String("addr", "phemex.com", "phemex feed address")
	flag.Parse()
	log.SetFlags(0)
	u := url.URL{Scheme: "wss", Host: *addr, Path: "ws"}

	accounts := data["CLIENTS"].([]interface{})
	for i := 0; i < len(accounts); i++ {
		item := accounts[i].(map[string]interface{})
		account := new(Account)
		account.ID = item["ID"].(string)
		account.Type = item["TYPE"].(string)
		account.hmac = hmac.New(crypto.SHA256.New, []byte(item["SECRET"].(string)))

		if client.Account == nil {
			client.Account = account
			client.Account.Accounts = map[string]*Account{}
		}

		account.Accounts = client.Account.Accounts
		account.Accounts[account.ID] = account
		account.Connect(u)

		if account.Type == "MAIN" {
			client.Account = account
		}

	}

	return client, paths
}

func MainSub() {
	client.Account.Auth().Subscribe("aop.subscribe", []interface{}{})
}

func (a *Account) Auth() *Account {

	seconds := 120
	time := int(time.Now().Unix())

	expiry := time + seconds
	byteMessage := []byte(fmt.Sprintf("%v%d", a.ID, expiry))

	a.hmac.Write(byteMessage)
	signature := fmt.Sprintf("%x", a.hmac.Sum(nil))

	a.hmac.Reset()

	message, err := json.Marshal(map[string]interface{}{
		"method": "user.auth",
		"params": []interface{}{
			"API",
			a.ID,
			signature,
			expiry},
		"id": 1234})

	if err != nil {
		panic("yike")
	}

	done := make(chan struct{})
	func() {
		defer close(done)
		err := a.Socket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
		}
		_, message, err := a.Socket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}()

	return a
}

func Subscribe() {
	client.Account.Subscribe("orderbook.subscribe", []interface{}{"BTCUSD"})
}

func (a *Account) Subscribe(method string, params []interface{}) {

	defer a.Socket.Close()

	message, err := json.Marshal(map[string]interface{}{
		"id":     1234,
		"method": method,
		"params": params})

	if err != nil {
		panic("yike")
	}

	done := make(chan struct{})

	func() {
		defer close(done)
		err := a.Socket.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			log.Println("write:", err)
			return
		}

		for {
			_, message, err := a.Socket.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
}

func (a *Account) Connect(u url.URL) {

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	a.Socket = c
}
