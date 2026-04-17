package userservice

import (
	"context"
	"fmt"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"

	"golang.org/x/crypto/bcrypt"
)

// contract
type UserResitory interface {
	IsPhoneNumberUnique(phonenumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
	FindByPhoneNumber(ctx context.Context, phonenumber string) (*entity.User, error)
	Get(ctx context.Context, id uint) (entity.User, error)
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
	// validate password
	if len(req.Password) < 8 {
		return dto.RegisterResponse{}, fmt.Errorf("Password length is at leaset 8 character ")
	}
	// password := hashPassword(req.Password)
	password := hashPassword(req.Password)

	// TODO: check the password with regex pattern
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    string(password),
	}
	newUser, err := s.repo.Register(user)
	if err != nil {
		return dto.RegisterResponse{
			Error: err.Error(),
		}, err
	}
	// return created user
	return dto.RegisterResponse{
		Name:        newUser.Name,
		PhoneNumber: newUser.PhoneNumber,
		Message:     " 🎉✨ You have successfully registered! 😃 ",
	}, nil
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	// check the exitance of phone number from repository
	// get the user by phonenumber
	user, err := s.repo.FindByPhoneNumber(ctx, req.PhoneNumber)
	// error handling
	const errMsg = `{"error":"❌phone number or password is not correct😐😑🙄"}`
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf(errMsg)
	}
	// campare user password with the req.password
	if !CheckPasswordHash(req.Password, user.Password) {
		return dto.LoginResponse{}, fmt.Errorf(errMsg)
	}
	// return ok
	return dto.LoginResponse{Token: "✅ You have successfully login!😎🖐😀"}, nil
}
func (s *Service) Get(ctx context.Context, id uint) (*entity.User, error) {
	// find user by id
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func hashPassword(plain string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hashed)
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
