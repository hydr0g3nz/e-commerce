package model

import (
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
)

type Product struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	// SubCategory    string            `json:"sub_category"`
	Variations     []domain.Variation `json:"variations"`
	Specifications map[string]string  `json:"specifications"`
	ReviewIDs      []string           `json:"review_ids"`
	Rating         float64            `json:"rating"`
	Images         []string           `json:"images"`
}

func ProductDomainToModel(product *domain.Product) *Product {
	return &Product{
		Model: Model{
			ID: product.ID,
		},
		Name:        product.Name,
		Description: product.Description,
		Brand:       product.Brand,
		Category:    product.Category,
		// SubCategory:    product.SubCategory,
		Variations:     product.Variations,
		Specifications: product.Specifications,
		ReviewIDs:      product.ReviewIDs,
		Rating:         product.Rating,
		Images:         product.Images,
	}
}

func ProductModelToDomain(product *Product) *domain.Product {
	return &domain.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Brand:       product.Brand,
		Category:    product.Category,
		// SubCategory:    product.SubCategory,
		Variations:     product.Variations,
		Specifications: product.Specifications,
		ReviewIDs:      product.ReviewIDs,
		Rating:         product.Rating,
		Images:         product.Images,
	}
}

func (p *Product) Map() map[string]interface{} {
	return map[string]interface{}{
		"_id":         p.ID,
		"name":        p.Name,
		"description": p.Description,
		"brand":       p.Brand,
		"category":    p.Category,
		// "sub_category":   p.SubCategory,
		"variations":     p.Variations,
		"specifications": p.Specifications,
		"review_ids":     p.ReviewIDs,
		"rating":         p.Rating,
		"images":         p.Images,
	}
}
