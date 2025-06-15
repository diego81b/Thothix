@echo off
setlocal EnableDelayedExpansion

:: Script per bump automatico di versione
:: Uso: version-bump.bat [major|minor|patch] [optional-description]

if "%1"=="" (
    echo ❌ Tipo di bump richiesto: major, minor, o patch
    echo 📖 Uso: version-bump.bat [major^|minor^|patch] [optional-description]
    echo.
    echo 📋 Esempi:
    echo   version-bump.bat patch "Fix authentication bug"
    echo   version-bump.bat minor "Add user management feature"
    echo   version-bump.bat major "Breaking API changes"
    exit /b 1
)

set "BUMP_TYPE=%1"
set "DESCRIPTION=%2"

:: Verifica che il tipo di bump sia valido
if not "%BUMP_TYPE%"=="major" if not "%BUMP_TYPE%"=="minor" if not "%BUMP_TYPE%"=="patch" (
    echo ❌ Tipo di bump non valido: %BUMP_TYPE%
    echo ✅ Tipi validi: major, minor, patch
    exit /b 1
)

echo 🔍 Lettura versione corrente...

:: Legge la versione corrente dal CHANGELOG
for /f "tokens=2 delims= " %%i in ('findstr /r "^## v[0-9]*\.[0-9]*\.[0-9]*" CHANGELOG.md ^| head -1') do (
    set "CURRENT_VERSION=%%i"
    goto :version_found
)

:version_found
if "%CURRENT_VERSION%"=="" (
    echo ❌ Impossibile trovare la versione corrente nel CHANGELOG
    exit /b 1
)

:: Rimuove il prefisso 'v' dalla versione
set "VERSION_NUMBERS=%CURRENT_VERSION:v=%"

:: Estrae major, minor, patch
for /f "tokens=1,2,3 delims=." %%a in ("%VERSION_NUMBERS%") do (
    set "MAJOR=%%a"
    set "MINOR=%%b"
    set "PATCH=%%c"
)

:: Calcola la nuova versione
if "%BUMP_TYPE%"=="major" (
    set /a "MAJOR=%MAJOR%+1"
    set "MINOR=0"
    set "PATCH=0"
) else if "%BUMP_TYPE%"=="minor" (
    set /a "MINOR=%MINOR%+1"
    set "PATCH=0"
) else if "%BUMP_TYPE%"=="patch" (
    set /a "PATCH=%PATCH%+1"
)

set "NEW_VERSION=v%MAJOR%.%MINOR%.%PATCH%"

echo 📈 Bump da %CURRENT_VERSION% a %NEW_VERSION% (%BUMP_TYPE%)

:: Se non è stata fornita una descrizione, la chiede
if "%DESCRIPTION%"=="" (
    set /p "DESCRIPTION=📝 Descrizione per questa release: "
)

if "%DESCRIPTION%"=="" (
    set "DESCRIPTION=Release %NEW_VERSION%"
)

:: Data corrente in formato ISO
for /f "tokens=1-3 delims=/" %%a in ('date /t') do (
    set "CURRENT_DATE=%%c-%%a-%%b"
)

echo.
echo 🚀 Creazione release %NEW_VERSION% - %DESCRIPTION%
echo 📅 Data: %CURRENT_DATE%
echo.

:: Crea un file temporaneo per il nuovo CHANGELOG
echo 🔄 Aggiornamento CHANGELOG...

:: Crea il nuovo contenuto del CHANGELOG
(
    echo # Changelog - Automazione e Qualità del Codice
    echo.
    echo ## [Unreleased]
    echo.
    echo ## %NEW_VERSION% - %DESCRIPTION% ^(%CURRENT_DATE%^)
    echo.
) > CHANGELOG_NEW.md

:: Copia il contenuto di [Unreleased] sotto la nuova versione, saltando la sezione [Unreleased]
powershell -Command "$content = Get-Content 'CHANGELOG.md'; $inUnreleased = $false; $foundFirstVersion = $false; foreach($line in $content) { if($line -match '^## \[Unreleased\]') { $inUnreleased = $true; continue } if($line -match '^## v[0-9]') { if($inUnreleased) { $inUnreleased = $false; $foundFirstVersion = $true } if($foundFirstVersion) { Add-Content 'CHANGELOG_NEW.md' $line } } elseif($inUnreleased -and $line.Trim() -ne '') { Add-Content 'CHANGELOG_NEW.md' $line } elseif($foundFirstVersion) { Add-Content 'CHANGELOG_NEW.md' $line } }"

:: Sostituisce il CHANGELOG originale
move CHANGELOG_NEW.md CHANGELOG.md

echo ✅ CHANGELOG aggiornato con la nuova versione

:: Crea Git tag
echo 🏷️  Creazione Git tag...
git add CHANGELOG.md
git commit -m "release: %NEW_VERSION% - %DESCRIPTION%"
git tag -a "%NEW_VERSION%" -m "Release %NEW_VERSION% - %DESCRIPTION%"

echo.
echo ✅ Bump di versione completato!
echo 📋 Nuova versione: %NEW_VERSION%
echo 🏷️  Tag Git creato: %NEW_VERSION%
echo.
echo 📤 Per pubblicare:
echo   git push origin main
echo   git push origin %NEW_VERSION%
