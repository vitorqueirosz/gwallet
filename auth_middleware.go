package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/vitorqueirosz/gwallet/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func getApiKeyFromHeader(headers http.Header) (string, error) {
	// ApiKey (token)
	authHeader := headers.Get("Authorization")

	authHeaderSplit := strings.Split(authHeader, " ")
	if len(authHeaderSplit) != 2 {
		return "", errors.New("Invalid authorization header")
	}

	if authHeaderSplit[0] != "ApiKey" {
		return "", errors.New("Malformed authorization header")
	}

	return authHeaderSplit[1], nil
}

func (apiCfg *apiConfig) authMiddleware(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := getApiKeyFromHeader(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error - %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("User not found - %v", err))
			return
		}

		handler(w, r, user)
	}
}
