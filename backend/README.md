# 🚀 Thothix Backend - Go API

[![Go](https://img.shields.io/badge/Go-1.23-blue?style=flat&logo=go)](https://golang.org) [![Gin](https://img.shields.io/badge/Gin-Web_Framework-green?style=flat)](https://github.com/gin-gonic/gin) [![GORM](https://img.shields.io/badge/GORM-ORM-yellow?style=flat)](https://gorm.io)

Modern Go REST API with **Vertical Slice Architecture**, Gin framework, GORM ORM, PostgreSQL, Clerk authentication, and HashiCorp Vault integration.

## 🏗️ Architecture

**Tech Stack**: Gin • GORM v2 • PostgreSQL 15 • Clerk Auth • HashiCorp Vault • Swagger/OpenAPI • Docker

**Vertical Slice Architecture**: Domain-driven organization with business features as independent slices, each containing their complete stack from handlers to data access.

### 🎯 Vertical Slice Architecture Benefits

- **High Cohesion**: Related functionality grouped together by business domain
- **Low Coupling**: Minimal dependencies between different domains
- **Parallel Development**: Teams can work on different slices independently
- **Feature-Focused**: Easy to add, modify, or remove complete business features
- **Clear Boundaries**: Each slice owns its complete vertical stack

### 🧩 Domain Structure

Each business domain (e.g., `users`, `projects`, `chats`) follows this pattern:

```
internal/users/                    # Business Domain Slice
├── domain/                        # Business Entities & Core Logic
│   ├── user.go                   # Domain models with business rules
│   └── user_test.go              # Domain logic tests
├── dto/                          # Domain-Specific DTOs
│   ├── user_dto.go               # Request/response contracts
│   └── user_dto_test.go          # DTO validation tests
├── service/                      # Business Logic Layer
│   ├── user_service.go           # Business operations
│   ├── user_service_interface.go # Service contract
│   └── user_service_test.go      # Business logic tests
├── handlers/                     # HTTP Controllers
│   ├── user_handler.go           # HTTP endpoints
│   ├── user_handler_test.go      # Controller tests
│   └── routes.go                 # Route registration
└── mappers/                      # Data Transformation
    ├── user_mapper.go            # Entity ↔ DTO mapping
    └── user_mapper_test.go       # Mapping tests
```

### 🌐 Shared Components

Cross-domain components that promote consistency and reusability:

```
internal/shared/
├── constants/                    # Common Constants
│   └── errors.go                 # Error codes & messages
├── dto/                          # Common DTOs
│   └── common_dto.go             # Shared response patterns
├── handlers/                     # Generic Handlers
│   ├── context_wrapper.go        # Context utilities
│   └── health.go                 # Health check endpoint
├── middleware/                   # Reusable Middleware
│   ├── clerk_auth.go             # Clerk authentication
│   ├── cors.go                   # CORS configuration
│   ├── logger.go                 # Logging utilities
│   └── context.go                # Context management
├── router/                       # Main Router
│   └── router.go                 # Route setup and registration
├── testing/                      # Generic Test Helpers
│   └── assertions.go             # Type-safe test assertions
└── README.md                     # Shared components documentation
```

**Key Features**:
- **Generic Test Assertions**: Type-safe helpers using Go generics for any domain
- **Common DTOs**: Shared response patterns and validation utilities
- **Reusable Middleware**: Authentication, CORS, logging, and context management
- **Centralized Routing**: Single point for route registration and middleware setup
- **Generic Handlers**: Health checks and utility endpoints
- **Common Constants**: Shared error codes and messages across all domains
- **Extensible**: Easy to add new shared components as needed

Example usage:
```go
// Generic test assertions for any domain
result := sharedTesting.AssertSuccessWithValue(t, response.Response)
sharedTesting.AssertValidationErrorWithCode(t, response.Response, constants.ValidationError)

// Common error codes
if req.Field == "" {
    return []dto.Error{{
        Code:    constants.ValidationError,
        Message: "Field is required",
    }}
}
```

### 🎯 Result Pattern Implementation

Thothix implements a **Result Pattern** with functional programming for type-safe error handling:

#### Core Components

```go
// Generic Response with lazy evaluation and pattern matching
Response[T] struct {
    producer func() Validation[T]
    result   *Exceptional[Validation[T]]
}

// Pattern matching replaces traditional if-err checks
result.Match(onException, onSuccess, onFailure)
```

#### Key Benefits

- **Type safety**: No null references, compile-time guarantees
- **Functional composition**: Pattern matching over if-err chains
- **Lazy evaluation**: Execution only when needed
- **Clear error types**: System exceptions vs validation errors

#### Usage Example

```go
// Service returns typed Response
func (s *UserService) GetUserByID(userID string) *dto.GetUserResponse {
    return dto.NewGetUserResponse(func() dto.Validation[*dto.UserResponse] {
        if userID == "" {
            return dto.Failure[*dto.UserResponse](
                dto.NewError("VALIDATION_ERROR", "User ID required", nil))
        }

        var user models.User
        if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                return dto.Invalid[*dto.UserResponse](
                    dto.NewError("USER_NOT_FOUND", "User not found", nil))
            }
            panic(err) // Auto-converted to Exception
        }

        return dto.Success(s.mapper.ModelToResponse(&user))
    })
}

// Handler with pattern matching
func (h *UserHandler) GetUser(c *gin.Context) {
    result := h.service.GetUserByID(userID)
    result.Match(
        func(err error) interface{} {
            c.JSON(500, dto.ManagedErrorResult(err)); return nil
        },
        func(user *dto.UserResponse) interface{} {
            c.JSON(200, user); return nil
        },
        func(errors []dto.Error) interface{} {
            c.JSON(400, dto.ErrorsToManagedResult(errors)); return nil
        },
    )
}
```

**Architecture**: `shared/dto/common_dto.go` (generic patterns) + `users/dto/user_dto.go` (domain DTOs)

---

## 🔧 Shared Components

The `internal/shared/` directory contains reusable components that can be used across all domain slices:

### Structure

```text
internal/shared/
├── dto/                      # Common DTOs and Result Pattern
│   ├── common_dto.go         # Generic Response, Validation, Error types
│   └── result_test.go        # Tests for Result Pattern
├── handlers/                 # Shared HTTP utilities
│   └── context_wrapper.go    # Gin context wrapper with standardized responses
├── router/                   # Application routing
│   └── router.go             # Main router setup and route registration
├── testing/                  # Generic test helpers
│   └── assertions.go         # Type-safe test assertions using Go generics
└── constants/                # Common constants and error codes
    └── errors.go             # Shared error codes and messages
```

### Key Components

#### Common DTOs (`shared/dto/common_dto.go`)
Generic Response patterns, validation types, and error structures:

```go
// Import in domain services
import "thothix-backend/internal/shared/dto"

// Use generic Response types
func (s *SomeService) SomeOperation(req *SomeRequest) *dto.Response[*SomeResponse] {
    return dto.NewResponse(func() dto.Validation[*SomeResponse] {
        // Business logic here
        return dto.Success(response)
    })
}
```

#### Context Wrapper (`shared/handlers/context_wrapper.go`)
Standardized HTTP response helpers:

```go
// Import in handlers
import sharedHandlers "thothix-backend/internal/shared/handlers"

func (h *SomeHandler) SomeEndpoint(c *gin.Context) {
    wrapper := sharedHandlers.WrapContext(c)

    // Standardized responses
    wrapper.SuccessResponse(data)
    wrapper.ValidationErrorResponse(errors, "Operation failed")
    wrapper.NotFoundErrorResponse("Resource", id)
}
```

#### Application Router (`shared/router/router.go`)
Centralized route registration and middleware setup:

```go
// Main application entry point
import "thothix-backend/internal/shared/router"

func main() {
    r := router.Setup(db, cfg)
    r.Run(":8080")
}
```

#### Generic Test Assertions (`shared/testing/assertions.go`)
Type-safe test helpers for any domain:

```go
// Example usage in any domain test
func TestSomeDomainService(t *testing.T) {
    response := domainService.SomeOperation(request)

    // Assert success and get the typed result
    result := sharedTesting.AssertSuccessWithValue(t, response.Response)
    assert.Equal(t, "expected", result.SomeField)

    // Assert validation error with specific code
    sharedTesting.AssertValidationErrorWithCode(t, response.Response, "VALIDATION_ERROR")
}
```

#### Common Constants (`shared/constants/errors.go`)
Shared error codes and messages:

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

### Benefits

- **Consistency**: Standardized patterns across all domains
- **Reusability**: Generic components that work with any domain type
- **Maintainability**: Single source of truth for shared functionality
- **Type Safety**: Go generics ensure compile-time type safety
- **DRY Principle**: Eliminates code duplication across domains

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+ (or Docker)
- Node.js 16+ (for automation scripts)

### Local Development

```bash
# 1. Clone and setup
git clone <repository-url>
cd Thothix

# 2. Install Node.js dependencies (for automation)
npm install

# 3. Setup environment
cp .env.example .env
# Edit .env with your configuration

# 4. Start services with npm scripts
npm run dev        # Full Docker environment
# OR for backend-only development:

# 5. Backend-only development
cd backend
go mod tidy
go run main.go
```

### Docker Development (Recommended)

```bash
# Complete environment (API + Database + Vault)
npm run dev

# Backend logs
npm run dev:logs backend

# Database operations
npm run db:status
npm run db:connect
```

## 📁 Project Structure

```txt
backend/
├── internal/                      # Private application code
│   ├── config/                   # Configuration management
│   │   └── config.go             # App config and environment variables
│   ├── database/                 # Database setup and migrations
│   │   └── database.go           # GORM setup and connection
│   ├── dto/                      # Common/Shared DTOs
│   │   ├── common_dto.go         # Result Pattern & generic types
│   │   └── result_test.go        # Common DTO tests
│   ├── models/                   # Shared data models and database schema
│   │   ├── base.go               # Base model with common fields
│   │   ├── chat.go               # Chat/Channel models
│   │   ├── message.go            # Message models
│   │   ├── migrate.go            # Database migrations
│   │   ├── project.go            # Project models
│   │   ├── roles.go              # Role and permission models
│   │   └── user.go               # Legacy user model (being phased out)
│   ├── middleware/               # Shared middleware
│   │   ├── clerk_auth.go         # Clerk authentication middleware
│   │   ├── context.go            # Request context management
│   │   ├── cors.go               # CORS configuration
│   │   ├── logger.go             # Logging middleware
│   │   └── rbac.go               # Role-based access control
│   ├── router/                   # Route definitions
│   │   └── router.go             # Gin router setup and route registration
│   ├── vault/                    # Vault integration
│   │   └── client.go             # Vault client and secret management
│   └── users/                    # Users Domain Slice (Vertical Slice)
│       ├── domain/               # User business entities
│       │   ├── user.go           # User domain model with business logic
│       │   └── user_test.go      # Domain logic and business rules tests
│       ├── dto/                  # User-specific DTOs
│       │   ├── user_dto.go       # User request/response contracts
│       │   └── user_dto_test.go  # User DTO validation tests
│       ├── service/              # User business logic
│       │   ├── user_service.go           # User business operations
│       │   ├── user_service_interface.go # User service contract
│       │   └── user_service_test.go      # User business logic tests
│       ├── handlers/             # User HTTP controllers
│       │   ├── user_handler.go   # User HTTP endpoints
│       │   ├── user_handler_test.go # User controller tests
│       │   └── routes.go         # User route registration
│       └── mappers/              # User data transformation
│           ├── user_mapper.go    # User entity ↔ DTO mapping
│           └── user_mapper_test.go # User mapping tests
├── tests/                        # Integration tests
│   └── integration/
│       └── user_integration_test.go # End-to-end user API tests
├── docs/                         # Generated Swagger documentation
│   ├── docs.go                   # Generated Swagger Go code
│   ├── swagger.json              # Generated Swagger JSON
│   └── swagger.yaml              # Generated Swagger YAML
├── main.go                       # Application entry point
├── go.mod                        # Go module dependencies
├── go.sum                        # Go module checksums
├── Dockerfile                    # Multi-stage Docker build
└── README.md                     # This file
```

### 🧩 Vertical Slice Architecture Layers

#### **1. Domain Layer** (`users/domain/`)

- **Business Entities**: Core domain models with embedded business logic
- **Domain Rules**: Validation and business constraints within entities
- **Pure Business Logic**: No external dependencies (database, HTTP, etc.)
- **Testable**: Domain logic tested in isolation

#### **2. DTO Layer** (`users/dto/`)

- **API Contracts**: Request/response structures for HTTP endpoints
- **Domain-Specific**: DTOs tailored to the specific business domain
- **Validation**: Input validation rules and constraints
- **Type Safety**: Strongly typed interfaces with the outside world

#### **3. Service Layer** (`users/service/`)

- **Business Operations**: Orchestrates domain entities and external dependencies
- **Use Cases**: Implements specific business scenarios (CreateUser, GetUser, etc.)
- **Dependency Injection**: Services implement interfaces for easy testing
- **Result Pattern**: Type-safe error handling with functional composition

#### **4. Handlers Layer** (`users/handlers/`)

- **HTTP Controllers**: Handle HTTP requests and responses
- **Route Registration**: Self-contained route setup in `routes.go`
- **Error Handling**: Standardized error responses using ContextWrapper
- **Clean Interface**: Uses only DTOs and services, no direct model access

#### **5. Mappers Layer** (`users/mappers/`)

- **Data Transformation**: Convert between domain entities and DTOs
- **Explicit Mapping**: Clear, testable transformation logic
- **Isolation**: Prevents data structure coupling between layers
- **Bi-directional**: Supports both entity → DTO and DTO → entity mapping

- **Services**: Business logic and domain rules
- **DTOs**: Data contracts for API input/output
- **Mappers**: Explicit transformation between DTOs and models
- **Dependency Injection**: Services implement interfaces for testability

#### **3. Data Layer** (`models/`, `database/`)

- **Models**: GORM entities representing database schema
- **Database**: Connection management and configuration
- **Migrations**: Schema versioning and updates

## 🎯 Vertical Slice Architecture Implementation

### Domain-Driven Organization

Each business feature is organized as a **vertical slice** containing all necessary layers:

```go
// Example: User domain slice structure
internal/users/
├── domain/     → Business entities with core logic
├── dto/        → API contracts specific to users
├── service/    → User business operations
├── handlers/   → User HTTP endpoints
└── mappers/    → User data transformation
```

### Modular Route Registration

Each domain slice registers its own routes independently:

```go
// users/handlers/routes.go
func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
    userService := service.NewUserService(db)
    userHandler := NewUserHandler(userService)

    users := router.Group("/users")
    users.GET("", userHandler.GetUsers)
    users.GET("/:id", userHandler.GetUserByID)
    users.POST("", userHandler.CreateUser)
    // ... more routes
}

// router/router.go - main router
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
    // ... middleware setup

    v1 := r.Group("/api/v1")
    protected := v1.Group("/").Use(middleware.ClerkAuthSDK())

    // Register domain slices
    userHandlers.RegisterUserRoutes(protected, db)
    // Future: projectHandlers.RegisterProjectRoutes(protected, db)
    // Future: chatHandlers.RegisterChatRoutes(protected, db)
}
```

### Benefits of Vertical Slice Architecture

1. **Feature Isolation**: Complete features can be developed, tested, and deployed independently
2. **Clear Ownership**: Each slice has clear boundaries and responsibilities
3. **Parallel Development**: Teams can work on different slices without conflicts
4. **Easy Feature Addition**: New domains follow the same established pattern
5. **Reduced Coupling**: Minimal dependencies between different business domains
6. **Single Responsibility**: Each slice handles one business concern completely

### Domain Layer Design

```go
// users/domain/user.go - Business entity with embedded logic
type User struct {
    models.BaseModel
    ClerkID   string `json:"clerk_id" gorm:"uniqueIndex;not null"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    Username  string `json:"username,omitempty"`
    AvatarURL string `json:"avatar_url,omitempty"`
}

// Domain methods embedded in the entity
func (u *User) TableName() string { return "users" }

func (u *User) SyncFromClerk(clerkData *ClerkUserData) {
    u.ClerkID = clerkData.ID
    u.Email = clerkData.Email
    u.Name = clerkData.FirstName + " " + clerkData.LastName
    // ... business logic
}
```

### Service Layer Pattern

```go
// users/service/user_service_interface.go - Contract
type UserServiceInterface interface {
    GetUserByID(userID string) *usersDto.GetUserResponse
    CreateUser(req *usersDto.CreateUserRequest) *usersDto.CreateUserResponse
    // ... other operations
}

// users/service/user_service.go - Implementation
type UserService struct {
    db *gorm.DB
}

func (s *UserService) GetUserByID(userID string) *usersDto.GetUserResponse {
    return usersDto.NewGetUserResponse(func() dto.Validation[*usersDto.UserResponse] {
        // Validation and business logic
        if userID == "" {
            return dto.Failure[*usersDto.UserResponse](
                dto.NewError("VALIDATION_ERROR", "User ID required", nil))
        }

        // Database operation
        var user domain.User
        if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
            // ... error handling
        }

        // Transform to DTO
        response := s.mapper.DomainToResponse(&user)
        return dto.Success(response)
    })
}
```

### Testing Strategy

Each layer has focused, isolated tests:

```go
// Domain tests - Pure business logic
func TestUser_SyncFromClerk(t *testing.T) {
    user := &domain.User{}
    clerkData := &domain.ClerkUserData{...}

    user.SyncFromClerk(clerkData)

    assert.Equal(t, clerkData.Email, user.Email)
}

