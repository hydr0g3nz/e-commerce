package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/hydr0g3nz/e-commerce/internal/adapters/model"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	refreshTokenCollection = "refresh_tokens"
	defaultTimeout         = 5 * time.Second
)

type AuthRepository struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Client) *AuthRepository {
	Db := db.Database("e-commerce")
	return &AuthRepository{db: Db}
}

// ensureIndexes creates the required indexes for the refresh tokens collection
func (r *AuthRepository) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(refreshTokenCollection)

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "uuid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "expires_at", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "expires_at", Value: 1},
			},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// StoreRefreshToken stores a new refresh token in MongoDB
func (r *AuthRepository) CreateRefreshToken(userId string, metadata *domain.TokenMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(refreshTokenCollection)

	refreshToken := model.RefreshToken{
		UUID:      metadata.Uuid,
		UserID:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Unix(metadata.Exp, 0),
		IsRevoked: false,
	}
	refreshToken.BeforeCreate()
	_, err := collection.InsertOne(ctx, refreshToken)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("refresh token already exists")
		}
		return err
	}

	return nil
}

// FetchRefreshToken retrieves a refresh token by UUID
func (r *AuthRepository) FetchRefreshToken(uuid string) (*model.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(refreshTokenCollection)

	var token model.RefreshToken
	err := collection.FindOne(ctx, bson.M{
		"uuid":       uuid,
		"is_revoked": false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&token)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("refresh token not found or expired")
		}
		return nil, err
	}

	return &token, nil
}

// DeleteRefreshToken marks a refresh token as revoked
func (r *AuthRepository) DeleteRefreshToken(uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(refreshTokenCollection)

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"uuid": uuid},
		bson.M{"$set": bson.M{"is_revoked": true}},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}

// DeleteUserRefreshTokens revokes all refresh tokens for a specific user
func (r *AuthRepository) DeleteUserRefreshTokens(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(refreshTokenCollection)

	_, err := collection.UpdateMany(
		ctx,
		bson.M{
			"user_id":    userId,
			"is_revoked": false,
		},
		bson.M{"$set": bson.M{"is_revoked": true}},
	)

	return err
}

// CleanupExpiredTokens removes expired tokens from the database
// This is optional as MongoDB TTL index will handle this automatically
func (r *AuthRepository) CleanupExpiredTokens() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(refreshTokenCollection)

	_, err := collection.DeleteMany(
		ctx,
		bson.M{"expires_at": bson.M{"$lt": time.Now()}},
	)

	return err
}

// GetUserActiveTokens retrieves all active refresh tokens for a user
func (r *AuthRepository) GetUserActiveTokens(userId string) ([]model.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(refreshTokenCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"user_id":    userId,
		"is_revoked": false,
		"expires_at": bson.M{"$gt": time.Now()},
	})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tokens []model.RefreshToken
	if err = cursor.All(ctx, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// UpdateRefreshTokenMetadata updates additional metadata for a refresh token
func (r *AuthRepository) UpdateRefreshTokenMetadata(uuid string, device, ipAddress string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection(refreshTokenCollection)

	update := bson.M{
		"$set": bson.M{
			"device":     device,
			"ip_address": ipAddress,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"uuid": uuid},
		update,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}
func (r *AuthRepository) FindByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection("users")

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection("users")
	userModel := model.UserDomainToModel(user)
	userModel.BeforeCreate()
	_, err := collection.InsertOne(ctx, userModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("email already registered")
		}
		return err
	}
	*user = *userModel.Domain()
	return nil
}

func (r *AuthRepository) EmailExists(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	collection := r.db.Collection("users")

	count, err := collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false
	}

	return count > 0
}
