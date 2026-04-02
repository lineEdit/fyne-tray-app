param(
    [string]$Action = 'dev',
    [switch]$Force
)

$APP_NAME = "TrayApp"
$CMD_PATH = "./cmd/app"
$OUTPUT_DIR = "./build/windows"
$localesSrc = "resources/locales"
$localesDst = "$OUTPUT_DIR/resources/locales"

# ===== Проверка окружения =====
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go not found in PATH" -ForegroundColor Red
    exit 1
}

Write-Host "🔨 Build: $Action" -ForegroundColor Cyan
$env:CGO_ENABLED = "1"

switch ($Action) {
    'release' {
        if (!(Test-Path $OUTPUT_DIR)) {
            New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
        }

        # Копируем локали
        if (Test-Path $localesSrc) {
            Copy-Item -Path $localesSrc -Destination $localesDst -Recurse -Force
            Write-Host "📦 Locales copied" -ForegroundColor Gray
        }

        Write-Host "🚀 Building release..." -ForegroundColor Cyan
        go build -ldflags="-H=windowsgui -s -w" -o "$OUTPUT_DIR/$APP_NAME.exe" $CMD_PATH
    }

    'dev' {
        if (!(Test-Path $OUTPUT_DIR)) {
            New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
        }

        if (Test-Path $localesSrc) {
            Copy-Item -Path $localesSrc -Destination $localesDst -Recurse -Force -ErrorAction SilentlyContinue
        }

        Write-Host "🔧 Building dev..." -ForegroundColor Cyan
        go build -v -o "$OUTPUT_DIR/$APP_NAME-dev.exe" $CMD_PATH
    }

    'run' {
        Write-Host "▶️ Running app..." -ForegroundColor Cyan
        go run $CMD_PATH
    }

    'clean' {
        Write-Host "🧹 Cleaning..." -ForegroundColor Yellow
        go clean -cache -modcache
        Remove-Item -Recurse -Force $OUTPUT_DIR -ErrorAction SilentlyContinue
        Write-Host "✅ Cleaned" -ForegroundColor Green
        exit 0
    }

    default {
        Write-Host "❌ Unknown action: $Action" -ForegroundColor Red
        Write-Host "Available: dev, release, run, clean" -ForegroundColor Yellow
        exit 1
    }
}

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Done" -ForegroundColor Green
} else {
    Write-Host "❌ Failed (exit code: $LASTEXITCODE)" -ForegroundColor Red
    exit 1
}