package adapters

import (
	"context"
	"fmt"

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
	categoryMap := category.Map()
	_, err := r.db.Collection("category").InsertOne(context.Background(), categoryMap)
	fmt.Println("save category", category)
	return err
}

func (r *CategoryRepository) GetByID(id string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Collection("category").FindOne(context.Background(), bson.M{"_id": id}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Update(c *domain.Category) error {
	category := model.CategoryDomainToModel(c)
	mCategory := category.Map()
	util.MapDeleteNilOrZero(mCategory)
	fmt.Printf("cat %+v\n", mCategory)
	_, err := r.db.Collection("category").UpdateOne(context.Background(), bson.M{"_id": category.ID}, bson.M{"$set": bson.M(mCategory)})
	return err
}

func (r *CategoryRepository) Delete(id string) error {
	_, err := r.db.Collection("category").DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *CategoryRepository) GetAll() ([]*domain.Category, error) {
	var categories []*domain.Category
	cursor, err := r.db.Collection("category").Find(context.Background(), bson.M{})
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
	_, err := r.db.Collection("category").UpdateOne(context.Background(), bson.M{"_id": categoryID}, bson.M{"$addToSet": bson.M{"product_ids": productID}})
	return err
}
