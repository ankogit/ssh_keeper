package screens

import (
	"fmt"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"ssh-keeper/internal/ui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConnectionsScreen представляет экран управления подключениями
type ConnectionsScreen struct {
	*BaseScreen
	list          list.Model
	searchInput   textinput.Model
	connectionSvc *services.ConnectionService
	allItems      []list.Item
}

// NewConnectionsScreen создает новый экран подключений
func NewConnectionsScreen() *ConnectionsScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Подключения")

	// Создаем сервис подключений
	connectionSvc := services.NewConnectionService()
	connections := connectionSvc.GetAllConnections()

	// Создаем элементы списка
	var listItems []list.Item
	for _, conn := range connections {
		listItems = append(listItems, components.NewConnectionItem(conn))
	}

	// Создаем список (компактный, без фильтрации)
	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false) // Отключаем встроенную фильтрацию
	l.SetShowHelp(false)         // Отключаем встроенную справку

	l.Styles.PaginationStyle = lipgloss.NewStyle().
		Margin(1, 0, 0, 0)

	// Создаем input для поиска
	searchInput := textinput.New()
	searchInput.Placeholder = "Поиск подключений..."
	searchInput.Focus()
	searchInput.CharLimit = 25
	searchInput.Width = 40 // Фиксируем ширину при создании

	return &ConnectionsScreen{
		BaseScreen:    baseScreen,
		list:          l,
		searchInput:   searchInput,
		connectionSvc: connectionSvc,
		allItems:      listItems,
	}
}

// Update обрабатывает обновления состояния
func (cs *ConnectionsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cs.SetSize(msg.Width, msg.Height)
		cs.list.SetSize(msg.Width-4, msg.Height-15) // Учитываем место для поиска
		return cs, nil

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			return cs, tea.Quit
		case "esc":
			// Возврат к главному меню
			return cs, ui.GoBackCmd()
		case "enter":
			// Подключиться к выбранному серверу
			fmt.Println("Подключение к выбранному серверу")
			cs.connectToSelected()
		case "ctrl+a":
			// TODO: Добавить новое подключение
		case "ctrl+e":
			fmt.Println("Редактирование выбранного подключения")

			// TODO: Редактировать выбранное подключение
		case "ctrl+d":
			fmt.Println("Удаление выбранного подключения")
			// TODO: Удалить выбранное подключение
		}
	}

	// Обновляем поиск
	cs.searchInput, cmd = cs.searchInput.Update(msg)

	// Фильтруем список при изменении поиска
	cs.filterList()

	// Обновляем список
	var listCmd tea.Cmd
	cs.list, listCmd = cs.list.Update(msg)
	if listCmd != nil {
		cmd = listCmd
	}

	// Обновляем базовый экран
	baseScreen, baseCmd := cs.BaseScreen.Update(msg)
	cs.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return cs, cmd
}

// View возвращает строку для отрисовки
func (cs *ConnectionsScreen) View() string {
	cs.updateContent()
	return cs.BaseScreen.View()
}

// updateContent обновляет содержимое экрана
func (cs *ConnectionsScreen) updateContent() {
	// Стили - принудительно фиксируем размер
	searchStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(styles.ColorWarning)).
		Padding(0, 1) // Отступы внутри

	instructionsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted)).
		Italic(styles.TextItalic).
		Width(80) // Фиксируем ширину для предотвращения переносов

	// Создаем поиск
	searchView := searchStyle.Render(cs.searchInput.View())

	// Получаем содержимое списка
	listContent := cs.list.View()

	// Инструкции - принудительно применяем стиль к каждой строке
	instructionsText := "↑/↓ нав. • Enter подкл. • Ctrl+E ред. • Ctrl+D удал. • Esc назад"
	instructions := instructionsStyle.Render(instructionsText)

	// Объединяем все
	content := lipgloss.JoinVertical(lipgloss.Left,
		searchView,
		"",
		listContent,
		"",
		instructions,
	)

	cs.SetContent(content)
}

// Init инициализирует экран
func (cs *ConnectionsScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (cs *ConnectionsScreen) GetName() string {
	return "connections"
}

// filterList фильтрует список по поисковому запросу
func (cs *ConnectionsScreen) filterList() {
	query := cs.searchInput.Value()
	if query == "" {
		// Показываем все элементы
		cs.list.SetItems(cs.allItems)
		return
	}

	// Фильтруем элементы
	var filteredItems []list.Item
	for _, item := range cs.allItems {
		if connItem, ok := item.(components.ConnectionItem); ok {
			// Поиск по названию, хосту и пользователю
			if strings.Contains(strings.ToLower(connItem.Title()), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(connItem.Description()), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(connItem.FilterValue()), strings.ToLower(query)) {
				filteredItems = append(filteredItems, item)
			}
		}
	}

	cs.list.SetItems(filteredItems)
}

// connectToSelected подключается к выбранному серверу
func (cs *ConnectionsScreen) connectToSelected() {
	selectedItem := cs.list.SelectedItem()
	if item, ok := selectedItem.(components.ConnectionItem); ok {
		conn := item.GetConnection()
		// TODO: Реализовать подключение к SSH
		// Пока просто показываем информацию
		cs.SetContent("Подключение к " + conn.Name + " (" + conn.Host + ")...")
	}
}
