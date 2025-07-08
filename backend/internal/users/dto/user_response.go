package dto

import (
	"thothix-backend/internal/shared/dto"
)

// === USER RESPONSE DTOs ===
// These DTOs wrap responses with validation and error handling

// GetUserResponse wraps a single UserDto response
type GetUserResponse struct {
	*dto.Response[*UserDto]
}

func NewGetUserResponse(producer func() dto.Validation[*UserDto]) *GetUserResponse {
	return &GetUserResponse{
		Response: dto.NewResponse(producer),
	}
}

// GetUsersResponse wraps a paginated list of UserDto
type GetUsersResponse = dto.ListResponse[UserDto]

func NewGetUsersResponse(producer func() dto.Validation[*UserListDto]) *GetUsersResponse {
	typedProducer := func() dto.Validation[*dto.PaginatedListResponse[UserDto]] {
		result := producer()
		if result.IsValid() {
			return dto.Valid[*dto.PaginatedListResponse[UserDto]](result.GetValue())
		}
		return dto.Invalid[*dto.PaginatedListResponse[UserDto]](result.GetErrors()...)
	}
	return dto.NewListResponse(typedProducer)
}

// CreateUserResponse wraps a newly created UserDto
type CreateUserResponse struct {
	*dto.Response[*UserDto]
}

func NewCreateUserResponse(producer func() dto.Validation[*UserDto]) *CreateUserResponse {
	return &CreateUserResponse{
		Response: dto.NewResponse(producer),
	}
}

// UpdateUserResponse wraps an updated UserDto
type UpdateUserResponse struct {
	*dto.Response[*UserDto]
}

func NewUpdateUserResponse(producer func() dto.Validation[*UserDto]) *UpdateUserResponse {
	return &UpdateUserResponse{
		Response: dto.NewResponse(producer),
	}
}

// DeleteUserResponse wraps a deletion confirmation message
type DeleteUserResponse struct {
	*dto.Response[string]
}

func NewDeleteUserResponse(producer func() dto.Validation[string]) *DeleteUserResponse {
	return &DeleteUserResponse{
		Response: dto.NewResponse(producer),
	}
}

// ClerkSyncUserResponse wraps the result of a Clerk user sync operation
type ClerkSyncUserResponse struct {
	*dto.Response[*ClerkUserSyncDto]
}

func NewClerkSyncUserResponse(producer func() dto.Validation[*ClerkUserSyncDto]) *ClerkSyncUserResponse {
	return &ClerkSyncUserResponse{
		Response: dto.NewResponse(producer),
	}
}
