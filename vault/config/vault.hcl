# Production Vault Configuration
ui = true

# Storage backend
storage "file" {
  path = "/vault/data"
}

# Network listener
listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = true
  # In production, enable TLS:
  # tls_cert_file = "/vault/config/vault.crt"
  # tls_key_file = "/vault/config/vault.key"
}

# Cluster configuration (for HA)
cluster_addr = "http://vault:8201"
api_addr = "http://vault:8200"

# Logging
log_level = "INFO"
log_file = "/vault/logs/vault.log"
log_rotate_duration = "24h"
log_rotate_max_files = 30

# Lease configuration
default_lease_ttl = "168h"   # 1 week
max_lease_ttl = "720h"       # 30 days

# Disable memory lock for Docker
disable_mlock = true

# Telemetry (optional)
telemetry {
  prometheus_retention_time = "30s"
  disable_hostname = true
}
