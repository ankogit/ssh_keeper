# Зависимости SSH Keeper

## Обзор зависимостей

SSH Keeper - это Go-приложение, которое использует различные типы зависимостей для работы и сборки.

## Типы зависимостей

### 1. Build Dependencies (зависимости для сборки)

#### Homebrew (macOS)

```ruby
depends_on "go" => :build
```

#### Debian/Ubuntu

```debian
Build-Depends: debhelper (>= 11), golang-go (>= 1.21), pkg-config
```

**Что это означает:**

- `go` - компилятор Go для сборки приложения
- `debhelper` - инструменты для создания Debian пакетов
- `pkg-config` - утилита для управления флагами компиляции

### 2. Runtime Dependencies (зависимости для работы)

#### Homebrew (macOS)

```ruby
# Go приложения обычно статически линкуются
# Поэтому runtime зависимости минимальны
```

#### Debian/Ubuntu

```debian
Depends: ${shlibs:Depends}, ${misc:Depends}
```

**Что это означает:**

- `${shlibs:Depends}` - автоматически определяемые системные библиотеки
- `${misc:Depends}` - другие зависимости, определенные автоматически

### 3. Optional Dependencies (опциональные зависимости)

#### Homebrew (macOS)

```ruby
depends_on "sshpass" => :optional  # Для подключений по паролю
depends_on "gnupg" => :optional    # Для управления GPG ключами
```

#### Debian/Ubuntu

```debian
Recommends: sshpass, gnupg
Suggests: openssh-client
```

**Что это означает:**

- `Recommends` - рекомендуемые пакеты (устанавливаются по умолчанию)
- `Suggests` - предлагаемые пакеты (не устанавливаются автоматически)

## Go Module Dependencies (go.mod)

```go
require (
    github.com/charmbracelet/bubbles v0.21.0    // TUI компоненты
    github.com/charmbracelet/bubbletea v1.3.4   // TUI фреймворк
    github.com/charmbracelet/lipgloss v1.1.0    // Стили для TUI
    github.com/creack/pty v1.1.24               // PTY поддержка
    github.com/google/uuid v1.6.0               // UUID генерация
    github.com/lucasb-eyer/go-colorful v1.2.0   // Цвета
    github.com/muesli/termenv v0.16.0           // Терминал
    golang.org/x/crypto v0.42.0                  // Криптография
)
```

## Системные зависимости

### macOS

- **Минимальные**: macOS 10.15+ (Catalina)
- **Рекомендуемые**: macOS 11+ (Big Sur)

### Linux

- **Минимальные**: glibc 2.17+
- **Рекомендуемые**: Ubuntu 18.04+, Debian 10+

### Windows

- **Минимальные**: Windows 10
- **Рекомендуемые**: Windows 11

## Проверка зависимостей

### Проверить установленные зависимости

```bash
# macOS
brew list

# Ubuntu/Debian
dpkg -l | grep ssh-keeper
apt list --installed | grep ssh-keeper
```

### Проверить зависимости пакета

```bash
# macOS
brew deps ssh-keeper

# Ubuntu/Debian
apt-cache depends ssh-keeper
```

## Устранение проблем с зависимостями

### Проблема: "go: command not found"

```bash
# macOS
brew install go

# Ubuntu/Debian
sudo apt install golang-go
```

### Проблема: "pkg-config: command not found"

```bash
# macOS
brew install pkg-config

# Ubuntu/Debian
sudo apt install pkg-config
```

### Проблема: Отсутствуют опциональные зависимости

```bash
# macOS
brew install sshpass gnupg

# Ubuntu/Debian
sudo apt install sshpass gnupg openssh-client
```

## Автоматическая установка зависимостей

### Скрипт проверки зависимостей

```bash
#!/bin/bash
# check-dependencies.sh

echo "🔍 Проверка зависимостей для SSH Keeper..."

# Проверяем Go
if ! command -v go &> /dev/null; then
    echo "❌ Go не установлен"
    echo "Установите: brew install go (macOS) или sudo apt install golang-go (Ubuntu)"
    exit 1
else
    echo "✅ Go установлен: $(go version)"
fi

# Проверяем опциональные зависимости
echo ""
echo "📦 Опциональные зависимости:"

if command -v sshpass &> /dev/null; then
    echo "✅ sshpass установлен"
else
    echo "⚠️  sshpass не установлен (для подключений по паролю)"
fi

if command -v gpg &> /dev/null; then
    echo "✅ GPG установлен"
else
    echo "⚠️  GPG не установлен (для управления ключами)"
fi

if command -v ssh &> /dev/null; then
    echo "✅ OpenSSH установлен"
else
    echo "⚠️  OpenSSH не установлен (для SSH подключений)"
fi

echo ""
echo "🎉 Проверка завершена!"
```

## Рекомендации по зависимостям

### Для разработчиков

1. **Всегда указывайте минимальные версии** зависимостей
2. **Используйте опциональные зависимости** для дополнительных функций
3. **Тестируйте на разных версиях** операционных систем
4. **Документируйте системные требования**

### Для пользователей

1. **Устанавливайте рекомендуемые зависимости** для полной функциональности
2. **Обновляйте зависимости** регулярно
3. **Проверяйте системные требования** перед установкой
4. **Сообщайте о проблемах** с зависимостями в Issues



