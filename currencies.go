package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vitorqueirosz/gwallet/internal/database"
)

type Currency struct {
	Name  string
	Code  string
	Price string
}

func (apiCfg *apiConfig) handleCreateCurrencies(w http.ResponseWriter, r *http.Request) {
	cb := []Currency{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cb)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding request body - %v", err))
		return
	}

	cc := []database.Currency{}
	for _, c := range cb {
		currency, err := apiCfg.DB.CreateCurrencies(r.Context(), database.CreateCurrenciesParams{
			ID:        uuid.New(),
			Name:      c.Name,
			Code:      c.Code,
			Price:     c.Price,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error creating currency - %v", err))
			return
		}
		cc = append(cc, currency)
	}

	respondWithJSON(w, 201, cc)
}

func (apiCfg *apiConfig) handleGetCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := apiCfg.DB.GetCurrencies(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching currencies - %v", err))
		return
	}
	respondWithJSON(w, 201, currencies)
}
