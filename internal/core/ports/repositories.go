package ports

import (
	"github.com/hydr0g3nz/e-commerce/internal/config"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
)

type CategoryRepository interface {
	Create(category *domain.Category) error
	GetByID(id string) (*domain.Category, error)
	GetAll() ([]*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id string) error
	AddProduct(categoryID string, productID string) error
	RemoveProduct(categoryID string, productID string) error
}
type ProductRepository interface {
	Config() *config.Config
	Create(product *domain.Product) error
	GetByID(id string) (*domain.Product, error)
	GetAll() ([]*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id string) error
	AddVariation(productID string, variation *domain.Variation) error
	RemoveVariation(productID string, variationID string) error
}
