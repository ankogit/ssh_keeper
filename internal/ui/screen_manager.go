package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Screen представляет интерфейс для всех экранов
type Screen interface {
	tea.Model
	GetName() string
}

// ScreenWithData представляет интерфейс для экранов, которые могут принимать данные
type ScreenWithData interface {
	Screen
	SetData(data interface{})
}

// ScreenFactory функция для создания экранов
type ScreenFactory func() Screen

// ScreenManager управляет переходами между экранами
type ScreenManager struct {
	screens   map[string]Screen
	factories map[string]ScreenFactory
	current   string
	history   []string
	mainMenu  string
	width     int
	height    int
}

// NewScreenManager создает новый менеджер экранов
func NewScreenManager() *ScreenManager {
	return &ScreenManager{
		screens:   make(map[string]Screen),
		factories: make(map[string]ScreenFactory),
		history:   make([]string, 0),
		mainMenu:  "main_menu",
	}
}

// RegisterScreen регистрирует экран в менеджере
func (sm *ScreenManager) RegisterScreen(name string, screen Screen) {
	sm.screens[name] = screen
}

// RegisterScreenFactory регистрирует фабрику экрана в менеджере
func (sm *ScreenManager) RegisterScreenFactory(name string, factory ScreenFactory) {
	sm.factories[name] = factory
}

// SetMainMenu устанавливает главное меню
func (sm *ScreenManager) SetMainMenu(name string) {
	sm.mainMenu = name
}

// NavigateTo переходит к указанному экрану
func (sm *ScreenManager) NavigateTo(screenName string) {
	// Проверяем, есть ли уже экран
	if _, exists := sm.screens[screenName]; exists {
		// Добавляем текущий экран в историю, если он не пустой
		if sm.current != "" {
			sm.history = append(sm.history, sm.current)
		}
		sm.current = screenName
		return
	}

	// Если экрана нет, но есть фабрика, создаем экран
	if factory, exists := sm.factories[screenName]; exists {
		screen := factory()
		sm.screens[screenName] = screen
		// Добавляем текущий экран в историю, если он не пустой
		if sm.current != "" {
			sm.history = append(sm.history, sm.current)
		}
		sm.current = screenName
	}
}

// NavigateToWithData переходит к указанному экрану с данными
func (sm *ScreenManager) NavigateToWithData(screenName string, data interface{}) {
	// Для экранов с данными всегда создаем новый экран
	if factory, exists := sm.factories[screenName]; exists {
		screen := factory()
		// Передаем данные сразу после создания
		if screenWithData, ok := screen.(ScreenWithData); ok && data != nil {
			screenWithData.SetData(data)
		}
		sm.screens[screenName] = screen
		// Добавляем текущий экран в историю, если он не пустой
		if sm.current != "" {
			sm.history = append(sm.history, sm.current)
		}
		sm.current = screenName
		return
	}

	// Если фабрики нет, проверяем существующий экран
	if existingScreen, exists := sm.screens[screenName]; exists {
		// Если экран уже существует, передаем ему данные
		if screenWithData, ok := existingScreen.(ScreenWithData); ok && data != nil {
			screenWithData.SetData(data)
		}
		// Добавляем текущий экран в историю, если он не пустой
		if sm.current != "" {
			sm.history = append(sm.history, sm.current)
		}
		sm.current = screenName
	}
}

// GoBack возвращается к предыдущему экрану
func (sm *ScreenManager) GoBack() {
	if len(sm.history) > 0 {
		// Берем последний экран из истории
		lastIndex := len(sm.history) - 1
		previousScreen := sm.history[lastIndex]

		// Удаляем из истории
		sm.history = sm.history[:lastIndex]

		// Переходим к предыдущему экрану
		sm.current = previousScreen
	} else {
		// Если истории нет, возвращаемся к главному меню
		sm.current = sm.mainMenu
	}
}

