package userrepository

import (
	"context"
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
	res, err := r.db.Exec(`insert into users(name, phone_number, password) values(?, ?, ?)`, user.Name, user.PhoneNumber, user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}
	id, _ := res.LastInsertId()
	user.ID = uint(id)
	return user, nil
}

// find user by phone number
func (r *mysqlUserRepository) FindByPhoneNumber(ctx context.Context, phonenumber string) (*entity.User, error) {
	var user entity.User
	// query
	query := `SELECT id, name, phone_number, password FROM users WHERE phone_number = ?`
	// query row context
	row := r.db.QueryRowContext(ctx, query, phonenumber)
	// scan row
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)

	// error handling
	// no user found
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("❌ user not found 🙄😑😐")
	}
	// failed to scan
	if err != nil {
		return nil, fmt.Errorf("❌ failure to scan: %w", err)
	}
	return &user, nil
}

// is phone number unique
func (d *mysqlUserRepository) IsPhoneNumberUnique(phonenumber string) (bool, error) {
	user := entity.User{}
	row := d.db.QueryRow("SELECT id, name, phone_number, password FROM users WHERE phone_number = ?", phonenumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("can't scan query resutl: %v", err)
	}
	return false, nil
}

func (r *mysqlUserRepository) Get(ctx context.Context, id uint) (entity.User, error) {
	query := `SELECT id, name, password, phone_number FROM users where id = ?`
	user := entity.User{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Password, &user.PhoneNumber); err != nil {
		return entity.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}
