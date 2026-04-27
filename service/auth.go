package service

import (
	"context"
	"errors"
	"fmt"
	contract "gameapp/contract"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*entity.User, error)
	validateUser(ctx context.Context, phonenumber, password string) (*entity.User, error)
}

type Auth struct {
	userRepo     contract.User
	tokenService Token
}

func NewAuth(userRepo contract.User) *Auth {
	tokenService := NewToken([]byte("game-app-secret-key"))
	return &Auth{
		userRepo:     userRepo,
		tokenService: *tokenService,
	}
}

// methods
// login
func (s *Auth) Login(ctx context.Context, req *dto.LoginRequest) (*string, error) {
	// validate password number
	if len(req.Password) == 0 || len(req.PhoneNumber) == 0 {
		return nil, errors.New("email and password is required")
	}
	// find by number
	user, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	// check password
	if !verifyPassword(user.Password, req.Password) {
		return nil, errors.New("Invalid credentials")
	}
	token, tErr := s.tokenService.Generate(user.ID)
	if tErr != nil {
		return nil, tErr
	}
	return &token, nil
}

// register user
func (s User) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// TODO verification number by otp

	if !phonenumber.IsValid(req.PhoneNumber) {
		return dto.RegisterResponse{}, fmt.Errorf("❌phone number is invalid")
	}
	// check uniqueness of number
	if isUnique, pErr := s.repo.IsPhoneNumberUnique(req.PhoneNumber); pErr != nil || !isUnique {
		if pErr != nil {
			return dto.RegisterResponse{}, fmt.Errorf("❌unexpected error in validation of number %v", pErr)
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

// validate user
func (s *Auth) validateUser(ctx context.Context, phoneNumber, password string) (*entity.User, error) {
	return nil, nil
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash)
}

// verify
func verifyPassword(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(plain), []byte(hash))
	return err == nil
}
