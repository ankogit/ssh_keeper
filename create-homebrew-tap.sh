#!/bin/bash

# Скрипт для создания Homebrew Tap репозитория

set -e

# Конфигурация
TAP_NAME="homebrew-ssh-keeper"
MAIN_REPO="yourusername/ssh-keeper"
GITHUB_USERNAME="yourusername"

echo "🚀 Создание Homebrew Tap для ssh-keeper"
echo "========================================"

# Проверяем, что мы в правильной директории
if [ ! -f "Formula/ssh-keeper.rb" ]; then
    echo "❌ Ошибка: Formula/ssh-keeper.rb не найден"
    echo "Запустите этот скрипт из корневой директории проекта"
    exit 1
fi

# Создаем директорию для tap
echo "📁 Создание директории для tap..."
mkdir -p "../$TAP_NAME"
cd "../$TAP_NAME"

# Инициализируем git репозиторий
echo "🔧 Инициализация git репозитория..."
git init
git branch -M main

# Копируем формулу
echo "📋 Копирование формулы..."
cp "../ssh_keeper/Formula/ssh-keeper.rb" "ssh-keeper.rb"

# Обновляем формулу с правильным URL
echo "🔗 Обновление URL в формуле..."
sed -i.bak "s|homepage \"https://github.com/yourusername/ssh-keeper\"|homepage \"https://github.com/$MAIN_REPO\"|g" ssh-keeper.rb
sed -i.bak "s|url \"https://github.com/yourusername/ssh-keeper/archive/v0.1.0.tar.gz\"|url \"https://github.com/$MAIN_REPO/archive/v0.1.0.tar.gz\"|g" ssh-keeper.rb

# Создаем README для tap
echo "📝 Создание README..."
cat > README.md << EOF
# Homebrew Tap for SSH Keeper

This tap provides the SSH Keeper CLI tool for managing SSH connections.

## Installation

\`\`\`bash
# Add this tap
brew tap $GITHUB_USERNAME/ssh-keeper

# Install ssh-keeper
brew install ssh-keeper
\`\`\`

## What is SSH Keeper?

SSH Keeper is a beautiful and secure CLI tool for managing SSH connections with a modern TUI interface.

## Features

- 🔍 Browse and search SSH connections
- ➕ Add new SSH connections easily  
- 🔐 Secure password storage
- 📤 Export/Import OpenSSH config compatibility
- 🎨 Beautiful UI with colors and animations
- ⚡ Fast and lightweight Go implementation

## Documentation

For more information, visit: https://github.com/$MAIN_REPO

## License

MIT License - see https://github.com/$MAIN_REPO/blob/main/LICENSE
EOF

# Создаем .gitignore
echo "📄 Создание .gitignore..."
cat > .gitignore << EOF
# macOS
.DS_Store

# Backup files
*.bak
EOF

# Добавляем файлы в git
echo "📦 Добавление файлов в git..."
git add .
git commit -m "Initial commit: Add ssh-keeper formula"

echo ""
echo "✅ Tap репозиторий создан!"
echo ""
echo "📋 Следующие шаги:"
echo "1. Создайте репозиторий на GitHub: https://github.com/new"
echo "   - Название: $TAP_NAME"
echo "   - Описание: Homebrew tap for SSH Keeper"
echo "   - Сделайте публичным"
echo ""
echo "2. Подключите удаленный репозиторий:"
echo "   cd ../$TAP_NAME"
echo "   git remote add origin https://github.com/$GITHUB_USERNAME/$TAP_NAME.git"
echo "   git push -u origin main"
echo ""
echo "3. После создания релиза в основном репозитории:"
echo "   - Обновите SHA256 в ssh-keeper.rb"
echo "   - Обновите версию в url"
echo "   - Сделайте commit и push"
echo ""
echo "4. Пользователи смогут устанавливать:"
echo "   brew tap $GITHUB_USERNAME/ssh-keeper"
echo "   brew install ssh-keeper"

