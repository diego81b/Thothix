@echo off
echo ====================================
echo Thothix Go Formatting Diagnostic
echo ====================================
echo.

echo 1. Checking Go installation...
go version
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    goto :end
)
echo.

echo 2. Checking gofumpt installation...
gofumpt --version
if %errorlevel% neq 0 (
    echo ERROR: gofumpt is not installed or not in PATH
    echo To install: go install mvdan.cc/gofumpt@latest
    goto :end
)
echo.

echo 3. Checking Go tools in GOPATH...
set GOPATH_BIN=%GOPATH%\bin
if exist "%GOPATH_BIN%\gofumpt.exe" (
    echo ✓ gofumpt found in GOPATH\bin
) else (
    echo ⚠ gofumpt not found in %GOPATH_BIN%
)
echo.

echo 4. Testing formatting on backend files...
cd /d "%~dp0..\backend"
gofumpt -l .
if %errorlevel% neq 0 (
    echo Some files need formatting. Running gofumpt...
    gofumpt -w .
    echo Formatting complete.
) else (
    echo ✓ All files are properly formatted
)
echo.

echo 5. Running golangci-lint...
golangci-lint --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠ golangci-lint not found. Please install it.
) else (
    golangci-lint run --timeout=30s
)
echo.

echo 6. VS Code settings check...
if exist "%~dp0..\.vscode\settings.json" (
    echo ✓ VS Code settings.json found
    findstr /i "gofumpt" "%~dp0..\.vscode\settings.json" >nul
    if %errorlevel% equ 0 (
        echo ✓ gofumpt configured in VS Code
    ) else (
        echo ⚠ gofumpt not found in VS Code settings
    )
) else (
    echo ⚠ VS Code settings.json not found
)
echo.

echo 7. Testing format-on-save simulation...
echo Creating test file...
cd /d "%~dp0.."
echo package main > test-diagnostic.go
echo. >> test-diagnostic.go
echo import( >> test-diagnostic.go
echo "fmt" >> test-diagnostic.go
echo  "log" >> test-diagnostic.go
echo ) >> test-diagnostic.go
echo. >> test-diagnostic.go
echo func main( ){ >> test-diagnostic.go
echo fmt.Println("test") >> test-diagnostic.go
echo } >> test-diagnostic.go

echo Before formatting:
type test-diagnostic.go
echo.
echo After formatting:
gofumpt -w test-diagnostic.go
type test-diagnostic.go
del test-diagnostic.go >nul 2>&1
echo.

:end
echo ====================================
echo Diagnostic complete.
echo.
echo If you're still experiencing issues:
echo 1. Restart VS Code completely
echo 2. Run "Go: Install/Update Tools" from Command Palette
echo 3. Check that no other formatters are interfering
echo 4. Verify .vscode/settings.json has "go.formatTool": "gofumpt"
echo ====================================
pause
