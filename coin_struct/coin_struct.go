package coin_struct

type Coin struct {
	ID             string  `json:"id"`
	DisplayName    string  `json:"display_name"`
	PriceBtc       float64 `json:"price_btc"`
	PriceEth       float64 `json:"price_eth"`
	PriceUsd       float64 `json:"price_usd"`
	QueryTimeStamp int64   `json:"query_timestamp"`
}
