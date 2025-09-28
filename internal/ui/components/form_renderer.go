package components

import (
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// FormRenderer отвечает за рендеринг полей формы
type FormRenderer struct {
	styles *FieldStyles
}

// NewFormRenderer создает новый рендерер формы
func NewFormRenderer() *FormRenderer {
	return &FormRenderer{
		styles: NewFieldStyles(),
	}
}

// RenderTextField рендерит текстовое поле
func (fr *FormRenderer) RenderTextField(
	label string,
	input textinput.Model,
	fieldIndex int,
	currentField int,
	errors map[string]string,
	errorKey string,
) string {
	labelText := fr.styles.LabelStyle.Render(label + ":")
	
	fieldStyle := fr.styles.GetFieldStyle(fieldIndex, currentField, true, errorKey, errors)
	fieldText := fieldStyle.Render(input.View())
	
	return lipgloss.JoinHorizontal(lipgloss.Center, labelText, fieldText)
}

// RenderPortField рендерит поле порта
func (fr *FormRenderer) RenderPortField(
	label string,
	input textinput.Model,
	fieldIndex int,
	currentField int,
	errors map[string]string,
	errorKey string,
) string {
	labelText := fr.styles.LabelStyle.Render(label + ":")
	
	fieldStyle := fr.styles.GetPortFieldStyle(fieldIndex, currentField, true, errorKey, errors)
	fieldText := fieldStyle.Render(input.View())
	
	return lipgloss.JoinHorizontal(lipgloss.Center, labelText, fieldText)
}

// RenderBoolField рендерит булевое поле
func (fr *FormRenderer) RenderBoolField(
	label string,
	boolField *BoolField,
) string {
	labelText := fr.styles.LabelStyle.Render(label + ":")
	fieldText := boolField.View()
	
	return lipgloss.JoinHorizontal(lipgloss.Center, labelText, fieldText)
}

// RenderInstructions рендерит инструкции
func (fr *FormRenderer) RenderInstructions(text string) string {
	return fr.styles.InstructionsStyle.Render(text)
}

// RenderScrollIndicator рендерит индикатор прокрутки
func (fr *FormRenderer) RenderScrollIndicator() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted)).
		Render("↓↓↓")
}
