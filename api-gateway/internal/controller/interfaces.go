package controller

import (
	"authjwt/internal/domain"
	"context"
	"github.com/labstack/echo/v4"
)

// Router — интерфейс для регистрации роутов.

type Router interface {
	RegisterRoutes(e *echo.Group)
}

type UserUseCase interface {
	CreateUser(ctx context.Context, user *domain.User, password string) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(h string) (bool, error)
}
