#!/bin/bash
# Database Verification Utility Script for Unix/Linux/macOS
# This script provides easy commands to verify database schema alignment

set -e

DB_NAME="thothix-db"
CONTAINER_NAME="postgres"

echo "=== Thothix Database Verification Utility ==="
echo ""

case "$1" in
    "check-basemodel")
        echo "Checking BaseModel columns alignment (should be 5 for all tables):"
        docker compose exec "$CONTAINER_NAME" psql -U postgres -d "$DB_NAME" -c "SELECT table_name, COUNT(*) as basemodel_columns FROM information_schema.columns WHERE table_schema = 'public' AND column_name IN ('id', 'created_by', 'created_at', 'updated_by', 'updated_at') GROUP BY table_name ORDER BY table_name;"
        ;;
    "list-tables")
        echo "Listing all tables in database:"
        docker compose exec "$CONTAINER_NAME" psql -U postgres -d "$DB_NAME" -c "\\d"
        ;;
    "check-table")
        if [ -z "$2" ]; then
            echo "Usage: $0 check-table <table_name>"
            exit 1
        fi
        echo "Checking table structure for: $2"
        docker compose exec "$CONTAINER_NAME" psql -U postgres -d "$DB_NAME" -c "\\d $2"
        ;;
    "missing-field")
        if [ -z "$2" ] || [ -z "$3" ]; then
            echo "Usage: $0 missing-field <table_name> <field_name>"
            exit 1
        fi
        echo "Checking if field '$3' exists in table '$2':"
        docker compose exec "$CONTAINER_NAME" psql -U postgres -d "$DB_NAME" -c "SELECT column_name FROM information_schema.columns WHERE table_name = '$2' AND column_name = '$3';"
        ;;
    "has-field")
        if [ -z "$2" ] || [ -z "$3" ]; then
            echo "Usage: $0 has-field <table_name> <field_name>"
            exit 1
        fi
        echo "Checking field '$3' in table '$2':"
        docker compose exec "$CONTAINER_NAME" psql -U postgres -d "$DB_NAME" -c "SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '$2' AND column_name = '$3';"
        ;;
    "connect")
        echo "Connecting to PostgreSQL database..."
        docker compose exec "$CONTAINER_NAME" psql -U postgres -d "$DB_NAME"
        ;;
    "status")
        echo "Database container status:"
        docker compose ps "$CONTAINER_NAME"
        echo ""
        echo "Database connection test:"
        docker compose exec "$CONTAINER_NAME" pg_isready -U postgres -d "$DB_NAME"
        ;;
    *)
        echo "Usage: $0 <command> [arguments]"
        echo ""
        echo "Available commands:"
        echo "  check-basemodel           - Verify BaseModel columns (id, created_by, created_at, updated_by, updated_at)"
        echo "  list-tables              - List all tables in the database"
        echo "  check-table <table>      - Show structure of specific table"
        echo "  missing-field <table> <field> - Check if field exists in table"
        echo "  has-field <table> <field>     - Show field details in table"
        echo "  connect                  - Connect to PostgreSQL interactively"
        echo "  status                   - Show database container status"
        echo ""
        echo "Examples:"
        echo "  $0 check-basemodel"
        echo "  $0 list-tables"
        echo "  $0 check-table users"
        echo "  $0 has-field users email"
        echo "  $0 connect"
        exit 1
        ;;
esac
