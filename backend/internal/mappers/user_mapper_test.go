package mappers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"thothix-backend/internal/dto"
	"thothix-backend/internal/models"
)

func TestUserMapper_ModelToResponse(t *testing.T) {
	mapper := NewUserMapper()

	tests := []struct {
		name     string
		input    *models.User
		expected *dto.UserResponse
	}{
		{
			name: "valid user",
			input: &models.User{
				BaseModel: models.BaseModel{
					ID:        "user123",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				ClerkID:   "clerk_123",
				Email:     "test@example.com",
				Name:      "Test User",
				Username:  "testuser",
				AvatarURL: "https://example.com/avatar.jpg",
				LastSync:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			expected: &dto.UserResponse{
				ID:        "user123",
				ClerkID:   "clerk_123",
				Email:     "test@example.com",
				Name:      "Test User",
				Username:  "testuser",
				AvatarURL: "https://example.com/avatar.jpg",
				CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				LastSync:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ModelToResponse(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserMapper_ModelsToResponses(t *testing.T) {
	mapper := NewUserMapper()

	tests := []struct {
		name     string
		input    []models.User
		expected []dto.UserResponse
	}{
		{
			name: "multiple users",
			input: []models.User{
				{
					BaseModel: models.BaseModel{ID: "user1"},
					ClerkID:   "clerk_1",
					Email:     "user1@example.com",
					Name:      "User One",
				},
				{
					BaseModel: models.BaseModel{ID: "user2"},
					ClerkID:   "clerk_2",
					Email:     "user2@example.com",
					Name:      "User Two",
				},
			},
			expected: []dto.UserResponse{
				{
					ID:      "user1",
					ClerkID: "clerk_1",
					Email:   "user1@example.com",
					Name:    "User One",
				},
				{
					ID:      "user2",
					ClerkID: "clerk_2",
					Email:   "user2@example.com",
					Name:    "User Two",
				},
			},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty slice",
			input:    []models.User{},
			expected: []dto.UserResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ModelsToResponses(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserMapper_CreateRequestToModel(t *testing.T) {
	mapper := NewUserMapper()

	tests := []struct {
		name     string
		input    *dto.CreateUserRequest
		expected *models.User
	}{
		{
			name: "valid request",
			input: &dto.CreateUserRequest{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Username:  "johndoe",
			},
			expected: &models.User{
				Email:      "test@example.com",
				Name:       "John Doe",
				Username:   "johndoe",
				SystemRole: models.RoleUser,
			},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.CreateRequestToModel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserMapper_ClerkSyncRequestToModel(t *testing.T) {
	mapper := NewUserMapper()

	tests := []struct {
		name     string
		input    *dto.ClerkUserSyncRequest
		expected *models.User
	}{
		{
			name: "valid request",
			input: &dto.ClerkUserSyncRequest{
				ClerkID:   "clerk_123",
				Email:     "test@example.com",
				Name:      "Test User",
				Username:  "testuser",
				AvatarURL: "https://example.com/avatar.jpg",
			},
			expected: &models.User{
				ClerkID:    "clerk_123",
				Email:      "test@example.com",
				Name:       "Test User",
				Username:   "testuser",
				AvatarURL:  "https://example.com/avatar.jpg",
				SystemRole: models.RoleUser,
			},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ClerkSyncRequestToModel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserMapper_UpdateRequestToMap(t *testing.T) {
	mapper := NewUserMapper()

	tests := []struct {
		name     string
		input    *dto.UpdateUserRequest
		expected map[string]interface{}
	}{
		{
			name: "all fields",
			input: &dto.UpdateUserRequest{
				Email:     stringPtr("new@example.com"),
				FirstName: stringPtr("New"),
				LastName:  stringPtr("Name"),
				Username:  stringPtr("newusername"),
				AvatarURL: stringPtr("https://example.com/new-avatar.jpg"),
			},
			expected: map[string]interface{}{
				"email":      "new@example.com",
				"name":       "New Name",
				"username":   "newusername",
				"avatar_url": "https://example.com/new-avatar.jpg",
			},
		},
		{
			name: "partial fields",
			input: &dto.UpdateUserRequest{
				Email:    stringPtr("partial@example.com"),
				Username: stringPtr("partialuser"),
			},
			expected: map[string]interface{}{
				"email":    "partial@example.com",
				"username": "partialuser",
			},
		},
		{
			name: "only first name",
			input: &dto.UpdateUserRequest{
				FirstName: stringPtr("OnlyFirst"),
			},
			expected: map[string]interface{}{
				"name": "OnlyFirst",
			},
		},
		{
			name: "only last name",
			input: &dto.UpdateUserRequest{
				LastName: stringPtr("OnlyLast"),
			},
			expected: map[string]interface{}{
				"name": " OnlyLast",
			},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty request",
			input:    &dto.UpdateUserRequest{},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.UpdateRequestToMap(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserMapper_ModelsToListResponse(t *testing.T) {
	mapper := NewUserMapper()

	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: "user1"},
			ClerkID:   "clerk_1",
			Email:     "user1@example.com",
			Name:      "User One",
		},
		{
			BaseModel: models.BaseModel{ID: "user2"},
			ClerkID:   "clerk_2",
			Email:     "user2@example.com",
			Name:      "User Two",
		},
	}

	result := mapper.ModelsToListResponse(users, 10, 1, 5)

	require.NotNil(t, result)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, int64(10), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 5, result.PerPage)
	assert.Equal(t, 2, result.TotalPages) // ceil(10/5)

	assert.Equal(t, "user1", result.Users[0].ID)
	assert.Equal(t, "user2", result.Users[1].ID)
}

func TestUserMapper_CreateSyncResponse(t *testing.T) {
	mapper := NewUserMapper()

	user := &models.User{
		BaseModel: models.BaseModel{ID: "user123"},
		ClerkID:   "clerk_123",
		Email:     "test@example.com",
		Name:      "Test User",
	}

	tests := []struct {
		name     string
		user     *models.User
		isNew    bool
		message  string
		expected *dto.ClerkUserSyncResponse
	}{
		{
			name:    "new user",
			user:    user,
			isNew:   true,
			message: "User created",
			expected: &dto.ClerkUserSyncResponse{
				User: dto.UserResponse{
					ID:      "user123",
					ClerkID: "clerk_123",
					Email:   "test@example.com",
					Name:    "Test User",
				},
				IsNew:   true,
				Message: "User created",
			},
		},
		{
			name:     "nil user",
			user:     nil,
			isNew:    false,
			message:  "Error",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.CreateSyncResponse(tt.user, tt.isNew, tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function for tests
func stringPtr(s string) *string {
	return &s
}
