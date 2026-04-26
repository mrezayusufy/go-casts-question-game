package contract

import (
	"context"
	"gameapp/entity"
)

type User interface {
	IsPhoneNumberUnique(phonenumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
	FindByPhoneNumber(ctx context.Context, phonenumber string) (*entity.User, error)
	Get(ctx context.Context, id uint) (entity.User, error)
}
