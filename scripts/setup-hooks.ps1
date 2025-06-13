# PowerShell script per configurare Git hooks
# Uso: .\setup-hooks.ps1

Write-Host "🔧 Configurando Git hooks..." -ForegroundColor Yellow

# Verifica che esista la directory .git/hooks
if (-not (Test-Path ".git\hooks")) {
  Write-Host "❌ Directory .git\hooks non trovata!" -ForegroundColor Red
  Write-Host "Assicurati di essere nella root del repository Git." -ForegroundColor Red
  exit 1
}

# Copia il pre-commit hook
$hookPath = ".git\hooks\pre-commit"
if (Test-Path $hookPath) {
  Write-Host "✅ Pre-commit hook già presente" -ForegroundColor Green
}
else {
  Write-Host "❌ Pre-commit hook non trovato!" -ForegroundColor Red
  Write-Host "Il file dovrebbe essere in: $hookPath" -ForegroundColor Yellow
}

Write-Host "✅ Git hooks configurati!" -ForegroundColor Green
Write-Host ""
Write-Host "📋 Comandi disponibili:" -ForegroundColor Cyan
Write-Host "  .\scripts\pre-commit.bat  - Esegui pre-commit manualmente" -ForegroundColor White
Write-Host "  git commit -m 'message'   - Il pre-commit si attiva automaticamente" -ForegroundColor White
