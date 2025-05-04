package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grez-lucas/boxer66-service/internal/config"
	"github.com/grez-lucas/boxer66-service/internal/repository"
	"github.com/grez-lucas/boxer66-service/users"
	"github.com/jackc/pgx/v5"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	queries := repository.New(conn)

	uService := users.NewUserService(ctx, queries)
	uHandlers := users.NewUserHandlers(uService)

	router := http.NewServeMux()
	router.HandleFunc("/", greet)
	router.HandleFunc("GET /users", uHandlers.GetUsers)

	server := http.Server{
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		Handler:           router,
	}

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()
}
