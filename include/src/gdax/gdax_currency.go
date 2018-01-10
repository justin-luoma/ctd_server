package gdax

import (
	"coin_struct"
	json2 "encoding/json"
	"errors"
	"exchange_api_status"
	"restful_query"
	"time"

	"github.com/golang/glog"
)

const dataOldDuration = 10 * time.Minute

var currencyTypes = map[string]string{
	"JPY": "fiat",
	"CAD": "fiat",
	"USD": "fiat",
	"EUR": "fiat",
	"GBP": "fiat",
}

type gdaxCurrencyStruct struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	MinSize string `json:"min_size,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

//noinspection ALL
func get_currency(id string) (*gdaxCurrencyStruct, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies/" + id)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	currency := gdaxCurrencyStruct{}
	if err := json2.Unmarshal(bodyBytes, &currency); err != nil {
		return &currency, err
	}

	return &currency, nil
}

func get_currencies_api() (*[]gdaxCurrencyStruct, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	var currencies []gdaxCurrencyStruct
	if err := json2.Unmarshal(bodyBytes, &currencies); err != nil {
		return &currencies, err
	}

	return &currencies, nil
}

type gdaxCurrencies struct {
	Currencies     []coin_struct.Coin
	queryTimestamp int64
}

func update_coins() *[]coin_struct.Coin {

	var tCoin coin_struct.Coin
	var tCoins []coin_struct.Coin

	currencies, err := get_currencies_api()
	if err != nil || len(*currencies) < 1 {
		tCoins = []coin_struct.Coin{}

		exchange_api_status.Update_Status("gdax", 0)
		return &tCoins
	}

	for _, currency := range *currencies {
		var isActive bool
		if currency.Status == "online" {
			isActive = true
		} else {
			isActive = false
		}
		var isFiat = false
		if _, ok := currencyTypes[currency.Id]; ok {
			isFiat = true
		}
		var statusMessage = ""
		if currency.Message != "" {
			statusMessage = currency.Message
		}

		tCoin = coin_struct.Coin{
			ID:            currency.Id,
			DisplayName:   currency.Name,
			IsActive:      isActive,
			IsFiat:        isFiat,
			StatusMessage: statusMessage,
		}

		tCoins = append(tCoins, tCoin)
	}

	exchange_api_status.Update_Status("gdax", 1)

	return &tCoins
}

func init_currencies() *gdaxCurrencies {
	var gC = gdaxCurrencies{
		queryTimestamp: time.Now().Unix(),
	}
	gC.Currencies = *update_coins()

	return &gC
}

func (gC *gdaxCurrencies) Test() *[]coin_struct.Coin {

	return &gC.Currencies
}

func (gC *gdaxCurrencies) update_data(force bool) {
	if force {
		gC.Currencies = *update_coins()

		return
	} else {
		dataAge := time.Since(time.Unix(gC.queryTimestamp, 0))
		if dataAge > dataOldDuration {
			gC.Currencies = *update_coins()
		}
	}
}

func (gC *gdaxCurrencies) get_coins() *[]coin_struct.Coin {
	var tmpCoin coin_struct.Coin
	var tmpCoins []coin_struct.Coin
	for _, currency := range gC.Currencies {
		if currency.IsActive && !currency.IsFiat {
			tmpCoin = currency
			tmpCoins = append(tmpCoins, tmpCoin)
		}
	}

	return &tmpCoins
}

func (gC *gdaxCurrencies) get_currencies() *[]coin_struct.Coin {
	return &gC.Currencies
}

func (gC *gdaxCurrencies) get_currency(id string) (*coin_struct.Coin, error) {
	for _, currency := range gC.Currencies {
		if id == currency.ID {
			return &currency, nil
		}
	}
	err := errors.New("invalid coin id: " + id)

	return nil, err
}