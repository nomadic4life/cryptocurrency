package phemex

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
)

type Socket struct {
	Client *Client // is this relevent?
	Conn   *websocket.Conn
	send   chan []byte // subscriber channel
}

func Subscribe() {
	// client.Account.Accounts[826079].Auth().Subscribe("aop.subscribe", []interface{}{}, true)
	// client.Account.Accounts[976380].Auth().Subscribe("aop.subscribe", []interface{}{}, true)
	// client.Account.Accounts[1929977].Auth().Subscribe("aop.subscribe", []interface{}{}, true)

	// client.Account.Accounts[826079].Subscribe("orderbook.subscribe", []interface{}{"BTCUSD"}, false)
	// client.Account.Accounts[976380].Subscribe("trade.subscribe", []interface{}{"BTCUSD"}, false)
	// client.Account.Accounts[1929977].Subscribe("market24h.subscribe", []interface{}{}, false)

	// client.Account.Accounts[826079].Subscribe("trade.subscribe", []interface{}{"BTCUSD"}, false)
	// client.Account.Accounts[826079].Subscribe("market24h.subscribe", []interface{}{}, false)

	// client.Account.Auth().Subscribe("orderbook.subscribe", []interface{}{"BTCUSD"})
	// client.Account.Subscribe("trade.subscribe", []interface{}{"BTCUSD"})
	// for key, value := range client.Account.Accounts {
	// 	fmt.Println(key)
	// 	value.Auth().Subscribe("aop.subscribe", []interface{}{}, true)
	// }
	client.Account.Auth().Subscribe("aop.subscribe", []interface{}{}, true)
}

// plays pingpong to keep socket connection alive
func keepAlive(c *websocket.Conn, id int, done chan struct{}) {
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(WebsocketTimeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer func() {
			ticker.Stop()
		}()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				id++
				p, _ := json.Marshal(map[string]interface{}{
					"id":     id,
					"method": "server.ping",
					"params": []string{},
				})
				c.WriteControl(websocket.PingMessage, p, time.Time{})

				// as suggested https://github.com/phemex/phemex-api-docs/blob/master/Public-API-en.md#1-session-management
				if time.Since(lastResponse) > 3*WebsocketTimeout {
					// errHandler(fmt.Errorf("last pong exceeded the timeout: %[1]v (%[2]v)", time.Since(lastResponse), id))
					fmt.Printf("last pong exceeded the timeout: %[1]v (%[2]v)", time.Since(lastResponse), id)
					return
				}
			}
		}
	}()
}

// Auth -> allows account to be authorized before subsribing to a channel
func (a *Account) Auth() *Account {
	fmt.Println("client", a)
	index, socket := a.Client.Subscribe()
	if socket == nil {
		fmt.Println("Max Connections, no more connections can be made.")
		return nil
	}

	seconds := 120
	time := int(time.Now().Unix())

	expiry := time + seconds
	byteMessage := []byte(fmt.Sprintf("%s%v", a.API_KEY, expiry))

	a.hmac.Write(byteMessage)
	signature := fmt.Sprintf("%x", a.hmac.Sum(nil))

	a.hmac.Reset()

	message, err := json.Marshal(map[string]interface{}{
		"method": "user.auth",
		"params": []interface{}{
			"API",
			a.API_KEY,
			signature,
			expiry},
		"id": a.ID})
	// no idea what to do with id prop

	if err != nil {
		panic("yike")
	}

	return func() *Account {
		err := socket.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
		}

		// blocking -> waiting to receive confirmation of Authorization
		message := <-socket.send

		// error -> panic: interface conversion: interface {} is nil, not map[string]interface {}
		if result(message) == "success" {
			client.ConnMap[index] += 1
			// instead of doing counter, should pass the index to the subscriber
			// and subscriber passes the index into the client subscribe and get
			// socket that way. by passing the check loop.
			// need to implement better error handeling.
			return a
		}

		return a
	}()
}

// Subscribe -> allows account to make a new subscribtion channel if not at max capacity
func (a *Account) Subscribe(method string, params []interface{}, auth bool) *Account {

	index, socket := a.Client.Subscribe()
	if socket == nil {
		fmt.Println("Max Connections, no more connections can be made.")
		return a
	}

	a.Listener()

	message, err := json.Marshal(map[string]interface{}{
		"id":     a.ID,
		"method": method,
		"params": params})

	if err != nil {
		panic("yike")
	}

	func() {
		err := socket.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
		}

		fmt.Println("there")

		// blocking -> waiting to receive confirmation of subscription
		message := <-socket.send

		if result(message) == "success" {
			fmt.Println("Made Sub: ->  ", string(message))
			a.Subscriptions[method] = index
			if auth != true {
				client.ConnMap[index] += 1
			}
		} else {
			fmt.Println("Unable to make Sub: ->  ", string(message))
		}
	}()

	return a
}

