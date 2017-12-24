package bittrex

import (
	"testing"
	"fmt"
	json2 "encoding/json"
	"github.com/toorop/go-bittrex"
)

func TestBittrex(t *testing.T) {
	b := bittrex.New("","")
	//markets, _ := b.GetCurrencies()
	markets, _ := b.GetMarketSummaries()
	json, err := json2.MarshalIndent(markets, "", " ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(json))
}
