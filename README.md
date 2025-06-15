# Thothix Messaging Platform

[![Docker](https://img.shields.io/badge/Docker-Ready-blue?style=flat&logo=docker)](https://docker.com)
[![Go](https://img.shields.io/badge/Go-1.23-blue?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue?style=flat&logo=postgresql)](https://postgresql.org)

Thothix is a modern enterprise messaging platform that enables project management, group chats, and private 1:1 conversations. Built with containerized microservices architecture.

## 🚀 GitHub Repository

```bash
# Clone the repository
git clone https://github.com/diego81b/Thothix.git
cd Thothix
```

## 📋 Table of Contents

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [HashiCorp Vault Integration](#hashicorp-vault-integration)
- [Docker Configuration](#docker-configuration)
- [Useful Docker Commands](#useful-docker-commands)
- [Database Verification Tools](#database-verification-tools)
- [Development](#development)
- [API Reference](#api-reference)
- [Contributing](#contributing)

## 📖 Additional Documentation

- **[⚡ Backend Documentation](./backend/README.md)** - Complete Go API documentation, data models, and development guide
- **[🔐 Vault Integration Guide](./docs/VAULT_INTEGRATION.md)** - Complete setup, troubleshooting & production guide
- **[🐳 Docker Guide](./docs/DOCKER_MODERNIZATION.md)** - Docker architecture updates and migration guide
- **[🌍 Node.js Automation](./docs/NODE_JS_GUIDE.md)** - Cross-platform automation with Node.js/Zx
- **[📝 Changelog](./CHANGELOG.md)** - Project history and version updates

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

> 🚀 **Recent Updates**: Docker configuration has been modernized with multi-stage builds and consistent naming conventions. See [Docker Modernization Guide](./docs/DOCKER_MODERNIZATION.md) for details.

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

## ⚡ Quick Start

### 1. Repository Clone

```bash
git clone <repository-url>
cd Thothix
```

### 2. Install Node.js Dependencies

```bash
# Install Node.js from https://nodejs.org/
npm install
```

### 3. Environment Configuration

```bash
# Copy and configure environment variables
cp .env.example .env
notepad .env  # or your preferred editor
```

**Configure your `.env` file with your credentials:**

📋 **Reference**: See [`.env.example`](./.env.example) for all available configuration options including:

- Database settings
- Clerk authentication keys
- Application configuration
- Vault integration settings

⚠️ **Security Note**: Never commit the `.env` file to version control. It's already included in `.gitignore`.

### 4. Development Commands

```bash
# Format code
npm run format

# Run pre-commit checks
npm run pre-commit

# Start development environment
npm run dev

# Check database
npm run db:status
```

**All commands are natively cross-platform with Node.js!** 🎯

**For Vault Integration (Optional):**

📖 **Full Setup Guide**: See [Vault Integration Guide](./docs/VAULT_INTEGRATION.md) for complete configuration.

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

📖 **[📋 Guida Completa all'Integrazione Vault →](./docs/VAULT_INTEGRATION.md)**

### Quick Start

1. **Abilita Vault**: Imposta `USE_VAULT=true` nel tuo `.env`
2. **Avvia servizi**: `docker compose up -d --build`
3. **Accedi a Vault UI**: <http://localhost:8200> (token: da tuo `.env`)

Vault gestisce automaticamente credenziali database, API keys Clerk e segreti applicazione.

Per setup completo, troubleshooting e configurazione produzione → **[VAULT_INTEGRATION.md](./docs/VAULT_INTEGRATION.md)**

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

📖 **Detailed Guide**: See [Docker Modernization Guide](./docs/DOCKER_MODERNIZATION.md) for migration details and best practices.

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
# Ricostruisci solo l'API
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

### Database Operations

The project includes cross-platform database utilities via npm scripts:

```bash
# Check database status
npm run db:status

# Connect to database interactively
npm run db:connect

# List all tables
npm run db:tables

# Verify BaseModel alignment (all tables should have 5 columns)
npm run db:check

# Advanced operations (direct Zx script usage)
npx zx scripts/db-verify.mjs check-table users
npx zx scripts/db-verify.mjs has-field users email
npx zx scripts/db-verify.mjs missing-field users created_by
```

**Cross-platform**: Same commands work on Windows, Linux, and macOS! 🌍

### Manual SQL Commands

To execute SQL commands directly:

```bash
# Direct database connection
docker-compose exec postgres psql -U postgres -d thothix-db

# Execute single command
docker-compose exec postgres psql -U postgres -d thothix-db -c "SELECT version();"
```

For more details on verification commands, see `docs/DB_MIGRATION.md`.

## 💻 Development

### Project Structure

```
thothix/
├── backend/                 # Go API + Gin (see backend/README.md)
├── frontend/                # Nuxt.js app
├── scripts/                 # Development automation (Node.js/Zx)
├── docker-compose.yml       # Development environment
└── README.md               # This file
```

### Backend Development

**For detailed backend documentation, API reference, data models, and development guide:**

📖 **[Backend Documentation →](./backend/README.md)**

Quick backend commands:
```bash
# Backend development via npm scripts
npm run format      # Format Go code
npm run lint        # Run golangci-lint
npm run test        # Run Go tests
npm run pre-commit  # Complete checks

# Direct backend development
cd backend
go mod tidy
go run main.go
```

### Development Tools

The project uses modern Node.js/Zx automation for all development tasks:

```bash
# Development workflow
npm run format      # Format Go code
npm run lint        # Run golangci-lint
npm run test        # Run Go tests
npm run pre-commit  # Complete pre-commit checks

# Environment management
npm run dev         # Start development environment
npm run dev:down    # Stop development environment
npm run staging     # Deploy to staging
npm run prod        # Deploy to production

# Database operations
npm run db:status   # Check database status
npm run db:check    # Verify BaseModel schema
```

**Modern Automation:**

- Cross-platform Node.js/Zx scripts work on Windows, Linux, macOS
- Zero logic duplication - single codebase for all platforms
- NPM scripts provide familiar, universal interface
- Integrated with VS Code tasks (Ctrl+Shift+P → "Tasks: Run Task")

For complete development details, see [NODE_JS_GUIDE.md](./docs/NODE_JS_GUIDE.md).

### Authentication Integration

Thothix uses Clerk for secure authentication. For complete setup and integration guide, see [CLERK_INTEGRATION.md](./docs/CLERK_INTEGRATION.md).

### Hot Reload

- **Frontend**: Nuxt with automatic hot reload
- **Backend**: Air for Go server auto-restart (via npm scripts)
- **Database**: Persists through Docker volumes

### Testing

```bash
# Go backend tests (via npm)
npm run test

# Go backend tests (via Docker)
docker compose exec thothix-api go test ./...

# Frontend tests
docker compose exec thothix-web npm run test

# Integration tests
docker compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📡 API Reference

**For complete API documentation, endpoints, and examples:**

📖 **[Backend API Documentation →](./backend/README.md#api-reference)**

### Quick API Access

- **Swagger UI**: `http://localhost:30000/swagger/index.html`
- **API Base**: `http://localhost:30000/api/v1`
- **Health Check**: `http://localhost:30000/health`

### Authentication & RBAC

Thothix uses **Clerk** for authentication and implements a simplified **Role-Based Access Control (RBAC)** system.

**For complete authentication and RBAC documentation:**
- 📖 **[Backend RBAC Guide →](./backend/RBAC_SIMPLIFIED.md)**
- 📖 **[Clerk Integration Guide →](./docs/CLERK_INTEGRATION.md)**

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
docker compose up -d --build

# Production (with integrated Vault - USE_VAULT automatically set to true)
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Using npm scripts for environment management
npm run dev         # Development
npm run staging     # Staging
npm run prod        # Production with Vault
```

### Vault Management

```bash
# Initialize Vault (production)
npx zx scripts/deploy.mjs prod vault init

# Open Vault UI
npx zx scripts/deploy.mjs prod vault ui

# Check Vault status
npx zx scripts/deploy.mjs prod vault status
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
