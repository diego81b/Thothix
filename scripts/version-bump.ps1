# Script PowerShell per bump automatico di versione
# Uso: .\version-bump.ps1 -BumpType [major|minor|patch] [-Description "Optional description"]

param(
  [Parameter(Mandatory = $true)]
  [ValidateSet("major", "minor", "patch")]
  [string]$BumpType,

  [Parameter(Mandatory = $false)]
  [string]$Description = ""
)

Write-Host "🔍 Lettura versione corrente..." -ForegroundColor Cyan

# Legge la versione corrente dal CHANGELOG
$changelogContent = Get-Content "CHANGELOG.md" -ErrorAction Stop
$versionLine = $changelogContent | Where-Object { $_ -match "^## v\d+\.\d+\.\d+" } | Select-Object -First 1

if (-not $versionLine) {
  Write-Host "❌ Impossibile trovare la versione corrente nel CHANGELOG" -ForegroundColor Red
  exit 1
}

# Estrae la versione corrente
if ($versionLine -match "^## (v\d+\.\d+\.\d+)") {
  $currentVersion = $matches[1]
  $versionNumbers = $currentVersion.Substring(1) # Rimuove la 'v'
}
else {
  Write-Host "❌ Formato versione non riconosciuto: $versionLine" -ForegroundColor Red
  exit 1
}

# Estrae major, minor, patch
$versionParts = $versionNumbers.Split('.')
$major = [int]$versionParts[0]
$minor = [int]$versionParts[1]
$patch = [int]$versionParts[2]

# Calcola la nuova versione
switch ($BumpType) {
  "major" {
    $major++
    $minor = 0
    $patch = 0
  }
  "minor" {
    $minor++
    $patch = 0
  }
  "patch" {
    $patch++
  }
}

$newVersion = "v$major.$minor.$patch"

Write-Host "📈 Bump da $currentVersion a $newVersion ($BumpType)" -ForegroundColor Green

# Se non è stata fornita una descrizione, la chiede
if ([string]::IsNullOrEmpty($Description)) {
  $Description = Read-Host "📝 Descrizione per questa release"
  if ([string]::IsNullOrEmpty($Description)) {
    $Description = "Release $newVersion"
  }
}

# Data corrente in formato ISO
$currentDate = Get-Date -Format "yyyy-MM-dd"

Write-Host ""
Write-Host "🚀 Creazione release $newVersion - $Description" -ForegroundColor Yellow
Write-Host "📅 Data: $currentDate" -ForegroundColor Yellow
Write-Host ""

# Aggiorna il CHANGELOG
Write-Host "🔄 Aggiornamento CHANGELOG..." -ForegroundColor Cyan

# Legge il contenuto corrente
$content = Get-Content "CHANGELOG.md"

# Crea il nuovo contenuto
$newContent = @()
$newContent += "# Changelog - Automazione e Qualità del Codice"
$newContent += ""
$newContent += "## [Unreleased]"
$newContent += ""
$newContent += "## $newVersion - $Description ($currentDate)"
$newContent += ""

# Copia il contenuto di [Unreleased] sotto la nuova versione
$inUnreleased = $false
$foundFirstVersion = $false

foreach ($line in $content) {
  if ($line -match "^## \[Unreleased\]") {
    $inUnreleased = $true
    continue
  }

  if ($line -match "^## v\d+") {
    if ($inUnreleased) {
      $inUnreleased = $false
      $foundFirstVersion = $true
    }
    if ($foundFirstVersion) {
      $newContent += $line
    }
  }
  elseif ($inUnreleased -and $line.Trim() -ne "") {
    $newContent += $line
  }
  elseif ($foundFirstVersion) {
    $newContent += $line
  }
}

# Scrive il nuovo CHANGELOG
$newContent | Out-File "CHANGELOG.md" -Encoding UTF8

Write-Host "✅ CHANGELOG aggiornato con la nuova versione" -ForegroundColor Green

# Verifica che Git sia pulito
$gitStatus = git status --porcelain
if ($gitStatus -and ($gitStatus | Where-Object { $_ -notmatch "CHANGELOG.md" })) {
  Write-Host "⚠️  Ci sono modifiche non committate. Commit prima di continuare." -ForegroundColor Yellow
  git status
  $continue = Read-Host "Continuare comunque? (y/N)"
  if ($continue -ne "y" -and $continue -ne "Y") {
    Write-Host "❌ Operazione annullata" -ForegroundColor Red
    exit 1
  }
}

# Crea Git tag
Write-Host "🏷️  Creazione Git tag..." -ForegroundColor Cyan

try {
  git add CHANGELOG.md
  git commit -m "release: $newVersion - $Description"
  git tag -a "$newVersion" -m "Release $newVersion - $Description"

  Write-Host ""
  Write-Host "✅ Bump di versione completato!" -ForegroundColor Green
  Write-Host "📋 Nuova versione: $newVersion" -ForegroundColor Green
  Write-Host "🏷️  Tag Git creato: $newVersion" -ForegroundColor Green
  Write-Host ""
  Write-Host "📤 Per pubblicare:" -ForegroundColor Yellow
  Write-Host "  git push origin main" -ForegroundColor Yellow
  Write-Host "  git push origin $newVersion" -ForegroundColor Yellow
}
catch {
  Write-Host "❌ Errore durante la creazione del tag Git: $_" -ForegroundColor Red
  exit 1
}
