package router

import (
	"context"
	"net/http"

	"github.com/grez-lucas/boxer66-service/internal/config"
	"github.com/grez-lucas/boxer66-service/internal/repository"
	"github.com/grez-lucas/boxer66-service/smtp"
	"github.com/grez-lucas/boxer66-service/users"
)

func NewRouter(ctx context.Context, cfg *config.Config, queries *repository.Queries) http.Handler {
	// Initialize services and handlers
	uService := users.NewUserService(ctx, queries)
	smtpService := smtp.NewSMTPService(cfg.SMTPConfig)
	uHandlers := users.NewUserHandlers(uService, smtpService)

	router := http.NewServeMux()

	router.HandleFunc("GET /users", uHandlers.GetUsers)
	router.HandleFunc("POST /login", uHandlers.Login)
	router.HandleFunc("POST /register", uHandlers.Register)
	router.HandleFunc("POST /verify-email", uHandlers.VerifyEmail)

	router.Handle("/api/", http.StripPrefix("/api", router))
	return router
}
