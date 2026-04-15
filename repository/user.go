package repository

import "gameapp/entity"

type UserResitory interface {
	IsPhoneNumberUnique(phonenumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
}
