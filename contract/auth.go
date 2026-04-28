package contract

import (
	"context"
	"gameapp/dto"
	"gameapp/entity"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, req dto.LoginRequest) (*entity.User, error)
	validateUser(ctx context.Context, phonenumber, password string) (*entity.User, error)
}
