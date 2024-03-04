package main

import (
	"context"
	"fmt"
	"food-menu/handlers"
	"food-menu/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
		userStorage = storage.UserPgStorage{DB: conn}
		userHandler = handlers.UserHandler{UserStorage: userStorage}
	)

	mux.HandleFunc("POST /auth/create-user", userHandler.HandleCreate)
	mux.HandleFunc("POST /auth/login", userHandler.HandleLogin)
	mux.HandleFunc("POST /auth/refresh", userHandler.HandleRefresh)

	server := &http.Server{
		Addr:         ":4000",
		Handler:      mux,
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
		IdleTimeout:  time.Second * 120,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}

}
