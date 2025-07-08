package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"thothix-backend/internal/shared/dto"
)

// TestAssertions provides assertion helpers for Response pattern matching
type TestAssertions struct {
	t *testing.T
}

// NewTestAssertions creates a new instance of TestAssertions
func NewTestAssertions(t *testing.T) *TestAssertions {
	return &TestAssertions{t: t}
}

// AssertOnSuccess asserts that the response is successful and executes the action
func AssertOnSuccess[T any](t *testing.T, response *dto.Response[T], action func(T)) {
	response.Match(
		func(exception error) interface{} {
			assert.Fail(t, "Expected success but got exception", exception.Error())
			return nil
		},
		func(success T) interface{} {
			action(success)
			return nil
		},
		func(errors []dto.Error) interface{} {
			assert.Fail(t, "Expected success but got validation errors", errors)
			return nil
		},
	)
}

// AssertOnValidationError asserts that the response has validation errors and executes the action
func AssertOnValidationError[T any](t *testing.T, response *dto.Response[T], action func([]dto.Error)) {
	response.Match(
		func(exception error) interface{} {
			assert.Fail(t, "Expected validation error but got exception", exception.Error())
			return nil
		},
		func(success T) interface{} {
			assert.Fail(t, "Expected validation error but got success", success)
			return nil
		},
		func(errors []dto.Error) interface{} {
			action(errors)
			return nil
		},
	)
}

// AssertOnException asserts that the response has an exception and executes the action
func AssertOnException[T any](t *testing.T, response *dto.Response[T], action func(error)) {
	response.Match(
		func(exception error) interface{} {
			action(exception)
			return nil
		},
		func(success T) interface{} {
			assert.Fail(t, "Expected exception but got success", success)
			return nil
		},
		func(errors []dto.Error) interface{} {
			assert.Fail(t, "Expected exception but got validation errors", errors)
			return nil
		},
	)
}

// AssertOnSuccessPaginated for paginated responses
func AssertOnSuccessPaginated[T any](t *testing.T, response *dto.Response[*dto.PaginatedListResponse[T]], action func(*dto.PaginatedListResponse[T])) {
	response.Match(
		func(exception error) interface{} {
			assert.Fail(t, "Expected success but got exception", exception.Error())
			return nil
		},
		func(success *dto.PaginatedListResponse[T]) interface{} {
			action(success)
			return nil
		},
		func(errors []dto.Error) interface{} {
			assert.Fail(t, "Expected success but got validation errors", errors)
			return nil
		},
	)
}

// AssertOnValidationErrorPaginated for paginated responses
func AssertOnValidationErrorPaginated[T any](t *testing.T, response *dto.Response[*dto.PaginatedListResponse[T]], action func([]dto.Error)) {
	response.Match(
		func(exception error) interface{} {
			assert.Fail(t, "Expected validation error but got exception", exception.Error())
			return nil
		},
		func(success *dto.PaginatedListResponse[T]) interface{} {
			assert.Fail(t, "Expected validation error but got success", success)
			return nil
		},
		func(errors []dto.Error) interface{} {
			action(errors)
			return nil
		},
	)
}

// Convenience functions that follow common patterns

// AssertSuccessWithValue asserts success and returns the value
func AssertSuccessWithValue[T any](t *testing.T, response *dto.Response[T]) T {
	var result T
	AssertOnSuccess(t, response, func(value T) {
		result = value
	})
	return result
}

// AssertValidationErrorWithCode asserts validation error with specific code
func AssertValidationErrorWithCode[T any](t *testing.T, response *dto.Response[T], expectedCode string) {
	AssertOnValidationError(t, response, func(errors []dto.Error) {
		assert.NotEmpty(t, errors, "Should have validation errors")
		found := false
		for _, err := range errors {
			if err.Code == expectedCode {
				found = true
				break
			}
		}
		assert.True(t, found, "Should contain expected error code: %s", expectedCode)
	})
}

// AssertExceptionWithMessage asserts exception and checks the error message
func AssertExceptionWithMessage[T any](t *testing.T, response *dto.Response[T], expectedMessage string) {
	AssertOnException(t, response, func(err error) {
		assert.Contains(t, err.Error(), expectedMessage, "Should contain expected error message")
	})
}

// AssertValidationErrorWithMessages asserts validation errors and checks for specific messages
func AssertValidationErrorWithMessages[T any](t *testing.T, response *dto.Response[T], expectedMessages []string) {
	AssertOnValidationError(t, response, func(errors []dto.Error) {
		assert.NotEmpty(t, errors, "Should have validation errors")
		for _, expectedMsg := range expectedMessages {
			found := false
			for _, err := range errors {
				if err.Message == expectedMsg {
					found = true
					break
				}
			}
			assert.True(t, found, "Should contain expected error message: %s", expectedMsg)
		}
	})
}

// AssertValidationErrorCount asserts the number of validation errors
func AssertValidationErrorCount[T any](t *testing.T, response *dto.Response[T], expectedCount int) {
	AssertOnValidationError(t, response, func(errors []dto.Error) {
		assert.Len(t, errors, expectedCount, "Should have expected number of validation errors")
	})
}

// AssertSuccessPaginatedWithValue asserts success and returns the paginated value
func AssertSuccessPaginatedWithValue[T any](t *testing.T, response *dto.Response[*dto.PaginatedListResponse[T]]) *dto.PaginatedListResponse[T] {
	var result *dto.PaginatedListResponse[T]
	AssertOnSuccessPaginated(t, response, func(data *dto.PaginatedListResponse[T]) {
		result = data
	})
	return result
}

// AssertPaginatedCount asserts the count of items in a paginated response
func AssertPaginatedCount[T any](t *testing.T, response *dto.Response[*dto.PaginatedListResponse[T]], expectedCount int) {
	AssertOnSuccessPaginated(t, response, func(data *dto.PaginatedListResponse[T]) {
		assert.Len(t, data.Items, expectedCount, "Should have expected number of items")
	})
}

// AssertSuccessAndTest asserts success and runs additional tests on the value
func AssertSuccessAndTest[T any](t *testing.T, response *dto.Response[T], testFunc func(*testing.T, T)) {
	AssertOnSuccess(t, response, func(value T) {
		testFunc(t, value)
	})
}
