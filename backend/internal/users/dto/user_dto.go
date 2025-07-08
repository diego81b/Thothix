package dto

import (
	"thothix-backend/internal/shared/dto"
)

// === USER-SPECIFIC REQUEST DTOs ===

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

// === USER-SPECIFIC RESPONSE DTOs ===

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	ClerkID   string `json:"clerk_id,omitempty"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UserListResponse represents paginated user list response
type UserListResponse = dto.PaginatedListResponse[UserResponse]

// NewUserListResponse creates a UserListResponse with proper pagination metadata
func NewUserListResponse(users []UserResponse, total int64, page, perPage int) *UserListResponse {
	return dto.NewPaginatedListResponse(users, total, page, perPage)
}

// ClerkUserSyncResponse represents the response after syncing with Clerk
type ClerkUserSyncResponse struct {
	User    UserResponse `json:"user"`
	IsNew   bool         `json:"is_new"`
	Message string       `json:"message"`
}

// === USER-SPECIFIC RESPONSE TYPES ===

type GetUserResponse struct {
	*dto.Response[*UserResponse]
}

func NewGetUserResponse(producer func() dto.Validation[*UserResponse]) *GetUserResponse {
	return &GetUserResponse{
		Response: dto.NewResponse(producer),
	}
}

type GetUsersResponse = dto.ListResponse[UserResponse]

func NewGetUsersResponse(producer func() dto.Validation[*UserListResponse]) *GetUsersResponse {
	typedProducer := func() dto.Validation[*dto.PaginatedListResponse[UserResponse]] {
		result := producer()
		if result.IsValid() {
			return dto.Valid[*dto.PaginatedListResponse[UserResponse]](result.GetValue())
		}
		return dto.Invalid[*dto.PaginatedListResponse[UserResponse]](result.GetErrors()...)
	}
	return dto.NewListResponse(typedProducer)
}

type CreateUserResponse struct {
	*dto.Response[*UserResponse]
}

func NewCreateUserResponse(producer func() dto.Validation[*UserResponse]) *CreateUserResponse {
	return &CreateUserResponse{
		Response: dto.NewResponse(producer),
	}
}

type UpdateUserResponse struct {
	*dto.Response[*UserResponse]
}

func NewUpdateUserResponse(producer func() dto.Validation[*UserResponse]) *UpdateUserResponse {
	return &UpdateUserResponse{
		Response: dto.NewResponse(producer),
	}
}

type DeleteUserResponse struct {
	*dto.Response[string]
}

func NewDeleteUserResponse(producer func() dto.Validation[string]) *DeleteUserResponse {
	return &DeleteUserResponse{
		Response: dto.NewResponse(producer),
	}
}

type ClerkSyncUserResponse struct {
	*dto.Response[*ClerkUserSyncResponse]
}

func NewClerkSyncUserResponse(producer func() dto.Validation[*ClerkUserSyncResponse]) *ClerkSyncUserResponse {
	return &ClerkSyncUserResponse{
		Response: dto.NewResponse(producer),
	}
}