// Service tests - Business operations with mocked dependencies
func TestUserService_GetUserByID(t *testing.T) {
    // Uses testcontainers for real database testing
    // Tests complete business logic flow
}

// Handler tests - HTTP endpoint behavior with mocked services
func TestUserHandler_GetUserByID(t *testing.T) {
    mockService := &MockUserService{}
    handler := NewUserHandler(mockService)
    // ... HTTP testing
}

// Integration tests - End-to-end API testing
func TestUserAPI_GetUserByID_Integration(t *testing.T) {
    // Complete API test with real database
}
```

## 🗄️ Data Models

### Core Models

All models inherit from `BaseModel` with audit fields:

```go
type BaseModel struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    CreatedBy *string   `json:"created_by,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedBy *string   `json:"updated_by,omitempty"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Model Definitions

#### **User Model**

- Synchronized with Clerk authentication
- Fields: `ClerkID`, `Email`, `FirstName`, `LastName`, `SystemRole`
- Relationships: Projects, Messages, Roles

#### **Project Model**

- Enterprise project management
- Fields: `Name`, `Description`, `Status`, `OwnerID`
- Relationships: Members, Channels

#### **Chat/Channel Model**

- Communication channels (public/private)
- Fields: `Name`, `Description`, `Type`, `ProjectID`, `IsPrivate`
- Relationships: Members, Messages

