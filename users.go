package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

type userController struct {
	users []User
}

func UserController() *userController {
	return &userController{}
}

func parseApiKeyFromUrl(w http.ResponseWriter, r *http.Request) (*uuid.UUID, error) {
	apiKeyStr := chi.URLParam(r, "apiKey")
	apiKey, err := uuid.Parse(apiKeyStr)
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (u *userController) getUserByApiKey(apiKey uuid.UUID) *User {
	var userFromApiKey *User
	for i, user := range u.users {
		if user.ApiKey == apiKey {
			userFromApiKey = &u.users[i]
		}
	}
	return userFromApiKey
}

func (u *userController) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	user := User{
		ID:     uuid.New(),
		Name:   params.Name,
		ApiKey: uuid.New(),
	}
	u.users = append(u.users, user)

	respondWithJSON(w, 201, user)
}

func (u *userController) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	apiKey, err := parseApiKeyFromUrl(w, r)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding request body - %v", err))
		return
	}

	user := u.getUserByApiKey(*apiKey)
	if user.Name == "" {
		respondWithError(w, 400, fmt.Sprint("Error fetching user by api key - user not found"))
		return
	}

	respondWithJSON(w, 200, user)
}

type assetHandler func(http.ResponseWriter, *http.Request, *[]Currency)

func assetMiddleware(handler assetHandler, c *[]Currency) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	}
}

func (u *userController) handleCreateUserAssets(w http.ResponseWriter, r *http.Request, currencies *[]Currency) {
	assets := []Asset{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&assets)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding body params - %v", err))
		return
	}

	apiKey, err := parseApiKeyFromUrl(w, r)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding request body - %v", err))
		return
	}

	ccMap := make(map[string]Currency)
	for _, c := range *currencies {
		ccMap[c.Code] = c
	}

	user := u.getUserByApiKey(*apiKey)

	for _, asset := range assets {
		_, ok := ccMap[asset.Code]
		if !ok {
			respondWithError(w, 400, fmt.Sprintf("Error updating user assets - currency code not valid - %v", asset.Code))
			return
		}
		user.Assets = append(user.Assets, asset)
	}

	respondWithJSON(w, 200, user)
}

// func (u *userController) handleGetUserBalance(w http.ResponseWriter, r *http.Request) {
// 	apiKey, err := parseApiKeyFromUrl(w, r)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Error decoding request body - %v", err))
// 		return
// 	}

// 	user := u.getUserByApiKey(*apiKey)

// 	for _, currency
// }
