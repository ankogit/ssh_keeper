package components

import (
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FieldType определяет тип поля
type FieldType int

const (
	FieldTypeText FieldType = iota
	FieldTypePort
	FieldTypePassword
	FieldTypeBool
	FieldTypeButton
)

// FieldConfig содержит конфигурацию поля
type FieldConfig struct {
	Name        string
	Label       string
	Required    bool
	Width       int
	MaxLength   int
	Placeholder string
	FieldType   FieldType
	Style       string // Стиль для кнопок: "default", "warning", "error", "success"
}

// FormField представляет универсальное поле формы
type FormField struct {
	config      FieldConfig
	input       textinput.Model
	boolField   *BoolField
	buttonField *ButtonField
	value       string
	hasError    bool
	focused     bool
	visible     bool
}

// NewFormField создает новое поле формы
func NewFormField(config FieldConfig) *FormField {
	field := &FormField{
		config:   config,
		hasError: false,
		focused:  false,
		visible:  true,
	}

	switch config.FieldType {
	case FieldTypeBool:
		field.boolField = NewBoolField(config.Label)
		field.boolField.SetWidth(config.Width)
	case FieldTypeButton:
		field.buttonField = NewButtonField(config.Label)
		field.buttonField.SetWidth(config.Width)
		if config.Style != "" {
			field.buttonField.SetStyle(config.Style)
		}
	default:
		field.input = textinput.New()
		field.input.Placeholder = config.Placeholder
		field.input.Width = config.Width

		// Устанавливаем максимальную длину, если указана
		if config.MaxLength > 0 {
			field.input.CharLimit = config.MaxLength
		} else {
			field.input.CharLimit = 100 // Значение по умолчанию
		}

		if config.FieldType == FieldTypePassword {
			field.input.EchoMode = textinput.EchoPassword
		}
	}

	return field
}

// Update обновляет состояние поля
func (ff *FormField) Update(msg interface{}) (interface{}, interface{}) {
	switch ff.config.FieldType {
	case FieldTypeBool:
		return ff.boolField.Update(msg)
	default:
		var cmd tea.Cmd
		ff.input, cmd = ff.input.Update(msg)
		return ff, cmd
	}
}

// Focus устанавливает фокус на поле
func (ff *FormField) Focus() {
	ff.focused = true
	switch ff.config.FieldType {
	case FieldTypeBool:
		ff.boolField.Focus()
	case FieldTypeButton:
		ff.buttonField.Focus()
	default:
		ff.input.Focus()
	}
}

// Blur убирает фокус с поля
func (ff *FormField) Blur() {
	ff.focused = false
	switch ff.config.FieldType {
	case FieldTypeBool:
		ff.boolField.Blur()
	case FieldTypeButton:
		ff.buttonField.Blur()
	default:
		ff.input.Blur()
	}
}

// Value возвращает значение поля
func (ff *FormField) Value() string {
	switch ff.config.FieldType {
	case FieldTypeBool:
		if ff.boolField.Value() {
			return "true"
		}
		return "false"
	case FieldTypeButton:
		return ff.buttonField.Value()
	default:
		return ff.input.Value()
	}
}

// SetValue устанавливает значение поля
func (ff *FormField) SetValue(value string) {
	switch ff.config.FieldType {
	case FieldTypeBool:
		ff.boolField.SetValue(value == "true")
	default:
		ff.input.SetValue(value)
	}
}

// SetError устанавливает ошибку для поля
func (ff *FormField) SetError(hasError bool) {
	ff.hasError = hasError
}

// HasError возвращает true если поле имеет ошибку
func (ff *FormField) HasError() bool {
	return ff.hasError
}

// IsFocused возвращает true если поле в фокусе
func (ff *FormField) IsFocused() bool {
	return ff.focused
}

// SetVisible устанавливает видимость поля
func (ff *FormField) SetVisible(visible bool) {
	ff.visible = visible
}

// IsVisible возвращает true если поле видимо
func (ff *FormField) IsVisible() bool {
	return ff.visible
}

// Validate валидирует поле
func (ff *FormField) Validate() bool {
	value := ff.Value()

	// Если поле не обязательное и пустое - валидно
	if !ff.config.Required && value == "" {
		return true
	}

	// Если поле обязательное и пустое - невалидно
	if ff.config.Required && value == "" {
		ff.SetError(true)
		return false
	}

	// Специальная валидация для порта
	if ff.config.FieldType == FieldTypePort && value != "" {
		// Проверяем что это число от 1 до 65535
		port := 0
		for _, r := range value {
			if r < '0' || r > '9' {
				ff.SetError(true)
				return false
			}
			port = port*10 + int(r-'0')
		}
		if port < 1 || port > 65535 {
			ff.SetError(true)
			return false
		}
	}

	ff.SetError(false)
	return true
}

// Render возвращает отрендеренное поле
func (ff *FormField) Render() string {
	// Стили
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorSecondary)).
		Bold(styles.TextBold).
		Width(15)

	// Определяем стиль поля в зависимости от типа
	var fieldStyle lipgloss.Style
	switch ff.config.FieldType {
	case FieldTypePort:
		fieldStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorGray)).
			Padding(0, 1).
			Width(ff.config.Width)
	default:
		fieldStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorGray)).
			Padding(0, 1).
			Width(ff.config.Width)
	}

	// Применяем стили в зависимости от состояния
	if ff.hasError {
		fieldStyle = fieldStyle.BorderForeground(lipgloss.Color(styles.ColorError))
	} else if ff.focused {
		fieldStyle = fieldStyle.BorderForeground(lipgloss.Color(styles.ColorWarning))
	}

	// Рендерим поле
	var fieldContent string
	switch ff.config.FieldType {
	case FieldTypeBool:
		fieldContent = ff.boolField.View()
	case FieldTypeButton:
		fieldContent = ff.buttonField.View()
	default:
		fieldContent = fieldStyle.Render(ff.input.View())
	}

	// Рендерим лейбл
	var labelContent string
	if ff.config.FieldType == FieldTypeButton {
		// Для кнопки не показываем лейбл отдельно
		labelContent = ""
	} else {
		labelContent = labelStyle.Render(ff.config.Label + ":")
	}

	if labelContent == "" {
		return fieldContent
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, labelContent, fieldContent)
}

// GetTextInput возвращает textinput.Model для текстовых полей
func (ff *FormField) GetTextInput() (textinput.Model, bool) {
	if ff.config.FieldType == FieldTypeBool {
		return textinput.Model{}, false
	}
	return ff.input, true
}

// SetTextInput устанавливает textinput.Model для текстовых полей
func (ff *FormField) SetTextInput(input textinput.Model) {
	if ff.config.FieldType != FieldTypeBool {
		ff.input = input
	}
}

// IsButton проверяет, является ли поле кнопкой
func (ff *FormField) IsButton() bool {
	return ff.config.FieldType == FieldTypeButton
}

// GetName возвращает имя поля
func (ff *FormField) GetName() string {
	return ff.config.Name
}
