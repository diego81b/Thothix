# 🌍 Node.js Development Guide for Thothix

## 📋 Overview

Thothix uses **Node.js with Zx** as the single, unified solution for all development scripts and automation. This approach provides true cross-platform compatibility without any need for separate Windows/Unix scripts or wrapper files.

## 🎯 Why Node.js + Zx?

### ✅ **Key Advantages**

- **Single script** for all platforms (Windows, Linux, macOS)
- **Familiar syntax** (Modern JavaScript)
- **Zero logic duplication**
- **Simplified maintenance**
- **Powerful** - Native shell commands with async/await
- **Mature ecosystem** - NPM, VS Code integration, CI/CD friendly
- **No wrappers needed** - npm scripts work everywhere

### ✅ **vs Other Solutions**

| Approach                          | Pros                        | Cons                   | Thothix Choice |
| --------------------------------- | --------------------------- | ---------------------- | -------------- |
| **Native scripts** (`.bat`/`.sh`) | Maximum compatibility       | Logic duplication      | ❌ Removed      |
| **Make**                          | Universal standard          | Extra setup on Windows | ❌ Removed      |
| **Task runners** (Task, etc.)     | Powerful                    | Additional tools       | ❌ Removed      |
| **Node.js/Zx**                    | Unified, powerful, familiar | Requires Node.js       | ✅ **Chosen**   |

## 🚀 Quick Start

### Installation

```bash
# 1. Install Node.js
# https://nodejs.org/

# 2. Install dependencies
npm install

# 3. Verify setup
npm run help
```

### Main Commands

```bash
# Development
npm run format      # Format Go code
npm run lint        # Run golangci-lint
npm run test        # Run Go tests
npm run pre-commit  # Format + lint + test

# Environments
npm run dev         # Start development
npm run staging     # Start staging
npm run prod        # Start production

# Database
npm run db:status   # Check database
npm run db:connect  # Connect to DB
npm run db:tables   # List tables
npm run db:check    # Check BaseModel
```

**That's it!** All commands work identically on every platform. 🎉

## 📁 Project Structure

```text
thothix/
├── package.json           # npm scripts configuration
├── scripts/
│   ├── dev.mjs           # Development automation (Zx)
│   ├── deploy.mjs        # Deployment automation (Zx)
│   └── db-verify.mjs     # Database operations (Zx)
└── backend/              # Go application code
```

## 🔧 Unified Scripts

### scripts/dev.mjs

Handles all development tasks:

```javascript
#!/usr/bin/env zx
import { $, argv, echo, cd } from 'zx'

const action = argv._[0] || 'help'

// Cross-platform format
const format = async () => {
  echo`📝 Formatting Go code...`
  await cd('backend')
  await $`gofmt -w .`
  echo`✅ Formatting completed`
}

// Cross-platform lint
const lint = async () => {
  echo`🔍 Running golangci-lint...`
  await cd('backend')
  await $`golangci-lint run --timeout=3m`
  echo`✅ Linting passed`
}
```

### scripts/deploy.mjs

Handles all deployment environments:

```javascript
#!/usr/bin/env zx
import { $, argv, echo } from 'zx'

const environment = argv._[0] || 'help'
const action = argv._[1] || 'up'

// Cross-platform Docker operations
const deployEnvironment = async (env, cmd) => {
  echo`🚀 Docker ${cmd} for ${env} environment...`

  if (env === 'dev') {
    await $`docker compose ${cmd} --build`
  } else {
    await $`docker compose -f docker-compose.yml -f docker-compose.${env}.yml ${cmd} --build`
  }
}
```

### scripts/db-verify.mjs

Handles all database operations:

```javascript
#!/usr/bin/env zx
import { $, argv, echo, question } from 'zx'

const action = argv._[0] || 'help'

// Cross-platform database connection
const checkDatabase = async () => {
  echo`🔍 Checking database status...`

  try {
    await $`docker compose exec postgres pg_isready -U thothix_user`
    echo`✅ Database is ready`
  } catch (error) {
    echo`❌ Database is not ready`
    process.exit(1)
  }
}
```

## 📦 package.json Configuration

All scripts are defined in `package.json` for maximum portability:

```json
{
  "name": "thothix-scripts",
  "scripts": {
    "format": "zx scripts/dev.mjs format",
    "lint": "zx scripts/dev.mjs lint",
    "test": "zx scripts/dev.mjs test",
    "pre-commit": "zx scripts/dev.mjs pre-commit",
    "dev": "zx scripts/deploy.mjs dev up",
    "dev:down": "zx scripts/deploy.mjs dev down",
    "staging": "zx scripts/deploy.mjs staging up",
    "prod": "zx scripts/deploy.mjs prod up",
    "db:status": "zx scripts/db-verify.mjs status",
    "db:connect": "zx scripts/db-verify.mjs connect"
  },
  "devDependencies": {
    "zx": "^8.1.4"
  }
}
```

## 🚀 Development Workflows

### Daily Development

```bash
# Start working
npm install              # Ensure dependencies
npm run pre-commit      # Check code quality
npm run dev             # Start development environment

# Make changes, then...
npm run format          # Format code
npm run lint           # Check linting
npm run test           # Run tests
npm run pre-commit     # Full check before commit

# Commit your changes
git add .
git commit -m "feat: your changes"
```

