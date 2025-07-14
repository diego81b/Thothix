package domain

import (
	"time"

	commonModels "thothix-backend/internal/common/models"
	sharedModels "thothix-backend/internal/shared/models"
)

// User represents a user entity in the users domain
type User struct {
	commonModels.BaseModel
	// Clerk user ID (optional - NULL for manually created users)
	ClerkID    *string               `json:"clerk_id" gorm:"uniqueIndex"` // NULL for manual users, unique for Clerk users
	Email      string                `json:"email"`
	Name       string                `json:"name"`
	Username   string                `json:"username"`
	AvatarURL  string                `json:"avatar_url"`
	SystemRole sharedModels.RoleType `json:"system_role" gorm:"default:'user'"` // Default system role
	LastSync   time.Time             `json:"last_sync"`                         // When we last synced with Clerk
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// ClerkUserData represents user data from Clerk
type ClerkUserData struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// SyncFromClerk updates user data from Clerk
func (u *User) SyncFromClerk(data ClerkUserData) {
	u.Email = data.Email
	u.Name = data.Name
	u.AvatarURL = data.AvatarURL
	u.LastSync = time.Now()
}
