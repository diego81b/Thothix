package dto

import (
	"testing"
)

// TestCreateUserRequest tests the CreateUserRequest structure
func TestCreateUserRequest(t *testing.T) {
	t.Run("CreateUserRequest has all required fields", func(t *testing.T) {
		createUser := CreateUserRequest{
			Email:     "test@example.com",
			Name:      "John Doe",
			Username:  "johndoe",
			ClerkID:   "clerk_123",
			FirstName: "John",
			LastName:  "Doe",
		}

		if createUser.Email != "test@example.com" {
			t.Errorf("Expected Email 'test@example.com', got '%s'", createUser.Email)
		}
		if createUser.Name != "John Doe" {
			t.Errorf("Expected Name 'John Doe', got '%s'", createUser.Name)
		}
		if createUser.Username != "johndoe" {
			t.Errorf("Expected Username 'johndoe', got '%s'", createUser.Username)
		}
		if createUser.ClerkID != "clerk_123" {
			t.Errorf("Expected ClerkID 'clerk_123', got '%s'", createUser.ClerkID)
		}
		if createUser.FirstName != "John" {
			t.Errorf("Expected FirstName 'John', got '%s'", createUser.FirstName)
		}
		if createUser.LastName != "Doe" {
			t.Errorf("Expected LastName 'Doe', got '%s'", createUser.LastName)
		}
	})
}

// TestUpdateUserRequest tests the UpdateUserRequest structure
func TestUpdateUserRequest(t *testing.T) {
	t.Run("UpdateUserRequest has all optional fields", func(t *testing.T) {
		email := "updated@example.com"
		name := "Updated Name"
		username := "updateduser"
		firstName := "UpdatedFirst"
		lastName := "UpdatedLast"
		avatarURL := "https://example.com/avatar.jpg"

		updateUser := UpdateUserRequest{
			Email:     &email,
			Name:      &name,
			Username:  &username,
			FirstName: &firstName,
			LastName:  &lastName,
			AvatarURL: &avatarURL,
		}

		if updateUser.Email == nil || *updateUser.Email != email {
			t.Errorf("Expected Email '%s', got %v", email, updateUser.Email)
		}
		if updateUser.Name == nil || *updateUser.Name != name {
			t.Errorf("Expected Name '%s', got %v", name, updateUser.Name)
		}
		if updateUser.Username == nil || *updateUser.Username != username {
			t.Errorf("Expected Username '%s', got %v", username, updateUser.Username)
		}
		if updateUser.FirstName == nil || *updateUser.FirstName != firstName {
			t.Errorf("Expected FirstName '%s', got %v", firstName, updateUser.FirstName)
		}
		if updateUser.LastName == nil || *updateUser.LastName != lastName {
			t.Errorf("Expected LastName '%s', got %v", lastName, updateUser.LastName)
		}
		if updateUser.AvatarURL == nil || *updateUser.AvatarURL != avatarURL {
			t.Errorf("Expected AvatarURL '%s', got %v", avatarURL, updateUser.AvatarURL)
		}
	})

	t.Run("UpdateUserRequest allows nil fields", func(t *testing.T) {
		updateUser := UpdateUserRequest{}

		if updateUser.Email != nil {
			t.Error("Expected Email to be nil")
		}
		if updateUser.Name != nil {
			t.Error("Expected Name to be nil")
		}
		if updateUser.Username != nil {
			t.Error("Expected Username to be nil")
		}
		if updateUser.FirstName != nil {
			t.Error("Expected FirstName to be nil")
		}
		if updateUser.LastName != nil {
			t.Error("Expected LastName to be nil")
		}
		if updateUser.AvatarURL != nil {
			t.Error("Expected AvatarURL to be nil")
		}
	})
}

// TestUserResponse tests the UserResponse structure
func TestUserResponse(t *testing.T) {
	t.Run("UserResponse has all required fields", func(t *testing.T) {
		user := UserResponse{
			ID:       "user_123",
			Email:    "user@example.com",
			Name:     "Test User",
			ClerkID:  "clerk_456",
			Username: "testuser",
		}

		if user.ID != "user_123" {
			t.Errorf("Expected ID 'user_123', got '%s'", user.ID)
		}
		if user.Email != "user@example.com" {
			t.Errorf("Expected Email 'user@example.com', got '%s'", user.Email)
		}
		if user.Name != "Test User" {
			t.Errorf("Expected Name 'Test User', got '%s'", user.Name)
		}
		if user.ClerkID != "clerk_456" {
			t.Errorf("Expected ClerkID 'clerk_456', got '%s'", user.ClerkID)
		}
		if user.Username != "testuser" {
			t.Errorf("Expected Username 'testuser', got '%s'", user.Username)
		}
	})
}

