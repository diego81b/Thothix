# Makefile per il progetto Thothix
# Uso: make <target>

.PHONY: help format lint test build clean dev

# Default target
help:
	@echo "🚀 Comandi disponibili per Thothix:"
	@echo ""
	@echo "  format     - Formatta tutto il codice Go"
	@echo "  lint       - Esegue golangci-lint"
	@echo "  test       - Esegue i test"
	@echo "  build      - Compila l'applicazione"
	@echo "  clean      - Pulisce i file temporanei"
	@echo "  dev        - Avvia in modalità sviluppo"
	@echo "  install    - Installa le dipendenze"

# Formattazione del codice
format:
	@echo "🔧 Formattazione del codice..."
	@cd backend && gofmt -w .
	@echo "✅ Formattazione completata"

# Linting
lint:
	@echo "🔍 Eseguendo golangci-lint..."
	@cd backend && golangci-lint run --timeout=3m
	@echo "✅ Linting completato"

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

# Pre-commit (formatta + lint + test)
pre-commit: format lint test
	@echo "🎉 Pre-commit completato con successo!"
	@echo "📝 Aggiungendo file modificati a git..."
	@git add backend/

# Setup git hooks
setup-hooks:
	@echo "🔧 Configurando Git hooks..."
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Git hooks configurati!"

# Commit con pre-checks automatici
commit: pre-commit
	@echo "🚀 Pronto per il commit!"
	@echo "Usa: git commit -m 'messaggio'"
