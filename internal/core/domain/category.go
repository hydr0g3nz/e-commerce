package domain

type Category struct {
	ID          string   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ProductIDs  []string `json:"product_ids"`
	ParentID    string   `json:"parent_id"`
	SubCategory []string `json:"sub_category"`
}
