package service

import (
	"gorm.io/gorm"

	"thothix-backend/internal/shared/dto"
	sharedMiddleware "thothix-backend/internal/shared/middleware"
	"thothix-backend/internal/users/domain"
	usersDto "thothix-backend/internal/users/dto"
	"thothix-backend/internal/users/mappers"
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
func (s *UserService) GetUserByID(userID string) *usersDto.GetUserResponse {
	return usersDto.NewGetUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		var validationErrors []dto.Error

		// Validation logic
		if userID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "User ID cannot be empty", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*usersDto.UserResponse](validationErrors...)
		}

		// Business logic - this can panic and will be caught by Try()
		var user domain.User
		if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return dto.Invalid[*usersDto.UserResponse](dto.NewError("USER_NOT_FOUND", "User not found", nil))
			}
			panic(err) // Gets converted to Exception by Try()
		}

		response := s.mapper.ModelToResponse(&user)
		return dto.Success(response)
	})
}

// GetUserByClerkID retrieves a user by Clerk ID using Response pattern
func (s *UserService) GetUserByClerkID(clerkID string) *usersDto.GetUserResponse {
	return usersDto.NewGetUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		var validationErrors []dto.Error

		// Validation logic
		if clerkID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Clerk ID cannot be empty", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*usersDto.UserResponse](validationErrors...)
		}

		// Business logic
		var user domain.User
		if err := s.db.Where("clerk_id = ?", clerkID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return dto.Invalid[*usersDto.UserResponse](dto.NewError("USER_NOT_FOUND", "User not found", nil))
			}
			panic(err)
		}

		response := s.mapper.ModelToResponse(&user)
		return dto.Success(response)
	})
}

// GetUsers retrieves paginated list of users using Response pattern
func (s *UserService) GetUsers(req *dto.PaginationRequest) *usersDto.GetUsersResponse {
	return usersDto.NewGetUsersResponse(func() dto.Validation[*usersDto.UserListResponse] {
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
			return dto.Failure[*usersDto.UserListResponse](validationErrors...)
		}

		// Business logic
		var users []domain.User
		var total int64

		// Count total users
		if err := s.db.Model(&domain.User{}).Count(&total).Error; err != nil {
			panic(err)
		}

		// Calculate offset
		offset := (req.Page - 1) * req.PerPage

		// Get users with pagination
		if err := s.db.Offset(offset).Limit(req.PerPage).Find(&users).Error; err != nil {
			panic(err)
		}

		// Convert to response DTOs
		userResponses := make([]usersDto.UserResponse, len(users))
		for i, user := range users {
			userResponse := s.mapper.ModelToResponse(&user)
			userResponses[i] = *userResponse
		}

		response := usersDto.NewUserListResponse(userResponses, total, req.Page, req.PerPage)
		return dto.Success(response)
	})
}

// CreateUser creates a new user using Response pattern
func (s *UserService) CreateUser(req *usersDto.CreateUserRequest) *usersDto.CreateUserResponse {
	return usersDto.NewCreateUserResponse(func() dto.Validation[*usersDto.UserResponse] {
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
			return dto.Failure[*usersDto.UserResponse](validationErrors...)
		}

		// Check if user already exists
		var existingUser domain.User
		if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return dto.Invalid[*usersDto.UserResponse](dto.NewError("CONFLICT", "User with this email already exists", nil))
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
func (s *UserService) UpdateUser(userID string, req *usersDto.UpdateUserRequest) *usersDto.UpdateUserResponse {
	return usersDto.NewUpdateUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		var validationErrors []dto.Error

		// Validation
		if userID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "User ID cannot be empty", nil))
		}

		if req == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Update user request cannot be nil", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*usersDto.UserResponse](validationErrors...)
		}

		// Get existing user
		var user domain.User
		if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return dto.Invalid[*usersDto.UserResponse](dto.NewError("USER_NOT_FOUND", "User not found", nil))
			}
			panic(err)
		}

		// Apply updates
		s.mapper.UpdateRequestToModel(&user, req)

		// Save changes
		if err := s.db.Save(&user).Error; err != nil {
			panic(err)
		}

		response := s.mapper.ModelToResponse(&user)
		return dto.Success(response)
	})
}

// DeleteUser deletes a user using Response pattern
func (s *UserService) DeleteUser(userID string) *usersDto.DeleteUserResponse {
	return usersDto.NewDeleteUserResponse(func() dto.Validation[string] {
		var validationErrors []dto.Error

		// Validation
		if userID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "User ID cannot be empty", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[string](validationErrors...)
		}

		// Check if user exists
		var user domain.User
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
func (s *UserService) SyncUserFromClerk(req *usersDto.ClerkUserSyncRequest) *usersDto.CreateUserResponse {
	return usersDto.NewCreateUserResponse(func() dto.Validation[*usersDto.UserResponse] {
		var validationErrors []dto.Error

		// Validation
		if req == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Clerk sync request cannot be nil", nil))
		}

		if req != nil && req.ClerkID == "" {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Clerk ID is required", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*usersDto.UserResponse](validationErrors...)
		}

		// Check if user already exists
		var existingUser domain.User
		if err := s.db.Where("clerk_id = ?", req.ClerkID).First(&existingUser).Error; err == nil {
			// Update existing user
			existingUser.SyncFromClerk(domain.ClerkUserData{
				Email:     req.Email,
				Name:      req.Name,
				AvatarURL: req.AvatarURL,
			})

			if err := s.db.Save(&existingUser).Error; err != nil {
				panic(err)
			}

			response := s.mapper.ModelToResponse(&existingUser)
			return dto.Success(response)
		}

		// Create new user
		user := s.mapper.ClerkSyncRequestToModel(req)
		if err := s.db.Create(user).Error; err != nil {
			panic(err)
		}

		response := s.mapper.ModelToResponse(user)
		return dto.Success(response)
	})
}

// ProcessClerkWebhook processes a Clerk webhook using Response pattern
func (s *UserService) ProcessClerkWebhook(userData *sharedMiddleware.UserWebhookData) *usersDto.ClerkSyncUserResponse {
	return usersDto.NewClerkSyncUserResponse(func() dto.Validation[*usersDto.ClerkUserSyncResponse] {
		var validationErrors []dto.Error

		// Validation
		if userData == nil {
			validationErrors = append(validationErrors, dto.NewError("VALIDATION_ERROR", "Webhook data cannot be nil", nil))
		}

		if len(validationErrors) > 0 {
			return dto.Failure[*usersDto.ClerkUserSyncResponse](validationErrors...)
		}

		// Process webhook logic here
		// Extract email from the first email address
		var email string
		if len(userData.EmailAddresses) > 0 {
			email = userData.EmailAddresses[0].EmailAddress
		}

		// Combine first and last name
		var name string
		if userData.FirstName != nil && userData.LastName != nil {
			name = *userData.FirstName + " " + *userData.LastName
		} else if userData.FirstName != nil {
			name = *userData.FirstName
		} else if userData.LastName != nil {
			name = *userData.LastName
		}

		response := &usersDto.ClerkUserSyncResponse{
			User: usersDto.UserResponse{
				ID:      userData.ID,
				ClerkID: userData.ID, // Clerk ID is the same as ID in webhook data
				Email:   email,
				Name:    name,
			},
			IsNew:   true,
			Message: "User synced from Clerk webhook",
		}

		return dto.Success(response)
	})
}
