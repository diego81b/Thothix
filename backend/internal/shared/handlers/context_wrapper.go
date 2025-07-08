package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"thothix-backend/internal/shared/dto"
)

// ContextWrapper wraps gin.Context to provide convenience methods for standardized error responses.
// This ensures consistent error handling and response format across all API endpoints.
type ContextWrapper struct {
	*gin.Context
}

// WrapContext creates a new ContextWrapper from a gin.Context.
// Use this to access enhanced error response methods.
func WrapContext(c *gin.Context) *ContextWrapper {
	return &ContextWrapper{Context: c}
}

// SystemErrorResponse sends a standardized system error response with logging.
// Use this for internal server errors, database errors, and other system-level failures.
func (c *ContextWrapper) SystemErrorResponse(err error, logMessage string, logArgs ...interface{}) {
	response := dto.LoggedSystemErrorResponse(err, logMessage, logArgs...)
	c.JSON(http.StatusInternalServerError, response)
}

// ValidationErrorResponse sends a standardized validation error response with logging.
// Use this for input validation failures, business rule violations, and data validation errors.
func (c *ContextWrapper) ValidationErrorResponse(errors []dto.Error, logMessage string, logArgs ...interface{}) {
	response := dto.LoggedValidationErrorResponse(errors, logMessage, logArgs...)
	c.JSON(http.StatusBadRequest, response)
}

// NotFoundErrorResponse sends a standardized not found error response with logging.
// Use this when a requested resource cannot be found.
func (c *ContextWrapper) NotFoundErrorResponse(resourceType, identifier string) {
	log.Printf("Resource not found: %s with identifier %s", resourceType, identifier)
	response := dto.ErrorViewModel{
		Success: false,
		Error:   "not_found",
		Message: resourceType + " not found",
	}
	c.JSON(http.StatusNotFound, response)
}

// UnauthorizedErrorResponse sends a standardized unauthorized error response.
// Use this when authentication is required but missing or invalid.
func (c *ContextWrapper) UnauthorizedErrorResponse(message string) {
	log.Printf("Unauthorized access attempt: %s", message)
	response := dto.ErrorViewModel{
		Success: false,
		Error:   "unauthorized",
		Message: message,
	}
	c.JSON(http.StatusUnauthorized, response)
}

// ForbiddenErrorResponse sends a standardized forbidden error response.
// Use this when the user is authenticated but lacks permission for the requested operation.
func (c *ContextWrapper) ForbiddenErrorResponse(message string) {
	log.Printf("Forbidden access attempt: %s", message)
	response := dto.ErrorViewModel{
		Success: false,
		Error:   "forbidden",
		Message: message,
	}
	c.JSON(http.StatusForbidden, response)
}

// BadRequestErrorResponse sends a standardized bad request error response.
// Use this for malformed requests, missing required parameters, etc.
func (c *ContextWrapper) BadRequestErrorResponse(message string) {
	log.Printf("Bad request: %s", message)
	response := dto.ErrorViewModel{
		Success: false,
		Error:   "bad_request",
		Message: message,
	}
	c.JSON(http.StatusBadRequest, response)
}

// ConflictErrorResponse sends a standardized conflict error response.
// Use this when a resource already exists (e.g., duplicate email, username).
func (c *ContextWrapper) ConflictErrorResponse(message string) {
	log.Printf("Conflict error: %s", message)
	response := dto.ErrorViewModel{
		Success: false,
		Error:   "conflict",
		Message: message,
	}
	c.JSON(http.StatusConflict, response)
}

// SuccessResponse sends a standardized success response with data.
// Use this for successful operations that return data.
func (c *ContextWrapper) SuccessResponse(data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// CreatedResponse sends a standardized created response with data.
// Use this for successful resource creation operations.
func (c *ContextWrapper) CreatedResponse(data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// NoContentResponse sends a standardized no content response.
// Use this for successful operations that don't return data (like DELETE).
func (c *ContextWrapper) NoContentResponse() {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Operation completed successfully",
	})
}

// DeletedResponse sends a standardized response for successful deletion operations.
// Use this for successful resource deletion operations.
func (c *ContextWrapper) DeletedResponse(message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})
}
