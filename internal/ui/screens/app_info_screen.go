package screens

import (
	"fmt"
	"runtime"

	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/styles"
	"ssh-keeper/internal/version"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppInfoScreen представляет экран информации о приложении
type AppInfoScreen struct {
	*BaseScreen
	viewport viewport.Model
}

// NewAppInfoScreen создает новый экран информации о приложении
func NewAppInfoScreen() *AppInfoScreen {
	baseScreen := NewBaseScreen("SSH Keeper - О приложении")

	// Создаем viewport для прокрутки
	vp := viewport.New(0, 0)

	return &AppInfoScreen{
		BaseScreen: baseScreen,
		viewport:   vp,
	}
}

// Update обрабатывает обновления состояния
func (ais *AppInfoScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ais.SetSize(msg.Width, msg.Height)
		// Инициализируем viewport при первом изменении размера
		if ais.viewport.Width == 0 && ais.viewport.Height == 0 {
			ais.viewport = viewport.New(msg.Width-4, msg.Height-12)
			ais.updateViewportContent()
		} else {
			ais.viewport.Width = msg.Width - 4
			ais.viewport.Height = msg.Height - 12
		}
		return ais, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q", "esc":
			// Возвращаемся к настройкам
			return ais, func() tea.Msg {
				return ui.NavigateToMsg{ScreenName: "settings"}
			}
		case "up", "k":
			ais.viewport.ScrollUp(1)
		case "down", "j":
			ais.viewport.ScrollDown(1)
		case "pageup":
			ais.viewport.PageUp()
		case "pagedown":
			ais.viewport.PageDown()
		}
	}

	// Обновляем содержимое viewport
	ais.updateViewportContent()

	// Обновляем базовый экран
	baseScreen, baseCmd := ais.BaseScreen.Update(msg)
	ais.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return ais, cmd
}

// View возвращает строку для отрисовки
func (ais *AppInfoScreen) View() string {
	// Подготавливаем содержимое viewport
	viewportContent := ais.viewport.View()

	// Добавляем индикатор прокрутки под viewport если не дошли до конца
	if !ais.viewport.AtBottom() {
		// Создаем более информативный индикатор прокрутки
		scrollIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorMuted)).
			Align(lipgloss.Center).
			Render("↓ Прокрутите вниз для просмотра всей информации ↓")

		// Добавляем индикатор под содержимое viewport
		viewportContent += "\n" + scrollIndicator
	}

	// Устанавливаем содержимое с индикатором
	ais.SetContent(viewportContent)
	return ais.BaseScreen.View()
}

// updateViewportContent обновляет содержимое viewport
func (ais *AppInfoScreen) updateViewportContent() {
	// Получаем версию приложения через пакет version
	appVersion := version.GetVersion()

	// Получаем информацию о системе
	goVersion := runtime.Version()
	goOS := runtime.GOOS
	goArch := runtime.GOARCH

	// Получаем время сборки через пакет version
	buildTime := version.GetBuildTime()

	// Создаем стили
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorPrimary)).
		Bold(true).
		Margin(0, 0, 2, 0).
		Align(lipgloss.Center)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorSecondary)).
		Bold(true).
		Margin(0, 0, 1, 0)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorText)).
		Margin(0, 0, 1, 0)

	contactStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted)).
		Margin(0, 0, 1, 0).
		Align(lipgloss.Center)

	// Создаем содержимое
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("[i] Информация о приложении"),

		sectionStyle.Render("Версия:"),
		infoStyle.Render(fmt.Sprintf("  SSH Keeper %s", appVersion)),

		sectionStyle.Render("Система:"),
		infoStyle.Render(fmt.Sprintf("  Операционная система: %s", goOS)),
		infoStyle.Render(fmt.Sprintf("  Архитектура: %s", goArch)),
		infoStyle.Render(fmt.Sprintf("  Go версия: %s", goVersion)),

		sectionStyle.Render("Сборка:"),
		infoStyle.Render(fmt.Sprintf("  Время сборки: %s", buildTime)),

		sectionStyle.Render("Описание:"),
		infoStyle.Render("  SSH Keeper - это безопасный менеджер SSH подключений"),
		infoStyle.Render("  с шифрованием данных и удобным интерфейсом."),

		sectionStyle.Render("Автор:"),
		infoStyle.Render("  Разработано @ankogit (https://github.com/ankogit) с <3 для удобного управления SSH подключениями"),

		contactStyle.Render(""),
		contactStyle.Render("Нажмите ESC для возврата к настройкам"),
	)

	ais.viewport.SetContent(content)
}

// Init инициализирует экран
func (ais *AppInfoScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (ais *AppInfoScreen) GetName() string {
	return "app_info"
}
