package gdax

import (
	"../restful_query"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"
)

var apiUrl string = "https://api.gdax.com/"
var supportedCoins = [3]string{"BTC", "ETH", "LTC"}

type GdaxCurrencies struct {
	Id      string  `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	MinSize float32 `json:"min_size,omitempty"`
	Status  string  `json:"status,omitempty"`
	Message string  `json:"message,omitempty"`
}

type GdaxProducts struct {
	Id     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

type GdaxStats struct {
	Product        string    `json:"id"`
	Open           string `json:"open"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Volume         string `json:"volume"`
	Last           string `json:"last"`
	Volume30Day    string `json:"volume_30day"`
	QueryTimeStamp int64     `json:"query_timestamp"`
}

type Coin struct {
	Id             string    `json:"id"`
	DisplayName    string    `json:"display_name"`
	PriceBtc       big.Float `json:"price_btc"`
	PriceEth       big.Float `json:"price_eth"`
	PriceUsd       big.Float `json:"price_usd"`
	QueryTimeStamp int64     `json:"query_timestamp"`
}

func get_currencies() (*[]GdaxCurrencies, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies")
	if err != nil {
		log.Fatalln(err)
	}
	var currencies []GdaxCurrencies
	json.Unmarshal(bodyBytes, &currencies)

	return &currencies, nil
}

func get_products() ([]GdaxProducts, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "products")
	if err != nil {
		log.Fatalln(err)
	}
	var products []GdaxProducts
	json.Unmarshal(bodyBytes, &products)

	return products, nil
}

func get_stat() /*(*[]GdaxStats, error)*/ {
	products, err := get_products()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(products)
	var stats []GdaxStats
	for i, product := range products {
		if product.Status == "online" {
			bodyBytes, err := restful_query.Get(apiUrl + "products/" + product.Id + "/stats")
			if err != nil {
				log.Fatalln(err)
			}
			stat := GdaxStats{Product: product.Id,
				QueryTimeStamp: time.Now().Unix()}
			json.Unmarshal(bodyBytes, &stat)
			stats = append(stats, stat)
			fmt.Println(stats[i])
		}
	}
	fmt.Println(stats)
}
