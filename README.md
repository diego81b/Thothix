# Thothix - Enterprise Messaging Platform

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
- [Docker Configuration](#docker-configuration)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Useful Docker Commands](#useful-docker-commands)
- [Development](#development)
- [API Reference](#api-reference)
- [Troubleshooting](#troubleshooting)

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
                       ┌─────────────────┐
                       │  thothix-minio  │
                       │ (File Storage)  │
                       │   Port: 30002   │
                       │ Console: 30003  │
                       └─────────────────┘
```

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
- [Kubernetes](https://kubernetes.io/docs/setup/) (for deployment)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) (for K8s management)
- Git

### Installation Verification

```bash
docker --version
docker-compose --version
kubectl version --client
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

### 3. Start the complete stack

```bash
# Start all services
docker-compose up -d --build

# Verify all containers are running
docker-compose ps
```

### 4. Database initialization

```bash
# Run database migrations
docker-compose exec thothix-api go run cmd/migrate/main.go

# Load sample data (optional)
docker-compose exec thothix-api go run cmd/seed/main.go
```

### 5. Access services

| Service | URL | Credentials |
|----------|-----|-------------|
| **API Swagger** | <http://localhost:30000/swagger/index.html> | - |
| **Thothix Web** | <http://localhost:30001> | - |
| **MinIO Console** | <http://localhost:30002> | admin/password123 |
| **pgAdmin** | <http://localhost:5432> | postgres/@Admin123 |

## ⚙️ Docker Configuration

### docker-compose.yml

```yaml
services:
  postgres:
    build:
      context: .
      dockerfile: Dockerfile.postgres
    image: thothix/postgres:17.5-thothix1.0
    container_name: thothix-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: '@Admin123'
      POSTGRES_DB: thothix-db
      POSTGRES_HOST_AUTH_METHOD: trust
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d thothix-db"]
      interval: 10s
      timeout: 5s
      retries: 5

  thothix-api:
    build:
      context: ./backend
      dockerfile: Dockerfile
    image: thothix/api:1.0.0
    container_name: thothix-api
    restart: unless-stopped
    ports:
      - "30000:30000"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: '@Admin123'
      DB_NAME: thothix-db
      PORT: 30000
      CLERK_SECRET_KEY: ${CLERK_SECRET_KEY}
      ENVIRONMENT: development
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - app-network

volumes:
  postgres_data:
    driver: local

networks:
  app-network:
    driver: bridge
```

### Dockerfile per PostgreSQL (Dockerfile.postgres)

```dockerfile
FROM postgres:17.5-alpine

# Metadata immagine
LABEL org.opencontainers.image.title="Thothix PostgreSQL Database"
LABEL org.opencontainers.image.description="Database PostgreSQL personalizzato per Thothix"
LABEL org.opencontainers.image.version="17.5-thothix1.0"
LABEL org.opencontainers.image.vendor="Thothix"

# Estensioni PostgreSQL
RUN apk add --no-cache postgresql-contrib

# Script di init
COPY db-init/ /docker-entrypoint-initdb.d/

EXPOSE 5432
```

### Dockerfile per API Go (Dockerfile)

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app/backend
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM registry.access.redhat.com/ubi9/ubi-micro
WORKDIR /app

# Metadata immagine
LABEL org.opencontainers.image.title="Thothix API"
LABEL org.opencontainers.image.description="API backend per la piattaforma Thothix"
LABEL org.opencontainers.image.version="1.0.0"

COPY --from=builder /app/backend/main .
EXPOSE 30000
ENTRYPOINT ["./main"]
```

### Variabili d'ambiente (.env)

```bash
# Applicazione
APP_NAME=thothix
APP_ENV=development
APP_URL=http://localhost:30001

# Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_NAME=thothix-db
DB_USER=postgres
DB_PASSWORD=@Admin123

# MinIO Storage
MINIO_ENDPOINT=thothix-minio:30002
MINIO_ACCESS_KEY=admin
MINIO_SECRET_KEY=@Admin123
MINIO_BUCKET=thothix-files
MINIO_USE_SSL=false

# API Backend
API_PORT=30000
API_HOST=localhost/api

# Sviluppo
DEBUG=true
LOG_LEVEL=debug
GIN_MODE=debug
```

## 🛠️ Comandi Docker Utili

### Comandi utili

```bash
# Avvia tutti i servizi in background
docker-compose up -d

# Ferma e rimuove container, network e volumi definiti
docker-compose down

# Build entrambe le immagini
docker-compose build

# Lista immagini Thothix
docker images | findstr thothix

# Push al registry (quando configurato)
docker push thothix/api:1.0.0
docker push thothix/postgres:17.5-thothix1.0

# Inspect metadata
docker inspect thothix/api:1.0.0 --format="{{json .Config.Labels}}"
docker inspect thothix/postgres:17.5-thothix1.0 --format="{{json .Config.Labels}}"
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
# Genera (o rigenera) la doc. OpenAPI dai commenti @swagger
cd backend
swag init -g main.go

# Riavvia solo il servizio API (con doc aggiornata)
docker-compose restart thothix-api
```

### Logging e debug

```bash
# Segui i log del servizio API
docker-compose logs -f thothix-api

# Segui i log del servizio Postgres
docker-compose logs -f thothix-postgres

# Esegui una shell interattiva nel container API
docker-compose exec thothix-api sh
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
chmod +x scripts/db-verify.sh
./scripts/db-verify.sh check-basemodel
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
│   ├── cmd/                # Entry points
│   ├── internal/           # Business logic
│   ├── pkg/               # Shared packages
│   ├── migrations/        # Database migrations
│   └── Dockerfile
├── frontend/              # Nuxt.js app
│   ├── components/       # Vue components
│   ├── pages/           # Route pages
│   ├── plugins/         # Nuxt plugins
│   └── Dockerfile
├── k8s/                 # Kubernetes manifests
├── docker-compose.yml
└── README.md
```

### Hot Reload

- **Frontend**: Nuxt with automatic hot reload
- **Backend**: Air for Go server auto-restart
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

#### Available Roles:
- **Admin**: Can manage the entire system
- **Manager**: Can manage everything except user management  
- **User**: Can participate in assigned projects and channels, create 1:1 chats
- **External**: Can only participate in public channels

#### Public/Private Channel Strategy:
- **Public Channels**: No explicit members in the `channel_members` table
- **Private Channels**: At least one member in the `channel_members` table

For more details see: [`backend/RBAC_SIMPLIFIED.md`](backend/RBAC_SIMPLIFIED.md)

### Main Endpoints:

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
