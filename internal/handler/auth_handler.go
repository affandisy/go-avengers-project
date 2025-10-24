package handler

import (
	"avenger/internal/domain"
	"avenger/internal/service"
	"avenger/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service service.UserService
}

func NewAuthHandler(s service.UserService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	if err := h.service.ValidateUser(user); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	if err := h.service.Register(&user); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var input domain.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	user, err := h.service.GetByEmail(input.Email)
	if err != nil || user == nil {
		http.Error(w, "invalid credential", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := utils.GenerateJWT(user.ID, user.Role)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
