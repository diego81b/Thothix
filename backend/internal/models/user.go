package models

import (
	"time"
)

// User represents a user entity
// swagger:model User
type User struct {
	BaseModel
	// Clerk user ID (primary identifier)
	ClerkID    string    `json:"clerk_id" gorm:"uniqueIndex;not null"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	AvatarURL  string    `json:"avatar_url"`
	SystemRole RoleType  `json:"system_role" gorm:"default:'user'"` // Default system role
	LastSync   time.Time `json:"last_sync"`                         // When we last synced with Clerk
}

// UserResponse represents the API response for a user
// swagger:model UserResponse
type UserResponse struct {
	ID         string    `json:"id"`
	ClerkID    string    `json:"clerk_id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatar_url"`
	SystemRole RoleType  `json:"system_role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		ClerkID:    u.ClerkID,
		Email:      u.Email,
		Name:       u.Name,
		AvatarURL:  u.AvatarURL,
		SystemRole: u.SystemRole,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
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
