package gdax

import (
	"math/big"
	"net/http"
)

var apiUrl string = "https://api.gdax.com/"

type GdaxCurrencyResponse struct {
	Collection []GdaxCurrency
}

type GdaxCurrency struct {
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
	Open           big.Float `json:"open"`
	High           big.Float `json:"high"`
	Low            big.Float `json:"low"`
	Volume         big.Float `json:"volume"`
	Last           big.Float `json:"last"`
	Volume30Day    big.Float `json:"volume_30day"`
	QueryTimeStamp int64     `json:"query_timestamp"`
}

func pull_currencies() {
	currencies := make([]GdaxCurrency, 0)
	response, err := http.Get("")
}
