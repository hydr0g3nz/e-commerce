package model

import (
	"github.com/google/uuid"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
)

type Category struct {
	// "category_id": "cat001",
	// "name": "Electronics",
	// "description": "Devices and gadgets",
	// "product_ids": ["987", "654", "321"]
	Model
	Name        string   `json:"name" bson:"name,omitempty"`
	Description string   `json:"description" bson:"description,omitempty"`
	ProductIDs  []string `json:"product_ids" bson:"product_ids,omitempty"`
}

func (c *Category) SetUUID() {
	id, _ := uuid.NewV7()
	c.ID = id.String()
}

func CategoryDomainToModel(d *domain.Category) Category {
	return Category{
		Model:       Model{ID: d.ID},
		Name:        d.Name,
		Description: d.Description,
		ProductIDs:  d.ProductIDs,
	}
}

func (c *Category) ModelToDomain() *domain.Category {
	return &domain.Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ProductIDs:  c.ProductIDs,
	}
}

func (c *Category) Map() map[string]interface{} {
	return map[string]interface{}{
		"_id":         c.ID,
		"created_at":  c.CreatedAt,
		"updated_at":  c.UpdatedAt,
		"deleted_at":  c.DeletedAt,
		"name":        c.Name,
		"description": c.Description,
		"product_ids": c.ProductIDs,
	}
}
