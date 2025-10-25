package service

import (
	"avenger/internal/domain"
	"avenger/internal/repository"
	"avenger/pkg/debug"
	"errors"
	"net/mail"
	"strings"
)

type UserService interface {
	Register(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	ValidateUser(user domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) Register(user *domain.User) error {
	debug.LogDebug("Registering new user: Email=%s, Role=%s", user.Email, user.Role)
	err := s.repo.Register(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique constraint") {
			debug.ErrorDebug("Duplicate email on registration: %s", user.Email)
			return errors.New("email already exists")
		}
		debug.LogDebug("Database error while registering user: %v", err)
		return errors.New("failed to register user")
	}

	debug.LogDebug("Successfully registered user with ID: %d", user.ID)
	return nil
}

func (s *userService) GetByEmail(email string) (*domain.User, error) {
	debug.LogDebug("Fetching user by email: %s", email)

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		debug.ErrorDebug("Database error while fetching user by email: %v", err)
		return nil, errors.New("failed to retrieve user")
	}

	if user == nil {
		debug.LogDebug("User not found for email: %s", email)
		return nil, nil
	}

	debug.LogDebug("Successfully fetched user")
	return user, nil
}

func (s *userService) ValidateUser(user domain.User) error {
	debug.LogDebug("Validating user data for: %s", user.Email)

	if user.Email == "" {
		return errors.New("email tidak boleh kosong")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.New("format email tidak valid")
	}
	if len(user.Password) < 8 {
		return errors.New("password minimal 8 karakter")
	}
	if len(user.FullName) < 6 || len(user.FullName) > 15 {
		return errors.New("full name minimal 6 dan maksimal 15 karakter")
	}
	if user.Age < 17 {
		return errors.New("umur minimal 17 tahun")
	}
	if user.Occupation == "" {
		return errors.New("occupation tidak boleh kosong")
	}
	if user.Role == "" {
		user.Role = "admin"
	}

	debug.LogDebug("User validation passed")
	return nil
}
