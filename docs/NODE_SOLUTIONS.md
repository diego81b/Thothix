# 🚀 Node.js Cross-Platform Solutions for Thothix

## 📋 Overview

Node.js offre eccellenti soluzioni multipiattaforma per script e task running, eliminando la necessità di mantenere versioni separate `.bat` e `.sh`.

## 🎯 Opzioni Node.js Disponibili

### **1. NPM Scripts** (Più Semplice)
- **File**: `package.json`
- **Pro**: Zero dependencies extra, supporto universale, sintassi familiare
- **Contro**: Limitato per script complessi
- **Ideale per**: Comandi semplici, progetti con Node.js

### **2. Zx (Google)** (Raccomandato)
- **Descrizione**: Write shell scripts in JavaScript
- **Pro**: JavaScript puro, cross-platform, potente, facile da leggere
- **Contro**: Richiede Node.js
- **Ideale per**: Script complessi, team JavaScript

### **3. Nx (Nrwl)** (Enterprise)
- **Descrizione**: Monorepo toolkit con task runner
- **Pro**: Molto potente, caching, dependency graph
- **Contro**: Overkill per progetti semplici
- **Ideale per**: Monorepo, progetti grandi

### **4. NPX Scripts** (Moderno)
- **Descrizione**: Script eseguibili con npx
- **Pro**: No install globale, versioning
- **Contro**: Richiede npm
- **Ideale per**: Script distribuibili

## 🎯 Raccomandazione per Thothix: **Zx**

### Perché Zx è Perfetto per Thothix

1. **Cross-platform nativo**: Stesso script su Windows/Linux/macOS
2. **JavaScript familiare**: Facile per team moderni
3. **Shell integration**: Può chiamare comandi shell nativi
4. **Zero config**: Funziona out-of-the-box
5. **Potente**: Supporta async/await, import, TypeScript

## 📦 Setup con Zx

### Installazione

```bash
# Globale (una volta)
npm install -g zx

# Oppure locale al progetto
npm install --save-dev zx
```

### Struttura Proposta

```text
scripts/
├── dev.mjs          # Script Zx unificato per sviluppo
├── deploy.mjs       # Script Zx unificato per deploy
├── db-verify.mjs    # Script Zx unificato per database
└── package.json     # NPM scripts per shortcuts
```

## 🔧 Implementazione Zx per Thothix

### scripts/dev.mjs

```javascript
#!/usr/bin/env zx

import { $, argv, echo, cd } from 'zx'

$.verbose = true

const action = argv._[0] || 'help'

echo`🔧 Thothix Development Script - Action: ${action}`

// Navigate to backend directory
await cd('backend')

const format = async () => {
  echo`📝 Formatting Go code...`
  await $`gofmt -w .`
  echo`✅ Formatting completed`
}

const lint = async () => {
  echo`🔍 Running golangci-lint...`
  await $`golangci-lint run --timeout=3m`
  echo`✅ Linting completed`
}

const test = async () => {
  echo`🧪 Running tests...`
  await $`go test ./...`
  echo`✅ Tests completed`
}

const preCommit = async () => {
  echo`🚀 Running pre-commit checks...`
  await format()
  echo`📋 Adding formatted files to git...`
  await cd('..')
  await $`git add backend/`
  await cd('backend')
  await lint()
  await test()
  echo`✅ Pre-commit checks completed`
}

switch (action) {
  case 'format':
    await format()
    break
  case 'lint':
    await lint()
    break
  case 'test':
    await test()
    break
  case 'pre-commit':
  case 'all':
    await preCommit()
    break
  case 'help':
  default:
    echo`
Usage: zx scripts/dev.mjs [action]

Actions:
  format      - Format Go code
  lint        - Run golangci-lint
  test        - Run tests
  pre-commit  - Run all checks
  all         - Same as pre-commit

Examples:
  zx scripts/dev.mjs format
  zx scripts/dev.mjs pre-commit
`
    break
}

echo`🎉 Script completed successfully!`
```

### scripts/deploy.mjs

```javascript
#!/usr/bin/env zx

import { $, argv, echo, fs, path } from 'zx'

const env = argv._[0] || 'help'
const cmd = argv._[1] || 'up'

if (env === 'help') {
  echo`
Usage: zx scripts/deploy.mjs [env] [command]

Environments:
  dev      - Development environment (.env)
  staging  - Staging environment (.env.staging)
  prod     - Production environment (.env.prod)

Commands:
  up       - Start services
  down     - Stop services
  logs     - Show logs
  status   - Show container status
