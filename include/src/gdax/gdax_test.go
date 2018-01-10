package gdax

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

/*
func TestPullCurrencies(t *testing.T) {
	get_stats()
}
*/

func TestCurrencies(t *testing.T) {
	gC := init_currencies()
	coins := gC.Test()

	jsonData, _ := json2.MarshalIndent(*coins, "", " ")
	fmt.Println(string(jsonData))
}