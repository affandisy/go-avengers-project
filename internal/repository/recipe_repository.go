package repository

import (
	"avenger/internal/domain"

	"gorm.io/gorm"
)

type RecipeRepository interface {
	GetAll() ([]domain.Recipe, error)
	Create(recipe *domain.Recipe) error
	Delete(id int) error
}

type recipeRepository struct {
	DB *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{DB: db}
}

func (r *recipeRepository) GetAll() ([]domain.Recipe, error) {
	var recipes []domain.Recipe
	err := r.DB.Find(&recipes).Error
	return recipes, err
}

func (r *recipeRepository) Create(recipe *domain.Recipe) error {
	return r.DB.Create(recipe).Error
}

func (r *recipeRepository) Delete(id int) error {
	return r.DB.Delete(&domain.Recipe{}, id).Error
}
