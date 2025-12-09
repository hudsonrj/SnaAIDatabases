@echo off
REM Wrapper para executar snip.exe ou usar go run como fallback
cd /d "%~dp0"

if exist "snip.exe" (
    snip.exe %*
) else (
    go run main.go %*
)

