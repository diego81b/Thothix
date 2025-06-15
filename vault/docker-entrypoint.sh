#!/bin/bash
set -e

# Script di entry point personalizzato per Thothix Vault
echo "🔐 Starting Thothix Vault..."

# Se è in modalità dev, avvia Vault in dev mode manualmente
if [ "${VAULT_DEV_MODE}" = "true" ]; then
    echo "🔧 Running in development mode"
    echo "   Root Token: ${VAULT_DEV_ROOT_TOKEN_ID}"
    echo "   Address: http://0.0.0.0:8200"

    # Avvia Vault in dev mode direttamente senza setcap issues
    exec vault server \
        -dev \
        -dev-root-token-id="${VAULT_DEV_ROOT_TOKEN_ID}" \
        -dev-listen-address="0.0.0.0:8200" \
        "$@"
fi

# Se non è in modalità dev, usa configurazione personalizzata
echo "🏢 Running in production mode"

# Verifica che la configurazione esista
if [ ! -f "/vault/config/vault.hcl" ]; then
    echo "❌ Configuration file /vault/config/vault.hcl not found!"
    exit 1
fi

# Avvia Vault con la configurazione personalizzata
echo "🚀 Starting Vault server with custom configuration..."
exec vault server -config=/vault/config/vault.hcl "$@"
