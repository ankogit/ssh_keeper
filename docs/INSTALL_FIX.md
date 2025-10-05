# ✅ Проблема исправлена!

## Что было не так:

В скрипте установки была ошибка в формировании URL - добавлялся лишний символ `v` в название файла.

**Неправильно:** `ssh-keeper-v0.1.0-darwin-arm64.tar.gz`  
**Правильно:** `ssh-keeper-0.1.0-darwin-arm64.tar.gz`

## Что исправлено:

- ✅ Исправлен URL в `scripts/install.sh`
- ✅ Исправлен URL в `scripts/install.ps1`
- ✅ Разделены переменные `VERSION` (0.1.0) и `VERSION_TAG` (v0.1.0)

## Теперь установка работает!

### macOS & Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.sh | bash
```

### Windows (PowerShell):

```powershell
iwr -useb https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.ps1 | iex
```

## Проверка:

URL теперь правильный: `https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-darwin-arm64.tar.gz`

Попробуйте установку снова! 🚀
