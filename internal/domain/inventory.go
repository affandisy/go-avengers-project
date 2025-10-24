package domain

type Inventory struct {
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required, min=3, max=100"`
	Code        string `json:"code" validate:"required, min=3, max=50"`
	Stock       int    `json:"stock" validate:"required,gte=0"`
	Description string `json:"description" validate:"max=500"`
	Status      string `json:"status" validate:"required,oneof=active broken"`
}
