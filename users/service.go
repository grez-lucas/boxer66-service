package users

import (
	"context"

	"github.com/grez-lucas/boxer66-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	ctx        context.Context
	repository *repository.Queries
}

func NewUserService(
	ctx context.Context,
	repository *repository.Queries,
) *UserService {
	return &UserService{
		ctx:        ctx,
		repository: repository,
	}
}

func (h *UserService) GetUsers() ([]repository.User, error) {
	users, err := h.repository.GetAllUsers(h.ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (h *UserService) CreateUser(email, password string) (*repository.User, error) {
	// Encrypt the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := h.repository.CreateUser(h.ctx, repository.CreateUserParams{
		Email:    email,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *UserService) Login(tokenStr string) bool {
	return false
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func comparePassword(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err
}
