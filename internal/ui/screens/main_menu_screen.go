package screens

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainMenuScreen представляет главное меню приложения
type MainMenuScreen struct {
	*BaseScreen
	list list.Model
}

// MenuItem представляет элемент меню
type MenuItem struct {
	title       string
	description string
}

// Title возвращает заголовок элемента меню
func (m MenuItem) Title() string {
	return m.title
}

// Description возвращает описание элемента меню
func (m MenuItem) Description() string {
	return m.description
}

// FilterValue возвращает значение для фильтрации
func (m MenuItem) FilterValue() string {
	return m.title
}

// NewMainMenuScreen создает новый экран главного меню
func NewMainMenuScreen() *MainMenuScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Главное меню")

	// Создаем элементы меню
	items := []list.Item{
		MenuItem{
			title:       "Подключения",
			description: "Просмотр и управление SSH подключениями",
		},
		MenuItem{
			title:       "Добавить подключение",
			description: "Создать новое SSH подключение",
		},
		MenuItem{
			title:       "Настройки",
			description: "Настройки приложения",
		},
		MenuItem{
			title:       "Справка",
			description: "Помощь по использованию приложения",
		},
		MenuItem{
			title:       "Выход",
			description: "Закрыть приложение",
		},
	}

	// Создаем список
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Выберите действие"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Margin(0, 0, 1, 0)

	l.Styles.PaginationStyle = lipgloss.NewStyle().
		Margin(1, 0, 0, 0)

	return &MainMenuScreen{
		BaseScreen: baseScreen,
		list:       l,
	}
}

// Update обрабатывает обновления состояния
func (mms *MainMenuScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mms.SetSize(msg.Width, msg.Height)

		mms.list.SetSize(msg.Width-4, msg.Height-16)
		return mms, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return mms, tea.Quit
		case "enter":
			// Обрабатываем выбор элемента
			selectedItem := mms.list.SelectedItem()
			if item, ok := selectedItem.(MenuItem); ok {
				switch item.title {
				case "Выход":
					return mms, tea.Quit
				case "Подключения":
					// TODO: Переход к экрану подключений
					mms.SetContent("Экран подключений будет реализован позже")
				case "Добавить подключение":
					// TODO: Переход к экрану добавления подключения
					mms.SetContent("Экран добавления подключения будет реализован позже")
				case "Настройки":
					// TODO: Переход к экрану настроек
					mms.SetContent("Экран настроек будет реализован позже")
				case "Справка":
					// TODO: Переход к экрану справки
					mms.SetContent("Экран справки будет реализован позже")
				}
			}
		}
	}

	// Обновляем список
	mms.list, cmd = mms.list.Update(msg)

	// Обновляем базовый экран
	baseScreen, baseCmd := mms.BaseScreen.Update(msg)
	mms.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return mms, cmd
}

// View возвращает строку для отрисовки
func (mms *MainMenuScreen) View() string {
	mms.updateContent()
	return mms.BaseScreen.View()
}

// updateContent обновляет содержимое экрана
func (mms *MainMenuScreen) updateContent() {
	// // Рендерим список
	listContent := mms.list.View()

	// // // Добавляем инструкции
	// instructions := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("#808080")).
	// 	Italic(true).
	// 	Align(lipgloss.Center).
	// 	Margin(1, 0, 0, 0).
	// 	MaxWidth(100).
	// 	Render("Используйте ↑/↓ для навигации, Enter для выбора, q для выхода")

	// // // Объединяем список и инструкции
	// content := lipgloss.JoinVertical(lipgloss.Left, listContent, instructions)
	// content := lipgloss.JoinVertical(lipgloss.Left, instructions)
	mms.SetContent(listContent)
}

// Init инициализирует экран
func (mms *MainMenuScreen) Init() tea.Cmd {
	return nil
}
