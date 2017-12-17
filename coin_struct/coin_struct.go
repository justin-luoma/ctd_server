package coin_struct

type Coin struct {
	ID             string  `json:"id,omitempty"`
	DisplayName    string  `json:"display_name,omitempty"`
	PriceBtc       float64 `json:"price_btc,omitempty"`
	PriceEth       float64 `json:"price_eth,omitempty"`
	PriceLtc       float64 `json:"price_ltc,omitempty"`
	PriceUsd       float64 `json:"price_usd,omitempty"`
	PriceEur       float64 `json:"price_eur,omitempty"`
	PriceGbp       float64 `json:"price_gbp,omitempty"`
	QueryTimeStamp int64   `json:"query_timestamp,omitempty"`
}
