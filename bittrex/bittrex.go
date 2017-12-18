package bittrex

import (
	"github.com/toorop/go-bittrex"
	"fmt"
	json2 "encoding/json"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func test()  {
	bittrex := bittrex.New(API_KEY, API_SECRET)

	markets, _ := bittrex.GetMarkets()
	json, err := json2.MarshalIndent(markets, "", " ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(json))
}