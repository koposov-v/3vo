package usecase

import (
	"authjwt/internal/domain"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}

type UserUseCase struct {
	repo      UserRepository
	jwtSecret string
}

func NewUserUseCase(repo UserRepository, jwtSecret string) UserUseCase {
	return UserUseCase{repo: repo, jwtSecret: jwtSecret}
}

func (u UserUseCase) CreateUser(ctx context.Context, user *domain.User, password string) error {
	if err := user.HashPassword(password); err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	return u.repo.CreateUser(ctx, user)
}

func (u UserUseCase) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	return u.repo.FindByUsername(ctx, username)
}

func (u UserUseCase) GenerateToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (u UserUseCase) ValidateToken(tokenStr string) (bool, error) {
	const prefix = "Bearer "
	if strings.HasPrefix(tokenStr, prefix) {
		tokenStr = tokenStr[len(prefix):]
	}

	tkn, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(u.jwtSecret), nil
	})
	if err != nil {
		return false, fmt.Errorf("invalid token: %w", err)
	}
	if !tkn.Valid {
		return false, fmt.Errorf("token is invalid or expired")
	}
	return true, nil
}
