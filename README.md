# Thothix Messaging Platform

[![GitHub Repository](https://img.shields.io/badge/GitHub-diego81b%2FThothix-blue?style=flat&logo=github)](https://github.com/diego81b/Thothix)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue?style=flat&logo=docker)](https://docker.com)
[![Go](https://img.shields.io/badge/Go-1.23-blue?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue?style=flat&logo=postgresql)](https://postgresql.org)

Thothix is a modern enterprise messaging platform that enables project management, group chats, and private 1:1 conversations. Built with containerized microservices architecture.

## 🚀 GitHub Repository

The project is available on GitHub: **[https://github.com/diego81b/Thothix](https://github.com/diego81b/Thothix)**

```bash
# Clone the repository
git clone https://github.com/diego81b/Thothix.git
cd Thothix
```

## 📋 Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start with Docker](#quick-start-with-docker)
- [🔐 HashiCorp Vault Integration](#hashicorp-vault-integration)
- [Docker Configuration](#docker-configuration)
- [Useful Docker Commands](#useful-docker-commands)
- [Database Verification Tools](#database-verification-tools)
- [Development](#development)
- [API Reference](#api-reference)
- [Contributing](#contributing)

### 📖 Additional Documentation

- **[🔐 Vault Integration Guide](./VAULT_INTEGRATION.md)** - Complete setup, troubleshooting & production guide
- **[� Docker Modernization Guide](./DOCKER_MODERNIZATION.md)** - Docker architecture updates and migration guide
- **[�📁 Scripts Documentation](./scripts/README.md)** - Development automation and tools
- **[🤖 Automation Guide](./AUTOMATION.md)** - Pre-commit hooks and formatting

## 🏗️ Architecture

### Technology Stack

- **Frontend**: Vue.js 3 + Nuxt.js + Nuxt UI
- **Backend**: Go + Gin Framework + GORM
- **Database**: PostgreSQL 15
- **File Storage**: MinIO (S3-compatible)
- **Orchestration**: Kubernetes
- **Containerization**: Docker

### Service Structure

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  thothix-web    │    │  thothix-api    │    │ thothix-postgres│
│   (Nuxt.js)     │───▶│     (Go/Gin)    │───▶│   (Database)    │
│   Port: 30001   │    │   Port: 30000   │    │   Port: 5432    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │  thothix-vault  │    │  thothix-minio  │
                       │ (Secrets Mgmt)  │    │ (File Storage)  │
                       │   Port: 8200    │    │   Port: 30002   │
                       └─────────────────┘    │ Console: 30003  │
                                              └─────────────────┘
```

> 🚀 **Recent Updates**: Docker configuration has been modernized with multi-stage builds and consistent naming conventions. See [Docker Modernization Guide](./DOCKER_MODERNIZATION.md) for details.

### Main Features

- 💼 **Project Management**: Creation and administration of enterprise projects
- 👥 **Group Chats**: Public and private channels for projects
- 💬 **1:1 Messages**: Private conversations between users
- 📁 **File Sharing**: Document upload and sharing via MinIO
- 🔐 **User Management**: Authentication and authorization system
- 📱 **Responsive**: Interface optimized for desktop and mobile

## 🚀 Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (version 20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 2.0+)
- Git

### Installation Verification

```bash
docker --version
docker compose version
```

## 🐳 Quick Start with Docker

### 1. Repository Clone

```bash
git clone <repository-url>
cd Thothix
```

### 2. Environment Configuration

```bash
# Copy and configure environment variables
cp .env.example .env
notepad .env
```

**Configure your `.env` file with your credentials:**

📋 **Reference**: See [`.env.example`](./.env.example) for all available configuration options including:

- Database settings
- Clerk authentication keys  
- Application configuration
- Vault integration settings

⚠️ **Security Note**: Never commit the `.env` file to version control. It's already included in `.gitignore`.

**For Vault Integration (Optional):**

📖 **Full Setup Guide**: See [Vault Integration Guide](./VAULT_INTEGRATION.md) for complete configuration.

Vault service is available but controlled by `USE_VAULT=true` in your `.env` file.

⚠️ **Security Note**: Never commit the `.env` file to version control. It's already included in `.gitignore`.

### 3. Start the complete stack

```bash
# Start all services (development mode)
docker compose up -d --build

# Verify all containers are running
docker compose ps
```

### 4. Database initialization

```bash
# Run database migrations
docker compose exec thothix-api go run cmd/migrate/main.go

# Load sample data (optional)
docker compose exec thothix-api go run cmd/seed/main.go
```

### 5. Access services

| Service           | URL                                         | Credentials        |
| ----------------- | ------------------------------------------- | ------------------ |
| **API Swagger**   | <http://localhost:30000/swagger/index.html> | -                  |
| **Thothix Web**   | <http://localhost:30001>                    | -                  |
| **MinIO Console** | <http://localhost:30002>                    | admin/password123  |
| **pgAdmin**       | <http://localhost:5432>                     | postgres/@Admin123 |

---

## 🔐 HashiCorp Vault Integration

Vault è integrato per gestire in modo sicuro segreti e configurazioni sensibili.

📖 **[📋 Guida Completa all'Integrazione Vault →](./VAULT_INTEGRATION.md)**

### Quick Start

1. **Abilita Vault**: Imposta `USE_VAULT=true` nel tuo `.env`
2. **Avvia servizi**: `docker compose up -d --build`
3. **Accedi a Vault UI**: <http://localhost:8200> (token: da tuo `.env`)

Vault gestisce automaticamente credenziali database, API keys Clerk e segreti applicazione.

Per setup completo, troubleshooting e configurazione produzione → **[VAULT_INTEGRATION.md](./VAULT_INTEGRATION.md)**

---

## ⚙️ Docker Configuration

### 📋 Configuration Files

**Docker Compose:**

- 🔧 [`docker-compose.yml`](./docker-compose.yml) - Development configuration
- 🚀 [`docker-compose.prod.yml`](./docker-compose.prod.yml) - Production configuration

**Dockerfiles:**

- 🐳 [`Dockerfile.backend`](./Dockerfile.backend) - API backend (multi-stage)
- 🐳 [`backend/Dockerfile.backend`](./backend/Dockerfile.backend) - Alternative backend location
- 🗄️ [`Dockerfile.postgres`](./Dockerfile.postgres) - PostgreSQL database (multi-stage)
- 🔐 [`Dockerfile.vault`](./Dockerfile.vault) - HashiCorp Vault (multi-stage)

**Environment:**

- ⚙️ [`.env.example`](./.env.example) - Environment variables template

### 🎯 Key Features

- **Multi-Stage Builds**: Optimized for development and production
- **Consistent Naming**: `-dev` suffix for development, `-prod` for production
- **Health Checks**: Built-in monitoring for all services
- **Volume Management**: Persistent data storage
- **Network Isolation**: Secure inter-service communication

📖 **Detailed Guide**: See [Docker Modernization Guide](./DOCKER_MODERNIZATION.md) for migration details and best practices.

## 🛠️ Useful Docker Commands

### Development Commands

```bash
# Avvia tutti i servizi in background (development mode)
docker compose up -d

# Avvia tutti i servizi in produzione
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Ferma e rimuove container, network e volumi definiti
docker compose down

# Build di tutte le immagini per development
docker compose build

# Build di tutte le immagini per production
docker compose -f docker-compose.yml -f docker-compose.prod.yml build

# Lista immagini Thothix
docker images | findstr thothix

# Push al registry (quando configurato)
docker push thothix/api:1.0.0-dev
docker push thothix/api:1.0.0-prod
docker push thothix/postgres:17.5-thothix1.0-dev
docker push thothix/postgres:17.5-thothix1.0-prod
docker push thothix/vault:1.15.0-thothix1.0-dev
docker push thothix/vault:1.15.0-thothix1.0-prod

# Inspect metadata delle immagini
docker inspect thothix/api:1.0.0-dev --format="{{json .Config.Labels}}"
docker inspect thothix/postgres:17.5-thothix1.0-prod --format="{{json .Config.Labels}}"
```

### Build e deployment immagini

```bash
# Ricostruisci solo l’API
docker-compose build thothix-api

# Ricostruisci solo il database
docker-compose build thothix-postgres

# Ricostruisci TUTTE le immagini
docker-compose build

# Rinomina / crea tag aggiuntivo
docker tag thothix/api:1.0.0  thothix/api:latest

# Pubblica su registry (se configurato)
docker push thothix/api:1.0.0
docker push thothix/postgres:17.5-thothix1.0
```

### Aggiornamento Swagger

```bash
```bash
# Genera (o rigenera) la doc. OpenAPI dai commenti @swagger
cd backend
swag init -g main.go

# Riavvia solo il servizio API (con doc aggiornata)
docker compose restart thothix-api
```

### Logging e debug

```bash
# Segui i log del servizio API
docker compose logs -f thothix-api

# Segui i log del servizio Postgres
docker compose logs -f postgres

# Segui i log del servizio Vault
docker compose logs -f vault

# Esegui una shell interattiva nel container API
docker compose exec thothix-api sh

# Esegui comandi Vault
docker compose exec vault vault status
```
```

### Gestione servizi

```bash
# Avvia tutti i servizi
docker-compose up -d

# Avvia servizi specifici
docker-compose up -d thothix-postgres thothix-minio

# Rebuild dell'API dopo modifiche
docker-compose up -d --build thothix-api

# Restart del frontend
docker-compose restart thothix-web
```

### Debugging

```bash
# Log dell'API Go
docker-compose logs -f thothix-api

# Log del frontend Nuxt
docker-compose logs -f thothix-web

# Accesso shell API container
docker-compose exec thothix-api sh

# Accesso shell database
docker-compose exec thothix-postgres psql -U thothix_user -d thothix
```

### Database operations

```bash
# Backup database
docker-compose exec thothix-postgres pg_dump -U thothix_user thothix > backup_$(date +%Y%m%d).sql

# Restore database
docker-compose exec -T thothix-postgres psql -U thothix_user -d thothix < backup.sql

# Reset database
docker-compose down -v
docker volume rm thothix_postgres_data
docker-compose up -d thothix-postgres
```

## 🔍 Database Verification Tools

### Utility Scripts

The project includes scripts to easily verify database alignment with Go models:

#### Windows (cmd/PowerShell)

```cmd
# Verify BaseModel alignment (all tables should have 5 columns)
scripts\db-verify.bat check-basemodel

# List all tables
scripts\db-verify.bat list-tables

# Check structure of a specific table
scripts\db-verify.bat check-table users

# Find tables missing a specific field
scripts\db-verify.bat missing-field updated_by

# Find tables that have a specific field
scripts\db-verify.bat has-field system_role

# Connect to database interactively
scripts\db-verify.bat connect

# Database status
scripts\db-verify.bat status
```

#### Linux/MacOS (bash)

```bash
# Same usage but with .sh extension
chmod +x scripts/db-verify.bat
.\scripts\db-verify.bat check-basemodel
```

### Manual SQL Commands

To execute SQL commands directly:

```bash
# Direct database connection
docker-compose exec postgres psql -U postgres -d thothix-db

# Execute single command
docker-compose exec postgres psql -U postgres -d thothix-db -c "SELECT version();"
```

For more details on verification commands, see `DB_MIGRATION.md`.

## 💻 Development

### Project Structure

```
thothix/
├── backend/                 # Go API + Gin
│   ├── internal/           # Business logic
│   │   ├── config/         # Application configuration
│   │   ├── database/       # Database setup and migrations
│   │   ├── handlers/       # HTTP handlers for APIs
│   │   ├── middleware/     # Custom middleware
│   │   ├── models/         # Data models
│   │   └── router/         # Route setup
│   ├── docs/              # Generated Swagger documentation
│   ├── main.go            # Application entry point
│   ├── go.mod            # Go dependencies
│   └── Dockerfile        # Docker configuration
├── frontend/              # Nuxt.js app
│   ├── components/       # Vue components
│   ├── pages/           # Route pages
│   ├── plugins/         # Nuxt plugins
│   └── Dockerfile
├── scripts/             # Development scripts
├── docker-compose.yml
└── README.md
```

### Data Models

- **User**: Platform users (synchronized with Clerk authentication)
- **Project**: Enterprise projects with members
- **Channel**: Communication channels (public/private)
- **Message**: Channel messages and direct messages
- **File**: Shared files in projects

### Backend Quick Development

For backend-specific development:

```bash
# Install dependencies
cd backend
go mod tidy

# Generate Swagger documentation
go install github.com/swaggo/swag/cmd/swag@latest
swag init

# Start the server (requires database)
go run main.go
```

### Development Tools

The project includes development automation scripts:

```bash
# Setup development environment (one-time)
.\scripts\setup-hooks.ps1

# Development workflow
.\scripts\dev.bat            # Start backend with hot reload
.\scripts\pre-commit.bat     # Run formatting and linting
.\scripts\db-verify.bat      # Verify database schema
```

For complete automation details, see [AUTOMATION.md](AUTOMATION.md).

### Authentication Integration

Thothix uses Clerk for secure authentication. For complete setup and integration guide, see [CLERK_INTEGRATION_COMPLETE.md](CLERK_INTEGRATION_COMPLETE.md).

### Hot Reload

- **Frontend**: Nuxt with automatic hot reload
- **Backend**: Air for Go server auto-restart (via dev.bat)
- **Database**: Persists through Docker volumes

### Testing

```bash
# Go backend tests
docker-compose exec thothix-api go test ./...

# Frontend tests
docker-compose exec thothix-web npm run test

# Integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📡 API Reference

The REST API is documented with Swagger UI and available at: `http://localhost:30000/swagger/index.html`

### Role-Based Access Control (RBAC) System

Thothix implements a simplified role-based access control (RBAC) system:

#### Available Roles

- **Admin**: Can manage the entire system
- **Manager**: Can manage everything except user management
- **User**: Can participate in assigned projects and channels, create 1:1 chats
- **External**: Can only participate in public channels

#### Public/Private Channel Strategy

- **Public Channels**: No explicit members in the `channel_members` table
- **Private Channels**: At least one member in the `channel_members` table

For more details see: [`backend/RBAC_SIMPLIFIED.md`](backend/RBAC_SIMPLIFIED.md)

### Main Endpoints

#### Authentication

- `POST /api/v1/auth/sync` - Sync user with Clerk
- `GET /api/v1/auth/me` - Current user information

#### Projects

- `GET /api/v1/projects` - List projects
- `POST /api/v1/projects` - Create project (Manager/Admin)
- `GET /api/v1/projects/{id}` - Project details

#### Channels

- `GET /api/v1/channels` - List accessible channels
- `POST /api/v1/channels` - Create channel (Manager/Admin)
- `POST /api/v1/channels/{id}/join` - Join public channel

#### Messages

- `GET /api/v1/channels/{id}/messages` - Channel messages
- `POST /api/v1/channels/{id}/messages` - Send message
- `POST /api/v1/messages/direct` - Direct message 1:1

#### Role Management (Admin Only)

- `POST /api/v1/roles` - Assign role
- `DELETE /api/v1/roles/{roleId}` - Revoke role

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

This project is distributed under the MIT License. See `LICENSE` for more information.

---

## 👥 Team

- **Diego** - [@diego81b](https://github.com/diego81b) - Initial development and architecture

## 📚 Roadmap

- [x] Base architecture setup
- [ ] JWT Authentication
- [ ] 1:1 Chat
- [ ] Group chat
- [ ] Project management
- [ ] File upload
- [ ] Real-time notifications
- [ ] Mobile app
- [ ] CI/CD pipeline

To contribute to the project, see the [developer guide](./docs/CONTRIBUTING.md).

---

## 🌍 Multi-Environment Management

Thothix supports multiple deployment environments using a single unified environment file.

### Environment Configuration

```bash
.env                    # Unified configuration for all environments
                       # Copy from .env.example and customize per environment
```

**Environment-specific settings are controlled by variables in your `.env` file:**

- `ENVIRONMENT=development|staging|production`
- `USE_VAULT=false|true` (recommended for staging/production)
- `GIN_MODE=debug|release`
- Different Clerk keys (test vs live)

### Deployment Commands

```bash
# Development (vault services available but only used if USE_VAULT=true)
docker-compose up -d --build

# Production (with integrated Vault - USE_VAULT automatically set to true)
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Using deployment script for any environment
.\scripts\deploy.bat dev up      # Development
.\scripts\deploy.bat prod up     # Production with Vault
```

### Vault Management

```bash
# Initialize Vault (production)
.\scripts\deploy.bat prod vault init

# Open Vault UI
.\scripts\deploy.bat prod vault ui

# Check Vault status
.\scripts\deploy.bat prod vault status
```

### Environment Variables

Each environment file should contain:

- **Database credentials** (different per environment)
- **Clerk API keys** (development vs production)
- **Application settings** (debug mode, logging level)
- **External service URLs** (different endpoints per environment)

### Security Notes

- ⚠️ **Never commit production secrets** to version control
- ✅ Use `.env.local` for local overrides (ignored by git)
- ✅ Store production secrets in secure vault systems
- ✅ Use different database names/passwords per environment

---
