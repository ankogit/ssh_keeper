#!/bin/bash

# Скрипт для создания Pull Request в Homebrew Core

set -e

# Конфигурация
HOMEBREW_CORE_REPO="https://github.com/Homebrew/homebrew-core.git"
FORMULA_NAME="ssh-keeper"
MAIN_REPO="yourusername/ssh-keeper"
VERSION="v0.1.0"

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 Подготовка Pull Request для Homebrew Core${NC}"
echo "=============================================="

# Проверяем, что мы в правильной директории
if [ ! -f "Formula/ssh-keeper-core.rb" ]; then
    echo -e "${RED}❌ Формула ssh-keeper-core.rb не найдена${NC}"
    exit 1
fi

# Создаем рабочую директорию
WORK_DIR=$(mktemp -d)
echo -e "${YELLOW}📁 Рабочая директория: $WORK_DIR${NC}"

# Клонируем Homebrew Core
echo -e "${BLUE}📥 Клонирование Homebrew Core...${NC}"
cd "$WORK_DIR"
git clone "$HOMEBREW_CORE_REPO"
cd homebrew-core

# Создаем новую ветку
BRANCH_NAME="ssh-keeper"
echo -e "${BLUE}🌿 Создание ветки: $BRANCH_NAME${NC}"
git checkout -b "$BRANCH_NAME"

# Копируем формулу
echo -e "${BLUE}📋 Копирование формулы...${NC}"
cp "/Users/anko/work/ssh_keeper/Formula/ssh-keeper-core.rb" "Formula/$FORMULA_NAME.rb"

# Обновляем формулу с правильным SHA256
echo -e "${YELLOW}🔐 Вычисление SHA256...${NC}"
ARCHIVE_URL="https://github.com/$MAIN_REPO/archive/$VERSION.tar.gz"
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"
curl -L "$ARCHIVE_URL" -o "archive.tar.gz"
SHA256=$(shasum -a 256 "archive.tar.gz" | cut -d' ' -f1)
cd "$WORK_DIR/homebrew-core"

# Обновляем SHA256 в формуле
sed -i.bak "s|sha256 \"abc123...\"|sha256 \"$SHA256\"|g" "Formula/$FORMULA_NAME.rb"

# Проверяем формулу
echo -e "${BLUE}🔍 Проверка формулы...${NC}"
brew audit --new-formula "Formula/$FORMULA_NAME.rb" || true

# Добавляем файлы
git add "Formula/$FORMULA_NAME.rb"
git commit -m "Add ssh-keeper formula

- Add SSH Keeper CLI tool for managing SSH connections
- Beautiful TUI interface built with Go
- Secure password storage and OpenSSH compatibility
- Homepage: https://github.com/$MAIN_REPO
- License: MIT"

echo -e "${GREEN}✅ Commit создан!${NC}"

# Показываем следующее шаги
echo ""
echo -e "${BLUE}📋 Следующие шаги:${NC}"
echo "1. Создайте форк Homebrew Core на GitHub"
echo "2. Добавьте ваш форк как remote:"
echo "   git remote add fork https://github.com/YOUR_USERNAME/homebrew-core.git"
echo "3. Push ветку в ваш форк:"
echo "   git push fork $BRANCH_NAME"
echo "4. Создайте Pull Request на GitHub"
echo ""
echo -e "${YELLOW}⚠️  Важные требования для Homebrew Core:${NC}"
echo "- Минимум 20 звезд на GitHub"
echo "- Стабильные релизы (не pre-release)"
echo "- Активная поддержка проекта"
echo "- Соответствие стандартам Homebrew"
echo "- Прохождение всех тестов"

# Очистка
cd /
rm -rf "$WORK_DIR"
rm -rf "$TEMP_DIR"

echo -e "${GREEN}🎉 Подготовка завершена!${NC}"



