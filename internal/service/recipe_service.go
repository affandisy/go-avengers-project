package service

import (
	"avenger/internal/domain"
	"avenger/internal/repository"
	"avenger/pkg/debug"
	"errors"
	"strings"

	"gorm.io/gorm"
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
	debug.LogDebug("Fetching all recipes")

	recipes, err := s.repo.GetAll()
	if err != nil {
		debug.ErrorDebug("Failed to fetch recipes")
		return nil, errors.New("failed to retrieve recipes from database")
	}

	debug.LogDebug("Successfully fetched %d recipes", len(recipes))
	return recipes, nil
}

func (s *recipeService) Create(recipe *domain.Recipe) error {
	debug.LogDebug("Creating new recipe")
	if recipe.CookTime <= 0 {
		debug.ErrorDebug("Invalid cook time")
		return errors.New("cook time must be greater than 0")
	}

	if recipe.Rating < 0 || recipe.Rating > 5 {
		debug.ErrorDebug("Invalid rating")
		return errors.New("rating must be between 0 and 5")
	}

	recipe.Name = strings.TrimSpace(recipe.Name)
	recipe.Description = strings.TrimSpace(recipe.Description)

	err := s.repo.Create(recipe)
	if err != nil {
		debug.ErrorDebug("Database error while creating recipe")
		return errors.New("failed to create recipe in database")
	}

	debug.LogDebug("Successfully created recipe")
	return nil
}

func (s *recipeService) Delete(id int) error {
	debug.LogDebug("Deleting recipe")

	if id <= 0 {
		debug.ErrorDebug("Invalid recipe ID for deletion %d", id)
		return errors.New("invalid recipe ID")
	}

	err := s.repo.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			debug.LogDebug("Recipe not found for deletion")
			return errors.New("recipe not found")
		}
		debug.ErrorDebug("Database error while deleting recipe ID")
		return errors.New("failed to delete recipe from database")
	}

	debug.LogDebug("Successfully delete recipe ID: %d", id)
	return nil
}
