package phemex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
)

type Socket struct {
	Client *Client
	Conn   *websocket.Conn
	send   chan []byte
}

func MainSub() {
	// client.Account.Auth().Subscribe("aop.subscribe", []interface{}{})
	fmt.Println("hello")
}

func Subscribe() {
	client.Account.Subscribe("orderbook.subscribe", []interface{}{"BTCUSD"})
	// client.Account.Subscribe("trade.subscribe", []interface{}{"BTCUSD"})
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

// Auth -> allows account to be authorized before subsribing to a channel
func (a *Account) Auth() *Account {

	_, socket := a.Client.Subscribe()
	if socket == nil {
		fmt.Println("Max Connections, no more connections can be made.")
		return nil
	}

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

	func() {
		err := socket.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
		}
	}()

	return a
}

// Subscribe -> allows account to make a new subscribtion channel if not at max capacity
func (a *Account) Subscribe(method string, params []interface{}) *Account {

	index, socket := a.Client.Subscribe()
	if socket == nil {
		fmt.Println("Max Connections, no more connections can be made.")
		return a
	}

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

		message := <-socket.send

		data := *parseMessage(message)

		if data["result"].(map[string]interface{})["status"].(string) == "success" {
			a.Subscriptions[method] = index
			client.ConnMap[index] += 1
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

					message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

					if strings.Contains(string(message), "{\"accounts\"") {
						data := *parseMessage(message)
						userID := int64(data["accounts"].([]interface{})[0].(map[string]interface{})["userID"].(float64))
						account := Conn.Account.Accounts[userID]
						account.receiver <- message
					} else if strings.Contains(string(message), "{\"results\"") {
						socket.send <- message

					} else {
						Conn.Account.receiver <- message
					}
				}
			}()

			return i, socket
		}
	}

	return -1, nil
}

// websocket to the hub
// hub to the account
