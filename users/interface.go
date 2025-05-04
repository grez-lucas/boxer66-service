package users

import (
	"net/http"

	"github.com/grez-lucas/boxer66-service/internal/repository"
)

type IUserHanlders interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type IUserService interface {
	GetUsers() ([]repository.User, error)
	Login(email string, requestPassword string) (*LoginResponse, error)
}
