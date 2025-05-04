package users

import (
	"context"
	"encoding/json"
	"errors"
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

func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var loginRequest LoginRequest
	if err := decoder.Decode(&loginRequest); err != nil {
		slog.Error("Failed to decode LoginResponse", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Password is invalid"))
			return
		}
		if errors.Is(err, ErrUserDoesntExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User does not exist"))
			return
		}
		slog.Error("An error occured", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
