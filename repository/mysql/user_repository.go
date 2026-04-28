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
		VALUES (?, ?, ?, NOW(), NOW());
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
	// query
	query := `UPDATE users SET name = ?, phone_number = ?, updated_at = NOW() WHERE id = ?;`
	_, err := r.db.Exec(query, user.Name, user.PhoneNumber, user.ID)
	return err
}

// delete
func (r *UserRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// Exists By PhoneNumber
func (r *UserRepository) ExistsByPhoneNumber(phoneNumber string) (bool, error) {
	// query
	query := `SELECT EXISTS(SELECT 1 FROM users where phone_number = ?);`
	var exists bool
	rErr := r.db.QueryRow(query, phoneNumber).Scan(&exists)
	return exists, rErr
}

// find by phone number
func (r *UserRepository) FindByPhoneNumber(phonenumber string) (*entity.User, error) {
	// query
	query := `SELECT id, name, phone_number FROM users WHERE phone_number = ?;`
	user := &entity.User{}
	err := r.db.QueryRow(query, phonenumber).Scan(&user.ID, &user.Name, &user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// find by id
func (r *UserRepository) FindByID(phonenumber string) (*entity.User, error) {
	// query
	query := `SELECT id, name, phone_number FROM users WHERE phone_number = ?;`
	user := &entity.User{}
	err := r.db.QueryRow(query, phonenumber).Scan(&user.ID, &user.Name, &user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return user, nil
}
