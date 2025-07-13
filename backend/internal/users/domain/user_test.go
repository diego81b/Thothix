package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	commonModels "thothix-backend/internal/common/models"
)

type UserDomainTestSuite struct {
	suite.Suite
}

func (suite *UserDomainTestSuite) TestUserTableName() {
	// Arrange & Act
	user := User{}
	tableName := user.TableName()

	// Assert
	assert.Equal(suite.T(), "users", tableName)
}

func (suite *UserDomainTestSuite) TestSyncFromClerk() {
	// Arrange
	clerkID := "clerk-123"
	user := &User{
		BaseModel: commonModels.BaseModel{
			ID:        "test-id",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ClerkID: &clerkID, // Use pointer
		Email:   "old@example.com",
		Name:    "Old Name",
	}

	clerkData := ClerkUserData{
		Email:     "new@example.com",
		Name:      "New Name",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	// Act
	user.SyncFromClerk(clerkData)

	// Assert
	assert.Equal(suite.T(), "new@example.com", user.Email)
	assert.Equal(suite.T(), "New Name", user.Name)
	assert.Equal(suite.T(), "https://example.com/avatar.jpg", user.AvatarURL)
	assert.NotZero(suite.T(), user.LastSync)
}

func (suite *UserDomainTestSuite) TestClerkUserData() {
	// Arrange & Act
	clerkData := ClerkUserData{
		Email:     "test@example.com",
		Name:      "Test User",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	// Assert
	assert.Equal(suite.T(), "test@example.com", clerkData.Email)
	assert.Equal(suite.T(), "Test User", clerkData.Name)
	assert.Equal(suite.T(), "https://example.com/avatar.jpg", clerkData.AvatarURL)
}

func TestUserDomainTestSuite(t *testing.T) {
	suite.Run(t, new(UserDomainTestSuite))
}
