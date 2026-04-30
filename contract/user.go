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
type UserRepositoryInterface interface {
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint) error
	ExistsByPhoneNumber(phonenumber string) (bool, error)
	FindByID(id uint) (*entity.User, error)
	FindByPhoneNumber(phonenumber string) (*entity.User, error)
}
