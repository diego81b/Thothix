# 🐳 Docker Modernization Guide

Questa guida documenta le modifiche apportate alla configurazione Docker di Thothix per uniformare e modernizzare l'architettura containerizzata.

## 📋 Indice

1. [Modifiche Apportate](#modifiche-apportate)
2. [Nuova Struttura](#nuova-struttura)
3. [Convenzioni di Naming](#convenzioni-di-naming)
4. [Dockerfile Multi-Stage](#dockerfile-multi-stage)
5. [Comandi Aggiornati](#comandi-aggiornati)
6. [Migration Guide](#migration-guide)

## 🔄 Modifiche Apportate

### 1. Uniformazione Dockerfile

**Prima:**

- `Dockerfile` (generico per backend)
- `Dockerfile.postgres`
- Vault gestito via immagine base

**Dopo:**

- `Dockerfile.backend` (backend specifico)
- `Dockerfile.postgres` (aggiornato multi-stage)
- `Dockerfile.vault` (nuovo, personalizzato)

### 2. Convenzioni di Naming

**Development Mode:**

- Immagini: `thothix/service:version-dev`
- Container: `thothix-service-dev`

**Production Mode:**

- Immagini: `thothix/service:version-prod`
- Container: `thothix-service-prod`

### 3. Multi-Stage Build

Tutti i Dockerfile ora supportano:

- **Target `dev`**: Ottimizzato per sviluppo
- **Target `prod`**: Ottimizzato per produzione

## 🏗️ Nuova Struttura

```
📁 Thothix/
├── 🐳 Dockerfile.backend (multi-stage: dev/prod)
├── 🐳 Dockerfile.postgres (multi-stage: dev/prod)
├── 🐳 Dockerfile.vault (multi-stage: dev/prod, con strumenti)
├── 📋 docker-compose.yml (dev: -dev suffixes, dev targets)
├── 📋 docker-compose.prod.yml (prod: -prod suffixes, prod targets)
└── 📁 backend/
    └── 🐳 Dockerfile.backend (alternativo)
```

### Dockerfile.vault Features

- **Base**: HashiCorp Vault 1.15.0
- **Development**: curl, jq, openssl, bash, wget
- **Production**: curl, jq, wget (essenziali)
- **Directory**: `/vault/scripts`, `/vault/config`
- **Permissions**: Corrette per utente `vault`

## 🏷️ Convenzioni di Naming

### Immagini Docker

| Servizio    | Development                            | Production                              |
| ----------- | -------------------------------------- | --------------------------------------- |
| Backend API | `thothix/api:1.0.0-dev`                | `thothix/api:1.0.0-prod`                |
| PostgreSQL  | `thothix/postgres:17.5-thothix1.0-dev` | `thothix/postgres:17.5-thothix1.0-prod` |
| Vault       | `thothix/vault:1.15.0-thothix1.0-dev`  | `thothix/vault:1.15.0-thothix1.0-prod`  |

### Container Names

| Servizio    | Development            | Production              |
| ----------- | ---------------------- | ----------------------- |
| Backend API | `thothix-api-dev`      | `thothix-api-prod`      |
| PostgreSQL  | `thothix-postgres-dev` | `thothix-postgres-prod` |
| Vault       | `thothix-vault-dev`    | `thothix-vault-prod`    |

## 🐳 Dockerfile Multi-Stage

All Dockerfiles now implement multi-stage builds with dedicated dev/prod targets.

📋 **Reference Files:**

- [`Dockerfile.backend`](./Dockerfile.backend) - API backend
- [`Dockerfile.postgres`](./Dockerfile.postgres) - PostgreSQL database
- [`Dockerfile.vault`](./Dockerfile.vault) - HashiCorp Vault

**Features:**

- Optimized builds for each environment
- Consistent labeling and metadata
- Target-specific configurations

## 📝 Comandi Aggiornati

### Development

```bash
# Avvia sviluppo
docker compose up -d

# Build sviluppo
docker compose build

# Logs sviluppo
docker compose logs -f thothix-api
```

### Production

```bash
# Avvia produzione
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Build produzione
docker compose -f docker-compose.yml -f docker-compose.prod.yml build

# Logs produzione
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f thothix-api
```

### Comandi Vault

```bash
# Status Vault
docker compose exec vault vault status

# Logs Vault
docker compose logs vault

# Shell Vault
docker compose exec vault sh
```

## 🔄 Migration Guide

### Per Sviluppatori Esistenti

1. **Ferma i servizi precedenti:**

   ```bash
   docker compose down
   docker system prune -f
   ```

2. **Rebuilda con nuova configurazione:**

   ```bash
   docker compose build
   docker compose up -d
   ```

3. **Verifica i nuovi nomi:**

   ```bash
   docker compose ps
   # Dovresti vedere: thothix-api-dev, thothix-postgres-dev, thothix-vault-dev
   ```

### Variabili d'Ambiente

📋 **Reference**: See [`.env.example`](./.env.example) for all configuration options.

**Key updates needed:**

- Vault connection settings
- Environment-specific configurations
- Service discovery settings

## ✅ Benefits

### 1. **Consistency**

- Naming uniforme tra dev/prod
- Struttura standardizzata
- Dockerfile coerenti

### 2. **Maintainability**

- Multi-stage builds
- Target specifici per ambiente
- Labels OCI standard

### 3. **Security**

- Immagini ottimizzate per produzione
- Strumenti di debug solo in dev
- Permissions corrette

### 4. **DevOps Ready**

- CI/CD friendly
- Registry ready
- Environment separation

## 🚨 Breaking Changes

### Docker Compose

- **Nomi container**: Ora includono suffisso `-dev`/`-prod`
- **Nomi immagini**: Ora includono suffisso `-dev`/`-prod`
- **Build targets**: Specificati esplicitamente

### Comando Updates

- `docker-compose` → `docker compose` (v2 syntax)
- Nuovi target di build
- Nuovi nomi container/immagine

## 🎯 Next Steps

1. **Update CI/CD**: Aggiorna pipeline per nuovi nomi
2. **Update Docs**: Aggiorna tutta la documentazione
3. **Update Scripts**: Aggiorna script di automazione
4. **Registry Push**: Push nuove immagini al registry

---

**✅ Docker Modernization Complete!**

La nuova configurazione è più robusta, manutenibile e pronta per la produzione.
