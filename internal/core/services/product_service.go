package services

import (
	"errors"
	"mime/multipart"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	if !product.IsCanCreate() {
		return errors.New(domain.ErrInvalidProduct)
	}
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
	if !variation.IsCanAdd() {
		return errors.New(domain.ErrInvalidVariation)
	}
	return s.repo.AddVariation(productID, variation)
}

func (s *ProductService) RemoveVariation(productID string, variationID string) error {
	return s.repo.RemoveVariation(productID, variationID)
}

func (s *ProductService) UploadImage(ctx *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	id, _ := uuid.NewV7()
	filename := id.String() + filepath.Ext(file.Filename)
	filepath := filepath.Join(s.repo.Config().Upload.UploadPath, "products", filename)
	if err := ctx.SaveFile(file, filepath); err != nil {
		return "", err
	}
	return filename, nil
}
