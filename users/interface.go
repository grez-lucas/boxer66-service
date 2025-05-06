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
	Login(email, requestPassword string) (user *repository.User, jwt string, err error)
	Register(email, password string) error
	VerifyEmailToken(email, token string) (user *repository.User, jwt string, err error)
	CreateUser(email, requestPassword string) (*repository.User, error)
}
