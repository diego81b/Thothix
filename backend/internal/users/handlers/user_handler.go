package handlers

import (
	"github.com/gin-gonic/gin"

	"thothix-backend/internal/shared/dto"
	"thothix-backend/internal/shared/handlers"
	usersDto "thothix-backend/internal/users/dto"
	"thothix-backend/internal/users/service"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} usersDto.UserListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	wrapper := handlers.WrapContext(c)

	var request dto.PaginationRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		wrapper.BadRequestErrorResponse("Invalid query parameters")
		return
	}

	// Set defaults
	if request.Page == 0 {
		request.Page = 1
	}
	if request.PerPage == 0 {
		request.PerPage = 20
	}

	// Get response from service using the Response pattern
	response := h.userService.GetUsers(&request)

	// Use Match pattern to handle all three cases with wrapper methods
	response.Match(
		// Exception case
		func(err error) interface{} {
			wrapper.SystemErrorResponse(err, "Failed to retrieve users list")
			return nil
		},
		// Success case
		func(result *usersDto.UserListDto) interface{} {
			wrapper.SuccessResponse(result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			wrapper.ValidationErrorResponse(errors, "Users list validation failed")
			return nil
		},
	)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} usersDto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	wrapper := handlers.WrapContext(c)
	userID := c.Param("id")

	// Get response from service
	response := h.userService.GetUserByID(userID)

	// Use Match pattern to handle all three cases with wrapper methods
	response.Match(
		// Exception case
		func(err error) interface{} {
			wrapper.SystemErrorResponse(err, "Failed to retrieve user with ID: %s", userID)
			return nil
		},
		// Success case
		func(result *usersDto.UserDto) interface{} {
			wrapper.SuccessResponse(result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			// Check if it's a not found error
			if len(errors) > 0 && errors[0].Code == "USER_NOT_FOUND" {
				wrapper.NotFoundErrorResponse("User", userID)
			} else {
				wrapper.ValidationErrorResponse(errors, "User retrieval validation failed for ID: %s", userID)
			}
			return nil
		},
	)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body usersDto.CreateUserRequest true "User data"
// @Success 201 {object} usersDto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	wrapper := handlers.WrapContext(c)

	var request usersDto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		wrapper.BadRequestErrorResponse("Invalid request payload")
		return
	}

	// Get response from service
	response := h.userService.CreateUser(&request)

	// Use Match pattern to handle all three cases with wrapper methods
	response.Match(
		// Exception case
		func(err error) interface{} {
			wrapper.SystemErrorResponse(err, "Failed to create user")
			return nil
		},
		// Success case
		func(result *usersDto.UserDto) interface{} {
			wrapper.CreatedResponse(result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			// Check if it's a conflict error
			if len(errors) > 0 && errors[0].Code == "CONFLICT" {
				wrapper.ConflictErrorResponse("User already exists with this email")
			} else {
				wrapper.ValidationErrorResponse(errors, "User creation validation failed")
			}
			return nil
		},
	)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body usersDto.UpdateUserRequest true "User data"
// @Success 200 {object} usersDto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	wrapper := handlers.WrapContext(c)
	userID := c.Param("id")

	var request usersDto.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		wrapper.BadRequestErrorResponse("Invalid request payload")
		return
	}

	// Get response from service
	response := h.userService.UpdateUser(userID, &request)

	// Use Match pattern to handle all three cases with wrapper methods
	response.Match(
		// Exception case
		func(err error) interface{} {
			wrapper.SystemErrorResponse(err, "Failed to update user with ID: %s", userID)
			return nil
		},
		// Success case
		func(result *usersDto.UserDto) interface{} {
			wrapper.SuccessResponse(result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			// Check if it's a not found error
			if len(errors) > 0 && errors[0].Code == "USER_NOT_FOUND" {
				wrapper.NotFoundErrorResponse("User", userID)
			} else {
				wrapper.ValidationErrorResponse(errors, "User update validation failed for ID: %s", userID)
			}
			return nil
		},
	)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	wrapper := handlers.WrapContext(c)
	userID := c.Param("id")

	// Get response from service
	response := h.userService.DeleteUser(userID)

	// Use Match pattern to handle all three cases with wrapper methods
	response.Match(
		// Exception case
		func(err error) interface{} {
			wrapper.SystemErrorResponse(err, "Failed to delete user with ID: %s", userID)
			return nil
		},
		// Success case
		func(result string) interface{} {
			wrapper.DeletedResponse(result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			// Check if it's a not found error
			if len(errors) > 0 && errors[0].Code == "USER_NOT_FOUND" {
				wrapper.NotFoundErrorResponse("User", userID)
			} else {
				wrapper.ValidationErrorResponse(errors, "User deletion validation failed for ID: %s", userID)
			}
			return nil
		},
	)
}
