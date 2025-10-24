package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password   string `gorm:"not null" json:"password,omitempty" validate:"required,min=8"`
	FullName   string `gorm:"not null" json:"full_name" validate:"required,min=6,max=15"`
	Age        int    `gorm:"not null" json:"age" validate:"required,gte=17"`
	Occupation string `gorm:"not null" json:"occupation" validate:"required"`
	Role       string `gorm:"not null;default:admin" json:"role" validate:"required,oneof=admin superadmin"`
}
