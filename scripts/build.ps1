param([string]$Action = 'dev')
$APP_NAME = "TrayApp"
$CMD_PATH = "./cmd/app"
$OUTPUT_DIR = "./build/windows"

# Генерация ресурсов
$iconSrc = "resources/icons/icon-on.ico"
$iconDst = "internal/resources/icons/icons.go"
if (Test-Path $iconSrc) {
    $dstDir = Split-Path $iconDst -Parent
    if (!(Test-Path $dstDir)) { New-Item -ItemType Directory -Force -Path $dstDir | Out-Null }
    if (!(Test-Path $iconDst)) {
        Write-Host "Generating resources..." -ForegroundColor Cyan
        fyne bundle -o $iconDst -package icons $iconSrc
    }
}

Write-Host "Build: $Action" -ForegroundColor Cyan
$env:CGO_ENABLED = "1"

if ($Action -eq 'release') {
    if (!(Test-Path $OUTPUT_DIR)) { New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null }
    go build -ldflags="-H=windowsgui -s -w" -o "$OUTPUT_DIR/$APP_NAME.exe" $CMD_PATH
}
elseif ($Action -eq 'dev') {
    if (!(Test-Path $OUTPUT_DIR)) { New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null }
    go build -v -o "$OUTPUT_DIR/$APP_NAME-dev.exe" $CMD_PATH
}
elseif ($Action -eq 'run') {
    go run $CMD_PATH
}
elseif ($Action -eq 'clean') {
    go clean -cache -modcache
    Remove-Item -Recurse -Force $OUTPUT_DIR -ErrorAction SilentlyContinue
}

if ($LASTEXITCODE -eq 0) { Write-Host "Done" -ForegroundColor Green } else { Write-Host "Failed" -ForegroundColor Red }