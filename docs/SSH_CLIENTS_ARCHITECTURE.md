# SSH Clients Architecture

## Обзор

Проект теперь использует раздельную архитектуру SSH клиентов для разных типов аутентификации:

- **KeyClient** - для аутентификации по SSH ключу
- **PasswordClient** - для аутентификации по паролю
- **ClientFactory** - фабрика для создания соответствующих клиентов

## Проблемы, которые были решены

### 1. Проблема с вставкой текста в nano (символы 0~ и 1~)

**Проблема**: При вставке текста в nano после SSH подключения появлялись лишние символы `0~` и `1~`.

**Решение**: Обновлена функция `restoreTerminal()` в `ssh_session_screen.go` для полного восстановления терминала:

```go
func restoreTerminal() {
    os.Stdout.WriteString("\033[?1049l") // Выход из альтернативного буфера
    os.Stdout.WriteString("\033[?25h")   // Показать курсор
    os.Stdout.WriteString("\033[?2004l") // Отключаем bracketed paste mode
    os.Stdout.WriteString("\033[?1l")    // Отключаем application cursor keys
    os.Stdout.WriteString("\033[?7h")    // Включаем auto wrap mode
    os.Stdout.WriteString("\033[?12l")   // Отключаем local echo
    os.Stdout.WriteString("\033[?25h")   // Показываем курсор
    os.Stdout.WriteString("\033[0m")     // Сбрасываем все атрибуты
}
```

### 2. Разделение SSH клиентов по типу аутентификации

**Проблема**: Один универсальный клиент для всех типов аутентификации был сложен в поддержке.

**Решение**: Созданы специализированные клиенты:

## Архитектура клиентов

### SSHClientInterface

```go
type SSHClientInterface interface {
    Connect() error
    TestConnection() error
    GetConnectionString() string
}
```

### KeyClient (`internal/ssh/key_client.go`)

Специализированный клиент для аутентификации по SSH ключу:

- Использует только публичные ключи
- Отключает аутентификацию по паролю
- Автоматически находит дефолтные ключи в `~/.ssh/`
- Поддерживает явно указанные ключи

**Особенности**:

- `PreferredAuthentications=publickey`
- `PubkeyAuthentication=yes`
- `PasswordAuthentication=no`

### PasswordClient (`internal/ssh/password_client.go`)

Специализированный клиент для аутентификации по паролю:

- Использует только пароли
- Отключает аутентификацию по ключу
- Поддерживает автоматическую передачу пароля через PTY
- Опционально поддерживает `sshpass`

**Особенности**:

- `PreferredAuthentications=password`
- `PubkeyAuthentication=no`
- `PasswordAuthentication=yes`

### ClientFactory (`internal/ssh/client_factory.go`)

Фабрика для создания соответствующих клиентов:

```go
func (cf *ClientFactory) CreateClient(conn *models.Connection) SSHClientInterface {
    if conn.HasPassword {
        return NewPasswordClient(conn)
    }
    return NewKeyClient(conn)
}
```

## Использование

### В SSH Session Screen

```go
// Создаем соответствующий SSH клиент на основе типа аутентификации
factory := ssh.NewClientFactory()
sshClient := factory.CreateClient(conn)

// Если это клиент с паролем, устанавливаем пароль
if passwordClient, ok := sshClient.(*ssh.PasswordClient); ok {
    passwordClient.SetPassword(password)
}
```

## Преимущества новой архитектуры

1. **Специализация**: Каждый клиент оптимизирован для своего типа аутентификации
2. **Безопасность**: Четкое разделение методов аутентификации
3. **Простота**: Упрощенная логика для каждого типа подключения
4. **Расширяемость**: Легко добавлять новые типы аутентификации
5. **Тестируемость**: Каждый клиент можно тестировать независимо
6. **Производительность**: Буферизованное копирование потоков с оптимальным размером буфера
7. **Надежность**: Улучшенная обработка ошибок и корректное завершение горутин

## Новые улучшения PTY

### Улучшенное копирование потоков

Используется буферизованное копирование через горутины с улучшенной обработкой ошибок:

```go
// Эффективный подход с буферизацией
go pty.CopyToStdout()        // Буферизованное копирование из PTY в stdout
go pty.CopyFromStdinWithInterrupt() // Буферизованное копирование из stdin в PTY
```

**Преимущества**:

- Буферизованное чтение/запись (4KB буфер)
- Корректная обработка ошибок без паники
- Неблокирующее выполнение
- Поддержка интерактивных приложений

### Методы PTY

- `StartSSH()` - запуск SSH с промежуточным PTY
- `StartSSHDirect()` - запуск SSH с прямым подключением к стандартным потокам
- `ConnectDirectly()` - подключение через горутины (для совместимости)
- `ConnectDirectlySync()` - синхронное подключение с автоматическим закрытием

### Преимущества прямого подключения

1. **Производительность**: Нет промежуточных буферов и горутин
2. **Надежность**: Прямая передача данных без потерь
3. **Простота**: Меньше кода и потенциальных точек отказа
4. **Совместимость**: Лучшая работа с интерактивными приложениями (nano, vim, etc.)

## Файлы

- `internal/ssh/key_client.go` - клиент для SSH ключей
- `internal/ssh/password_client.go` - клиент для паролей
- `internal/ssh/client_factory.go` - фабрика клиентов
- `internal/ssh/pty.go` - улучшенная реализация PTY с прямым подключением
- `internal/ui/screens/ssh_session_screen.go` - обновленный экран сессии
