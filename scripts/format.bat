@echo off
REM Script per formattare tutto il codice Go nel progetto Thothix
REM Uso: scripts\format.bat

echo 🔧 Formattazione del codice Go in corso...

REM Cambia nella directory backend
cd /d "%~dp0\..\backend"

echo 📁 Directory corrente: %CD%

REM 1. gofmt - Formattazione base Go
echo 🎨 Eseguendo gofmt...
gofmt -w . 2>nul
if %errorlevel% equ 0 (
    echo ✅ gofmt completato
) else (
    echo ❌ gofmt fallito o non trovato
)

REM 2. goimports - Gestione import + formattazione
echo 📦 Eseguendo goimports...
goimports -w . 2>nul
if %errorlevel% equ 0 (
    echo ✅ goimports completato
) else (
    echo ⚠️  goimports non trovato, installazione in corso...
    go install golang.org/x/tools/cmd/goimports@latest
    goimports -w .
    if %errorlevel% equ 0 (
        echo ✅ goimports installato e completato
    ) else (
        echo ❌ Errore nell'installazione di goimports
    )
)

REM 3. gofumpt - Formattazione più rigida
echo 🎯 Eseguendo gofumpt...
gofumpt -w . 2>nul
if %errorlevel% equ 0 (
    echo ✅ gofumpt completato
) else (
    echo ⚠️  gofumpt non trovato, installazione in corso...
    go install mvdan.cc/gofumpt@latest
    gofumpt -w .
    if %errorlevel% equ 0 (
        echo ✅ gofumpt installato e completato
    ) else (
        echo ❌ Errore nell'installazione di gofumpt
    )
)

echo.
echo 🎉 Formattazione completata con successo!
echo 💡 Ora puoi eseguire 'golangci-lint run' per verificare la qualità del codice

pause