// GetCurrentScreen возвращает текущий экран
func (sm *ScreenManager) GetCurrentScreen() Screen {
	if screen, exists := sm.screens[sm.current]; exists {
		return screen
	}
	return nil
}

// GetCurrentScreenName возвращает имя текущего экрана
func (sm *ScreenManager) GetCurrentScreenName() string {
	return sm.current
}

// SetCurrentScreen устанавливает текущий экран
func (sm *ScreenManager) SetCurrentScreen(screenName string) {
	sm.current = screenName
}

// Update обрабатывает обновления состояния менеджера
func (sm *ScreenManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentScreen := sm.GetCurrentScreen()
	if currentScreen == nil {
		return sm, tea.Quit
	}

	// Обрабатываем специальные команды навигации
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Сохраняем размеры экрана
		sm.width = msg.Width
		sm.height = msg.Height
		// Передаем сообщение текущему экрану
		updatedScreen, cmd := currentScreen.Update(msg)
		sm.screens[sm.current] = updatedScreen.(Screen)
		return sm, cmd
	case NavigateToMsg:
		if msg.Data != nil {
			sm.NavigateToWithData(msg.ScreenName, msg.Data)
		} else {
			sm.NavigateTo(msg.ScreenName)
		}
		// Передаем размеры экрана новому экрану, если они есть
		if sm.width > 0 && sm.height > 0 {
			if newScreen := sm.GetCurrentScreen(); newScreen != nil {
				newScreen.Update(tea.WindowSizeMsg{Width: sm.width, Height: sm.height})
			}
		}
		// Передаем NavigateToMsg новому экрану для обработки
		if newScreen := sm.GetCurrentScreen(); newScreen != nil {
			newScreen.Update(msg)
		}
		return sm, nil
	case GoBackMsg:
		sm.GoBack()
		// Передаем размеры экрана новому экрану, если они есть
		if sm.width > 0 && sm.height > 0 {
			if newScreen := sm.GetCurrentScreen(); newScreen != nil {
				newScreen.Update(tea.WindowSizeMsg{Width: sm.width, Height: sm.height})
			}
		}
		// Передаем GoBackMsg новому экрану для обработки
		if newScreen := sm.GetCurrentScreen(); newScreen != nil {
			newScreen.Update(msg)
		}
		return sm, nil
	}

	// Передаем сообщение текущему экрану
	updatedScreen, cmd := currentScreen.Update(msg)

	// Обновляем экран в менеджере
	sm.screens[sm.current] = updatedScreen.(Screen)

	return sm, cmd
}

// View возвращает строку для отрисовки
func (sm *ScreenManager) View() string {
	currentScreen := sm.GetCurrentScreen()
	if currentScreen == nil {
		return "No screen available"
	}
	return currentScreen.View()
}

// Init инициализирует менеджер экранов
func (sm *ScreenManager) Init() tea.Cmd {
	currentScreen := sm.GetCurrentScreen()
	if currentScreen == nil {
		return nil
	}
	return currentScreen.Init()
}

// NavigateToMsg сообщение для навигации к экрану
type NavigateToMsg struct {
	ScreenName string
	Data       interface{} // Данные для передачи в экран
}

// GoBackMsg сообщение для возврата к предыдущему экрану
type GoBackMsg struct{}

// NavigateToCmd создает команду для навигации к экрану
func NavigateToCmd(screenName string) tea.Cmd {
	return func() tea.Msg {
		return NavigateToMsg{ScreenName: screenName}
	}
}

// NavigateToWithDataCmd создает команду для навигации к экрану с данными
func NavigateToWithDataCmd(screenName string, data interface{}) tea.Cmd {
	return func() tea.Msg {
		return NavigateToMsg{ScreenName: screenName, Data: data}
	}
}

// GoBackCmd создает команду для возврата к предыдущему экрану
func GoBackCmd() tea.Cmd {
	return func() tea.Msg {
		return GoBackMsg{}
	}
}
