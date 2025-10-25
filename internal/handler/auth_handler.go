package handler

import (
	"avenger/internal/domain"
	"avenger/internal/service"
	"avenger/pkg/utils"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service  service.UserService
	validate *validator.Validate
}

func NewAuthHandler(s service.UserService) *AuthHandler {
	return &AuthHandler{service: s, validate: validator.New()}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", map[string]string{
			"body": "Request body must be valid JSON",
		})
		return
	}

	if err := h.service.ValidateUser(user); err != nil {
		slog.Warn("User validation failed", slog.Any("error", err))
		writeError(w, http.StatusBadRequest, "Invalid request body", map[string]string{
			"validation": err.Error(),
		})
		return
	}

	if len(user.Password) < 8 {
		writeError(w, http.StatusBadRequest, "Validation failed", map[string]string{
			"password": "Password must be at least 8 characters",
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "Failed to process password", nil)
		return
	}
	user.Password = string(hashed)

	if err := h.service.Register(&user); err != nil {
		slog.Error("Failed to register user", slog.Any("error", err))
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			writeError(w, http.StatusConflict, "Email already registered", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to register user", nil)
		return
	}

	user.Password = ""

	writeJSON(w, http.StatusCreated, Response{
		Message: "User registered successfully",
		Data:    user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var input domain.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", map[string]string{
			"body": "Request body must be valid JSON",
		})
		return
	}

	if input.Email == "" || input.Password == "" {
		writeError(w, http.StatusBadRequest, "Validation failed", map[string]string{
			"credentials": "Email and password are required",
		})
		return
	}

	user, err := h.service.GetByEmail(input.Email)
	if err != nil {
		slog.Error("Failed to get user by email", slog.String("email", input.Email), slog.Any("error", err))
		writeError(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if user == nil {
		slog.Warn("Login attempt with non-existent email", slog.String("email", input.Email))
		writeError(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		slog.Warn("Login attempt with incorrect password", slog.String("email", input.Email))
		writeError(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	token, err := utils.GenerateJWT(int(user.ID), user.Role)
	if err != nil {
		slog.Error("Failed to generate JWT token", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "Failed to generate authentication token", nil)
		return
	}

	slog.Info("User logged in successfully", slog.String("email", user.Email), slog.String("role", user.Role))

	writeJSON(w, http.StatusOK, Response{
		Message: "Login successsful",
		Data: map[string]any{
			"token": token,
			"user": map[string]any{
				"id":        user.ID,
				"email":     user.Email,
				"full_name": user.FullName,
				"role":      user.Role,
			},
		},
	})
}
