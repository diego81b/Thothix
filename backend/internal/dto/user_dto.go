package dto

import (
	"time"
)

// CreateUserRequest represents the request payload for creating a new user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
}

// UpdateUserRequest represents the request payload for updating user information
type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Username  *string `json:"username,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID        string    `json:"id"`
	ClerkID   string    `json:"clerk_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastSync  time.Time `json:"last_sync"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}

// ClerkUserSyncRequest represents the request for syncing user data from Clerk
type ClerkUserSyncRequest struct {
	ClerkID   string `json:"clerk_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

// ClerkUserSyncResponse represents the response after syncing with Clerk
type ClerkUserSyncResponse struct {
	User    UserResponse `json:"user"`
	IsNew   bool         `json:"is_new"`
	Message string       `json:"message"`
}

// UserSearchRequest represents search criteria for users
type UserSearchRequest struct {
	Query    string `json:"query" form:"query"`
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Page     int    `json:"page" form:"page" validate:"min=1"`
	PerPage  int    `json:"per_page" form:"per_page" validate:"min=1,max=100"`
}

// GetUsersRequest represents request parameters for getting users with pagination
type GetUsersRequest struct {
	Page    int `json:"page" form:"page" validate:"min=1"`
	PerPage int `json:"per_page" form:"per_page" validate:"min=1,max=100"`
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse represents generic success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
