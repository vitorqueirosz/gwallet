package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

// ROUTES

// Register user
// Create Currencies - memory for now, once we have the scrapper we're going to pushed them to the db
// Update currencies by scrapping in x interval
// Register assets (currencies)
// Get total balance
// Get currencies
// Get currency by id

// User
// name
// api_key
// currencies

// User -> asset
// name
// code
// amount
// user_id

// Currency
// name
// code
// price

func main() {
	fmt.Println("")

	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the env")
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/healthz", handleReadiness)

	userController := UserController()
	router.Post("/users", userController.handleCreateUser)
	router.Get("/users/{apiKey}", userController.handleGetUserByApiKey)

	currencyController := CurrencyController()
	router.Post("/currencies", currencyController.createCurrencies)

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Printf("Server staring on port %v", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}