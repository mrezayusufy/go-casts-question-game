package userrepository

import (
	"context"
	"database/sql"
	"fmt"
	"gameapp/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByPhoneNumber(ctx context.Context, phonenumber string) (*entity.User, error)
	FindByID(ctx context.Context, userID uint) (*entity.User, error)
}

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
	scanner := userScanner{user: &user}
	// query
	query := `SELECT id, name, phone_number, password FROM users WHERE phone_number = ?`
	// query row context
	row := r.db.QueryRowContext(ctx, query, phonenumber)
	// scan row
	err := scanner.ScanRow(row)

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
	scanner := &userScanner{user: &user}
	row := d.db.QueryRow("SELECT id, name, phone_number, password FROM users WHERE phone_number = ?", phonenumber)
	err := scanner.ScanRow(row)
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
	scanner := userScanner{user: &user}
	row := r.db.QueryRowContext(ctx, query, id)
	if err := scanner.ScanRow(row); err != nil {
		return entity.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

type userScanner struct {
	user *entity.User
}

func (us *userScanner) ScanRow(row *sql.Row) error {
	return row.Scan(
		&us.user.ID,
		&us.user.PhoneNumber,
		&us.user.Name,
		&us.user.Password,
	)
}
