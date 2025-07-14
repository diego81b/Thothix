package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"thothix-backend/internal/shared/dto"
	sharedTesting "thothix-backend/internal/shared/testing"
	"thothix-backend/internal/users/domain"
	usersDto "thothix-backend/internal/users/dto"
)

type UserServiceTestSuite struct {
	suite.Suite
	container *sharedTesting.PostgresTestContainer
}

func (suite *UserServiceTestSuite) SetupSuite() {
	// Get the shared test container (initialized once per test package)
	suite.container = sharedTesting.GetSharedTestContainer(
		suite.T(),
		"users/service",
		[]interface{}{&domain.User{}},
	)
}

func (suite *UserServiceTestSuite) TearDownSuite() {
	// No cleanup needed - handled by t.Cleanup in getTestContainer
}

func (suite *UserServiceTestSuite) TestNewUserService() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Act
		service := NewUserService(db)

		// Assert
		assert.NotNil(suite.T(), service)
		assert.NotNil(suite.T(), service.db)
		assert.NotNil(suite.T(), service.mapper)
	})
}

func (suite *UserServiceTestSuite) TestGetUserByID_Success() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Arrange - use unique IDs for this test
		testName := "TestGetUserByID_Success"
		user := &domain.User{
			ClerkID: stringPtr("clerk-" + testName), // Use helper function
			Email:   "test-" + testName + "@example.com",
			Name:    "Test User " + testName,
		}
		user.ID = uuid.New().String()
		err := db.Create(user).Error
		assert.NoError(suite.T(), err)

		// Create service with transaction DB
		service := NewUserService(db)

		// Act
		response := service.GetUserByID(user.ID)

		// Assert
		userResponse := sharedTesting.AssertSuccessWithValue(suite.T(), response.Response)
		assert.Equal(suite.T(), "test-"+testName+"@example.com", userResponse.Email)
	})
}

func (suite *UserServiceTestSuite) TestGetUserByID_EmptyID() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Create service with transaction DB
		service := NewUserService(db)

		// Act
		response := service.GetUserByID("")

		// Assert
		sharedTesting.AssertValidationErrorWithCode(suite.T(), response.Response, "VALIDATION_ERROR")
	})
}

func (suite *UserServiceTestSuite) TestGetUserByID_NotFound() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Create service with transaction DB
		service := NewUserService(db)

		// Act
		response := service.GetUserByID(uuid.New().String())

		// Assert
		sharedTesting.AssertValidationErrorWithCode(suite.T(), response.Response, "USER_NOT_FOUND")
	})
}

func (suite *UserServiceTestSuite) TestGetUserByClerkID_Success() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Arrange - use unique data for this test
		testName := "TestGetUserByClerkID_Success"
		user := suite.generateUniqueTestUser(testName)
		user.ID = uuid.New().String()
		err := db.Create(user).Error
		assert.NoError(suite.T(), err)

		// Create service with transaction DB
		service := NewUserService(db)

		// Act
		response := service.GetUserByClerkID(*user.ClerkID) // Dereference pointer

		// Assert
		userResponse := sharedTesting.AssertSuccessWithValue(suite.T(), response.Response)
		assert.Equal(suite.T(), user.Email, userResponse.Email)
	})
}

func (suite *UserServiceTestSuite) TestGetUsers_Success() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Arrange - use unique data for this test
		testName := "TestGetUsers_Success"
		users := []domain.User{
			{ClerkID: stringPtr("clerk-1-" + testName), Email: "user1-" + testName + "@example.com", Name: "User 1 " + testName},
			{ClerkID: stringPtr("clerk-2-" + testName), Email: "user2-" + testName + "@example.com", Name: "User 2 " + testName},
			{ClerkID: stringPtr("clerk-3-" + testName), Email: "user3-" + testName + "@example.com", Name: "User 3 " + testName},
		}
		for _, user := range users {
			user.ID = uuid.New().String()
			err := db.Create(&user).Error
			assert.NoError(suite.T(), err)
		}

		// Create service with transaction DB
		service := NewUserService(db)

		req := &dto.PaginationRequest{
			Page:    1,
			PerPage: 2,
		}

		// Act
		response := service.GetUsers(req)

		// Assert
		usersResponse := sharedTesting.AssertSuccessPaginatedWithValue(suite.T(), response.Response)
		assert.Len(suite.T(), usersResponse.Items, 2)
		assert.Equal(suite.T(), int64(3), usersResponse.Total)
	})
}

func (suite *UserServiceTestSuite) TestCreateUser_Success() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Create service with transaction DB
		service := NewUserService(db)

		// Arrange - use unique data for this test
		testName := "TestCreateUser_Success"
		req := &usersDto.CreateUserRequest{
			Email: "test-" + testName + "@example.com",
			Name:  "Test User " + testName,
			// ClerkID removed - CreateUserRequest is for manual user creation only
		}

		// Act
		response := service.CreateUser(req)

		// Assert
		userResponse := sharedTesting.AssertSuccessWithValue(suite.T(), response.Response)
		assert.Equal(suite.T(), "test-"+testName+"@example.com", userResponse.Email)
		assert.Equal(suite.T(), "Test User "+testName, userResponse.Name)
	})
}

