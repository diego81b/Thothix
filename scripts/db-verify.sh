#!/bin/bash

# Database Verification Utility Script
# This script provides easy commands to verify database schema alignment

DB_NAME="thothix-db"
CONTAINER_NAME="postgres"

echo "=== Thothix Database Verification Utility ==="
echo ""

# Function to execute SQL commands
execute_sql() {
    docker-compose exec $CONTAINER_NAME psql -U postgres -d $DB_NAME -c "$1"
}

# Function to execute PostgreSQL commands
execute_psql() {
    docker-compose exec $CONTAINER_NAME psql -U postgres -d $DB_NAME -c "$1"
}

case "$1" in
    "check-basemodel")
        echo "Checking BaseModel columns alignment (should be 5 for all tables):"
        execute_sql "
        SELECT table_name, COUNT(*) as basemodel_columns 
        FROM information_schema.columns 
        WHERE table_schema = 'public' 
          AND column_name IN ('id', 'created_by', 'created_at', 'updated_by', 'updated_at') 
        GROUP BY table_name 
        ORDER BY table_name;"
        ;;
    
    "list-tables")
        echo "Listing all tables in database:"
        execute_psql "\d"
        ;;
    
    "check-table")
        if [ -z "$2" ]; then
            echo "Usage: $0 check-table <table_name>"
            exit 1
        fi
        echo "Checking structure of table: $2"
        execute_psql "\d $2"
        ;;
    
    "missing-field")
        if [ -z "$2" ]; then
            echo "Usage: $0 missing-field <field_name>"
            exit 1
        fi
        echo "Tables missing field '$2':"
        execute_sql "
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
          AND table_name NOT IN (
            SELECT table_name 
            FROM information_schema.columns 
            WHERE column_name = '$2' AND table_schema = 'public'
          )
        ORDER BY table_name;"
        ;;
    
    "has-field")
        if [ -z "$2" ]; then
            echo "Usage: $0 has-field <field_name>"
            exit 1
        fi
        echo "Tables that have field '$2':"
        execute_sql "
        SELECT table_name 
        FROM information_schema.columns 
        WHERE table_schema = 'public' AND column_name = '$2' 
        GROUP BY table_name
        ORDER BY table_name;"
        ;;
    
    "connect")
        echo "Connecting to database (use \q to exit):"
        docker-compose exec $CONTAINER_NAME psql -U postgres -d $DB_NAME
        ;;
    
    "status")
        echo "Database connection status:"
        execute_sql "SELECT version();"
        echo ""
        echo "Database size:"
        execute_sql "SELECT pg_size_pretty(pg_database_size('$DB_NAME')) as database_size;"
        ;;
    
    *)
        echo "Usage: $0 {check-basemodel|list-tables|check-table|missing-field|has-field|connect|status}"
        echo ""
        echo "Commands:"
        echo "  check-basemodel           Check if all tables have BaseModel columns"
        echo "  list-tables               List all tables in database"
        echo "  check-table <name>        Show structure of specific table"
        echo "  missing-field <field>     Show tables missing a specific field"
        echo "  has-field <field>         Show tables that have a specific field"
        echo "  connect                   Connect to database interactively"
        echo "  status                    Show database connection and size info"
        echo ""
        echo "Examples:"
        echo "  $0 check-basemodel"
        echo "  $0 check-table users"
        echo "  $0 missing-field updated_by"
        echo "  $0 has-field system_role"
        ;;
esac
