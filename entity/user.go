package entity

import "time"

type User struct {
	ID          uint      `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Profile struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

func (u *User) ToProfile() *Profile {
	return &Profile{
		ID:          u.ID,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
	}
}
