package usecase

import (
	"authjwt/internal/domain"
	"context"
	"github.com/pkg/errors"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
}

type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) UserUseCase {
	return UserUseCase{repo: repo}
}

func (u UserUseCase) CreateUser(ctx context.Context, user *domain.User, password string) error {
	if err := user.HashPassword(password); err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	return u.repo.CreateUser(ctx, user)
}
