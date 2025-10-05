# SSH Keeper - Инструкция по установке

## 🚀 Быстрая установка

### macOS Apple Silicon (M1/M2/M3)

```bash
# Скачать и установить
curl -L -O https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-darwin-arm64.tar.gz
tar -xzf ssh-keeper-0.1.0-darwin-arm64.tar.gz
chmod +x ssh-keeper-0.1.0-darwin-arm64
sudo mv ssh-keeper-0.1.0-darwin-arm64 /usr/local/bin/ssh-keeper
ssh-keeper
```

### macOS Intel

```bash
# Скачать и установить
curl -L -O https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-darwin-amd64.tar.gz
tar -xzf ssh-keeper-0.1.0-darwin-amd64.tar.gz
chmod +x ssh-keeper-0.1.0-darwin-amd64
sudo mv ssh-keeper-0.1.0-darwin-amd64 /usr/local/bin/ssh-keeper
ssh-keeper
```

### Linux

```bash
# Скачать и установить
curl -L -O https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-linux-amd64.tar.gz
tar -xzf ssh-keeper-0.1.0-linux-amd64.tar.gz
chmod +x ssh-keeper-0.1.0-linux-amd64
sudo mv ssh-keeper-0.1.0-linux-amd64 /usr/local/bin/ssh-keeper
ssh-keeper
```

### Windows

```powershell
# Скачать архив
Invoke-WebRequest -Uri "https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-windows-amd64.zip" -OutFile "ssh-keeper-0.1.0-windows-amd64.zip"

# Распаковать архив
Expand-Archive -Path "ssh-keeper-0.1.0-windows-amd64.zip" -DestinationPath "ssh-keeper"

# Перейти в папку и запустить
cd ssh-keeper
.\ssh-keeper-0.1.0-windows-amd64.exe
```

## 🔧 Альтернативная установка (в домашнюю папку)

### macOS/Linux

```bash
# Создать папку для бинарников
mkdir -p ~/bin

# Скачать и установить
curl -L -O https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-darwin-arm64.tar.gz
tar -xzf ssh-keeper-0.1.0-darwin-arm64.tar.gz
mv ssh-keeper-0.1.0-darwin-arm64 ~/bin/ssh-keeper

# Добавить в PATH
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Запустить
ssh-keeper
```

## 📋 Проверка установки

```bash
# Проверить что установлен
which ssh-keeper

# Проверить версию
ssh-keeper --version

# Запустить
ssh-keeper
```

## 🎯 Первый запуск

1. **Установите мастер-пароль** для защиты подключений
2. **Добавьте первое подключение** через "➕ Add Connection"
3. **Настройте подключение**:
   - Имя подключения
   - Хост (IP или домен)
   - Порт (обычно 22)
   - Пользователь
   - Тип аутентификации (пароль/SSH ключ)

## 🚨 Решение проблем

### macOS: "Cannot be opened because it is from an unidentified developer"
```bash
sudo xattr -rd com.apple.quarantine ssh-keeper-0.1.0-darwin-arm64
```

### Linux: "Permission denied"
```bash
chmod +x ssh-keeper-0.1.0-linux-amd64
```

### Windows: "Windows protected your PC"
- Нажмите "More info" → "Run anyway"

## 🗑️ Удаление

### macOS/Linux
```bash
sudo rm /usr/local/bin/ssh-keeper
# или
rm ~/bin/ssh-keeper
```

### Windows
```powershell
Remove-Item "C:\path\to\ssh-keeper-0.1.0-windows-amd64.exe"
```

## 📚 Дополнительная информация

- **Документация**: [README.md](README.md)
- **Релиз**: [GitHub Releases](https://github.com/ankogit/ssh_keeper/releases)
- **Исходный код**: [GitHub Repository](https://github.com/ankogit/ssh_keeper)

---

**Готово! Наслаждайтесь использованием SSH Keeper! 🚀**