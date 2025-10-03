# Настройка репозитория для распространения ssh-keeper

## Что было создано

### 1. Homebrew Formula (macOS)

- **Файл**: `Formula/ssh-keeper.rb`
- **Назначение**: Позволяет устанавливать через `brew install`
- **Использование**:
  ```bash
  brew tap yourusername/ssh-keeper
  brew install ssh-keeper
  ```

### 2. Debian Package (Ubuntu/Debian)

- **Файлы**: `debian/control`, `debian/rules`, `debian/changelog`
- **Скрипт сборки**: `build-deb.sh`
- **Назначение**: Позволяет устанавливать через `apt install`

### 3. GitHub Actions

- **Файл**: `.github/workflows/release.yml`
- **Назначение**: Автоматическая сборка и создание релизов при создании тегов

### 4. Скрипты установки

- **Файл**: `install.sh` - универсальный скрипт установки
- **Файл**: `INSTALL.md` - подробные инструкции по установке

## Шаги для настройки

### 1. Обновите информацию в файлах

Замените `yourusername` на ваш GitHub username в следующих файлах:

- `Formula/ssh-keeper.rb` (строка 3)
- `debian/control` (строки 3, 5, 6)
- `debian/changelog` (строка 5)
- `.github/workflows/release.yml` (строка 3)
- `install.sh` (строка 7)
- `INSTALL.md` (все вхождения)

### 2. Создайте GitHub репозиторий

```bash
# Если репозиторий еще не создан
git init
git add .
git commit -m "Initial commit with package support"
git branch -M main
git remote add origin https://github.com/yourusername/ssh-keeper.git
git push -u origin main
```

### 3. Настройте Homebrew tap

```bash
# Создайте tap репозиторий
# Название должно быть: homebrew-ssh-keeper
# URL: https://github.com/yourusername/homebrew-ssh-keeper

# Скопируйте формулу в tap репозиторий
cp Formula/ssh-keeper.rb ../homebrew-ssh-keeper/ssh-keeper.rb
```

### 4. Создайте первый релиз

```bash
# Создайте тег для релиза
git tag v0.1.0
git push origin v0.1.0

# GitHub Actions автоматически создаст релиз
```

### 5. Обновите SHA256 в Homebrew формуле

После создания релиза:

1. Скачайте macOS ARM64 архив
2. Вычислите SHA256: `shasum -a 256 ssh-keeper-0.1.0-darwin-arm64.tar.gz`
3. Обновите значение в `Formula/ssh-keeper.rb`

## Команды для тестирования

### Сборка пакетов

```bash
# Создать все релизные файлы
make github-release

# Создать только Debian пакет
make deb

# Создать только Homebrew формулу
make homebrew-formula
```

### Тестирование установки

```bash
# Тест скрипта установки
curl -fsSL https://raw.githubusercontent.com/yourusername/ssh-keeper/main/install.sh | bash
```

## Структура файлов

```
ssh_keeper/
├── Formula/
│   └── ssh-keeper.rb          # Homebrew formula
├── debian/
│   ├── control                # Debian package metadata
│   ├── rules                  # Debian build rules
│   └── changelog              # Debian changelog
├── .github/
│   └── workflows/
│       └── release.yml        # GitHub Actions workflow
├── build-deb.sh               # Debian package build script
├── install.sh                 # Universal installation script
├── INSTALL.md                 # Installation documentation
└── Makefile                   # Updated with new targets
```

## Что нужно сделать дополнительно

1. **Обновить README.md** с инструкциями по установке
2. **Создать LICENSE файл** (MIT рекомендуется)
3. **Настроить GitHub Pages** для документации (опционально)
4. **Добавить CI/CD** для автоматического тестирования
5. **Создать ISSUE_TEMPLATE** и **PULL_REQUEST_TEMPLATE**

## Проверка работы

После настройки проверьте:

1. **Homebrew**: `brew install yourusername/ssh-keeper/ssh-keeper`
2. **apt**: `sudo dpkg -i ssh-keeper_0.1.0_amd64.deb`
3. **Скрипт**: `curl -fsSL https://raw.githubusercontent.com/yourusername/ssh-keeper/main/install.sh | bash`

## Поддержка пользователей

Создайте раздел в README.md:

````markdown
## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/ssh-keeper/main/install.sh | bash
```
````

### Package Managers

#### macOS (Homebrew)

```bash
brew tap yourusername/ssh-keeper
brew install ssh-keeper
```

#### Ubuntu/Debian

```bash
# Download latest .deb from releases page
sudo dpkg -i ssh-keeper_*.deb
```

```

```

