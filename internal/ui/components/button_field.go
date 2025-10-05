package components

import (
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ButtonField представляет кнопку в форме
type ButtonField struct {
	label   string
	focused bool
	width   int
	style   string // Тип стиля: "default", "warning", "error", "success"
}

// NewButtonField создает новое поле кнопки
func NewButtonField(label string) *ButtonField {
	return &ButtonField{
		label:   label,
		focused: false,
		width:   20,
	}
}

// SetWidth устанавливает ширину кнопки
func (bf *ButtonField) SetWidth(width int) {
	bf.width = width
}

// SetStyle устанавливает стиль кнопки
func (bf *ButtonField) SetStyle(style string) {
	bf.style = style
}

// Focus устанавливает фокус на кнопку
func (bf *ButtonField) Focus() {
	bf.focused = true
}

// Blur убирает фокус с кнопки
func (bf *ButtonField) Blur() {
	bf.focused = false
}

// IsFocused проверяет, имеет ли кнопка фокус
func (bf *ButtonField) IsFocused() bool {
	return bf.focused
}

// SetLabel устанавливает текст кнопки
func (bf *ButtonField) SetLabel(label string) {
	bf.label = label
}

// GetLabel возвращает текст кнопки
func (bf *ButtonField) GetLabel() string {
	return bf.label
}

// View отображает кнопку
func (bf *ButtonField) View() string {
	// Определяем цвет в зависимости от стиля
	var borderColor string
	var textColor string = styles.ColorMuted

	switch bf.style {
	case "warning":
		borderColor = styles.ColorWarning
	case "error":
		borderColor = styles.ColorError
	case "success":
		borderColor = styles.ColorSuccess
	default:
		borderColor = styles.ColorSecondary
	}

	// Создаем стили для кнопки
	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(textColor)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Width(bf.width).
		Align(lipgloss.Center)

	// Если кнопка в фокусе, делаем текст ярче
	if bf.focused {
		buttonStyle = buttonStyle.Foreground(lipgloss.Color(styles.ColorText))
	}

	return buttonStyle.Render(bf.label)
}

// Update обрабатывает обновления состояния кнопки
func (bf *ButtonField) Update(msg interface{}) (*ButtonField, interface{}) {
	// Кнопка не обрабатывает ввод, только отображается
	return bf, nil
}

// Value возвращает значение кнопки (всегда пустая строка)
func (bf *ButtonField) Value() string {
	return ""
}

// SetValue устанавливает значение кнопки (игнорируется)
func (bf *ButtonField) SetValue(value string) {
	// Кнопка не имеет значения
}

// Validate валидирует кнопку (всегда true)
func (bf *ButtonField) Validate() bool {
	return true
}

// IsVisible проверяет, видима ли кнопка
func (bf *ButtonField) IsVisible() bool {
	return true
}

// SetVisible устанавливает видимость кнопки
func (bf *ButtonField) SetVisible(visible bool) {
	// Кнопка всегда видима
}

// GetWidth возвращает ширину кнопки
func (bf *ButtonField) GetWidth() int {
	return bf.width
}

// SetError устанавливает ошибку для кнопки (игнорируется)
func (bf *ButtonField) SetError(hasError bool) {
	// Кнопка не может иметь ошибку
}

// HasError проверяет, есть ли ошибка у кнопки
func (bf *ButtonField) HasError() bool {
	return false
}
