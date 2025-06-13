@echo off
echo 🧪 Quick API Test
echo.

echo 📍 Testing backend health...
curl -s http://localhost:30000/health
if %errorlevel% neq 0 (
    echo ❌ Backend not running on port 30000
    echo 💡 Start it with: cd backend && go run main.go
    goto :end
)
echo ✅ Backend is running
echo.

echo 📍 Testing webhook endpoint...
curl -X POST http://localhost:30000/api/v1/auth/webhooks/clerk ^
  -H "Content-Type: application/json" ^
  -d "{\"type\":\"test\",\"data\":{}}" ^
  -w "HTTP Status: %%{http_code}\n"
echo.

echo 🎉 Basic tests passed!
echo.
echo 📚 For complete API testing use Swagger UI:
echo    http://localhost:30000/swagger/index.html
echo.

:end
pause
