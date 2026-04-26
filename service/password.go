package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	Hash(password string) (string, error)
	Validate(plain, hash string) bool
}
type password struct {
	cost int
}

func NewPassword(cost int) *password {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &password{cost: cost}
}

// methods
// hash
func (s *password) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("failure to hash password %w", err)
	}

	return string(hash), nil
}

// verify
func (s *password) Verify(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(plain), []byte(hash))
	return err == nil
}
