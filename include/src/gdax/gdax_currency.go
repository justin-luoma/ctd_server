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

const currencyDataOldDuration = 10 * time.Minute

var currencyTypes = map[string]string{
	"JPY": "fiat",
	"CAD": "fiat",
	"USD": "fiat",
	"EUR": "fiat",
	"GBP": "fiat",
}

type gdaxCurrencyApi struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	MinSize string `json:"min_size,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

//noinspection ALL
func get_currency(id string) (*gdaxCurrencyApi, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies/" + id)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	currency := gdaxCurrencyApi{}
	if err := json2.Unmarshal(bodyBytes, &currency); err != nil {
		return &currency, err
	}

	return &currency, nil
}

func get_currencies_api() (*[]gdaxCurrencyApi, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "currencies")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	var currencies []gdaxCurrencyApi
	if err := json2.Unmarshal(bodyBytes, &currencies); err != nil {
		return &currencies, err
	}

	return &currencies, nil
}

type gdaxCurrencies struct {
	Currencies     []coin_struct.Coin
	queryTimestamp int64
}

func init_currencies() *gdaxCurrencies {
	var gC = gdaxCurrencies{}
	gC.update_coins()

	return &gC
}

func (gC *gdaxCurrencies) update_coins() {
	var timestamp = time.Now().Unix()
	var tCoin coin_struct.Coin
	var tCoins []coin_struct.Coin

	currencies, err := get_currencies_api()
	if err != nil || len(*currencies) < 1 {
		tCoins = []coin_struct.Coin{}

		exchange_api_status.Update_Status("gdax", 0)
		return
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

	*gC = gdaxCurrencies{
		Currencies:     tCoins,
		queryTimestamp: timestamp,
	}
}

func (gC *gdaxCurrencies) Test() *[]coin_struct.Coin {

	return &gC.Currencies
}

func (gC *gdaxCurrencies) update_data(force bool) {
	if force {
		gC.update_coins()

		return
	} else {
		dataAge := time.Since(time.Unix(gC.queryTimestamp, 0))
		if dataAge > currencyDataOldDuration {
			gC.update_coins()
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