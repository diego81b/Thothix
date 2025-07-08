package domain

import (
	"time"

	"thothix-backend/internal/models"
)

// User represents a user entity in the users domain
type User struct {
	models.BaseModel
	// Clerk user ID (primary identifier)
	ClerkID    string          `json:"clerk_id" gorm:"uniqueIndex;not null"`
	Email      string          `json:"email"`
	Name       string          `json:"name"`
	Username   string          `json:"username"`
	AvatarURL  string          `json:"avatar_url"`
	SystemRole models.RoleType `json:"system_role" gorm:"default:'user'"` // Default system role
	LastSync   time.Time       `json:"last_sync"`                         // When we last synced with Clerk
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// SyncFromClerk updates user data from Clerk API
func (u *User) SyncFromClerk(clerkUser ClerkUserData) {
	u.Email = clerkUser.Email
	u.Name = clerkUser.Name
	u.AvatarURL = clerkUser.AvatarURL
	u.LastSync = time.Now()
}

// ClerkUserData represents data from Clerk API
type ClerkUserData struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}
