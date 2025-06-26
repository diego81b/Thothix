package services

import (
	"gorm.io/gorm"

	"thothix-backend/internal/dto"
	"thothix-backend/internal/mappers"
	"thothix-backend/internal/middleware"
	"thothix-backend/internal/models"
)

type UserService struct {
	db     *gorm.DB
	mapper *mappers.UserMapper
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db:     db,
		mapper: mappers.NewUserMapper(),
	}
}

// GetUserByID retrieves a user by ID using Response pattern with lazy evaluation
func (s *UserService) GetUserByID(userID string) *dto.GetUserResponse {
	return dto.NewGetUserResponse(func() dto.Validation[*dto.UserResponse] {
		var validationErrors []dto.Error

		// Validation logic
		if userID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "User ID cannot be empty", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*dto.UserResponse](validationErrors...)
		}

		// Business logic - this can panic and will be caught by Try()
		var user models.User
		if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return dto.Invalid[*dto.UserResponse](dto.NewError("USER_NOT_FOUND", "User not found", nil))
			}
			panic(err) // Gets converted to Exception by Try()
		}

		response := s.mapper.ModelToResponse(&user)
		return dto.Success(response)
	})
}

// GetUserByClerkID retrieves a user by Clerk ID using Response pattern
func (s *UserService) GetUserByClerkID(clerkID string) *dto.GetUserResponse {
	return dto.NewGetUserResponse(func() dto.Validation[*dto.UserResponse] {
		var validationErrors []dto.Error

		// Validation logic
		if clerkID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Clerk ID cannot be empty", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*dto.UserResponse](validationErrors...)
		}

		// Business logic
		var user models.User
		if err := s.db.Where("clerk_id = ?", clerkID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return dto.Invalid[*dto.UserResponse](dto.NewError("USER_NOT_FOUND", "User not found", nil))
			}
			panic(err)
		}

		response := s.mapper.ModelToResponse(&user)
		return dto.Success(response)
	})
}

// GetUsers retrieves paginated list of users using Response pattern
func (s *UserService) GetUsers(req *dto.PaginationRequest) *dto.GetUsersResponse {
	return dto.NewGetUsersResponse(func() dto.Validation[*dto.UserListResponse] {
		var validationErrors []dto.Error

		// Validation
		if req == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Pagination request cannot be nil", nil))
		}

		if req != nil && req.Page < 1 {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Page must be greater than 0", nil))
		}

		if req != nil && (req.PerPage < 1 || req.PerPage > 100) {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Per page must be between 1 and 100", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*dto.UserListResponse](validationErrors...)
		}

		// Business logic
		var users []models.User
		var total int64

		// Count total users
		if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
			panic(err)
		}

		// Calculate offset
		offset := (req.Page - 1) * req.PerPage

		// Get users with pagination
		if err := s.db.Offset(offset).Limit(req.PerPage).Find(&users).Error; err != nil {
			panic(err)
		}

		// Convert to response DTOs
		userResponses := make([]dto.UserResponse, len(users))
		for i, user := range users {
			userResponse := s.mapper.ModelToResponse(&user)
			userResponses[i] = *userResponse
		}

		// Calculate pagination metadata
		totalPages := int(total) / req.PerPage
		if int(total)%req.PerPage > 0 {
			totalPages++
		}

		response := dto.NewUserListResponse(userResponses, total, req.Page, req.PerPage)

		return dto.Success(response)
	})
}

// CreateUser creates a new user using Response pattern
func (s *UserService) CreateUser(req *dto.CreateUserRequest) *dto.CreateUserResponse {
	return dto.NewCreateUserResponse(func() dto.Validation[*dto.UserResponse] {
		var validationErrors []dto.Error

		// Validation
		if req == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Create user request cannot be nil", nil))
		}

		if req != nil && req.Email == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Email is required", nil))
		}

		if req != nil && req.Name == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Name is required", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*dto.UserResponse](validationErrors...)
		}

		// Check if user already exists
		var existingUser models.User
		if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			// User already exists - this is a validation error
			return dto.Invalid[*dto.UserResponse](dto.NewError("CONFLICT", "User with email already exists", nil))
		} else if err != gorm.ErrRecordNotFound {
			// Database error
			panic(err)
		}

		// Create new user
		user := s.mapper.CreateRequestToModel(req)
		if err := s.db.Create(user).Error; err != nil {
			panic(err)
		}

		response := s.mapper.ModelToResponse(user)
		return dto.Success(response)
	})
}