#### **Message Model**

- Channel messages and direct messages
- Fields: `Content`, `SenderID`, `ChannelID`, `RecipientID`
- Support for both channel and 1:1 messages

#### **Role Model**

- Role-based access control
- Fields: `UserID`, `ProjectID`, `ChannelID`, `Role`
- Granular permissions per resource

### Database Schema Verification

```bash
# Verify BaseModel implementation across all tables
npm run db:check

# Check specific table structure
npx zx scripts/db-verify.mjs check-table users

# Verify specific field exists
npx zx scripts/db-verify.mjs has-field users system_role
```

## 📡 API Reference

### Live Documentation

- **Swagger UI**: `http://localhost:30000/swagger/index.html`
- **API Base URL**: `http://localhost:30000/api/v1`

### Generate Documentation

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
cd backend
swag init

# Regenerate after API changes
npm run docs:generate  # Via npm script
```

### Core Endpoints

#### Authentication

- `POST /auth/sync` - Sync user with Clerk
- `GET /auth/me` - Current user information

#### Projects

- `GET /projects` - List accessible projects
- `POST /projects` - Create project (Manager/Admin)
- `GET /projects/{id}` - Project details
- `PUT /projects/{id}` - Update project
- `DELETE /projects/{id}` - Delete project

#### Channels/Chats

- `GET /channels` - List accessible channels
- `POST /channels` - Create channel (Manager/Admin)
- `GET /channels/{id}` - Channel details
- `POST /channels/{id}/join` - Join public channel
- `DELETE /channels/{id}/leave` - Leave channel

#### Messages

- `GET /channels/{id}/messages` - Channel messages
- `POST /channels/{id}/messages` - Send channel message
- `POST /messages/direct` - Send direct message
- `GET /messages/direct/{userId}` - Direct message history

#### Role Management (Admin Only)

- `GET /roles` - List all roles
- `POST /roles` - Assign role
- `PUT /roles/{id}` - Update role
- `DELETE /roles/{id}` - Revoke role

### HTTP Status Codes

- `200` - Success (GET, PUT, PATCH)
- `201` - Created (POST)
- `204` - No Content (DELETE)
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (authentication required)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

## 🔐 Authentication & RBAC

### Clerk Integration

Thothix uses Clerk for secure authentication:

```go
// Middleware validates Clerk JWT tokens
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractTokenFromHeader(c)
        clerkUser, err := validateClerkToken(token)
        // Set user context...
    }
}
```

### Role-Based Access Control (RBAC)

#### System Roles

- **Admin**: Full system access
- **Manager**: Project and channel management
- **User**: Standard user access
- **External**: Limited public access

#### Resource Roles

- **Project Member**: Access to project channels
- **Channel Member**: Access to specific private channels

For complete RBAC documentation, see [`RBAC_SIMPLIFIED.md`](./RBAC_SIMPLIFIED.md).

### Public/Private Channel Strategy

- **Public Channels**: No explicit membership required
- **Private Channels**: Explicit membership in `channel_members` table

## 💻 Development

### Environment Setup

```bash
# Install Go dependencies
go mod tidy

