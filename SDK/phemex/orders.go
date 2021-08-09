package phemex

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var DefaultOrder map[string]interface{} = map[string]interface{}{
	// DEFAULT KEY VALUES:
	"symbol":         "BTCUSD",
	"ordtype":        "Limit",
	"timeInForce":    "PostOnly",
	"reduceOnly":     false,
	"closeOnTrigger": false}

func CreateOrder(side string, orderQty, priceEp int64, append ...map[string]interface{}) *Response {

	// DEFAULT KEY VALUES:
	//	-> "symbol": 			"BTCUSD"
	//	-> "clOrdID": 			<UUID> - Generated
	//	-> "ordtype": 			"Limit"
	//	-> "timeInForce": 		"PostOnly"
	//	-> "reduceOnly": 		false
	//	-> "closeOnTrigger": 	false

	// GENERATED KEY VALUES:
	// 	-> "clOrdID": 			<UUID> - Generated

	// INPUT KEY VALUES:
	//	-> "side": 				["Buy", Sell]
	//	-> "orderQty": 			<QUANTITY>
	//	-> "priceEp": 			<PRICE>

	order := map[string]interface{}{
		"side":     side,
		"orderQty": orderQty,
		"priceEp":  priceEp,
		"clOrdID":  uuid.NewString()}

	for key, value := range DefaultOrder {
		order[key] = value
	}

	if append != nil {
		for key, value := range append[0] {
			if key != "side" && key != "orderQty" && key != "priceEp" {
				fmt.Print(key)
				order[key] = value
			}
		}
	}

	return Send("POST", "/orders", nil, order)
}

func AmendOrder(id string, amend map[string]string) *Response {

	query := map[string]string{
		"symbol":  "BTCUSD",
		"orderID": id}

	for key, value := range amend {
		if key != "orderID" {
			query[key] = value
		}
	}

	return Send("PUT", "/orders/replace", query, nil)
}

func CancelOrders(ids []string, symbol string) *Response {

	query := map[string]string{
		"symbol": "BTCUSD"}

	path := "/orders/cancel"

	if symbol != "" {
		query["symbol"] = symbol
	}

	if ids == nil {
		path = "/orders/all"
	} else if len(ids) > 1 {
		path = "/orders"
		query["orderID"] = strings.Join(ids, ",")
		// query["untriggered"] = "false"
		// query["untriggered"] = "true"
	} else {
		query["orderID"] = ids[0]
	}

	return Send("DELETE", path, query, nil)
}

func GetOrders(id, symbol string) *Response {

	path := "/orders/activeList"
	query := map[string]string{
		"symbol": "BTCUSD"}

	if symbol != "" {
		query["symbol"] = symbol
	}

	if id != "" {
		path = "/orders/active"
		query["orderID"] = id
	}

	return Send("GET", path, query, nil)
}

// need to add some validation and sanitation for query and body keys
// need to create enums for enum keys
// reconsider inputs for cancel order and get orders

// 	:: ENUMS :: 	-> order input <-
// side -> ["Buy", "Sell"]
// orderType -> ["Market", "Limit", "StopLimit", "MarketIfTouched", "LimitIfTouched"]
// timeInForce -> ["GoodtillCancel", "ImmediateOrCancel", "FillOrKill", "PostOnly"]
// triggerType -> ["ByMarkPrice", "ByLastPrice"]
// pegPriceType -> ["TrailingStopPeg", "TrailingTakeProfitPeg"]
// symbol -> ["BTCUSD",...] -> <Trading Symbols>

// 	:: Bool ::		-> order input <-
//	-> reduceOnly
//	-> closeOnTrigger

// 	:: String ::	-> order input <-
//	-> symbol
//	-> clOrdID

// 	:: Price ::		-> order input <-
//	-> priceEp
//	-> stopPxEP
//	-> takeProfitEp
//	-> stopLossEp
//	-> pegOffsetValueEp

// 	:: Required ::	-> order input <-
//	-> symbol
//	-> clOrdID
//	-> side
//	-> orderQty -> int64

//  :: Quantity ::	 -> order input <-
//	-> orderQty

//	:: Amend Order ::	-> order input <-
//	-> orderID
//	-> origClOrdID

//  :: Required Amend :: -> order input <-
//	-> symbol
//	-> orderID

// :: Cancel Order :: -> order input <-
//	-> symbol
//	-> orderID -> single -> []bulk
//	-> untriggered -> for bulk
//	-> text
