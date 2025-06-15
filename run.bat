@echo off
REM Thothix Universal Runner - Windows
REM This script automatically installs Node.js dependencies and runs npm scripts

echo 🔧 Thothix Universal Runner (Windows)

REM Check if Node.js is installed
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Node.js is not installed. Please install Node.js 16+ from https://nodejs.org/
    exit /b 1
)

REM Check if we have arguments
if "%~1"=="" (
    echo 📖 Usage: run.bat [script] [args...]
    echo 📋 Available scripts:
    npm run --silent 2>nul
    if %errorlevel% neq 0 (
        echo   Installing dependencies first...
        npm install
        npm run --silent
    )
    exit /b 0
)

REM Auto-install dependencies if node_modules doesn't exist
if not exist "node_modules" (
    echo 📦 Installing Node.js dependencies...
    npm install
    if %errorlevel% neq 0 (
        echo ❌ Failed to install dependencies
        exit /b 1
    )
    echo ✅ Dependencies installed
)

REM Run the npm script with all arguments
echo 🚀 Running: npm run %*
npm run %*
