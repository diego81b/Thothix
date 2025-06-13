# Thothix Backend

This directory contains the Go backend for the Thothix enterprise messaging platform.

> **Main Documentation**: See the [main README](../README.md) for complete setup instructions, Docker usage, and project overview.

## Quick Backend Development

For backend-specific development:

```bash
# Install dependencies
go mod tidy

# Generate Swagger documentation
go install github.com/swaggo/swag/cmd/swag@latest
swag init

# Start the server (requires database)
go run main.go
```

## Project Structure

```
backend/
├── internal/
│   ├── config/         # Application configuration
│   ├── database/       # Database setup and migrations
│   ├── handlers/       # HTTP handlers for APIs
│   ├── middleware/     # Custom middleware
│   ├── models/         # Data models
│   └── router/         # Route setup
├── docs/              # Generated Swagger documentation
├── main.go            # Application entry point
├── go.mod            # Go dependencies
└── Dockerfile        # Docker configuration
```

## Data Models

- **User**: Platform users (synchronized with Clerk)
- **Project**: Enterprise projects with members
- **Channel**: Communication channels (public/private)
- **Message**: Channel messages and direct messages
- **File**: Shared files in projects

## Key Features

- **Clerk Integration**: JWT authentication with user sync
- **RBAC System**: Admin, Manager, User, External roles
- **RESTful API**: Documented with Swagger/OpenAPI
- **Docker Ready**: Multi-stage build with health checks

## API Documentation

- **Swagger UI**: `http://localhost:30000/swagger/index.html`
- **Health Check**: `GET /health`
- **RBAC Details**: See [RBAC_SIMPLIFIED.md](RBAC_SIMPLIFIED.md)

For complete setup instructions, Docker usage, and deployment, see the [main project README](../README.md).
