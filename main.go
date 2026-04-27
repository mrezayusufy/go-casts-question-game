package main

import (
	"encoding/json"
	"fmt"
	"gameapp/database"
	"gameapp/dto"
	repository "gameapp/repository"
	UserRepository "gameapp/repository/user"
	service "gameapp/service"
	"io"
	"log"
	"net/http"
)

func main() {
	const (
		port    = "8080"
		address = "localhost:" + port
	)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/users/register", userRegisterHandler)
	http.HandleFunc("/users/login", LoginHandler)
	http.HandleFunc("/users/profile", UserProfileHandler)
	log.Printf("you are listening to %s...", address)
	http.ListenAndServe(address, nil)
}
func homeHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Fprintf(res, "invalid method")
		return
	}
	fmt.Fprintf(res, "Hello user")

}
func UserProfileHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, `{"error": "%s method is not allowed"}`, req.Method)

		return
	}
	// config
	cfg := database.Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		DBName:   "gameapp_db",
	}
	// new connection
	db, dErr := database.NewConn(cfg)
	if dErr != nil {
		log.Println("error in connecting to mysq", dErr)
		return
	}
	defer db.Close()
	ctx := req.Context()
	// create repository
	userRepo := repository.NewUser(db)
	passwordService := service.NewPassword(10)
	tokenService := service.NewToken([]byte("question-game-app-secret"))
	service.NewAuth(userRepo, *passwordService, *tokenService)

}
func LoginHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, `{"error": "method not allowed"}`)

		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, `{"error": "Content-Type must be application/json"}`)
		return
	}
	if req.Body == nil {
		fmt.Fprintf(w, `{"error": "empty body"}`)

		return
	}
	defer req.Body.Close()
	var request dto.LoginRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err == io.EOF {
			fmt.Fprintf(w, `{"error": "request body is empty"}`)
		} else {
			fmt.Fprintf(w, `{"error": "invalid JSON: %s"}`, err.Error())
		}
		return
	}
	// config
	cfg := database.Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		DBName:   "gameapp_db",
	}
	// new connection
	db, dErr := database.NewConn(cfg)
	if dErr != nil {
		log.Println("error in connecting to mysq", dErr)
		return
	}
	defer db.Close()
	ctx := req.Context()
	// create repository
	userRepo := repository.NewUser(db)
	// create service and inject repository into service
	passwordService := service.NewPassword(10)
	tokenService := service.NewToken([]byte("question-game-app-secret"))
	authService := service.NewAuth(userRepo, *passwordService, *tokenService)
	token, ucErr := authService.Login(ctx, &request)
	if ucErr != nil {
		fmt.Fprint(w, ucErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)

}
func userRegisterHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		fmt.Fprintf(res, "You are not Registered")
		return
	}
	if req.Body == nil {
		fmt.Fprintf(res, "empty_body %s", "request body is required")
		return
	}
	defer req.Body.Close()

	if req.Header.Get("Content-Type") != "application/json" {
		res.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(res, `{"error": "Content-Type must be application/json"}`)
		return
	}

	var request dto.RegisterRequest
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		if err == io.EOF {
			fmt.Fprintf(res, `{"error": "request body is empty"}`)
		} else {
			fmt.Fprintf(res, `{"error": "invalid JSON: %s"}`, err.Error())
		}
		return
	}
	// config
	cfg := database.Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		DBName:   "gameapp_db",
	}
	// new connection
	db, dErr := database.NewConn(cfg)
	if dErr != nil {
		log.Println("error in connecting to mysq", dErr)
		return
	}
	defer db.Close()
	// create repository
	userRepo := UserRepository.New(db)
	// create service and inject repository into service
	userService := service.NewUser(userRepo)
	userCreated, ucErr := userService.Register(request)
	if ucErr != nil {
		fmt.Fprint(res, ucErr)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(userCreated)
}
