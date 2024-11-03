package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/hydr0g3nz/e-commerce/internal/adapters/model"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	orderCollection   = "orders"
	productCollection = "product"
)

type OrderRepository struct {
	db *mongo.Database
}

func NewOrderRepository(db *mongo.Client) *OrderRepository {
	database := db.Database("e-commerce")
	return &OrderRepository{db: database}
}

func (r *OrderRepository) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(orderCollection)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "date", Value: 1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	collection := r.db.Collection(orderCollection)
	m := model.DomainOrderToModel(order)
	m.BeforeCreate()
	_, err := collection.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID, status string) error {
	collection := r.db.Collection(orderCollection)

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": orderID},
		update,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("order not found")
	}

	return nil
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, orderID string) (*domain.Order, error) {
	collection := r.db.Collection(orderCollection)

	var order domain.Order
	err := collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

// GetUserOrders retrieves all orders for a specific user
func (r *OrderRepository) GetUserOrders(ctx context.Context, userID string) ([]domain.Order, error) {
	collection := r.db.Collection(orderCollection)

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
