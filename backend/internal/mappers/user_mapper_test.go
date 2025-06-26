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

	// Test data
	now := time.Now()
	user := &models.User{
		BaseModel: models.BaseModel{
			ID:        "1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		ClerkID:  "clerk_123",
		Email:    "test@example.com",
		Name:     "John Doe",
		Username: "johndoe",
	}

	// Execute
	response := mapper.ModelToResponse(user)

	// Verify
	require.NotNil(t, response)
	assert.Equal(t, "1", response.ID)
	assert.Equal(t, "clerk_123", response.ClerkID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "John Doe", response.Name)
	assert.Equal(t, "johndoe", response.Username)
}

func TestUserMapper_ModelToResponse_NilInput(t *testing.T) {
	mapper := NewUserMapper()

	// Execute
	response := mapper.ModelToResponse(nil)

	// Verify
	assert.Nil(t, response)
}

func TestUserMapper_CreateRequestToModel(t *testing.T) {
	mapper := NewUserMapper()

	// Test data
	request := &dto.CreateUserRequest{
		ClerkID:   "clerk_456", // Note: This field is not used by CreateRequestToModel
		Email:     "jane@example.com",
		Name:      "Jane Smith", // Note: This field is not used by CreateRequestToModel
		Username:  "janesmith",
		FirstName: "Jane",
		LastName:  "Smith",
	}

	// Execute
	user := mapper.CreateRequestToModel(request)

	// Verify
	require.NotNil(t, user)
	assert.Equal(t, "", user.ClerkID) // ClerkID is not set by CreateRequestToModel
	assert.Equal(t, "jane@example.com", user.Email)
	assert.Equal(t, "Jane Smith", user.Name) // Name is constructed from FirstName + " " + LastName
	assert.Equal(t, "janesmith", user.Username)
	assert.Equal(t, models.RoleUser, user.SystemRole) // Default system role
	assert.Equal(t, "", user.ID)                      // Should be empty for new model
}

func TestUserMapper_CreateRequestToModel_NilInput(t *testing.T) {
	mapper := NewUserMapper()

	// Execute
	user := mapper.CreateRequestToModel(nil)

	// Verify
	assert.Nil(t, user)
}

func TestUserMapper_CreateRequestToModel_EmptyNames(t *testing.T) {
	mapper := NewUserMapper()

	// Test data with empty first/last names
	request := &dto.CreateUserRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "",
		LastName:  "",
	}

	// Execute
	user := mapper.CreateRequestToModel(request)

	// Verify
	require.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, " ", user.Name) // FirstName + " " + LastName = "" + " " + "" = " "
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, models.RoleUser, user.SystemRole)
}

func TestUserMapper_CreateRequestToModel_OnlyFirstName(t *testing.T) {
	mapper := NewUserMapper()

	// Test data with only first name
	request := &dto.CreateUserRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "John",
		LastName:  "",
	}

	// Execute
	user := mapper.CreateRequestToModel(request)

	// Verify
	require.NotNil(t, user)
	assert.Equal(t, "John ", user.Name) // "John" + " " + "" = "John "
}

func TestUserMapper_UpdateRequestToMap(t *testing.T) {
	mapper := NewUserMapper()

	// Test data
	request := &dto.UpdateUserRequest{
		Email:     stringPtr("updated@example.com"),
		FirstName: stringPtr("Updated"),
		LastName:  stringPtr("Name"),
		Username:  stringPtr("updateduser"),
	}

	// Execute
	updateMap := mapper.UpdateRequestToMap(request)

	// Verify
	require.NotNil(t, updateMap)
	assert.Equal(t, "updated@example.com", updateMap["email"])
	assert.Equal(t, "Updated Name", updateMap["name"])
	assert.Equal(t, "updateduser", updateMap["username"])
}

func TestUserMapper_UpdateRequestToMap_PartialUpdate(t *testing.T) {
	mapper := NewUserMapper()

	// Test data - only some fields set
	request := &dto.UpdateUserRequest{
		Email:     stringPtr("partial@example.com"),
		FirstName: stringPtr("Partial"),
		// LastName and Username are nil
	}

	// Execute
	updateMap := mapper.UpdateRequestToMap(request)

	// Verify
	require.NotNil(t, updateMap)
	assert.Equal(t, "partial@example.com", updateMap["email"])
	assert.Equal(t, "Partial", updateMap["name"])
	_, hasUsername := updateMap["username"]
	assert.False(t, hasUsername)
}

func TestUserMapper_UpdateRequestToMap_NilInput(t *testing.T) {
	mapper := NewUserMapper()

	// Execute
	updateMap := mapper.UpdateRequestToMap(nil)

	// Verify
	assert.Nil(t, updateMap)
}

