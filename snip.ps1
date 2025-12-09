# Script wrapper para executar Snip
# Se snip.exe não funcionar, usa go run como fallback

$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$exePath = Join-Path $scriptPath "snip.exe"
$mainPath = Join-Path $scriptPath "main.go"

if (Test-Path $exePath) {
    try {
        & $exePath $args
        exit $LASTEXITCODE
    } catch {
        Write-Host "snip.exe não pôde ser executado, usando go run..." -ForegroundColor Yellow
    }
}

# Fallback para go run
if (Test-Path $mainPath) {
    Push-Location $scriptPath
    go run main.go $args
    Pop-Location
} else {
    Write-Host "Erro: Não foi possível encontrar snip.exe ou main.go" -ForegroundColor Red
    exit 1
}

