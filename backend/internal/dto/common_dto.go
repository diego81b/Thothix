package dto

import (
	"fmt"
	"log"
)

// Error represents a validation/business error with structured information.
// Used to provide detailed error context in validation failures and business rule violations.
type Error struct {
	Code    string            `json:"code"`              // Error code (e.g., "VALIDATION_ERROR", "USER_NOT_FOUND")
	Message string            `json:"message"`           // Human-readable error message
	Details map[string]string `json:"details,omitempty"` // Additional context (field names, values, etc.)
}

// NewError creates a new Error instance with the specified code, message, and optional details.
// This is the preferred way to create structured errors throughout the application.
func NewError(code, message string, details map[string]string) Error {
	return Error{Code: code, Message: message, Details: details}
}

// Error implements the error interface, providing a string representation of the structured error.
// Returns a formatted string with code, message, and optional details for better debugging.
func (e Error) Error() string {
	if e.Message != "" {
		if len(e.Details) > 0 {
			return fmt.Sprintf("%s: %s (details: %v)", e.Code, e.Message, e.Details)
		}
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Code
}

// Exceptional holds either a value T or an Exception.
// This represents the result of an operation that might throw an exception.
type Exceptional[T any] struct {
	value T
	err   error
}

// NewExceptional creates a new Exceptional containing a successful value.
// Use this when an operation completes successfully without exceptions.
func NewExceptional[T any](value T) Exceptional[T] {
	return Exceptional[T]{value: value}
}

// NewExceptionalError creates a new Exceptional containing an exception/error.
// Use this when an operation fails with an exception that should be handled gracefully.
func NewExceptionalError[T any](err error) Exceptional[T] {
	return Exceptional[T]{err: err}
}

// Match for Exceptional - implements functional pattern matching.
// Executes onException if this contains an error, or onSuccess if it contains a value.
// Returns interface{} to allow flexible return types in pattern matching scenarios.
func (e Exceptional[T]) Match(
	onException func(error) interface{},
	onSuccess func(T) interface{},
) interface{} {
	if e.err != nil {
		return onException(e.err)
	}
	return onSuccess(e.value)
}

// IsSuccess returns true if this Exceptional contains a value (no error).
func (e Exceptional[T]) IsSuccess() bool {
	return e.err == nil
}

// GetError returns the error if present, nil otherwise.
func (e Exceptional[T]) GetError() error {
	return e.err
}

// GetValue returns the value. Should only be called if IsSuccess() returns true.
func (e Exceptional[T]) GetValue() T {
	return e.value
}

// Validation holds either a value T or validation errors.
// This represents the result of validation logic - either success with data or failure with errors.
type Validation[T any] struct {
	value  T
	errors []Error
}

// Valid creates a new Validation containing a successful validation result.
// Use this when validation passes and you have a valid value to return.
func Valid[T any](value T) Validation[T] {
	return Validation[T]{value: value}
}

// Invalid creates a new Validation containing validation errors.
// Use this when validation fails and you need to return structured error information.
func Invalid[T any](errors ...Error) Validation[T] {
	return Validation[T]{errors: errors}
}

// Match for Validation - implements functional pattern matching for validation results.
// Executes onInvalid if validation failed (with error list), or onValid if validation succeeded.
// Returns interface{} to support flexible return types in different matching contexts.
func (v Validation[T]) Match(
	onInvalid func([]Error) interface{},
	onValid func(T) interface{},
) interface{} {
	if len(v.errors) > 0 {
		return onInvalid(v.errors)
	}
	return onValid(v.value)
}

// IsValid returns true if validation succeeded (no errors).
func (v Validation[T]) IsValid() bool {
	return len(v.errors) == 0
}

// GetValue returns the validated value. Should only be called if IsValid() returns true.
func (v Validation[T]) GetValue() T {
	return v.value
}

// GetErrors returns the validation errors. Empty slice if validation succeeded.
func (v Validation[T]) GetErrors() []Error {
	return v.errors
}

// Response represents the main response wrapper with lazy evaluation.
// The producer function is executed only when Match() is called, enabling deferred computation.
type Response[T any] struct {
	producer func() Validation[T]        // Lazy producer function
	result   *Exceptional[Validation[T]] // Cached result after first execution
}

// NewResponse creates a new Response with lazy producer function.
// The producer will be executed only when Match() is called for the first time.
func NewResponse[T any](producer func() Validation[T]) *Response[T] {
	return &Response[T]{producer: producer}
}

// Try executes a function and catches panics as errors.
// Converts any panic into an Exceptional error, enabling safe execution of potentially dangerous code.
// Returns Exceptional[T] containing either the successful result or the caught panic as an error.
func Try[T any](fn func() T) Exceptional[T] {
	var result T
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("panic: %v", r)
				}
			}
		}()

		result = fn()
	}()

	if err != nil {
		return NewExceptionalError[T](err)
	}

	return NewExceptional(result)
}

