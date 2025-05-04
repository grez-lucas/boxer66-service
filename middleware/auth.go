package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grez-lucas/boxer66-service/internal/config"
	"github.com/grez-lucas/boxer66-service/internal/repository"
)

type ContextKey string

var ContextUserKey ContextKey = "user"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("jwt-token")

		token, err := ValidateJWT(tokenStr)
		if err != nil {
			slog.Error("Failed to validate JWT", slog.Any("error", err))
			writeUnauthorized(w)
			return
		}

		if !token.Valid {
			slog.Error("Token is invalid", slog.String("token", tokenStr))
			writeUnauthorized(w)
			return
		}

		// check claims
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			slog.Error("Token has invalid claims")
			writeUnauthorized(w)
			return
		}

		// Add the userID to the request context for later use
		userID := claims["userID"].(int32)
		ctx := context.WithValue(context.Background(), ContextUserKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateJWT(user *repository.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"userID":    user.ID,
	}

	secret := config.LoadConfig().JWTSecret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	jwtSecret := config.LoadConfig().JWTSecret

	return jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["&alg"])
		}
		return []byte(jwtSecret), nil
	})
}

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}
