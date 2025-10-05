# CI/CD Pipeline Documentation

## 🚀 GitHub Actions Workflows

### 1. **CI Pipeline** (`.github/workflows/ci.yml`)

**Триггеры:**

- Push в `main` или `develop` ветки
- Pull Request в `main` ветку

**Задачи:**

- ✅ **Тестирование** - запуск тестов на Ubuntu, macOS, Windows
- ✅ **Линтинг** - проверка кода с golangci-lint
- ✅ **Безопасность** - сканирование уязвимостей с Gosec
- ✅ **Сборка** - проверка сборки для всех платформ

**Environment для тестов:**

```bash
DEBUG=true
ENV=development
APP_SIGNATURE=ssh-keeper-sig-test
LOG_LEVEL=debug
```

### 2. **Release Pipeline** (`.github/workflows/release.yml`)

**Триггеры:**

- Push тега `v*` (например, `v0.1.0`)
- Manual dispatch с указанием версии

**Задачи:**

- ✅ **Сборка** - для всех платформ (Linux, macOS Intel/ARM, Windows)
- ✅ **Архивирование** - создание .tar.gz и .zip файлов
- ✅ **Релиз** - автоматическое создание GitHub Release
- ✅ **Артефакты** - загрузка бинарников в релиз

**Environment для production:**

```bash
DEBUG=false
ENV=production
APP_SIGNATURE=ssh-keeper-sig-prod-{version}
LOG_LEVEL=info
```

### 3. **Update Release** (`.github/workflows/update-release.yml`)

**Триггеры:**

- Manual dispatch для обновления существующего релиза

**Задачи:**

- ✅ **Пересборка** - обновление артефактов для указанной версии
- ✅ **Обновление релиза** - замена файлов в существующем релизе

## 🔧 Environment Configuration

### Development (.env)

```bash
DEBUG=true
ENV=development
APP_SIGNATURE=ssh-keeper-sig-dev
LOG_LEVEL=debug
```

### Production (env.production)

```bash
DEBUG=false
ENV=production
APP_SIGNATURE=ssh-keeper-sig-prod
LOG_LEVEL=info
```

### Test (CI)

```bash
DEBUG=true
ENV=development
APP_SIGNATURE=ssh-keeper-sig-test
LOG_LEVEL=debug
```

## 📦 Build Process

### 1. **Автоматическая сборка**

```bash
# Для каждой платформы
GOOS=linux GOARCH=amd64 go build \
  -ldflags "-X main.version=0.1.0" \
  -o ssh-keeper-linux-amd64 ./cmd/ssh-keeper
```

### 2. **Создание архивов**

```bash
# Linux/macOS
tar -czf ssh-keeper-0.1.0-linux-amd64.tar.gz ssh-keeper-linux-amd64

# Windows
zip ssh-keeper-0.1.0-windows-amd64.zip ssh-keeper-windows-amd64.exe
```

### 3. **Загрузка в релиз**

- Автоматическое создание GitHub Release
- Загрузка всех артефактов
- Генерация release notes

## 🎯 Как использовать

### Создание релиза

**Способ 1: Через тег**

```bash
git tag v0.1.0
git push origin v0.1.0
```

**Способ 2: Manual dispatch**

1. Перейти в Actions → "Build and Release"
2. Нажать "Run workflow"
3. Указать версию (например, `v0.1.0`)

### Обновление релиза

1. Перейти в Actions → "Update Release Assets"
2. Нажать "Run workflow"
3. Указать версию для обновления

### Мониторинг

**Статус сборки:**

- Actions tab в GitHub репозитории
- Уведомления в PR
- Статус коммитов

**Артефакты:**

- Releases page в GitHub
- Download links в README
- Автоматические ссылки на скачивание

## 🔒 Security Features

### 1. **Code Scanning**

- Gosec security scanner
- SARIF reports
- GitHub Security tab

### 2. **Dependency Scanning**

- `go mod verify`
- Vulnerability detection
- Automatic updates

### 3. **Build Security**

- Production environment variables
- Secure APP_SIGNATURE generation
- No secrets in logs

## ⚙️ GitHub Environment Setup

### Required Environment Variables

Для работы CI/CD pipeline нужно настроить environment variables в GitHub:

1. **Перейти в Settings → Environments**
2. **Создать/настроить environment "Production"**
3. **Добавить переменную:**

| Variable        | Value                       | Description                           |
| --------------- | --------------------------- | ------------------------------------- |
| `APP_SIGNATURE` | `ssh-keeper-prod-signature` | Production app signature for security |

### Setup Instructions

```bash
# 1. Go to repository Settings
# 2. Navigate to Environments
# 3. Create/Edit "Production" environment
# 4. Add environment variable:
#    Name: APP_SIGNATURE
#    Value: ssh-keeper-prod-signature
# 5. Save changes
```

**Важно:** `APP_SIGNATURE` используется для защиты master password в system keyring. В production это должно быть уникальное значение.

## 📊 Pipeline Status

| Workflow       | Status    | Triggers | Duration |
| -------------- | --------- | -------- | -------- |
| CI             | ✅ Active | Push/PR  | ~5 min   |
| Release        | ✅ Active | Tags     | ~10 min  |
| Update Release | ✅ Active | Manual   | ~8 min   |

## 🛠️ Troubleshooting

### Common Issues

**1. Build failures**

- Проверить Go version compatibility
- Убедиться что все зависимости доступны
- Проверить .env файл

**2. Release failures**

- Проверить права доступа к репозиторию
- Убедиться что тег существует
- Проверить GitHub token

**3. Test failures**

- Проверить тестовые данные
- Убедиться что все сервисы доступны
- Проверить environment variables

### Debug Commands

```bash
# Локальная сборка
make build-all

# Тестирование
make test

# Проверка линтера
make lint

# Проверка релиза
make release
```

---

**Pipeline готов к использованию! 🚀**
