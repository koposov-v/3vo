package domain

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

func (u *User) HashPassword(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) VerifyPassword(input string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input)) == nil
}
