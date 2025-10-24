package service

import (
	"avenger/internal/domain"
	"avenger/internal/repository"
	"errors"
	"net/mail"
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
	return s.repo.Register(user)
}

func (s *userService) GetByEmail(email string) (*domain.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) ValidateUser(user domain.User) error {
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

	return nil
}
