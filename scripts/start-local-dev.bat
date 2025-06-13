@echo off
echo 🚀 Thothix Local Development with Clerk Integration
echo.

echo 📋 Prerequisites Check:
echo 1. Ngrok installed? (Download from https://ngrok.com/download)
echo 2. Clerk project configured? (dashboard.clerk.com)
echo 3. Backend environment ready?
echo.

echo 🔧 Starting Local Development Environment...
echo.

echo Step 1: Starting backend server...
start cmd /k "cd /d %~dp0\..\backend && echo 🔧 Starting Thothix Backend on port 30000... && go run main.go"

echo Step 2: Waiting for backend to start...
timeout /t 5 /nobreak > nul

echo Step 3: Starting ngrok tunnel with custom URL...
echo 🌐 Using ngrok URL: https://flying-mullet-socially.ngrok-free.app
echo 📝 Webhook URL: https://flying-mullet-socially.ngrok-free.app/api/v1/auth/webhooks/clerk
echo.
start cmd /k "echo 🌐 Starting ngrok tunnel with custom URL... && ngrok http --url=flying-mullet-socially.ngrok-free.app 30000"

echo.
echo ✅ Development environment started!
echo.
echo 📋 Next Steps:
echo 1. Your ngrok URL is ready: https://flying-mullet-socially.ngrok-free.app
echo 2. Go to Clerk Dashboard → Webhooks → Add Endpoint
echo 3. Use this exact webhook URL: https://flying-mullet-socially.ngrok-free.app/api/v1/auth/webhooks/clerk
echo 4. Select events: user.created, user.updated, user.deleted
echo 5. Copy the webhook signing secret to your .env file
echo.
echo 🧪 Test APIs:
echo - Health: http://localhost:30000/health
echo - Health (via ngrok): https://flying-mullet-socially.ngrok-free.app/health
echo - Swagger: http://localhost:30000/swagger/index.html
echo - Swagger (via ngrok): https://flying-mullet-socially.ngrok-free.app/swagger/index.html
echo - Sync User: POST http://localhost:30000/api/v1/auth/sync
echo - Webhook: POST https://flying-mullet-socially.ngrok-free.app/api/v1/auth/webhooks/clerk
echo.
pause