// TestUserListResponse tests the UserListResponse structure
func TestUserListResponse(t *testing.T) {
	t.Run("UserListResponse contains users and pagination", func(t *testing.T) {
		users := []UserResponse{
			{
				ID:       "user_1",
				Email:    "user1@example.com",
				Name:     "User One",
				ClerkID:  "clerk_1",
				Username: "user1",
			},
			{
				ID:       "user_2",
				Email:    "user2@example.com",
				Name:     "User Two",
				ClerkID:  "clerk_2",
				Username: "user2",
			},
		}

		response := UserListResponse{
			Users: users,
			PaginationMeta: PaginationMeta{
				Total:      2,
				Page:       1,
				PerPage:    10,
				TotalPages: 1,
			},
		}

		if len(response.Users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(response.Users))
		}

		if response.Users[0].Email != "user1@example.com" {
			t.Errorf("Expected first user email 'user1@example.com', got '%s'", response.Users[0].Email)
		}

		if response.Total != 2 {
			t.Errorf("Expected total 2, got %d", response.Total)
		}
	})
}

// TestClerkUserSyncRequest tests the ClerkUserSyncRequest structure
func TestClerkUserSyncRequest(t *testing.T) {
	t.Run("ClerkUserSyncRequest has all required fields", func(t *testing.T) {
		clerkUser := ClerkUserSyncRequest{
			ClerkID:   "clerk_789",
			Email:     "clerk@example.com",
			Name:      "Clerk User",
			Username:  "clerkuser",
			AvatarURL: "https://example.com/avatar.jpg",
		}

		if clerkUser.ClerkID != "clerk_789" {
			t.Errorf("Expected ClerkID 'clerk_789', got '%s'", clerkUser.ClerkID)
		}
		if clerkUser.Email != "clerk@example.com" {
			t.Errorf("Expected Email 'clerk@example.com', got '%s'", clerkUser.Email)
		}
		if clerkUser.Name != "Clerk User" {
			t.Errorf("Expected Name 'Clerk User', got '%s'", clerkUser.Name)
		}
		if clerkUser.Username != "clerkuser" {
			t.Errorf("Expected Username 'clerkuser', got '%s'", clerkUser.Username)
		}
		if clerkUser.AvatarURL != "https://example.com/avatar.jpg" {
			t.Errorf("Expected AvatarURL 'https://example.com/avatar.jpg', got '%s'", clerkUser.AvatarURL)
		}
	})
}

// TestClerkUserSyncResponse tests the ClerkUserSyncResponse structure
func TestClerkUserSyncResponse(t *testing.T) {
	t.Run("ClerkUserSyncResponse has all required fields", func(t *testing.T) {
		user := UserResponse{
			ID:       "user_123",
			Email:    "user@example.com",
			Name:     "Test User",
			ClerkID:  "clerk_456",
			Username: "testuser",
		}

		syncResponse := ClerkUserSyncResponse{
			User:    user,
			IsNew:   true,
			Message: "User synchronized successfully",
		}

		if syncResponse.User.ID != user.ID {
			t.Errorf("Expected User ID '%s', got '%s'", user.ID, syncResponse.User.ID)
		}
		if !syncResponse.IsNew {
			t.Error("Expected IsNew to be true")
		}
		if syncResponse.Message != "User synchronized successfully" {
			t.Errorf("Expected Message 'User synchronized successfully', got '%s'", syncResponse.Message)
		}
	})
}

// TestResponseTypes tests the response wrapper types
func TestResponseTypes(t *testing.T) {
	t.Run("NewGetUserResponse creates valid response", func(t *testing.T) {
		user := &UserResponse{
			ID:    "user_123",
			Email: "test@example.com",
			Name:  "Test User",
		}

		response := NewGetUserResponse(func() Validation[*UserResponse] {
			return Valid(user)
		})

		if response == nil {
			t.Error("Expected response to be created")
		}
		if response.Response == nil {
			t.Error("Expected response.Response to be created")
		}
	})

	t.Run("NewCreateUserResponse creates valid response", func(t *testing.T) {
		user := &UserResponse{
			ID:    "user_456",
			Email: "create@example.com",
			Name:  "Created User",
		}

		response := NewCreateUserResponse(func() Validation[*UserResponse] {
			return Valid(user)
		})

		if response == nil {
			t.Error("Expected response to be created")
		}
		if response.Response == nil {
			t.Error("Expected response.Response to be created")
		}
	})

	t.Run("NewGetUsersResponse creates valid response", func(t *testing.T) {
		userList := &UserListResponse{
			Users: []UserResponse{
				{ID: "user_1", Email: "user1@example.com", Name: "User One"},
			},
		}

		response := NewGetUsersResponse(func() Validation[*UserListResponse] {
			return Valid(userList)
		})

		if response == nil {
			t.Error("Expected response to be created")
		}
		if response.Response == nil {
			t.Error("Expected response.Response to be created")
		}
	})
}
