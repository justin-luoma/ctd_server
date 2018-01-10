package bitfinex

import "sync"

const apiUrl = "https://api.bitfinex.com/v1/"

var apiCalls = 0
var apiCallTime int64

var bitfinexDataSet struct{
	sync.RWMutex
	Coin map[string]interface{}
}{Coin: make(map[string]interface{})}

