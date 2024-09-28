package domain

type Category struct {
	// "category_id": "cat001",
	// "name": "Electronics",
	// "description": "Devices and gadgets",
	// "product_ids": ["987", "654", "321"]
	ID          string   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ProductIDs  []string `json:"product_ids"`
}
