#!/bin/bash
# Initialize HashiCorp Vault with Thothix secrets

set -e

# Configuration
VAULT_ADDR=${VAULT_ADDR:-http://vault:8200}
VAULT_MOUNT=${VAULT_MOUNT:-thothix}
MAX_RETRIES=30
RETRY_INTERVAL=5

echo "🔐 Initializing Thothix secrets in Vault..."
echo "   Vault Address: $VAULT_ADDR"
echo "   Mount Path: $VAULT_MOUNT"

# Wait for Vault to be ready
echo "🔐 Waiting for Vault to be ready..."
for i in $(seq 1 $MAX_RETRIES); do
  if vault status >/dev/null 2>&1; then
    echo "✅ Vault is ready!"
    break
  fi
  echo "   Attempt $i/$MAX_RETRIES - Vault not ready, waiting $RETRY_INTERVAL seconds..."
  sleep $RETRY_INTERVAL

  if [ $i -eq $MAX_RETRIES ]; then
    echo "❌ Vault failed to become ready after $MAX_RETRIES attempts"
    exit 1
  fi
done

# Check if secrets engine already exists
if vault secrets list | grep -q "^$VAULT_MOUNT/"; then
  echo "✅ KV secrets engine '$VAULT_MOUNT/' already exists"
else
  echo "🔐 Enabling KV secrets engine at '$VAULT_MOUNT/'..."
  vault secrets enable -path=$VAULT_MOUNT kv-v2
fi

# Environment-specific secrets
ENVIRONMENT=${ENVIRONMENT:-development}

echo "🔐 Creating database secrets for $ENVIRONMENT..."
vault kv put $VAULT_MOUNT/database \
  host="postgres" \
  port="5432" \
  username="postgres" \
  password="${DB_PASSWORD:-secure_db_password_$(date +%s)}" \
  database="thothix-${ENVIRONMENT}"

echo "🔐 Creating Clerk secrets for $ENVIRONMENT..."
if [ "$ENVIRONMENT" = "production" ]; then
  vault kv put $VAULT_MOUNT/clerk \
    secret_key="${CLERK_SECRET_KEY:-sk_live_change_me_in_production}" \
    webhook_secret="${CLERK_WEBHOOK_SECRET:-whsec_change_me_in_production}" \
    publishable_key="${CLERK_PUBLISHABLE_KEY:-pk_live_change_me_in_production}"
else
  vault kv put $VAULT_MOUNT/clerk \
    secret_key="${CLERK_SECRET_KEY:-sk_test_development_key}" \
    webhook_secret="${CLERK_WEBHOOK_SECRET:-whsec_development_secret}" \
    publishable_key="${CLERK_PUBLISHABLE_KEY:-pk_test_development_key}"
fi

echo "🔐 Creating application secrets for $ENVIRONMENT..."
vault kv put $VAULT_MOUNT/app \
  jwt_secret="${JWT_SECRET:-jwt_secret_$(openssl rand -hex 32)}" \
  encryption_key="${ENCRYPTION_KEY:-$(openssl rand -hex 32)}" \
  environment="$ENVIRONMENT" \
  debug_mode="${DEBUG_MODE:-false}"

echo "🔐 Creating policies..."
vault policy write thothix-app - <<EOF
# Read-only access to application secrets
path "$VAULT_MOUNT/data/*" {
  capabilities = ["read"]
}

# Allow token renewal
path "auth/token/renew-self" {
  capabilities = ["update"]
}

# Allow token lookup
path "auth/token/lookup-self" {
  capabilities = ["read"]
}
EOF

echo "🔐 Creating read-only policy for monitoring..."
vault policy write thothix-readonly - <<EOF
# Read-only access for monitoring/debugging
path "$VAULT_MOUNT/data/*" {
  capabilities = ["read"]
}
EOF

echo "🔐 Creating application tokens..."
echo "Creating app token with thothix-app policy..."
APP_TOKEN=$(vault token create \
  -policy=thothix-app \
  -ttl=8760h \
  -renewable=true \
  -display-name="thothix-app-token" \
  -format=json | jq -r .auth.client_token)

echo "Creating readonly token for monitoring..."
READONLY_TOKEN=$(vault token create \
  -policy=thothix-readonly \
  -ttl=168h \
  -renewable=true \
  -display-name="thothix-readonly-token" \
  -format=json | jq -r .auth.client_token)

echo ""
echo "✅ Vault initialized successfully!"
echo ""
echo "📋 Configuration Summary:"
echo "   Environment: $ENVIRONMENT"
echo "   Mount Path: $VAULT_MOUNT"
echo "   Vault UI: http://localhost:8200"
echo ""
echo "🔑 Tokens created:"
echo "   App Token: $APP_TOKEN"
echo "   Readonly Token: $READONLY_TOKEN"
echo ""
echo "⚠️  IMPORTANT: Save these tokens securely!"
echo "   Add VAULT_APP_TOKEN=$APP_TOKEN to your .env file"
echo ""
echo "🔍 Test secret retrieval:"
echo "   vault kv get $VAULT_MOUNT/database"
echo "   vault kv get $VAULT_MOUNT/clerk"
echo "   vault kv get $VAULT_MOUNT/app"
