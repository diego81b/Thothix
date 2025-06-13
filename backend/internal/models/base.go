package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel provides common fields for all entities
type BaseModel struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to set ID, timestamps and creator automatically
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = generateID()
	}
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now

	// Set CreatedBy from context if not already set
	if b.CreatedBy == "" {
		if userID := tx.Statement.Context.Value("user_id"); userID != nil {
			if uid, ok := userID.(string); ok {
				b.CreatedBy = uid
			}
		}
	}
	return nil
}

// BeforeUpdate hook to update timestamps
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}

// generateID generates a unique ID using UUID
func generateID() string {
	return uuid.New().String()
}
