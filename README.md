# Thothix - Messaggistica Aziendale

Thothix è una piattaforma di messaggistica aziendale moderna che permette la gestione di progetti, chat di gruppo e conversazioni private 1:1. Costruita con architettura a microservizi containerizzata.

## 📋 Indice

- [Architettura](#architettura)
- [Prerequisiti](#prerequisiti)
- [Avvio Rapido con Docker](#avvio-rapido-con-docker)
- [Configurazione Docker](#configurazione-docker)
- [Deployment Kubernetes](#deployment-kubernetes)
- [Comandi Docker Utili](#comandi-docker-utili)
- [Sviluppo](#sviluppo)
- [API Reference](#api-reference)
- [Troubleshooting](#troubleshooting)

## 🏗️ Architettura

### Stack Tecnologico

- **Frontend**: Vue.js 3 + Nuxt.js + Nuxt UI
- **Backend**: Go + Gin Framework + GORM
- **Database**: PostgreSQL 15
- **File Storage**: MinIO (S3-compatible)
- **Orchestrazione**: Kubernetes
- **Containerizzazione**: Docker

### Struttura dei Servizi

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

### Funzionalità Principali

- 💼 **Gestione Progetti**: Creazione e amministrazione progetti aziendali
- 👥 **Chat di Gruppo**: Canali pubblici e privati per progetti
- 💬 **Messaggi 1:1**: Conversazioni private tra utenti
- 📁 **Condivisione File**: Upload e condivisione documenti tramite MinIO
- 🔐 **Gestione Utenti**: Sistema di autenticazione e autorizzazione
- 📱 **Responsive**: Interfaccia ottimizzata per desktop e mobile

## 🚀 Prerequisiti

- [Docker](https://docs.docker.com/get-docker/) (versione 20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (versione 2.0+)
- [Kubernetes](https://kubernetes.io/docs/setup/) (per deployment)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) (per gestione K8s)
- Git

### Verifica installazione

```bash
docker --version
docker-compose --version
kubectl version --client
```

## 🐳 Avvio Rapido con Docker

### 1. Clone del repository

```bash
git clone <repository-url>
cd Thothix
```

### 2. Configurazione ambiente

```bash
# Copia e configura le variabili d'ambiente
cp .env.example .env
notepad .env
```

### 3. Avvio completo dell'stack

```bash
# Avvia tutti i servizi
docker-compose up -d --build

# Verifica che tutti i container siano attivi
docker-compose ps
```

### 4. Inizializzazione database

```bash
# Esegui le migrazioni del database
docker-compose exec thothix-api go run cmd/migrate/main.go

# Carica dati di esempio (opzionale)
docker-compose exec thothix-api go run cmd/seed/main.go
```

### 5. Accesso ai servizi

| Servizio | URL | Credenziali |
|----------|-----|-------------|
| **API Swagger** | <http://localhost:30000/swagger/index.html> | - |
| **Thothix Web** | <http://localhost:30001> | - |
| **MinIO Console** | <http://localhost:30002> | admin/password123 |
| **pgAdmin** | <http://localhost:5432> | postgres/@Admin123 |

## ⚙️ Configurazione Docker

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

## 💻 Sviluppo

### Struttura del progetto

```
thothix/
├── backend/                 # API Go + Gin
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

### Hot reload

- **Frontend**: Nuxt con hot reload automatico
- **Backend**: Air per auto-restart del server Go
- **Database**: Persiste attraverso volumi Docker

### Testing

```bash
# Test Go backend
docker-compose exec thothix-api go test ./...

# Test frontend
docker-compose exec thothix-web npm run test

# Test di integrazione
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📡 API Reference

TODO
---

## 📚 Roadmap

- [x] Setup base architettura
- [ ] Autenticazione JWT
- [ ] Chat 1:1
- [ ] Chat di gruppo
- [ ] Gestione progetti
- [ ] Upload file
- [ ] Notifiche real-time
- [ ] Mobile app
- [ ] CI/CD pipeline

Per contribuire al progetto, consulta la [guida per sviluppatori](./docs/CONTRIBUTING.md).
