package coincap

import (
	"coin_struct"
	"encoding/json"
	"log"
	"restful_query"
	"time"

	"github.com/jinzhu/copier"
)

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
	PriceEur       float64 `json:"price_eur"`
	PriceUsd       float64 `json:"price_usd"`
	QueryTimeStamp int64   `json:"query_timestamp"`
}

type Coin struct {
	Id             string  `json:"id"`
	DisplayName    string  `json:"display_name"`
	PriceBtc       float64 `json:"price_btc,omitempty"`
	PriceEth       float64 `json:"price_eth,omitempty"`
	PriceUsd       float64 `json:"price_usd,omitempty"`
	QueryTimeStamp int64   `json:"query_timestamp"`
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
	bodyBytes, err := restful_query.Get(apiUrl + "page/" + product)
	if err != nil {
		log.Fatalln(err)
	}
	page := CoinCapPage{Id: product,
		QueryTimeStamp: time.Now().Unix()}
	json.Unmarshal(bodyBytes, &page)
	return page
}

func GetCoinCapCoin(id string) coin_struct.Coin {
	page := get_page(id)
	coinCapCoin := coin_struct.Coin{}
	copier.Copy(&coinCapCoin, &page)
	return coinCapCoin
}