`
  process.exit(0)
}

// Environment configuration
const configs = {
  dev: {
    envFile: '.env',
    composeFiles: ['-f', 'docker-compose.yml']
  },
  staging: {
    envFile: '.env.staging',
    composeFiles: ['-f', 'docker-compose.yml', '-f', 'docker-compose.staging.yml']
  },
  prod: {
    envFile: '.env.prod',
    composeFiles: ['-f', 'docker-compose.yml', '-f', 'docker-compose.prod.yml']
  }
}

const config = configs[env]
if (!config) {
  echo`❌ Invalid environment: ${env}`
  process.exit(1)
}

// Check environment file exists
if (!await fs.pathExists(config.envFile)) {
  echo`❌ Environment file ${config.envFile} not found!`
  echo`📋 Please copy .env.example to ${config.envFile} and configure it`
  process.exit(1)
}

echo`🚀 Thothix Deployment - Environment: ${env}, Command: ${cmd}`

const dockerCompose = async (args) => {
  return $`docker compose ${config.composeFiles} --env-file=${config.envFile} ${args}`
}

switch (cmd) {
  case 'up':
    echo`📦 Starting ${env} environment...`
    await dockerCompose`up -d --build`
    echo`✅ ${env} environment started successfully`
    await dockerCompose`ps`
    break

  case 'down':
    echo`🛑 Stopping ${env} environment...`
    await dockerCompose`down`
    echo`✅ ${env} environment stopped`
    break

  case 'logs':
    echo`📋 Showing logs...`
    await dockerCompose`logs -f`
    break

  case 'status':
    echo`📊 Container status for ${env} environment:`
    await dockerCompose`ps`
    break

  default:
    echo`❌ Invalid command: ${cmd}`
    process.exit(1)
}
```

### package.json (Shortcuts)

```json
{
  "name": "thothix-scripts",
  "version": "1.0.0",
  "private": true,
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
    "zx": "^8.0.0"
  }
}
```

## 🚀 Vantaggi Soluzione Node.js/Zx

### ✅ **Unificazione Completa**
- **Un solo script** per tutte le piattaforme
- **Nessuna duplicazione** di logica
- **Manutenzione semplificata**

### ✅ **Developer Experience**
- **Sintassi familiare** (JavaScript)
- **Async/await** per operazioni complesse
- **Error handling** moderno
- **Import/export** per modularità

### ✅ **Potenza**
- **Shell commands** nativi con `$`
- **File system operations** cross-platform
- **Process management** avanzato
- **Environment detection** automatico

### ✅ **Integrazione**
- **NPM scripts** per shortcuts
- **VS Code** task integration
- **Git hooks** supportati
- **CI/CD** friendly

## 🔄 Migrazione Graduale

### Fase 1: Aggiungere Zx accanto agli esistenti
```text
scripts/
├── dev.bat         # Esistente Windows
├── dev.sh          # Esistente Unix
├── dev.mjs         # Nuovo Zx unificato
└── package.json    # NPM shortcuts
```

### Fase 2: Wrapper che preferisce Zx
```javascript
// Wrapper che cerca prima Zx, poi fallback a nativi
if (hasNodeJS && hasZx) {
  await $`zx scripts/dev.mjs ${action}`
} else {
  await $`scripts/dev.${isWindows ? 'bat' : 'sh'} ${action}`
}
```

### Fase 3: Deprecazione graduale scripts nativi
- Mantenere per compatibilità
- Documentare Zx come standard
- Team migration su base volontaria

## 💡 Alternative Node.js

### **NPM Scripts Only** (Minimal)
```json
{
  "scripts": {
    "format": "cd backend && gofmt -w .",
    "lint": "cd backend && golangci-lint run --timeout=3m",
    "pre-commit": "npm run format && git add backend/ && npm run lint && npm run test",
    "dev": "docker compose up -d --build"
  }
}
```

### **Nx Integration** (Advanced)
```json
{
  "scripts": {
    "nx": "nx",
    "format": "nx run-many --target=format --all",
    "lint": "nx run-many --target=lint --all",
    "pre-commit": "nx affected --target=pre-commit"
  }
}
```

## 🎯 Raccomandazione Finale

**Zx è perfetto per Thothix** perché:

1. **Zero breaking changes** - Scripts esistenti rimangono
2. **Gradual adoption** - Team può migrare quando pronto
3. **Universal solution** - Un script, tutte le piattaforme
4. **Modern tooling** - JavaScript ecosystem familiare
5. **Powerful yet simple** - Più potente di shell, più semplice di tool complessi

Vuoi che implementi la migrazione con Zx mantenendo la compatibilità esistente?
