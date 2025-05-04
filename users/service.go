package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/grez-lucas/boxer66-service/internal/repository"
	"github.com/grez-lucas/boxer66-service/middleware"
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

var (
	ErrUserDoesntExist = errors.New("user does not exist")
	ErrInvalidPassword = errors.New("password is invalid")
)

func (h *UserService) GetUsers() ([]repository.User, error) {
	users, err := h.repository.GetAllUsers(h.ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (h *UserService) GetUserByEmail(email string) (*repository.User, error) {
	user, err := h.repository.GetUserByEmail(h.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserDoesntExist
		} else {
			return nil, err
		}
	}
	return &user, nil
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

func (h *UserService) Login(email string, requestPassword string) (*LoginResponse, error) {
	// Get the user
	user, err := h.repository.GetUserByEmail(h.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserDoesntExist
		}
		return nil, err
	}

	// Compare his request's password vs the hashedpassword
	if err := comparePassword(user.Password, requestPassword); err != nil {
		return nil, ErrInvalidPassword
	}

	token, err := middleware.CreateJWT(&user)
	if err != nil {
		return nil, err
	}
	resp := LoginResponse{Token: token, UserID: user.ID}

	return &resp, nil
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
