package gdax

import (
	"../coin_struct"
	"../exchange_api_status"
	"../restful_query"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/jinzhu/copier"
)

var apiUrl string = "https://api.gdax.com/"
var supportedCoins = [3]string{"BTC","ETH","LTC"}
var productsIds []string

var gdaxDataset = struct {
	sync.RWMutex
	Coin map[string]*coin_struct.Coin
}{Coin: make(map[string]*coin_struct.Coin)}

type GdaxCurrencies struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	MinSize string `json:"min_size,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type GdaxProducts struct {
	Id            string `json:"id,omitempty"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	Status        string `json:"status,omitempty"`
}

type GdaxStatsQuery struct {
	Product        string `json:"id"`
	Open           string `json:"open"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Volume         string `json:"volume"`
	Last           string `json:"last"`
	Volume30Day    string `json:"volume_30day"`
	QueryTimeStamp int64  `json:"query_timestamp"`
}

type GdaxStats struct {
	Product        string  `json:"id"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Volume         float64 `json:"volume"`
	Last           float64 `json:"last"`
	Volume30Day    float64 `json:"volume_30day"`
	QueryTimeStamp int64   `json:"query_timestamp"`
}

func init() {
	flag.Parse()
	build_gdax_dataset()
}

func get_currency(id string) (*GdaxCurrencies, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies/" + id)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	currencies := GdaxCurrencies{}
	json.Unmarshal(bodyBytes, &currencies)

	return &currencies, nil
}

func get_products() ([]GdaxProducts, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "products")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	var products []GdaxProducts
	json.Unmarshal(bodyBytes, &products)

	return products, nil
}

func get_online_products() ([]GdaxProducts, []GdaxProducts, error) {
	products, err := get_products()
	if err != nil {
		glog.Errorln(err)
		return nil, nil, err
	}

	var onlineProducts []GdaxProducts
	var offlineProducts []GdaxProducts

	for _, product := range products {
		if product.Status == "online" {
			onlineProducts = append(onlineProducts, product)
			productsIds = append(productsIds, product.Id)
		} else {
			offlineProducts = append(offlineProducts, product)
			glog.Warningln("skipped GDAX product: " + product.Id + " status: " + product.Status)
		}
	}

	return onlineProducts, offlineProducts, nil
}

func get_stats() ([]GdaxStats, error) {
	onlineProducts, _, err := get_online_products()
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	var stats []GdaxStats

	for _, product := range onlineProducts {
		bodyBytes, err := restful_query.Get(apiUrl + "products/" + product.Id + "/stats")
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		var queryStat = GdaxStatsQuery{
			Product:        product.Id,
			QueryTimeStamp: time.Now().Unix(),
		}
		json.Unmarshal(bodyBytes, &queryStat)

		var floatStat GdaxStats
		convert_stats(&queryStat, &floatStat)
		stats = append(stats, floatStat)
	}

	return stats, nil
}

func get_product_stats(productId string) (GdaxStats, error) {
	var productQStat = GdaxStatsQuery{
		Product:        productId,
		QueryTimeStamp: time.Now().Unix()}

	var productStats GdaxStats

	bodyBytes, err := restful_query.Get(apiUrl + "products/" + productId + "/stats")
	if err != nil {
		glog.Errorln(err)
		return productStats, err
	}
	json.Unmarshal(bodyBytes, &productQStat)

	convert_stats(&productQStat, &productStats)

	return productStats, nil
}

func convert_stats(stringStatsStruct *GdaxStatsQuery, floatStatsStruct *GdaxStats) {
	floatStatsStruct.Product = stringStatsStruct.Product
	floatStatsStruct.QueryTimeStamp = stringStatsStruct.QueryTimeStamp
	fOpen, errO := strconv.ParseFloat(stringStatsStruct.Open, 64)
	fLast, errL := strconv.ParseFloat(stringStatsStruct.Last, 64)
	fHigh, errH := strconv.ParseFloat(stringStatsStruct.High, 64)
	fVolume, errV := strconv.ParseFloat(stringStatsStruct.Volume, 64)
	f30Day, err3 := strconv.ParseFloat(stringStatsStruct.Volume30Day, 64)
	fLow, errLow := strconv.ParseFloat(stringStatsStruct.Low, 64)
	errs := []error{errO, errL, errH, errV, err3, errLow}
	for _, err := range errs {
		if err != nil {
			glog.Errorln("Unable to convert GDAX strings to floats: ", err)
		}
	}
	floatStatsStruct.Open = fOpen
	floatStatsStruct.Last = fLast
	floatStatsStruct.High = fHigh
	floatStatsStruct.Volume = fVolume
	floatStatsStruct.Volume30Day = f30Day
	floatStatsStruct.Low = fLow
}

