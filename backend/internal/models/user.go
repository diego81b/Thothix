package models

import (
	"time"
)

// User represents a user entity
// swagger:model User
type User struct {
	BaseModel
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	AvatarURL  string   `json:"avatar_url"`
	SystemRole RoleType `json:"system_role" gorm:"default:'user'"` // Default system role
}

// UserResponse represents the API response for a user
// swagger:model UserResponse
type UserResponse struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	SystemRole RoleType  `json:"system_role"`
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatar_url"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Name:       u.Name,
		AvatarURL:  u.AvatarURL,
		SystemRole: u.SystemRole,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
