# 🤖 Automazione Pre-Commit per Thothix

Questo documento spiega come utilizzare il sistema di automazione per formattazione, linting e test prima di ogni commit.

## 📁 **File Coinvolti**

```text
scripts/dev.bat           # Script Windows unificato per sviluppo
.vscode/tasks.json        # Task VS Code
.golangci.yml             # Configurazione linting
Makefile                  # Target Make
```

## 🎯 **Cosa Viene Automatizzato**

Il sistema automatizza:

1. **Formattazione del codice** (gofmt, gofumpt)
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

### **1. Script Manuale** (Raccomandato)

Per eseguire i check manualmente prima del commit:

```bash
# Windows - Script unificato
.\scripts\dev.bat format      # Solo formattazione
.\scripts\dev.bat lint        # Solo linting
.\scripts\dev.bat pre-commit  # Pre-commit completo
.\scripts\dev.bat all         # Equivalente a pre-commit
```

### **2. VS Code Tasks**

Nel Command Palette (`Ctrl+Shift+P`):

- **"Tasks: Run Task"** → **"Dev: Pre-commit"**
- **"Tasks: Run Task"** → **"Dev: Format"**
- **"Tasks: Run Task"** → **"Dev: Lint"**

### **3. Makefile** (se hai Make installato)

```bash
make pre-commit    # Esegue formattazione + lint + test
make commit        # Pre-commit + prepara per git commit
```

## 📋 **Workflow Consigliato**

### **Setup Iniziale** (una sola volta)

```bash
# 1. Verifica che tutto funzioni
.\scripts\dev.bat all
```

### **Workflow Quotidiano**

```bash
# 1. Modifica il codice
# 2. Esegui pre-commit check prima del commit
.\scripts\dev.bat pre-commit

# 3. Se tutto è corretto, procedi con il commit
git add .
git commit -m "Aggiunta nuova funzionalità"
```

## ⚙️ **Configurazione**

### **Personalizzazione del Pre-commit**

Modifica lo script `.\scripts\dev.bat` per:

- Disabilitare i test (rimuovi la chiamata a `go test`)
- Cambiare il timeout del linting
- Aggiungere altri check personalizzati

### **Disabilitare Temporaneamente**

```bash
# Salta i check pre-commit per un commit urgente
git commit --no-verify -m "Commit urgente"

# Oppure esegui solo formattazione senza linting
.\scripts\dev.bat format
```

## 🛠️ **Tool Utilizzati**

- **gofmt**: Formattazione base Go
- **gofumpt**: Formattazione più rigorosa
- **golangci-lint**: Linting completo
- **go test**: Esecuzione test

## 🔧 Troubleshooting

### **Errore: Tool non trovati**

```bash
# Installa i tool mancanti
go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### **Linting troppo rigoroso**

Modifica `.golangci.yml` per rilassare le regole o aggiungi esclusioni.

## 🔧 **Risoluzione Problemi Formattazione**

### **Diagnostica Automatica**

Se riscontri problemi di formattazione persistenti, usa lo script di sviluppo per diagnosticare:

```bash
.\scripts\dev.bat format
```

Questo comando:

- ✅ Verifica installazione Go e gofmt
- ✅ Applica la formattazione corretta sui file
- ✅ Mostra eventuali errori di formattazione
- ✅ Funziona con la configurazione VS Code ottimizzata

### **Configurazione VS Code Ottimizzata**

La configurazione in `.vscode/settings.json` è stata ottimizzata per evitare conflitti:

- **Formatter unico**: Solo `gofmt` (formato standard Go)
- **Import organizzazione**: Gestita automaticamente da gofmt
- **Azioni automatiche**: Minimizzate per prevenire interferenze

### **Soluzioni Comuni**

**Problema**: VS Code non formatta al salvataggio
**Soluzione**:

1. Riavvia VS Code completamente
2. Esegui "Go: Install/Update Tools" da Command Palette
3. Verifica che `.\scripts\dev.bat format` funzioni correttamente

**Problema**: Formattazione inconsistente
**Soluzione**:

1. Rimuovi configurazioni user che potrebbero sovrascrivere
2. Usa solo gli script forniti per formattazione manuale
3. Assicurati che non ci siano estensioni conflittuali

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
3. **Script di Correzione**: `.\scripts\dev.bat format` per correggere tutti i file
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

### **Problema: Conflitto goimports vs gofumpt**

**Sintomo**: Errori di formattazione persistenti anche dopo aver eseguito i tool.

**Causa**: `goimports` e `gofumpt` hanno regole diverse per l'organizzazione degli import.

**Soluzione Implementata**:

1. **Uso esclusivo di gofumpt**: Rimosso `goimports` da tutti gli script
2. **gofumpt include import organization**: Non serve più `goimports`
3. **VS Code configurato con gofumpt**: Formattazione consistente nell'IDE
4. **Tutti i workflow aggiornati**: Scripts, tasks, Git hooks, Makefile

**Comando Unico**:

```bash
gofumpt -w .  # Include formattazione + organizzazione import
```

**Tool Consolidati**:

- ❌ ~~gofmt + goimports + gofumpt~~ (conflitti)
- ✅ **Solo gofumpt** (tutto incluso, zero conflitti)
