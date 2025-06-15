# 🚀 Thothix Backend - Go API

[![Go](https://img.shields.io/badge/Go-1.23-blue?style=flat&logo=go)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-Web_Framework-green?style=flat)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-ORM-yellow?style=flat)](https://gorm.io)

Thothix backend is a modern Go REST API built with Gin framework, GORM ORM, and PostgreSQL, featuring Clerk authentication and HashiCorp Vault integration.

## 📋 Table of Contents

- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Data Models](#data-models)
- [API Reference](#api-reference)
- [Authentication & RBAC](#authentication--rbac)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)

## 🏗️ Architecture

### Technology Stack

- **Framework**: Gin Web Framework
- **ORM**: GORM v2
- **Database**: PostgreSQL 15
- **Authentication**: Clerk Integration
- **Secrets**: HashiCorp Vault
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker multi-stage builds

### Service Dependencies

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Thothix Web    │    │  Thothix API    │    │   PostgreSQL    │
│   (Frontend)    │───▶│   (Backend)     │───▶│   (Database)    │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │  HashiCorp      │
                       │  Vault          │
                       │  (Secrets)      │
                       └─────────────────┘
```

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

```
backend/
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   │   └── config.go      # App config and environment variables
│   ├── database/          # Database setup and migrations
│   │   └── database.go    # GORM setup and connection
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

### Unit Tests

```bash
# Run all tests
npm run test
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/handlers
go test ./internal/models
```

### Integration Tests

```bash
# Start test environment
docker compose -f docker-compose.test.yml up -d

# Run integration tests
go test -tags=integration ./...

# Cleanup test environment
docker compose -f docker-compose.test.yml down -v
```

### Test Structure

```
backend/
├── internal/
│   ├── handlers/
│   │   ├── auth_test.go
│   │   ├── projects_test.go
│   │   └── messages_test.go
│   ├── models/
│   │   ├── user_test.go
│   │   └── project_test.go
│   └── middleware/
│       └── auth_test.go
└── tests/
    ├── integration/
    └── fixtures/
```

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
