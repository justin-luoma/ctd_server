package bitstamp

import (
	"coin_struct"
	"decimal_math"
	"exchange_api_status"
	"errors"
	"flag"
	"github.com/golang/glog"
	"strings"
	"sync"
	"time"
)

const apiUrl = "https://www.bitstamp.net/api/"

var currencyTypes = map[string]string{
	"BTC": "crypto",
	"BCH": "crypto",
	"ETH": "crypto",
	"LTC": "crypto",
	"XRP": "crypto",
	"USD": "fiat",
	"EUR": "fiat",
}

var bitstampDataSet = struct {
	sync.RWMutex
	Coin map[string]interface{}
}{Coin: make(map[string]interface{})}

func init() {
	flag.Parse()
	build_bitstamp_dataset()
}

func is_valid_coin(coinId string) bool {
	if _, ok := currencyTypes[coinId]; ok && currencyTypes[coinId] == "crypto" {
		return true
	} else {
		return false
	}
}

func valid_product_stats(baseCurrency string, quoteCurrency string) bool {
	bitstampDataSet.RLock()
	defer bitstampDataSet.RUnlock()

	if _, ok := bitstampDataSet.Coin[baseCurrency].(map[string]interface{})[quoteCurrency]; ok {
		return true
	} else {
		return false
	}
}

func is_data_old(coinId string, maxAgeSeconds int) bool {
	var dataOld bool = false

	bitstampDataSet.RLock()
	defer bitstampDataSet.RUnlock()
	coin := bitstampDataSet.Coin[coinId].(map[string]interface{})

	for quote := range currencyTypes {
		if !valid_product_stats(coinId, quote) {
			continue
		}
		dataAge := time.Since(
			time.Unix(coin[quote].(map[string]interface{})["QueryTimestamp"].(int64), 0)).Seconds()
		if int(dataAge) > maxAgeSeconds {
			dataOld = true
			return dataOld
		}
	}

	return dataOld
}

func build_json_struct(coinId string) *map[string]interface{} {
	bitstampDataSet.RLock()
	defer bitstampDataSet.RUnlock()
	/*
		tell go that bitstampDataSet.Coin[coinId] is a map[string]interface{}
		and set bitstampData equal the coin we want
	*/
	bitstampData := bitstampDataSet.Coin[coinId].(map[string]interface{})
	/*
		jsonData will hold all of the data for the coin
		since id and displayname are static we can set them outside the loop
	*/
	jsonData := map[string]interface{}{
		"id":           coinId,
		"display_name": bitstampData["DisplayName"],
	}
	for quote := range currencyTypes {
		if !valid_product_stats(coinId, quote) {
			continue
		}
		/*
			quoteData hold the data for the current quote currency in the loop,
			while quoteTmp is hold the structure for the json formatted data
		*/
		quoteData := bitstampData[quote].(map[string]interface{})
		quoteTmp := map[string]interface{}{
			strings.ToLower(quote) + "_price":           quoteData["Price"],
			strings.ToLower(quote) + "_24_hour_change":  quoteData["Delta"],
			strings.ToLower(quote) + "_query_timestamp": quoteData["QueryTimestamp"],
		}

		for k, v := range quoteTmp {
			jsonData[k] = v
		}
	}
	/*jsonString, _ := json2.Marshal(jsonData)

	fmt.Println(string(jsonString))*/

	return &jsonData
}

//noinspection ALL
func build_bitstamp_dataset() {
	uP, _, err := get_online_products()
	if err != nil && uP == nil {
		glog.Errorln("Unable to initialize Bitstamp package, check error log")
		exchange_api_status.Update_Status("bitstamp", 0)
		return
	} else {
		exchange_api_status.Update_Status("bitstamp", 1)

		if err != nil {
			glog.Errorln("Unable to get Bitstamp currencies")
			exchange_api_status.Update_Status("bitstamp", 0)
			return
		}

		for c, t := range currencyTypes {
			//we only want to
			if t == "crypto" {
				update_coin_data(c, &uP)
			}
		}

		//TEST SECTION//
		//_ = build_json_struct("BTC")
		/*bitstampDataSet.RLock()
		jsonData, _ := json2.MarshalIndent(bitstampDataSet.Coin, "", " ")
		bitstampDataSet.RUnlock()
		fmt.Println(string(jsonData))*/
	}
}

