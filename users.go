package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID
	Name       string
	ApiKey     uuid.UUID
	Currencies []Currency
}

type userController struct {
	users []User
}

func UserController() *userController {
	return &userController{}
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
	apiKeyStr := chi.URLParam(r, "apiKey")
	apiKey, err := uuid.Parse(apiKeyStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing apiKey - apiKey not valid - %v", err))
	}

	user := User{}
	for _, u := range u.users {
		if u.ApiKey == apiKey {
			user = u
		}
	}

	if user.Name == "" {
		respondWithError(w, 400, fmt.Sprint("Error fetching user by api key - user not found"))
		return
	}

	respondWithJSON(w, 200, user)
}

// func (u *userController) handleCreateUserAssets(w http.ResponseWriter, r *http.Request) {
// 	type body struct {

// 	}

// 	apiKey, err := uuid.Parse(apiKeyStr)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Error parsing apiKey - apiKey not valid - %v", err))
// 	}

// 	user := User{}
// 	for _, u := range u.users {
// 		if u.ApiKey == apiKey {
// 			user = u
// 		}
// 	}

// 	if user.Name == "" {
// 		respondWithError(w, 400, fmt.Sprint("Error fetching user by api key - user not found"))
// 		return
// 	}

// 	respondWithJSON(w, 200, user)
// }
