package main

import (
	"database/sql"
	"fmt"
	"gameapp/config"
	"gameapp/handler"
	"gameapp/middleware"
	"gameapp/repository/mysql"
	"gameapp/service"
	"log"
	"net/http"
)

func main() {
	const (
		port    = "8080"
		address = "localhost:" + port
	)
	cfg := config.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	db, sErr := sql.Open("mysql", dsn)
	if sErr != nil {
		log.Fatalf("failed to connect to database: %v", sErr)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database  %v", err)
	}

	userRepo := mysql.NewUserRepository(db)

	authService := service.NewAuth(userRepo, cfg.JWTSecret)

	authHandler := handler.NewAuth(authService)

	mux := handler.NewSimpleMux()

	mux.HandleFunc("POST", "/api/auth/register", authHandler.Register)
	mux.HandleFunc("POST", "/api/auth/login", authHandler.Login)
	// protected routes with middleware
	protectedMux := handler.NewSimpleMux()
	protectedMux.HandleFunc("GET", "/api/profile", authHandler.GetProfile)
	protectedMux.HandleFunc("PUT", "/api/profile", authHandler.UpdateProfile)
	protectedMux.HandleFunc("POST", "/api/change-password", authHandler.ChangePassword)

	protectedHandler := middleware.AuthMiddleware(cfg.JWTSecret)(protectedMux)
	// combined router
	combinedMux := handler.NewSimpleMux()
	combinedMux.HandleFunc("POST", "/api/auth/register", authHandler.Register)
	combinedMux.HandleFunc("POST", "/api/auth/login", authHandler.Login)
	combinedMux.Handle("GET", "/api/profile", protectedHandler)
	combinedMux.Handle("PUT", "/api/profile", protectedHandler)
	combinedMux.Handle("POST", "/api/change-password", protectedHandler)
	log.Printf("you are listening to %s...", address)
	log.Printf("API Endpoints:")
	log.Printf("  POST   /api/auth/register")
	log.Printf("  POST   /api/auth/login")
	log.Printf("  GET    /api/profile (protected)")
	log.Printf("  PUT    /api/profile (protected)")
	log.Printf("  POST   /api/change-password (protected)")
	if err := http.ListenAndServe(address, combinedMux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
