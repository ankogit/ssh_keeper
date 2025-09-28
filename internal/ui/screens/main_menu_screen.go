package screens

import (
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainMenuScreen представляет главное меню приложения
type MainMenuScreen struct {
	*BaseScreen
	list      list.Model
	config    ui.MenuConfig
	menuItems []ui.MenuItem
}

// NewMainMenuScreen создает новый экран главного меню
func NewMainMenuScreen() *MainMenuScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Главное меню")
	config := ui.DefaultMenuConfig()

	// Создаем элементы меню с действиями
	menuItems := make([]ui.MenuItem, len(config.Items))
	for i, itemConfig := range config.Items {
		menuItems[i] = ui.NewMenuItem(itemConfig)
	}

	// Создаем список
	l := list.New(convertToListItem(menuItems), list.NewDefaultDelegate(), 0, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	l.Styles.PaginationStyle = lipgloss.NewStyle().
		Margin(styles.ListPaginationMargin, 0, 0, 0)

	return &MainMenuScreen{
		BaseScreen: baseScreen,
		list:       l,
		config:     config,
		menuItems:  menuItems,
	}
}

// NewMainMenuScreenWithConfig создает новый экран главного меню с пользовательской конфигурацией
func NewMainMenuScreenWithConfig(config ui.MenuConfig) *MainMenuScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Главное меню")

	// Создаем элементы меню с действиями
	menuItems := make([]ui.MenuItem, len(config.Items))
	for i, itemConfig := range config.Items {
		menuItems[i] = ui.NewMenuItem(itemConfig)
	}

	// Создаем список
	l := list.New(convertToListItem(menuItems), list.NewDefaultDelegate(), 0, 0)
	// l.Title = config.Title
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	// l.Styles.Title = lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color(styles.ColorPrimary)).
	// 	Bold(styles.TextBold).
	// 	Margin(0, 0, styles.ListTitleMargin, 0)

	l.Styles.PaginationStyle = lipgloss.NewStyle().
		Margin(styles.ListPaginationMargin, 0, 0, 0)

	return &MainMenuScreen{
		BaseScreen: baseScreen,
		list:       l,
		config:     config,
		menuItems:  menuItems,
	}
}

// convertToListItem конвертирует MenuItem в list.Item
func convertToListItem(menuItems []ui.MenuItem) []list.Item {
	items := make([]list.Item, len(menuItems))
	for i, item := range menuItems {
		items[i] = item
	}
	return items
}

// Update обрабатывает обновления состояния
func (mms *MainMenuScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mms.SetSize(msg.Width, msg.Height)
		mms.list.SetSize(msg.Width-4, msg.Height-10)
		return mms, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return mms, tea.Quit
		case "enter":
			// Обрабатываем выбор элемента
			selectedItem := mms.list.SelectedItem()
			if item, ok := selectedItem.(ui.MenuItem); ok {
				// Выполняем действие элемента меню
				actionCmd := item.Execute()
				if actionCmd != nil {
					cmd = actionCmd
				}
			}
		default:
			// Проверяем горячие клавиши
			for _, menuItem := range mms.menuItems {
				if menuItem.GetShortcut() == msg.String() {
					actionCmd := menuItem.Execute()
					if actionCmd != nil {
						cmd = actionCmd
					}
					break
				}
			}
		}
	}

	// Обновляем список
	var listCmd tea.Cmd
	mms.list, listCmd = mms.list.Update(msg)
	if listCmd != nil {
		cmd = listCmd
	}

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

	// Создаем стили
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Margin(0, 0, 1, 0)

	// Создаем заголовок
	header := headerStyle.Render("Выберите действие:")

	// // Создаем список подключений
	// var items []string
	// for i, conn := range cs.connections {
	// 	item := itemStyle.Render(lipgloss.NewStyle().
	// 		Foreground(lipgloss.Color("#04B575")).
	// 		Render("• ") + conn)
	// 	items = append(items, item)
	// 	if i < len(cs.connections)-1 {
	// 		items = append(items, "")
	// 	}
	// }

	// // Создаем инструкции
	// instructions := instructionsStyle.Render("Нажмите 'Esq' для возврата к главному меню, 'q' для выхода")

	// // Объединяем все части
	// content := lipgloss.JoinVertical(lipgloss.Left,
	// 	header,
	// 	"",
	// 	lipgloss.JoinVertical(lipgloss.Left, items...),
	// 	"",
	// 	instructions,
	// )
	// // // Добавляем инструкции
	// instructions := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("#808080")).
	// 	Italic(true).
	// 	Align(lipgloss.Center).
	// 	Margin(1, 0, 0, 0).
	// 	MaxWidth(100).
	// 	Render("Используйте ↑/↓ для навигации, Enter для выбора, q для выхода")

	// // // Объединяем список и инструкции
	content := lipgloss.JoinVertical(lipgloss.Left, header, listContent)
	// content := lipgloss.JoinVertical(lipgloss.Left, instructions)
	mms.SetContent(content)
}

// Init инициализирует экран
func (mms *MainMenuScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (mms *MainMenuScreen) GetName() string {
	return "main_menu"
}
