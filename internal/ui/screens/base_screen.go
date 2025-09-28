package screens

import (
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

	// Создаем стили для компонентов
	containerStyle := lipgloss.NewStyle().
		Width(bs.width - 10).
		// Height(bs.height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		// Background(lipgloss.Color("#1A1A1A")).
		Foreground(lipgloss.Color("#FFFFFF"))

	// Заголовок
	headerStyle := lipgloss.NewStyle().
		MaxWidth(bs.width). // Учитываем рамку контейнера
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		// Background(lipgloss.Color("#2D2D2D")).
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Align(lipgloss.Center)

	// Контент с рамкой
	contentStyle := lipgloss.NewStyle().
		Width(bs.width-12). // Учитываем рамку контейнера (2 символа с каждой стороны)
		// Height(bs.height-4). // Учитываем заголовок (1 строка) + рамку контейнера (2 строки) + отступ (1 строка)
		Padding(1, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#04B575")).
		// Background(lipgloss.Color("#2D2D2D")).
		Foreground(lipgloss.Color("#FFFFFF"))

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
		case "q", "ctrl+c":
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
