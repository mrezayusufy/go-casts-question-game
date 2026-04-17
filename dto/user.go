package dto

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
type RegisterResponse struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	Error       string `json:"error"`
}
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}
