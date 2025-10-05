package screens

import (
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// App представляет основное приложение SSH Keeper
type App struct {
	*ui.ScreenManager
}

// NewApp создает новое приложение SSH Keeper с менеджером экранов
func NewApp() *App {
	// Создаем менеджер экранов
	manager := ui.NewScreenManager()

	// Создаем экраны
	welcome := NewWelcomeScreen()
	mainMenu := CreateMenuWithActions() // Используем меню с действиями
	connections := NewConnectionsScreen()
	addConnection := NewAddConnectionScreen()
	settings := NewSettingsScreen()
	exportScreen := NewExportScreen()
	importScreen := NewImportScreen()

	// Регистрируем экраны
	manager.RegisterScreen("welcome", welcome)
	manager.RegisterScreen("main_menu", mainMenu)
	manager.RegisterScreen("connections", connections)
	manager.RegisterScreen("add_connection", addConnection)
	manager.RegisterScreen("settings", settings)
	manager.RegisterScreen("export", exportScreen)
	manager.RegisterScreen("import", importScreen)

	// Регистрируем фабрики экранов (для динамического создания)
	manager.RegisterScreenFactory("edit_connection", func() ui.Screen {
		return NewEditConnectionScreenEmpty()
	})

	// Определяем начальный экран на основе состояния мастер-пароля и настроек
	var initialScreen string
	if services.IsMasterPasswordInitializedWithSignature() {
		// Мастер-пароль установлен - проверяем настройку запроса пароля при запуске
		requirePassword, err := services.GetRequirePasswordOnStartupWithSignature()
		if err != nil {
			// Если ошибка получения настройки, используем значение по умолчанию (true)
			requirePassword = true
		}

		if requirePassword {
			// Пользователь хочет вводить пароль при каждом запуске
			initialScreen = "welcome"
		} else {
			// Пользователь не хочет вводить пароль - переходим к главному меню
			initialScreen = "main_menu"
		}
	} else {
		// Мастер-пароль не установлен - показываем экран приветствия
		initialScreen = "welcome"
	}

	// Устанавливаем начальный экран
	manager.SetCurrentScreen(initialScreen)

	return &App{
		ScreenManager: manager,
	}
}

// CreateMenuWithActions создает главное меню с действиями для навигации
func CreateMenuWithActions() *MainMenuScreen {
	// Создаем конфигурацию меню с действиями
	config := ui.MenuConfig{
		Title: "Выберите действие",
		Items: []ui.MenuItemConfig{
			{
				Title:       "Подключения",
				Description: "Просмотр и управление SSH подключениями",
				Shortcut:    "1",
				Action: func() tea.Cmd {
					return ui.NavigateToCmd("connections")
				},
			},
			{
				Title:       "Добавить подключение",
				Description: "Создать новое SSH подключение",
				Shortcut:    "2",
				Action: func() tea.Cmd {
					return ui.NavigateToCmd("add_connection")
				},
			},
			{
				Title:       "Настройки",
				Description: "Настройки приложения",
				Shortcut:    "3",
				Action: func() tea.Cmd {
					return ui.NavigateToCmd("settings")
				},
			},
			{
				Title:       "Справка",
				Description: "Помощь по использованию приложения",
				Shortcut:    "4",
				Action: func() tea.Cmd {
					// TODO: Реализовать экран справки
					return nil
				},
			},
			{
				Title:       "Выход",
				Description: "Закрыть приложение",
				Shortcut:    "ctrl+q",
				Action: func() tea.Cmd {
					return tea.Quit
				},
			},
		},
		ShowBack: false,
		ShowQuit: true,
	}

	return NewMainMenuScreenWithConfig(config)
}

// Update обрабатывает обновления состояния приложения
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return a.ScreenManager.Update(msg)
}

// View возвращает строку для отрисовки
func (a *App) View() string {
	return a.ScreenManager.View()
}

// Init инициализирует приложение
func (a *App) Init() tea.Cmd {
	return a.ScreenManager.Init()
}
