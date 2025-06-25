package mappers

import (
	"math"

	"thothix-backend/internal/dto"
	"thothix-backend/internal/models"
)

// UserMapper handles conversion between User models and DTOs
type UserMapper struct{}

// NewUserMapper creates a new UserMapper instance
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ModelToResponse converts a User model to UserResponse DTO
func (m *UserMapper) ModelToResponse(user *models.User) *dto.UserResponse {
	if user == nil {
		return nil
	}

	return &dto.UserResponse{
		ID:       user.ID,
		ClerkID:  user.ClerkID,
		Email:    user.Email,
		Name:     user.Name,
		Username: user.Username,
	}
}

// ModelsToResponses converts a slice of User models to UserResponse DTOs
func (m *UserMapper) ModelsToResponses(users []models.User) []dto.UserResponse {
	if users == nil {
		return nil
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response := m.ModelToResponse(&user)
		if response != nil {
			responses[i] = *response
		}
	}

	return responses
}

// CreateRequestToModel converts CreateUserRequest DTO to User model
func (m *UserMapper) CreateRequestToModel(req *dto.CreateUserRequest) *models.User {
	if req == nil {
		return nil
	}

	return &models.User{
		Email:      req.Email,
		Name:       req.FirstName + " " + req.LastName,
		Username:   req.Username,
		SystemRole: models.RoleUser,
	}
}

// ClerkSyncRequestToModel converts ClerkUserSyncRequest DTO to User model
func (m *UserMapper) ClerkSyncRequestToModel(req *dto.ClerkUserSyncRequest) *models.User {
	if req == nil {
		return nil
	}

	return &models.User{
		ClerkID:    req.ClerkID,
		Email:      req.Email,
		Name:       req.Name,
		Username:   req.Username,
		AvatarURL:  req.AvatarURL,
		SystemRole: models.RoleUser,
	}
}

// UpdateRequestToMap converts UpdateUserRequest DTO to a map for GORM updates
func (m *UserMapper) UpdateRequestToMap(req *dto.UpdateUserRequest) map[string]interface{} {
	if req == nil {
		return nil
	}

	updates := make(map[string]interface{})

	if req.Email != nil {
		updates["email"] = *req.Email
	}

	if req.FirstName != nil || req.LastName != nil {
		// If either first or last name is being updated, we need to reconstruct the full name
		// Note: This assumes we have access to the current user data to fill missing parts
		// In a real implementation, you might want to store first/last name separately
		name := ""
		if req.FirstName != nil {
			name = *req.FirstName
		}
		if req.LastName != nil {
			if req.FirstName != nil {
				// Both first and last name provided
				name += " " + *req.LastName
			} else {
				// Only last name provided - add space prefix to maintain format
				name = " " + *req.LastName
			}
		}
		if name != "" {
			updates["name"] = name
		}
	}

	if req.Username != nil {
		updates["username"] = *req.Username
	}

	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}

	return updates
}

// ModelsToListResponse converts models and pagination info to UserListResponse
func (m *UserMapper) ModelsToListResponse(users []models.User, total int64, page, perPage int) *dto.UserListResponse {
	userResponses := m.ModelsToResponses(users)
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.UserListResponse{
		Users: userResponses,
		PaginationMeta: dto.PaginationMeta{
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		},
	}
}

// CreateSyncResponse creates a ClerkUserSyncResponse
func (m *UserMapper) CreateSyncResponse(user *models.User, isNew bool, message string) *dto.ClerkUserSyncResponse {
	if user == nil {
		return nil
	}

	userResponse := m.ModelToResponse(user)
	return &dto.ClerkUserSyncResponse{
		User:    *userResponse,
		IsNew:   isNew,
		Message: message,
	}
}
