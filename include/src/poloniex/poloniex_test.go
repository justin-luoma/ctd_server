package poloniex

import (
	"flag"
	"fmt"
	"testing"
	json2 "encoding/json"
)

func init() {
	if testing.Verbose() {
		flag.Set("v", "2")
		flag.Set("logtostderr", "true")
		flag.Set("stderrthreshold", "INFO")
	}
}

/*func TestGetTicker(t *testing.T) {
	data, err := get_ticker()
	if err != nil {
		glog.Errorln(err)
	}
	tmp := data["ETH_ZEC"].(map[string]interface{})["last"]
	fmt.Println(tmp)

jsonData, err := json2.MarshalIndent(data, "", " ")
	if err != nil {
		glog.Errorln(err)
	}

	fmt.Println(string(jsonData))

}*/

/*
func TestGetCurrencies(t *testing.T) {
	data, err := get_currencies()
	if err != nil {
		glog.Errorln(err)
	}
	tmp := data["BTC"]

	fmt.Println(tmp)
}*/

func TestPoloniex(t *testing.T)  {
	//TestInit()
	/*fmt.Printf("%s: %t\n", "BTC", is_valid_coin("BTC"))
	fmt.Printf("%s: %t\n", "LTC", is_valid_coin("LTC"))
	fmt.Printf("%s: %t\n", "ETH", is_valid_coin("ETH"))
	fmt.Printf("%s: %t\n", "VTC", is_valid_coin("VTC"))*/

	/*for i := 1; i <= 10; i++ {
		fmt.Printf("%s: %t\n", "Is data old", is_data_old("BTC", 5))
		time.Sleep(time.Second)
	}*/

	/*jsonData, err := Get_Coin_Stats("BTC")
	if err != nil {
		fmt.Println(err)
	}

	out, err := json2.MarshalIndent(jsonData, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))*/
	/*poloniexCurrencies, err := get_currencies()
	if err != nil {
		glog.Errorln(err)
		//return nil, err
	}
	for _, data := range poloniexCurrencies {
		data := data.(map[string]interface{})
		fmt.Printf("%T", data["name"])
	}*/
	coins, err := Get_Coins()
	if err != nil {
		fmt.Println(err)
	}

	out, err := json2.MarshalIndent(coins, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
}