### Environment Management

```bash
# Development
npm run dev              # Start development
npm run dev:down         # Stop development
npm run dev:logs         # View logs

# Staging
npm run staging          # Deploy to staging
npm run staging:down     # Stop staging

# Production
npm run prod             # Deploy to production
npm run prod:down        # Stop production
```

### Database Operations

```bash
# Quick checks
npm run db:status        # Check if database is running
npm run db:tables        # List all tables

# Interactive
npm run db:connect       # Connect to database shell
npm run db:check         # Verify BaseModel implementation
```

## 🔬 Advanced Usage

### Running Specific Commands

```bash
# Run specific Zx scripts directly
npx zx scripts/dev.mjs format
npx zx scripts/deploy.mjs dev up
npx zx scripts/db-verify.mjs status

# Chain operations
npm run format && npm run lint && npm run test
```

### Custom Environment Variables

```bash
# Set environment variables
ENVIRONMENT=staging npm run dev
LOG_LEVEL=debug npm run dev:logs
```

### Integration with IDE

VS Code tasks automatically use npm scripts:

```json
{
  "label": "Dev: Format",
  "type": "shell",
  "command": "npm",
  "args": ["run", "format"],
  "group": "build"
}
```

## 🎯 Best Practices

### 1. Always Use npm scripts

```bash
# ✅ Good - Cross-platform
npm run format

# ❌ Avoid - Platform-specific
./scripts/format.sh
.\scripts\format.bat
```

### 2. Leverage Zx Features

```javascript
// ✅ Good - Use Zx built-ins
import { $, echo, cd, question } from 'zx'

// ✅ Good - Error handling
try {
  await $`command-that-might-fail`
} catch (error) {
  echo`❌ Command failed: ${error.message}`
  process.exit(1)
}
```

### 3. Keep Scripts Simple

```javascript
// ✅ Good - Clear, focused scripts
const format = async () => {
  echo`📝 Formatting Go code...`
  await cd('backend')
  await $`gofmt -w .`
  echo`✅ Formatting completed`
}

// ❌ Avoid - Complex monolithic scripts
const doEverything = async () => {
  // 200 lines of mixed responsibilities
}
```

## 🛠️ Troubleshooting

### Node.js Not Found

```bash
# Check Node.js installation
node --version

# Install Node.js from https://nodejs.org/
# Minimum version: Node.js 16+
```

### Zx Not Available

```bash
# Install project dependencies
npm install

# Or install Zx globally
npm install -g zx
```

### Script Execution Issues

```bash
# Check script permissions (Unix)
chmod +x scripts/*.mjs

# Run with explicit Node.js
node --loader=zx/esm scripts/dev.mjs format

# Debug mode
DEBUG=1 npm run format
```

### Docker Issues

```bash
# Verify Docker is running
docker --version
docker compose version

# Check container status
npm run dev:logs
```

## 🔄 Migration from Old Scripts

If you have old `.bat` or `.sh` scripts, here's how to migrate:

### 1. Identify Logic

```bash
# Old: format.bat
@echo off
cd backend
gofmt -w .

# Old: format.sh
#!/bin/bash
cd backend
gofmt -w .
```

### 2. Convert to Zx

```javascript
// New: scripts/dev.mjs
const format = async () => {
  echo`📝 Formatting Go code...`
  await cd('backend')
  await $`gofmt -w .`
  echo`✅ Formatting completed`
}
```

### 3. Add npm Script

```json
{
  "scripts": {
    "format": "zx scripts/dev.mjs format"
  }
}
```

### 4. Remove Old Files

```bash
# Clean up old scripts
rm scripts/*.bat scripts/*.sh
```

## 📈 Performance Considerations

### Script Startup Time

```bash
# Zx scripts have minimal startup overhead
time npm run format
# real    0m2.1s (includes Node.js startup)

# Compare with native scripts
time ./format.sh
# real    0m1.8s (minimal difference)
```

### Memory Usage

- Zx: ~50MB RAM per script execution
- Native scripts: ~5MB RAM per script execution
- **Trade-off**: Slightly higher resource usage for significant development benefits

### CI/CD Performance

```yaml
# GitHub Actions example
- name: Install dependencies
  run: npm ci

- name: Run pre-commit checks
  run: npm run pre-commit
  # Cross-platform, no matrix needed!
```

## 🎉 Summary

**Thothix now uses a single, modern automation solution:**

- ✅ **One codebase** - Node.js/Zx scripts work everywhere
- ✅ **No duplication** - Same logic across all platforms
- ✅ **Industry standard** - npm scripts are universally supported
- ✅ **Simple workflow** - `npm run <command>` is all you need
- ✅ **Powerful features** - Full JavaScript ecosystem available
- ✅ **Easy maintenance** - Single source of truth for all automation

**Start developing with confidence - it just works! 🚀**

---

📝 **Next Steps**: Check out the main [README.md](../README.md) for project overview and [CHANGELOG.md](../CHANGELOG.md) for recent updates.
