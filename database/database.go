package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host, Port, User, Password, DBName string
}

func NewConn(cfg Config) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failure to open database: %w\n", err)
	}
	// open conn
	db.SetMaxOpenConns(25)
	// max idle conn
	db.SetMaxIdleConns(25)
	// max life time 3 min
	db.SetConnMaxLifetime(5 * time.Minute)
	if pErr := db.Ping(); pErr != nil {
		return nil, fmt.Errorf("failure to ping database: %w\n", pErr)
	}
	return db, nil
}
