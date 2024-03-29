package coin_struct

type Coin struct {
	ID               string  `json:"id,omitempty"`
	DisplayName      string  `json:"display_name,omitempty"`
	PriceBtc         float64 `json:"price_btc,omitempty"`
	PriceEth         float64 `json:"price_eth,omitempty"`
	PriceLtc         float64 `json:"price_ltc,omitempty"`
	PriceUsd         float64 `json:"price_usd,omitempty"`
	PriceEur         float64 `json:"price_eur,omitempty"`
	PriceGbp         float64 `json:"price_gbp,omitempty"`
	DayDeltaPriceUsd float64 `json:"24hour_price_usd,omitempty"`
	DayDeltaPriceEur float64 `json:"24hour_price_eur,omitempty"`
	DayDeltaPriceBgp float64 `json:"24hour_price_gbp,omitempty"`
	DayDeltaPriceBtc float64 `json:"24hour_price_btc,omitempty"`
	QueryTimeStamp   int64   `json:"query_timestamp,omitempty"`
	IsActive         bool    `json:"is_active"`
	StatusMessage    string  `json:"status_message,omitempty"`
	IsFrozen         bool    `json:"is_frozen"`
	IsFiat           bool    `json:"is_fiat,omitempty"`
}
