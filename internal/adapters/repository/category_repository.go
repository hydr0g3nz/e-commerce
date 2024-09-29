package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/hydr0g3nz/e-commerce/internal/adapters/model"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/pkg/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository struct {
	db *mongo.Database
}

func NewCategoryRepository(db *mongo.Client) *CategoryRepository {
	Db := db.Database("e-commerce")
	return &CategoryRepository{db: Db}
}

func (r *CategoryRepository) Create(c *domain.Category) error {
	category := model.CategoryDomainToModel(c)
	category.BeforeCreate()
	_, err := r.db.Collection("category").InsertOne(context.Background(), category)
	fmt.Println("save category", category)
	return err
}

func (r *CategoryRepository) GetByID(id string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Collection("category").FindOne(context.Background(), bson.M{"_id": id, "deleted_at": nil}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Update(c *domain.Category) error {
	category := model.CategoryDomainToModel(c)
	category.BeforeUpdate()
	mCategory := category.Map()
	util.MapDeleteNilOrZero(mCategory)
	_, err := r.db.Collection("category").UpdateOne(context.Background(), bson.M{"_id": category.ID}, bson.M{"$set": bson.M(mCategory)})
	return err
}

func (r *CategoryRepository) Delete(id string) error {
	_, err := r.db.Collection("category").UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}

func (r *CategoryRepository) GetAll() ([]*domain.Category, error) {
	var categories []*domain.Category
	cursor, err := r.db.Collection("category").Find(context.Background(), bson.M{"deleted_at": nil})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) AddProduct(categoryID string, productID string) error {
	update := bson.M{
		"$addToSet": bson.M{"product_ids": productID},
		"$set":      bson.M{"updated_at": time.Now()},
	}
	_, err := r.db.Collection("category").UpdateOne(context.Background(), bson.M{"_id": categoryID}, update)
	return err
}
func (r *CategoryRepository) RemoveProduct(categoryID string, productID string) error {
	update := bson.M{
		"$pull": bson.M{"product_ids": productID},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err := r.db.Collection("category").UpdateOne(context.Background(), bson.M{"_id": categoryID}, update)
	return err
}
