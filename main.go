package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/vitorqueirosz/gwallet/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("")

	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the env")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the env")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("DATABASE CONNECTION FAILED")
	}

	db := database.New(conn)
	apiConfig := apiConfig{
		DB: db,
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

	router.Post("/users", apiConfig.handleCreateUser)
	router.Get("/users", apiConfig.authMiddleware(apiConfig.handleGetUserByApiKey))
	router.Post("/users/assets", apiConfig.authMiddleware(apiConfig.handleCreateUserAssets))
	router.Get("/users/balance", apiConfig.authMiddleware(apiConfig.handleGetUserBalance))

	router.Post("/currencies", apiConfig.handleCreateCurrencies)
	router.Get("/currencies", apiConfig.handleGetCurrencies)

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Printf("Server staring on port %v", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
