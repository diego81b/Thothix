package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"thothix-backend/internal/shared/dto"
)

type UserDTOTestSuite struct {
	suite.Suite
}

func (suite *UserDTOTestSuite) TestCreateUserRequest() {
	// Arrange & Act
	req := CreateUserRequest{
		Email:     "test@example.com",
		Name:      "Test User",
		Username:  "testuser",
		ClerkID:   "clerk-123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Assert
	assert.Equal(suite.T(), "test@example.com", req.Email)
	assert.Equal(suite.T(), "Test User", req.Name)
	assert.Equal(suite.T(), "testuser", req.Username)
	assert.Equal(suite.T(), "clerk-123", req.ClerkID)
	assert.Equal(suite.T(), "Test", req.FirstName)
	assert.Equal(suite.T(), "User", req.LastName)
}

func (suite *UserDTOTestSuite) TestUpdateUserRequest() {
	// Arrange
	email := "updated@example.com"
	name := "Updated Name"

	// Act
	req := UpdateUserRequest{
		Email: &email,
		Name:  &name,
	}

	// Assert
	assert.NotNil(suite.T(), req.Email)
	assert.Equal(suite.T(), "updated@example.com", *req.Email)
	assert.NotNil(suite.T(), req.Name)
	assert.Equal(suite.T(), "Updated Name", *req.Name)
}

func (suite *UserDTOTestSuite) TestUserResponse() {
	// Arrange & Act
	response := UserResponse{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		ClerkID:   "clerk-123",
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-01-01T00:00:00Z",
	}

	// Assert
	assert.Equal(suite.T(), "test-id", response.ID)
	assert.Equal(suite.T(), "test@example.com", response.Email)
	assert.Equal(suite.T(), "Test User", response.Name)
	assert.Equal(suite.T(), "clerk-123", response.ClerkID)
	assert.Equal(suite.T(), "testuser", response.Username)
	assert.Equal(suite.T(), "https://example.com/avatar.jpg", response.AvatarURL)
}

func (suite *UserDTOTestSuite) TestNewUserListResponse() {
	// Arrange
	users := []UserResponse{
		{ID: "1", Email: "user1@example.com", Name: "User 1"},
		{ID: "2", Email: "user2@example.com", Name: "User 2"},
	}

	// Act
	response := NewUserListResponse(users, 2, 1, 10)

	// Assert
	assert.NotNil(suite.T(), response)
	assert.Len(suite.T(), response.Items, 2)
	assert.Equal(suite.T(), int64(2), response.Total)
	assert.Equal(suite.T(), 1, response.Page)
	assert.Equal(suite.T(), 10, response.PerPage)
	assert.Equal(suite.T(), 1, response.TotalPages)
}

func (suite *UserDTOTestSuite) TestClerkUserSyncRequest() {
	// Arrange & Act
	req := ClerkUserSyncRequest{
		ClerkID:   "clerk-123",
		Email:     "test@example.com",
		Name:      "Test User",
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	// Assert
	assert.Equal(suite.T(), "clerk-123", req.ClerkID)
	assert.Equal(suite.T(), "test@example.com", req.Email)
	assert.Equal(suite.T(), "Test User", req.Name)
	assert.Equal(suite.T(), "testuser", req.Username)
	assert.Equal(suite.T(), "https://example.com/avatar.jpg", req.AvatarURL)
}

func (suite *UserDTOTestSuite) TestClerkUserSyncResponse() {
	// Arrange & Act
	userResponse := UserResponse{
		ID:      "test-id",
		Email:   "test@example.com",
		Name:    "Test User",
		ClerkID: "clerk-123",
	}

	response := ClerkUserSyncResponse{
		User:    userResponse,
		IsNew:   true,
		Message: "User synchronized successfully",
	}

	// Assert
	assert.Equal(suite.T(), userResponse, response.User)
	assert.True(suite.T(), response.IsNew)
	assert.Equal(suite.T(), "User synchronized successfully", response.Message)
}

func (suite *UserDTOTestSuite) TestNewGetUserResponse() {
	// Arrange
	userResponse := &UserResponse{
		ID:    "test-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	producer := func() dto.Validation[*UserResponse] {
		return dto.Success(userResponse)
	}

	// Act
	response := NewGetUserResponse(producer)

	// Assert
	assert.NotNil(suite.T(), response)
	assert.NotNil(suite.T(), response.Response)
}

func (suite *UserDTOTestSuite) TestNewCreateUserResponse() {
	// Arrange
	userResponse := &UserResponse{
		ID:    "test-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	producer := func() dto.Validation[*UserResponse] {
		return dto.Success(userResponse)
	}

	// Act
	response := NewCreateUserResponse(producer)

	// Assert
	assert.NotNil(suite.T(), response)
	assert.NotNil(suite.T(), response.Response)
}

func TestUserDTOTestSuite(t *testing.T) {
	suite.Run(t, new(UserDTOTestSuite))
}
