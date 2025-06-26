package dto

import "fmt"

// Error represents a validation/business error
type Error struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func NewError(code, message string, details map[string]string) Error {
	return Error{Code: code, Message: message, Details: details}
}

func (e Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Code
}

// Exceptional holds either a value T or an Exception
type Exceptional[T any] struct {
	value T
	err   error
	isErr bool
}

func NewExceptional[T any](value T) Exceptional[T] {
	return Exceptional[T]{value: value, isErr: false}
}

func NewExceptionalError[T any](err error) Exceptional[T] {
	return Exceptional[T]{err: err, isErr: true}
}

// Match for Exceptional - returns interface{}
func (e Exceptional[T]) Match(
	onException func(error) interface{},
	onSuccess func(T) interface{},
) interface{} {
	if e.isErr {
		return onException(e.err)
	}
	return onSuccess(e.value)
}

// Validation holds either a value T or validation errors
type Validation[T any] struct {
	value   T
	errors  []Error
	isValid bool
}

func Valid[T any](value T) Validation[T] {
	return Validation[T]{value: value, isValid: true}
}

func Invalid[T any](errors ...Error) Validation[T] {
	return Validation[T]{errors: errors, isValid: false}
}

// Match for Validation - returns interface{}
func (v Validation[T]) Match(
	onInvalid func([]Error) interface{},
	onValid func(T) interface{},
) interface{} {
	if !v.isValid {
		return onInvalid(v.errors)
	}
	return onValid(v.value)
}

// Response represents the main response wrapper with lazy evaluation
type Response[T any] struct {
	producer func() Validation[T]
	result   *Exceptional[Validation[T]]
}

// NewResponse creates a new Response with lazy producer
func NewResponse[T any](producer func() Validation[T]) *Response[T] {
	return &Response[T]{producer: producer}
}

// Try executes a function and catches panics as errors
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

// Match provides pattern matching - returns interface{} instead of TResult
func (r *Response[T]) Match(
	onException func(error) interface{},
	onSuccess func(T) interface{},
	onFailure func([]Error) interface{},
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
func Success[T any](value T) Validation[T] {
	return Valid(value)
}

func Failure[T any](errors ...Error) Validation[T] {
	return Invalid[T](errors...)
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

// PaginationRequest represents generic pagination request parameters
type PaginationRequest struct {
	Page    int `json:"page" form:"page" validate:"min=1"`
	PerPage int `json:"per_page" form:"per_page" validate:"min=1,max=100"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ManagedErrorResult creates a managed error result from an error
func ManagedErrorResult(err error) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"error":   "internal_server_error",
		"message": err.Error(),
	}
}

// Helper functions for HTTP responses
func ToManagedErrorResult(errors []Error) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"errors":  errors,
	}
}

// ToManagedErrorResult converts a slice of errors to managed error result
func ErrorsToManagedResult(errors []Error) map[string]interface{} {
	return ToManagedErrorResult(errors)
}
