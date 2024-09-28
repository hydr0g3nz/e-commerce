package services

import (
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/ports"
)

type CategoryService struct {
	repo ports.CategoryRepository
}

func NewCategoryService(repo ports.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategoryAll() ([]*domain.Category, error) {
	return s.repo.GetAll()
}
func (s *CategoryService) GetCategory(id string) (*domain.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) CreateCategory(category *domain.Category) error {
	return s.repo.Create(category)
}

func (s *CategoryService) UpdateCategory(category *domain.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) DeleteCategory(id string) error {
	return s.repo.Delete(id)
}
func (s *CategoryService) AddProduct(categoryID string, productID string) error {
	return s.repo.AddProduct(categoryID, productID)
}
