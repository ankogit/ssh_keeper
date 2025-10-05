# SSH Keeper - Работа с SSH ключами

## Обзор

SSH Keeper полностью поддерживает различные типы аутентификации SSH, включая использование SSH ключей. Система может работать с:

- **SSH ключами** (приоритетный метод)
- **Паролями** (резервный метод)
- **Комбинированной аутентификацией** (ключ + пароль)
- **Дефолтными ключами** системы

## Типы аутентификации

### 1. 🔑 Аутентификация по SSH ключу

```go
conn := &models.Connection{
    Name:        "Production Server",
    Host:        "prod.example.com",
    Port:        22,
    User:        "admin",
    KeyPath:     "~/.ssh/prod_key",  // Путь к приватному ключу
    HasPassword: false,
}
```

**Особенности:**

- Поддерживаются относительные пути (`~/.ssh/key`)
- Поддерживаются абсолютные пути (`/home/user/.ssh/key`)
- Автоматическое расширение `~` в абсолютный путь
- Совместимость с OpenSSH форматом

### 2. 🔐 Аутентификация по паролю

```go
conn := &models.Connection{
    Name:        "Test Server",
    Host:        "test.example.com",
    Port:        22,
    User:        "user",
    KeyPath:     "",                 // Пустой путь
    HasPassword: true,
    Password:    "secret_password",  // Пароль (зашифрован)
}
```

**Особенности:**

- Пароли автоматически шифруются
- Безопасное хранение в конфиге
- Используется AES-256-GCM шифрование

### 3. 🔑🔐 Комбинированная аутентификация

```go
conn := &models.Connection{
    Name:        "Secure Server",
    Host:        "secure.example.com",
    Port:        22,
    User:        "root",
    KeyPath:     "~/.ssh/secure_key", // Основной ключ
    HasPassword: true,
    Password:    "backup_password",   // Резервный пароль
}
```

**Особенности:**

- Сначала пробуется SSH ключ
- При неудаче запрашивается пароль
- Максимальная безопасность

### 4. 🔑 Дефолтные ключи системы

```go
conn := &models.Connection{
    Name:        "Default Server",
    Host:        "default.example.com",
    Port:        22,
    User:        "ubuntu",
    KeyPath:     "",                 // Пустой путь = дефолтные ключи
    HasPassword: false,
}
```

**Особенности:**

- Использует стандартные ключи (`~/.ssh/id_rsa`, `~/.ssh/id_ed25519`)
- Не требует указания конкретного ключа
- Удобно для стандартных настроек

## Формат конфига OpenSSH

SSH Keeper генерирует конфиг в стандартном формате OpenSSH:

```ssh
# Production Server
Host prod.example.com
    HostName prod.example.com
    User admin
    IdentityFile ~/.ssh/prod_key
    StrictHostKeyChecking ask
    ServerAliveInterval 60
    ServerAliveCountMax 3
    # SSH Keeper ID: 20251003033843
    # Created: 2025-10-03T03:38:43+03:00
    # Updated: 2025-10-03T03:38:43+03:00

# Test Server with Password
Host test.example.com
    HostName test.example.com
    User user
    Password encrypted_password_here
    StrictHostKeyChecking ask
    ServerAliveInterval 60
    ServerAliveCountMax 3
    # SSH Keeper ID: 20251003033844
    # Created: 2025-10-03T03:38:44+03:00
    # Updated: 2025-10-03T03:38:44+03:00
```

## Использование в коде

### Создание подключения с ключом

```go
// Создание подключения
conn := models.NewConnection("My Server", "192.168.1.100", "user")
conn.KeyPath = "~/.ssh/my_key"
conn.HasPassword = false

// Добавление в сервис
service := services.NewConnectionService(configPath, masterKey)
err := service.AddConnection(conn)
```

### Подключение через SSH клиент

```go
// Создание SSH клиента
factory := ssh.NewClientFactory()
sshClient := factory.CreateClient(conn)

// Подключение
err := sshClient.Connect()
```

### Построение SSH аргументов

```go
// В KeyClient.buildSSHArgs()
if kc.connection.KeyPath != "" {
    keyPath, err := filepath.Abs(kc.connection.KeyPath)
    if err == nil {
        args = append(args, "-i", keyPath)
    }
}
```

## Поддерживаемые форматы ключей

SSH Keeper поддерживает все стандартные форматы SSH ключей:

- **RSA** (`~/.ssh/id_rsa`)
- **Ed25519** (`~/.ssh/id_ed25519`)
- **ECDSA** (`~/.ssh/id_ecdsa`)
- **DSA** (`~/.ssh/id_dsa`) - устаревший

## Безопасность

### Шифрование паролей

- Используется AES-256-GCM
- Уникальный nonce для каждого пароля
- Мастер-ключ для расшифровки

### Хранение ключей

- SSH ключи хранятся в исходном виде (не шифруются)
- Пути к ключам сохраняются в конфиге
- Поддержка защищенных паролем ключей

### Права доступа

- Конфиг сохраняется с правами 600
- Только владелец может читать/писать
- Автоматическое создание директории

## Примеры использования

### Импорт существующего SSH конфига

```go
// Импорт из стандартного SSH конфига
err := service.ImportConfig("~/.ssh/config")
```

### Экспорт в SSH конфиг

```go
// Экспорт в стандартный формат
err := service.ExportConfig("~/.ssh/ssh-keeper-config")
```

### Работа с различными серверами

```go
// Продакшн сервер с ключом
prodConn := &models.Connection{
    Name:    "Production",
    Host:    "prod.company.com",
    User:    "deploy",
    KeyPath: "~/.ssh/prod_deploy_key",
}

// Тестовый сервер с паролем
testConn := &models.Connection{
    Name:        "Testing",
    Host:        "test.company.com",
    User:        "tester",
    HasPassword: true,
    Password:    "test123",
}

// Локальный сервер с дефолтными ключами
localConn := &models.Connection{
    Name: "Local Dev",
    Host: "localhost",
    User: "developer",
    // KeyPath пустой = дефолтные ключи
}
```

## Совместимость

### OpenSSH

- Полная совместимость с `ssh`, `scp`, `rsync`
- Поддержка всех стандартных опций
- Импорт/экспорт конфигов

### Другие SSH клиенты

- PuTTY (с конвертацией ключей)
- WinSCP
- VS Code Remote SSH

## Рекомендации

1. **Используйте SSH ключи** вместо паролей для продакшн серверов
2. **Защищайте приватные ключи** паролем
3. **Регулярно ротируйте ключи** для безопасности
4. **Используйте Ed25519** вместо RSA для новых ключей
5. **Храните резервные копии** ключей в безопасном месте



