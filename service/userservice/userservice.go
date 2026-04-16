package userservice

import (
	"fmt"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"
)

// contract
type UserResitory interface {
	IsPhoneNumberUnique(phonenumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
}
type Service struct {
	repo UserResitory
}

// concret object
func New(repo UserResitory) *Service {
	return &Service{
		repo: repo,
	}
}
func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// TODO verification phone number by otp
	// Validate phone number

	if !phonenumber.IsValid(req.PhoneNumber) {
		return dto.RegisterResponse{}, fmt.Errorf("❌phone number is invalid")
	}
	// check uniqueness of phone number
	if isUnique, pErr := s.repo.IsPhoneNumberUnique(req.PhoneNumber); pErr != nil || !isUnique {
		if pErr != nil {
			return dto.RegisterResponse{}, fmt.Errorf("❌unexpected error in validation of phone number %v", pErr)
		}
		if !isUnique {
			return dto.RegisterResponse{}, fmt.Errorf("❌phone is not unique")
		}

	}
	// validate name
	if len(req.Name) < 3 {
		return dto.RegisterResponse{}, fmt.Errorf("Name length should be greater than 3")
	}
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
	}
	newUser, err := s.repo.Register(user)
	if err != nil {
		return dto.RegisterResponse{}, err
	}
	// create new user in storage
	// return created user
	return dto.RegisterResponse{
		User: newUser,
	}, nil
}
