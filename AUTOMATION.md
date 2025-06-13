# 🤖 Automazione Pre-Commit per Thothix

Questo documento spiega come utilizzare il sistema di automazione per formattazione, linting e test prima di ogni commit.

## 🎯 **Cosa Viene Automatizzato**

Il sistema automatizza:

1. **Formattazione del codice** (gofmt, goimports, gofumpt)
2. **Riaggiunta automatica** dei file formattati al commit
3. **Linting** (golangci-lint)
4. **Test** (go test - opzionale)
5. **Blocco del commit** se ci sono errori

### 🔧 **Processo Automatico**

Quando esegui `git commit`, il sistema:

1. Formatta automaticamente tutto il codice Go
2. Aggiunge i file formattati al commit corrente
3. Esegue golangci-lint per verificare la qualità
4. Esegue i test (se abilitati)
5. Procede con il commit solo se tutto è corretto

## 🚀 **Metodi di Utilizzo**

### **1. Git Hook Automatico** (Raccomandato)

Il pre-commit hook si attiva automaticamente ad ogni `git commit`:

```bash
# Configura una sola volta
.\scripts\setup-hooks.ps1

# Poi ogni commit attiverà automaticamente i check
git add .
git commit -m "Il tuo messaggio"
```

### **2. Script Manuale**

Per eseguire i check manualmente prima del commit:

```bash
# Windows
.\scripts\pre-commit.bat

# Oppure con PowerShell
.\scripts\setup-hooks.ps1
```

### **3. VS Code Tasks**

Nel Command Palette (`Ctrl+Shift+P`):

- **"Tasks: Run Task"** → **"Go: Pre-commit"**
- **"Tasks: Run Task"** → **"Git: Setup Hooks"**

### **4. Makefile** (se hai Make installato)

```bash
make pre-commit    # Esegue formattazione + lint + test
make setup-hooks   # Configura Git hooks
make commit        # Pre-commit + prepara per git commit
```

## 📋 **Workflow Consigliato**

### **Setup Iniziale** (una sola volta)

```bash
# 1. Configura i Git hooks
.\scripts\setup-hooks.ps1

# 2. Verifica che tutto funzioni
.\scripts\pre-commit.bat
```

### **Workflow Quotidiano**

```bash
# 1. Modifica il codice
# 2. Commit normale - l'automazione si attiva automaticamente
git add .
git commit -m "Aggiunta nuova funzionalità"

# Se il pre-commit fallisce, risolvi gli errori e riprova
git commit -m "Fix dopo linting"
```

## ⚙️ **Configurazione**

### **Personalizzazione del Pre-commit**

Modifica `.git/hooks/pre-commit` per:

- Disabilitare i test (commenta la sezione test)
- Cambiare il timeout del linting
- Aggiungere altri check personalizzati

### **Disabilitare Temporaneamente**

```bash
# Salta il pre-commit hook per un commit urgente
git commit --no-verify -m "Commit urgente"
```

## 🛠️ **Tool Utilizzati**

- **gofmt**: Formattazione base Go
- **goimports**: Gestione automatica import
- **gofumpt**: Formattazione più rigorosa
- **golangci-lint**: Linting completo
- **go test**: Esecuzione test

## 📁 **File Coinvolti**

```
.git/hooks/pre-commit          # Git hook principale
scripts/pre-commit.bat         # Script Windows manuale
scripts/setup-hooks.ps1        # Setup PowerShell
.vscode/tasks.json            # Task VS Code
.golangci.yml                 # Configurazione linting
Makefile                      # Target Make
```

## 🔧 **Troubleshooting**

### **Errore: Hook non eseguibile**

```bash
# Su Linux/macOS
chmod +x .git/hooks/pre-commit

# Su Windows
.\scripts\setup-hooks.ps1
```

### **Errore: Tool non trovati**

```bash
# Installa i tool mancanti
go install golang.org/x/tools/cmd/goimports@latest
go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### **Linting troppo rigoroso**

Modifica `.golangci.yml` per rilassare le regole o aggiungi esclusioni.

## ✅ **Vantaggi**

- ✅ **Codice sempre formattato** correttamente
- ✅ **Nessun errore di linting** nel repository
- ✅ **Test sempre passanti** prima del commit
- ✅ **Processo automatico** e trasparente
- ✅ **Configurabile** secondo le necessità del team

Il sistema garantisce che il repository mantenga sempre un alto standard di qualità del codice! 🎉

### **Problemi di Formattazione Risolti Automaticamente**

Il sistema risolve automaticamente:

- **Spaziatura degli import** Go secondo le convenzioni gofumpt
- **Ordinamento degli import** con gruppi separati correttamente
- **Formattazione del codice** secondo gli standard Go
- **Aggiunta automatica** dei file corretti al commit

Se vedi errori come "File is not properly formatted (gofumpt)", il sistema li corregge automaticamente e riaggiunge i file al commit.

### **Problema: VS Code Rompe la Formattazione Go**

**Sintomo**: VS Code aggiunge spazi extra nelle sezioni import causando errori `gofumpt`.

**Soluzione**:

1. **Configurazione VS Code Corretta**: Il file `.vscode/settings.json` è configurato per evitare questo problema
2. **Task Manuale**: Usa `Ctrl+Shift+P` → "Tasks: Run Task" → "Go: Format Current File"
3. **Script di Correzione**: `.\scripts\fix-formatting.bat` per correggere tutti i file
4. **Formattazione Automatica**: Il pre-commit hook corregge automaticamente questi errori

**Formato Import Corretto**:

```go
import (
    // Librerie standard Go
    "fmt"
    "net/http"

    // Moduli locali
    "thothix-backend/internal/models"

    // Librerie esterne
    "github.com/gin-gonic/gin"
)
```

**Prevenzione**: Le impostazioni VS Code sono configurate per:

- Disabilitare format-on-paste/type
- Usare `goimports` come formattatore primario
- Trim automatico whitespace
- Rilevamento indentazione disabilitato
