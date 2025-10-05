# Homebrew Tap Setup for SSH Keeper

## Создание собственного Homebrew Tap

### 1. Создать репозиторий для tap

```bash
# Создать новый репозиторий на GitHub
# Название: homebrew-ssh-keeper
# Описание: Homebrew tap for SSH Keeper
```

### 2. Настроить формулу

```bash
# Клонировать tap репозиторий
git clone https://github.com/ankogit/homebrew-ssh-keeper.git
cd homebrew-ssh-keeper

# Создать формулу
mkdir -p Formula
cp ../ssh_keeper/Formula/ssh-keeper-homebrew.rb Formula/ssh-keeper.rb

# Рассчитать SHA256 для архива
curl -L -O https://github.com/ankogit/ssh_keeper/archive/v0.1.0.tar.gz
shasum -a 256 v0.1.0.tar.gz
# Заменить YOUR_SHA256_HERE на полученное значение

# Закоммитить
git add Formula/ssh-keeper.rb
git commit -m "Add SSH Keeper formula"
git push origin main
```

### 3. Установка через tap

```bash
# Добавить tap
brew tap ankogit/ssh-keeper

# Установить SSH Keeper
brew install ssh-keeper

# Обновить
brew upgrade ssh-keeper
```

## Альтернативные пакетные менеджеры

### Snap (Linux)

```bash
# Создать snapcraft.yaml
name: ssh-keeper
version: '0.1.0'
summary: SSH Connection Manager with secure password storage
description: |
  A beautiful and secure CLI tool for managing SSH connections
  with a modern TUI interface.

grade: stable
confinement: strict

apps:
  ssh-keeper:
    command: ssh-keeper
    plugs: [network, home]

parts:
  ssh-keeper:
    plugin: dump
    source: https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-linux-amd64.tar.gz
    source-type: tar
```

### Chocolatey (Windows)

```xml
<!-- ssh-keeper.nuspec -->
<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>ssh-keeper</id>
    <version>0.1.0</version>
    <title>SSH Keeper</title>
    <authors>SSH Keeper Contributors</authors>
    <description>A beautiful and secure CLI tool for managing SSH connections</description>
    <projectUrl>https://github.com/ankogit/ssh_keeper</projectUrl>
    <licenseUrl>https://github.com/ankogit/ssh_keeper/blob/main/LICENSE</licenseUrl>
    <requireLicenseAcceptance>false</requireLicenseAcceptance>
  </metadata>
  <files>
    <file src="tools\**" target="tools" />
  </files>
</package>
```

### Scoop (Windows)

```json
{
  "version": "0.1.0",
  "description": "A beautiful and secure CLI tool for managing SSH connections",
  "homepage": "https://github.com/ankogit/ssh_keeper",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-windows-amd64.zip",
      "hash": "YOUR_HASH_HERE",
      "extract_dir": "ssh-keeper-0.1.0-windows-amd64"
    }
  },
  "bin": "ssh-keeper.exe",
  "checkver": "github",
  "autoupdate": {
    "architecture": {
      "64bit": {
        "url": "https://github.com/ankogit/ssh_keeper/releases/download/$version/ssh-keeper-$version-windows-amd64.zip"
      }
    }
  }
}
```

## Упрощенная установка

### Однострочные команды

**macOS & Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
iwr -useb https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.ps1 | iex
```

**Windows (CMD):**

```cmd
powershell -Command "iwr -useb https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.ps1 | iex"
```

### Через пакетные менеджеры

**Homebrew (macOS):**

```bash
brew tap ankogit/ssh-keeper
brew install ssh-keeper
```

**Snap (Linux):**

```bash
sudo snap install ssh-keeper
```

**Chocolatey (Windows):**

```cmd
choco install ssh-keeper
```

**Scoop (Windows):**

```powershell
scoop install ssh-keeper
```

## Преимущества упрощенной установки

1. **Автоматическое определение платформы** - скрипт сам определяет ОС и архитектуру
2. **Однострочная команда** - как у Docker, Node.js, и других популярных инструментов
3. **Автоматическая установка в PATH** - готов к использованию сразу после установки
4. **Проверка установки** - автоматическая верификация успешной установки
5. **Красивый вывод** - цветной ASCII-арт и информативные сообщения
6. **Обработка ошибок** - понятные сообщения об ошибках и их решения
7. **Кроссплатформенность** - работает на всех поддерживаемых платформах