//noinspection ALL
func update_coin_data(coinId string, onlineProducts *[]BitstampProducts) {
	glog.V(2).Infoln("update_coin_data " + coinId)

	if is_valid_coin(coinId) {

		bitstampDataSet.Lock()
		defer bitstampDataSet.Unlock()

		coin := bitstampDataSet.Coin

		currencyNames := make(map[string]string)

		for _, currency := range *onlineProducts {
			splitName := strings.Split(currency.Name, "/")
			BaseCurrency := strings.Trim(strings.Split(currency.Description, "/")[0]," ")
			currencyNames[splitName[0]] = BaseCurrency
		}

		for _, product := range *onlineProducts {
			//only want products with the specified coin: coinId=BTC productId=BTC-USD/BTC-EUR
			if strings.HasPrefix(product.Name, coinId) {

				stats, err := get_product_stats(product.URLSymbol)
				if err != nil {
					glog.Warningln("Failed to retriece stats for product: " + product.URLSymbol)
					continue
				}
				delta := decimal_math.Calculate_Percent_Change_Float(stats.Open, stats.Last)
				/*
					build the structure for the coin:
					{
					 "DisplayName": "Bitcoin",
					 "USD": {
					  "Delta": -5.63,
					  "Price": 16336.62,
					  "QueryTimestamp": 1513825686
					 }
					}
				*/
				QuoteCurrency := strings.Trim(strings.Split(product.Name, "/")[1]," ")
				coinData := map[string]interface{}{
					"DisplayName": currencyNames[coinId],
					QuoteCurrency: map[string]interface{}{
						"Price":          stats.Last,
						"Delta":          delta,
						"QueryTimestamp": stats.Timestamp,
					},
				}

				if coin[coinId] == nil {
					coin[coinId] = coinData
				} else {
					/*
						since coin is a map[string]interface{} and an interface
						can be anything we havc tell go what it is
						in this case it's another map[string]interface{}
						but as far as go knows it could be map[stuct]string or anything else
					*/
					tmp := coin[coinId].(map[string]interface{})
					for k, v := range coinData {
						tmp[k] = v
					}
				}
			}
		}

	} else {
		glog.Warningln("update_coin_data: invalid Coin ID: " + coinId)
		return
	}
}

//noinspection ALL
func Get_Coins() ([]coin_struct.Coin, error) {
	var coin coin_struct.Coin
	var coins []coin_struct.Coin

	var err error

	currencyNames := make(map[string]string)
	uP, _, err := get_online_products()
	for _, currency := range uP {
		splitName := strings.Split(currency.Name, "/")
		BaseCurrency := strings.Trim(strings.Split(currency.Description, "/")[0]," ")
		currencyNames[splitName[0]] = BaseCurrency
	}

	for c, t := range currencyTypes {
		if err != nil {
			glog.Error(err)
			glog.Error("invalid response for Bitstamp currency: " + c)
		} else if t == "crypto" {
			coin.ID = c
			coin.DisplayName = currencyNames[c]
			coin.IsActive = true
			coins = append(coins, coin)
		} else {
			continue
		}
	}

	return coins, err
}

//noinspection ALL
func Get_Coin_Stats(coinId string) (*map[string]interface{}, error) {
	if !is_valid_coin(coinId) {
		err := errors.New("invalid coinId id: " + coinId)
		return nil, err
	}

	if is_data_old(coinId, 10) {
		uP, _, err := get_online_products()
		if err != nil || uP == nil {
			glog.Errorln("Bitstamp package offline, check error log")
			exchange_api_status.Update_Status("bitstamp", 0)
			return nil, errors.New("Bitstamp API is down")
		} else {
			exchange_api_status.Update_Status("bitstamp", 1)
			update_coin_data(coinId, &uP)
		}
	}

	jsonData := build_json_struct(coinId)

	return jsonData, nil
}
