package gdax

import (
	json2 "encoding/json"
	"../restful_query"
	"github.com/golang/glog"
)

type GdaxProducts struct {
	Id            string `json:"id,omitempty"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	Status        string `json:"status,omitempty"`
}

//noinspection ALL
func get_products() ([]GdaxProducts, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "products")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	var products []GdaxProducts
	json2.Unmarshal(bodyBytes, &products)

	return products, nil
}

//noinspection ALL
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
			onlineProductIds = append(onlineProductIds, product.Id)
		} else {
			offlineProducts = append(offlineProducts, product)
			glog.Warningln("skipped GDAX product: " + product.Id + " status: " + product.Status)
		}
	}

	return onlineProducts, offlineProducts, nil
}
