# Makefile per il progetto Thothix - Node.js/Zx Integration
# Uso: make <target>

.PHONY: help format lint test build clean dev pre-commit install node-check

# Default target
help: node-check
	@echo "🚀 Comandi disponibili per Thothix (Node.js/Zx):"
	@echo ""
	@echo "  install    - Installa le dipendenze Node.js"
	@echo "  format     - Formatta tutto il codice Go"
	@echo "  lint       - Esegue golangci-lint"
	@echo "  test       - Esegue i test"
	@echo "  pre-commit - Esegue format + lint + test (equivale a dev script)"
	@echo "  build      - Compila l'applicazione"
	@echo "  clean      - Pulisce i file temporanei"
	@echo "  dev        - Avvia in modalità sviluppo"
	@echo ""
	@echo "💡 Raccomandato: usa direttamente Node.js/Zx:"
	@echo "  npm run format"
	@echo "  npm run pre-commit"
	@echo "  npm run dev"
	@echo "  ./run format      # Wrapper universale"

# Check Node.js availability
node-check:
	@which node >/dev/null || (echo "❌ Node.js not found! Install: https://nodejs.org/" && exit 1)
	@test -f package.json || (echo "❌ package.json not found!" && exit 1)

# Install Node.js dependencies
install: node-check
	@echo "� Installing Node.js dependencies..."
	@npm install

# Formattazione del codice (via Node.js/Zx)
format: node-check
	@echo "🔧 Formattazione del codice (via Node.js/Zx)..."
	@npm run format

# Linting (via Node.js/Zx)
lint: node-check
	@echo "🔍 Eseguendo golangci-lint (via Node.js/Zx)..."
	@npm run lint

# Test (via Node.js/Zx)
test: node-check
	@echo "🧪 Eseguendo test (via Node.js/Zx)..."
	@npm run test

# Pre-commit completo (via Node.js/Zx)
pre-commit: node-check
	@echo "� Pre-commit checks (via Node.js/Zx)..."
	@npm run pre-commit

# Test
test:
	@echo "🧪 Eseguendo i test..."
	@cd backend && go test ./...
	@echo "✅ Test completati"

# Build
build:
	@echo "🏗️  Compilazione..."
	@cd backend && go build -o bin/thothix-backend .
	@echo "✅ Compilazione completata"

# Clean
clean:
	@echo "🧹 Pulizia..."
	@cd backend && rm -rf bin/ tmp/
	@echo "✅ Pulizia completata"

# Sviluppo
dev:
	@echo "🚀 Avvio in modalità sviluppo..."
	@cd backend && go run main.go

# Installazione dipendenze
install:
	@echo "📦 Installazione strumenti di sviluppo..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Installazione completata"

# Pre-commit (unified script)
pre-commit:
	@echo "🔧 Eseguendo pre-commit con script unificato..."
	@.\scripts\dev.bat pre-commit
	@echo "🎉 Pre-commit completato con successo!"

# Commit con pre-checks automatici
commit: pre-commit
	@echo "🚀 Pronto per il commit!"
	@echo "Usa: git commit -m 'messaggio'"