# Install development tools
go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### Development Workflow

```bash
# Format code
npm run format      # Cross-platform via npm
gofmt -w .         # Direct Go command

# Lint code
npm run lint       # Cross-platform via npm
golangci-lint run  # Direct Go command

# Run tests
npm run test       # Cross-platform via npm
go test ./...      # Direct Go command

# Pre-commit checks (format + lint + test)
npm run pre-commit

# Generate API documentation
swag init
```

### Clean Architecture Development Guidelines

When adding new features, follow this order:

#### **1. Define DTOs** (`internal/dto/`)

```go
// Define request/response contracts first
type CreateProjectRequest struct {
    Name        string `json:"name" binding:"required,min=2,max=100"`
    Description string `json:"description" binding:"max=500"`
}

type ProjectResponse struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### **2. Create/Update Models** (`internal/models/`)

```go
// Database entities (if needed)
type Project struct {
    BaseModel
    Name        string `gorm:"size:100;not null"`
    Description string `gorm:"size:500"`
    OwnerID     string `gorm:"size:50;not null"`
}
```

#### **3. Implement Mappers** (`internal/mappers/`)

```go
// DTO ↔ Model transformation
type ProjectMapperInterface interface {
    ModelToResponse(project *models.Project) *dto.ProjectResponse
    CreateRequestToModel(req *dto.CreateProjectRequest) *models.Project
}

