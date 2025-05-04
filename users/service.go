package users

import (
	"context"

	"github.com/grez-lucas/boxer66-service/internal/repository"
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
