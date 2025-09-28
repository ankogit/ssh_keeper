package components

import (
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

// FieldStyles содержит стили для полей формы
type FieldStyles struct {
	FieldStyle        lipgloss.Style
	LabelStyle        lipgloss.Style
	InstructionsStyle lipgloss.Style
	PortFieldStyle    lipgloss.Style
}

// NewFieldStyles создает новый набор стилей для полей
func NewFieldStyles() *FieldStyles {
	return &FieldStyles{
		FieldStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorGray)).
			Padding(0, 1).
			Width(50),

		LabelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorSecondary)).
			Bold(styles.TextBold).
			Width(15),

		InstructionsStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorMuted)).
			Italic(styles.TextItalic),

		PortFieldStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorGray)).
			Padding(0, 1).
			Width(15),
	}
}

// GetFieldStyle возвращает стиль поля в зависимости от состояния
func (fs *FieldStyles) GetFieldStyle(fieldIndex int, currentField int, hasError bool, errorKey string, errors map[string]string) lipgloss.Style {
	baseStyle := fs.FieldStyle

	// Проверяем ошибки
	if hasError {
		if _, exists := errors[errorKey]; exists {
			return baseStyle.
				BorderForeground(lipgloss.Color(styles.ColorError))
		}
	}

	// Проверяем активное поле
	if currentField == fieldIndex {
		return baseStyle.
			BorderForeground(lipgloss.Color(styles.ColorWarning))
	}

	return baseStyle
}

// GetPortFieldStyle возвращает стиль поля порта в зависимости от состояния
func (fs *FieldStyles) GetPortFieldStyle(fieldIndex int, currentField int, hasError bool, errorKey string, errors map[string]string) lipgloss.Style {
	baseStyle := fs.PortFieldStyle

	// Проверяем ошибки
	if hasError {
		if _, exists := errors[errorKey]; exists {
			return baseStyle.
				BorderForeground(lipgloss.Color(styles.ColorError))
		}
	}

	// Проверяем активное поле
	if currentField == fieldIndex {
		return baseStyle.
			BorderForeground(lipgloss.Color(styles.ColorWarning))
	}

	return baseStyle
}
