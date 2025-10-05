#!/bin/bash

# Скрипт для проверки формулы перед отправкой в Homebrew Core

set -e

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

FORMULA_FILE="Formula/ssh-keeper-core.rb"

echo -e "${BLUE}🔍 Проверка формулы для Homebrew Core${NC}"
echo "========================================"

# Проверяем, что формула существует
if [ ! -f "$FORMULA_FILE" ]; then
    echo -e "${RED}❌ Формула $FORMULA_FILE не найдена${NC}"
    exit 1
fi

echo -e "${YELLOW}📋 Проверка синтаксиса Ruby...${NC}"
if ruby -c "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ Синтаксис Ruby корректен${NC}"
else
    echo -e "${RED}❌ Ошибка в синтаксисе Ruby${NC}"
    exit 1
fi

echo -e "${YELLOW}🔍 Проверка стиля формулы...${NC}"
if brew style "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ Стиль формулы корректен${NC}"
else
    echo -e "${RED}❌ Ошибки в стиле формулы${NC}"
    echo -e "${YELLOW}Исправьте ошибки стиля и попробуйте снова${NC}"
fi

echo -e "${YELLOW}🔍 Проверка формулы (audit)...${NC}"
if brew audit --new-formula "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ Формула прошла audit${NC}"
else
    echo -e "${RED}❌ Формула не прошла audit${NC}"
    echo -e "${YELLOW}Исправьте ошибки и попробуйте снова${NC}"
fi

echo -e "${YELLOW}🧪 Тестирование формулы...${NC}"
if brew test "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ Тесты прошли успешно${NC}"
else
    echo -e "${RED}❌ Тесты не прошли${NC}"
    echo -e "${YELLOW}Исправьте ошибки в тестах${NC}"
fi

echo ""
echo -e "${BLUE}📊 Результаты проверки:${NC}"
echo "=========================="

# Проверяем основные требования
echo -e "${YELLOW}📋 Основные требования:${NC}"

# Проверяем наличие описания
if grep -q 'desc "' "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ Описание присутствует${NC}"
else
    echo -e "${RED}❌ Отсутствует описание${NC}"
fi

# Проверяем наличие homepage
if grep -q 'homepage "' "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ Homepage указан${NC}"
else
    echo -e "${RED}❌ Отсутствует homepage${NC}"
fi

# Проверяем наличие URL
if grep -q 'url "' "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ URL указан${NC}"
else
    echo -e "${RED}❌ Отсутствует URL${NC}"
fi

# Проверяем наличие SHA256
if grep -q 'sha256 "' "$FORMULA_FILE" && ! grep -q 'sha256 ""' "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ SHA256 указан${NC}"
else
    echo -e "${RED}❌ Отсутствует или пустой SHA256${NC}"
fi

# Проверяем наличие лицензии
if grep -q 'license "' "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ Лицензия указана${NC}"
else
    echo -e "${RED}❌ Отсутствует лицензия${NC}"
fi

# Проверяем наличие тестов
if grep -q 'test do' "$FORMULA_FILE"; then
    echo -e "${GREEN}✅ Тесты присутствуют${NC}"
else
    echo -e "${RED}❌ Отсутствуют тесты${NC}"
fi

echo ""
echo -e "${BLUE}🎯 Рекомендации:${NC}"
echo "=================="
echo "1. Убедитесь, что у проекта минимум 20 звезд на GitHub"
echo "2. Создайте стабильный релиз (не pre-release)"
echo "3. Проверьте, что все зависимости минимальны"
echo "4. Убедитесь, что тесты проходят"
echo "5. Документируйте все изменения"

echo ""
echo -e "${GREEN}🎉 Проверка завершена!${NC}"



