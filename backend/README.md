# Thothix Backend - Enterprise Messaging Platform

This is the backend for the Thothix enterprise messaging platform, developed in Go with the Gin framework and Clerk integration for authentication.

## Data Models

The backend uses the following main models:

- **User**: Platform users (synchronized with Clerk)
- **Project**: Enterprise projects 
- **ProjectMember**: Members assigned to projects with roles
- **Channel**: Communication channels within projects
- **ChannelMember**: Channel members
- **Message**: Messages (in channels or direct between users)
- **File**: Shared files in projects

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
├── main.go            # Application entry point
├── go.mod            # Go dependencies
├── Dockerfile        # Docker configuration
└── .env.example      # Environment variables template
```

## Quick Start

1. **Configure environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your Clerk keys
   ```

2. **Start with Docker**:
   ```bash
   cd ..
   docker-compose up -d
   ```

3. **Verify it's working**:
   ```bash
   curl http://localhost:30000/health
   ```

## API Endpoints

- **Health**: `GET /health`
- **Swagger**: `GET /swagger/index.html`
- **API**: `GET /api/v1/auth/sync`, `GET /api/v1/auth/me`
- **Users**: `GET /api/v1/users`, `PUT /api/v1/users/me`
- **Projects**: `GET /api/v1/projects` (TODO)
- **Channels**: `GET /api/v1/channels` (TODO)
- **Messages**: `GET /api/v1/channels/{id}/messages` (TODO)

## Development

For local development:

```bash
# Install dependencies
go mod tidy

# Generate Swagger documentation
go install github.com/swaggo/swag/cmd/swag@latest
swag init

# Start the server
go run main.go
```

## Clerk Authentication

The backend uses Clerk for authentication. Users are automatically synchronized to the local database on first access.

Authentication flow:
1. User authenticates via Clerk in the frontend
2. Frontend sends Clerk token in API requests
3. `ClerkAuth` middleware verifies token with Clerk
4. `/auth/sync` endpoint creates/updates user in local database

## TODO

- [ ] Implement project handlers
- [ ] Implement channel handlers  
- [ ] Implement message handlers
- [ ] Add WebSocket for real-time
- [ ] Add file upload support (MinIO)
- [ ] Add unit tests
- [ ] Improve API documentation
