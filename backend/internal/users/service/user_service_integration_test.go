package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"thothix-backend/internal/shared/dto"
	sharedTesting "thothix-backend/internal/shared/testing"
	"thothix-backend/internal/users/domain"
	usersDto "thothix-backend/internal/users/dto"
)

// UserServiceIntegrationTestSuite tests the service layer integration with real database
// These tests validate complex service operations that involve multiple database operations
type UserServiceIntegrationTestSuite struct {
	suite.Suite
	container *sharedTesting.PostgresTestContainer
}

func (suite *UserServiceIntegrationTestSuite) SetupSuite() {
	suite.container = sharedTesting.GetSharedTestContainer(
		suite.T(),
		"users/service_integration",
		[]interface{}{&domain.User{}},
	)
}

func (suite *UserServiceIntegrationTestSuite) TearDownSuite() {
	// Cleanup handled by t.Cleanup in GetSharedTestContainer
}

// TestBulkUserOperations tests service performance with multiple users
// This validates that the service can handle bulk operations efficiently
func (suite *UserServiceIntegrationTestSuite) TestBulkUserOperations() {
	suite.container.WithTransaction(func(db *gorm.DB) {
		service := NewUserService(db)
		testName := "BulkUserOperations"

		// Create multiple users
		userIDs := make([]string, 10)
		for i := 0; i < 10; i++ {
			idx := fmt.Sprintf("%d", i+1)
			req := &usersDto.CreateUserRequest{
				Email: "bulk-" + testName + "-" + idx + "@example.com",
				Name:  "Bulk User " + testName + " " + idx,
				// ClerkID removed - CreateUserRequest is for manual user creation only
			}

			response := service.CreateUser(req)
			userResponse := sharedTesting.AssertSuccessWithValue(suite.T(), response.Response)
			userIDs[i] = userResponse.ID
		}

		// Test pagination with large dataset
		paginationReq := &dto.PaginationRequest{
			Page:    1,
			PerPage: 5,
		}

		response := service.GetUsers(paginationReq)
		usersResponse := sharedTesting.AssertSuccessPaginatedWithValue(suite.T(), response.Response)

		assert.Len(suite.T(), usersResponse.Items, 5)
		assert.Equal(suite.T(), int64(10), usersResponse.Total)
		assert.Equal(suite.T(), 1, usersResponse.Page)
		assert.Equal(suite.T(), 5, usersResponse.PerPage)

		// Test second page
		paginationReq.Page = 2
		response = service.GetUsers(paginationReq)
		usersResponse = sharedTesting.AssertSuccessPaginatedWithValue(suite.T(), response.Response)

		assert.Len(suite.T(), usersResponse.Items, 5)
		assert.Equal(suite.T(), int64(10), usersResponse.Total)
		assert.Equal(suite.T(), 2, usersResponse.Page)
	})
}

// TestUserDataConsistency tests complex scenarios that involve multiple service operations
// This validates data consistency across operations
func (suite *UserServiceIntegrationTestSuite) TestUserDataConsistency() {
	suite.container.WithTransaction(func(db *gorm.DB) {
		service := NewUserService(db)
		testName := "DataConsistency"

		// Create user
		createReq := &usersDto.CreateUserRequest{
			Email: "consistency-" + testName + "@example.com",
			Name:  "Consistency User " + testName,
		}

		createResponse := service.CreateUser(createReq)
		user := sharedTesting.AssertSuccessWithValue(suite.T(), createResponse.Response)

		// Update same user multiple times and verify consistency
		updates := []struct {
			email string
			name  string
		}{
			{"update1-" + testName + "@example.com", "Update 1 " + testName},
			{"update2-" + testName + "@example.com", "Update 2 " + testName},
			{"update3-" + testName + "@example.com", "Update 3 " + testName},
		}

		for i, update := range updates {
			updateReq := &usersDto.UpdateUserRequest{
				Email: &update.email,
				Name:  &update.name,
			}

			updateResponse := service.UpdateUser(user.ID, updateReq)
			updatedUser := sharedTesting.AssertSuccessWithValue(suite.T(), updateResponse.Response)

			// Verify update was applied
			assert.Equal(suite.T(), update.email, updatedUser.Email)
			assert.Equal(suite.T(), update.name, updatedUser.Name)

			// Verify by fetching the user again
			getResponse := service.GetUserByID(user.ID)
			fetchedUser := sharedTesting.AssertSuccessWithValue(suite.T(), getResponse.Response)

			assert.Equal(suite.T(), update.email, fetchedUser.Email)
			assert.Equal(suite.T(), update.name, fetchedUser.Name)

			// Log progress
			suite.T().Logf("Completed update %d/%d: %s", i+1, len(updates), update.email)
		}

		// Final verification: get all users and verify our user is in the list with final state
		paginationReq := &dto.PaginationRequest{
			Page:    1,
			PerPage: 10,
		}

		listResponse := service.GetUsers(paginationReq)
		usersList := sharedTesting.AssertSuccessPaginatedWithValue(suite.T(), listResponse.Response)

		// Find our user in the list
		var foundUser *usersDto.UserDto
		for _, u := range usersList.Items {
			if u.ID == user.ID {
				foundUser = &u
				break
			}
		}

		assert.NotNil(suite.T(), foundUser, "User should be found in the list")
		assert.Equal(suite.T(), updates[len(updates)-1].email, foundUser.Email)
		assert.Equal(suite.T(), updates[len(updates)-1].name, foundUser.Name)
	})
}

func TestUserServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceIntegrationTestSuite))
}
