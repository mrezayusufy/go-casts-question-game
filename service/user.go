package service

import (
	"context"
	"gameapp/contract"
	"gameapp/entity"
)

// contract

type User struct {
	repo contract.User
}

// constructor injection
func NewUser(repo contract.User) *User {
	return &User{
		repo: repo,
	}
}
func (s *User) Get(ctx context.Context, id uint) (*entity.User, error) {
	// find user by id
	user, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
