package dto

import "gameapp/entity"

type RegisterRequest struct {
	Name        string
	PhoneNumber string
}
type RegisterResponse struct {
	User entity.User
}
