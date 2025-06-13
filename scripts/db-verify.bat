@echo off
REM Database Verification Utility Script for Windows
REM This script provides easy commands to verify database schema alignment

set DB_NAME=thothix-db
set CONTAINER_NAME=postgres

echo === Thothix Database Verification Utility ===
echo.

if "%1"=="check-basemodel" goto check_basemodel
if "%1"=="list-tables" goto list_tables
if "%1"=="check-table" goto check_table
if "%1"=="missing-field" goto missing_field
if "%1"=="has-field" goto has_field
if "%1"=="connect" goto connect
if "%1"=="status" goto status
goto usage

:check_basemodel
echo Checking BaseModel columns alignment (should be 5 for all tables):
docker-compose exec %CONTAINER_NAME% psql -U postgres -d %DB_NAME% -c "SELECT table_name, COUNT(*) as basemodel_columns FROM information_schema.columns WHERE table_schema = 'public' AND column_name IN ('id', 'created_by', 'created_at', 'updated_by', 'updated_at') GROUP BY table_name ORDER BY table_name;"
goto end

:list_tables
echo Listing all tables in database:
docker-compose exec %CONTAINER_NAME% psql -U postgres -d %DB_NAME% -c "\d"
goto end

:check_table
if "%2"=="" (
    echo Usage: %0 check-table ^<table_name^>
    goto end
)
echo Checking structure of table: %2
docker-compose exec %CONTAINER_NAME% psql -U postgres -d %DB_NAME% -c "\d %2"
goto end

:missing_field
if "%2"=="" (
    echo Usage: %0 missing-field ^<field_name^>
    goto end
)
echo Tables missing field '%2':
docker-compose exec %CONTAINER_NAME% psql -U postgres -d %DB_NAME% -c "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' AND table_name NOT IN (SELECT table_name FROM information_schema.columns WHERE column_name = '%2' AND table_schema = 'public') ORDER BY table_name;"
goto end

:has_field
if "%2"=="" (
    echo Usage: %0 has-field ^<field_name^>
    goto end
)
echo Tables that have field '%2':
docker-compose exec %CONTAINER_NAME% psql -U postgres -d %DB_NAME% -c "SELECT table_name FROM information_schema.columns WHERE table_schema = 'public' AND column_name = '%2' GROUP BY table_name ORDER BY table_name;"
goto end

:connect
echo Connecting to database (use \q to exit):
docker-compose exec %CONTAINER_NAME% psql -U postgres -d %DB_NAME%
goto end

:status
echo Database connection status:
docker-compose exec %CONTAINER_NAME% psql -U postgres -d %DB_NAME% -c "SELECT version();"
echo.
echo Database size:
docker-compose exec %CONTAINER_NAME% psql -U postgres -d %DB_NAME% -c "SELECT pg_size_pretty(pg_database_size('%DB_NAME%')) as database_size;"
goto end

:usage
echo Usage: %0 {check-basemodel^|list-tables^|check-table^|missing-field^|has-field^|connect^|status}
echo.
echo Commands:
echo   check-basemodel           Check if all tables have BaseModel columns
echo   list-tables               List all tables in database
echo   check-table ^<name^>        Show structure of specific table
echo   missing-field ^<field^>     Show tables missing a specific field
echo   has-field ^<field^>         Show tables that have a specific field
echo   connect                   Connect to database interactively
echo   status                    Show database connection and size info
echo.
echo Examples:
echo   %0 check-basemodel
echo   %0 check-table users
echo   %0 missing-field updated_by
echo   %0 has-field system_role

:end
