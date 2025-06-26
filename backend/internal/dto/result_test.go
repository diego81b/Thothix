package dto

import (
	"errors"
	"testing"
)

// TestError tests the Error type
func TestError(t *testing.T) {
	t.Run("NewError creates error with all fields", func(t *testing.T) {
		details := map[string]string{"field": "value"}
		err := NewError("VALIDATION_ERROR", "Field is required", details)

		if err.Code != "VALIDATION_ERROR" {
			t.Errorf("Expected Code 'VALIDATION_ERROR', got '%s'", err.Code)
		}
		if err.Message != "Field is required" {
			t.Errorf("Expected Message 'Field is required', got '%s'", err.Message)
		}
		if err.Details["field"] != "value" {
			t.Errorf("Expected Details['field'] 'value', got '%s'", err.Details["field"])
		}
	})

	t.Run("Error implements error interface", func(t *testing.T) {
		err := NewError("TEST_ERROR", "Test message", nil)
		expectedString := "TEST_ERROR: Test message"

		if err.Error() != expectedString {
			t.Errorf("Expected '%s', got '%s'", expectedString, err.Error())
		}
	})

	t.Run("Error with no message returns just code", func(t *testing.T) {
		err := NewError("SIMPLE_ERROR", "", nil)

		if err.Error() != "SIMPLE_ERROR" {
			t.Errorf("Expected 'SIMPLE_ERROR', got '%s'", err.Error())
		}
	})
}

// TestExceptional tests the Exceptional type
func TestExceptional(t *testing.T) {
	t.Run("NewExceptional creates success exceptional", func(t *testing.T) {
		value := "test value"
		exceptional := NewExceptional(value)

		result := exceptional.Match(
			func(error) interface{} { t.Error("Should not be called"); return nil },
			func(v string) interface{} { return v },
		)

		if result != value {
			t.Errorf("Expected '%s', got '%v'", value, result)
		}
	})

	t.Run("NewExceptionalError creates error exceptional", func(t *testing.T) {
		testError := errors.New("test error")
		exceptional := NewExceptionalError[string](testError)

		result := exceptional.Match(
			func(err error) interface{} { return err.Error() },
			func(string) interface{} { t.Error("Should not be called"); return nil },
		)

		if result != "test error" {
			t.Errorf("Expected 'test error', got '%v'", result)
		}
	})
}

// TestValidation tests the Validation type
func TestValidation(t *testing.T) {
	t.Run("Valid creates valid validation", func(t *testing.T) {
		value := "valid value"
		validation := Valid(value)

		result := validation.Match(
			func([]Error) interface{} { t.Error("Should not be called"); return nil },
			func(v string) interface{} { return v },
		)

		if result != value {
			t.Errorf("Expected '%s', got '%v'", value, result)
		}
	})

	t.Run("Invalid creates invalid validation", func(t *testing.T) {
		err1 := NewError("ERROR1", "First error", nil)
		err2 := NewError("ERROR2", "Second error", nil)
		validation := Invalid[string](err1, err2)

		result := validation.Match(
			func(errors []Error) interface{} { return len(errors) },
			func(string) interface{} { t.Error("Should not be called"); return nil },
		)

		if result != 2 {
			t.Errorf("Expected 2 errors, got %v", result)
		}
	})
}

// TestResponse tests the Response type with lazy evaluation
func TestResponse(t *testing.T) {
	t.Run("Response with successful producer", func(t *testing.T) {
		value := "response value"
		response := NewResponse(func() Validation[string] {
			return Valid(value)
		})

		result := response.Match(
			func(error) interface{} { t.Error("Should not be called"); return nil },
			func(v string) interface{} { return v },
			func([]Error) interface{} { t.Error("Should not be called"); return nil },
		)

		if result != value {
			t.Errorf("Expected '%s', got '%v'", value, result)
		}
	})

	t.Run("Response with failing producer", func(t *testing.T) {
		err := NewError("VALIDATION_FAILED", "Validation failed", nil)
		response := NewResponse(func() Validation[string] {
			return Invalid[string](err)
		})

		result := response.Match(
			func(error) interface{} { t.Error("Should not be called"); return nil },
			func(string) interface{} { t.Error("Should not be called"); return nil },
			func(errors []Error) interface{} { return errors[0].Message },
		)

		if result != "Validation failed" {
			t.Errorf("Expected 'Validation failed', got '%v'", result)
		}
	})

	t.Run("Response lazy evaluation - producer called only once", func(t *testing.T) {
		callCount := 0
		response := NewResponse(func() Validation[string] {
			callCount++
			return Valid("test")
		})

		// Call Match multiple times
		response.Match(
			func(error) interface{} { return nil },
			func(string) interface{} { return "ok" },
			func([]Error) interface{} { return nil },
		)

		response.Match(
			func(error) interface{} { return nil },
			func(string) interface{} { return "ok" },
			func([]Error) interface{} { return nil },
		)

		if callCount != 1 {
			t.Errorf("Expected producer to be called once, but was called %d times", callCount)
		}
	})
}

