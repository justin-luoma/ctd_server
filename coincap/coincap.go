package coincap

import (
	"../restful_query"
	"encoding/json"
	"log"
	"math/big"
	"time"
	"github.com/golang/glog"
	"flag"
)

func init()  {
	flag.Parse()
}

var apiUrl string = "https://coincap.io/"

type CoinCapMap struct {
	Aliases []interface{} `json:"aliases"`
	Name    string        `json:"name,omitempty"`
	Symbol  string        `json:"symbol,omitempty"`
}

type CoinCapPage struct {
	Id             string  `json:"id"`
	DisplayName    string  `json:"display_name"`
	Cap24HrChange  float64 `json:"cap24hrChange"`
	PriceBtc       float64 `json:"price_btc"`
	PriceEth       float64 `json:"price_eth"`
	PriceUsd       float64 `json:"price_usd"`
	QueryTimeStamp int64   `json:"query_timestamp"`
}

type Coin struct {
	Id             string    `json:"id"`
	DisplayName    string    `json:"display_name"`
	PriceBtc       big.Float `json:"price_btc"`
	PriceEth       big.Float `json:"price_eth"`
	PriceUsd       big.Float `json:"price_usd"`
	QueryTimeStamp int64     `json:"query_timestamp"`
}

func get_map() ([]CoinCapMap, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "map")
	if err != nil {
		log.Fatalln(err)
	}
	var coinMap []CoinCapMap
	json.Unmarshal(bodyBytes, &coinMap)

	return coinMap, nil
}

func get_page(product string) CoinCapPage {
	glog.V(2).Infoln("Testing")
	bodyBytes, err := restful_query.Get(apiUrl + "page/" + product)
	if err != nil {
		log.Fatalln(err)
	}
	page := CoinCapPage{Id: product,
		QueryTimeStamp: time.Now().Unix()}
	json.Unmarshal(bodyBytes, &page)
	return page
}

func GetCoinCapCoin(id string) *CoinCapPage {
	page := get_page(id)

	return &page
}