// Write mapper tests first (TDD approach)
func TestProjectMapper_ModelToResponse(t *testing.T) {
    // Test implementation...
}
```

#### **4. Create Service Interface** (`internal/services/`)

```go
// Define business logic contract
type ProjectServiceInterface interface {
    CreateProject(req *dto.CreateProjectRequest) (*dto.ProjectResponse, error)
    GetProject(id string) (*dto.ProjectResponse, error)
    UpdateProject(id string, req *dto.UpdateProjectRequest) (*dto.ProjectResponse, error)
}
```

#### **5. Implement Service** (`internal/services/`)

```go
// Business logic implementation
type ProjectService struct {
    db     *gorm.DB
    mapper mappers.ProjectMapperInterface
}

// Write service tests with real PostgreSQL
func TestProjectService_CreateProject(t *testing.T) {
    db := setupTestDB(t) // Uses testcontainers
    service := NewProjectService(db)
    // Test implementation...
}
```

#### **6. Create Handlers** (`internal/handlers/`)

```go
// HTTP controllers using only DTOs and services
type ProjectHandler struct {
    projectService services.ProjectServiceInterface
}

func NewProjectHandler(service services.ProjectServiceInterface) *ProjectHandler {
    return &ProjectHandler{projectService: service}
}

// routes.go - Route registration
func RegisterProjectRoutes(router *gin.RouterGroup, db *gorm.DB) {
    projectService := service.NewProjectService(db)
    projectHandler := NewProjectHandler(projectService)

    projects := router.Group("/projects")
    projects.GET("", projectHandler.GetProjects)
    projects.GET("/:id", projectHandler.GetProjectByID)
    projects.POST("", projectHandler.CreateProject)
}
```

#### **7. Register Routes in Main Router**

```go
// internal/router/router.go
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
    // ... existing setup

    v1 := r.Group("/api/v1")
    protected := v1.Group("/").Use(middleware.ClerkAuthSDK())

    // Register domain slices
    userHandlers.RegisterUserRoutes(protected, db)
    projectHandlers.RegisterProjectRoutes(protected, db)  // Add new domain

    return r
}
```

#### **8. Add Tests for Each Layer

```bash
# Create test files
touch internal/projects/domain/project_test.go
touch internal/projects/dto/project_dto_test.go
touch internal/projects/service/project_service_test.go
touch internal/projects/handlers/project_handler_test.go
touch internal/projects/mappers/project_mapper_test.go
touch tests/integration/project_integration_test.go
```

### Key Principles for New Domains

1. **Self-Contained**: Each domain slice should be independent
2. **Consistent Structure**: Follow the established 5-layer pattern
3. **Clear Contracts**: Define interfaces for services and mappers
4. **Comprehensive Testing**: Test each layer in isolation
5. **Route Registration**: Use the modular registration pattern
6. **Error Handling**: Use the Result Pattern consistently
7. **Validation**: Implement input validation in DTOs

This pattern ensures consistency, maintainability, and allows teams to work on different domains independently.
