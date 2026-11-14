package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Ganesh-12-spec/Linkboard/internal/service"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) AuthHandler {
	return AuthHandler{
		userService: userService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Name, email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(
		r.Context(),
		req.Name,
		req.Email,
		req.Password,
	)
	if err != nil {
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}
