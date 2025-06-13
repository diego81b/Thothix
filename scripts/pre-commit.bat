@echo off
REM Pre-commit script per Windows
REM Automatizza formattazione e linting

echo 🔧 Eseguendo pre-commit checks...

REM Vai nella directory backend
cd /d "%~dp0\..\backend"
if errorlevel 1 (
    echo ❌ Impossibile accedere alla directory backend
    exit /b 1
)

REM Formatta il codice
echo 📝 Formattazione del codice...
REM Using basic gofmt only to avoid import formatting conflicts
gofmt -w .

REM Aggiungi automaticamente i file formattati
echo 📋 Aggiungendo file formattati...
cd /d "%~dp0\.."
git add backend/
cd /d "%~dp0\..\backend"

REM Esegui linting
echo 🔍 Eseguendo golangci-lint...
golangci-lint run --timeout=3m
if errorlevel 1 (
    echo ❌ Linting fallito! Il commit è stato bloccato.
    echo Risolvi gli errori e riprova.
    exit /b 1
)

REM Esegui i test
echo 🧪 Eseguendo i test...
go test ./...
if errorlevel 1 (
    echo ❌ Test falliti! Il commit è stato bloccato.
    echo Risolvi i test e riprova.
    exit /b 1
)

REM Torna alla directory principale e aggiungi i file
cd /d "%~dp0\.."
git add backend/

echo ✅ Pre-commit checks completati con successo!
echo 🚀 Procedendo con il commit...
