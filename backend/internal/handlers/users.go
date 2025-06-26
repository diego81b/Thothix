package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"thothix-backend/internal/dto"
	"thothix-backend/internal/services"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
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
// @Success 200 {object} dto.UserListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	var request dto.PaginationRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
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

	// Use Match pattern to handle all three cases
	response.Match(
		// Exception case
		func(err error) interface{} {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		},
		// Success case
		func(result *dto.UserListResponse) interface{} {
			c.JSON(http.StatusOK, result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			c.JSON(http.StatusBadRequest, dto.ToManagedErrorResult(errors))
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
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	// Get response from service
	response := h.userService.GetUserByID(userID)

	// Use Match pattern to handle all three cases
	response.Match(
		// Exception case
		func(err error) interface{} {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		},
		// Success case
		func(result *dto.UserResponse) interface{} {
			c.JSON(http.StatusOK, result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			// Check if it's a not found error
			if len(errors) > 0 && errors[0].Code == "USER_NOT_FOUND" {
				c.JSON(http.StatusNotFound, dto.ToManagedErrorResult(errors))
			} else {
				c.JSON(http.StatusBadRequest, dto.ToManagedErrorResult(errors))
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
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get response from service
	response := h.userService.CreateUser(&request)

	// Use Match pattern to handle all three cases
	response.Match(
		// Exception case
		func(err error) interface{} {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		},
		// Success case
		func(result *dto.UserResponse) interface{} {
			c.JSON(http.StatusCreated, result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			// Check if it's a conflict error
			if len(errors) > 0 && errors[0].Code == "CONFLICT" {
				c.JSON(http.StatusConflict, dto.ToManagedErrorResult(errors))
			} else {
				c.JSON(http.StatusBadRequest, dto.ToManagedErrorResult(errors))
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
// @Param user body dto.UpdateUserRequest true "User data"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var request dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get response from service
	response := h.userService.UpdateUser(userID, &request)

	// Use Match pattern to handle all three cases
	response.Match(
		// Exception case
		func(err error) interface{} {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		},
		// Success case
		func(result *dto.UserResponse) interface{} {
			c.JSON(http.StatusOK, result)
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			// Check if it's a not found error
			if len(errors) > 0 && errors[0].Code == "USER_NOT_FOUND" {
				c.JSON(http.StatusNotFound, dto.ToManagedErrorResult(errors))
			} else {
				c.JSON(http.StatusBadRequest, dto.ToManagedErrorResult(errors))
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
	userID := c.Param("id")

	// Get response from service
	response := h.userService.DeleteUser(userID)

	// Use Match pattern to handle all three cases
	response.Match(
		// Exception case
		func(err error) interface{} {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		},
		// Success case
		func(result string) interface{} {
			c.JSON(http.StatusOK, gin.H{"message": result})
			return nil
		},
		// Validation failure case
		func(errors []dto.Error) interface{} {
			// Check if it's a not found error
			if len(errors) > 0 && errors[0].Code == "USER_NOT_FOUND" {
				c.JSON(http.StatusNotFound, dto.ToManagedErrorResult(errors))
			} else {
				c.JSON(http.StatusBadRequest, dto.ToManagedErrorResult(errors))
			}
			return nil
		},
	)
}
