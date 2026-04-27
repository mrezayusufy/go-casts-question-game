package mysql

import (
	"database/sql"
	"gameapp/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// create
func (r *UserRepository) Create(user *entity.User) error {
	// query
	query := `
		INSERT INTO users (phone_number, name, password, created_at, updated_at) 
		VALUES (?, ?, ?, NOW(), NOW())
	`
	// db execute and get result
	result, err := r.db.Exec(query, user.PhoneNumber, user.Name, user.Password)
	// error handling
	if err != nil {
		return err
	}
	// get last id
	id, lErr := result.LastInsertId()
	if lErr != nil {
		return lErr
	}
	// update last id
	user.ID = uint(id)
	return nil
}

// update
func (r *UserRepository) Update(user *entity.User) error {

	return nil
}

// delete
// Exists By PhoneNumber
func (r *UserRepository) ExistsByPhoneNumber(phoneNumber string) (bool, error) {
	// query
	query := `SELECT EXISTS(SELECT 1 FROM users where phone_number = ?);`
	var exists bool
	rErr := r.db.QueryRow(query, phoneNumber).Scan(&exists)
	return exists, rErr
}

// find by phone number
// find by id