// Match provides pattern matching for Response with lazy evaluation and exception handling.
// Executes the producer function only once (lazy evaluation) and handles three scenarios:
// - onException: Called if producer panics or throws exception
// - onSuccess: Called if producer succeeds and validation passes
// - onFailure: Called if producer succeeds but validation fails
// Returns interface{} to support flexible return types in pattern matching.
func (r *Response[T]) Match(
	onException func(error) any,
	onSuccess func(T) any,
	onFailure func([]Error) any,
) interface{} {
	// Lazy evaluation - execute producer only when needed
	if r.result == nil {
		result := Try(r.producer)
		r.result = &result
	}

	return r.result.Match(
		onException,
		func(validation Validation[T]) interface{} {
			return validation.Match(onFailure, onSuccess)
		},
	)
}

// Factory methods

// Success creates a successful Validation[T] containing the provided value.
// This is an alias for Valid() that provides semantic clarity in factory method contexts.
func Success[T any](value T) Validation[T] {
	return Valid(value)
}

// Failure creates a failed Validation[T] containing the provided validation errors.
// This is an alias for Invalid() that provides semantic clarity in factory method contexts.
func Failure[T any](errors ...Error) Validation[T] {
	return Invalid[T](errors...)
}

// PaginationMeta represents pagination metadata for list responses.
// Contains information about total count, current page, and page size calculations.
type PaginationMeta struct {
	Total      int64 `json:"total"`       // Total number of items across all pages
	Page       int   `json:"page"`        // Current page number (1-based)
	PerPage    int   `json:"per_page"`    // Number of items per page
	TotalPages int   `json:"total_pages"` // Total number of pages available
	// HasNext e HasPrevious rimossi per semplificazione
}

// PaginatedListResponse represents a generic paginated list response for any type T.
// This provides a standardized structure for all paginated API responses.
type PaginatedListResponse[T any] struct {
	Items []T `json:"items"`
	PaginationMeta
}

// NewPaginatedListResponse creates a new paginated list response with the provided items and pagination metadata.
// Automatically calculates pagination metadata based on the provided parameters.
func NewPaginatedListResponse[T any](items []T, total int64, page int, perPage int) *PaginatedListResponse[T] {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	if totalPages == 0 {
		totalPages = 1 // Ensure at least 1 page even with 0 items
	}

	return &PaginatedListResponse[T]{
		Items: items,
		PaginationMeta: PaginationMeta{
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		},
	}
}

// ErrorResponse represents a standard error response for HTTP APIs.
// Provides a consistent structure for returning errors to clients.
type ErrorResponse struct {
	Error   string `json:"error"`   // Error code or type
	Message string `json:"message"` // Human-readable error description
}

// ErrorViewModel represents a standardized error response structure for all API endpoints.
// Ensures consistent error format across the entire application.
type ErrorViewModel struct {
	Success bool    `json:"success"`          // Always false for error responses
	Error   string  `json:"error"`            // Error code or type
	Message string  `json:"message"`          // Human-readable error description
	Errors  []Error `json:"errors,omitempty"` // Detailed validation errors (optional)
}

// Helper functions for HTTP responses

// SystemErrorResponse creates a standardized HTTP error response from a system error.
// Returns a standardized ErrorViewModel for internal server error responses.
func SystemErrorResponse(err error) ErrorViewModel {
	return ErrorViewModel{
		Success: false,
		Error:   "internal_server_error",
		Message: err.Error(),
	}
}

// ValidationErrorResponse converts validation errors to a standardized HTTP response.
// Returns a standardized ErrorViewModel for validation error responses.
func ValidationErrorResponse(errors []Error) ErrorViewModel {
	return ErrorViewModel{
		Success: false,
		Error:   "validation_error",
		Message: "One or more validation errors occurred",
		Errors:  errors,
	}
}

// LoggedSystemErrorResponse creates a system error response with automatic logging.
// Logs the error with the provided context message and returns a standardized ErrorViewModel.
func LoggedSystemErrorResponse(err error, logMessage string, logArgs ...interface{}) ErrorViewModel {
	log.Printf(logMessage+": %v", append(logArgs, err)...)
	return SystemErrorResponse(err)
}

// LoggedValidationErrorResponse creates a validation error response with automatic logging.
// Logs the validation errors with the provided context message and returns a standardized ErrorViewModel.
func LoggedValidationErrorResponse(errors []Error, logMessage string, logArgs ...interface{}) ErrorViewModel {
	log.Printf(logMessage+": %v", append(logArgs, errors)...)
	return ValidationErrorResponse(errors)
}

// Generic list response type for consistency across all list endpoints
type ListResponse[T any] struct {
	*Response[*PaginatedListResponse[T]]
}

// NewListResponse creates a new generic list response with lazy evaluation.
// Provides a standardized way to create paginated list responses for any type T.
func NewListResponse[T any](producer func() Validation[*PaginatedListResponse[T]]) *ListResponse[T] {
	return &ListResponse[T]{
		Response: NewResponse(producer),
	}
}

// PaginationRequest represents generic pagination request parameters.
// Used in API endpoints that support paginated results.
type PaginationRequest struct {
	Page    int `json:"page" form:"page" validate:"min=1"`                 // Page number (minimum 1)
	PerPage int `json:"per_page" form:"per_page" validate:"min=1,max=100"` // Items per page (1-100)
}
