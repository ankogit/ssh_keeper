#!/bin/bash

# Скрипт для обновления Homebrew формулы после релиза

set -e

# Конфигурация
TAP_REPO="../homebrew-ssh-keeper"
MAIN_REPO="yourusername/ssh-keeper"
VERSION=""

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

usage() {
    echo "Использование: $0 <version>"
    echo "Пример: $0 v0.1.0"
    exit 1
}

if [ $# -eq 0 ]; then
    usage
fi

VERSION=$1

echo -e "${BLUE}🔄 Обновление Homebrew формулы для версии $VERSION${NC}"
echo "================================================"

# Проверяем, что tap репозиторий существует
if [ ! -d "$TAP_REPO" ]; then
    echo -e "${RED}❌ Tap репозиторий не найден: $TAP_REPO${NC}"
    echo -e "${YELLOW}Сначала запустите: ./scripts/create-homebrew-tap.sh${NC}"
    exit 1
fi

# Переходим в tap репозиторий
cd "$TAP_REPO"

# Проверяем, что формула существует
if [ ! -f "ssh-keeper.rb" ]; then
    echo -e "${RED}❌ Формула ssh-keeper.rb не найдена${NC}"
    exit 1
fi

echo -e "${YELLOW}📥 Скачивание архива для вычисления SHA256...${NC}"

# Создаем временную директорию
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Скачиваем архив
ARCHIVE_URL="https://github.com/$MAIN_REPO/archive/$VERSION.tar.gz"
ARCHIVE_NAME="ssh-keeper-$VERSION.tar.gz"

echo -e "${BLUE}Скачивание: $ARCHIVE_URL${NC}"
curl -L "$ARCHIVE_URL" -o "$ARCHIVE_NAME"

# Вычисляем SHA256
SHA256=$(shasum -a 256 "$ARCHIVE_NAME" | cut -d' ' -f1)
echo -e "${GREEN}✅ SHA256: $SHA256${NC}"

# Очищаем временную директорию
cd /
rm -rf "$TEMP_DIR"

# Возвращаемся в tap репозиторий
cd "$TAP_REPO"

echo -e "${YELLOW}📝 Обновление формулы...${NC}"

# Создаем резервную копию
cp ssh-keeper.rb ssh-keeper.rb.bak

# Обновляем версию и SHA256 в формуле
sed -i.tmp "s|url \"https://github.com/$MAIN_REPO/archive/v[^\"]*\.tar\.gz\"|url \"https://github.com/$MAIN_REPO/archive/$VERSION.tar.gz\"|g" ssh-keeper.rb
sed -i.tmp "s|sha256 \"[^\"]*\"|sha256 \"$SHA256\"|g" ssh-keeper.rb

# Удаляем временный файл
rm -f ssh-keeper.rb.tmp

echo -e "${GREEN}✅ Формула обновлена!${NC}"

# Показываем изменения
echo -e "${BLUE}📋 Изменения в формуле:${NC}"
echo "----------------------------------------"
diff ssh-keeper.rb.bak ssh-keeper.rb || true

# Предлагаем сделать commit
echo ""
echo -e "${YELLOW}📦 Хотите сделать commit и push? (y/N):${NC}"
read -r response

if [[ "$response" =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}📝 Создание commit...${NC}"
    git add ssh-keeper.rb
    git commit -m "Update ssh-keeper to $VERSION"
    
    echo -e "${BLUE}🚀 Push в GitHub...${NC}"
    git push origin main
    
    echo -e "${GREEN}✅ Формула обновлена и опубликована!${NC}"
    echo ""
    echo -e "${BLUE}🎉 Пользователи могут теперь обновить:${NC}"
    echo "brew tap yourusername/ssh-keeper"
    echo "brew upgrade ssh-keeper"
else
    echo -e "${YELLOW}⚠️  Изменения сохранены, но не закоммичены${NC}"
    echo -e "${BLUE}Для публикации выполните:${NC}"
    echo "cd $TAP_REPO"
    echo "git add ssh-keeper.rb"
    echo "git commit -m \"Update ssh-keeper to $VERSION\""
    echo "git push origin main"
fi

# Удаляем резервную копию
rm -f ssh-keeper.rb.bak



