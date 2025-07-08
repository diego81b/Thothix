package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"thothix-backend/internal/config"
	"thothix-backend/internal/shared/router"
	sharedTesting "thothix-backend/internal/shared/testing"
	"thothix-backend/internal/users/domain"
	usersDto "thothix-backend/internal/users/dto"
)

// setupTestContainer creates a shared test container for e2e tests
func setupTestContainer(t *testing.T) *sharedTesting.PostgresTestContainer {
	gin.SetMode(gin.TestMode)
	return sharedTesting.NewPostgresTestContainer(
		t,
		[]interface{}{&domain.User{}},
	)
}

// TestCompleteUserLifecycle tests the entire user lifecycle in one flow
// This is valuable as an integration test because it validates:
// 1. Request routing and middleware
// 2. Multi-step business flows
// 3. Data consistency across operations
func TestCompleteUserLifecycle(t *testing.T) {
	container := setupTestContainer(t)

	container.WithTransaction(func(db *gorm.DB) {
		// Setup router with transaction DB - using test router without authentication
		cfg := &config.Config{Environment: "test"}
		testRouter := router.SetupTestRouter(db, cfg)

		testName := "CompleteUserLifecycle"
		userID := ""

		// Step 1: Create User
		createReq := usersDto.CreateUserRequest{
			Email: "lifecycle-" + testName + "@example.com",
			Name:  "Lifecycle User " + testName,
			// ClerkID removed - CreateUserRequest is for manual user creation only
		}

		reqBody, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Debug: Print response body if test fails
		if w.Code != http.StatusCreated {
			t.Logf("Create user failed. Status: %d, Body: %s", w.Code, w.Body.String())
		}

		var createResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		assert.NoError(t, err)

		// Debug: Print response structure if data is nil
		if createResponse["data"] == nil {
			t.Logf("Response data is nil. Full response: %+v", createResponse)
		}

		data := createResponse["data"].(map[string]interface{})
		userID = data["id"].(string)
		assert.NotEmpty(t, userID)
		assert.Equal(t, createReq.Email, data["email"])

		// Step 2: Get User (verify creation)
		req, _ = http.NewRequest("GET", "/api/v1/users/"+userID, nil)
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var getResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &getResponse)
		assert.NoError(t, err)

		userData := getResponse["data"].(map[string]interface{})
		assert.Equal(t, userID, userData["id"])
		assert.Equal(t, createReq.Email, userData["email"])

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

		assert.Equal(t, http.StatusOK, w.Code)

		var updateResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
		assert.NoError(t, err)

		updatedData := updateResponse["data"].(map[string]interface{})
		assert.Equal(t, newEmail, updatedData["email"])
		assert.Equal(t, newName, updatedData["name"])

		// Step 4: Delete User
		req, _ = http.NewRequest("DELETE", "/api/v1/users/"+userID, nil)
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Step 5: Verify deletion (should return 404)
		req, _ = http.NewRequest("GET", "/api/v1/users/"+userID, nil)
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestUserPaginationEndToEnd tests pagination flow which involves:
// 1. Multiple database operations
// 2. Query parameter parsing
// 3. Response formatting
// This is valuable as integration test because it validates the complete pagination pipeline
func TestUserPaginationEndToEnd(t *testing.T) {
	container := setupTestContainer(t)

	container.WithTransaction(func(db *gorm.DB) {
		// Setup router with transaction DB - using test router without authentication
		cfg := &config.Config{Environment: "test"}
		testRouter := router.SetupTestRouter(db, cfg)

		testName := "PaginationEndToEnd"

		// Create multiple users via API (tests creation + data consistency)
		userIDs := make([]string, 5)
		for i := range 5 {
			// Fixed: use proper string conversion instead of rune casting
			createReq := usersDto.CreateUserRequest{
				Email: fmt.Sprintf("pagination-%s-%d@example.com", testName, i+1),
				Name:  fmt.Sprintf("Pagination User %s %d", testName, i+1),
			}

			reqBody, _ := json.Marshal(createReq)
			req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			data := response["data"].(map[string]interface{})
			userIDs[i] = data["id"].(string)
		}

		// Test pagination - page 1
		req, _ := http.NewRequest("GET", "/api/v1/users?page=1&per_page=3", nil)
		w := httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var paginationResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &paginationResponse)
		assert.NoError(t, err)

		data := paginationResponse["data"].(map[string]interface{})
		items := data["items"].([]interface{})

		assert.Len(t, items, 3)
		assert.Equal(t, float64(5), data["total"]) // JSON numbers are float64
		assert.Equal(t, float64(1), data["page"])
		assert.Equal(t, float64(3), data["per_page"])

		// Test pagination - page 2
		req, _ = http.NewRequest("GET", "/api/v1/users?page=2&per_page=3", nil)
		w = httptest.NewRecorder()
		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &paginationResponse)
		assert.NoError(t, err)

		data = paginationResponse["data"].(map[string]interface{})
		items = data["items"].([]interface{})

		assert.Len(t, items, 2) // Remaining 2 users
		assert.Equal(t, float64(5), data["total"])
		assert.Equal(t, float64(2), data["page"])
		assert.Equal(t, float64(3), data["per_page"])
	})
}
