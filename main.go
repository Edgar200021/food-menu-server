package main

import (
	"context"
	"fmt"
	"food-menu/handlers"
	"food-menu/middlewares"
	"food-menu/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	conn, err := pgxpool.New(ctx, os.Getenv("DB_CONNECTION"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	mux := http.NewServeMux()
	var (
		userStorage    = storage.UserPgStorage{DB: conn}
		productStorage = storage.ProductPgStorage{DB: conn}
		userHandler    = handlers.UserHandler{UserStorage: userStorage}
		productHandler = handlers.ProductHandler{ProductStorage: productStorage}
	)

	mux.HandleFunc("POST /api/v1/auth/sign-up", userHandler.HandleCreate)
	mux.HandleFunc("POST /api/v1/auth/login", userHandler.HandleLogin)
	mux.HandleFunc("POST /api/v1/auth/refresh", userHandler.HandleRefresh)

	mux.HandleFunc("/api/v1/products/{id}", middlewares.AuthRequired(productHandler.HandleGetProduct, &userStorage))
	mux.HandleFunc("/api/v1/products", middlewares.AuthRequired(productHandler.HandleGetProducts, &userStorage))
	mux.HandleFunc("POST /api/v1/products", middlewares.AuthRequired(productHandler.HandleCreateProduct, &userStorage))

	handler := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodHead},
		Debug:            true,
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedHeaders:   []string{"*"},
		//OptionsPassthrough: true,
		//ExposedHeaders:     []string{},

	}).Handler(mux)

	server := &http.Server{
		Addr:         ":4000",
		Handler:      handler,
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
		IdleTimeout:  time.Second * 120,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
