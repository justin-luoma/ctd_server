package gdax

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

func TestPullCurrencies(t *testing.T) {
	get_stats()
}
