package gdax

import (
	json2 "encoding/json"
	"flag"
	"fmt"
	"testing"
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

// func TestCurrencies(t *testing.T) {
// 	gC := init_currencies()
// 	coins := gC.Test()

// 	jsonData, _ := json2.MarshalIndent(*coins, "", " ")
// 	fmt.Println(string(jsonData))
// }

func TestProducts(t *testing.T) {
	gP := init_products()
	products := gP.Test()

	jsonData, _ := json2.MarshalIndent(*products, "", " ")
	fmt.Println(string(jsonData))
}
