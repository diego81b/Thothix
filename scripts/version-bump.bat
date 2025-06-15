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

echo 🔍 Reading current version from CHANGELOG...

:: Find first version in CHANGELOG (skip [Unreleased])
for /f "tokens=2 delims= " %%a in ('findstr "^## v" CHANGELOG.md') do (
    set "TEMP_VERSION=%%a"
    :: Skip [Unreleased] and only process real versions
    if not "!TEMP_VERSION!"=="[Unreleased]" (
        set "CURRENT_VERSION=!TEMP_VERSION!"
        goto :version_found
    )
)
:version_found
echo 📋 Current version found: %CURRENT_VERSION%

:: Extract version numbers (remove 'v')
set "VERSION_CLEAN=%CURRENT_VERSION:v=%"

:: Parse version using delimiters
for /f "tokens=1,2,3 delims=." %%a in ("%VERSION_CLEAN%") do (
    set "MAJOR=%%a"
    set "MINOR=%%b"
    set "PATCH=%%c"
)

echo 🔢 Parsed version: MAJOR=%MAJOR%, MINOR=%MINOR%, PATCH=%PATCH%

:: Calculate new version
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

echo 📈 Bump from %CURRENT_VERSION% to %NEW_VERSION% (%BUMP_TYPE%)

:: If no description provided, create a default one
if "%DESCRIPTION%"=="" (
    set "DESCRIPTION=Release %NEW_VERSION%"
)

echo 📝 Description: %DESCRIPTION%

:: Current date in ISO format
for /f "tokens=2 delims==" %%i in ('wmic os get localdatetime /value') do set datetime=%%i
set "CURRENT_DATE=%datetime:~0,4%-%datetime:~4,2%-%datetime:~6,2%"

echo 🚀 Creating release %NEW_VERSION% - %DESCRIPTION%
echo 📅 Date: %CURRENT_DATE%

:: Get complete commit message (subject + body) before making any changes
echo 📝 Capturing latest commit details...
git log -1 --pretty=format:"%B" > temp_commit_msg.txt

:: Update CHANGELOG...
echo 🔄 Updating CHANGELOG with new version entry...

:: Create backup
copy CHANGELOG.md CHANGELOG_BACKUP.md >nul

:: Create new CHANGELOG header with complete commit message
(
    echo # Changelog - Automation and Code Quality
    echo.
    echo ## [Unreleased]
    echo.
    echo ## %NEW_VERSION% - %DESCRIPTION% ^(%CURRENT_DATE%^)
    echo.
    type temp_commit_msg.txt
    echo.
) > CHANGELOG_NEW.md

:: Copy all existing versions (skip header line)
powershell -Command "$content = Get-Content 'CHANGELOG.md'; $skipHeader = $true; foreach($line in $content) { if($skipHeader -and $line -match '^# ') { $skipHeader = $false; continue } if(-not $skipHeader) { Add-Content 'CHANGELOG_NEW.md' $line } }"

:: Replace original CHANGELOG
move CHANGELOG_NEW.md CHANGELOG.md

:: Clean up temporary files
del temp_commit_msg.txt
if exist CHANGELOG.md del CHANGELOG_BACKUP.md

echo ✅ CHANGELOG updated with new version

:: Create Git commit and tag
echo 🏷️  Creating Git commit and tag...
git add CHANGELOG.md
git commit -m "release: %NEW_VERSION% - %DESCRIPTION%"
git tag -a "%NEW_VERSION%" -m "Release %NEW_VERSION% - %DESCRIPTION%"

echo.
echo ✅ Version bump completed!
echo 📋 New version: %NEW_VERSION%
echo 🏷️  Git tag created: %NEW_VERSION%
echo.
echo 📤 To publish:
echo   git push origin main
echo   git push origin %NEW_VERSION%
