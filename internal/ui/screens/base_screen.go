package screens

import (
	"ssh-keeper/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BaseScreen представляет базовый шаблон для всех экранов
type BaseScreen struct {
	title   string
	content string
	width   int
	height  int
}

// NewBaseScreen создает новый базовый экран
func NewBaseScreen(title string) *BaseScreen {
	return &BaseScreen{
		title:   title,
		content: "",
		width:   0,
		height:  0,
	}
}

// SetContent устанавливает содержимое экрана
func (bs *BaseScreen) SetContent(content string) {
	bs.content = content
}

// SetSize устанавливает размеры экрана
func (bs *BaseScreen) SetSize(width, height int) {
	bs.width = width
	bs.height = height
}

// Render отрисовывает экран
func (bs *BaseScreen) Render() string {
	if bs.width == 0 || bs.height == 0 {
		return "Loading..."
	}
	if bs.width < 80 {
		return "Resize the window to at least 80 characters wide..."
	}
	width := bs.width
	// height := bs.height

	// Создаем стили для компонентов
	containerStyle := lipgloss.NewStyle().
		Width(width - styles.ContainerBorderWidth).
		// Height(bs.height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(styles.ColorPrimary)).
		// Background(lipgloss.Color(styles.ColorBackground)).
		Foreground(lipgloss.Color(styles.ColorText))

	// Заголовок
	headerStyle := lipgloss.NewStyle().
		MaxWidth(width). // Учитываем рамку контейнера
		Bold(styles.TextBold).
		Foreground(lipgloss.Color(styles.ColorPrimary)).
		// Background(lipgloss.Color(styles.ColorContainer)).
		Padding(0, styles.HeaderPadding).
		Margin(0, 0, styles.HeaderMargin, 0).
		Align(lipgloss.Center)

	// Контент с рамкой
	contentStyle := lipgloss.NewStyle().
		Width(width-styles.ContentBorderWidth). // Учитываем рамку контейнера (2 символа с каждой стороны)
		// Height(bs.height-4). // Учитываем заголовок (1 строка) + рамку контейнера (2 строки) + отступ (1 строка)
		Padding(styles.ContentPadding, styles.ContentPadding).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(styles.ColorSecondary)).
		// Background(lipgloss.Color(styles.ColorContainer)).
		Foreground(lipgloss.Color(styles.ColorText))

	// Создаем заголовок
	header := headerStyle.Render(bs.title)

	// Создаем контент
	content := contentStyle.Render(bs.content)

	// Объединяем компоненты
	fullContent := lipgloss.JoinVertical(lipgloss.Center, header, content)

	return containerStyle.Render(fullContent)
}

// Update обрабатывает обновления состояния экрана
func (bs *BaseScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		bs.SetSize(msg.Width, msg.Height)
		return bs, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return bs, tea.Quit
		}
	}
	return bs, nil
}

// View возвращает строку для отрисовки
func (bs *BaseScreen) View() string {
	return bs.Render()
}

// Init инициализирует экран
func (bs *BaseScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (bs *BaseScreen) GetName() string {
	return "base"
}
