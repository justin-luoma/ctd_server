package restful_query

import (
	"fmt"
	"log"
	"testing"
)

func TestQuery(t *testing.T) {
	body, err := Get("https://api.gdax.com/LTC-EUR/stats")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(body)
	body, err = Get("https://api.gdax.com/BTC-EUR/stats")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(body)
}
