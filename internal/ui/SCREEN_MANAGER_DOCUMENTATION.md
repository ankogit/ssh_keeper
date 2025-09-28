# Менеджер экранов SSH Keeper

## Обзор

Менеджер экранов предоставляет удобный способ управления переходами между различными экранами приложения. Он поддерживает навигацию, историю переходов и унифицированный интерфейс для всех экранов.

## Основные компоненты

### 1. Интерфейс Screen

Все экраны должны реализовывать интерфейс `Screen`:

```go
type Screen interface {
    tea.Model
    GetName() string
}
```

### 2. ScreenManager

Менеджер экранов управляет переходами и состоянием:

```go
type ScreenManager struct {
    screens    map[string]Screen
    current    string
    history    []string
    mainMenu   string
}
```

### 3. MenuItem с действиями

Элементы меню теперь поддерживают функции выполнения:

```go
type MenuItemConfig struct {
    Title       string
    Description string
    Action      MenuAction
    Shortcut    string
}
```

## Использование

### Создание менеджера экранов

```go
manager := ui.NewScreenManager()

// Регистрация экранов
manager.RegisterScreen("main_menu", mainMenuScreen)
manager.RegisterScreen("connections", connectionsScreen)
manager.RegisterScreen("settings", settingsScreen)

// Установка текущего экрана
manager.SetCurrentScreen("main_menu")
```

### Создание меню с действиями

```go
config := ui.MenuConfig{
    Title: "Выберите действие",
    Items: []ui.MenuItemConfig{
        {
            Title:       "Подключения",
            Description: "Просмотр SSH подключений",
            Shortcut:    "1",
            Action: func() tea.Cmd {
                return ui.NavigateToCmd("connections")
            },
        },
        {
            Title:       "Настройки",
            Description: "Настройки приложения",
            Shortcut:    "2",
            Action: func() tea.Cmd {
                return ui.NavigateToCmd("settings")
            },
        },
    },
}

menu := NewMainMenuScreenWithConfig(config)
```

### Навигация между экранами

```go
// Переход к экрану
cmd := ui.NavigateToCmd("connections")

// Возврат к предыдущему экрану
cmd := ui.GoBackCmd()
```

## Создание нового экрана

### 1. Создайте структуру экрана

```go
type MyScreen struct {
    *BaseScreen
    // дополнительные поля
}
```

### 2. Реализуйте необходимые методы

```go
func (ms *MyScreen) GetName() string {
    return "my_screen"
}

func (ms *MyScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // обработка сообщений
}

func (ms *MyScreen) View() string {
    // отрисовка экрана
}

func (ms *MyScreen) Init() tea.Cmd {
    return nil
}
```

### 3. Зарегистрируйте экран

```go
myScreen := NewMyScreen()
manager.RegisterScreen("my_screen", myScreen)
```

## Примеры экранов

В проекте уже созданы примеры экранов:

- `ConnectionsScreen` - экран управления подключениями
- `SettingsScreen` - экран настроек
- `App` - основное приложение с менеджером экранов

## Горячие клавиши

Меню поддерживает горячие клавиши:

- `1-9` - быстрый переход к пунктам меню
- `q` - выход из приложения
- `Esc` - возврат к предыдущему экрану (в некоторых экранах)

## Сообщения навигации

- `NavigateToMsg` - переход к указанному экрану
- `GoBackMsg` - возврат к предыдущему экрану

## Преимущества

1. **Унификация** - все экраны следуют единому интерфейсу
2. **Навигация** - простое управление переходами между экранами
3. **История** - автоматическое отслеживание истории переходов
4. **Гибкость** - легко добавлять новые экраны и действия
5. **Горячие клавиши** - быстрый доступ к функциям
