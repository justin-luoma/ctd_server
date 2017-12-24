package bittrex

import (
	"github.com/toorop/go-bittrex"
	"sync"
	"github.com/golang/glog"
	"strings"
	"../decimal_math"
	"time"
	"fmt"
	json2 "encoding/json"
	"../coin_struct"
	"errors"
	"flag"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

var bittrexDataSet = struct {
	sync.RWMutex
	Coin map[string]interface{}
}{Coin: make(map[string]interface{})}

var baseCurrencies = map[string]string{
	"BTC": "base",
	"ETH": "base",
	"USDT": "base",
}

var b = bittrex.New(API_KEY, API_SECRET)

func init() {
	flag.Parse()
	build_data_set()
}

func Tinit() /*(map[string]interface{})*/ {
	build_data_set()

	//TEST
	/*bittrexDataSet.RLock()
	bittrexData := bittrexDataSet.Coin["BTC"].(map[string]interface{})
	defer bittrexDataSet.RUnlock()
	for k, v := range bittrexData {
		switch v.(type) {
		case map[string]interface{}:
			fmt.Println(k)
		}
	}*/
	jsonData := build_json_data("BTC")
	rtJson, _ := json2.MarshalIndent(jsonData, "", " ")

	fmt.Println(string(rtJson))
}

func is_valid_coin(coinId string) bool {
	for currency := range baseCurrencies {
		if currency == coinId {
			return true
		}
	}

	return false
}

func is_valid_market(baseCurrency string, marketCurrency string) bool {
	bittrexDataSet.RLock()
	defer bittrexDataSet.RUnlock()

	if _, ok := bittrexDataSet.Coin[baseCurrency].(map[string]interface{})[marketCurrency]; ok {
		return true
	} else {
		return false
	}
}

func is_data_old(coinId string, maxAgeSeconds int) bool {
	var dataOld bool = false

	bittrexDataSet.RLock()
	defer bittrexDataSet.RUnlock()

	coin := bittrexDataSet.Coin[coinId].(map[string]interface{})

	for market, v := range coin {
		switch v.(type) {
		case map[string]interface{}:
			dataAge := time.Since(
				time.Unix(coin[market].(map[string]interface{})["QueryTimeStamp"].(int64), 0)).Seconds()
			if int(dataAge) > maxAgeSeconds {
				dataOld = true
				return dataOld
			}
		}
	}

	return dataOld
}

func build_json_data(coinId string) *map[string]interface{} {
	bittrexDataSet.RLock()
	defer bittrexDataSet.RUnlock()

	bittrexData := bittrexDataSet.Coin[coinId].(map[string]interface{})

	jsonData := map[string]interface{}{
		"id":	coinId,
		"display_name":	bittrexData["DisplayName"],
	}

	for market, v := range bittrexData {
		switch v.(type) {
		case map[string]interface{}:

			marketData := bittrexData[market].(map[string]interface{})
			marketTmp := map[string]interface{}{
				strings.ToLower(market) + "_price": marketData["Price"],
				strings.ToLower(market) + "_24_hour_change": marketData["Delta"],
				strings.ToLower(market) + "_query_timestamp": marketData["QueryTimeStamp"],
			}

			for k, v := range marketTmp {
				jsonData[k] = v
			}

		}
	}

	return &jsonData
}

func build_data_set() {
	bittrexDataSet.Lock()
	defer bittrexDataSet.Unlock()

	coin := bittrexDataSet.Coin

	currencies, err := b.GetCurrencies()
	if err != nil {
		glog.Fatalln(err)
	}

	currencyNames := make(map[string]string)

	for _, currency := range currencies {
		currencyNames[currency.Currency] = currency.CurrencyLong
	}

	bittrexMarkets, err := b.GetMarketSummaries()
	if err != nil {
		glog.Fatalln(err)
	}

	for _, market := range bittrexMarkets {
		/*splitProduct := strings.Split(productId, "-")
		baseCurrency := splitProduct[0]
		quoteCurrency := splitProduct[1]*/
		splitStr := strings.Split(market.MarketName, "-")
		baseCurrency := splitStr[0]
		marketCurrency := splitStr[1]

		delta := decimal_math.Calculate_Percent_Change_Decimal(market.PrevDay, market.Last)
		timeStamp, _ := time.Parse(time.RFC3339, market.TimeStamp + "Z")
		priceFloat := decimal_math.Convert_To_Float64(market.Last)

		coinData := map[string]interface{}{
			"DisplayName": currencyNames[baseCurrency],
			marketCurrency: map[string]interface{}{
				"Price": priceFloat,
				"Delta": delta,
				"QueryTimeStamp": timeStamp.Unix(),
			},
		}

		if coin[baseCurrency] == nil {
			coin[baseCurrency] = coinData
		} else {
			tmp := coin[baseCurrency].(map[string]interface{})
			for k, v := range coinData {
				tmp[k] = v
			}
		}
	}
}

func Get_Coins() ([]coin_struct.Coin, error) {
	var coin coin_struct.Coin
	var coins []coin_struct.Coin

	var err error

	currencies, err := b.GetCurrencies()
	if err != nil {
		glog.Fatalln(err)
	}

	for _, currency := range currencies {
		coin.ID = currency.Currency
		coin.DisplayName = currency.CurrencyLong
		coin.IsActive = currency.IsActive
		coin.StatusMessage = currency.Notice
		coins = append(coins, coin)
	}

	return coins, err
}

func Get_Coin_Stats(coinId string) (*map[string]interface{}, error) {
	if !is_valid_coin(coinId) {
		err := errors.New("invalid coinId id: " + coinId)
		return nil, err
	}
	if is_data_old(coinId, 10) {
		build_data_set()
	}

	jsonData := build_json_data(coinId)

	return jsonData, nil
}