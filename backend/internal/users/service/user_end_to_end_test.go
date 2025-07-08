package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"thothix-backend/internal/config"
	"thothix-backend/internal/shared/router"
	sharedTesting "thothix-backend/internal/shared/testing"
	"thothix-backend/internal/users/domain"
	usersDto "thothix-backend/internal/users/dto"
)

// UserEndToEndTestSuite tests complete user flows that span multiple components
// These tests validate the entire request pipeline: Router -> Middleware -> Handler -> Service -> DB
type UserEndToEndTestSuite struct {
	suite.Suite
	container *sharedTesting.PostgresTestContainer
}

func (suite *UserEndToEndTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.container = sharedTesting.GetSharedTestContainer(
		suite.T(),
		"integration/user_e2e",
		[]interface{}{&domain.User{}},
	)
}

func (suite *UserEndToEndTestSuite) TearDownSuite() {
	// Cleanup handled by t.Cleanup in GetSharedTestContainer
}

// TestCompleteUserLifecycle tests the entire user lifecycle in one flow
// This is valuable as an integration test because it validates:
// 1. Request routing and middleware
// 2. Multi-step business flows
// 3. Data consistency across operations
func (suite *UserEndToEndTestSuite) TestCompleteUserLifecycle() {
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
func (suite *UserEndToEndTestSuite) TestUserPaginationEndToEnd() {
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

func TestUserEndToEndTestSuite(t *testing.T) {
	suite.Run(t, new(UserEndToEndTestSuite))
}
