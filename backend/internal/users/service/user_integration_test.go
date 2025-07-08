package service

import (
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

func (suite *UserIntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.container = sharedTesting.GetSharedTestContainer(
		suite.T(),
		"users/integration",
		[]interface{}{&domain.User{}},
	)
}

func (suite *UserIntegrationTestSuite) TearDownSuite() {
	// Cleanup handled by t.Cleanup in GetSharedTestContainer
}

// =====================================
// SERVICE INTEGRATION TESTS
// =====================================

// TestBulkUserOperations tests service performance with multiple users
// This validates that the service can handle bulk operations efficiently
func (suite *UserIntegrationTestSuite) TestBulkUserOperations() {
	suite.container.WithTransaction(func(db *gorm.DB) {
		service := NewUserService(db)
		testName := "BulkUserOperations"

		// Create multiple users
		userIDs := make([]string, 10)
		for i := 0; i < 10; i++ {
			req := &usersDto.CreateUserRequest{
				Email:   "bulk-" + testName + "-" + string(rune('1'+i)) + "@example.com",
				Name:    "Bulk User " + testName + " " + string(rune('1'+i)),
				ClerkID: "clerk-bulk-" + testName + "-" + string(rune('1'+i)),
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
func (suite *UserIntegrationTestSuite) TestUserDataConsistency() {
	suite.container.WithTransaction(func(db *gorm.DB) {
		service := NewUserService(db)
		testName := "DataConsistency"

		// Create user
		createReq := &usersDto.CreateUserRequest{
			Email:   "consistency-" + testName + "@example.com",
			Name:    "Consistency User " + testName,
			ClerkID: "clerk-consistency-" + testName,
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

			// Verify ClerkID remains unchanged
			assert.Equal(suite.T(), createReq.ClerkID, fetchedUser.ClerkID)

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
		var foundUser *usersDto.UserResponse
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

// =====================================
// END-TO-END TESTS
// =====================================

// TestCompleteUserLifecycle tests the entire user lifecycle in one flow
// This is valuable as an integration test because it validates:
// 1. Request routing and middleware
// 2. Multi-step business flows
// 3. Data consistency across operations
func (suite *UserIntegrationTestSuite) TestCompleteUserLifecycle() {
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Setup router with transaction DB
		cfg := &config.Config{Environment: "test"}
		testRouter := router.Setup(db, cfg)

		testName := "CompleteUserLifecycle"
		userID := ""

		// Step 1: Create User
		createReq := usersDto.CreateUserRequest{
			Email:   "lifecycle-" + testName + "@example.com",
			Name:    "Lifecycle User " + testName,
			ClerkID: "clerk-lifecycle-" + testName,
		}

		reqBody, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusCreated, w.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		assert.NoError(suite.T(), err)

		data := createResponse["data"].(map[string]interface{})
		userID = data["id"].(string)
		assert.NotEmpty(suite.T(), userID)
		assert.Equal(suite.T(), createReq.Email, data["email"])

		// Step 2: Get User (verify creation)
		req, _ = http.NewRequest("GET", "/api/v1/users/"+userID, nil)
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var getResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &getResponse)
		assert.NoError(suite.T(), err)

		userData := getResponse["data"].(map[string]interface{})
		assert.Equal(suite.T(), userID, userData["id"])
		assert.Equal(suite.T(), createReq.Email, userData["email"])

		// Step 3: Update User
		newEmail := "updated-lifecycle-" + testName + "@example.com"
		newName := "Updated Lifecycle User " + testName
		updateReq := usersDto.UpdateUserRequest{
			Email: &newEmail,
			Name:  &newName,
		}

		reqBody, _ = json.Marshal(updateReq)
		req, _ = http.NewRequest("PUT", "/api/v1/users/"+userID, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var updateResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
		assert.NoError(suite.T(), err)

		updatedData := updateResponse["data"].(map[string]interface{})
		assert.Equal(suite.T(), newEmail, updatedData["email"])
		assert.Equal(suite.T(), newName, updatedData["name"])

		// Step 4: Delete User
		req, _ = http.NewRequest("DELETE", "/api/v1/users/"+userID, nil)
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		// Step 5: Verify deletion (should return 404)
		req, _ = http.NewRequest("GET", "/api/v1/users/"+userID, nil)
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	})
}

// TestUserPaginationEndToEnd tests pagination flow which involves:
// 1. Multiple database operations
// 2. Query parameter parsing
// 3. Response formatting
// This is valuable as integration test because it validates the complete pagination pipeline
func (suite *UserIntegrationTestSuite) TestUserPaginationEndToEnd() {
	suite.container.WithTransaction(func(db *gorm.DB) {
		// Setup router with transaction DB
		cfg := &config.Config{Environment: "test"}
		testRouter := router.Setup(db, cfg)

		testName := "PaginationEndToEnd"

		// Create multiple users via API (tests creation + data consistency)
		userIDs := make([]string, 5)
		for i := 0; i < 5; i++ {
			createReq := usersDto.CreateUserRequest{
				Email:   "pagination-" + testName + "-" + string(rune('1'+i)) + "@example.com",
				Name:    "Pagination User " + testName + " " + string(rune('1'+i)),
				ClerkID: "clerk-pagination-" + testName + "-" + string(rune('1'+i)),
			}

			reqBody, _ := json.Marshal(createReq)
			req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusCreated, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			data := response["data"].(map[string]interface{})
			userIDs[i] = data["id"].(string)
		}

		// Test pagination - page 1
		req, _ := http.NewRequest("GET", "/api/v1/users?page=1&per_page=3", nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var paginationResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &paginationResponse)
		assert.NoError(suite.T(), err)

		data := paginationResponse["data"].(map[string]interface{})
		items := data["items"].([]interface{})

		assert.Len(suite.T(), items, 3)
		assert.Equal(suite.T(), float64(5), data["total"]) // JSON numbers are float64
		assert.Equal(suite.T(), float64(1), data["page"])
		assert.Equal(suite.T(), float64(3), data["per_page"])

		// Test pagination - page 2
		req, _ = http.NewRequest("GET", "/api/v1/users?page=2&per_page=3", nil)
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &paginationResponse)
		assert.NoError(suite.T(), err)

		data = paginationResponse["data"].(map[string]interface{})
		items = data["items"].([]interface{})

		assert.Len(suite.T(), items, 2) // Remaining 2 users
		assert.Equal(suite.T(), float64(5), data["total"])
		assert.Equal(suite.T(), float64(2), data["page"])
		assert.Equal(suite.T(), float64(3), data["per_page"])
	})
}

func TestUserIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserIntegrationTestSuite))
}
