package phemex

import (
	"crypto"
	"crypto/hmac"
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
	receiver      chan []byte
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

	accounts := data["CLIENTS"].([]interface{})
	for i := 0; i < len(accounts); i++ {
		item := accounts[i].(map[string]interface{})
		account := new(Account)
		account.ID = int64(item["ID"].(float64))
		account.API_KEY = item["API_KEY"].(string)
		account.Type = item["TYPE"].(string)
		account.hmac = hmac.New(crypto.SHA256.New, []byte(item["SECRET"].(string)))
		account.Subscriptions = make(map[string]int)

		if client.Account == nil {
			client.Account = account
			client.Account.Accounts = map[int64]*Account{}
		}

		account.Accounts = client.Account.Accounts
		account.Accounts[account.ID] = account
		account.Client = client

		if account.Type == "MAIN" {
			// considering that anything has to do with public info is in a public
			// account and that the public account is at the top of the tree?
			client.Account = account
			// main listener. not sure if I should start here. but I am.
			account.Listener()
		}
	}

	return client, paths
}

func (a *Account) Listener() {
	if a.receiver != nil {
		return
	}

	a.receiver = make(chan []byte, 100)

	go func() {
		for {
			message := <-a.receiver
			fmt.Println("\n from account listener: \n", a.ID, string(message))
			// need a message handler
			// handler can be run in a go routine
			// many will routines will be spend up to
			// handle many incoming messages
			// to prevent any blocking
			// implement a counter to maintain order
			// of an incoming message.
			// need a channel to kill listener and break from loop
			// and to close all account channels.
		}
	}()
}

func (Conn *Client) activeAccount(message []byte) *Account {
	data := parseMessage(message)
	for {
		switch v := data.(type) {
		case map[string]interface{}:
			if val, ok := v["position_info"]; ok {
				data = val
			} else if val, ok := v["accounts"]; ok {
				data = val
			} else if val, ok := v["userID"]; ok {
				data = val
			}

		case []interface{}:
			data = v[0]

		case string:
			userID, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("error")
				return nil
			}
			return Conn.Account.Accounts[int64(userID)]

		case float64:
			userID := v
			return Conn.Account.Accounts[int64(userID)]
		}
	}
}

func result(message []byte) string {
	data := parseMessage(message)

	for {
		switch v := data.(type) {

		case map[string]interface{}:
			if val, ok := v["result"]; ok {
				data = val
			} else if val, ok := v["status"]; ok {
				data = val
			}

		case string:
			return v
		default:
			return "failure"
		}
	}
}

// not relevent function yet.
func Message(msg []byte) {
	if strings.Contains(string(msg), "{\"error\"") {

	} else if strings.Contains(string(msg), "{\"book\"") {

	} else if strings.Contains(string(msg), "{\"trades\"") {

	} else if strings.Contains(string(msg), "{\"kline\"") {

	} else if strings.Contains(string(msg), "{\"accounts\"") {

	} else if strings.Contains(string(msg), "{\"market24h\"") {

	} else if strings.Contains(string(msg), "{\"id\"") {

	} else if strings.Contains(string(msg), "{\"result\"") {

	} else if strings.Contains(string(msg), "{\"status\"") {

	}
}
