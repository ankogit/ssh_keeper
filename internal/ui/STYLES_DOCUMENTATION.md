# Styles Documentation

## Обзор

Файл `styles.go` содержит все стили и цветовую схему для SSH Keeper TUI приложения. Стили построены на основе библиотеки `lipgloss` и обеспечивают единообразный, красивый интерфейс.

## Цветовая схема

### Основные цвета

- **Primary Color** (`#7D56F4`) - Основной фиолетовый цвет для акцентов
- **Secondary Color** (`#04B575`) - Зеленый цвет для успешных действий
- **Accent Color** (`#F25D94`) - Розовый цвет для выделения
- **Warning Color** (`#FFA500`) - Оранжевый для предупреждений
- **Error Color** (`#FF6B6B`) - Красный для ошибок
- **Success Color** (`#51CF66`) - Зеленый для успеха

### Текстовые цвета

- **Text Primary** (`#FFFFFF`) - Основной белый текст
- **Text Secondary** (`#B0B0B0`) - Вторичный серый текст
- **Text Muted** (`#808080`) - Приглушенный серый текст
- **Text Inverse** (`#000000`) - Инверсный черный текст

### Фоновые цвета

- **Background Primary** (`#1A1A1A`) - Основной темный фон
- **Background Secondary** (`#2D2D2D`) - Вторичный фон
- **Background Tertiary** (`#404040`) - Третичный фон
- **Background Highlight** (`#3D3D3D`) - Подсветка

## Стили компонентов

### Основной контейнер

```go
AppStyle = lipgloss.NewStyle().
    Margin(1, 2).
    Padding(1, 2).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(primaryColor).
    Background(bgPrimary).
    Foreground(textPrimary)
```

### Заголовки

- **HeaderStyle** - Стиль для заголовков секций
- **TitleStyle** - Стиль для основных заголовков
- **SubtitleStyle** - Стиль для подзаголовков

### Кнопки

- **ButtonStyle** - Обычное состояние кнопки
- **ButtonHoverStyle** - Состояние при наведении
- **ButtonActiveStyle** - Активное состояние

### Поля ввода

- **InputStyle** - Обычное поле ввода
- **InputFocusedStyle** - Поле ввода в фокусе

### Списки

- **ListStyle** - Контейнер списка
- **ListItemStyle** - Элемент списка
- **ListItemHoverStyle** - Элемент при наведении
- **ListItemSelectedStyle** - Выбранный элемент

### Статусы

- **SuccessStyle** - Успешные сообщения
- **WarningStyle** - Предупреждения
- **ErrorStyle** - Ошибки
- **InfoStyle** - Информационные сообщения

## Адаптивность

### Поддержка тем

- **IsDarkTheme()** - Определяет темную тему терминала
- **GetAdaptiveColor()** - Возвращает адаптивный цвет
- **GetTerminalProfile()** - Получает профиль терминала

### Градиенты

- **CreateGradient()** - Создает градиентный эффект для текста
- **RenderLogo()** - Рендерит логотип с градиентом

## Что нужно доработать

### 1. Поддержка светлой темы

```go
// TODO: Добавить полную поддержку светлой темы
func GetLightThemeColors() ColorScheme {
    return ColorScheme{
        Primary:   "#6B46C1",
        Secondary: "#059669",
        // ... остальные цвета для светлой темы
    }
}
```

### 2. Анимации и переходы

```go
// TODO: Добавить плавные переходы между состояниями
func (s *Style) WithTransition(duration time.Duration) lipgloss.Style {
    // Реализация анимаций
}
```

### 3. Кастомные темы

```go
// TODO: Система пользовательских тем
type Theme struct {
    Name   string
    Colors ColorScheme
}

func LoadCustomTheme(filename string) (*Theme, error) {
    // Загрузка темы из файла
}
```

### 4. Адаптивные размеры

```go
// TODO: Адаптивные размеры для разных размеров терминала
func GetResponsiveStyle(width, height int) lipgloss.Style {
    // Логика адаптации под размер терминала
}
```

### 5. Иконки и символы

```go
// TODO: Расширенная библиотека иконок
const (
    IconConnection = "🔗"
    IconKey        = "🔑"
    IconPassword   = "🔒"
    IconServer     = "🖥️"
    // ... больше иконок
)
```

### 6. Стили для специальных состояний

```go
// TODO: Стили для загрузки, ошибок сети, etc.
var (
    LoadingStyle    = lipgloss.NewStyle().Foreground(warningColor)
    NetworkErrorStyle = lipgloss.NewStyle().Foreground(errorColor)
    DisabledStyle   = lipgloss.NewStyle().Foreground(textMuted)
)
```

### 7. Конфигурируемость

```go
// TODO: Возможность настройки стилей через конфиг
type StyleConfig struct {
    Theme        string `yaml:"theme"`
    FontSize     int    `yaml:"font_size"`
    BorderStyle  string `yaml:"border_style"`
    ColorScheme  string `yaml:"color_scheme"`
}
```

### 8. Производительность

```go
// TODO: Кэширование стилей для лучшей производительности
var styleCache = make(map[string]lipgloss.Style)

func GetCachedStyle(key string, builder func() lipgloss.Style) lipgloss.Style {
    if style, exists := styleCache[key]; exists {
        return style
    }
    style := builder()
    styleCache[key] = style
    return style
}
```

## Использование

### Базовое использование

```go
// Создание стилизованного текста
text := TitleStyle.Render("Заголовок")

// Создание кнопки
button := ButtonStyle.Render("Нажми меня")

// Создание контейнера
container := AppStyle.Render(content)
```

### Адаптивные цвета

```go
// Использование адаптивных цветов
color := GetAdaptiveColor(lightColor, darkColor)
style := lipgloss.NewStyle().Foreground(color)
```

### Градиенты

```go
// Создание градиентного текста
gradientText := CreateGradient("SSH Keeper", startColor, endColor)
```

## Рекомендации

1. **Используйте константы** - Все цвета и стили определены как константы
2. **Следуйте принципам** - Используйте семантические названия стилей
3. **Тестируйте на разных терминалах** - Убедитесь в совместимости
4. **Документируйте изменения** - Обновляйте документацию при изменениях
5. **Оптимизируйте производительность** - Используйте кэширование где возможно
