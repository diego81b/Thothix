#!/bin/bash
# Universal Go development script for Thothix
# Usage: ./dev.sh [format|lint|pre-commit|all]

set -e

ACTION=${1:-all}

echo "🔧 Thothix Development Script - Action: $ACTION"

# Change to backend directory
cd "$(dirname "$0")/../backend"

format() {
    echo "📝 Formatting Go code..."
    # gofumpt includes import organization and stricter formatting
    gofumpt -w .
    echo "✅ Formatting completed"
}

lint() {
    echo "🔍 Running golangci-lint..."
    golangci-lint run --timeout=3m
    echo "✅ Linting passed"
}

pre_commit() {
    echo "🚀 Running pre-commit checks..."
    format
    echo "📋 Adding formatted files to git..."
    cd "$(dirname "$0")/.."
    git add backend/
    cd "$(dirname "$0")/../backend"
    lint
    echo "🧪 Running tests..."
    go test ./...
    echo "✅ Pre-commit checks completed"
}

case $ACTION in
    format)
        format
        ;;
    lint)
        lint
        ;;
    pre-commit)
        pre_commit
        ;;
    all)
        pre_commit
        ;;
    *)
        echo "❌ Invalid action. Use: format, lint, pre-commit, or all"
        exit 1
        ;;
esac

echo "🎉 Script completed successfully!"
