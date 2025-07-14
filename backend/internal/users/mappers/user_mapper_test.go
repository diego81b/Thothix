package mappers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	commonModels "thothix-backend/internal/common/models"
	"thothix-backend/internal/users/domain"
	usersDto "thothix-backend/internal/users/dto"
)

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

type UserMapperTestSuite struct {
	suite.Suite
	mapper *UserMapper
}

func (suite *UserMapperTestSuite) SetupSuite() {
	suite.mapper = NewUserMapper()
}

func (suite *UserMapperTestSuite) TestNewUserMapper() {
	// Act
	mapper := NewUserMapper()

	// Assert
	assert.NotNil(suite.T(), mapper)
}

func (suite *UserMapperTestSuite) TestModelToDto() {
	// Arrange
	now := time.Now()
	user := &domain.User{
		BaseModel: commonModels.BaseModel{
			ID:        "test-id",
			CreatedAt: now,
			UpdatedAt: now,
		},
		ClerkID:   stringPtr("clerk-123"), // Use helper function for pointer
		Email:     "test@example.com",
		Name:      "Test User",
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	// Act
	dto := suite.mapper.ModelToDto(user)

	// Assert
	assert.NotNil(suite.T(), dto)
	assert.Equal(suite.T(), "test-id", dto.ID)
	assert.Equal(suite.T(), "clerk-123", dto.ClerkID)
	assert.Equal(suite.T(), "test@example.com", dto.Email)
	assert.Equal(suite.T(), "Test User", dto.Name)
	assert.Equal(suite.T(), "testuser", dto.Username)
	assert.Equal(suite.T(), "https://example.com/avatar.jpg", dto.AvatarURL)
	assert.Equal(suite.T(), now.Format(time.RFC3339), dto.CreatedAt)
	assert.Equal(suite.T(), now.Format(time.RFC3339), dto.UpdatedAt)
}

func (suite *UserMapperTestSuite) TestModelToDto_NilInput() {
	// Act
	dto := suite.mapper.ModelToDto(nil)

	// Assert
	assert.Nil(suite.T(), dto)
}

func (suite *UserMapperTestSuite) TestModelsToDtos() {
	// Arrange
	now := time.Now()
	users := []domain.User{
		{
			BaseModel: commonModels.BaseModel{
				ID:        "test-id-1",
				CreatedAt: now,
				UpdatedAt: now,
			},
			Email: "user1@example.com",
			Name:  "User 1",
		},
		{
			BaseModel: commonModels.BaseModel{
				ID:        "test-id-2",
				CreatedAt: now,
				UpdatedAt: now,
			},
			Email: "user2@example.com",
			Name:  "User 2",
		},
	}

	// Act
	dtos := suite.mapper.ModelsToDtos(users)

	// Assert
	assert.NotNil(suite.T(), dtos)
	assert.Len(suite.T(), dtos, 2)
	assert.Equal(suite.T(), "user1@example.com", dtos[0].Email)
	assert.Equal(suite.T(), "user2@example.com", dtos[1].Email)
}

func (suite *UserMapperTestSuite) TestModelsToDtos_NilInput() {
	// Act
	dtos := suite.mapper.ModelsToDtos(nil)

	// Assert
	assert.Nil(suite.T(), dtos)
}

func (suite *UserMapperTestSuite) TestCreateRequestToModel() {
	// Arrange
	req := &usersDto.CreateUserRequest{
		Email:     "test@example.com",
		Name:      "Test User",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}

	// Act
	user := suite.mapper.CreateRequestToModel(req)

	// Assert
	assert.NotNil(suite.T(), user)
	assert.NotEmpty(suite.T(), user.ID)
	assert.Nil(suite.T(), user.ClerkID) // ClerkID should be nil for manual user creation
	assert.Equal(suite.T(), "test@example.com", user.Email)
	assert.Equal(suite.T(), "Test User", user.Name)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.NotZero(suite.T(), user.CreatedAt)
	assert.NotZero(suite.T(), user.UpdatedAt)
}

func (suite *UserMapperTestSuite) TestCreateRequestToModel_NilInput() {
	// Act
	user := suite.mapper.CreateRequestToModel(nil)

	// Assert
	assert.Nil(suite.T(), user)
}

func (suite *UserMapperTestSuite) TestUpdateRequestToModel() {
	// Arrange
	now := time.Now()
	user := &domain.User{
		BaseModel: commonModels.BaseModel{
			ID:        "test-id",
			CreatedAt: now,
			UpdatedAt: now,
		},
		Email: "old@example.com",
		Name:  "Old Name",
	}

	newEmail := "new@example.com"
	newName := "New Name"
	newUsername := "newuser"
	newAvatarURL := "https://example.com/new-avatar.jpg"

	req := &usersDto.UpdateUserRequest{
		Email:     &newEmail,
		Name:      &newName,
		Username:  &newUsername,
		AvatarURL: &newAvatarURL,
	}

	// Act
	suite.mapper.UpdateRequestToModel(user, req)

	// Assert
	assert.Equal(suite.T(), "new@example.com", user.Email)
	assert.Equal(suite.T(), "New Name", user.Name)
	assert.Equal(suite.T(), "newuser", user.Username)
	assert.Equal(suite.T(), "https://example.com/new-avatar.jpg", user.AvatarURL)
	assert.True(suite.T(), user.UpdatedAt.After(now) || user.UpdatedAt.Equal(now))
}

func (suite *UserMapperTestSuite) TestUpdateRequestToModel_NilInputs() {
	// Act
	suite.mapper.UpdateRequestToModel(nil, nil)

	// Assert - should not panic
}

func (suite *UserMapperTestSuite) TestClerkSyncRequestToModel() {
	// Arrange
	req := &usersDto.ClerkUserSyncRequest{
		ClerkID:   "clerk-123",
		Email:     "test@example.com",
		Name:      "Test User",
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	// Act
	user := suite.mapper.ClerkSyncRequestToModel(req)

	// Assert
	assert.NotNil(suite.T(), user)
	assert.NotEmpty(suite.T(), user.ID)
	assert.NotNil(suite.T(), user.ClerkID)
	assert.Equal(suite.T(), "clerk-123", *user.ClerkID) // Dereference pointer
	assert.Equal(suite.T(), "test@example.com", user.Email)
	assert.Equal(suite.T(), "Test User", user.Name)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.Equal(suite.T(), "https://example.com/avatar.jpg", user.AvatarURL)
	assert.NotZero(suite.T(), user.CreatedAt)
	assert.NotZero(suite.T(), user.UpdatedAt)
	assert.NotZero(suite.T(), user.LastSync)
}

func (suite *UserMapperTestSuite) TestClerkSyncRequestToModel_NilInput() {
	// Act
	user := suite.mapper.ClerkSyncRequestToModel(nil)

	// Assert
	assert.Nil(suite.T(), user)
}

func TestUserMapperTestSuite(t *testing.T) {
	suite.Run(t, new(UserMapperTestSuite))
}
