# Shared Components

This directory contains shared components that can be used across all domain slices in the vertical slice architecture.

## Structure

- `constants/` - Common constants, error codes, and default values
- `testing/` - Generic test assertion helpers using Go generics

## Components

### Testing Assertions (`testing/assertions.go`)

Generic assertion helpers for Response pattern matching that work with any domain type:

```go
// Example usage in any domain test
func TestSomeDomainService(t *testing.T) {
    response := domainService.SomeOperation(request)

    // Assert success and get the typed result
    result := sharedTesting.AssertSuccessWithValue(t, response.Response)
    assert.Equal(t, "expected", result.SomeField)

    // Assert validation error with specific code
    sharedTesting.AssertValidationErrorWithCode(t, response.Response, "VALIDATION_ERROR")

    // Assert paginated results
    paginatedResult := sharedTesting.AssertSuccessPaginatedWithValue(t, response.Response)
    assert.Len(t, paginatedResult.Items, 2)
}
```

### Constants (`constants/errors.go`)

Common error codes and messages used across domains:

```go
// Example usage in any domain service
import "thothix-backend/internal/shared/constants"

func (s *SomeService) Validate(req *SomeRequest) []dto.Error {
    if req.Field == "" {
        return []dto.Error{{
            Code:    constants.ValidationError,
            Message: "Field is required",
        }}
    }
    return nil
}
```

## Usage Guidelines

1. **Testing**: Use the generic assertion helpers instead of domain-specific ones
2. **Constants**: Use shared constants for common error codes and messages
3. **Extensibility**: Add new shared components here when they're needed by multiple domains
4. **Imports**: Import shared components with descriptive aliases when needed

## Benefits

- **Consistency**: Common error codes and patterns across all domains
- **Reusability**: Generic components that work with any domain type
- **Maintainability**: Single source of truth for shared functionality
- **Type Safety**: Go generics ensure type safety across different domain types
