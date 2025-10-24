package domain

import "gorm.io/gorm"

type Recipe struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	CookTime    int     `gorm:"not null" json:"cook_time"`
	Rating      float64 `gorm:"not null" json:"rating"`
}
