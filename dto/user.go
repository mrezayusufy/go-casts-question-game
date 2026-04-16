package dto

import "gameapp/entity"

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}
type RegisterResponse struct {
	User entity.User
}
