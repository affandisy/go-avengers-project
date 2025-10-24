package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID         int    `gorm:"unique;not null" json:"id"`
	Email      string `gorm:"not null" json:"email"`
	Password   string `gorm:"not null" json:"password"`
	FullName   string `gorm:"not null" json:"full_name"`
	Age        int    `gorm:"not null" json:"age"`
	Occupation string `gorm:"not null" json:"occupation"`
	Role       string `gorm:"not null;default:admin" json:"role"`
}
