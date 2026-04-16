package userrepository

import (
	"database/sql"
	"fmt"
	"gameapp/entity"
)

// struct
type mysqlUserRepository struct {
	db *sql.DB
}

// concret object
func New(db *sql.DB) *mysqlUserRepository {
	return &mysqlUserRepository{db: db}
}

// methods
// save
func (r *mysqlUserRepository) Register(user entity.User) (entity.User, error) {
	res, err := r.db.Exec(`insert into users(name, phone_number) values(?, ?)`, user.Name, user.PhoneNumber)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}
	id, _ := res.LastInsertId()
	user.ID = uint(id)
	return user, nil
}

// is phone number unique
func (d *mysqlUserRepository) IsPhoneNumberUnique(phonenumber string) (bool, error) {
	user := entity.User{}
	var createdAt []uint8
	row := d.db.QueryRow("SELECT * FROM users WHERE phone_number = ?", phonenumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("can't scan query resutl: %v", err)
	}
	return false, nil
}
