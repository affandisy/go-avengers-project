package service

import (
	"avenger/internal/domain"
	"avenger/internal/repository"
	"avenger/pkg/debug"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-playground/validator"
)

type InventoryService interface {
	GetAll() ([]domain.Inventory, error)
	GetByID(id int) (*domain.Inventory, error)
	Create(inv domain.Inventory) (int, error)
	Update(id int, inv domain.Inventory) error
	Delete(id int) error
}

type inventoryService struct {
	repo     repository.InventoryRepository
	validate *validator.Validate
}

func NewInventoryService(r repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: r, validate: validator.New()}
}

func (s *inventoryService) GetAll() ([]domain.Inventory, error) {
	debug.LogDebug("Fetching all inventories")

	inventories, err := s.repo.GetAll()
	if err != nil {
		debug.ErrorDebug("Failed to fetch inventory: %v", err)
		return nil, errors.New("failed to retrieve from database")
	}

	debug.LogDebug("Successfully fetched %d inventories", len(inventories))
	return inventories, nil
}

func (s *inventoryService) GetByID(id int) (*domain.Inventory, error) {
	debug.LogDebug("Fetching inventory with ID: %d", id)
	if id <= 0 {
		debug.ErrorDebug("invalid intentory ID: %d", id)
		return nil, errors.New("invalid inventory ID")
	}

	inventory, err := s.repo.GetByID(id)
	if err != nil {
		debug.LogDebug("Database error while fetching inventory ID %d: %v", id, err)
		return nil, errors.New("failed to retrieve inventory from database")
	}

	if inventory == nil {
		debug.LogDebug("inventory not found for ID: %d", id)
		return nil, errors.New("inventory not found")
	}

	debug.LogDebug("Successfully fetched inventory: %d", id)
	return inventory, nil
}

func (s *inventoryService) Create(inv domain.Inventory) (int, error) {
	debug.LogDebug("Creating new inventory")
	if err := s.validate.Struct(inv); err != nil {
		debug.ErrorDebug("validation error", err)
		return 0, errors.New("invalid inventory data")
	}

	if inv.Stock < 0 {
		debug.ErrorDebug("invalid stock value")
		return 0, errors.New("stock cannot be negative")
	}

	if inv.Status != "active" && inv.Status != "broken" {
		debug.ErrorDebug("invalid status value")
		return 0, errors.New("status must be active or broken")
	}

	inv.Code = strings.ToUpper(strings.TrimSpace(inv.Code))
	inv.Name = strings.TrimSpace(inv.Name)
	inv.Description = strings.TrimSpace(inv.Description)

	id, err := s.repo.Create(inv)
	if err != nil {
		if err == sql.ErrConnDone {
			debug.ErrorDebug("Duplicate inventory code")
			return 0, errors.New("inventory code already exists")
		}
		debug.ErrorDebug("Database error while creating inventory")
		return 0, errors.New("failed to create inventory in database")
	}

	debug.LogDebug("successfully created inventory")
	return id, nil
}

func (s *inventoryService) Update(id int, inv domain.Inventory) error {
	debug.LogDebug("Updating inventory ID")
	if id <= 0 {
		debug.ErrorDebug("invalid inventory id for update")
		return errors.New("invalid inventory id")
	}

	inv.ID = id
	if err := s.validate.Struct(inv); err != nil {
		debug.ErrorDebug("validation failed for update")
		return errors.New("invalid inventory data" + err.Error())
	}

	if inv.Stock < 0 {
		debug.ErrorDebug("invalid stock value for update")
		return errors.New("stock cannot be negative")
	}

	if inv.Status != "active" && inv.Status != "broken" {
		debug.ErrorDebug("Invalid status value for update: %s", inv.Status)
		return errors.New("status must be either 'active' or 'broken'")
	}

	inv.Code = strings.ToUpper(strings.TrimSpace(inv.Code))
	inv.Name = strings.TrimSpace(inv.Name)
	inv.Description = strings.TrimSpace(inv.Description)

	err := s.repo.Update(id, inv)
	if err != nil {
		if err == sql.ErrNoRows {
			debug.ErrorDebug("Inventory not found for update: ID %d", id)
			return errors.New("inventory not found")
		}
		if err == sql.ErrConnDone {
			debug.ErrorDebug("Duplicate inventory code on update: %s", inv.Code)
			return errors.New("inventory code already exists")
		}
		debug.LogDebug("Database error while updating inventory ID %d: %v", id, err)
		return errors.New("failed to update inventory in database")
	}

	debug.LogDebug("Successfully updated inventory ID: %d", id)
	return nil
}

func (s *inventoryService) Delete(id int) error {
	debug.LogDebug("Deleting inventory")

	if id <= 0 {
		debug.ErrorDebug("invalid inventory for deletion")
		return errors.New("invalid inventory id")
	}

	err := s.repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			debug.ErrorDebug("inventory not found for deletion")
			return errors.New("inventory not found")
		}
		debug.LogDebug("database error while deleting")
		return errors.New("failed to delete inventory from database")
	}

	debug.LogDebug("Successfully deleted inventory")

	return nil
}
