package service

import (
	"avenger/internal/domain"
	"avenger/internal/repository"
)

type InventoryService interface {
	GetAll() ([]domain.Inventory, error)
	GetByID(id int) (*domain.Inventory, error)
	Create(inv domain.Inventory) error
	Update(id int, inv domain.Inventory) error
	Delete(id int) error
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(r repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: r}
}

func (s *inventoryService) GetAll() ([]domain.Inventory, error) {
	return s.repo.GetAll()
}

func (s *inventoryService) GetByID(id int) (*domain.Inventory, error) {
	return s.repo.GetByID(id)
}

func (s *inventoryService) Create(inv domain.Inventory) error {
	return s.repo.Create(inv)
}

func (s *inventoryService) Update(id int, inv domain.Inventory) error {
	return s.repo.Update(id, inv)
}

func (s *inventoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
