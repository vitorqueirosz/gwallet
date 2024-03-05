package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Currency struct {
	Name  string
	Code  string
	Price string
}

type currencyController struct {
	currencies []Currency
}

func CurrencyController() *currencyController {
	return &currencyController{}
}

func (c *currencyController) createCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies := []Currency{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&currencies)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding request body - %v", err))
		return
	}

	respondWithJSON(w, 201, currencies)
}
