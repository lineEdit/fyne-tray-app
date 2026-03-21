# 🗂️ Tray App

Приложение для Windows, работающее в системном трее, созданное на **Go** с использованием фреймворка **Fyne**.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Fyne](https://img.shields.io/badge/Fyne-2.x-00A8E0?style=flat&logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows-blue?style=flat)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

---

## 📋 Оглавление

- [Возможности](#-возможности)
- [Требования](#-требования)
- [Структура проекта](#-структура-проекта)
- [Быстрый старт](#-быстрый-старт)
- [Сборка через скрипт](#-сборка-через-скрипт)
- [Ручная сборка](#-ручная-сборка)
- [Конфигурация](#-конфигурация)
- [Устранение неполадок](#-устранение-неполадок)
- [Лицензия](#-лицензия)

---

## ✨ Возможности

- 🖥️ Работа в системном трее Windows
- 🎨 Графический интерфейс на базе Fyne
- ⚙️ Настройки через встроенное меню
- 📝 Логирование в файл
- 🔄 Автообновление (опционально)
- 🌐 Локализация (RU/EN)

---

## 🛠 Требования

### Обязательные

| Компонент | Версия | Примечание |
|-----------|--------|------------|
| **Go** | 1.21+ | [Скачать](https://go.dev/dl/) |
| **GCC (MinGW-w64)** | 8.1+ | Требуется для CGO |
| **PowerShell** | 5.1+ | Для скриптов сборки |

### Установка MinGW-w64 (Windows)

```powershell
# Вариант 1: Через WinLibs (рекомендуется)
# 1. Скачайте с https://winlibs.com/ (MSVCRT runtime, x86_64)
# 2. Распакуйте в C:\mingw64
# 3. Добавьте в PATH:
[Environment]::SetEnvironmentVariable("Path", "C:\mingw64\mingw64\bin;" + $env:Path, "User")

# Вариант 2: Через MSYS2
winget install -e --id MSYS2.MSYS2
# Затем в терминале MSYS2:
# pacman -S mingw-w64-x86_64-toolchain