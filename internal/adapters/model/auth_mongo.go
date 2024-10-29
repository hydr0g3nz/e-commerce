package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
)

type RefreshToken struct {
	ID        string    `bson:"_id"`
	UUID      string    `bson:"uuid"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	IsRevoked bool      `bson:"is_revoked"`
}

func (r *RefreshToken) SetID() {
	id, _ := uuid.NewV7()
	r.ID = id.String()
}
func (r *RefreshToken) BeforeCreate() {
	r.SetID()
}
func UserDomainToModel(user *domain.User) *User {
	return &User{
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Name:     user.Name,
		Address:  user.Address,
	}
}

type User struct {
	Model    `bson:",inline"`
	Email    string           `bson:"email"`
	Password string           `bson:"password"`
	Role     string           `bson:"role"`
	Name     string           `bson:"name"`
	Address  []domain.Address `bson:"address"`
}

func (u *User) BeforeCreate() {
	u.Model.BeforeCreate()
	if u.Role == "" {
		u.Role = "user"
	}
}

func (u *User) Domain() *domain.User {
	return &domain.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
		Name:     u.Name,
		Address:  u.Address,
	}
}
