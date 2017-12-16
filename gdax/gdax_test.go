package gdax

import (
	"testing"
	"fmt"
	"log"
)

func TestPullCurrencies(t *testing.T) {
	rt, err := pull_currencies()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rt)
}
