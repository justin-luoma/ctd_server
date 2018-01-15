package bittrex

import (
	"coin_struct"
	"errors"
	"exchange_api_status"
	"time"

	"github.com/golang/glog"
)

const currencyDataOldDuration = 10 * time.Minute

type bittrexCurrencies struct {
 	Currencies []coin_struct.Coin
 	queryTimestamp int64
 }

 func init_currencies() (bC *bittrexCurrencies) {
 	bC.update_currencies()

 	return
 }

 func (bC *bittrexCurrencies) update_currencies() {
 	var tCoin coin_struct.Coin
 	var tCoins coin_struct.Coin
 	var timestamp = time.Now().Unix()

 	currencies, err := b.GetCurrencies()
 	if err != nil {
 		exchange_api_status.Update_Status("bittrex", 0)

 		return
	}

	for _, currency := range currencies {
		tCoin = coin_struct.Coin{
			ID: currency.Currency,
			DisplayName: currency.CurrencyLong,
			IsActive: currency.IsActive,
			StatusMessage: currency.Notice
		}
		tCoins = append(tCoins, tCoin)
	}

	exchange_api_status.Update_Status("bittrex", 1)
	 *bC = bittrexCurrencies{
	 	Currencies: tCoins,
	 	queryTimestamp: timestamp,
	 }
 }

 func (bC *bittrexCurrencies) update_data(force bool) {
 	if force {
		glog.Infoln("force update bittrex currencies called")
		bC.update_currencies()
		return
	} else {
		dataAge := time.Since(time.Unix(bC.queryTimestamp, 0))
		if dataAge > currencyDataOldDuration {
			glog.Infoln("bittrex currency data old, updating")
			bC.update_currencies()
		}
	}
 }

 func (bC *bittrexCurrencies) get_coins() (coins *[]coin_struct.Coin) {
 	for _, currency := range bC.Currencies {
 		if currency.IsActive {
 			*coins = append(*coins, currency)
		}
	}

	return
 }

 func (bC *bittrexCurrencies) get_currencies() *[]coin_struct.Coin {
 	return &bC.Currencies
 }

 func (bC *bittrexCurrencies) get_currency(id string) (*coin_struct.Coin, error) {
 	for _, currency := range bC.Currencies {
 		if id == currency.ID {
 			return &currency, nil
		}
	}
	err := errors.New("invalid coin id: " + id)

	return nil, err
 }
