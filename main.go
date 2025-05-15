package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grez-lucas/boxer66-service/internal/config"
	"github.com/grez-lucas/boxer66-service/internal/repository"
	"github.com/grez-lucas/boxer66-service/internal/router"
	"github.com/grez-lucas/boxer66-service/middleware"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	queries := repository.New(conn)

	chain := middleware.CreateStack(middleware.Logging, middleware.Cors)

	router := router.NewRouter(ctx, cfg, queries)

	server := http.Server{
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		Handler:           chain(router),
	}

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()
}
