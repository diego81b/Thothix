#!/usr/bin/env zx

// Thothix Database Verification Script - Cross-platform with Zx
// Usage: zx scripts/db-verify.mjs [command] [args...]

import { $, argv, echo, path } from 'zx'

// Configure for Windows compatibility
if (process.platform === 'win32') {
  $.shell = 'cmd.exe'
  $.prefix = ''
}

$.verbose = true

const command = argv._[0] || 'help'
const arg1 = argv._[1]
const arg2 = argv._[2]

const DB_NAME = 'thothix-db'
const CONTAINER_NAME = 'postgres'

echo`=== Thothix Database Verification Utility ===`
echo``

// Helper function for database queries
const dbQuery = async (query) => {
  return $`docker compose exec ${CONTAINER_NAME} psql -U postgres -d ${DB_NAME} -c ${query}`
}

const dbCommand = async (pgCommand) => {
  return $`docker compose exec ${CONTAINER_NAME} psql -U postgres -d ${DB_NAME} ${pgCommand}`
}

// Main command execution
switch (command) {
  case 'check-basemodel':
    echo`Checking BaseModel columns alignment (should be 5 for all tables):`
    await dbQuery(`"SELECT table_name, COUNT(*) as basemodel_columns FROM information_schema.columns WHERE table_schema = 'public' AND column_name IN ('id', 'created_by', 'created_at', 'updated_by', 'updated_at') GROUP BY table_name ORDER BY table_name;"`)
    break

  case 'list-tables':
    echo`Listing all tables in database:`
    await dbCommand(`-c "\\d"`)
    break

  case 'check-table':
    if (!arg1) {
      echo`Usage: zx scripts/db-verify.mjs check-table <table_name>`
      process.exit(1)
    }
    echo`Checking table structure for: ${arg1}`
    await dbCommand(`-c "\\d ${arg1}"`)
    break

  case 'missing-field':
    if (!arg1 || !arg2) {
      echo`Usage: zx scripts/db-verify.mjs missing-field <table_name> <field_name>`
      process.exit(1)
    }
    echo`Checking if field '${arg2}' exists in table '${arg1}':`
    await dbQuery(`"SELECT column_name FROM information_schema.columns WHERE table_name = '${arg1}' AND column_name = '${arg2}';"`)
    break

  case 'has-field':
    if (!arg1 || !arg2) {
      echo`Usage: zx scripts/db-verify.mjs has-field <table_name> <field_name>`
      process.exit(1)
    }
    echo`Checking field '${arg2}' in table '${arg1}':`
    await dbQuery(`"SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '${arg1}' AND column_name = '${arg2}';"`)
    break

  case 'connect':
    echo`Connecting to PostgreSQL database...`
    await $`docker compose exec ${CONTAINER_NAME} psql -U postgres -d ${DB_NAME}`
    break

  case 'status':
    echo`Database container status:`
    await $`docker compose ps ${CONTAINER_NAME}`
    echo``
    echo`Database connection test:`
    await $`docker compose exec ${CONTAINER_NAME} pg_isready -U postgres -d ${DB_NAME}`
    break

  case 'help':
  default:
    echo`Usage: zx scripts/db-verify.mjs <command> [arguments]

Available commands:
  check-basemodel           - Verify BaseModel columns (id, created_by, created_at, updated_by, updated_at)
  list-tables              - List all tables in the database
  check-table <table>      - Show structure of specific table
  missing-field <table> <field> - Check if field exists in table
  has-field <table> <field>     - Show field details in table
  connect                  - Connect to PostgreSQL interactively
  status                   - Show database container status

Examples:
  zx scripts/db-verify.mjs check-basemodel
  zx scripts/db-verify.mjs list-tables
  zx scripts/db-verify.mjs check-table users
  zx scripts/db-verify.mjs has-field users email
  zx scripts/db-verify.mjs connect

NPM shortcuts:
  npm run db:status        # Check database status
  npm run db:connect       # Connect interactively
  npm run db:tables        # List all tables
  npm run db:check         # Check BaseModel columns
`
    process.exit(command === 'help' ? 0 : 1)
}

echo`🎉 Database operation completed successfully!`
