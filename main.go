package main

import (
	"bitstamp"
	"bittrex"
	"coin_struct"
	"coincap"
	"encoding/json"
	"exchange_api_status"
	"flag"
	"gdax"
	"log"
	"net/http"
	_ "net/http/pprof"
	"poloniex"
	"scheduler"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var flagEnableProfile = flag.Bool("profile", false, "enable profiling with pprof at localhost:6060")

func init() {
	flag.Parse()
	if *flagEnableProfile {
		flag.Set("v", "10")
		flag.Set("logtostderr", "true")
		flag.Set("stderrthreshold", "INFO")
	}
	flag.Parse()
}

func main() {
	exchange_api_status.Init()

	bitstamp.Init()
	bittrex.Init()
	gdax.Init()
	poloniex.Init()



	exchange_api_status.Start_Exchange_Monitoring()
	scheduler.Start_Scheduler()

	router := mux.NewRouter()
	router.HandleFunc("/coin/{id}", GetCoin).
		Methods("GET")

	router.StrictSlash(true)

	router.HandleFunc("/{exchange}/coins", get_exchange_coins).
		Methods("GET")

	router.HandleFunc("/{exchange}/coin/{id}", get_exchange_coin).
		Methods("GET")

	if *flagEnableProfile {
		go http.ListenAndServe(":6060", nil)
	}

	log.Fatal(http.ListenAndServe(":8000", router))
}

func get_exchange_coins(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var coins []coin_struct.Coin
	var err error
	switch params["exchange"] {
	case "gdax":
		coins = *gdax.Get_Coins()
	case "poloniex":
		coins, err = poloniex.Get_Coins()
	case "bittrex":
		coins, err = bittrex.Get_Coins()
	case "bitstamp":
		coins, err = bitstamp.Get_Coins()
	default:
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err != nil {
		glog.Errorln(err)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(coins)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func get_exchange_coin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var jsonData *map[string]interface{}
	var err error
	switch params["exchange"] {
	case "gdax":
		jsonData, err = gdax.Get_Coin_Stats(strings.ToUpper(params["id"]))
	case "poloniex":
		jsonData, err = poloniex.Get_Coin_Stats(strings.ToUpper(params["id"]))
	case "bittrex":
		jsonData, err = bittrex.Get_Coin_Stats(strings.ToUpper(params["id"]))
	case "bitstamp":
		jsonData, err = bitstamp.Get_Coin_Stats(strings.ToUpper(params["id"]))
	}
	if err != nil {
		glog.Errorln(err)
		if strings.HasPrefix(err.Error(), "invalid coinId id:") {
			http.Error(w, error.Error(err), http.StatusBadRequest)
			return
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jsonData)
	if err != nil {
		glog.Errorln(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}

func GetCoin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	coin := coin_struct.Coin{}
	coin = coincap.GetCoinCapCoin(params["id"])
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(coin)
	if err != nil {
		log.Println(err)
	}
}
