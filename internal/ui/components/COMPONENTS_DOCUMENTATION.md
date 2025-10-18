# SSH Keeper - Документация по компонентам

## Обзор

Компоненты SSH Keeper представляют собой переиспользуемые UI элементы, построенные на базе Bubble Tea и Lipgloss. Все компоненты находятся в пакете `components` и следуют единому интерфейсу.

## Структура компонентов

```
internal/ui/components/
├── bool_field.go           # Булевое поле для выбора
├── connection_item.go      # Элемент списка подключений
├── field_constants.go      # Константы имен полей
├── form_field.go          # Универсальное поле формы
├── form_manager.go        # Менеджер формы
├── field_navigator.go     # Навигация по полям (устарел)
├── field_styles.go        # Стили полей (устарел)
└── form_renderer.go       # Рендеринг формы (устарел)
```

## Основные компоненты

### 1. FormManager

**Файл:** `form_manager.go`

Центральный менеджер для управления формой и всеми ее полями.

**Основные возможности:**

- Управление коллекцией полей формы
- Навигация между полями (Tab/Shift+Tab)
- Валидация всех полей
- Рендеринг всей формы
- Получение значений полей

**Структура:**

```go
type FormManager struct {
    fields      map[string]*FormField  // Поля формы по имени
    fieldOrder  []string              // Порядок полей
    currentField string               // Текущее активное поле
}
```

**Ключевые методы:**

- `AddField(config FieldConfig)` - добавление поля
- `GetField(name string)` - получение поля по имени
- `NextField()` - переход к следующему полю
- `PrevField()` - переход к предыдущему полю
- `IsLastField()` - проверка, является ли поле последним
- `UpdateFocus()` - обновление фокуса
- `ValidateAll()` - валидация всех полей
- `GetValues()` - получение всех значений
- `RenderForm()` - рендеринг формы

**Пример использования:**

```go
formManager := components.NewFormManager()

// Добавление поля
formManager.AddField(components.FieldConfig{
    Name:        "username",
    Label:       "Имя пользователя",
    Required:    true,
    Width:       50,
    MaxLength:   50,
    Placeholder: "Введите имя",
    FieldType:   components.FieldTypeText,
})

// Навигация
formManager.NextField()

// Валидация
if formManager.ValidateAll() {
    values := formManager.GetValues()
    // обработка значений
}
```

### 2. FormField

**Файл:** `form_field.go`

Универсальный компонент поля формы, поддерживающий различные типы полей.

**Поддерживаемые типы:**

- `FieldTypeText` - текстовое поле
- `FieldTypePassword` - поле пароля
- `FieldTypePort` - поле порта
- `FieldTypeBool` - булевое поле

**Структура:**

```go
type FormField struct {
    config    FieldConfig      // Конфигурация поля
    input     textinput.Model  // Текстовое поле
    boolField *BoolField      // Булевое поле
    value     string          // Текущее значение
    hasError  bool            // Есть ли ошибка валидации
    focused   bool            // Активно ли поле
    visible   bool            // Видимо ли поле
}
```

**Ключевые методы:**

- `Update(msg interface{})` - обновление состояния
- `View()` - рендеринг поля
- `Focus()` - установка фокуса
- `Blur()` - снятие фокуса
- `Value()` - получение значения
- `SetValue(value string)` - установка значения
- `HasError()` - проверка ошибки
- `IsFocused()` - проверка фокуса
- `SetVisible(visible bool)` - установка видимости
- `IsVisible()` - проверка видимости
- `Validate()` - валидация поля
- `GetTextInput()` - получение textinput.Model
- `SetTextInput(input textinput.Model)` - установка textinput.Model

**Конфигурация поля:**

```go
type FieldConfig struct {
    Name        string    // Имя поля
    Label       string    // Лейбл поля
    Required    bool      // Обязательное ли поле
    Width       int       // Ширина поля
    MaxLength   int       // Максимальная длина
    Placeholder string    // Плейсхолдер
    FieldType   FieldType // Тип поля
}
```

### 3. BoolField

**Файл:** `bool_field.go`

Специализированный компонент для булевых значений с визуальным отображением.

**Визуальное представление:**

- `● Да` - значение true (зеленый цвет)
- `○ Нет` - значение false (серый цвет)

**Структура:**

```go
type BoolField struct {
    label  string    // Лейбл поля
    value  bool      // Текущее значение
    focused bool     // Активно ли поле
    width  int       // Ширина поля
    height int       // Высота поля
    styles lipgloss.Style // Стили
}
```

**Ключевые методы:**

- `NewBoolField(label string)` - создание поля
- `Update(msg tea.Msg)` - обновление состояния
- `View()` - рендеринг поля
- `Focus()` - установка фокуса
- `Blur()` - снятие фокуса
- `Value()` - получение значения
- `SetValue(value bool)` - установка значения
- `Toggle()` - переключение значения
- `SetWidth(width int)` - установка ширины

**Обработка клавиш:**

- `←` или `h` - установка false
- `→` или `l` - установка true
- `Space` - переключение значения
- `Enter` - игнорируется (для навигации)

### 4. ConnectionItem

**Файл:** `connection_item.go`

Элемент списка для отображения SSH подключений.

**Структура:**

```go
type ConnectionItem struct {
    Connection models.Connection // Данные подключения
}
```

**Ключевые методы:**

