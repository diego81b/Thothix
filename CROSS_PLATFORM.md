# 🌍 Cross-Platform Development Guide

## 📋 Overview

Thothix supporta sviluppo su **Windows**, **Linux** e **macOS** attraverso:

1. **Makefile** universale (raccomandato per tutti i sistemi)
2. **Script nativi** per ogni piattaforma
3. **VS Code tasks** multipiattaforma

## 🚀 Quick Start per Piattaforma

### Windows
```cmd
# Wrapper universale (raccomandato)
.\run.bat pre-commit

# Script nativo Windows
.\scripts\dev.bat pre-commit

# Task runner moderno (se installato)
task pre-commit

# Makefile (se hai Make installato)
make pre-commit
```

### Linux/macOS
```bash
# Wrapper universale (raccomandato)
chmod +x run && ./run pre-commit

# Script nativo Unix
chmod +x scripts/dev.sh && ./scripts/dev.sh pre-commit

# Task runner moderno (se installato)
task pre-commit

# Makefile (disponibile di default)
make pre-commit
```

### VS Code (tutte le piattaforme)
- `Ctrl+Shift+P` → "Tasks: Run Task" → "Dev: Pre-commit"

## 📁 Struttura Scripts

```text
scripts/
├── dev.bat         # Windows development tasks
├── dev.sh          # Unix development tasks
├── deploy.bat      # Windows Docker deployment
├── deploy.sh       # Unix Docker deployment
├── db-verify.bat   # Windows database utilities
└── db-verify.sh    # Unix database utilities

# Universal wrappers
├── run.bat         # Windows universal wrapper
├── run             # Unix universal wrapper
├── Taskfile.yml    # Modern task runner (all platforms)
└── Makefile        # Traditional make targets (all platforms)
```

## 🎯 Raccomandazioni per Utilizzo

### Opzione 1: Wrapper Universale (Più Semplice)
```bash
# Rileva automaticamente il miglior tool disponibile
./run pre-commit    # Unix
.\run.bat pre-commit # Windows
```

### Opzione 2: Task Runner Moderno (Più Potente)
```bash
# Installa Task: https://taskfile.dev/installation/
task pre-commit     # Stesso comando su tutti i sistemi
task --list         # Mostra tutti i task disponibili
```

### Opzione 3: Makefile Tradizionale (Più Compatibile)
```bash
make pre-commit     # Funziona su Linux/macOS, Windows con Make
```

### Opzione 4: Script Nativi (Massima Compatibilità)
```bash
./scripts/dev.sh pre-commit    # Unix
.\scripts\dev.bat pre-commit   # Windows
```

## ⚙️ Configurazione per Sistema

### Windows
- **Shell**: `cmd.exe` o PowerShell
- **Dependencies**: Go, golangci-lint, git
- **Script**: `dev.bat`

### Linux/Ubuntu
```bash
# Installa dipendenze
sudo apt update
sudo apt install golang-go git make

# Installa golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Rendi eseguibile lo script
chmod +x scripts/dev.sh
```

### macOS
```bash
# Installa dipendenze con Homebrew
brew install go git make golangci-lint

# Rendi eseguibile lo script
chmod +x scripts/dev.sh
```

## 🎯 Comandi Unificati

Tutti i sistemi supportano gli stessi comandi:

| Azione     | Windows                        | Unix                          | Make              | VS Code Task      |
| ---------- | ------------------------------ | ----------------------------- | ----------------- | ----------------- |
| Format     | `.\scripts\dev.bat format`     | `./scripts/dev.sh format`     | `make format`     | "Dev: Format"     |
| Lint       | `.\scripts\dev.bat lint`       | `./scripts/dev.sh lint`       | `make lint`       | "Dev: Lint"       |
| Pre-commit | `.\scripts\dev.bat pre-commit` | `./scripts/dev.sh pre-commit` | `make pre-commit` | "Dev: Pre-commit" |
| All        | `.\scripts\dev.bat all`        | `./scripts/dev.sh all`        | `make pre-commit` | "Dev: Pre-commit" |

## 🐳 Docker (Multipiattaforma)

Docker funziona identicamente su tutte le piattaforme:

```bash
# Sviluppo
docker compose up -d --build

# Produzione
docker compose -f docker-compose.prod.yml up -d
```

## 🔧 Setup IDE per Team Misto

### VS Code Settings (multipiattaforma)

Il file `.vscode/settings.json` funziona su tutti i sistemi:

```json
{
    "go.formatTool": "gofmt",
    "go.lintTool": "golangci-lint",
    "[go]": {
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true
        }
    }
}
```

### VS Code Tasks (multipiattaforma)

I task in `.vscode/tasks.json` si adattano automaticamente:

```json
{
    "label": "Dev: Pre-commit",
    "type": "shell",
    "command": "${workspaceFolder}/scripts/dev.${shellName === 'cmd' ? 'bat' : 'sh'}",
    "args": ["pre-commit"],
    "group": "build"
}
```

## 📝 Best Practices Team Misto

### 1. Makefile come Standard
```makefile
# Funziona su tutti i sistemi
make pre-commit
make format
make lint
```

### 2. Script Nativi come Fallback
- **Windows**: `dev.bat` per compatibilità cmd.exe
- **Unix**: `dev.sh` per bash/zsh

### 3. Docker come Equalizer
```bash
# Stesso ambiente su tutti i sistemi
docker compose exec backend ./scripts/dev.sh pre-commit
```

### 4. Git Hooks Universali

Crea `.git/hooks/pre-commit` che rileva il sistema:

```bash
#!/bin/sh
# Universal pre-commit hook

if [ -x "scripts/dev.sh" ]; then
    ./scripts/dev.sh pre-commit
elif [ -x "scripts/dev.bat" ]; then
    ./scripts/dev.bat pre-commit
else
    make pre-commit
fi
```

## 🚨 Troubleshooting Cross-Platform

### Windows Subsystem for Linux (WSL)
```bash
# In WSL, usa gli script Unix
./scripts/dev.sh pre-commit
```

### Git Line Endings
```bash
# Configura una sola volta
git config core.autocrlf true  # Windows
git config core.autocrlf input # Unix
```

### Path Separators
Scripts gestiscono automaticamente `\` (Windows) vs `/` (Unix).

## 🎯 Raccomandazioni per Team

1. **Sviluppatori Windows**: Usare `dev.bat` o Makefile
2. **Sviluppatori Unix**: Usare `dev.sh` o Makefile
3. **VS Code**: Usare i task integrati
4. **CI/CD**: Usare sempre Makefile per consistenza
5. **Docker**: Usare per uniformità ambiente

La filosofia **single-version policy** è mantenuta: un comando, un comportamento, su tutte le piattaforme! 🌍
