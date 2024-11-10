package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hydr0g3nz/e-commerce/internal/adapters/dto"
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

func (s *ProductService) SetProductList(product []dto.ProductListPage) error {
	if len(product) == 0 {
		productDbList, err := s.GetAll()
		if err != nil {
			return err
		}
		productList := []dto.ProductListPage{}
		for _, productDb := range productDbList {
			sale := 0
			image1 := ""
			image2 := ""
			slug := ""
			price := 0
			for _, variation := range productDb.Variations {
				if variation.Sale > 0 && int(variation.Sale) > sale {
					sale = int(variation.Sale)
				}
				if len(variation.Images) > 1 {
					image1 = variation.Images[0]
					image2 = variation.Images[1]
				}
				slug = variation.Sku
				price = int(variation.Price)
			}
			productList = append(productList, dto.ProductListPage{
				ID:          productDb.ID,
				Name:        productDb.Name,
				Slug:        slug,
				VariantsNum: len(productDb.Variations),
				Price:       price,
				Sale:        sale,
				Image1:      image1,
				Image2:      image2,
				Category:    productDb.Category,
			})
		}
		product = productList
	}
	return s.repo.SetProductList(context.Background(), product)
}

func (s *ProductService) GetProductList(ctx context.Context) ([]dto.ProductListPage, error) {
	return s.repo.GetProductList(ctx)
}

func (s *ProductService) InitProductList() error {
	return s.SetProductList([]dto.ProductListPage{})
}

func (s *ProductService) SetProductHeroList() error {
	productDbList, err := s.GetAll()
	if err != nil {
		return err
	}
	productList := []dto.ProductListPage{}
	for _, productDb := range productDbList {
		sale := 0
		image1 := ""
		image2 := ""
		slug := ""
		price := 0
		for _, variation := range productDb.Variations {
			if variation.Sale > 0 && int(variation.Sale) > sale {
				sale = int(variation.Sale)
				if len(variation.Images) > 1 {
					image1 = variation.Images[0]
					image2 = variation.Images[1]
				}
			}
			slug = variation.Sku
			price = int(variation.Price)
		}
		productList = append(productList, dto.ProductListPage{
			ID:          productDb.ID,
			Name:        productDb.Name,
			Slug:        slug,
			VariantsNum: len(productDb.Variations),
			Price:       price,
			Sale:        sale,
			Image1:      image1,
			Image2:      image2,
			Category:    productDb.Category,
		})
	}
	sort.Slice(productList, func(i, j int) bool {
		return productList[i].Sale > productList[j].Sale
	})
	return s.repo.SetProductHeroList(context.Background(), productList[:4])
}
func (s *ProductService) GetProductHeroList(ctx context.Context) ([]dto.ProductListPage, error) {
	return s.repo.GetProductHeroList(ctx)
}
func (s *ProductService) InitProductHeroList() error {
	return s.SetProductHeroList()
}
