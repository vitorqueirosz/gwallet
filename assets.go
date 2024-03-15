package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vitorqueirosz/gwallet/internal/database"
)

type Asset struct {
	Code   string
	Amount string
}

func (apiCfg *apiConfig) handleCreateUserAssets(w http.ResponseWriter, r *http.Request, user database.User) {
	assets := []Asset{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&assets)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding body params - %v", err))
		return
	}

	currencies, err := apiCfg.DB.GetCurrencies(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching currencies - %v", err))
		return
	}

	ccMap := make(map[string]database.Currency)
	for _, c := range currencies {
		ccMap[c.Code] = c
	}

	ua := []database.Asset{}
	for _, asset := range assets {
		cc, ok := ccMap[asset.Code]
		if !ok {
			respondWithError(w, 400, fmt.Sprintf("Error updating user assets - currency code not valid - %v", asset.Code))
			return
		}
		a, err := apiCfg.DB.CreateUserAsset(r.Context(), database.CreateUserAssetParams{
			ID:         uuid.New(),
			CurrencyID: cc.ID,
			UserID:     user.ID,
			Amount:     asset.Amount,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error creating user asset - %v", err))
			return
		}
		ua = append(ua, a)
	}

	respondWithJSON(w, 201, ua)
}

func (apiCfg *apiConfig) handleGetUserBalance(w http.ResponseWriter, r *http.Request, user database.User) {
	assets, err := apiCfg.DB.GetUserAssets(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting user assets - %v", err))
		return
	}

	fiatAmount := decimal.New(0, 0)
	for _, a := range assets {
		am, err := decimal.NewFromString(a.Amount)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing asset amount to decimal - %v", err))
			return
		}
		cp, err := decimal.NewFromString(a.CurrencyPrice)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing currency price to decimal - %v", err))
			return
		}

		assetInFiat := am.Mul(cp)
		fiatAmount = fiatAmount.Add(assetInFiat)
	}

	type response struct {
		EstimateFiatBalance string
	}

	respondWithJSON(w, 200, response{
		EstimateFiatBalance: fiatAmount.String(),
	})
}
