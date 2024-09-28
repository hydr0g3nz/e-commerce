package adapters

import (
	"context"
	"fmt"

	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
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

func (r *CategoryRepository) Create(category *domain.Category) error {
	_, err := r.db.Collection("category").InsertOne(context.Background(), category)
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

func (r *CategoryRepository) Update(category *domain.Category) error {
	_, err := r.db.Collection("category").ReplaceOne(context.Background(), bson.M{"_id": category.ID}, category)
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