func Init() {
	glog.V(2).Infoln("GDAX Init called")
	/*stats, err := get_product_stats("BTC-USD")
	if err != nil {
		glog.Errorln(err)
	}
	fmt.Println(stats)*/
}

func build_gdax_dataset() {
	/*gdaxDataset.Lock()
	defer gdaxDataset.Unlock()*/

	uP, _, err := get_online_products()
	if err != nil || uP == nil {
		glog.Errorln("Unable to initialize GDAX package, check error log")
		exchange_api_status.Update_Status("gdax", 0)
	} else {
		exchange_api_status.Update_Status("gdax", 1)

		for _, coin := range supportedCoins {
			update_coin_data(coin)
		}
	}
}

func update_coin_data(coinId string) {
	glog.V(2).Infoln("update_coin_data " + coinId)
	var productsToUpdate []string
	for _, productsId := range productsIds {
		if strings.HasPrefix(productsId, coinId) {
			productsToUpdate = append(productsToUpdate, productsId)
		}
	}

	var btcCoin = &coin_struct.Coin{}
	var ethCoin = &coin_struct.Coin{}
	var ltcCoin = &coin_struct.Coin{}

	for _, productId := range productsToUpdate {
		baseCurrency := strings.Split(productId, "-")[0]
		quoteCurrency := strings.Split(productId, "-")[1]
		stats, err := get_product_stats(productId)
		var queryTime = time.Now().Unix()
		currency, err := get_currency(baseCurrency)
		if err != nil || currency == nil {
			glog.Warningln("Unable to get currency info for: " + baseCurrency)
			break
		}

		// THIS MUST BE UPDATED IF GDAX ADDS COINS
		switch baseCurrency {
		case "BTC":
			btcCoin.ID = "BTC"
			btcCoin.PriceBtc = 1
			btcCoin.QueryTimeStamp = queryTime
			btcCoin.DisplayName = currency.Name
			switch quoteCurrency {
			case "USD":
				btcCoin.PriceUsd = stats.Last
			case "EUR":
				btcCoin.PriceEur = stats.Last
			case "GBP":
				btcCoin.PriceGbp = stats.Last
			}
			fmt.Println(btcCoin)
		case "ETH":
			ethCoin.ID = "ETH"
			ethCoin.PriceEth = 1
			ethCoin.QueryTimeStamp = queryTime
			ethCoin.DisplayName = "Ethereum"
			switch quoteCurrency {
			case "USD":
				ethCoin.PriceUsd = stats.Last
			case "EUR":
				ethCoin.PriceEur = stats.Last
			case "BTC":
				ethCoin.PriceBtc = stats.Last
			}
		case "LTC":
			ltcCoin.ID = "LTC"
			ltcCoin.PriceLtc = 1
			ltcCoin.QueryTimeStamp = queryTime
			ltcCoin.DisplayName = currency.Name
			switch quoteCurrency {
			case "USD":
				ltcCoin.PriceUsd = stats.Last
			case "EUR":
				ltcCoin.PriceEur = stats.Last
			case "BTC":
				ltcCoin.PriceBtc = stats.Last
			}
		}

	}


	gdaxDataset.Lock()
	defer gdaxDataset.Unlock()

	coin := gdaxDataset.Coin

	switch coinId {
	case "BTC":
		coin["BTC"] = btcCoin
	case "ETH":
		coin["ETH"] = ethCoin
	case "LTC":
		coin["LTC"] = ltcCoin
	}
}

func Get_Coin_Stats(coin string) (coin_struct.Coin, error) {
	//check to see if current data is old
	gdaxDataset.RLock()
	dataAge := time.Since(time.Unix(gdaxDataset.Coin[coin].QueryTimeStamp, 0)).Seconds()
	gdaxDataset.RUnlock()
	if dataAge >= 5 {
		update_coin_data(coin)
	}
	var coinData coin_struct.Coin
	gdaxDataset.RLock()
	copier.Copy(&coinData, gdaxDataset.Coin[coin])
	gdaxDataset.RUnlock()

	return coinData, nil
}
