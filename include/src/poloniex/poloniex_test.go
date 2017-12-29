package main

import (
	"flag"
	"testing"

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
	TestInit()
}