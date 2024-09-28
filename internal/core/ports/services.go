package ports

import "github.com/hydr0g3nz/e-commerce/internal/core/domain"

type CategoryService interface {
	Create(category *domain.Category) error
	GetAll() ([]*domain.Category, error)
	GetByID(id string) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id string) error
	AddProduct(categoryID string, productID string) error
}
