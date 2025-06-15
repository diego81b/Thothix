# 🚀 Modern Cross-Platform Task Runners for Thothix

## 📋 Overview

Alternative moderne al Makefile che funzionano nativamente su Windows/Linux/macOS senza installazioni aggiuntive.

## 🎯 Opzioni Valutate

### **1. Just (Raccomandato)**
- **Sintassi**: Simile a Makefile ma più moderna
- **Installazione**: Singolo binario cross-platform
- **Pro**: Zero dependencies, sintassi pulita, supporto Windows nativo
- **Contro**: Tool relativamente nuovo

### **2. Task (Go-based)**
- **Sintassi**: YAML-based task runner
- **Installazione**: Singolo binario Go
- **Pro**: YAML familiare, molto potente, attiva community
- **Contro**: Più verboso di Makefile

### **3. NPM Scripts**
- **Sintassi**: JSON-based in package.json
- **Installazione**: Richiede Node.js
- **Pro**: Universalmente supportato, sintassi familiare
- **Contro**: Richiede Node.js ecosystem

### **4. PowerShell Core**
- **Sintassi**: PowerShell script
- **Installazione**: Disponibile su tutti i sistemi
- **Pro**: Potentissimo, Microsoft supportato
- **Contro**: Sintassi complessa per task semplici

## 🎯 Raccomandazione per Thothix

### **Strategia Ibrida** (Mantiene compatibilità)

1. **Makefile** - Per Linux/macOS (standard)
2. **Task** - Per Windows e multipiattaforma (moderno)
3. **Scripts nativi** - Fallback per compatibilità massima

## 📦 Setup Task Runner

### Installazione Task (https://taskfile.dev)

```bash
# Windows (PowerShell)
winget install Task.Task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

# macOS
brew install go-task/tap/go-task

# Go (tutti i sistemi)
go install github.com/go-task/task/v3/cmd/task@latest
```

### Taskfile.yml per Thothix

```yaml
version: '3'

tasks:
  format:
    desc: Format Go code
    cmds:
      - cd backend && gofmt -w .
    silent: false

  lint:
    desc: Run golangci-lint
    cmds:
      - cd backend && golangci-lint run --timeout=3m
    silent: false

  test:
    desc: Run tests
    cmds:
      - cd backend && go test ./...
    silent: false

  pre-commit:
    desc: Run pre-commit checks (format + lint + test)
    deps: [format]
    cmds:
      - git add backend/
      - task: lint
      - task: test

  dev:
    desc: Start development environment
    cmds:
      - docker compose up -d --build

  deploy:
    desc: Deploy to environment
    cmds:
      - |
        if [ "{{.ENV}}" = "" ]; then
          echo "Usage: task deploy ENV=dev|staging|prod"
          exit 1
        fi
        {{if eq OS "windows"}}
        scripts\deploy.bat {{.ENV}} up
        {{else}}
        scripts/deploy.sh {{.ENV}} up
        {{end}}

  db-verify:
    desc: Verify database schema
    cmds:
      - |
        {{if eq OS "windows"}}
        scripts\db-verify.bat {{.ACTION}}
        {{else}}
        scripts/db-verify.sh {{.ACTION}}
        {{end}}
```

## 🔄 Migrazione Graduale

### Fase 1: Aggiungere Task accanto a Makefile
- Mantenere Makefile esistente
- Aggiungere Taskfile.yml
- Documentare entrambe le opzioni

### Fase 2: Script wrapper universali
- Creare script di rilevamento OS
- Chiamare automaticamente Make o Task

### Fase 3: Unificazione (opzionale)
- Scegliere un solo standard basato su feedback team
