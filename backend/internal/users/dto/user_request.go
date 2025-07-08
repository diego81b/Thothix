package dto

// === USER REQUEST DTOs ===
// These DTOs represent incoming requests from clients/APIs

// CreateUserRequest represents the request payload for creating a new user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Username  string `json:"username,omitempty"`
	ClerkID   string `json:"clerk_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// UpdateUserRequest represents the request payload for updating user information
type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Name      *string `json:"name,omitempty"`
	Username  *string `json:"username,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// ClerkUserSyncRequest represents the request for syncing user data from Clerk
type ClerkUserSyncRequest struct {
	ClerkID   string `json:"clerk_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}
