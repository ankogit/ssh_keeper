package screens

import (
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SettingsScreen представляет экран настроек
type SettingsScreen struct {
	*BaseScreen
	settings map[string]string
}

// NewSettingsScreen создает новый экран настроек
func NewSettingsScreen() *SettingsScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Настройки")

	// Пример настроек
	settings := map[string]string{
		"Master Key Timeout": "1 час",
		"SSH Path":           "/usr/bin/ssh",
		"Export Format":      "OpenSSH",
		"Default Port":       "22",
		"Theme":              "default",
	}

	return &SettingsScreen{
		BaseScreen: baseScreen,
		settings:   settings,
	}
}

// Update обрабатывает обновления состояния
func (ss *SettingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ss.SetSize(msg.Width, msg.Height)
		return ss, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return ss, tea.Quit
		case "esc":
			// Возврат к главному меню
			return ss, ui.GoBackCmd()
		}
	}

	// Обновляем базовый экран
	baseScreen, baseCmd := ss.BaseScreen.Update(msg)
	ss.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return ss, cmd
}

// View возвращает строку для отрисовки
func (ss *SettingsScreen) View() string {
	ss.updateContent()
	return ss.BaseScreen.View()
}

// updateContent обновляет содержимое экрана
func (ss *SettingsScreen) updateContent() {
	// Создаем стили
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorPrimary)).
		Bold(styles.TextBold).
		Margin(0, 0, 1, 0)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorSecondary)).
		Bold(styles.TextBold).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorText))

	instructionsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted)).
		Italic(styles.TextItalic).
		Margin(1, 0, 0, 0)

	// Создаем заголовок
	header := headerStyle.Render("Настройки приложения:")

	// Создаем список настроек
	var settingsList []string
	for key, value := range ss.settings {
		setting := keyStyle.Render(key+":") + " " + valueStyle.Render(value)
		settingsList = append(settingsList, setting)
	}

	// Создаем инструкция и
	instructions := instructionsStyle.Render("Нажмите 'Esc' для возврата к главному меню, 'q' для выхода")

	// Объединяем все части
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		lipgloss.JoinVertical(lipgloss.Left, settingsList...),
		"",
		instructions,
	)

	ss.SetContent(content)
}

// Init инициализирует экран
func (ss *SettingsScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (ss *SettingsScreen) GetName() string {
	return "settings"
}
