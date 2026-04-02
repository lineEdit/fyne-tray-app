param(
    [string]$Action = 'dev',
    [switch]$Force
)

$APP_NAME = "TrayApp"
$CMD_PATH = "./cmd/app"
$OUTPUT_DIR = "./build/windows"
$localesSrc = "resources/locales"
$localesDst = "$OUTPUT_DIR/resources/locales"

# ===== Проверка Go =====
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Go not found in PATH" -ForegroundColor Red
    exit 1
}

# ===== Проверка GCC =====
$env:CGO_ENABLED = "1"
$env:CC = "gcc"
$env:CXX = "g++"

if (!(Get-Command gcc -ErrorAction SilentlyContinue)) {
    Write-Host "❌ GCC (MinGW) not found in PATH" -ForegroundColor Red
    Write-Host ""
    Write-Host "📥 Installation instructions:" -ForegroundColor Yellow
    Write-Host "1. Download WinLibs: https://winlibs.com/" -ForegroundColor Yellow
    Write-Host "2. Extract to C:\mingw64" -ForegroundColor Yellow
    Write-Host "3. Add C:\mingw64\mingw64\bin to PATH" -ForegroundColor Yellow
    Write-Host "4. Restart PowerShell" -ForegroundColor Yellow
    Write-Host ""
    exit 1
}

# Проверка архитектуры GCC
$gccArch = gcc -dumpmachine 2>$null
if ($gccArch -notmatch "x86_64") {
    Write-Host "⚠️ Warning: GCC architecture may not match Go ($gccArch)" -ForegroundColor Yellow
}

Write-Host "🔨 Build: $Action" -ForegroundColor Cyan
Write-Host "🔧 CGO: Enabled (CC=$env:CC, GCC=$gccArch)" -ForegroundColor Gray

# ===== Сборка =====
switch ($Action) {
    'release' {
        if (!(Test-Path $OUTPUT_DIR)) {
            New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
        }

        if (Test-Path $localesSrc) {
            New-Item -ItemType Directory -Path $localesDst -Force | Out-Null
            Copy-Item -Path "$localesSrc\*.json" -Destination $localesDst -Force
            Write-Host "📦 Locales copied: $((Get-ChildItem $localesDst -Filter *.json).Count) files" -ForegroundColor Gray
        }

        Write-Host "🚀 Building release..." -ForegroundColor Cyan
        go build -ldflags="-H=windowsgui -s -w" -o "$OUTPUT_DIR/$APP_NAME.exe" $CMD_PATH
    }

    'dev' {
        if (!(Test-Path $OUTPUT_DIR)) {
            New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
        }

        if (Test-Path $localesSrc) {
            New-Item -ItemType Directory -Path $localesDst -Force | Out-Null
            Copy-Item -Path "$localesSrc\*.json" -Destination $localesDst -Force -ErrorAction SilentlyContinue
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
        go clean -cache -modcache -i
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
    Write-Host ""
    Write-Host "💡 Troubleshooting:" -ForegroundColor Yellow
    Write-Host "1. Run: go clean -cache -modcache -i" -ForegroundColor Yellow
    Write-Host "2. Check: gcc --version" -ForegroundColor Yellow
    Write-Host "3. See: build-debug.log for details" -ForegroundColor Yellow
    exit 1
}