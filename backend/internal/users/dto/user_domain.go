package dto

import (
	"thothix-backend/internal/shared/dto"
)

// === USER DOMAIN DTOs ===
// These DTOs represent domain objects for mapping models

// UserDto represents the user data structure (domain mapping)
type UserDto struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	ClerkID   string `json:"clerk_id,omitempty"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UserListDto represents paginated user list data
type UserListDto = dto.PaginatedListResponse[UserDto]

// NewUserListDto creates a UserListDto with proper pagination metadata
func NewUserListDto(users []UserDto, total int64, page, perPage int) *UserListDto {
	return dto.NewPaginatedListResponse(users, total, page, perPage)
}

// ClerkUserSyncDto represents the data after syncing with Clerk
type ClerkUserSyncDto struct {
	User    UserDto `json:"user"`
	IsNew   bool    `json:"is_new"`
	Message string  `json:"message"`
}
