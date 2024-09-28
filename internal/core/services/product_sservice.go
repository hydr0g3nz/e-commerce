package services

import (
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/ports"
)

type ProductService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(product *domain.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetAll() ([]*domain.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id string) (*domain.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *domain.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *ProductService) AddVariation(productID string, variation *domain.Variation) error {
	return s.repo.AddVariation(productID, variation)
}

func (s *ProductService) RemoveVariation(productID string, variationID string) error {
	return s.repo.RemoveVariation(productID, variationID)
}
