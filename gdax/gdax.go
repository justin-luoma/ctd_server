package gdax

import (
	"../restful_query"
	"encoding/json"
	"errors"
	"flag"
	"github.com/golang/glog"
	"time"
	"../coin_struct"
)

var apiUrl string = "https://api.gdax.com/"
var supportedCoins [] string

func init() {
	flag.Parse()
	currencies, err := get_currencies()
	if err != nil {
		glog.Errorln("Failed initialize GDAX package")
	}
}

type GdaxCurrencies struct {
	Id      string  `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	MinSize float32 `json:"min_size,omitempty"`
	Status  string  `json:"status,omitempty"`
	Message string  `json:"message,omitempty"`
}

type GdaxProducts struct {
	Id     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

type GdaxStats struct {
	Product        string `json:"id"`
	Open           string `json:"open"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Volume         string `json:"volume"`
	Last           string `json:"last"`
	Volume30Day    string `json:"volume_30day"`
	QueryTimeStamp int64  `json:"query_timestamp"`
}

func get_currencies() ([]GdaxCurrencies, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	var currencies []GdaxCurrencies
	json.Unmarshal(bodyBytes, &currencies)

	return currencies, nil
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

func

func get_stats() ([]GdaxStats, error) {
	products, err := get_products()
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	var stats []GdaxStats
	var productErr error = nil

	for _, product := range products {
		if product.Status == "online" {
			bodyBytes, err := restful_query.Get(apiUrl + "products/" + product.Id + "/stats")
			if err != nil {
				glog.Errorln(err)
				return nil, err
			}
			stat := GdaxStats{Product: product.Id,
				QueryTimeStamp: time.Now().Unix()}
			json.Unmarshal(bodyBytes, &stat)
			stats = append(stats, stat)
		} else {
			glog.Warningln("skipped GDAX product: " + product.Id + " status: " + product.Status)
			productErr = errors.New("product: " + product.Id + " status: " + product.Status)
		}
	}

	if productErr == nil {
		return stats, nil
	} else {
		return stats, productErr
	}
}

func Get_Coin_Stats(id string) coin_struct.Coin {
	gdaxCoin := coin_struct.Coin{}

}
