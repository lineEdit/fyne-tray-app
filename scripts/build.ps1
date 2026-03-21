param(
    [string]$Action = 'dev',
    [switch]$Force
)

$APP_NAME = "TrayApp"
$CMD_PATH = "./cmd/app"
$OUTPUT_DIR = "./build/windows"

# ===== Проверка окружения =====
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go not found in PATH" -ForegroundColor Red
    exit 1
}

# ===== Генерация ресурсов =====
$iconSrc = "resources/icons/Icon-on.ico"
$iconDst = "internal/resources/icons/icons.go"
$localesSrc = "resources/locales"
$localesDst = "$OUTPUT_DIR/resources/locales"

if (Test-Path $iconSrc) {
    $dstDir = Split-Path $iconDst -Parent
    if (!(Test-Path $dstDir)) {
        New-Item -ItemType Directory -Force -Path $dstDir | Out-Null
    }

    if (!(Test-Path $iconDst) -or $Force) {
        Write-Host "📦 Generating icon resources..." -ForegroundColor Cyan
        go run fyne.io/tools/cmd/fyne@latest bundle -o $iconDst -package icons $iconSrc
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ Failed to generate icon resources" -ForegroundColor Red
            exit 1
        }
        Write-Host "✅ Icon resources generated" -ForegroundColor Green
    }
} else {
    Write-Host "⚠️ Icon not found: $iconSrc" -ForegroundColor Yellow
}

Write-Host "🔨 Build: $Action" -ForegroundColor Cyan
$env:CGO_ENABLED = "1"

# ===== Действия =====
switch ($Action) {
    'release' {
        if (!(Test-Path $OUTPUT_DIR)) {
            New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
        }

        # Копируем локали для релиза
        if (Test-Path $localesSrc) {
            Copy-Item -Path $localesSrc -Destination $localesDst -Recurse -Force
            Write-Host "📦 Locales copied to $localesDst" -ForegroundColor Gray
        }

        Write-Host "🚀 Building release..." -ForegroundColor Cyan
        go build -ldflags="-H=windowsgui -s -w" -o "$OUTPUT_DIR/$APP_NAME.exe" $CMD_PATH
    }

    'dev' {
        if (!(Test-Path $OUTPUT_DIR)) {
            New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
        }

        # Копируем локали для dev
        if (Test-Path $localesSrc) {
            Copy-Item -Path $localesSrc -Destination $localesDst -Recurse -Force
            Write-Host "📦 Locales copied to $localesDst" -ForegroundColor Gray
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

# ===== Финальная проверка =====
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Done" -ForegroundColor Green
} else {
    Write-Host "❌ Failed (exit code: $LASTEXITCODE)" -ForegroundColor Red
    exit 1
}