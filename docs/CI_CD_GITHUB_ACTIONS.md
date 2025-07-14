# GitHub Actions & CI/CD Documentation

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [GitHub Secrets Configuration](#github-secrets-configuration)
- [Available Workflows](#available-workflows)
- [CI/CD Strategy](#cicd-strategy)
- [Local Development](#local-development)
- [Troubleshooting](#troubleshooting)
- [Vault Integration](#vault-integration)
- [Best Practices](#best-practices)

## Overview

This project uses GitHub Actions for CI/CD with a focus on:
- Fast feedback through optimized workflows
- Flexible secret management (Vault + GitHub Secrets)
- Environment alignment between local and CI
- Simplified maintenance and debugging

## Prerequisites

### Local Development Environment
- **Node.js**: v22.13.0+ (matches CI environment)
- **Go**: v1.24.3+ (matches CI environment)
- **npm**: Latest version compatible with Node 22

### Environment Alignment
All CI/CD workflows are configured to match your local environment:
- All workflows use Node.js v22
- All workflows use Go v1.24.3
- Package.json requires Node >=22.0.0

## GitHub Secrets Configuration

### Vault Configuration (Recommended for Production)

If you have a Vault instance deployed in the cloud:

- `VAULT_ADDR`: The URL of your Vault server (e.g., `https://vault.yourcompany.com`)
- `VAULT_TOKEN`: A valid Vault token with read access to the `thothix` mount

### Fallback Environment Variables

If Vault is not available, configure these GitHub Secrets:

- `DB_USER`: Database username (default: `postgres`)
- `DB_PASSWORD`: Database password (default: `postgres`)
- `DB_NAME`: Database name (default: `thothix_test`)
- `DB_HOST`: Database host (default: `localhost`)
- `DB_PORT`: Database port (default: `5432`)
- `CLERK_SECRET_KEY`: Clerk authentication secret key
- `CLERK_WEBHOOK_SECRET`: Clerk webhook secret

### How to Configure GitHub Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret with its corresponding value

## Available Workflows

### 1. **`basic-checks.yml`** - Fast Feedback ⚡
**Trigger**: Push/PR to `main` or `dev`  
**Duration**: 3-5 minutes  
**Purpose**: Quick validation for immediate feedback

**Features**:
- Go tests with testcontainers (PostgreSQL)
- Go linting with golangci-lint
- npm scripts validation
- Vault connectivity check
- Automatic fallback to GitHub Secrets

**Jobs**:
- `go-tests`: Tests with Vault or fallback secrets
- `go-lint`: Code quality validation
- `scripts-validation`: npm scripts verification

### 2. **`ci.yml`** - Complete Pipeline 🔄
**Trigger**: Push/PR to `main` or `dev`  
**Duration**: 10-15 minutes  
**Purpose**: Comprehensive testing and building

**Features**:
- Complete backend testing and building
- Docker image creation
- Security scanning
- Docker Compose integration tests

### 3. **`auto-version-bump.yml`** - Versioning 📦
**Trigger**: PR merged to `main`  
**Purpose**: Automated version management

**Features**:
- Bumps package.json version
- Creates PR with version update
- Extracts feature authors from git log

### 4. **`debug.yml`** - Troubleshooting 🔧
**Trigger**: Manual dispatch  
**Purpose**: CI debugging and environment inspection

**Features**:
- Three debug levels: basic, detailed, full
- Environment variable inspection
- Dependency validation
- Workflow troubleshooting

## CI/CD Strategy

### ✅ **Optimized Approach**

```
Local Development:
npm run pre-commit → Push → basic-checks.yml (fast) → ci.yml (complete)
                                     ↓
                            If problems → debug.yml (manual)
```

### **What We Keep (Valuable)**

1. **`basic-checks.yml`** - Fast Feedback
   - Quick validation (3-5 min vs 15 min)
   - Immediate developer feedback
   - Separate Go tests, lint, npm validation

2. **`debug.yml`** - Troubleshooting Tool
   - Manual debugging for CI failures
   - Environment inspection capabilities
   - Uses same setup as other workflows

3. **npm scripts** - Local Testing
   - Part of codebase (always in sync)
   - `npm run pre-commit`, `npm run test`, etc.
   - No external maintenance required

### **What We Removed (Problematic)**

- **Duplicate local scripts**: Avoided drift and maintenance issues
- **Redundant workflows**: Consolidated where possible
- **Manual sync requirements**: Everything auto-aligned

### **Benefits**

1. **No duplication**: All logic in workflows only
2. **Always in sync**: npm scripts part of codebase
3. **Simple maintenance**: Update workflow = everything updates
4. **Clear separation**: Local (npm) vs CI (workflows) vs Debug (manual)

## Local Development

### Before Pushing
```bash
# Quick local validation
npm run pre-commit

# Individual checks
npm run test          # Run tests
npm run lint          # Check code quality
npm run format        # Format code

# Full validation
npm run pre-commit:full
```

### Secret Management Options

#### Option 1: Use Local Vault (Recommended)
```bash
# Start Vault in dev mode
npm run vault:start

# Initialize with your secrets
npm run vault:init

# Your application reads from Vault
```

#### Option 2: Use .env File
Create a `.env` file from the template:
```bash
cp .env.example .env
# Edit .env with your values
```

**Note**: The `.env` file is in `.gitignore` and never committed.

#### Option 3: Export Environment Variables
```bash
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=your_password
# ... other variables
```

## Troubleshooting

### Local Testing
Use existing npm scripts to validate your environment:
```bash
# Test core functionality
npm run pre-commit

# Individual checks
npm run format
npm run lint
npm run test

# Full validation
npm run pre-commit:full
```

### CI Debugging
1. Go to Actions tab in GitHub
2. Run "Debug Workflow" manually
3. Select appropriate debug level
4. Review output for issues

### Common Issues

#### Vault Connection Issues
**Symptoms**: "Cannot connect to Vault" in CI logs
**Solutions**:
1. Verify `VAULT_ADDR` is correct and accessible
2. Check `VAULT_TOKEN` has required permissions
3. Ensure Vault allows connections from GitHub IPs

#### Node.js Version Mismatch
**Symptoms**: CI fails with compatibility errors
**Solution**: Update local Node.js to v22.13.0+

#### Go Version Mismatch
**Symptoms**: Go build/test failures in CI
**Solution**: Update local Go to v1.24.3+

#### Missing Secrets
**Symptoms**: Tests fail due to missing configuration
**Solutions**:
1. Verify all required GitHub Secrets are configured
2. Check secret names match exactly (case-sensitive)
3. Ensure secrets don't contain extra whitespace

#### npm Script Failures
**Symptoms**: npm commands timeout or fail
**Solution**: Run `npm run pre-commit:full` locally first

## Vault Integration

### Production Setup

For production environments with Vault:

1. Deploy Vault to your cloud provider
2. Configure the `thothix` mount path with KV v2 engine
3. Store secrets under the path `secret/data/thothix`
4. Create a token with appropriate policies
5. Configure `VAULT_ADDR` and `VAULT_TOKEN` in GitHub Secrets

#### Example Vault Policy
```hcl
path "secret/data/thothix" {
  capabilities = ["read"]
}

path "secret/metadata/thothix" {
  capabilities = ["read"]
}
```

### CI/CD Behavior

The GitHub Actions workflow automatically:

1. **Checks Vault connectivity** - Tests if `VAULT_ADDR` and `VAULT_TOKEN` work
2. **Uses Vault if available** - Reads secrets from Vault mount
3. **Falls back gracefully** - Uses GitHub Secrets if Vault unavailable
4. **Shows clear status** - Logs whether Vault is being used

You'll see in CI logs:
- ✅ Vault is accessible
- ⚠️ Running tests without Vault integration

### Development Workflow

```bash
# For Vault users
npm run vault:start     # Start local Vault
npm run vault:init      # Initialize secrets
npm run dev            # Start development

# For .env users
cp .env.example .env   # Create local config
# Edit .env with your values
npm run dev           # Start development
```

## Best Practices

### Development
1. **Test locally first**: Use `npm run pre-commit` before pushing
2. **Keep versions aligned**: Local environment should match CI
3. **Use templates**: Copy from `.env.example` for local setup
4. **Choose your secret strategy**: Vault for production, .env for development

### CI/CD
1. **Configure secrets properly**: Either Vault or GitHub Secrets
2. **Monitor workflows**: Check Actions tab for failures
3. **Use debug workflow**: For persistent CI issues
4. **Review logs**: Always check failure details

### Maintenance
When updating versions:
1. Update all workflow files consistently
2. Update package.json engines field
3. Update this documentation
4. Test locally with npm scripts

### Security
1. **Never commit secrets**: Use `.gitignore` and Vault
2. **Use Vault for production**: More secure than env vars
3. **Rotate tokens regularly**: Update GitHub Secrets when needed
4. **Review secret access**: Ensure minimal required permissions

---

For questions or issues, check the GitHub Actions logs or run the debug workflow for detailed diagnostics.
