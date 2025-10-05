# Установка SSH Keeper на macOS Apple Silicon

## 🍎 Для macOS с чипом Apple Silicon (M1/M2/M3)

### Способ 1: Быстрая установка (рекомендуется)

```bash
# 1. Скачайте архив для Apple Silicon
curl -L -O https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-darwin-arm64.tar.gz

# 2. Распакуйте архив
tar -xzf ssh-keeper-0.1.0-darwin-arm64.tar.gz

# 3. Сделайте файл исполняемым
chmod +x ssh-keeper-0.1.0-darwin-arm64

# 4. Переместите в системную папку (требует пароль)
sudo mv ssh-keeper-0.1.0-darwin-arm64 /usr/local/bin/ssh-keeper

# 5. Запустите SSH Keeper
ssh-keeper
```

### Способ 2: Установка в домашнюю папку

```bash
# 1. Скачайте архив
curl -L -O https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-darwin-arm64.tar.gz

# 2. Распакуйте архив
tar -xzf ssh-keeper-0.1.0-darwin-arm64.tar.gz

# 3. Создайте папку для бинарников (если не существует)
mkdir -p ~/bin

# 4. Переместите файл
mv ssh-keeper-0.1.0-darwin-arm64 ~/bin/ssh-keeper

# 5. Добавьте ~/bin в PATH (добавьте в ~/.zshrc или ~/.bash_profile)
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc

# 6. Перезагрузите терминал или выполните
source ~/.zshrc

# 7. Запустите SSH Keeper
ssh-keeper
```

### Способ 3: Через Homebrew (если доступно)

```bash
# Если у вас есть Homebrew tap для SSH Keeper
brew install ankogit/ssh-keeper/ssh-keeper
```

## 🔧 Первый запуск

После установки при первом запуске:

1. **Установите мастер-пароль** - это пароль для защиты всех ваших SSH подключений
2. **Добавьте первое подключение** через меню "➕ Add Connection"
3. **Настройте подключение**:
   - Имя подключения (например: "Мой сервер")
   - Хост (IP адрес или домен)
   - Порт (обычно 22)
   - Пользователь
   - Выберите тип аутентификации (пароль или SSH ключ)

## 🎯 Проверка установки

```bash
# Проверьте что SSH Keeper установлен
which ssh-keeper

# Проверьте версию
ssh-keeper --version

# Запустите приложение
ssh-keeper
```

## 🚨 Возможные проблемы

### Проблема: "Permission denied"
```bash
# Решение: сделайте файл исполняемым
chmod +x ssh-keeper-0.1.0-darwin-arm64
```

### Проблема: "Command not found"
```bash
# Решение: добавьте путь к бинарнику в PATH
export PATH="/path/to/ssh-keeper:$PATH"
```

### Проблема: "Cannot be opened because it is from an unidentified developer"
```bash
# Решение: разрешите выполнение в настройках безопасности
# Или выполните команду:
sudo xattr -rd com.apple.quarantine ssh-keeper-0.1.0-darwin-arm64
```

## 🗑️ Удаление

```bash
# Удалите бинарник
sudo rm /usr/local/bin/ssh-keeper

# Или если установлен в домашнюю папку
rm ~/bin/ssh-keeper

# Удалите конфигурацию (опционально)
rm -rf ~/.ssh-keeper
```

## 📱 Использование

После установки SSH Keeper предоставляет красивый TUI интерфейс:

- **🔍 View Connections** - просмотр подключений
- **➕ Add Connection** - добавление нового подключения  
- **⚙️ Settings** - настройки приложения
- **📤 Export** - экспорт в OpenSSH config
- **📥 Import** - импорт из OpenSSH config

### Горячие клавиши:
- `↑/↓` - навигация по меню
- `Enter` - выбор пункта
- `Esc` - назад
- `Q` - выход

---

**Готово! Теперь вы можете управлять SSH подключениями с помощью SSH Keeper! 🚀**
