package gdax

import (
	"../coin_struct"
	"../decimal_math"
	"../exchange_api_status"
	"errors"
	"flag"
	"github.com/golang/glog"
	"strings"
	"sync"
	"time"
)

const apiUrl = "https://api.gdax.com/"

var currencyTypes = map[string]string{
	"BTC": "crypto",
	"BCH": "crypto",
	"ETH": "crypto",
	"LTC": "crypto",
	"USD": "fiat",
	"EUR": "fiat",
	"GBP": "fiat",
}

var onlineProductIds []string

var gdaxDataSet = struct {
	sync.RWMutex
	Coin map[string]interface{}
}{Coin: make(map[string]interface{})}

func init() {
	flag.Parse()
	build_gdax_dataset()
}

func check_online_status() bool {
	_, _, err := get_online_products()
	if err != nil {
		return false
	}

	return true
}

func is_valid_coin(coinId string) bool {
	if _, ok := currencyTypes[coinId]; ok && currencyTypes[coinId] == "crypto" {
		return true
	} else {
		return false
	}
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

	for quote, _ := range currencyTypes {
		if !valid_product_stats(coinId, quote) {
			continue
		}
		dataAge := time.Since(
			time.Unix(coin[quote].(map[string]interface{})[quote+"QueryTimestamp"].(int64), 0)).Seconds()
		if dataAge > 10 {
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
	for quote, _ := range currencyTypes {
		if !valid_product_stats(coinId, quote) {
			continue
		}
		/*
			quoteData hold the data for the current quote currency in the loop,
			while quoteTmp is hold the structure for the json formatted data
		*/
		quoteData := gdaxData[quote].(map[string]interface{})
		quoteTmp := map[string]interface{}{
			strings.ToLower(quote) + "_price":           quoteData["Price"+quote],
			strings.ToLower(quote) + "_24_hour_change":  quoteData["Delta"+quote],
			strings.ToLower(quote) + "_query_timestamp": quoteData[quote+"QueryTimestamp"],
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
func build_gdax_dataset() {
	uP, _, err := get_online_products()
	if err != nil && uP == nil {
		glog.Errorln("Unable to initialize GDAX package, check error log")
		exchange_api_status.Update_Status("gdax", 0)
		return
	} else {
		exchange_api_status.Update_Status("gdax", 1)

		currencies, err := get_currencies()
		if err != nil {
			glog.Errorln("Unable to get GDAX currencies")
			exchange_api_status.Update_Status("gdax", 0)
			return
		}

		for _, currency := range *currencies {
			//we only want to
			if currency.Status == "online" && currencyTypes[currency.Id] == "crypto" {
				update_coin_data(currency.Id, currencies, &uP)
			}
		}

		//TEST SECTION//
		//_ = build_json_struct("BTC")
	}
}

//noinspection ALL
func update_coin_data(coinId string, currencies *[]GdaxCurrencies, onlineProducts *[]GdaxProducts) {
	glog.V(2).Infoln("update_coin_data " + coinId)

	if is_valid_coin(coinId) {

		gdaxDataSet.Lock()
		defer gdaxDataSet.Unlock()

		coin := gdaxDataSet.Coin

		currencyNames := make(map[string]string)

		for _, currency := range *currencies {
			currencyNames[currency.Id] = currency.Name
		}

		for _, product := range *onlineProducts {
			//only want products with the specified coin: coinId=BTC productId=BTC-USD/BTC-EUR
			if strings.HasPrefix(product.Id, coinId) {

				stats, err := get_product_stats(product.Id)
				if err != nil {
					glog.Warningln("Failed to retriece stats for product: " + product.Id)
					continue
				}
				delta := decimal_math.Calculate_Percent_Change_Float(stats.Open, stats.Last)
				/*
					build the structure for the coin:
					{
					 "DisplayName": "Bitcoin",
					 "USD": {
					  "DeltaUSD": -5.63,
					  "PriceUSD": 16336.62,
					  "USDQueryTimestamp": 1513825686
					 }
					}
				*/
				coinData := map[string]interface{}{
					"DisplayName": currencyNames[coinId],
					product.QuoteCurrency: map[string]interface{}{
						"Price" + product.QuoteCurrency:          stats.Last,
						"Delta" + product.QuoteCurrency:          delta,
						product.QuoteCurrency + "QueryTimestamp": stats.QueryTimeStamp,
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

	for c, _ := range currencyTypes {
		var currency, err = get_currency(c)
		if err != nil {
			glog.Error(err)
			glog.Error("invalid response for GDAX currency: " + c)
		} else {
			coin.ID = currency.Id
			coin.DisplayName = currency.Name
			coins = append(coins, coin)
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
			glog.Errorln("GDAX package offline, check error log")
			exchange_api_status.Update_Status("gdax", 0)
			return nil, errors.New("GDAX API is down")
		} else {
			exchange_api_status.Update_Status("gdax", 1)

			currencies, err := get_currencies()
			if err != nil {
				glog.Errorln("Unable to get GDAX currencies")
				exchange_api_status.Update_Status("gdax", 0)
				return nil, errors.New("GDAX API is down")
			}
			update_coin_data(coinId, currencies, &uP)
		}
	}

	jsonData := build_json_struct(coinId)

	return jsonData, nil
}
