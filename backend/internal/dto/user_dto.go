package dto

// === USER-SPECIFIC DTOs ===

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

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	ClerkID  string `json:"clerk_id,omitempty"`
	Username string `json:"username,omitempty"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	PaginationMeta
}

// ClerkUserSyncRequest represents the request for syncing user data from Clerk
type ClerkUserSyncRequest struct {
	ClerkID   string `json:"clerk_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Name      string `json:"name" validate:"required"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// ClerkUserSyncResponse represents the response after syncing with Clerk
type ClerkUserSyncResponse struct {
	User    UserResponse `json:"user"`
	IsNew   bool         `json:"is_new"`
	Message string       `json:"message"`
}

// === USER-SPECIFIC RESPONSE TYPES ===

type GetUserResponse struct {
	*Response[*UserResponse]
}

func NewGetUserResponse(producer func() Validation[*UserResponse]) *GetUserResponse {
	return &GetUserResponse{
		Response: NewResponse(producer),
	}
}

type GetUsersResponse struct {
	*Response[*UserListResponse]
}

func NewGetUsersResponse(producer func() Validation[*UserListResponse]) *GetUsersResponse {
	return &GetUsersResponse{
		Response: NewResponse(producer),
	}
}

type CreateUserResponse struct {
	*Response[*UserResponse]
}

func NewCreateUserResponse(producer func() Validation[*UserResponse]) *CreateUserResponse {
	return &CreateUserResponse{
		Response: NewResponse(producer),
	}
}

type UpdateUserResponse struct {
	*Response[*UserResponse]
}

func NewUpdateUserResponse(producer func() Validation[*UserResponse]) *UpdateUserResponse {
	return &UpdateUserResponse{
		Response: NewResponse(producer),
	}
}

type DeleteUserResponse struct {
	*Response[string]
}

func NewDeleteUserResponse(producer func() Validation[string]) *DeleteUserResponse {
	return &DeleteUserResponse{
		Response: NewResponse(producer),
	}
}

type ClerkSyncUserResponse struct {
	*Response[*ClerkUserSyncResponse]
}

func NewClerkSyncUserResponse(producer func() Validation[*ClerkUserSyncResponse]) *ClerkSyncUserResponse {
	return &ClerkSyncUserResponse{
		Response: NewResponse(producer),
	}
}
