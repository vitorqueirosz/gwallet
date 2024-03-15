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
