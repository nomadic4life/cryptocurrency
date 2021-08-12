package phemex

import (
	_ "crypto/sha256"
	"fmt"
	"time"
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
