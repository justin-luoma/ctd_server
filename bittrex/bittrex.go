package bittrex

import (
	json2 "encoding/json"
	"fmt"
	"github.com/toorop/go-bittrex"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func test() {
	bittrex := bittrex.New(API_KEY, API_SECRET)

	//markets, _ := bittrex.GetMarkets()
	//markets, _ := bittrex.GetMarketSummary("USDT-BTC")
	markets, _ := bittrex.GetMarketSummaries()
	json, err := json2.MarshalIndent(markets, "", " ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(json))
}
