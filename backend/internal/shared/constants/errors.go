package constants

// Common error codes used across different domains
const (
	// Generic errors
	ValidationError = "VALIDATION_ERROR"
	InternalError   = "INTERNAL_ERROR"
	NotFoundError   = "NOT_FOUND"
	ConflictError   = "CONFLICT"

	// User specific errors
	UserNotFoundError = "USER_NOT_FOUND"
	UserExistsError   = "USER_EXISTS"

	// Authorization errors
	UnauthorizedError = "UNAUTHORIZED"
	ForbiddenError    = "FORBIDDEN"

	// Common messages
	CreatedSuccessfully = "Created successfully"
	UpdatedSuccessfully = "Updated successfully"
	DeletedSuccessfully = "Deleted successfully"
)

// Common pagination defaults
const (
	DefaultPage    = 1
	DefaultPerPage = 10
	MaxPerPage     = 100
)
