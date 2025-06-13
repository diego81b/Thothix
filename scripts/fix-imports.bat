@echo off
echo 🔧 Fixing all Go import formatting issues...

cd /d "%~dp0\..\backend"

echo 📝 Step 1: Running gofumpt on all files...
gofumpt -w .

echo 🔍 Step 2: Checking for remaining formatting issues...
gofumpt -l .

if errorlevel 1 (
    echo ⚠️  Some files still need formatting. Running again...
    gofumpt -w .
) else (
    echo ✅ All files are properly formatted
)

echo 📋 Step 3: Adding formatted files to git...
cd /d "%~dp0\.."
git add backend/

echo 🎉 Import formatting fix completed!
echo 💡 Now save any Go file in VS Code to test format-on-save
pause
