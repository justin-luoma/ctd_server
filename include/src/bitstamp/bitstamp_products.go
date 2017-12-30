package bitstamp

import (
	json2 "encoding/json"
	"restful_query"
	"github.com/golang/glog"
)

type BitstampProducts struct {
	BaseDecimals    int    `json:"base_decimals"`
	MinimumOrder    string `json:"minimum_order"`
	Name            string `json:"name"`
	CounterDecimals int    `json:"counter_decimals"`
	Trading         string `json:"trading"`
	URLSymbol       string `json:"url_symbol"`
	Description     string `json:"description"`
}
type BitstampProduct struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
}

//noinspection ALL
func get_products() ([]BitstampProducts, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "v2/trading-pairs-info/")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	var products []BitstampProducts
	json2.Unmarshal(bodyBytes, &products)

	return products, nil
}

//noinspection ALL
func get_online_products() ([]BitstampProducts, []BitstampProducts, error) {
	products, err := get_products()
	if err != nil {
		glog.Errorln(err)
		return nil, nil, err
	}

	var onlineProducts []BitstampProducts
	var offlineProducts []BitstampProducts

	//This will need to be fixed if we find a way to check if a coin is "online"

	/*
	for _, product := range products {
		if product.Status == "online" {
			onlineProducts = append(onlineProducts, product)
			onlineProductIds = append(onlineProductIds, product.Id)
		} else {
			offlineProducts = append(offlineProducts, product)
			glog.Warningln("skipped Bitstamp product: " + product.Id + " status: " + product.Status)
		}
	}*/
	onlineProducts = products
	return onlineProducts, offlineProducts, nil
}
