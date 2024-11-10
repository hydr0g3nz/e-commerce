package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hydr0g3nz/e-commerce/internal/adapters/dto"
	"github.com/hydr0g3nz/e-commerce/internal/adapters/model"
	"github.com/hydr0g3nz/e-commerce/internal/config"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/pkg/redis"
	"github.com/hydr0g3nz/e-commerce/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CacheKeyProductList     = "product-list"
	CacheKeyProductHeroList = "product-hero-list"
)

type ProductRepository struct {
	db         *mongo.Database
	cache      *redis.RedisClient
	cfg        *config.Config
	locks      map[string]*sync.Mutex
	locksMutex sync.RWMutex
}

func (r *ProductRepository) Config() *config.Config {
	return r.cfg
}
func NewProductRepository(cfg *config.Config, db *mongo.Client, cache *redis.RedisClient) *ProductRepository {
	Db := db.Database("e-commerce")
	return &ProductRepository{db: Db, cache: cache, cfg: cfg, locks: make(map[string]*sync.Mutex)}
}

func (r *ProductRepository) Create(p *domain.Product) error {
	product := model.ProductDomainToModel(p)
	product.BeforeCreate()
	_, err := r.db.Collection("product").InsertOne(context.Background(), product)
	return err
}
func (r *ProductRepository) GetByID(id string) (*domain.Product, error) {
	var product model.Product
	err := r.db.Collection("product").FindOne(context.Background(), bson.M{"_id": id, "deleted_at": nil}).Decode(&product)
	if err != nil {
		return nil, err
	}
	productD := model.ProductModelToDomain(&product)
	return productD, nil
}
func (r *ProductRepository) Update(p *domain.Product) error {
	product := model.ProductDomainToModel(p)
	product.BeforeUpdate()
	productMap := product.Map()
	util.MapDeleteNilOrZero(productMap)
	_, err := r.db.Collection("product").UpdateOne(context.Background(), bson.M{"_id": product.ID}, bson.M{"$set": bson.M(productMap)})
	return err
}
func (r *ProductRepository) Delete(id string) error {
	_, err := r.db.Collection("product").UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}
func (r *ProductRepository) GetAll() ([]*domain.Product, error) {
	var products []*model.Product
	cursor, err := r.db.Collection("product").Find(context.Background(), bson.M{"deleted_at": nil})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &products)
	if err != nil {
		return nil, err
	}
	productDomainList := model.ProductListModelToDomainList(products)
	return productDomainList, nil
}

func (r *ProductRepository) AddVariation(productID string, variation *domain.Variation) error {
	update := bson.M{
		"$addToSet": bson.M{"variations": variation},
		"$set":      bson.M{"updated_at": time.Now()},
	}
	_, err := r.db.Collection("product").UpdateOne(context.Background(), bson.M{"_id": productID}, update)
	return err
}

func (r *ProductRepository) RemoveVariation(productID string, variationID string) error {
	update := bson.M{
		"$pull": bson.M{"variations": bson.M{"sku": variationID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.db.Collection("product").UpdateOne(context.Background(), bson.M{"_id": productID}, update)
	return err
}
func (r *ProductRepository) GetProductBySku(ctx context.Context, productId, sku string) (*domain.Product, error) {
	collection := r.db.Collection(productCollection)

	var product domain.Product
	err := collection.FindOne(ctx, bson.M{
		"_id":            productId,
		"variations.sku": sku,
	}).Decode(&product)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) getProductLock(productId string) *sync.Mutex {
	r.locksMutex.RLock()
	if mu, exists := r.locks[productId]; exists {
		r.locksMutex.RUnlock()
		return mu
	}
	r.locksMutex.RUnlock()

	// If we get here, we need to create a new mutex
	r.locksMutex.Lock()
	defer r.locksMutex.Unlock()

	// Double-check in case another goroutine created it
	if mu, exists := r.locks[productId]; exists {
		return mu
	}

	// Create new mutex
	mu := &sync.Mutex{}
	r.locks[productId] = mu
	return mu
}

func (r *ProductRepository) ReserveStock(ctx context.Context, productId, sku string, quantity int) error {
	// Get the mutex for this product
	mu := r.getProductLock(productId)

	// Lock the mutex
	mu.Lock()
	defer mu.Unlock()
	collection := r.db.Collection("product")

	filter := bson.M{
		"_id":              productId,
		"variations.sku":   sku,
		"variations.stock": bson.M{"$gte": quantity}, // Check stock in filter
	}

	update := bson.M{
		"$inc": bson.M{
			"variations.$[elem].stock": -quantity, // Changed from $ to $[elem]
		},
	}

	// Configure array filters correctly
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.sku": sku},
		},
	})

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("insufficient stock or product not found")
	}

	return nil
}

// Cleanup method to remove unused locks (optional)
func (r *ProductRepository) CleanupLock(productId string) {
	r.locksMutex.Lock()
	defer r.locksMutex.Unlock()
	delete(r.locks, productId)
}

func (r *ProductRepository) ReleaseStock(ctx context.Context, productId, sku string, quantity int) error {
	collection := r.db.Collection(productCollection)

	update := bson.M{
		"$inc": bson.M{
			"variations.$.stock": quantity,
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.sku": sku},
		},
	})

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": productId, "variations.sku": sku},
		update,
		arrayFilters,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}

// Additional helper methods for product repository

func (r *ProductRepository) UpdateProductPrice(ctx context.Context, sku string, price float64) error {
	collection := r.db.Collection(productCollection)

	update := bson.M{
		"$set": bson.M{
			"variations.$.price": price,
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.sku": sku},
		},
	})

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"variations.sku": sku},
		update,
		arrayFilters,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *ProductRepository) UpdateSale(ctx context.Context, sku string, salePercentage float32) error {
	collection := r.db.Collection(productCollection)

	update := bson.M{
		"$set": bson.M{
			"variations.$.sale": salePercentage,
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.sku": sku},
		},
	})

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"variations.sku": sku},
		update,
		arrayFilters,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *ProductRepository) SetProductList(ctx context.Context, product []dto.ProductListPage) error {
	return r.cache.Set(ctx, CacheKeyProductList, product, time.Hour*24)
}
func (r *ProductRepository) SetProductHeroList(ctx context.Context, product []dto.ProductListPage) error {
	return r.cache.Set(ctx, CacheKeyProductHeroList, product, time.Hour*24)
}
func (r *ProductRepository) GetProductList(ctx context.Context) ([]dto.ProductListPage, error) {
	var product []dto.ProductListPage
	err := r.cache.Get(ctx, CacheKeyProductList, &product)
	return product, err
}
func (r *ProductRepository) GetProductHeroList(ctx context.Context) ([]dto.ProductListPage, error) {
	var product []dto.ProductListPage
	err := r.cache.Get(ctx, CacheKeyProductHeroList, &product)
	return product, err
}
