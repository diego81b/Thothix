# 🚀 Thothix Backend - Go API

[![Go](https://img.shields.io/badge/Go-1.23-blue?style=flat&logo=go)](https://golang.org) [![Gin](https://img.shields.io/badge/Gin-Web_Framework-green?style=flat)](https://github.com/gin-gonic/gin) [![GORM](https://img.shields.io/badge/GORM-ORM-yellow?style=flat)](https://gorm.io)

Modern Go REST API with Gin framework, GORM ORM, PostgreSQL, Clerk authentication, and HashiCorp Vault integration.

## 🏗️ Architecture

**Tech Stack**: Gin • GORM v2 • PostgreSQL 15 • Clerk Auth • HashiCorp Vault • Swagger/OpenAPI • Docker

**Clean Architecture**: Presentation Layer (handlers/middleware/router) → Business Layer (services/dto/mappers) → Data Layer (models/database)

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

**Architecture**: `common_dto.go` (generic patterns) + `user_dto.go` (domain DTOs)

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
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   │   └── config.go      # App config and environment variables
│   ├── database/          # Database setup and migrations
│   │   └── database.go    # GORM setup and connection
│   ├── dto/               # Data Transfer Objects (Clean Architecture)
│   │   ├── user_dto.go    # User request/response DTOs
│   │   ├── project_dto.go # Project DTOs
│   │   ├── chat_dto.go    # Chat/Channel DTOs
│   │   └── common_dto.go  # Common DTOs (error, pagination)
│   ├── services/          # Business Logic Layer (Clean Architecture)
│   │   ├── user_service.go          # User business logic
│   │   ├── user_service_interface.go # User service contract
│   │   ├── user_service_test.go     # User service tests
│   │   ├── project_service.go       # Project business logic
│   │   └── chat_service.go          # Chat business logic
│   ├── mappers/           # DTO ↔ Model Mapping (Clean Architecture)
│   │   ├── user_mapper.go          # User DTO/Model mapping
│   │   ├── user_mapper_test.go     # User mapper tests
│   │   ├── project_mapper.go       # Project DTO/Model mapping
│   │   └── chat_mapper.go          # Chat DTO/Model mapping
│   ├── handlers/          # HTTP handlers (controllers)
│   │   ├── auth.go        # Authentication endpoints
│   │   ├── chats.go       # Chat/Channel management
│   │   ├── health.go      # Health check endpoint
│   │   ├── messages.go    # Message endpoints
│   │   ├── projects.go    # Project management
│   │   ├── roles.go       # Role-based access control
│   │   └── users.go       # User management
│   ├── middleware/        # Custom middleware
│   │   ├── auth.go        # Clerk authentication middleware
│   │   ├── context.go     # Request context management
│   │   ├── cors.go        # CORS configuration
│   │   ├── logger.go      # Logging middleware
│   │   └── rbac.go        # Role-based access control
│   ├── models/            # Data models and database schema
│   │   ├── base.go        # Base model with common fields
│   │   ├── chat.go        # Chat/Channel models
│   │   ├── message.go     # Message models
│   │   ├── migrate.go     # Database migrations
│   │   ├── project.go     # Project models
│   │   ├── roles.go       # Role and permission models
│   │   └── user.go        # User models
│   ├── router/            # Route definitions
│   │   └── router.go      # Gin router setup and routes
│   └── vault/             # Vault integration
│       └── client.go      # Vault client and secret management
├── docs/                  # Generated Swagger documentation
│   ├── docs.go           # Generated Swagger Go code
│   ├── swagger.json      # Generated Swagger JSON
│   └── swagger.yaml      # Generated Swagger YAML
├── main.go               # Application entry point
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
├── Dockerfile           # Multi-stage Docker build
└── README.md           # This file
```

### Clean Architecture Layers

#### **1. Presentation Layer** (`handlers/`, `middleware/`, `router/`)

- **Handlers**: HTTP request/response controllers
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)
- **Router**: Route definitions and middleware setup
- **Input/Output**: Uses DTOs exclusively, no direct model access

#### **2. Business Layer** (`services/`, `dto/`, `mappers/`)

- **Services**: Business logic and domain rules
- **DTOs**: Data contracts for API input/output
- **Mappers**: Explicit transformation between DTOs and models
- **Dependency Injection**: Services implement interfaces for testability

#### **3. Data Layer** (`models/`, `database/`)

- **Models**: GORM entities representing database schema
- **Database**: Connection management and configuration
- **Migrations**: Schema versioning and updates

## 🎯 Clean Architecture Implementation

### Data Transfer Objects (DTOs)

DTOs define the API contract and decouple external interfaces from internal models

### Service Layer

Services contain business logic and implement interfaces for dependency injection

### Mapper Layer

Mappers handle explicit transformation between DTOs and models

### Handler Layer

Handlers use only DTOs and services, never accessing models or database directly

### Dependency Injection

Services are injected into handlers via interfaces

### Benefits of This Architecture

1. **Separation of Concerns**: Clear boundaries between layers
2. **Testability**: Services implement interfaces, easy to mock
3. **Maintainability**: Changes to models don't affect handlers
4. **API Stability**: DTOs provide stable external contracts
5. **Validation**: Input validation at DTO level
6. **Type Safety**: Explicit mapping prevents data leaks

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

func (h *ProjectHandler) CreateProject(c *gin.Context) {
    var req dto.CreateProjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, dto.ErrorResponse{Error: "validation_failed"})
        return
    }

    project, err := h.projectService.CreateProject(&req)
    if err != nil {
        c.JSON(500, dto.ErrorResponse{Error: "creation_failed"})
        return
    }

    c.JSON(201, project)
}
```

#### **7. Update Router** (`internal/router/`)

```go
// Dependency injection and route setup
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
    // Initialize services
    projectService := services.NewProjectService(db)

    // Initialize handlers with service interfaces
    projectHandler := handlers.NewProjectHandler(projectService)

    // Setup routes
    projects := protected.Group("/projects")
    projects.POST("", projectHandler.CreateProject)
}
```

#### **8. Update Swagger Documentation**

```go
// Add Swagger comments to handlers
// @Summary Create project
// @Tags projects
// @Accept json
// @Produce json
// @Param project body dto.CreateProjectRequest true "Project data"
// @Success 201 {object} dto.ProjectResponse
// @Router /api/v1/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) { ... }
```

### Hot Reload Development

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Start with hot reload
air

# Or via Docker with hot reload
npm run dev
```

### Environment Variables

Key environment variables for development:

```bash
# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=thothix_dev
POSTGRES_USER=thothix_user
POSTGRES_PASSWORD=secure_password

# Clerk Authentication
CLERK_SECRET_KEY=sk_test_...
CLERK_PUBLISHABLE_KEY=pk_test_...

# Vault Integration
USE_VAULT=true
VAULT_ADDR=http://localhost:8200
VAULT_ROOT_TOKEN=thothix-secure-root-token

# Application
PORT=8080
GIN_MODE=debug
LOG_LEVEL=debug
```

### Database Operations

```bash
# Database status and connection
npm run db:status
npm run db:connect

# Schema verification
npm run db:check           # Verify BaseModel fields
npm run db:tables         # List all tables

# Advanced database operations
npx zx scripts/db-verify.mjs check-table users
npx zx scripts/db-verify.mjs has-field projects owner_id
npx zx scripts/db-verify.mjs missing-field channels is_private
```

## 🧪 Testing

### Test Architecture

The testing strategy follows the clean architecture layers:

```
┌─────────────────────────────────────────────────────────────┐
│                    INTEGRATION TESTS                        │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │
│  │   Handler   │────│   Service   │────│  Database   │      │
│  │   Tests     │    │   Tests     │    │(PostgreSQL  │      │
│  │ (E2E API)   │    │(with real   │    │Testcont.)   │      │
│  │             │    │ database)   │    │             │      │
│  └─────────────┘    └─────────────┘    └─────────────┘      │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                      UNIT TESTS                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │
│  │   Mapper    │────│   DTO       │────│  Business   │      │
│  │   Tests     │    │Validation   │    │   Logic     │      │
│  │ (Pure Go)   │    │   Tests     │    │   Tests     │      │
│  └─────────────┘    └─────────────┘    └─────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### Unit Tests

```bash
# Run all tests
npm run test
go test ./...

# Run tests with debug output
npm run test --debug

# Run specific test packages
go test ./internal/mappers/...     # Mapper tests (no database)
go test ./internal/services/...    # Service tests (with PostgreSQL)
```

### Integration Tests with PostgreSQL

Services use **testcontainers** with real PostgreSQL for realistic testing:

```go
func setupTestDB(t *testing.T) *gorm.DB {
    ctx := context.Background()

    // Create PostgreSQL container
    postgresContainer, err := pgtest.Run(ctx,
        "docker.io/postgres:16-alpine",
        pgtest.WithDatabase("testdb"),
        pgtest.WithUsername("testuser"),
        pgtest.WithPassword("testpass"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(30*time.Second)),
    )
    require.NoError(t, err)

    // Auto-cleanup after test
    t.Cleanup(func() {
        postgresContainer.Terminate(ctx)
    })

    // Connect with GORM
    connStr, _ := postgresContainer.ConnectionString(ctx, "sslmode=disable")
    db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
    require.NoError(t, err)

    // Auto-migrate schema
    db.AutoMigrate(&models.User{})

    return db
}
```

### Test Structure

```txt
backend/
├── internal/
│   ├── mappers/
│   │   ├── user_mapper.go
│   │   └── user_mapper_test.go          # Pure unit tests
│   ├── services/
│   │   ├── user_service.go
│   │   └── user_service_test.go         # Integration tests with PostgreSQL
│   ├── dto/
│   │   └── user_dto_test.go             # Validation tests
│   └── handlers/
│       └── user_handler_test.go         # API endpoint tests
└── tests/
    ├── integration/                     # Full E2E tests
    ├── fixtures/                        # Test data
    └── testcontainers/                  # Container configurations
```

### Test Commands

```bash
# Standard test run
npm run test

# Debug mode with race detection and coverage
npm run test -- --debug

# Test specific layers
go test ./internal/mappers/...     # Unit tests (fast)
go test ./internal/services/...    # Integration tests (with DB)
go test ./internal/handlers/...    # API tests

# Test with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Continuous testing (with file watching)
air test
```

### Test Benefits

1. **Realistic Testing**: Uses same PostgreSQL as production
2. **Isolation**: Each test gets a fresh database container
3. **Fast Feedback**: Parallel test execution
4. **Coverage**: Tests all architectural layers
5. **CI/CD Ready**: Works in any Docker environment

## 🚀 Deployment

### Docker Build

```bash
# Local build
docker build -t thothix/backend:latest .

# Multi-stage production build
docker build --target prod -t thothix/backend:prod .

# Development build with hot reload
docker build --target dev -t thothix/backend:dev .
```

### Environment Deployment

```bash
# Development
npm run dev

# Staging
npm run staging

# Production
npm run prod
```

### Health Checks

The API includes health check endpoints:

- `GET /health` - Basic health status
- `GET /health/ready` - Readiness check (includes database)
- `GET /health/live` - Liveness check

### Monitoring

```bash
# Container logs
npm run dev:logs backend

# Database logs
npm run dev:logs postgres

# Vault logs
npm run dev:logs vault

# Resource usage
docker stats thothix-backend-dev
```

## 🔧 Configuration

### Gin Framework Configuration

```go
// Production mode
gin.SetMode(gin.ReleaseMode)

// Development mode with detailed logging
gin.SetMode(gin.DebugMode)
```

### GORM Configuration

```go
// Database connection with optimizations
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
    NamingStrategy: schema.NamingStrategy{
        SingularTable: false,
    },
})
```

### Middleware Stack

1. **CORS** - Cross-origin resource sharing
2. **Logger** - Request/response logging
3. **Recovery** - Panic recovery
4. **Auth** - Clerk authentication
5. **RBAC** - Role-based access control
6. **Context** - Request context management

## 📚 Additional Resources

- **[Main Project README](../README.md)** - Overall project documentation
- **[Vault Integration](../docs/VAULT_INTEGRATION.md)** - HashiCorp Vault setup
- **[Clerk Integration](../docs/CLERK_INTEGRATION.md)** - Authentication setup
- **[Node.js Development Guide](../docs/NODE_JS_GUIDE.md)** - Cross-platform automation
- **[Database Migration](../docs/DB_MIGRATION.md)** - Database schema management

---

## 🤝 Contributing

1. Follow Go best practices and project conventions
2. Write tests for new features
3. Update Swagger documentation for API changes
4. Use the provided npm scripts for consistency
5. Ensure all pre-commit checks pass

```bash
# Before committing
npm run pre-commit
```

---

**Built with ❤️ using Go, Gin, GORM, and modern DevOps practices.**
