package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
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
func (s *ProductService) DeleteImage(filename string) error {
	filename = filepath.Join(s.repo.Config().Upload.UploadPath, "products", filename)
	// Check if file exists
	_, err := os.Stat(filename)
	fmt.Println(s.repo.Config().Upload.UploadPath + filename)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not found")
		}
		return fmt.Errorf("error checking file: %v", err)
	}
	// Delete the file
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}
	return nil
}
