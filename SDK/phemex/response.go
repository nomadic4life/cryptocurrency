package phemex

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	data   []byte                 // rawData, blob
	output map[string]interface{} // data
	req    *Request
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
