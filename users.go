package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/vitorqueirosz/gwallet/internal/database"
)

type Asset struct {
	Code   string
	Amount string
}

type User struct {
	ID     uuid.UUID
	Name   string
	ApiKey uuid.UUID
	Assets []Asset
}

func parseApiKeyFromUrl(w http.ResponseWriter, r *http.Request) (*uuid.UUID, error) {
	apiKeyStr := chi.URLParam(r, "apiKey")
	apiKey, err := uuid.Parse(apiKeyStr)
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Name string
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding body params - %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating new user - %v", err))
		return
	}

	respondWithJSON(w, 201, user)
}

func (apiCfg *apiConfig) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, user)
}

type assetHandler func(http.ResponseWriter, *http.Request, *[]Currency)

func assetMiddleware(handler assetHandler, c *[]Currency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	}
}

// func (apiCfg *apiConfig) handleCreateUserAssets(w http.ResponseWriter, r *http.Request, currencies *[]Currency) {
// 	assets := []Asset{}

// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&assets)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Error decoding body params - %v", err))
// 		return
// 	}

// 	apiKey, err := parseApiKeyFromUrl(w, r)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Error decoding request body - %v", err))
// 		return
// 	}

// 	ccMap := make(map[string]Currency)
// 	for _, c := range *currencies {
// 		ccMap[c.Code] = c
// 	}

// 	user := u.getUserByApiKey(*apiKey)

// 	for _, asset := range assets {
// 		_, ok := ccMap[asset.Code]
// 		if !ok {
// 			respondWithError(w, 400, fmt.Sprintf("Error updating user assets - currency code not valid - %v", asset.Code))
// 			return
// 		}
// 		user.Assets = append(user.Assets, asset)
// 	}

// 	respondWithJSON(w, 200, user)
// }

// func (apiCfg *apiConfig) handleGetUserBalance(w http.ResponseWriter, r *http.Request, currencies *[]Currency) {
// 	apiKey, err := parseApiKeyFromUrl(w, r)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Error decoding request body - %v", err))
// 		return
// 	}

// 	ccMap := make(map[string]Currency)
// 	for _, c := range *currencies {
// 		ccMap[c.Code] = c
// 	}

// 	user := u.getUserByApiKey(*apiKey)

// 	fiatAmount := decimal.Decimal{}
// 	for _, asset := range user.Assets {
// 		c, ok := ccMap[asset.Code]
// 		if !ok {
// 			respondWithError(w, 400, fmt.Sprintf("Error getting user balance - currency code not valid - %v", asset.Code))
// 			return
// 		}

// 		dp, _ := decimal.NewFromString(c.Price)
// 		da, _ := decimal.NewFromString(asset.Amount)

// 		assetAmount := da.Mul(dp)
// 		fiatAmount = fiatAmount.Add(assetAmount)
// 	}

// 	type response struct {
// 		EstimateFiatBalance string
// 	}

// 	respondWithJSON(w, 200, response{
// 		EstimateFiatBalance: fiatAmount.String(),
// 	})
// }