// TestTry tests the Try function
func TestTry(t *testing.T) {
	t.Run("Try with successful function", func(t *testing.T) {
		result := Try(func() string {
			return "success"
		})

		value := result.Match(
			func(error) interface{} { t.Error("Should not be called"); return nil },
			func(v string) interface{} { return v },
		)

		if value != "success" {
			t.Errorf("Expected 'success', got '%v'", value)
		}
	})
}

// TestFactoryMethods tests the factory methods
func TestFactoryMethods(t *testing.T) {
	t.Run("Success creates valid validation", func(t *testing.T) {
		value := "success value"
		validation := Success(value)

		result := validation.Match(
			func([]Error) interface{} { t.Error("Should not be called"); return nil },
			func(v string) interface{} { return v },
		)

		if result != value {
			t.Errorf("Expected '%s', got '%v'", value, result)
		}
	})

	t.Run("Failure creates invalid validation", func(t *testing.T) {
		err := NewError("FAILURE", "Failure message", nil)
		validation := Failure[string](err)

		result := validation.Match(
			func(errors []Error) interface{} { return errors[0].Message },
			func(string) interface{} { t.Error("Should not be called"); return nil },
		)

		if result != "Failure message" {
			t.Errorf("Expected 'Failure message', got '%v'", result)
		}
	})
}

// TestPaginationMeta tests the PaginationMeta type
func TestPaginationMeta(t *testing.T) {
	t.Run("PaginationMeta has all required fields", func(t *testing.T) {
		meta := PaginationMeta{
			Total:      100,
			Page:       2,
			PerPage:    10,
			TotalPages: 10,
		}

		if meta.Total != 100 {
			t.Errorf("Expected Total 100, got %d", meta.Total)
		}
		if meta.Page != 2 {
			t.Errorf("Expected Page 2, got %d", meta.Page)
		}
		if meta.PerPage != 10 {
			t.Errorf("Expected PerPage 10, got %d", meta.PerPage)
		}
		if meta.TotalPages != 10 {
			t.Errorf("Expected TotalPages 10, got %d", meta.TotalPages)
		}
	})
}

// TestPaginationRequest tests the PaginationRequest type
func TestPaginationRequest(t *testing.T) {
	t.Run("PaginationRequest has all required fields", func(t *testing.T) {
		req := PaginationRequest{
			Page:    1,
			PerPage: 20,
		}

		if req.Page != 1 {
			t.Errorf("Expected Page 1, got %d", req.Page)
		}
		if req.PerPage != 20 {
			t.Errorf("Expected PerPage 20, got %d", req.PerPage)
		}
	})
}

// TestHelperFunctions tests the helper functions
func TestHelperFunctions(t *testing.T) {
	t.Run("ValidationErrorResponse creates correct structure", func(t *testing.T) {
		errors := []Error{
			NewError("ERROR1", "First error", nil),
			NewError("ERROR2", "Second error", nil),
		}

		result := ValidationErrorResponse(errors)

		if result.Success != false {
			t.Errorf("Expected success to be false, got %v", result.Success)
		}

		if result.Error != "validation_error" {
			t.Errorf("Expected error to be 'validation_error', got %s", result.Error)
		}

		if len(result.Errors) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(result.Errors))
		}
	})

	t.Run("SystemErrorResponse creates correct structure", func(t *testing.T) {
		testErr := NewError("TEST_ERROR", "Test error message", nil)

		result := SystemErrorResponse(testErr)

		if result.Success != false {
			t.Errorf("Expected success to be false, got %v", result.Success)
		}

		if result.Error != "internal_server_error" {
			t.Errorf("Expected error to be 'internal_server_error', got %s", result.Error)
		}

		if result.Message != testErr.Error() {
			t.Errorf("Expected message to be '%s', got %s", testErr.Error(), result.Message)
		}
	})
}

// TestPaginatedListResponse tests the generic PaginatedListResponse type
func TestPaginatedListResponse(t *testing.T) {
	t.Run("PaginatedListResponse with items", func(t *testing.T) {
		items := []string{"item1", "item2", "item3"}
		response := NewPaginatedListResponse(items, 10, 1, 3)

		if len(response.Items) != 3 {
			t.Errorf("Expected 3 items, got %d", len(response.Items))
		}

		if response.Total != 10 {
			t.Errorf("Expected Total 10, got %d", response.Total)
		}

		if response.Page != 1 {
			t.Errorf("Expected Page 1, got %d", response.Page)
		}

		if response.PerPage != 3 {
			t.Errorf("Expected PerPage 3, got %d", response.PerPage)
		}

		if response.TotalPages != 4 {
			t.Errorf("Expected TotalPages 4, got %d", response.TotalPages)
		}
	})

	t.Run("PaginatedListResponse last page", func(t *testing.T) {
		items := []string{"item1"}
		response := NewPaginatedListResponse(items, 10, 4, 3)

		if response.TotalPages != 4 {
			t.Errorf("Expected TotalPages 4, got %d", response.TotalPages)
		}
	})

	t.Run("PaginatedListResponse empty items", func(t *testing.T) {
		items := []string{}
		response := NewPaginatedListResponse(items, 0, 1, 10)

		if response.TotalPages != 1 {
			t.Errorf("Expected TotalPages 1 for empty list, got %d", response.TotalPages)
		}
	})
}
