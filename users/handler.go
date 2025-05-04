package users

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type UserHandlers struct {
	ctx     context.Context
	service IUserService
}

func NewUserHandlers(
	service IUserService,
) *UserHandlers {
	return &UserHandlers{
		ctx:     context.Background(),
		service: service,
	}
}

func (h *UserHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Failed to get users from service", slog.Any("error", err))
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Failed to encode users", slog.Any("error", err))
		return
	}
}
