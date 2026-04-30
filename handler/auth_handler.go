package handler

import (
	"encoding/json"
	"errors"
	"gameapp/dto"
	"gameapp/service"
	"net/http"
)

type Auth struct {
	authService *service.Auth
}

func NewAuth(as *service.Auth) *Auth {
	return &Auth{authService: as}
}

// methods
// register
func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")

		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	if req.PhoneNumber == "" || req.Password == "" || req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "all fields are required")

		return
	}
	if len(req.Password) < 8 {
		h.writeError(w, http.StatusBadRequest, "password must be at least 8 charachters")

		return
	}
	response, rErr := h.authService.Register(req)
	if rErr != nil {
		if errors.Is(rErr, service.ErrUserExists) {
			h.writeError(w, http.StatusConflict, "user already exists")

			return
		}
		h.writeError(w, http.StatusInternalServerError, "internal server error")

		return
	}
	h.writeJSON(w, http.StatusCreated, response)
}

// login
// get profile
// update profile
// change password
// write json
func (h *Auth) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// write error
func (h *Auth) writeError(w http.ResponseWriter, status int, msg string) {

}

// get user id from context
// get path params
