# 🚀 Упрощенная установка SSH Keeper

Теперь установка SSH Keeper стала такой же простой, как у Docker!

## Однострочная установка

### macOS & Linux

```bash
curl -fsSL https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.sh | bash
```

### Windows (PowerShell)

```powershell
iwr -useb https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.ps1 | iex
```

## Что делает скрипт установки

✅ **Автоматически определяет** вашу операционную систему и архитектуру  
✅ **Скачивает** правильную версию для вашей платформы  
✅ **Устанавливает** в системную папку (`/usr/local/bin` на Unix, `~/bin` на Windows)  
✅ **Добавляет в PATH** автоматически  
✅ **Проверяет установку** и показывает статус  
✅ **Красивый вывод** с ASCII-артом и цветными сообщениями

## После установки

```bash
# Запустить SSH Keeper
ssh-keeper

# Проверить версию
ssh-keeper --version

# Показать справку
ssh-keeper --help
```

## Альтернативные способы установки

### Homebrew (macOS)

```bash
brew tap ankogit/ssh-keeper
brew install ssh-keeper
```

### Snap (Linux)

```bash
sudo snap install ssh-keeper
```

### Chocolatey (Windows)

```cmd
choco install ssh-keeper
```

### Scoop (Windows)

```powershell
scoop install ssh-keeper
```

## Преимущества

- 🎯 **Одна команда** - никаких сложных инструкций
- 🔍 **Автоопределение** - скрипт сам знает что скачивать
- 🛡️ **Безопасность** - проверка целостности и подписей
- 🌍 **Кроссплатформенность** - работает везде
- 🎨 **Красивый интерфейс** - как у профессиональных инструментов

---

**Теперь установка SSH Keeper такая же простая, как `curl | bash` у Docker! 🐳**
