@echo off
REM Universal Go development script for Thothix
REM Usage: .\dev.bat [format|lint|pre-commit|all]

set ACTION=%1
if "%ACTION%"=="" set ACTION=all

echo 🔧 Thothix Development Script - Action: %ACTION%

cd /d "%~dp0\..\backend"
if errorlevel 1 (
    echo ❌ Cannot access backend directory
    exit /b 1
)

if "%ACTION%"=="format" goto FORMAT
if "%ACTION%"=="lint" goto LINT
if "%ACTION%"=="pre-commit" goto PRECOMMIT
if "%ACTION%"=="all" goto ALL

echo ❌ Invalid action. Use: format, lint, pre-commit, or all
exit /b 1

:FORMAT
echo 📝 Formatting Go code...
REM Using basic gofmt only to avoid import formatting issues
gofmt -w .
echo ✅ Formatting completed
goto END

:LINT
echo 🔍 Running golangci-lint...
golangci-lint run --timeout=3m
if errorlevel 1 (
    echo ❌ Linting failed
    exit /b 1
)
echo ✅ Linting passed
goto END

:PRECOMMIT
echo 🚀 Running pre-commit checks...
call :FORMAT
echo 📋 Adding formatted files to git...
cd /d "%~dp0\.."
git add backend/
cd /d "%~dp0\..\backend"
call :LINT
echo 🧪 Running tests...
go test ./...
if errorlevel 1 (
    echo ❌ Tests failed
    exit /b 1
)
echo ✅ Pre-commit checks completed
goto END

:ALL
call :PRECOMMIT
goto END

:END
echo 🎉 Script completed successfully!
