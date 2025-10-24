package repository

import (
	"avenger/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Register(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Register(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
