package main

// func test() {
// 	jsonFile, err := os.Open("./config.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	defer jsonFile.Close()

// 	byteValue, err := ioutil.ReadAll(jsonFile)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	type Secret struct {
// 		Key string `json:"SECRET"`
// 	}

// 	client := make(map[string]interface{})

// 	json.Unmarshal(byteValue, &client)

// 	fmt.Println(client)

// 	// client.hmac = hmac.New(crypto.SHA256.New, []byte(client.Secret))
// }
