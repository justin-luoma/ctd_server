package gdax

import (
	"coin_struct"
	"decimal_math"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
)

const apiUrl = "https://api.gdax.com/"

//var GdaxCurrencies []coin_struct.Coin

var onlineProductIds []string

var gdaxDataSet = struct {
	sync.RWMutex
	Coin map[string]interface{}
}{Coin: make(map[string]interface{})}

var gC *gdaxCurrencies

var gP *gdaxProducts

var gS *gdaxStats

func Init() {
	gC = init_currencies()
	gP = init_products()
	gS = init_stats()
	build_gdax_dataset()
}

// func check_online_status() bool {
// 	_, _, err := get_online_products()
// 	if err != nil {
// 		return false
// 	}

// 	return true
// }

func is_valid_coin(coinId string) bool {
	for _, coin := range *gC.get_coins() {
		if coinId == coin.ID {
			return true
		}
	}

	return false
}

func valid_product_stats(baseCurrency string, quoteCurrency string) bool {
	gdaxDataSet.RLock()
	defer gdaxDataSet.RUnlock()

	if _, ok := gdaxDataSet.Coin[baseCurrency].(map[string]interface{})[quoteCurrency]; ok {
		return true
	} else {
		return false
	}
}

func is_data_old(coinId string, maxAgeSeconds int) bool {
	var dataOld bool = false

	gdaxDataSet.RLock()
	defer gdaxDataSet.RUnlock()
	coin := gdaxDataSet.Coin[coinId].(map[string]interface{})

	for _, quote := range *gC.get_currencies() {
		if !valid_product_stats(coinId, quote.ID) {
			continue
		}
		dataAge := time.Since(
			time.Unix(coin[quote.ID].(map[string]interface{})["QueryTimestamp"].(int64), 0)).Seconds()
		if int(dataAge) > maxAgeSeconds {
			dataOld = true
			return dataOld
		}
	}

	return dataOld
}

func build_json_struct(coinId string) *map[string]interface{} {
	gdaxDataSet.RLock()
	defer gdaxDataSet.RUnlock()
	/*
		tell go that gdaxDataSet.Coin[coinId] is a map[string]interface{}
		and set gdaxData equal the coin we want
	*/
	gdaxData := gdaxDataSet.Coin[coinId].(map[string]interface{})
	/*
		jsonData will hold all of the data for the coin
		since id and displayname are static we can set them outside the loop
	*/
	jsonData := map[string]interface{}{
		"id":           coinId,
		"display_name": gdaxData["DisplayName"],
	}
	for _, quote := range *gC.get_currencies() {
		if !valid_product_stats(coinId, quote.ID) {
			continue
		}
		/*
			quoteData hold the data for the current quote currency in the loop,
			while quoteTmp is hold the structure for the json formatted data
		*/
		quoteData := gdaxData[quote.ID].(map[string]interface{})
		quoteTmp := map[string]interface{}{
			strings.ToLower(quote.ID) + "_price":           quoteData["Price"],
			strings.ToLower(quote.ID) + "_24_hour_change":  quoteData["Delta"],
			strings.ToLower(quote.ID) + "_query_timestamp": quoteData["QueryTimestamp"],
		}

		for k, v := range quoteTmp {
			jsonData[k] = v
		}
	}
	/*jsonString, _ := json2.Marshal(jsonData)

	fmt.Println(string(jsonString))*/

	return &jsonData
}

func build_gdax_dataset() {
	currencies := gC.get_coins()

	for _, currency := range *currencies {
		update_coin_data(currency.ID, currency.DisplayName)
	}
}

func update_coin_data(coinId string, coinName string) {
	glog.V(2).Infoln("update_coin_data " + coinId)

	if is_valid_coin(coinId) {

		gdaxDataSet.Lock()
		defer gdaxDataSet.Unlock()

		coin := gdaxDataSet.Coin

		for _, product := range gP.Products {
			//only want products with the specified coin: coinId=BTC productId=BTC-USD/BTC-EUR
			if strings.HasPrefix(product.Id, coinId) {

				stats := gS.get_product_stats(product.Id)
				if stats == nil {
					glog.Warningln("Failed to retriece stats for product: " + product.Id)
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
				coinData := map[string]interface{}{
					"DisplayName": coinName,
					product.QuoteCurrency: map[string]interface{}{
						"Price":          stats.Last,
						"Delta":          delta,
						"QueryTimestamp": stats.QueryTimestamp,
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

func Get_Coins() *[]coin_struct.Coin {
	return gC.get_coins()
}

func Get_Coin_Stats(coinId string) (*map[string]interface{}, error) {
	if !is_valid_coin(coinId) {
		err := errors.New("invalid coinId id: " + coinId)
		return nil, err
	}

	/*if is_data_old(coinId, 10) {
		currency, _ := gC.get_currency(coinId)
		update_coin_data(coinId, currency.DisplayName)
	}*/

	jsonData := build_json_struct(coinId)

	return jsonData, nil
}

func Update_Data(force bool) {
	gC.update_data(force)
	gP.update_data(force)
	gS.update_data(force)
}