func TestUserMapper_UpdateRequestToMap_AllNilFields(t *testing.T) {
	mapper := NewUserMapper()

	// Test data - all fields are nil
	request := &dto.UpdateUserRequest{
		Email:     nil,
		FirstName: nil,
		LastName:  nil,
		Username:  nil,
	}

	// Execute
	updateMap := mapper.UpdateRequestToMap(request)

	// Verify
	require.NotNil(t, updateMap)
	assert.Empty(t, updateMap)
}

func TestUserMapper_ClerkSyncRequestToModel(t *testing.T) {
	mapper := NewUserMapper()

	// Test data
	request := &dto.ClerkUserSyncRequest{
		ClerkID:   "clerk_sync_123",
		Email:     "sync@example.com",
		Name:      "Sync User",
		Username:  "syncuser",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	// Execute
	user := mapper.ClerkSyncRequestToModel(request)

	// Verify
	require.NotNil(t, user)
	assert.Equal(t, "clerk_sync_123", user.ClerkID)
	assert.Equal(t, "sync@example.com", user.Email)
	assert.Equal(t, "Sync User", user.Name)
	assert.Equal(t, "syncuser", user.Username)
	assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarURL)
	assert.Equal(t, models.RoleUser, user.SystemRole)
}

func TestUserMapper_ClerkSyncRequestToModel_NilInput(t *testing.T) {
	mapper := NewUserMapper()

	// Execute
	user := mapper.ClerkSyncRequestToModel(nil)

	// Verify
	assert.Nil(t, user)
}

func TestUserMapper_ModelsToResponses(t *testing.T) {
	mapper := NewUserMapper()

	// Test data
	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: "1"},
			ClerkID:   "clerk_1",
			Email:     "user1@example.com",
			Name:      "User One",
			Username:  "user1",
		},
		{
			BaseModel: models.BaseModel{ID: "2"},
			ClerkID:   "clerk_2",
			Email:     "user2@example.com",
			Name:      "User Two",
			Username:  "user2",
		},
	}

	// Execute
	responses := mapper.ModelsToResponses(users)

	// Verify
	require.NotNil(t, responses)
	assert.Len(t, responses, 2)

	assert.Equal(t, "1", responses[0].ID)
	assert.Equal(t, "clerk_1", responses[0].ClerkID)
	assert.Equal(t, "user1@example.com", responses[0].Email)

	assert.Equal(t, "2", responses[1].ID)
	assert.Equal(t, "clerk_2", responses[1].ClerkID)
	assert.Equal(t, "user2@example.com", responses[1].Email)
}

func TestUserMapper_ModelsToResponses_NilInput(t *testing.T) {
	mapper := NewUserMapper()

	// Execute
	responses := mapper.ModelsToResponses(nil)

	// Verify
	assert.Nil(t, responses)
}

func TestUserMapper_ModelsToListResponse(t *testing.T) {
	mapper := NewUserMapper()

	// Test data
	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: "1"},
			ClerkID:   "clerk_1",
			Email:     "user1@example.com",
			Name:      "User One",
		},
	}

	// Execute
	listResponse := mapper.ModelsToListResponse(users, 25, 2, 10)

	// Verify
	require.NotNil(t, listResponse)
	assert.Len(t, listResponse.Items, 1)
	assert.Equal(t, int64(25), listResponse.PaginationMeta.Total)
	assert.Equal(t, 2, listResponse.PaginationMeta.Page)
	assert.Equal(t, 10, listResponse.PaginationMeta.PerPage)
	assert.Equal(t, 3, listResponse.PaginationMeta.TotalPages) // 25/10 = 3 pages
}

func TestUserMapper_CreateSyncResponse(t *testing.T) {
	mapper := NewUserMapper()

	// Test data
	user := &models.User{
		BaseModel: models.BaseModel{ID: "1"},
		ClerkID:   "clerk_123",
		Email:     "sync@example.com",
		Name:      "Sync User",
	}

	// Execute
	syncResponse := mapper.CreateSyncResponse(user, true, "User created successfully")

	// Verify
	require.NotNil(t, syncResponse)
	assert.Equal(t, "1", syncResponse.User.ID)
	assert.Equal(t, "clerk_123", syncResponse.User.ClerkID)
	assert.True(t, syncResponse.IsNew)
	assert.Equal(t, "User created successfully", syncResponse.Message)
}

func TestUserMapper_CreateSyncResponse_NilInput(t *testing.T) {
	mapper := NewUserMapper()

	// Execute
	syncResponse := mapper.CreateSyncResponse(nil, false, "test")

	// Verify
	assert.Nil(t, syncResponse)
}

// Helper function to create string pointers for optional fields
func stringPtr(s string) *string {
	return &s
}
