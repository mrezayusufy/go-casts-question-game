package dto

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Name        string `json:"name"`
}
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}
type UpdateProfileRequest struct {
	Name string `json:"name"`
}
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type UserProfileDTO struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type MessageResponse struct {
	Message string `json:"message"`
}
type RegisterResponse struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	Error       string `json:"error"`
}
