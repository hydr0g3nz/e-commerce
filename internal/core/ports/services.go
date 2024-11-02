package ports

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
)

type CategoryService interface {
	Create(category *domain.Category) error
	GetAll() ([]*domain.Category, error)
	GetByID(id string) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id string) error
	AddProduct(categoryID string, productID string) error
	RemoveProduct(categoryID string, productID string) error
}
type ProductService interface {
	Create(product *domain.Product) error
	GetAll() ([]*domain.Product, error)
	GetByID(id string) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id string) error
	AddVariation(productID string, variation *domain.Variation) error
	RemoveVariation(productID string, variationID string) error
	UploadImage(ctx *fiber.Ctx, file *multipart.FileHeader) (string, error)
	DeleteImage(filename string) (string, error)
}

type AuthService interface {
	ValidateCredentials(email, password string) (*domain.User, error)
	Login(credentials *domain.UserCredentials) (*domain.User, *domain.TokenResponse, error)
	RefreshAccessToken(refreshToken string) (*domain.TokenResponse, error)
	GenerateTokenPair(userId string, role string) (*domain.TokenResponse, error)
	Register(request *domain.User) (*domain.TokenResponse, error)
}
