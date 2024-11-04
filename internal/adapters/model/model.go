package model

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        string     `bson:"_id"`
	UpdatedAt time.Time  `bson:"updated_at"`
	CreatedAt time.Time  `bson:"created_at"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

func (m *Model) SetCreatedAt() {
	m.CreatedAt = time.Now()
}
func (m *Model) SetUpdatedAt() {
	m.UpdatedAt = time.Now()
}
func (m *Model) SetDeletedAt() {
	m.DeletedAt = new(time.Time)
	*m.DeletedAt = time.Now()
}
func (m *Model) SetID() {
	id, _ := uuid.NewV7()
	m.ID = id.String()
}
func (m *Model) BeforeCreate() {
	m.SetCreatedAt()
	m.SetUpdatedAt()
	m.SetID()
}
func (m *Model) BeforeUpdate() {
	m.SetUpdatedAt()
}
func (m *Model) BeforeDelete() {
	m.SetDeletedAt()
}
