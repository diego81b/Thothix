@echo off
setlocal enabledelayedexpansion

if "%1"=="" (
    echo Usage: %0 [dev^|staging^|prod] [command]
    echo.
    echo Environments:
    echo   dev      - Development environment ^(.env^)
    echo   staging  - Staging environment ^(.env.staging^)
    echo   prod     - Production environment ^(.env.prod^)
    echo.
    echo Commands:
    echo   up       - Start services
    echo   down     - Stop services
    echo   logs     - Show logs
    echo   status   - Show container status
    exit /b 1
)

set ENV=%1
set CMD=%2

if "%CMD%"=="" set CMD=up

REM Set environment file
if "%ENV%"=="dev" (
    set ENV_FILE=.env
    set COMPOSE_FILES=docker-compose.yml
) else if "%ENV%"=="staging" (
    set ENV_FILE=.env.staging
    set COMPOSE_FILES=docker-compose.yml
) else if "%ENV%"=="prod" (
    set ENV_FILE=.env.prod
    set COMPOSE_FILES=docker-compose.yml docker-compose.prod.yml
) else (
    echo Error: Unknown environment '%ENV%'
    exit /b 1
)

echo 🌍 Environment: %ENV%
echo 📄 Config: %ENV_FILE%
echo.

REM Execute command
if "%CMD%"=="up" (
    echo 🚀 Starting Thothix in %ENV% environment...
    docker-compose --env-file %ENV_FILE% up -d --build
) else if "%CMD%"=="down" (
    echo 🛑 Stopping Thothix in %ENV% environment...
    docker-compose --env-file %ENV_FILE% down
) else if "%CMD%"=="logs" (
    echo 📋 Showing logs for %ENV% environment...
    docker-compose --env-file %ENV_FILE% logs -f
) else if "%CMD%"=="status" (
    echo 📊 Container status for %ENV% environment...
    docker-compose --env-file %ENV_FILE% ps
) else (
    echo Error: Unknown command '%CMD%'
    exit /b 1
)
