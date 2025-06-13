# Thothix - Messaggistica Aziendale

[![GitHub Repository](https://img.shields.io/badge/GitHub-diego81b%2FThothix-blue?style=flat&logo=github)](https://github.com/diego81b/Thothix)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue?style=flat&logo=docker)](https://docker.com)
[![Go](https://img.shields.io/badge/Go-1.23-blue?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue?style=flat&logo=postgresql)](https://postgresql.org)

Thothix è una piattaforma di messaggistica aziendale moderna che permette la gestione di progetti, chat di gruppo e conversazioni private 1:1. Costruita con architettura a microservizi containerizzata.

## 🚀 Repository GitHub

Il progetto è disponibile su GitHub: **[https://github.com/diego81b/Thothix](https://github.com/diego81b/Thothix)**

```bash
# Clona il repository
git clone https://github.com/diego81b/Thothix.git
cd Thothix
```

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

## 🔍 Database Verification Tools

### Utility Scripts

Il progetto include script per verificare facilmente l'allineamento del database con i modelli Go:

#### Windows (cmd/PowerShell)

```cmd
# Verifica allineamento BaseModel (tutte le tabelle dovrebbero avere 5 colonne)
scripts\db-verify.bat check-basemodel

# Lista tutte le tabelle
scripts\db-verify.bat list-tables

# Controlla struttura di una tabella specifica
scripts\db-verify.bat check-table users

# Trova tabelle che mancano di un campo specifico
scripts\db-verify.bat missing-field updated_by

# Trova tabelle che hanno un campo specifico
scripts\db-verify.bat has-field system_role

# Connetti al database interattivamente
scripts\db-verify.bat connect

# Stato del database
scripts\db-verify.bat status
```

#### Linux/MacOS (bash)

```bash
# Stesso utilizzo ma con estensione .sh
chmod +x scripts/db-verify.sh
./scripts/db-verify.sh check-basemodel
```

### Comandi SQL Manuali

Per eseguire comandi SQL direttamente:

```bash
# Connessione diretta al database
docker-compose exec postgres psql -U postgres -d thothix-db

# Esecuzione comando singolo
docker-compose exec postgres psql -U postgres -d thothix-db -c "SELECT version();"
```

Per più dettagli sui comandi di verifica, consulta `DB_MIGRATION.md`.

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

L'API REST è documentata con Swagger UI e disponibile all'indirizzo: `http://localhost:30000/swagger/index.html`

### Sistema di Ruoli e Permessi (RBAC)

Thothix implementa un sistema di controllo accessi basato su ruoli (RBAC) semplificato:

#### Ruoli Disponibili:
- **Admin**: Può gestire tutto il sistema
- **Manager**: Può gestire tutto tranne la gestione degli utenti  
- **User**: Può partecipare ai progetti e canali assegnati, creare chat 1:1
- **External**: Può solo partecipare ai canali pubblici

#### Strategia Canali Pubblici/Privati:
- **Canali Pubblici**: Nessun membro esplicito nella tabella `channel_members`
- **Canali Privati**: Almeno un membro nella tabella `channel_members`

Per maggiori dettagli consultare: [`backend/RBAC_SIMPLIFIED.md`](backend/RBAC_SIMPLIFIED.md)

### Endpoint Principali:

#### Autenticazione
- `POST /api/v1/auth/sync` - Sincronizza utente con Clerk
- `GET /api/v1/auth/me` - Informazioni utente corrente

#### Progetti
- `GET /api/v1/projects` - Lista progetti
- `POST /api/v1/projects` - Crea progetto (Manager/Admin)
- `GET /api/v1/projects/{id}` - Dettagli progetto

#### Canali
- `GET /api/v1/channels` - Lista canali accessibili
- `POST /api/v1/channels` - Crea canale (Manager/Admin)
- `POST /api/v1/channels/{id}/join` - Unisciti a canale pubblico

#### Messaggi
- `GET /api/v1/channels/{id}/messages` - Messaggi del canale
- `POST /api/v1/channels/{id}/messages` - Invia messaggio
- `POST /api/v1/messages/direct` - Messaggio diretto 1:1

#### Gestione Ruoli (Solo Admin)
- `POST /api/v1/roles` - Assegna ruolo
- `DELETE /api/v1/roles/{roleId}` - Revoca ruolo

---

## 🤝 Contribuire

1. Fai fork del repository
2. Crea un branch per la tua feature (`git checkout -b feature/AmazingFeature`)
3. Committa le modifiche (`git commit -m 'Add some AmazingFeature'`)
4. Push al branch (`git push origin feature/AmazingFeature`)
5. Apri una Pull Request

---

## 📄 Licenza

Questo progetto è distribuito sotto licenza MIT. Vedi `LICENSE` per maggiori informazioni.

---

## 👥 Team

- **Diego** - [@diego81b](https://github.com/diego81b) - Sviluppo iniziale e architettura

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