func (suite *UserServiceTestSuite) TestCreateUser_ValidationError() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Create service with transaction DB
		service := NewUserService(db)

		// Arrange
		req := &usersDto.CreateUserRequest{
			Email: "", // Invalid: empty email
			Name:  "Test User",
		}

		// Act
		response := service.CreateUser(req)

		// Assert
		sharedTesting.AssertValidationErrorWithCode(suite.T(), response.Response, "VALIDATION_ERROR")
	})
}

func (suite *UserServiceTestSuite) TestCreateUser_DuplicateEmail() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Arrange - use unique data for this test
		testName := "TestCreateUser_DuplicateEmail"
		existingUser := suite.generateUniqueTestUser(testName)
		// Generate a valid UUID for the existing user
		existingUser.ID = uuid.New().String()
		err := db.Create(existingUser).Error
		assert.NoError(suite.T(), err)

		// Create service with transaction DB
		service := NewUserService(db)

		req := &usersDto.CreateUserRequest{
			Email: existingUser.Email, // Same email as existing user
			Name:  "New User " + testName,
			// ClerkID removed - CreateUserRequest is for manual user creation only
		}

		// Act
		response := service.CreateUser(req)

		// Assert
		sharedTesting.AssertValidationErrorWithCode(suite.T(), response.Response, "CONFLICT")
	})
}

func (suite *UserServiceTestSuite) TestUpdateUser_Success() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Arrange - use unique data for this test
		testName := "TestUpdateUser_Success"
		user := suite.generateUniqueTestUser(testName)
		user.ID = uuid.New().String()
		err := db.Create(user).Error
		assert.NoError(suite.T(), err)

		// Create service with transaction DB
		service := NewUserService(db)

		newEmail := "new-" + testName + "@example.com"
		newName := "New Name " + testName
		req := &usersDto.UpdateUserRequest{
			Email: &newEmail,
			Name:  &newName,
		}

		// Act
		response := service.UpdateUser(user.ID, req)

		// Assert
		userResponse := sharedTesting.AssertSuccessWithValue(suite.T(), response.Response)
		assert.Equal(suite.T(), newEmail, userResponse.Email)
		assert.Equal(suite.T(), newName, userResponse.Name)
	})
}

func (suite *UserServiceTestSuite) TestDeleteUser_Success() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Arrange - use unique data for this test
		testName := "TestDeleteUser_Success"
		user := suite.generateUniqueTestUser(testName)
		user.ID = uuid.New().String()
		err := db.Create(user).Error
		assert.NoError(suite.T(), err)

		// Create service with transaction DB
		service := NewUserService(db)

		// Act
		response := service.DeleteUser(user.ID)

		// Assert
		message := sharedTesting.AssertSuccessWithValue(suite.T(), response.Response)
		assert.Equal(suite.T(), "User deleted successfully", message)

		// Verify user is deleted
		var count int64
		db.Model(&domain.User{}).Where("id = ?", user.ID).Count(&count)
		assert.Equal(suite.T(), int64(0), count)
	})
}

func (suite *UserServiceTestSuite) TestSyncUserFromClerk_NewUser() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Create service with transaction DB
		service := NewUserService(db)

		// Arrange - use unique data for this test
		testName := "TestSyncUserFromClerk_NewUser"
		req := &usersDto.ClerkUserSyncRequest{
			ClerkID:   "clerk-" + testName,
			Email:     "test-" + testName + "@example.com",
			Name:      "Test User " + testName,
			Username:  "testuser" + testName,
			AvatarURL: "https://example.com/avatar-" + testName + ".jpg",
		}

		// Act
		response := service.SyncUserFromClerk(req)

		// Assert
		userResponse := sharedTesting.AssertSuccessWithValue(suite.T(), response.Response)
		assert.Equal(suite.T(), "test-"+testName+"@example.com", userResponse.Email)
		assert.Equal(suite.T(), "Test User "+testName, userResponse.Name)
	})
}

func (suite *UserServiceTestSuite) TestSyncUserFromClerk_ExistingUser() {
	// Use transaction for test isolation
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Arrange - use unique data for this test
		testName := "TestSyncUserFromClerk_ExistingUser"
		existingUser := suite.generateUniqueTestUser(testName)
		existingUser.ID = uuid.New().String()
		existingUser.Email = "old-" + testName + "@example.com"
		existingUser.Name = "Old Name " + testName
		err := db.Create(existingUser).Error
		assert.NoError(suite.T(), err)

		// Create service with transaction DB
		service := NewUserService(db)

		req := &usersDto.ClerkUserSyncRequest{
			ClerkID:   *existingUser.ClerkID, // Dereference pointer to get string
			Email:     "new-" + testName + "@example.com",
			Name:      "New Name " + testName,
			Username:  "newuser" + testName,
			AvatarURL: "https://example.com/new-avatar-" + testName + ".jpg",
		}

		// Act
		response := service.SyncUserFromClerk(req)

		// Assert
		userResponse := sharedTesting.AssertSuccessWithValue(suite.T(), response.Response)
		assert.Equal(suite.T(), "new-"+testName+"@example.com", userResponse.Email)
		assert.Equal(suite.T(), "New Name "+testName, userResponse.Name)
	})
}

// generateUniqueTestUser creates a unique user for testing based on test name
func (suite *UserServiceTestSuite) generateUniqueTestUser(testIdentifier string) *domain.User {
	clerkID := "clerk-" + testIdentifier
	return &domain.User{
		ClerkID: &clerkID, // Convert string to *string
		Email:   "test-" + testIdentifier + "@example.com",
		Name:    "Test User " + testIdentifier,
	}
}

// stringPtr is a helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
