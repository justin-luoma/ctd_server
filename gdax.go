package main

import (
	"encoding/json"
	"math/big"
)

var apiUrl string = "https://api.gdax.com/"

type gdaxCurrency struct {
	id string `json:"id,omitempty"`
	name string `json:"name,omitempty"`
	min_size float32 `json:"min_size,omitempty"`
	status string `json:"status,omitempty"`
	message string `json"message,omitempty"`
}

type gdaxProducts struct {
	id string `json:"id,omitempty"`
	status string `json:"status,omitempty"`
}

type gdaxStats struct {
	last big.Float `json:"last,omitempty"`
	timestamp string
}

