#!/bin/bash
# Thothix Deployment Script for Unix/Linux/macOS
# Usage: ./deploy.sh [dev|staging|prod] [command] [options]

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 [dev|staging|prod] [command] [options]"
    echo ""
    echo "Environments:"
    echo "  dev      - Development environment (.env)"
    echo "  staging  - Staging environment (.env.staging)"
    echo "  prod     - Production environment (.env.prod) with Vault"
    echo ""
    echo "Commands:"
    echo "  up       - Start services"
    echo "  down     - Stop services"
    echo "  logs     - Show logs"
    echo "  status   - Show container status"
    echo "  vault    - Vault-specific commands (init, ui, status)"
    echo ""
    echo "Note: Vault is now integrated in all environments"
    exit 1
fi

ENV="$1"
CMD="${2:-up}"
OPT="$3"

# Navigate to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Set environment file and compose files
case "$ENV" in
    "dev")
        ENV_FILE=".env"
        COMPOSE_FILES="-f docker-compose.yml"
        ;;
    "staging")
        ENV_FILE=".env.staging"
        COMPOSE_FILES="-f docker-compose.yml -f docker-compose.staging.yml"
        ;;
    "prod")
        ENV_FILE=".env.prod"
        COMPOSE_FILES="-f docker-compose.yml -f docker-compose.prod.yml"
        ;;
    *)
        echo "❌ Invalid environment: $ENV. Use dev, staging, or prod"
        exit 1
        ;;
esac

# Check if environment file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "❌ Environment file $ENV_FILE not found!"
    echo "📋 Please copy .env.example to $ENV_FILE and configure it"
    exit 1
fi

echo "🚀 Thothix Deployment - Environment: $ENV, Command: $CMD"

case "$CMD" in
    "up")
        echo "📦 Starting $ENV environment..."
        docker compose $COMPOSE_FILES --env-file="$ENV_FILE" up -d --build
        echo "✅ $ENV environment started successfully"
        echo "🔍 Container status:"
        docker compose $COMPOSE_FILES ps
        ;;
    "down")
        echo "🛑 Stopping $ENV environment..."
        docker compose $COMPOSE_FILES --env-file="$ENV_FILE" down
        echo "✅ $ENV environment stopped"
        ;;
    "logs")
        if [ -n "$OPT" ]; then
            echo "📋 Showing logs for service: $OPT"
            docker compose $COMPOSE_FILES --env-file="$ENV_FILE" logs -f "$OPT"
        else
            echo "📋 Showing all logs..."
            docker compose $COMPOSE_FILES --env-file="$ENV_FILE" logs -f
        fi
        ;;
    "status")
        echo "📊 Container status for $ENV environment:"
        docker compose $COMPOSE_FILES ps
        echo ""
        echo "🔍 Resource usage:"
        docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
        ;;
    "vault")
        case "$OPT" in
            "init")
                echo "🔐 Initializing Vault..."
                docker compose $COMPOSE_FILES --env-file="$ENV_FILE" exec vault vault operator init
                ;;
            "ui")
                echo "🌐 Opening Vault UI..."
                VAULT_ADDR=$(grep VAULT_ADDR "$ENV_FILE" | cut -d'=' -f2 | tr -d '"')
                echo "Vault UI available at: $VAULT_ADDR/ui"
                ;;
            "status")
                echo "🔍 Vault status:"
                docker compose $COMPOSE_FILES --env-file="$ENV_FILE" exec vault vault status
                ;;
            *)
                echo "❌ Invalid vault command. Use: init, ui, status"
                exit 1
                ;;
        esac
        ;;
    *)
        echo "❌ Invalid command: $CMD. Use: up, down, logs, status, vault"
        exit 1
        ;;
esac

echo "🎉 Operation completed successfully!"
