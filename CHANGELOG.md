# Changelog - Automazione e Qualità del Codice

## [Unreleased]

### Infrastructure

- feat: create simplified version-bump.bat script
  - Implemented pure Windows batch script for version bumping
  - Supports semantic versioning (major/minor/patch)
  - Automatic version parsing from CHANGELOG.md
  - No external dependencies (PowerShell, Unix tools)
  - Simplified single-file approach following project guidelines
  - **Impact**: Functional version management tool ready for testing and refinement

- feat: simplify script architecture with single-version policy
  - Removed duplicate version-bump scripts (.ps1, .sh)
  - Standardized on .bat files for Windows primary development
  - Updated Copilot instructions with script development guidelines
  - Cleaned scripts README.md documentation
  - Updated VS Code tasks to use simplified scripts
  - **Impact**: Reduced maintenance overhead, eliminated confusion, improved developer experience

### Documentation

- docs: enhance Copilot instructions with mandatory CHANGELOG updates
  - Added comprehensive CHANGELOG guidelines with examples and best practices
  - Made CHANGELOG updates mandatory in pre-commit checklist
  - Included detailed formatting standards for CHANGELOG entries
  - Added release process documentation with semantic versioning
  - **Impact**: Ensures consistent and detailed change tracking for all commits

## v1.3.0 - Consolidamento Documentazione (2025-06-14)

### Documentation Updates

- docs: consolidate Clerk documentation into single comprehensive file
  - Merged CLERK_INTEGRATION.md and CLERK_WEBHOOK_SETUP.md into CLERK_INTEGRATION_COMPLETE.md
  - Removed duplicate and obsolete documentation files
  - Updated all references to consolidated documentation
  - **Impact**: Simplified documentation structure with single source of truth for Clerk integration

- docs: simplify README structure and remove redundancies
  - Integrated backend documentation into main README
  - Removed duplicate setup instructions
  - Updated project structure documentation
  - **Impact**: Cleaner, more maintainable documentation

## v1.2.0 - Automazione Completa (2025-06-13)

### ✨ Nuove Funzionalità

- **Sistema di automazione pre-commit completo**
- **Git hooks automatici** per formattazione e linting
- **Script cross-platform** (Windows/Unix) per sviluppo
- **VS Code tasks** integrate per workflow di sviluppo

### 🔧 Strumenti di Formattazione

- **gofmt**: Formattazione base Go
- **goimports**: Gestione automatica import
- **gofumpt**: Formattazione rigorosa per CI/CD
- **golangci-lint**: Linting configurato con regole rilassate

### 🛠️ Configurazione Migliorata

- **`.golangci.yml`** ottimizzato per produttività sviluppatori
- **Makefile** con target per tutte le operazioni comuni
- **VS Code settings** per auto-formattazione
- **Scripts PowerShell/Batch** per setup automatico

### 🐛 Problemi Risolti

- ✅ Errori di formattazione `gofumpt` risolti automaticamente
- ✅ Spaziatura import Go corretta secondo convenzioni
- ✅ Aggiunta automatica file formattati al commit
- ✅ Hook pre-commit robusto con gestione errori
- ✅ **NUOVO**: Problema VS Code che rompeva formattazione import Go
- ✅ **NUOVO**: Configurazione VS Code ottimizzata per evitare spazi extra
- ✅ **NUOVO**: Task VS Code per formattazione file singolo
- ✅ **NUOVO**: Script batch per correzione formattazione massive
- ✅ **RISOLTO**: Conflitto tra goimports e gofumpt che causava errori di formattazione persistenti

### 🧹 Pulizia e Ottimizzazione Script

- ✅ **NUOVO**: Script unificato `dev.bat/sh` con azioni multiple (format|lint|pre-commit|all)
- ✅ **RIMOSSI**: Script duplicati `format.bat/sh`, `fix-formatting.bat`
- ✅ **CONSOLIDATO**: Tutte le funzionalità di sviluppo in script unici
- ✅ **SEMPLIFICATO**: Workflow sviluppo con comandi chiari e intuitivi

### 📚 Documentazione

- **AUTOMATION.md**: Guida completa all'automazione
- **README.md**: Aggiornato con setup sviluppo
- **backend/README.md**: Documentazione strumenti sviluppo
- **Troubleshooting**: Sezioni per risoluzione problemi comuni
