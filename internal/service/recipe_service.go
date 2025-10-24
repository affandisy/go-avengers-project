package service

import (
	"avenger/internal/domain"
	"avenger/internal/repository"
)

type RecipeService interface {
	GetAll() ([]domain.Recipe, error)
	Create(recipe *domain.Recipe) error
	Delete(id int) error
}

type recipeService struct {
	repo repository.RecipeRepository
}

func NewRecipeService(r repository.RecipeRepository) RecipeService {
	return &recipeService{repo: r}
}

func (s *recipeService) GetAll() ([]domain.Recipe, error) {
	return s.repo.GetAll()
}

func (s *recipeService) Create(recipe *domain.Recipe) error {
	return s.repo.Create(recipe)
}

func (s *recipeService) Delete(id int) error {
	return s.repo.Delete(id)
}
