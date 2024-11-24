package ports

import (
	"context"

	"github.com/hydr0g3nz/e-commerce/internal/adapters/dto"
	"github.com/hydr0g3nz/e-commerce/internal/adapters/model"
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
	GetProductBySku(ctx context.Context, productId, sku string) (*domain.Product, error)
	ReserveStock(ctx context.Context, productId, sku string, quantity int) error
	ReleaseStock(ctx context.Context, productId, sku string, quantity int) error
	SetProductList(ctx context.Context, product []dto.ProductListPage) error
	GetProductList(ctx context.Context) ([]dto.ProductListPage, error)
	SetProductHeroList(ctx context.Context, product []dto.ProductListPage) error
	GetProductHeroList(ctx context.Context) ([]dto.ProductListPage, error)
	GetProductsCategoryDelegate(ctx context.Context) (map[string]*domain.Product, error)
	GetCacheProductsCategoryDelegate(ctx context.Context) (map[string]*domain.Product, error)
	SetProductCategoryDelegate(ctx context.Context, product map[string]*domain.Product) error
	GetByCategory(ctx context.Context, category string) ([]*domain.Product, error)
}

type AuthRepository interface {
	// CreateUser(user *domain.User) error
	// GetUserByEmail(email string) (*domain.User, error)
	// GetUserByID(id string) (*domain.User, error)
	// UpdateUser(user *domain.User) error
	// DeleteUser(id string) error
	FindByEmail(email string) (*domain.User, error)
	CreateRefreshToken(userId string, metadata *domain.TokenMetadata) error
	FetchRefreshToken(userId string) (*model.RefreshToken, error)
	EmailExists(email string) bool
	Create(user *domain.User) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) (string, error)
	UpdateStatus(ctx context.Context, orderID, status string) error
}
