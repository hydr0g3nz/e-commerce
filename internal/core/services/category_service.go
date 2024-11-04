package services

import (
	"fmt"

	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/ports"
)

type CategoryService struct {
	repo ports.CategoryRepository
}

func NewCategoryService(repo ports.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]*domain.Category, error) {
	return s.repo.GetAll()
}
func (s *CategoryService) GetByID(id string) (*domain.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Create(category *domain.Category) error {
	return s.repo.Create(category)
}

func (s *CategoryService) Update(category *domain.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) Delete(id string) error {
	return s.repo.Delete(id)
}
func (s *CategoryService) AddProduct(categoryID string, productID string) error {
	return s.repo.AddProduct(categoryID, productID)
}

func (s *CategoryService) RemoveProduct(categoryID string, productID string) error {
	fmt.Println("categoryID", categoryID, "productID", productID)
	return s.repo.RemoveProduct(categoryID, productID)
}
