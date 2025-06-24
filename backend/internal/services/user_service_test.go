package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pgtest "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"thothix-backend/internal/dto"
	"thothix-backend/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()

	// Create PostgreSQL container
	dbName := "testdb"
	dbUser := "testuser"
	dbPassword := "testpass"

	postgresContainer, err := pgtest.Run(ctx,
		"docker.io/postgres:16-alpine",
		pgtest.WithDatabase(dbName),
		pgtest.WithUsername(dbUser),
		pgtest.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	// Clean up container after test
	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %s", err)
		}
	})

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Connect to database with silent logging for tests
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

func TestUserService_SyncUserFromClerk(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	tests := []struct {
		name          string
		existingUser  *models.User
		request       *dto.ClerkUserSyncRequest
		expectedIsNew bool
		expectError   bool
	}{
		{
			name:         "create new user",
			existingUser: nil,
			request: &dto.ClerkUserSyncRequest{
				ClerkID:   "clerk_123",
				Email:     "test@example.com",
				Name:      "Test User",
				Username:  "testuser",
				AvatarURL: "https://example.com/avatar.jpg",
			},
			expectedIsNew: true,
			expectError:   false,
		},
		{
			name: "update existing user",
			existingUser: &models.User{
				ClerkID:  "clerk_123",
				Email:    "old@example.com",
				Name:     "Old Name",
				Username: "olduser",
				LastSync: time.Now().Add(-time.Hour),
			},
			request: &dto.ClerkUserSyncRequest{
				ClerkID:   "clerk_123",
				Email:     "new@example.com",
				Name:      "New Name",
				Username:  "newuser",
				AvatarURL: "https://example.com/avatar.jpg",
			},
			expectedIsNew: false,
			expectError:   false,
		},
		{
			name: "no changes needed",
			existingUser: &models.User{
				ClerkID:   "clerk_123",
				Email:     "same@example.com",
				Name:      "Same Name",
				Username:  "sameuser",
				AvatarURL: "https://example.com/avatar.jpg",
				LastSync:  time.Now().Add(-time.Hour),
			},
			request: &dto.ClerkUserSyncRequest{
				ClerkID:   "clerk_123",
				Email:     "same@example.com",
				Name:      "Same Name",
				Username:  "sameuser",
				AvatarURL: "https://example.com/avatar.jpg",
			},
			expectedIsNew: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.existingUser != nil {
				err := db.Create(tt.existingUser).Error
				require.NoError(t, err)
			}

			// Execute
			response, err := service.SyncUserFromClerk(tt.request)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, response)
				assert.Equal(t, tt.expectedIsNew, response.IsNew)
				assert.Equal(t, tt.request.Email, response.User.Email)
				assert.Equal(t, tt.request.Name, response.User.Name)
				assert.Equal(t, tt.request.Username, response.User.Username)
				assert.Equal(t, tt.request.AvatarURL, response.User.AvatarURL)
			}

			// Cleanup
			db.Exec("DELETE FROM users")
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	// Setup test data
	user := &models.User{
		ClerkID:  "clerk_123",
		Email:    "test@example.com",
		Name:     "Test User",
		Username: "testuser",
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		expectError bool
	}{
		{
			name:        "existing user",
			userID:      user.ID, // Use the generated ID
			expectError: false,
		},
		{
			name:        "non-existing user",
			userID:      "nonexistent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.GetUserByID(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, response)
				assert.Equal(t, tt.userID, response.ID)
				assert.Equal(t, user.Email, response.Email)
			}
		})
	}
}

