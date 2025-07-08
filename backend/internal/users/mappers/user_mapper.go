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

// ModelToDto converts a User model to UserDto
func (m *UserMapper) ModelToDto(user *domain.User) *usersDto.UserDto {
	if user == nil {
		return nil
	}

	return &usersDto.UserDto{
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

// ModelsToDtos converts a slice of User models to UserDto DTOs
func (m *UserMapper) ModelsToDtos(users []domain.User) []usersDto.UserDto {
	if users == nil {
		return nil
	}

	dtos := make([]usersDto.UserDto, len(users))
	for i, user := range users {
		dto := m.ModelToDto(&user)
		if dto != nil {
			dtos[i] = *dto
		}
	}

	return dtos
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
