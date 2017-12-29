package main

import (
	"coin_struct"
	"decimal_math"
	json2 "encoding/json"
	"exchange_api_status"
	"flag"
	"fmt"
	"restful_query"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

const apiUrl = "https://poloniex.com/public"

var apiCalls = 0
var apiCallTime int64

var poloniexDataSet = struct {
	sync.RWMutex
	Coin map[string]interface{}
}{Coin: make(map[string]interface{})}

var currencies = make(map[string]string)

func init() {
	flag.Parse()
	build_data_set()
}

func TestInit() {
	fmt.Println("poloniex test init")
}

func api_call_wrapper(url string) (*[]byte, error) {
	var bodyBytes []byte
	var err error
	switch apiCalls {
	case 0:
		apiCallTime = time.Now().Unix()
		bodyBytes, err = restful_query.Get(url)
		apiCalls++
	case 6:
		timePassed := time.Since(time.Unix(apiCallTime, 0)).Seconds()
		if timePassed <= 1 {
			time.Sleep(time.Second)
		}
		apiCallTime = time.Now().Unix()
		bodyBytes, err = restful_query.Get(url)
		apiCalls = 1
	default:
		bodyBytes, err = restful_query.Get(url)
		apiCalls++
	}

	return &bodyBytes, err
}

func get_currencies() (map[string]interface{}, error) {
	bodyBytes, err := api_call_wrapper(apiUrl + "?command=returnCurrencies")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	/*
		currencyData structure
		{
			"BTC": {
				"id": 28,
				"name": "Bitcoin",
				"txFee": "0.00050000",
				"minConf": 1,
				"depositAddress": null,
				"disabled": 0,
				"delisted": 0,
				"frozen": 0
			},
			"ETH": {
				"id": 267,
				"name": "Ethereum",
				"txFee": "0.00500000",
				"minConf": 35,
				"depositAddress": null,
				"disabled": 0,
				"delisted": 0,
				"frozen": 0
			},
		}
	*/
	var currencyData map[string]interface{}
	err = json2.Unmarshal(*bodyBytes, &currencyData)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	return currencyData, nil
}

func get_ticker() (map[string]interface{}, error) {
	bodyBytes, err := api_call_wrapper(apiUrl + "?command=returnTicker")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	/*
		tickerData structure
		{
			"BTC_VTC": {
				"id": 100,
				"last": "0.00050451",
				"lowestAsk": "0.00050700",
				"highestBid": "0.00050452",
				"percentChange": "0.06749751",
				"baseVolume": "201.18918617",
				"quoteVolume": "418982.01792770",
				"isFrozen": "0",
				"high24hr": "0.00051842",
				"low24hr": "0.00043091"
			},
			"USDT_BTC": {
				"id": 121,
				"last": "14655.00000000",
				"lowestAsk": "14655.00000000",
				"highestBid": "14640.11308236",
				"percentChange": "0.03936170",
				"baseVolume": "99380539.61794215",
				"quoteVolume": "6990.21453774",
				"isFrozen": "0",
				"high24hr": "15100.00000000",
				"low24hr": "13335.00000000"
			},
		}
	*/
	var tickerData map[string]interface{}
	err = json2.Unmarshal(*bodyBytes, &tickerData)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	return tickerData, nil
}

func build_data_set() {
	poloniexDataSet.Lock()
	defer poloniexDataSet.Unlock()

	coin := poloniexDataSet.Coin

	queryTime := time.Now().Unix()
	tickerData, err := get_ticker()
	if err != nil {
		glog.Errorln(err)
	}
	if len(tickerData) == 0 {
		glog.Errorln("Poloniex didn't return any data")
		exchange_api_status.Update_Status("poloniex", 0)
		return
	}

	exchange_api_status.Update_Status("poloniex", 1)

	currencyData, err := get_currencies()
	if err != nil {
		glog.Errorln(err)
	}

	for currencyPair, pairData := range tickerData {
		//Split currencyPair ("BTC_LTC") into base and market values ("BTC" and "LTC")
		splitPair := strings.Split(currencyPair, "_")
		baseCurrency := splitPair[0]
		marketCurrency := splitPair[1]

		//keep track of base currencies for use later in valid market verification
		currencies[baseCurrency] = "base"

		//type assertion for getting the currency info for the base name
		baseDisplayname := currencyData[baseCurrency].(map[string]interface{})["name"]

		/*
			priceFloat takes the string pairData["last"] and converts it to a
			float64 with the required type assertion if there is an error, skip
			the entire currency pair and report the error
		*/
		priceFloat, err := decimal_math.Convert_String_To_Float64(
			pairData.(map[string]interface{})["last"].(string),
			8,
			false)
		if err != nil {
			glog.Warningf("%s%s%s\n%s\n",
				"Skipping currency pair: ", currencyPair, "due to float conversion error:", err)
			continue
		}

		//same as priceFloat but for delta
		delta, err := decimal_math.Convert_String_To_Float64(
			pairData.(map[string]interface{})["percentChange"].(string),
			2,
			true)
		if err != nil {
			glog.Warningf("%s%s%s\n%s\n",
				"Skipping currency pair: ", currencyPair, "due to float conversion error:", err)
			continue
		}

		/*
			Build the underlying data structure for coin data
			if coin["BTC"] the data looks like the following:
			"DisplayName": "Bitcoin",
			"ETH": {
				"Delta": 1.55,
				"Price": 0.05107932,
				"QueryTimeStamp": 1514582938
			},
			"VTC": {
				"Delta": -0.19,
				"Price": 0.0004863,
				"QueryTimeStamp": 1514582938
			},
		*/
		coinData := map[string]interface{}{
			"DisplayName": baseDisplayname,
			marketCurrency: map[string]interface{}{
				"Price":          priceFloat,
				"Delta":          delta,
				"QueryTimeStamp": queryTime,
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

	//TESTING
	//fmt.Println(currencies)
	/*jsonData, err := json2.MarshalIndent(coin, "", " ")
	fmt.Println(string(jsonData), err)*/
}

//TODO write Get_Coins function
/*
func Get_Coins() ([]coin_struct.Coin, error) {

}*/

//TODO write Get_Coin_Stats function
/*
func Get_Coin_Stats(coinId string) (*map[string]interface{}, error) {

}
 */
