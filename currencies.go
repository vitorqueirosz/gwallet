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

func (c *currencyController) handleCreateCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies := []Currency{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&currencies)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding request body - %v", err))
		return
	}

	c.currencies = currencies

	// for _, c := range currencies {
	// 	parsedStr := strings.ReplaceAll(c.Price, ",", "")
	// 	fmt.Println(parsedStr)
	// 	fmt.Println(strconv.ParseFloat(parsedStr, 64))
	// }

	respondWithJSON(w, 201, currencies)
}

func (c *currencyController) handleGetCurrencies(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 201, c.currencies)
}
