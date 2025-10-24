package repository

import (
	"avenger/internal/domain"
	"database/sql"
	"strings"
)

type InventoryRepository interface {
	GetAll() ([]domain.Inventory, error)
	GetByID(id int) (*domain.Inventory, error)
	Create(inv domain.Inventory) (int, error)
	Update(id int, inv domain.Inventory) error
	Delete(id int) error
}

type inventoryRepository struct {
	DB *sql.DB
}

func NewInventoryRepository(db *sql.DB) InventoryRepository {
	return &inventoryRepository{DB: db}
}

func (r *inventoryRepository) GetAll() ([]domain.Inventory, error) {
	rows, err := r.DB.Query("SELECT id, name, code, stock, description, status FROM inventories ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		err := rows.Scan(&inv.ID, &inv.Name, &inv.Code, &inv.Stock, &inv.Description, &inv.Status)
		if err != nil {
			return nil, err
		}
		list = append(list, inv)
	}

	return list, nil
}

func (r *inventoryRepository) GetByID(id int) (*domain.Inventory, error) {
	row := r.DB.QueryRow("SELECT id, name, code, stock, description, status FROM inventories WHERE id = $1", id)
	var inv domain.Inventory
	err := row.Scan(&inv.ID, &inv.Name, &inv.Code, &inv.Stock, &inv.Description, &inv.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &inv, nil
}

func (r *inventoryRepository) Create(inv domain.Inventory) (int, error) {
	query := `
	INSERT INTO inventories (name, code, stock, description, status)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	var id int
	err := r.DB.QueryRow(
		query,
		inv.Name,
		inv.Code,
		inv.Stock,
		inv.Description,
		inv.Status,
	).Scan(&id)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return 0, sql.ErrConnDone
		}
		return 0, err
	}

	return id, nil
}

func (r *inventoryRepository) Update(id int, inv domain.Inventory) error {
	result, err := r.DB.Exec(`UPDATE inventories SET name=$1, code=$2, stock=$3, description=$4, status=$5 WHERE id=$6`, inv.Name, inv.Code, inv.Stock, inv.Description, inv.Status, id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return sql.ErrConnDone
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *inventoryRepository) Delete(id int) error {
	result, err := r.DB.Exec("DELETE FROM inventories WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
