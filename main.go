package main

import (
	"encoding/json"
	"fmt"
	"gameapp/database"
	"gameapp/dto"
	UserRepository "gameapp/repository/user"
	UserService "gameapp/service/userservice"
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
	// create service and inject repository in service
	userService := UserService.New(userRepo)
	userCreated, ucErr := userService.Register(request)
	if ucErr != nil {
		fmt.Fprint(res, ucErr)
		return
	}
	fmt.Fprintf(res, "{'name':'%s','phone_number':'%s'}", userCreated.User.Name, userCreated.User.PhoneNumber)
}