func TestUserService_GetUserByClerkID(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	// Setup test data - create 2 users with different ClerkIDs
	user1 := &models.User{
		ClerkID:  "clerk_123",
		Email:    "user1@example.com",
		Name:     "Test User 1",
		Username: "testuser1",
	}
	err := db.Create(user1).Error
	require.NoError(t, err)

	user2 := &models.User{
		ClerkID:  "clerk_456",
		Email:    "user2@example.com",
		Name:     "Test User 2",
		Username: "testuser2",
	}
	err = db.Create(user2).Error
	require.NoError(t, err)

	tests := []struct {
		name         string
		clerkID      string
		expectError  bool
		expectedUser *models.User
	}{
		{
			name:         "existing user 1",
			clerkID:      "clerk_123",
			expectError:  false,
			expectedUser: user1,
		},
		{
			name:         "existing user 2",
			clerkID:      "clerk_456",
			expectError:  false,
			expectedUser: user2,
		},
		{
			name:         "non-existing user",
			clerkID:      "clerk_nonexistent",
			expectError:  true,
			expectedUser: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.GetUserByClerkID(tt.clerkID)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
				assert.Nil(t, response, "Expected nil response but got: %v", response)
			} else {
				assert.NoError(t, err, "Unexpected error: %v", err)
				require.NotNil(t, response, "Expected response but got nil")
				assert.Equal(t, tt.clerkID, response.ClerkID)
				assert.Equal(t, tt.expectedUser.Email, response.Email)
				assert.Equal(t, tt.expectedUser.Name, response.Name)
				assert.Equal(t, tt.expectedUser.Username, response.Username)
				// Verify we got the right user by checking the generated ID matches
				assert.Equal(t, tt.expectedUser.ID, response.ID)
			}
		})
	}
}

func TestUserService_GetUsers(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	// Setup test data
	users := []models.User{
		{
			ClerkID: "clerk_1",
			Email:   "user1@example.com",
			Name:    "User One",
		},
		{
			ClerkID: "clerk_2",
			Email:   "user2@example.com",
			Name:    "User Two",
		},
		{
			ClerkID: "clerk_3",
			Email:   "user3@example.com",
			Name:    "User Three",
		},
	}

	for _, user := range users {
		err := db.Create(&user).Error
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		request       *dto.GetUsersRequest
		expectedCount int
		expectedTotal int64
		expectedPages int
	}{
		{
			name: "first page",
			request: &dto.GetUsersRequest{
				Page:    1,
				PerPage: 2,
			},
			expectedCount: 2,
			expectedTotal: 3,
			expectedPages: 2,
		},
		{
			name: "second page",
			request: &dto.GetUsersRequest{
				Page:    2,
				PerPage: 2,
			},
			expectedCount: 1,
			expectedTotal: 3,
			expectedPages: 2,
		},
		{
			name: "all users",
			request: &dto.GetUsersRequest{
				Page:    1,
				PerPage: 10,
			},
			expectedCount: 3,
			expectedTotal: 3,
			expectedPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.GetUsers(tt.request)

			assert.NoError(t, err)
			require.NotNil(t, response)
			assert.Len(t, response.Users, tt.expectedCount)
			assert.Equal(t, tt.expectedTotal, response.Total)
			assert.Equal(t, tt.request.Page, response.Page)
			assert.Equal(t, tt.request.PerPage, response.PerPage)
			assert.Equal(t, tt.expectedPages, response.TotalPages)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	service := NewUserService(db)

	// Setup test data
	user := &models.User{
		ClerkID:  "clerk_123",
		Email:    "old@example.com",
		Name:     "Old Name",
		Username: "olduser",
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		request     *dto.UpdateUserRequest
		expectError bool
		errorMsg    string
	}{
		{
			name:   "valid update",
			userID: user.ID, // Use the generated ID
			request: &dto.UpdateUserRequest{
				Email:     stringPtr("new@example.com"),
				FirstName: stringPtr("New"),
				LastName:  stringPtr("Name"),
				Username:  stringPtr("newuser"),
			},
			expectError: false,
		},
		{
			name:   "partial update",
			userID: user.ID, // Use the generated ID
			request: &dto.UpdateUserRequest{
				Email: stringPtr("partial@example.com"),
			},
			expectError: false,
		},
		{
			name:        "empty update",
			userID:      user.ID, // Use the generated ID
			request:     &dto.UpdateUserRequest{},
			expectError: true,
			errorMsg:    "no fields to update",
		},
		{
			name:   "non-existing user",
			userID: "nonexistent",
			request: &dto.UpdateUserRequest{
				Email: stringPtr("test@example.com"),
			},
			expectError: true,
			errorMsg:    "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.UpdateUser(tt.userID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, response)
				assert.Equal(t, tt.userID, response.ID)

				if tt.request.Email != nil {
					assert.Equal(t, *tt.request.Email, response.Email)
				}
				if tt.request.Username != nil {
					assert.Equal(t, *tt.request.Username, response.Username)
				}
			}
		})
	}
}

// Helper function for tests
func stringPtr(s string) *string {
	return &s
}
