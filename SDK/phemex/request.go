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
