package handler

import (
	"encoding/json"
	"errors"
	"gameapp/dto"
	"gameapp/middleware"
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
func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")

		return
	}
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	if req.PhoneNumber == "" || req.Password == "" {
		h.writeError(w, http.StatusBadRequest, "email and password is required")

		return
	}
	token, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrWrongPassword) {
			h.writeError(w, http.StatusUnauthorized, "invalid credentials")

			return
		}
		h.writeJSON(w, http.StatusInternalServerError, "internal server error")

		return
	}
	h.writeJSON(w, http.StatusOK, dto.LoginResponse{
		Token: *token,
	})

}

// get profile
func (h *Auth) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")

		return
	}
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.authService.GetProfile(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			h.writeError(w, http.StatusNotFound, "user not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "internal service error")

		return
	}
	h.writeJSON(w, http.StatusOK, profile)
}

// update profile
func (h *Auth) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")

		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")

		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "name is required")

		return
	}

	profile, err := h.authService.UpdateProfile(userID, req.Name)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "internal server error")

		return
	}

	h.writeJSON(w, http.StatusOK, profile)
}

// change password
func (h *Auth) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method post note allowed")

		return
	}
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeError(w, http.StatusUnauthorized, "unauthorized")

		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		h.writeError(w, http.StatusBadRequest, "old password and new password are required")

		return
	}
	if len(req.NewPassword) < 8 {
		h.writeError(w, http.StatusBadRequest, "new password must be at least 8 charachters")

		return
	}
	cErr := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if cErr != nil {
		if errors.Is(cErr, service.ErrWrongPassword) {
			h.writeError(w, http.StatusBadRequest, "password or phone number is wrong")

			return
		}
		h.writeError(w, http.StatusInternalServerError, cErr.Error())

		return
	}

	h.writeJSON(w, http.StatusOK, dto.MessageResponse{Message: "Password changed successfully!✨🎉✔😀"})

}

// write json
func (h *Auth) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// write error
func (h *Auth) writeError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, dto.ErrorResponse{Error: msg})
}
