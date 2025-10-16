package domain

type Inventory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Stock       int    `json:"stock"`
	Description int    `json:"description"`
	Status      string `json:"status"`
}
