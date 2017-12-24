package bittrex

import (
	"github.com/toorop/go-bittrex"
	"sync"
	"github.com/golang/glog"
	"strings"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

var bittrexDataSet = struct {
	sync.RWMutex
	Coin map[string]interface{}
}{Coin: make(map[string]interface{})}

var b = bittrex.New(API_KEY, API_SECRET)

func build_data_set() {
	bittrexDataSet.Lock()
	defer bittrexDataSet.Unlock()

	coin := bittrexDataSet.Coin

	currencies, err := b.GetCurrencies()
	if err != nil {
		glog.Fatalln(err)
	}

	currencyNames := make(map[string]string)

	for _, currency := range currencies{
		currencyNames[currency.Currency] = currency.CurrencyLong
	}

	markets, err := b.GetMarketSummaries()
	if err != nil {
		glog.Fatalln(err)
	}

	for _, market := range markets {
		/*splitProduct := strings.Split(productId, "-")
		baseCurrency := splitProduct[0]
		quoteCurrency := splitProduct[1]*/
		splitStr := strings.Split(market.MarketName, "-")
		baseCurrency := splitStr[0]
		marketCurrency := splitStr[1]

		coinData := map[string]interface{}{
			"DisplayName": currencyNames[]
		}
	}
}
