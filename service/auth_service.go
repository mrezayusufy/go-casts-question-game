package service

import (
	"errors"
	"fmt"
	contract "gameapp/contract"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrWrongPassword      = errors.New("wrong password")
)

type Auth struct {
	repo         contract.UserRepositoryInterface
	tokenService Token
}

func NewAuth(repo contract.UserRepositoryInterface, jwtSecret string) *Auth {
	tokenService := NewToken([]byte(jwtSecret))
	return &Auth{
		repo:         repo,
		tokenService: *tokenService,
	}
}

// methods
// login
func (s *Auth) Login(req *dto.LoginRequest) (*string, error) {
	// validate password number
	if len(req.Password) == 0 || len(req.PhoneNumber) == 0 {
		return nil, errors.New("email and password is required")
	}
	// find by number
	user, err := s.repo.FindByPhoneNumber(req.PhoneNumber)
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

// get profile
func (s *Auth) GetProfile(id uint) (*entity.Profile, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user.ToProfile(), nil
}

// update profile
func (s *Auth) UpdateProfile(id uint, name string) (*entity.Profile, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	user.Name = name
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user.ToProfile(), nil
}
func (s *Auth) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	// Verify old password
	if cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); cErr != nil {
		return cErr
	}

	hashedPassword := hashPassword(newPassword)

	user.Password = string(hashedPassword)
	return s.repo.Update(user)
}

// register user
func (s *Auth) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// TODO verification number by otp

	if !phonenumber.IsValid(req.PhoneNumber) {
		return dto.RegisterResponse{}, fmt.Errorf("❌phone number is invalid")
	}
	// check uniqueness of number
	if userExists, pErr := s.repo.ExistsByPhoneNumber(req.PhoneNumber); pErr != nil || userExists {
		if pErr != nil {
			return dto.RegisterResponse{}, fmt.Errorf("❌unexpected error in validation of number %v", pErr)
		}
		if userExists {
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
	err := s.repo.Create(&user)
	if err != nil {
		return dto.RegisterResponse{
			Error: err.Error(),
		}, err
	}
	// return created user
	return dto.RegisterResponse{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Message:     " 🎉✨ You have successfully registered! 😃 ",
	}, nil
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