// UpdateUser updates an existing user using Response pattern
func (s *UserService) UpdateUser(userID string, req *dto.UpdateUserRequest) *dto.UpdateUserResponse {
	return dto.NewUpdateUserResponse(func() dto.Validation[*dto.UserResponse] {
		var validationErrors []dto.Error

		// Validation
		if userID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "User ID cannot be empty", nil))
		}

		if req == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Update user request cannot be nil", nil))
		}

		// Check if at least one field is provided for update
		if req != nil && req.Email == nil && req.Name == nil && req.Username == nil && req.FirstName == nil && req.LastName == nil && req.AvatarURL == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "no fields to update", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*dto.UserResponse](validationErrors...)
		}

		// Find existing user
		var user models.User
		if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return dto.Invalid[*dto.UserResponse](dto.NewError("USER_NOT_FOUND", "User not found", nil))
			}
			panic(err)
		}

		// Update user fields
		if req.Email != nil {
			user.Email = *req.Email
		}
		if req.Name != nil {
			user.Name = *req.Name
		}
		if req.Username != nil {
			user.Username = *req.Username
		}

		// Save changes
		if err := s.db.Save(&user).Error; err != nil {
			panic(err)
		}

		response := s.mapper.ModelToResponse(&user)
		return dto.Success(response)
	})
}

// DeleteUser deletes a user using Response pattern
func (s *UserService) DeleteUser(userID string) *dto.DeleteUserResponse {
	return dto.NewDeleteUserResponse(func() dto.Validation[string] {
		var validationErrors []dto.Error

		// Validation
		if userID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "User ID cannot be empty", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[string](validationErrors...)
		}

		// Find user to delete
		var user models.User
		if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return dto.Invalid[string](dto.NewError("USER_NOT_FOUND", "User not found", nil))
			}
			panic(err)
		}

		// Delete user
		if err := s.db.Delete(&user).Error; err != nil {
			panic(err)
		}

		return dto.Success("User deleted successfully")
	})
}

// SyncUserFromClerk synchronizes a user from Clerk using Response pattern
func (s *UserService) SyncUserFromClerk(req *dto.ClerkUserSyncRequest) *dto.CreateUserResponse {
	return dto.NewCreateUserResponse(func() dto.Validation[*dto.UserResponse] {
		var validationErrors []dto.Error

		// Validation
		if req == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Clerk user sync request cannot be nil", nil))
		}

		if req != nil && req.ClerkID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Clerk ID is required", nil))
		}

		if req != nil && req.Email == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Email is required", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*dto.UserResponse](validationErrors...)
		}

		// Check if user already exists by Clerk ID
		var existingUser models.User
		if err := s.db.Where("clerk_id = ?", req.ClerkID).First(&existingUser).Error; err == nil {
			// User exists, update it
			existingUser.Email = req.Email
			existingUser.Name = req.Name
			if req.Username != "" {
				existingUser.Username = req.Username
			}

			if err := s.db.Save(&existingUser).Error; err != nil {
				panic(err)
			}

			response := s.mapper.ModelToResponse(&existingUser)
			return dto.Success(response)
		} else if err != gorm.ErrRecordNotFound {
			// Database error
			panic(err)
		}

		// Create new user
		user := &models.User{
			ClerkID:  req.ClerkID,
			Email:    req.Email,
			Name:     req.Name,
			Username: req.Username,
		}

		if err := s.db.Create(user).Error; err != nil {
			panic(err)
		}

		response := s.mapper.ModelToResponse(user)
		return dto.Success(response)
	})
}

// ProcessClerkWebhook processes a Clerk webhook using Response pattern
func (s *UserService) ProcessClerkWebhook(userData *middleware.UserWebhookData) *dto.ClerkSyncUserResponse {
	return dto.NewClerkSyncUserResponse(func() dto.Validation[*dto.ClerkUserSyncResponse] {
		var validationErrors []dto.Error

		// Validation
		if userData == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "User webhook data cannot be nil", nil))
		}

		if userData != nil && userData.ID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "User ID is required", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*dto.ClerkUserSyncResponse](validationErrors...)
		}

		// Convert webhook data to sync request
		email := ""
		if len(userData.EmailAddresses) > 0 {
			email = userData.EmailAddresses[0].EmailAddress
		}

		// Handle optional string fields
		var name string
		if userData.FirstName != nil && userData.LastName != nil {
			name = *userData.FirstName + " " + *userData.LastName
		} else if userData.FirstName != nil {
			name = *userData.FirstName
		} else if userData.LastName != nil {
			name = *userData.LastName
		}

		var username string
		if userData.Username != nil {
			username = *userData.Username
		}

		var avatarURL string
		if userData.ImageURL != nil {
			avatarURL = *userData.ImageURL
		}

		syncReq := &dto.ClerkUserSyncRequest{
			ClerkID:   userData.ID,
			Email:     email,
			Name:      name,
			Username:  username,
			AvatarURL: avatarURL,
		}

		// Sync the user - extract result from the Response pattern
		userResponse := s.SyncUserFromClerk(syncReq)

		// Use Match to extract the result
		var result *dto.ClerkUserSyncResponse
		userResponse.Match(
			func(err error) interface{} {
				panic(err) // Re-panic system errors
			},
			func(user *dto.UserResponse) interface{} {
				result = &dto.ClerkUserSyncResponse{
					User:    *user,
					IsNew:   true, // We'll determine this later
					Message: "User synchronized successfully",
				}
				return result
			},
			func(errors []dto.Error) interface{} {
				panic(errors[0]) // Convert validation error to panic for exceptional handling
			},
		)

		return dto.Success(result)
	})
}
