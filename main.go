package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"time"
)

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
type Coin struct {
	ID             string
	PriceBTC       float32 `json:"price_btc"`
	PriceETH       float32 `json:"price_eth"`
	PriceUSD       float32 `json:"price_usd"`
	QueryTimeStamp int64   `json:"query_timestamp"`
}

//var people []Person

// our main function
func main() {
	/*people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})*/

	router := mux.NewRouter()
	/*router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")*/
	router.HandleFunc("/coin/{id}", GetCoin).Methods("GET")

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
func GetCoin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	url := fmt.Sprintf("http://coincap.io/page/%s", params["id"])

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		coin := Coin{QueryTimeStamp:time.Now().Unix()}
		err := json.NewDecoder(response.Body).Decode(&coin)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(coin)
		err = json.NewEncoder(w).Encode(coin)
		if err != nil{
			log.Println(err)
		}
	}
}