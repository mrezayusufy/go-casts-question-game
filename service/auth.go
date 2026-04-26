package service

import (
	"context"
	"errors"
	contract "gameapp/contract"
	"gameapp/dto"
	"gameapp/entity"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*entity.User, error)
	validateUser(ctx context.Context, phonenumber, password string) (*entity.User, error)
}

type Auth struct {
	userRepo        contract.User
	passwordService password
	tokenService    Token
}

func NewAuth(userRepo contract.User, ps password, ts Token) *Auth {

	return &Auth{
		userRepo:        userRepo,
		passwordService: ps,
		tokenService:    ts,
	}
}

// methods
// login
func (s *Auth) Login(ctx context.Context, req *dto.LoginRequest) (*string, error) {
	// validate password phone number
	if len(req.Password) == 0 || len(req.PhoneNumber) == 0 {
		return nil, errors.New("email and password is required")
	}
	// find by phone number
	user, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	// check password
	if !s.passwordService.Verify(user.Password, req.Password) {
		return nil, errors.New("Invalid credentials")
	}
	token, tErr := s.tokenService.Generate(user.ID)
	if tErr != nil {
		return nil, tErr
	}
	return &token, nil
}

// validate user
func (s *Auth) validateUser(ctx context.Context, phoneNumber, password string) (*entity.User, error) {
	return nil, nil
}
