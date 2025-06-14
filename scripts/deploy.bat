@echo off
setlocal enabledelayedexpansion

if "%1"=="" (
    echo Usage: %0 [dev^|staging^|prod] [command] [options]
    echo.
    echo Environments:
    echo   dev      - Development environment ^(.env^)
    echo   staging  - Staging environment ^(.env.staging^)
    echo   prod     - Production environment ^(.env.prod^) with Vault
    echo.
    echo Commands:
    echo   up       - Start services
    echo   down     - Stop services
    echo   logs     - Show logs
    echo   status   - Show container status
    echo   vault    - Vault-specific commands ^(init, ui, status^)
    echo.
    echo Options:
    echo   --vault  - Force enable Vault for dev/staging
    exit /b 1
)

set ENV=%1
set CMD=%2
set OPT=%3

if "%CMD%"=="" set CMD=up

REM Set environment file and compose files
if "%ENV%"=="dev" (
    set ENV_FILE=.env
    if "%OPT%"=="--vault" (
        set COMPOSE_FILES=-f docker-compose.yml -f docker-compose.vault.yml
        set ENV_FILE=.env.vault
    ) else (
        set COMPOSE_FILES=-f docker-compose.yml
    )
) else if "%ENV%"=="staging" (
    set ENV_FILE=.env.staging
    if "%OPT%"=="--vault" (
        set COMPOSE_FILES=-f docker-compose.yml -f docker-compose.vault.yml
    ) else (
        set COMPOSE_FILES=-f docker-compose.yml
    )
) else if "%ENV%"=="prod" (
    set ENV_FILE=.env.prod
    set COMPOSE_FILES=-f docker-compose.yml -f docker-compose.prod.yml
) else if "%ENV%"=="prod" (
    set ENV_FILE=.env.prod
    set COMPOSE_FILES=docker-compose.yml docker-compose.prod.yml
) else (
    echo Error: Unknown environment '%ENV%'
    exit /b 1
)

echo 🌍 Environment: %ENV%
echo 📄 Config: %ENV_FILE%
echo 🐳 Compose: %COMPOSE_FILES%
echo.

REM Execute command
if "%CMD%"=="vault" (
    if "%OPT%"=="init" (
        echo 🔐 Initializing Vault...
        docker-compose %COMPOSE_FILES% --env-file %ENV_FILE% up vault-init
    ) else if "%OPT%"=="ui" (
        echo 🌐 Opening Vault UI...
        start http://localhost:8200
    ) else if "%OPT%"=="status" (
        echo 📊 Vault status...
        docker-compose %COMPOSE_FILES% --env-file %ENV_FILE% exec vault vault status
    ) else (
        echo Available vault commands: init, ui, status
    )
) else if "%CMD%"=="up" (
    echo 🚀 Starting Thothix in %ENV% environment...
    docker-compose %COMPOSE_FILES% --env-file %ENV_FILE% up -d --build
) else if "%CMD%"=="down" (
    echo 🛑 Stopping Thothix in %ENV% environment...
    docker-compose %COMPOSE_FILES% --env-file %ENV_FILE% down
) else if "%CMD%"=="logs" (
    echo 📋 Showing logs for %ENV% environment...
    docker-compose %COMPOSE_FILES% --env-file %ENV_FILE% logs -f
) else if "%CMD%"=="status" (
    echo 📊 Container status for %ENV% environment...
    docker-compose %COMPOSE_FILES% --env-file %ENV_FILE% ps
) else (
    echo Error: Unknown command '%CMD%'
    exit /b 1
)
