# SSH Keeper - Конфигурация и шифрование

## Обзор

SSH Keeper теперь поддерживает сохранение и загрузку подключений из файла конфигурации в формате OpenSSH. Все пароли автоматически шифруются для безопасности.

## Основные возможности

### 1. Формат конфига OpenSSH

- Совместимость с стандартным форматом SSH конфигурации
- Поддержка всех основных SSH опций
- Возможность импорта/экспорта конфигов

### 2. Шифрование паролей

- Все пароли автоматически шифруются при сохранении
- Используется AES-GCM шифрование
- Мастер-ключ для расшифровки

### 3. Автоматическое сохранение

- Все изменения автоматически сохраняются в файл
- Поддержка импорта из существующих конфигов
- Экспорт в стандартный формат OpenSSH

## Структура файлов

```
internal/
├── models/
│   ├── connection.go      # Модель подключения
│   ├── config.go         # Основная конфигурация
│   └── ssh_config.go     # Модели для SSH конфига
├── services/
│   ├── connection_service.go    # Основной сервис подключений
│   ├── ssh_config_service.go    # Сервис для работы с SSH конфигом
│   ├── encryption_service.go    # Сервис шифрования
│   └── global.go               # Глобальные сервисы
```

## Использование

### Инициализация

```go
// Создание сервиса с путем к конфигу и мастер-ключом
configPath := "~/.ssh-keeper/config"
masterKey := "your-secure-master-key"
service := services.NewConnectionService(configPath, masterKey)
```

### Добавление подключения

```go
conn := models.NewConnection("My Server", "192.168.1.100", "user")
conn.Password = "secret_password"
conn.HasPassword = true

err := service.AddConnection(conn) // Автоматически сохраняется
```

### Импорт конфига

```go
err := service.ImportConfig("/path/to/existing/ssh/config")
```

### Экспорт конфига

```go
err := service.ExportConfig("/path/to/export/config")
```

## Формат конфига

```ssh
# SSH Keeper Configuration File
# Generated on 2024-01-15T10:30:00Z
# Version: 1.0

# Global Settings
ServerAliveInterval 60
ServerAliveCountMax 3

# My Server
Host my-server
    HostName 192.168.1.100
    Port 22
    User user
    Password encrypted_password_here
    StrictHostKeyChecking ask
    # SSH Keeper ID: unique_id
    # Created: 2024-01-15T10:30:00Z
    # Updated: 2024-01-15T10:30:00Z
```

## Безопасность

### Шифрование паролей

- Используется AES-256-GCM
- Мастер-ключ выводится из SHA-256 хеша
- Каждый пароль имеет уникальный nonce

### Хранение конфига

- Конфиг сохраняется в `~/.ssh-keeper/config`
- Права доступа: 600 (только владелец)
- Автоматическое создание директории

## Совместимость

### Импорт из OpenSSH

- Поддержка стандартного формата SSH конфига
- Автоматическое определение зашифрованных паролей
- Сохранение всех SSH опций

### Экспорт в OpenSSH

- Генерация стандартного SSH конфига
- Совместимость с ssh, scp, rsync
- Сохранение метаданных SSH Keeper

## Примеры

### Создание подключения с паролем

```go
conn := &models.Connection{
    Name:        "Production Server",
    Host:        "prod.example.com",
    Port:        22,
    User:        "admin",
    Password:    "secure_password",
    HasPassword: true,
}
service.AddConnection(conn)
```

### Создание подключения с ключом

```go
conn := &models.Connection{
    Name:        "Staging Server",
    Host:        "staging.example.com",
    Port:        22,
    User:        "deploy",
    KeyPath:     "~/.ssh/staging_key",
    HasPassword: false,
}
service.AddConnection(conn)
```

## Будущие улучшения

1. **Улучшенная безопасность**

   - Поддержка аппаратных токенов
   - Интеграция с системными хранилищами ключей

2. **Расширенная совместимость**

   - Поддержка PuTTY конфигов
   - Импорт из других SSH менеджеров

3. **Облачная синхронизация**
   - Синхронизация между устройствами
   - Облачное резервное копирование



