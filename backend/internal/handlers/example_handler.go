package handlers

import (
	"github.com/gin-gonic/gin"
	"thothix-backend/internal/dto"
)

// ExampleHandler shows how to use the ContextWrapper for standardized responses.
// This demonstrates all the available response methods.
type ExampleHandler struct{}

// NewExampleHandler creates a new ExampleHandler instance.
func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

// ExampleSuccess demonstrates a successful response.
// @Summary Example success response
// @Description Shows how to return a success response with data
// @Tags examples
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/examples/success [get]
func (h *ExampleHandler) ExampleSuccess(c *gin.Context) {
	ctx := WrapContext(c)

	data := map[string]interface{}{
		"message":   "This is a successful response",
		"timestamp": "2025-06-26T10:00:00Z",
	}

	ctx.SuccessResponse(data)
}

// ExampleValidationError demonstrates a validation error response.
// @Summary Example validation error response
// @Description Shows how to return validation errors
// @Tags examples
// @Accept json
// @Produce json
// @Failure 400 {object} dto.ErrorViewModel
// @Router /api/v1/examples/validation-error [get]
func (h *ExampleHandler) ExampleValidationError(c *gin.Context) {
	ctx := WrapContext(c)

	errors := []dto.Error{
		dto.NewError("FIELD_REQUIRED", "Name is required", map[string]string{"field": "name"}),
		dto.NewError("FIELD_TOO_SHORT", "Password must be at least 8 characters", map[string]string{"field": "password", "min_length": "8"}),
	}

	ctx.ValidationErrorResponse(errors, "Validation failed for user input")
}

// ExampleSystemError demonstrates a system error response.
// @Summary Example system error response
// @Description Shows how to return system errors
// @Tags examples
// @Accept json
// @Produce json
// @Failure 500 {object} dto.ErrorViewModel
// @Router /api/v1/examples/system-error [get]
func (h *ExampleHandler) ExampleSystemError(c *gin.Context) {
	ctx := WrapContext(c)

	// Simulate a system error
	err := dto.NewError("DATABASE_ERROR", "Connection to database failed", nil)

	ctx.SystemErrorResponse(err, "Database connection error in example handler")
}

// ExampleNotFound demonstrates a not found error response.
// @Summary Example not found response
// @Description Shows how to return not found errors
// @Tags examples
// @Accept json
// @Produce json
// @Failure 404 {object} dto.ErrorViewModel
// @Router /api/v1/examples/not-found [get]
func (h *ExampleHandler) ExampleNotFound(c *gin.Context) {
	ctx := WrapContext(c)

	ctx.NotFoundErrorResponse("User", "12345")
}

// ExampleUnauthorized demonstrates an unauthorized error response.
// @Summary Example unauthorized response
// @Description Shows how to return unauthorized errors
// @Tags examples
// @Accept json
// @Produce json
// @Failure 401 {object} dto.ErrorViewModel
// @Router /api/v1/examples/unauthorized [get]
func (h *ExampleHandler) ExampleUnauthorized(c *gin.Context) {
	ctx := WrapContext(c)

	ctx.UnauthorizedErrorResponse("Authentication token is missing or invalid")
}

// ExampleForbidden demonstrates a forbidden error response.
// @Summary Example forbidden response
// @Description Shows how to return forbidden errors
// @Tags examples
// @Accept json
// @Produce json
// @Failure 403 {object} dto.ErrorViewModel
// @Router /api/v1/examples/forbidden [get]
func (h *ExampleHandler) ExampleForbidden(c *gin.Context) {
	ctx := WrapContext(c)

	ctx.ForbiddenErrorResponse("You don't have permission to access this resource")
}

// ExampleBadRequest demonstrates a bad request error response.
// @Summary Example bad request response
// @Description Shows how to return bad request errors
// @Tags examples
// @Accept json
// @Produce json
// @Failure 400 {object} dto.ErrorViewModel
// @Router /api/v1/examples/bad-request [get]
func (h *ExampleHandler) ExampleBadRequest(c *gin.Context) {
	ctx := WrapContext(c)

	ctx.BadRequestErrorResponse("Invalid JSON format in request body")
}

// ExampleCreated demonstrates a created response.
// @Summary Example created response
// @Description Shows how to return a created response
// @Tags examples
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/examples/created [post]
func (h *ExampleHandler) ExampleCreated(c *gin.Context) {
	ctx := WrapContext(c)

	data := map[string]interface{}{
		"id":      "12345",
		"message": "Resource created successfully",
	}

	ctx.CreatedResponse(data)
}

// ExampleNoContent demonstrates a no content response.
// @Summary Example no content response
// @Description Shows how to return a no content response
// @Tags examples
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/examples/no-content [delete]
func (h *ExampleHandler) ExampleNoContent(c *gin.Context) {
	ctx := WrapContext(c)

	ctx.NoContentResponse()
}
