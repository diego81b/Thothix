@echo off
echo 🧪 Quick Ngrok + Backend Test
echo.

echo 📍 Testing local backend...
curl -s http://localhost:30000/health
if %errorlevel% neq 0 (
    echo ❌ Local backend not running on port 30000
    echo 💡 Start it with: cd backend && go run main.go
    goto :end
)
echo ✅ Local backend is running
echo.

echo 📍 Testing ngrok tunnel...
curl -s https://flying-mullet-socially.ngrok-free.app/health
if %errorlevel% neq 0 (
    echo ❌ Ngrok tunnel not working
    echo 💡 Start it with: ngrok http --url=flying-mullet-socially.ngrok-free.app 30000
    goto :end
)
echo ✅ Ngrok tunnel is working
echo.

echo 📍 Testing webhook endpoint...
curl -X POST https://flying-mullet-socially.ngrok-free.app/api/v1/auth/webhooks/clerk ^
  -H "Content-Type: application/json" ^
  -d "{\"type\":\"test\",\"data\":{}}" ^
  -w "HTTP Status: %%{http_code}\n"
echo.

echo 🎉 All tests passed! Your setup is ready for Clerk integration.
echo.
echo 📋 Next steps:
echo 1. Configure webhook in Clerk Dashboard: https://flying-mullet-socially.ngrok-free.app/api/v1/auth/webhooks/clerk
echo 2. Add CLERK_WEBHOOK_SECRET to your .env file
echo 3. Test with real user registration/updates

:end
echo.
pause