- `NewConnectionItem(conn models.Connection)` - создание элемента
- `Title()` - заголовок элемента
- `Description()` - описание элемента
- `FilterValue()` - значение для фильтрации
- `GetConnection()` - получение подключения
- `RenderCustomItem()` - кастомное отображение

**Отображение:**

```
Production Server (prod.example.com:22) | пользователь: admin | [K]
```

**Иконки аутентификации:**

- [K] - SSH ключ
- [P] - Пароль
- [?] - Неизвестно

## Константы

### FieldConstants

**Файл:** `field_constants.go`

Константы для имен полей формы:

```go
const (
    FieldNameName     = "name"
    FieldNameHost     = "host"
    FieldNamePort     = "port"
    FieldNameUser     = "user"
    FieldNameAuth     = "auth"
    FieldNamePassword = "password"
    FieldNameKey      = "key"
)
```

## Устаревшие компоненты

### FieldNavigator

**Файл:** `field_navigator.go`

**Статус:** Устарел

**Причина:** Функциональность перенесена в `FormManager`

**Замена:** Используйте `FormManager.NextField()`, `FormManager.PrevField()`

### FieldStyles

**Файл:** `field_styles.go`

**Статус:** Устарел

**Причина:** Стили перенесены в `FormField`

**Замена:** Используйте встроенные стили в `FormField`

### FormRenderer

**Файл:** `form_renderer.go`

**Статус:** Устарел

**Причина:** Функциональность перенесена в `FormManager`

**Замена:** Используйте `FormManager.RenderForm()`

## Интеграция с экранами

### Использование в AddConnectionScreen

```go
// Создание менеджера формы
formManager := components.NewFormManager()

// Добавление полей
formManager.AddField(components.FieldConfig{
    Name:        components.FieldNameName,
    Label:       "Название",
    Required:    true,
    Width:       50,
    MaxLength:   50,
    Placeholder: "Название подключения",
    FieldType:   components.FieldTypeText,
})

// Обработка навигации
case "tab":
    formManager.NextField()
case "shift+tab":
    formManager.PrevField()

// Валидация и сохранение
if formManager.ValidateAll() {
    values := formManager.GetValues()
    // создание подключения
}
```

### Использование в ConnectionsScreen

```go
// Создание элементов списка
var listItems []list.Item
for _, conn := range connections {
    listItems = append(listItems, components.NewConnectionItem(conn))
}

// Обработка выбора
if item, ok := selectedItem.(components.ConnectionItem); ok {
    conn := item.GetConnection()
    // подключение к серверу
}
```

## Стилизация

### Цветовая схема

- **Основной текст** - `styles.ColorText`
- **Вторичный текст** - `styles.ColorSecondary`
- **Приглушенный текст** - `styles.ColorMuted`
- **Ошибки** - `styles.ColorError`
- **Предупреждения** - `styles.ColorWarning`
- **Успех** - `styles.ColorSuccess`

### Стили полей

- **Обычное состояние** - синяя рамка
- **Фокус** - оранжевая рамка
- **Ошибка** - красная рамка
- **Булевое поле** - зеленый/серый текст

## Расширение компонентов

### Добавление нового типа поля

1. **Добавить тип** в `FieldType`:

```go
const (
    FieldTypeText     = "text"
    FieldTypePassword = "password"
    FieldTypePort     = "port"
    FieldTypeBool     = "bool"
    FieldTypeNew      = "new"  // новый тип
)
```

2. **Обновить FormField** для поддержки нового типа:

```go
switch ff.config.FieldType {
case FieldTypeNew:
    // обработка нового типа
default:
    // существующие типы
}
```

3. **Добавить валидацию** если необходимо:

```go
func (ff *FormField) Validate() bool {
    switch ff.config.FieldType {
    case FieldTypeNew:
        // валидация нового типа
    default:
        // существующая валидация
    }
}
```

### Создание нового компонента

1. **Создать структуру** компонента
2. **Реализовать методы** `Update`, `View`, `Init`
3. **Добавить стили** если необходимо
4. **Интегрировать** с существующими компонентами

## Примеры использования

### Создание простой формы

```go
// Создание менеджера
formManager := components.NewFormManager()

// Добавление полей
formManager.AddField(components.FieldConfig{
    Name:        "username",
    Label:       "Имя пользователя",
    Required:    true,
    Width:       50,
    FieldType:   components.FieldTypeText,
})

formManager.AddField(components.FieldConfig{
    Name:        "usePassword",
    Label:       "Использовать пароль",
    Required:    true,
    Width:       20,
    FieldType:   components.FieldTypeBool,
})

// Рендеринг формы
formContent := formManager.RenderForm()
```

### Обработка навигации

```go
func (screen *MyScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            screen.formManager.NextField()
        case "shift+tab":
            screen.formManager.PrevField()
        case "enter":
            if screen.formManager.IsLastField() {
                screen.saveForm()
            } else {
                screen.formManager.NextField()
            }
        }
    }

    // Обновление фокуса
    screen.formManager.UpdateFocus()

    return screen, nil
}
```

## Заключение

Система компонентов SSH Keeper обеспечивает:

- **Переиспользование** - компоненты можно использовать в разных экранах
- **Консистентность** - единообразный интерфейс и поведение
- **Расширяемость** - легко добавлять новые типы полей и компонентов
- **Производительность** - эффективная обработка событий и рендеринг
- **Удобство** - простой API для создания форм и интерфейсов

Компоненты спроектированы для максимальной гибкости и простоты использования при сохранении производительности и консистентности интерфейса.
