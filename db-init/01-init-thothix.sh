#!/bin/bash
set -e

# Script di inizializzazione personalizzato per Thothix Database
echo "🚀 Inizializzazione database Thothix..."

# Crea estensioni utili
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Abilita estensioni utili
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pg_trgm";
    CREATE EXTENSION IF NOT EXISTS "unaccent";
    
    -- Crea indici per le ricerche full-text (opzionale)
    -- CREATE INDEX IF NOT EXISTS idx_messages_content_fts ON messages USING gin(to_tsvector('english', content));
    
    -- Configura timezone
    SET timezone = 'UTC';
    
    -- Log inizializzazione completata
    SELECT 'Database Thothix inizializzato con successo!' as status;
EOSQL

echo "✅ Inizializzazione database Thothix completata!"
