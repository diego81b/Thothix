package mappers

import (
	"time"

	"github.com/google/uuid"
	"thothix-backend/internal/users/domain"
	usersDto "thothix-backend/internal/users/dto"
)

// UserMapper handles conversion between User models and DTOs
type UserMapper struct{}

// NewUserMapper creates a new UserMapper instance
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ModelToResponse converts a User model to UserResponse DTO
func (m *UserMapper) ModelToResponse(user *domain.User) *usersDto.UserResponse {
	if user == nil {
		return nil
	}

	return &usersDto.UserResponse{
		ID:        user.ID,
		ClerkID:   user.ClerkID,
		Email:     user.Email,
		Name:      user.Name,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

// ModelsToResponses converts a slice of User models to UserResponse DTOs
func (m *UserMapper) ModelsToResponses(users []domain.User) []usersDto.UserResponse {
	if users == nil {
		return nil
	}

	responses := make([]usersDto.UserResponse, len(users))
	for i, user := range users {
		response := m.ModelToResponse(&user)
		if response != nil {
			responses[i] = *response
		}
	}

	return responses
}

// CreateRequestToModel converts CreateUserRequest DTO to User model
func (m *UserMapper) CreateRequestToModel(req *usersDto.CreateUserRequest) *domain.User {
	if req == nil {
		return nil
	}

	user := &domain.User{
		ClerkID:  req.ClerkID,
		Email:    req.Email,
		Name:     req.Name,
		Username: req.Username,
	}

	// Set base model fields
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return user
}

// UpdateRequestToModel applies UpdateUserRequest changes to existing User model
func (m *UserMapper) UpdateRequestToModel(user *domain.User, req *usersDto.UpdateUserRequest) {
	if user == nil || req == nil {
		return
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}

	user.UpdatedAt = time.Now()
}

// ClerkSyncRequestToModel converts ClerkUserSyncRequest to User model
func (m *UserMapper) ClerkSyncRequestToModel(req *usersDto.ClerkUserSyncRequest) *domain.User {
	if req == nil {
		return nil
	}

	user := &domain.User{
		ClerkID:   req.ClerkID,
		Email:     req.Email,
		Name:      req.Name,
		Username:  req.Username,
		AvatarURL: req.AvatarURL,
		LastSync:  time.Now(),
	}

	// Set base model fields
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return user
}
