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

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = 5 * time.Second
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = true
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Paths []string

type Account struct {
	Client        *Client
	Type          string
	ID            int64
	API_KEY       string
	hmac          hash.Hash
	Subscriptions map[string]int
	Accounts      map[int64]*Account
}

type Client struct {
	conn     http.Client
	HostHTTP string `json:"HOSTHTTP"`
	HostWSS  string `json:"HOSTWSS"`
	WSS      url.URL
	ConnMap  map[int]int
	Sockets  [5]*Socket
	Account  *Account
}

type Socket struct {
	Hub  *Hub
	Conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	// Registered Accounts.
	Accounts map[*Account]bool

	// Inbound messages from the Sockets.
	response chan []byte

	// Register requests from the Sockets or Accounts?.
	register chan *Account

	// Unregister requests from Sockets or Accounts?.
	unregister chan *Account
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

func (a *Account) Send(method, path string, query map[string]string, body map[string]interface{}) *Response {
	request := new(Request)
	response := new(Response)
	request.setPath(method, path)
	request.setQuery(query)
	request.setBody(body)
	request.setRequest()
	request.sign(a)
	request.send(response)

	return response
}

func setupPaths() *Paths {
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

	return paths
}

func readConfig() []byte {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	return byteValue
}

func setupClient() (*Client, *Paths) {
	paths := setupPaths()

	data := make(map[string]interface{})

	json.Unmarshal(readConfig(), &data)

	addr := flag.String("addr", "phemex.com", "phemex feed address")
	flag.Parse()
	log.SetFlags(0)

	client := new(Client)
	client.conn = *http.DefaultClient
	client.HostHTTP = data["HOSTHTTP"].(string)
	client.HostWSS = data["HOSTWSS"].(string)
	client.WSS = url.URL{Scheme: "wss", Host: *addr, Path: "ws"}

	client.ConnMap = make(map[int]int)
	client.ConnMap[0] = 0
	client.ConnMap[1] = 0
	client.ConnMap[2] = 0
	client.ConnMap[3] = 0
	client.ConnMap[4] = 0

	// client.Hub = new(Hub)
	// client.Hub.broadcast = make(chan []byte)
	// client.Hub.register = make(chan *Account)
	// client.Hub.unregister = make(chan *Account)
	// client.Hub.Accounts = make(map[*Account]bool)

	accounts := data["CLIENTS"].([]interface{})
	for i := 0; i < len(accounts); i++ {
		item := accounts[i].(map[string]interface{})
		account := new(Account)
		account.ID = int64(item["ID"].(float64))
		account.API_KEY = item["API_KEY"].(string)
		account.Type = item["TYPE"].(string)
		account.hmac = hmac.New(crypto.SHA256.New, []byte(item["SECRET"].(string)))

		if client.Account == nil {
			client.Account = account
			client.Account.Accounts = map[int64]*Account{}
		}

		account.Accounts = client.Account.Accounts
		account.Accounts[account.ID] = account
		account.Client = client

		if account.Type == "MAIN" {
			client.Account = account
		}
	}

	fmt.Print(client.Account.Accounts)

	return client, paths
}

func GetAccounts() {
	client.GetAccounts()
}

func (c *Client) GetAccounts() {

	res := c.Account.Send("GET", "/phemex-user/users/children", nil, nil).HandleResponse(JSON)
	users := res.output["data"].([]interface{})
	for i := 0; i < len(users); i++ {
		user := users[i]
		fmt.Println(int64(user.(map[string]interface{})["userId"].(float64)))
	}
}

// func MainSub() {
// 	// client.Account.Auth().Subscribe("aop.subscribe", []interface{}{})
// 	fmt.Println("hello")
// }

// func Subscribe() {
// 	client.Account.Subscribe("orderbook.subscribe", []interface{}{"BTCUSD"})
// 	// client.Account.Subscribe("trade.subscribe", []interface{}{"BTCUSD"})
// }

// // plays pingpong to keep socket connection alive
// func keepAlive(c *websocket.Conn, id int64, done chan struct{}) {
// 	rand.Seed(time.Now().UnixNano())
// 	ticker := time.NewTicker(WebsocketTimeout)

// 	lastResponse := time.Now()
// 	c.SetPongHandler(func(msg string) error {
// 		lastResponse = time.Now()
// 		return nil
// 	})

// 	go func() {
// 		defer func() {
// 			ticker.Stop()
// 		}()
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case <-ticker.C:
// 				id++
// 				p, _ := json.Marshal(map[string]interface{}{
// 					"id":     id,
// 					"method": "server.ping",
// 					"params": []string{},
// 				})
// 				c.WriteControl(websocket.PingMessage, p, time.Time{})

// 				// as suggested https://github.com/phemex/phemex-api-docs/blob/master/Public-API-en.md#1-session-management
// 				if time.Since(lastResponse) > 3*WebsocketTimeout {
// 					// errHandler(fmt.Errorf("last pong exceeded the timeout: %[1]v (%[2]v)", time.Since(lastResponse), id))
// 					fmt.Printf("last pong exceeded the timeout: %[1]v (%[2]v)", time.Since(lastResponse), id)
// 					return
// 				}
// 			}
// 		}
// 	}()
// }

// func (h *Hub) run() {

// 	go func() {
// 		for {
// 			select {
// 			case account := <-h.register:
// 				h.Accounts[account] = true
// 			case account := <-h.unregister:
// 				if _, ok := h.Accounts[account]; ok {
// 					delete(h.Accounts, account)
// 					close(s.send)
// 				}
// 			// comes from socket connection
// 			case message := <-h.response:
// 				// send to main or one of the sub accounts
// 				// extract the id or account id

// 				// send to the account so the account can handle data
// 				select {
// 				case s.send <- message:
// 				default:
// 					close(s.send)
// 					delete(h.Accounts, account)
// 				}

// 			}
// 		}
// 	}()
// }

// func Message(msg []byte) {
// 	if strings.Contains(string(msg), "{\"error\"") {

// 	} else if strings.Contains(string(msg), "{\"book\"") {

// 	} else if strings.Contains(string(msg), "{\"trades\"") {

// 	} else if strings.Contains(string(msg), "{\"kline\"") {

// 	} else if strings.Contains(string(msg), "{\"accounts\"") {

// 	} else if strings.Contains(string(msg), "{\"market24h\"") {

// 	} else if strings.Contains(string(msg), "{\"id\"") {

// 	} else if strings.Contains(string(msg), "{\"result\"") {

// 	} else if strings.Contains(string(msg), "{\"status\"") {

// 	}

// }

// // Auth -> allows account to be authorized before subsribing to a channel
// func (a *Account) Auth() *Account {

// 	_, socket := a.Client.Subscribe()
// 	if socket == nil {
// 		fmt.Println("Max Connections, no more connections can be made.")
// 		return nil
// 	}

// 	seconds := 120
// 	time := int(time.Now().Unix())

// 	expiry := time + seconds
// 	byteMessage := []byte(fmt.Sprintf("%v%d", a.ID, expiry))

// 	a.hmac.Write(byteMessage)
// 	signature := fmt.Sprintf("%x", a.hmac.Sum(nil))

// 	a.hmac.Reset()

// 	message, err := json.Marshal(map[string]interface{}{
// 		"method": "user.auth",
// 		"params": []interface{}{
// 			"API",
// 			a.ID,
// 			signature,
// 			expiry},
// 		"id": 1234})

// 	if err != nil {
// 		panic("yike")
// 	}

// 	func() {
// 		err := socket.WriteMessage(websocket.TextMessage, message)
// 		if err != nil {
// 			log.Println("write:", err)
// 		}
// 	}()

// 	return a
// }

// // Subscribe -> allows account to make a new subscribtion channel if not at max capacity
// func (a *Account) Subscribe(method string, params []interface{}) *Account {

// 	index, socket := a.Client.Subscribe()
// 	if socket == nil {
// 		fmt.Println("Max Connections, no more connections can be made.")
// 		return a
// 	}

// 	message, err := json.Marshal(map[string]interface{}{
// 		"id":     1234,
// 		"method": method,
// 		"params": params})

// 	if err != nil {
// 		panic("yike")
// 	}

// 	func() {
// 		err := socket.WriteMessage(websocket.TextMessage, message)
// 		// if succesful connection update connection
// 		client.ConnMap[index] += 1

// 		if err != nil {
// 			log.Println("write:", err)
// 		}
// 	}()

// 	return a
// }

// // Subscribe -> returns an available socket connection so Account can make a subscription channel or returns none if  at max capacity.
// func (Conn *Client) Subscribe() (int, *websocket.Conn) {
// 	// fmt.Println(Conn)
// 	// fmt.Println(Conn.ConnMap)
// 	for key, value := range Conn.ConnMap {
// 		if value < 20 && value > 0 {
// 			fmt.Println(key, Conn.Sockets[key])
// 			return key, Conn.Sockets[key]
// 		}
// 	}

// 	return Conn.Connect()
// }

// // Connect -> makes a new socket connection or none if at max capacity
// func (Conn *Client) Connect() (int, *websocket.Conn) {

// 	for i := 0; i < len(Conn.Sockets); i++ {
// 		if Conn.Sockets[i] == nil {

// 			done := make(chan struct{})
// 			socket, _, err := websocket.DefaultDialer.Dial(Conn.WSS.String(), nil)
// 			if err != nil {
// 				log.Fatal("dial:", err)
// 				return -1, nil
// 			}

// 			if WebsocketKeepalive {
// 				// keepAlive(a.Socket, *s.id, done, errHandler)
// 				keepAlive(socket, 1, done)
// 			}

// 			go func() {

// 				defer func() {
// 					// Conn.Hub.unregister <-
// 					socket.Close() // not sure if this goes here? don't fully completly understand this concept
// 					close(done)
// 				}()

// 				for {
// 					_, message, err := socket.ReadMessage()
// 					if err != nil {
// 						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 							log.Printf("error: %v", err)
// 						} else {
// 							log.Println("read:", err)
// 						}
// 						return
// 					}

// 					// log.Printf("recv: %s", message)
// 					// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
// 					// parse message here
// 					// send message to correct account by account id or to main by default

// 					if strings.Contains(string(message), "{\"accounts\"") {
// 						// extract id
// 						// find account from id
// 						// send message to acount
// 					} else if strings.Contains(string(message), "{\"results\"") {
// 						// if success extract id
// 						// update socket counter
// 						// update account reference to socket and subscribtion to refer for unsubing
// 					} else {
// 						// send to main
// 					}
// 				}
// 			}()

// 			Conn.Sockets[i] = socket

// 			return i, socket
// 		}
// 	}

// 	return -1, nil
// }

// websocket to the hub
// hub to the account
