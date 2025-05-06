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

func WriteJSON(w http.ResponseWriter, v any, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	resp := StatusResponse{Status: APIResponseStatusError, Message: message}
	WriteJSON(w, resp, statusCode)
}

func WriteSuccess(w http.ResponseWriter, mesg string, statusCode int) {
	w.WriteHeader(statusCode)
	resp := StatusResponse{Status: APIResponseStatusSuccess, Message: mesg}
	WriteJSON(w, resp, statusCode)
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

	user, token, err := h.service.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			WriteError(w, "Password is invalid", http.StatusBadRequest)
			return
		}
		if errors.Is(err, ErrUserDoesntExist) {
			WriteError(w, "User does not exist", http.StatusNotFound)
			return
		}
		slog.Error("An error occured", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		UserID: user.ID,
		Token:  token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		slog.Error("Failed to decode registerRequest", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.Register(registerRequest.Email, registerRequest.Password); err != nil {
		slog.Error("Failed to register user", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteSuccess(w, "Verification code sent to email", http.StatusAccepted)
}

func (h *UserHandlers) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var verifyEmailRequest VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&verifyEmailRequest); err != nil {
		slog.Error("Failed to decode verifyEmailRequest", slog.Any("error", err))
		// TODO: Check for missing required fields
		w.WriteHeader(http.StatusBadRequest)
	}

	user, jwt, err := h.service.VerifyEmailToken(verifyEmailRequest.Email, verifyEmailRequest.Token)
	if err != nil {
		slog.Error(
			"Error verifying email",
			slog.String("email", verifyEmailRequest.Email),
			slog.String("token", verifyEmailRequest.Token),
			slog.Any("error", err),
		)
		if errors.Is(err, ErrTokenIsExpired) {
			WriteError(w, "Token is expired", http.StatusBadRequest)
			return
		}
		if errors.Is(err, ErrInvalidToken) {
			WriteError(w, "Token is invalid", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := VerifyEmailResponse{
		Token:  jwt,
		UserID: user.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
