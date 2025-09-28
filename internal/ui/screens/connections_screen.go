package screens

import (
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConnectionsScreen представляет экран управления подключениями
type ConnectionsScreen struct {
	*BaseScreen
	connections []string
}

// NewConnectionsScreen создает новый экран подключений
func NewConnectionsScreen() *ConnectionsScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Подключения")

	// Пример подключений
	connections := []string{
		"server1.example.com (192.168.1.100)",
		"server2.example.com (192.168.1.101)",
		"dev-server.local (10.0.0.50)",
		"prod-server.example.com (203.0.113.10)",
	}

	return &ConnectionsScreen{
		BaseScreen:  baseScreen,
		connections: connections,
	}
}

// Update обрабатывает обновления состояния
func (cs *ConnectionsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cs.SetSize(msg.Width, msg.Height)
		return cs, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return cs, tea.Quit
		case "esc":
			// Возврат к главному меню
			return cs, ui.GoBackCmd()
		}
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
	// Создаем стили
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorPrimary)).
		Bold(styles.TextBold).
		Margin(0, 0, 1, 0)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorText)).
		Margin(0, 0, 0, 2)

	instructionsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted)).
		Italic(styles.TextItalic).
		Margin(1, 0, 0, 0)

	// Создаем заголовок
	header := headerStyle.Render("Доступные SSH подключения:")

	// Создаем список подключений
	var items []string
	for i, conn := range cs.connections {
		item := itemStyle.Render(lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorSecondary)).
			Render("• ") + conn)
		items = append(items, item)
		if i < len(cs.connections)-1 {
			items = append(items, "")
		}
	}

	// Создаем инструкции
	instructions := instructionsStyle.Render("Нажмите 'Esc' для возврата к главному меню, 'q' для выхода")

	// Объединяем все части
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		lipgloss.JoinVertical(lipgloss.Left, items...),
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
