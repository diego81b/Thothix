# 🔐 Vault Integration Guide

Complete secret management for Thothix using HashiCorp Vault.

## ⚡ Quick Start

```bash
cp .env.example .env     # Copy template
# Set USE_VAULT=true
npm run vault:init       # Initialize + sync
npm run dev             # Start services
```

**Vault UI**: <http://localhost:8200>

## � Comment-Based Sync System

Only sections with `:folder_name` comments sync to Vault:

### ✅ Synced Sections

```bash
# :database - Database credentials
DB_USER=postgres
DB_PASSWORD=secret123
DB_NAME=thothix_db

# :clerk - Authentication tokens
CLERK_SECRET_KEY=sk_test_...
CLERK_WEBHOOK_SECRET=whsec_...

# :app - Application secrets
JWT_SECRET=your_jwt_secret
ENCRYPTION_KEY=your_encryption_key
```

### ❌ NOT Synced (No prefix)

```bash
# Application Configuration
PORT=3000
ENVIRONMENT=development
USE_VAULT=true

# Vault Configuration
VAULT_ADDR=http://vault:8200
VAULT_ROOT_TOKEN=your_token
```

**Logic**: `# :folder_name - Description` → `vault/thothix/folder_name` in Vault

## 🎯 Commands

| Command                 | Purpose                 | Usage              |
| ----------------------- | ----------------------- | ------------------ |
| `npm run vault:init`    | First-time setup + sync | Initial setup      |
| `npm run vault:sync`    | Sync secrets only       | After .env changes |
| `npm run vault:cleanup` | Clean temp files        | If files remain    |

### Direct Script Usage

```bash
zx scripts/vault.mjs --init     # Initialize
zx scripts/vault.mjs --sync     # Sync only
zx scripts/vault.mjs --cleanup  # Cleanup
```

## 🚨 Common Issues

### Vault Won't Start

```bash
docker compose logs vault
docker compose restart vault
```

### Secrets Not Syncing

```bash
# Check comment format
grep -n "^#.*:" .env

# Manual sync with output
npm run vault:sync
```

### Token Issues

```bash
# Check token
echo $VAULT_ROOT_TOKEN

# Reinitialize
npm run vault:init
```

### Temp Files Remain

```bash
npm run vault:cleanup  # Auto cleanup
```

## 📋 Vault Structure

```text
thothix/
├── database/    # Database credentials
├── clerk/       # Authentication tokens
└── app/         # Application secrets
```

## 🔧 Configuration Files

- **Environment Template**: [`.env.example`](../.env.example)
- **Vault Config**: [`vault.hcl`](../vault.hcl)
- **Sync Script**: [`scripts/vault.mjs`](../scripts/vault.mjs)
- **Docker Setup**: [`docker-compose.yml`](../docker-compose.yml)

## 🛡️ Security Features

- ✅ Only marked sections sync to Vault
- ✅ Temp files auto-cleaned and Git-ignored
- ✅ Secrets never exposed in logs
- ✅ Vault sealed by default in production

---

**System**: Modern unified vault management v3.0
**Status**: ✅ Fully functional and tested


### Secret Management Strategy

- ✅ **Comment-based sync**: Only `:section` marked secrets go to Vault
- ✅ **Configuration separation**: App configs stay in `.env`, secrets in Vault
- ✅ **Environment flexibility**: Can use Vault or local `.env` per environment
- ✅ **Security by default**: Sensitive data automatically excluded from Git
- ✅ **Easy migration**: Add/remove Vault without changing application code

### Environment Management

- ✅ **One Vault instance per environment** (dev/staging/prod)
- ✅ **Separate mount points** per environment (`thothix-dev`, `thothix-prod`)
- ✅ **Environment-specific policies** and tokens
- ✅ **Backup Vault data** before major changes

### Monitoring

- ✅ **Check Vault health** in monitoring systems
- ✅ **Alert on token expiration** before they expire
- ✅ **Monitor secret access patterns** for anomalies
- ✅ **Review audit logs** regularly

---

## 📚 Additional Resources

### Quick Command Reference

```bash
# Complete Development Workflow
npm run dev              # 🚀 Start all services (auto-init Vault)
npm run dev:logs         # 📋 Monitor all service logs
npm run dev:down         # 🛑 Stop all services

# Vault Operations
npm run vault:init       # 🏗️  Full Vault setup + sync all secrets
npm run vault:sync       # 🔄 Update Vault with current .env secrets

# Database Management
npm run db:status        # ✅ Check database connection health
npm run db:tables        # 📋 List all database tables
npm run db:connect       # 🔗 Open interactive database shell

# Code Quality
npm run format           # 🎨 Format all code (Go, JS, etc.)
npm run lint             # 🔍 Run all linters and checks
npm run pre-commit       # ✅ Run pre-commit validation

# Maintenance
npm run vault:cleanup          # 🧹 Clean temporary files
```

### Links

- 🔗 **HashiCorp Vault Documentation**: <https://developer.hashicorp.com/vault>
- 🔗 **Docker Compose Reference**: [`docker-compose.yml`](../docker-compose.yml)
- 🔗 **Environment Template**: [`.env.example`](../.env.example)
- 🔗 **Vault Configuration**: [`vault.hcl`](../vault.hcl)
- 🔗 **Vault Management Script**: [`scripts/vault.mjs`](../scripts/vault.mjs)

### Related Documentation

- 📖 **[Docker Modernization](./DOCKER_MODERNIZATION.md)** - Container setup
- 📖 **[Clerk Integration](./CLERK_INTEGRATION.md)** - Authentication setup
- 📖 **[Database Migration](./DB_MIGRATION.md)** - Database configuration

### Support

For issues:

1. Check this troubleshooting guide
2. Review container logs: `docker compose logs`
3. Verify `.env` configuration
4. Test with Vault UI at <http://localhost:8200>

---

**Last Updated**: June 18, 2025
**System**: Modern unified vault management v3.0
**Status**: ✅ Fully functional and tested
