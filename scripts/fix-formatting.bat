@echo off
REM Script per correggere formattazione Go automaticamente
REM Uso: .\fix-formatting.bat

echo 🔧 Correzione automatica formattazione Go...

cd /d "%~dp0\..\backend"

echo 📝 Applicando goimports...
goimports -w .

echo 🎯 Applicando gofumpt...
gofumpt -w .

echo 🔍 Verificando con golangci-lint...
golangci-lint run

if errorlevel 1 (
    echo ⚠️  Alcuni errori di linting rilevati. Controllare manualmente.
) else (
    echo ✅ Formattazione corretta! Tutti i file sono conformi.
)

echo 📋 Aggiungendo file corretti a git...
cd /d "%~dp0\.."
git add backend/

echo 🎉 Formattazione completata!
