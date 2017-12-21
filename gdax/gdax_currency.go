package gdax

import (
	"../restful_query"
	json2 "encoding/json"
	"github.com/golang/glog"
)

type GdaxCurrencies struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	MinSize string `json:"min_size,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

//noinspection ALL
func get_currency(id string) (*GdaxCurrencies, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies/" + id)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	currency := GdaxCurrencies{}
	json2.Unmarshal(bodyBytes, &currency)

	return &currency, nil
}

func get_currencies() (*[]GdaxCurrencies, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	var currencies []GdaxCurrencies
	json2.Unmarshal(bodyBytes, &currencies)

	return &currencies, nil
}