// Subscribe -> returns an available socket connection so Account can make a subscription channel or returns none if  at max capacity.
func (Conn *Client) Subscribe() (int, *Socket) {
	for key, value := range Conn.ConnMap {
		if value < 20 && value > 0 {
			return key, Conn.Sockets[key]
		}
	}

	return Conn.Connect()
}

// Connect -> makes a new socket connection or none if at max capacity
func (Conn *Client) Connect() (int, *Socket) {

	for i := 0; i < len(Conn.Sockets); i++ {
		if Conn.Sockets[i] == nil {

			done := make(chan struct{})
			connection, _, err := websocket.DefaultDialer.Dial(Conn.WSS.String(), nil)
			if err != nil {
				log.Fatal("dial:", err)
				return -1, nil
			}

			socket := new(Socket)
			socket.Conn = connection
			socket.Client = client
			socket.send = make(chan []byte)
			Conn.Sockets[i] = socket

			if WebsocketKeepalive {
				// keepAlive(a.Socket, *s.id, done, errHandler)
				keepAlive(socket.Conn, i, done)
			}

			go func() {

				defer func() {
					close(done)
					socket.Conn.Close()
				}()

				for {
					_, message, err := socket.Conn.ReadMessage()
					if err != nil {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
							log.Printf("error: %v", err)
						} else {
							log.Println("read:", err)
						}
						return
					}

					// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

					// need some error handling if message contains an error
					if strings.Contains(string(message), "\"result\"") {
						socket.send <- message

					} else if strings.Contains(string(message), "\"accounts\"") || strings.Contains(string(message), "\"userID\"") {
						// probably send to a pool worker routine and buffer channel
						// so this method doesn't cause any blocking
						account := Conn.activeAccount(message)
						account.receiver <- message

					} else {
						fmt.Println("General:")
						Conn.Account.receiver <- message
					}
				}
			}()

			return i, socket
		}
	}

	return -1, nil
}

// // The message types are defined in RFC 6455, section 11.8.
// const (
// 	// TextMessage denotes a text data message. The text message payload is
// 	// interpreted as UTF-8 encoded text data.
// 	TextMessage = 1

// 	// BinaryMessage denotes a binary data message.
// 	BinaryMessage = 2

// 	// CloseMessage denotes a close control message. The optional message
// 	// payload contains a numeric code and text. Use the FormatCloseMessage
// 	// function to format a close message payload.
// 	CloseMessage = 8

// 	// PingMessage denotes a ping control message. The optional message payload
// 	// is UTF-8 encoded text.
// 	PingMessage = 9

// 	// PongMessage denotes a pong control message. The optional message payload
// 	// is UTF-8 encoded text.
// 	PongMessage = 10
// )

// Client makes a socket connection
//	-> listen for any incoming responses in a go Routine
//	-> go routine routes the incoming message based on
//	->	-> public general message
//	-> 	-> private account message (aop.subscribe -> only one auth connection per socket)
//	-> 	-> subscribing confirmation message base from result == success

// Client.Subscribe
// 	-> returns a socket that is available
// 	-> or creates and returns a new socekt
// 	-> or returns nil if at max capacity

// Account.Subscribe
//	-> connects to an available socket that subscribes to any subscription channel
//	-> generates a listener for account if one is not currently running
//	-> waits for a response and then checks if successful connection or not
//	-> if succesful updates connection map and account subscription map.
//	-> then returns

// Account.Auth
//	-> connects to an available socket that subscribes to any subscription channel
//	-> generates a signature and expirey
//	-> sends message
//	-> waits for response and checks if auth was succesful
//	-> if succesfuly updates connnection map
//	-> returns

// NOTES:
// only one auth connection can be made per a socekt connection
// if multiple auth connections were made only the last will be connected.
// overwritting any previous auth connection.
// Need to implement some type of check to either prevent any new auth connection
// or unsub first updating local changes and then resub with new auth.
// subbing to "aop.subscribe" streams data if (active order?) or position.
// Need to implement reconnection if not receiving any pong messages.
// Need a better way of initiating main listener.
// Need to create public account and listener.
// Need to spin up more go routines when receiving massive amount of messages in a short given time.
// Can create a writter go routine for subbing.
// Right now subbing is a blocking action. which prevents the socket writter being used concurrently.
// Can also implment subbing async with channels to prevent socket writter being used concurrently.
// Not sure which methode will be best. but will stick with blocking method for now.
// When authorizing maybe send bool with a channel and remove the udpate state
// from auth and keep only in subbing. will be cleaner and state is only updated in one location.
// listener needs a handler passed into it so it can pass off messages to be handled in anyway fit.
// Thinking about revising the keepAlive function that I took from another repo,
// so I can fully understand how it works. and make improvements on it.
