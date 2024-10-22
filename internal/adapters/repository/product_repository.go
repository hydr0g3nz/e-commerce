package adapters

import (
	"context"
	"time"

	"github.com/hydr0g3nz/e-commerce/internal/adapters/model"
	"github.com/hydr0g3nz/e-commerce/internal/config"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/pkg/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	db  *mongo.Database
	cfg *config.Config
}

func (r *ProductRepository) Config() *config.Config {
	return r.cfg
}
func NewProductRepository(cfg *config.Config, db *mongo.Client) *ProductRepository {
	Db := db.Database("e-commerce")
	return &ProductRepository{db: Db, cfg: cfg}
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