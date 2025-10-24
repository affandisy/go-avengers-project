package domain

import "gorm.io/gorm"

type Recipe struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name" validate:"required,min=3,max=100"`
	Description string  `gorm:"not null" json:"description" validate:"required,min=10,max=1000"`
	CookTime    int     `gorm:"not null" json:"cook_time" validate:"required,gt=0"`
	Rating      float64 `gorm:"not null" json:"rating" validate:"required,gte=0,lte=5"`
}
