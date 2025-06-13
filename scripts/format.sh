#!/bin/bash

# Script per formattare tutto il codice Go nel progetto Thothix
# Uso: ./scripts/format.sh

set -e

echo "🔧 Formattazione del codice Go in corso..."

# Cambia nella directory backend
cd "$(dirname "$0")/../backend"

echo "📁 Directory corrente: $(pwd)"

# 1. gofmt - Formattazione base Go
echo "🎨 Eseguendo gofmt..."
if command -v gofmt &> /dev/null; then
    gofmt -w .
    echo "✅ gofmt completato"
else
    echo "❌ gofmt non trovato"
fi

# 2. goimports - Gestione import + formattazione
echo "📦 Eseguendo goimports..."
if command -v goimports &> /dev/null; then
    goimports -w .
    echo "✅ goimports completato"
else
    echo "⚠️  goimports non trovato, installazione in corso..."
    go install golang.org/x/tools/cmd/goimports@latest
    goimports -w .
    echo "✅ goimports installato e completato"
fi

# 3. gofumpt - Formattazione più rigida (opzionale)
echo "🎯 Eseguendo gofumpt..."
if command -v gofumpt &> /dev/null; then
    gofumpt -w .
    echo "✅ gofumpt completato"
else
    echo "⚠️  gofumpt non trovato, installazione in corso..."
    go install mvdan.cc/gofumpt@latest
    gofumpt -w .
    echo "✅ gofumpt installato e completato"
fi

echo ""
echo "🎉 Formattazione completata con successo!"
echo "💡 Ora puoi eseguire 'golangci-lint run' per verificare la qualità del codice"
