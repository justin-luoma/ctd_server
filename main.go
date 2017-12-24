package main

import (
	"./coin_struct"
	"./coincap"
	"./exchange_api_status"
	"./gdax"
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"./bittrex"
)

func init() {
	flag.Parse()
}

/*type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}*/
/*type Coin struct {
	ID             string  `json:"id"`
	DisplayName    string  `json:"display_name"`
	Cap24HrChange  float64 `json:"cap24hrChange"`
	PriceBtc       float64 `json:"price_btc"`
	PriceEth       float64 `json:"price_eth"`
	PriceUsd       float64 `json:"price_usd"`
	QueryTimeStamp int64   `json:"query_timestamp"`
}*/

//var people []Person

// our main function
func main() {
	exchange_api_status.Start_Exchange_Monitoring()
	/*people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})*/

	router := mux.NewRouter()
	/*router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")*/
	router.HandleFunc("/coin/{id}", GetCoin).Methods("GET")

	router.HandleFunc("/gdax/coins", get_gdax_coins).Methods("GET")
	router.HandleFunc("/gdax/coin/{id}", get_gdax_coin).Methods("GET")

	router.HandleFunc("/bittrex/coins", get_bittrex_coins).Methods("GET")
	router.HandleFunc("/bittrex/coin/{id}", get_bittrex_coin).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

/*
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}*/

func get_gdax_coins(w http.ResponseWriter, r *http.Request) {
	coins, err := gdax.Get_Coins()
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

func get_gdax_coin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	jsonData, err := gdax.Get_Coin_Stats(strings.ToUpper(params["id"]))
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

func get_bittrex_coins(w http.ResponseWriter, r *http.Request) {
	coins, err := bittrex.Get_Coins()
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

func get_bittrex_coin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	jsonData, err := bittrex.Get_Coin_Stats(strings.ToUpper(params["id"]))
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

	/*
		url := fmt.Sprintf("http://coincap.io/page/%s", params["id"])

		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			coin := Coin{QueryTimeStamp: time.Now().Unix()}
			err := json.NewDecoder(response.Body).Decode(&coin)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(coin)
			if err != nil {
				log.Println(err)
			}
		}*/
}
