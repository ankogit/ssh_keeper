package components

import (
	"ssh-keeper/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BoolField представляет компонент для выбора булевого значения
type BoolField struct {
	value   bool
	focused bool
	label   string
	width   int
}

// NewBoolField создает новый булевый компонент
func NewBoolField(label string) *BoolField {
	return &BoolField{
		value:   false,
		focused: false,
		label:   label,
		width:   20,
	}
}

// SetValue устанавливает значение
func (bf *BoolField) SetValue(value bool) {
	bf.value = value
}

// Value возвращает текущее значение
func (bf *BoolField) Value() bool {
	return bf.value
}

// Focus устанавливает фокус
func (bf *BoolField) Focus() {
	bf.focused = true
}

// Blur убирает фокус
func (bf *BoolField) Blur() {
	bf.focused = false
}

// Focused возвращает состояние фокуса
func (bf *BoolField) Focused() bool {
	return bf.focused
}

// Toggle переключает значение
func (bf *BoolField) Toggle() {
	bf.value = !bf.value
}

// Update обрабатывает обновления
func (bf *BoolField) Update(msg tea.Msg) (*BoolField, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if bf.focused {
			switch msg.String() {
			case "left", "h":
				bf.value = false
			case "right", "l":
				bf.value = true
			case "space":
				bf.Toggle()
			}
		}
	}

	return bf, cmd
}

// View возвращает строку для отрисовки
func (bf *BoolField) View() string {
	// Символы для отображения
	const (
		filledCircle = "●" // Закрашенный кружок
		emptyCircle  = "○" // Пустой кружок
		trueText     = "Да"
		falseText    = "Нет"
	)

	// Определяем символ и текст
	var symbol, text string
	if bf.value {
		symbol = filledCircle
		text = trueText
	} else {
		symbol = emptyCircle
		text = falseText
	}

	// Стили
	baseStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(bf.width)

	// Применяем стили в зависимости от состояния
	var style lipgloss.Style
	if bf.focused {
		// Активное поле - оранжевая рамка
		style = baseStyle.
			BorderForeground(lipgloss.Color(styles.ColorWarning)) // Оранжевый
	} else {
		// Неактивное поле - обычная рамка
		style = baseStyle.
			BorderForeground(lipgloss.Color(styles.ColorGray)) // Серый
	}

	// Стили для символа
	symbolStyle := lipgloss.NewStyle()
	if bf.value {
		// Зеленый цвет для закрашенного кружка
		symbolStyle = symbolStyle.Foreground(lipgloss.Color(styles.ColorSuccess)) // Зеленый
	} else {
		// Серый цвет для пустого кружка
		symbolStyle = symbolStyle.Foreground(lipgloss.Color(styles.ColorMuted)) // Серый
	}

	// Объединяем символ и текст
	content := symbolStyle.Render(symbol) + " " + text

	return style.Render(content)
}

// SetWidth устанавливает ширину компонента
func (bf *BoolField) SetWidth(width int) {
	bf.width = width
}
