# 🔐 Guida Completa all'Integrazione di HashiCorp Vault

Questa guida ti accompagna passo dopo passo nell'integrazione di HashiCorp Vault per la gestione sicura dei segreti in Thothix.

## 📋 Indice

1. [Quick Start](#quick-start) ⚡
2. [Cos'è Vault e Perché Usarlo](#cosa-è-vault)
3. [Setup Iniziale](#setup-iniziale)
4. [Configurazione Sviluppo](#configurazione-sviluppo)
5. [Configurazione Produzione](#configurazione-produzione)
6. [Gestione dei Segreti](#gestione-dei-segreti)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)

## ⚡ Quick Start

**Vuoi iniziare subito? Segui questi 4 passi:**

### 1. Abilita Vault nel .env
```bash
cp .env.example .env
# Modifica .env e imposta:
USE_VAULT=true
```

### 2. Avvia i servizi
```bash
docker-compose up -d --build
```

### 3. Accedi a Vault UI
- **URL**: http://localhost:8200
- **Token**: `thothix-dev-root-token` (dal tuo .env)

### 4. Verifica l'integrazione
```bash
# Controlla che Vault sia attivo
docker-compose logs vault-init

# Verifica che l'app legga da Vault
docker-compose logs thothix-api | findstr vault
```

**✅ Fatto!** Vault è ora integrato e gestisce automaticamente:
- **Database**: Credenziali PostgreSQL
- **Clerk**: API keys e webhook secrets  
- **App**: JWT secrets e encryption keys

📖 **Per configurazione avanzata, troubleshooting e produzione continua a leggere...**

---

## 🎯 Cos'è Vault

HashiCorp Vault è un tool per:
- **Gestione sicura dei segreti** (password, API keys, certificati)
- **Crittografia** dei dati sensibili
- **Controllo accessi** granulare
- **Audit logging** completo
- **Rotazione automatica** delle credenziali

### Vantaggi per Thothix:
- ✅ **Nessun segreto hardcoded** nel codice
- ✅ **Centralizzazione** di tutte le credenziali
- ✅ **Sicurezza** enterprise-grade
- ✅ **Audit trail** completo
- ✅ **Separazione** dev/staging/prod

## 🚀 Setup Iniziale

### 1. Prerequisiti

Assicurati di avere:
- Docker e Docker Compose installati
- File `.env` configurato (copia da `.env.example`)
- Porte 8200 (Vault) e 30000 (API) disponibili

### 2. Configurazione Base

```bash
# 1. Copia il template di configurazione
cp .env.example .env

# 2. Modifica il file .env
notepad .env
```

### 3. Configurazione Minimale nel .env

```bash
# Abilita Vault
USE_VAULT=true

# Configurazione Vault per sviluppo
VAULT_ADDR=http://vault:8200
VAULT_ROOT_TOKEN=thothix-dev-root-token
VAULT_APP_TOKEN=will-be-set-after-init
VAULT_MOUNT=thothix
VAULT_DEV_MODE=true

# I tuoi segreti attuali (verranno migrati in Vault)
POSTGRES_PASSWORD=your_secure_password
CLERK_SECRET_KEY=sk_test_your_clerk_key
CLERK_WEBHOOK_SECRET=whsec_your_webhook_secret
DB_PASSWORD=your_secure_password
```

## 💻 Configurazione Sviluppo

### 1. Avvio Completo

```bash
# Avvia tutti i servizi incluso Vault
docker-compose up -d --build

# Verifica che tutti i container siano attivi
docker-compose ps
```

### 2. Verifica Vault

```bash
# Controlla che Vault sia healthy
docker-compose exec vault vault status

# Dovrebbe mostrare:
# - Sealed: false
# - Cluster Mode: standalone
# - Version: 1.15.0
```

### 3. Accesso Vault UI

1. Apri browser su: `http://localhost:8200`
2. Seleziona "Token" come metodo di login
3. Usa il token dal tuo `.env`: `thothix-dev-root-token`
4. Dovresti vedere la dashboard di Vault

### 4. Verifica Inizializzazione

Lo script di init dovrebbe aver creato:

```
Secrets Engines:
├── thothix/ (KV Version 2)
    ├── data/database/
    ├── data/clerk/
    └── data/app/
```

## 🏭 Configurazione Produzione

### 1. Differenze Produzione vs Sviluppo

| Aspetto | Sviluppo | Produzione |
|---------|----------|------------|
| **Vault Mode** | Dev Mode (in-memory) | Production (persistent) |
| **TLS** | Disabilitato | Abilitato con certificati |
| **Autenticazione** | Root token | Policy-based tokens |
| **Storage** | Locale/temporaneo | Volumi persistenti |
| **Backup** | Non necessario | Schedulato |

### 2. Configurazione .env Produzione

```bash
# Production settings
USE_VAULT=true
ENVIRONMENT=production
GIN_MODE=release

# Vault production
VAULT_ADDR=http://vault:8200
VAULT_ROOT_TOKEN=your-secure-production-root-token
VAULT_APP_TOKEN=your-production-app-token

# Database production
POSTGRES_PASSWORD=very-secure-production-password
POSTGRES_DB=thothix-prod
```

### 3. Avvio Produzione

```bash
# Usa la configurazione produzione
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Verifica che Vault sia in modalità produzione
docker-compose exec vault vault status
```

### 4. Setup Sicurezza Produzione

```bash
# 1. Cambia il root token (dopo il primo avvio)
docker-compose exec vault vault auth -method=token

# 2. Crea policy dedicata per l'app
docker-compose exec vault vault policy write thothix-app - <<EOF
path "thothix/data/*" {
  capabilities = ["read"]
}
EOF

# 3. Genera token con policy limitata
docker-compose exec vault vault token create -policy=thothix-app
# Usa questo token come VAULT_APP_TOKEN
```

## 🔧 Gestione dei Segreti

### 1. Struttura dei Segreti

Vault organizza i segreti in percorsi logici:

```
thothix/
├── data/database/
│   ├── host=postgres
│   ├── port=5432
│   ├── username=postgres
│   ├── password=secure_password
│   └── database=thothix-db
│
├── data/clerk/
│   ├── secret_key=sk_test_...
│   ├── webhook_secret=whsec_...
│   └── publishable_key=pk_test_...
│
└── data/app/
    ├── jwt_secret=auto_generated_32_chars
    ├── encryption_key=auto_generated_32_chars
    └── environment=development
```

### 2. Lettura Segreti

```bash
# Lista tutti i path dei segreti
docker-compose exec vault vault kv list thothix/data

# Leggi segreti database
docker-compose exec vault vault kv get thothix/data/database

# Leggi segreti Clerk
docker-compose exec vault vault kv get thothix/data/clerk

# Leggi segreti app
docker-compose exec vault vault kv get thothix/data/app
```

### 3. Modifica Segreti

```bash
# Aggiorna password database
docker-compose exec vault vault kv put thothix/data/database \\
  host=postgres \\
  port=5432 \\
  username=postgres \\
  password=new_secure_password \\
  database=thothix-db

# Aggiorna chiavi Clerk
docker-compose exec vault vault kv put thothix/data/clerk \\
  secret_key=sk_live_new_production_key \\
  webhook_secret=whsec_new_webhook_secret \\
  publishable_key=pk_live_new_public_key

# Riavvia l'API per ricaricare i segreti
docker-compose restart thothix-api
```

### 4. Backup e Restore

```bash
# Backup completo (JSON)
docker-compose exec vault vault kv get -format=json thothix/data > vault_backup.json

# Backup specifico
docker-compose exec vault vault kv get -field=password thothix/data/database

# Restore da backup (manuale)
# Editare vault_backup.json e fare put per ogni segreto
```

## 🛠️ Troubleshooting

### 1. Vault Non Si Avvia

```bash
# Controlla logs
docker-compose logs vault

# Errori comuni:
# - Porta 8200 già in uso
# - Permessi su volumi Docker
# - Configurazione malformata
```

### 2. App Non Si Connette a Vault

```bash
# Verifica connettività
docker-compose exec thothix-api curl -s http://vault:8200/v1/sys/health

# Controlla logs app
docker-compose logs thothix-api | grep -i vault

# Verifica token
docker-compose exec vault vault token lookup $VAULT_APP_TOKEN
```

### 3. Segreti Non Trovati

```bash
# Verifica che esistano
docker-compose exec vault vault kv list thothix/data

# Controlla path esatto
docker-compose exec vault vault kv get thothix/data/database

# Re-inizializza se necessario (ATTENZIONE: cancella tutto)
docker-compose down
docker volume rm thothix_vault_dev_data
docker-compose up -d --build
```

### 4. Reset Completo

```bash
# ATTENZIONE: Cancella tutti i dati di Vault!

# 1. Ferma tutto
docker-compose down

# 2. Rimuovi volumi Vault
docker volume rm thothix_vault_dev_data thothix_vault_dev_logs

# 3. Riavvia
docker-compose up -d --build

# 4. Verifica reinizializzazione
docker-compose logs vault-init
```

## 📚 Best Practices

### 1. Sicurezza

- ✅ **Mai committare** token di produzione nel codice
- ✅ **Ruota i token** regolarmente
- ✅ **Usa policy** granulari per limitare accessi
- ✅ **Abilita TLS** in produzione
- ✅ **Backup** regolari dei segreti

### 2. Organizzazione

- ✅ **Separa** segreti per ambiente (dev/staging/prod)
- ✅ **Usa naming** consistente per i path
- ✅ **Documenta** ogni segreto e il suo scopo
- ✅ **Versiona** le modifiche ai segreti

### 3. Operazioni

- ✅ **Monitora** lo stato di Vault
- ✅ **Log** tutte le operazioni sui segreti
- ✅ **Testa** il failover su .env se Vault è down
- ✅ **Automatizza** il deployment di Vault

### 4. Sviluppo

- ✅ **Usa dev mode** solo per sviluppo locale
- ✅ **Testa** sempre con Vault abilitato prima del deploy
- ✅ **Fallback** graceful su file .env
- ✅ **Cache** i segreti in modo sicuro

## 🔄 Workflow Completo

### Sviluppo Locale
1. `cp .env.example .env`
2. Configura segreti base nel .env
3. Imposta `USE_VAULT=true`
4. `docker-compose up -d --build`
5. Accedi a Vault UI per gestire segreti

### Staging/Produzione
1. Configura .env per l'ambiente target
2. Genera token dedicati con policy limitate
3. Usa `docker-compose.prod.yml` per produzione
4. Backup regolari dei segreti
5. Monitor dello stato di Vault

---

## 📞 Supporto

Per problemi o domande sull'integrazione Vault:

1. **Controlla** questa guida e la sezione troubleshooting
2. **Verifica** i logs con `docker-compose logs`
3. **Consulta** la documentazione ufficiale di [Vault](https://developer.hashicorp.com/vault/docs)
4. **Apri** un issue nel repository del progetto

---

**Buona integrazione con Vault! 🔐**
