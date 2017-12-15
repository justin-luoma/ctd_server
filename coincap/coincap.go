package coincap

import (
	"encoding/json"
	"net/http"
	"fmt"
	"time"
	"log"
)

var apiUrl string = "https://coincap.io"

type CoinCapPage struct {
	ID             string  `json:"id"`
	DisplayName    string  `json:"display_name"`
	Cap24HrChange  float64 `json:"cap24hrChange"`
	PriceBtc       float64 `json:"price_btc"`
	PriceEur       float64 `json:"price_eur"`
	PriceUsd       float64 `json:"price_usd"`
	QueryTimeStamp int64   `json:"query_timestamp"`
}

func GetCoinCapCoin(id string) *CoinCapPage {
	coin := CoinCapPage{QueryTimeStamp: time.Now().Unix()}

	url := fmt.Sprintf("%s/page/%s",apiUrl, id)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		err := json.NewDecoder(response.Body).Decode(&coin)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(coin)
	}
	return &coin
}