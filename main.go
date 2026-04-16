package main

import (
	"gameapp/database"
	"gameapp/dto"
	userrepository "gameapp/repository/user"
	UserService "gameapp/service/userservice"
	"log"
)

func main() {
	// load configuration
	cfg := database.Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		DBName:   "gameapp_db",
	}
	// 2. create database connect
	db, err := database.NewConn(cfg)
	if err != nil {
		log.Fatalf("failure to connect database: %v", err)
	}
	defer db.Close()
	// 3. inject into repository
	userRepo := userrepository.New(db)
	// 4. inject repository into service
	userService := UserService.New(userRepo)
	// 5. use service to call create user
	response, rErr := userService.Register(dto.RegisterRequest{
		Name:        "Hasan",
		PhoneNumber: "09020072667",
	})
	if rErr != nil {
		log.Fatalf("failure to create user: %v", rErr)
	}
	log.Printf("user: %+v", response)
}
