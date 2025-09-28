package components

import (
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

// FormManager управляет полями формы
type FormManager struct {
	fields       map[string]*FormField
	fieldOrder   []string
	currentField string
}

// NewFormManager создает новый менеджер формы
func NewFormManager() *FormManager {
	return &FormManager{
		fields:       make(map[string]*FormField),
		fieldOrder:   make([]string, 0),
		currentField: "",
	}
}

// AddField добавляет поле в форму
func (fm *FormManager) AddField(config FieldConfig) {
	field := NewFormField(config)
	fm.fields[config.Name] = field
	fm.fieldOrder = append(fm.fieldOrder, config.Name)

	// Устанавливаем первое поле как текущее
	if fm.currentField == "" {
		fm.currentField = config.Name
	}
}

// GetField возвращает поле по имени
func (fm *FormManager) GetField(name string) *FormField {
	return fm.fields[name]
}

// GetCurrentField возвращает текущее поле
func (fm *FormManager) GetCurrentField() string {
	return fm.currentField
}

// SetCurrentField устанавливает текущее поле
func (fm *FormManager) SetCurrentField(name string) {
	if _, exists := fm.fields[name]; exists {
		fm.currentField = name
	}
}

// NextField переходит к следующему видимому полю
func (fm *FormManager) NextField() {
	currentIndex := fm.getCurrentFieldIndex()
	for i := 0; i < len(fm.fieldOrder); i++ {
		nextIndex := (currentIndex + 1 + i) % len(fm.fieldOrder)
		fieldName := fm.fieldOrder[nextIndex]
		if fm.fields[fieldName].IsVisible() {
			fm.currentField = fieldName
			return
		}
	}
}

// PrevField переходит к предыдущему видимому полю
func (fm *FormManager) PrevField() {
	currentIndex := fm.getCurrentFieldIndex()
	for i := 0; i < len(fm.fieldOrder); i++ {
		prevIndex := (currentIndex - 1 - i + len(fm.fieldOrder)) % len(fm.fieldOrder)
		fieldName := fm.fieldOrder[prevIndex]
		if fm.fields[fieldName].IsVisible() {
			fm.currentField = fieldName
			return
		}
	}
}

// IsLastField проверяет, является ли текущее поле последним видимым
func (fm *FormManager) IsLastField() bool {
	// Находим последнее видимое поле
	for i := len(fm.fieldOrder) - 1; i >= 0; i-- {
		fieldName := fm.fieldOrder[i]
		if fm.fields[fieldName].IsVisible() {
			return fm.currentField == fieldName
		}
	}
	return false
}

// UpdateFocus обновляет фокус полей
func (fm *FormManager) UpdateFocus() {
	// Убираем фокус со всех полей
	for _, field := range fm.fields {
		field.Blur()
	}

	// Устанавливаем фокус на текущее поле
	if currentField, exists := fm.fields[fm.currentField]; exists {
		currentField.Focus()
	}
}

// GetCurrentFieldModel возвращает текущее поле для прямого доступа
func (fm *FormManager) GetCurrentFieldModel() *FormField {
	if currentField, exists := fm.fields[fm.currentField]; exists {
		return currentField
	}
	return nil
}

// ValidateAll валидирует все поля
func (fm *FormManager) ValidateAll() bool {
	allValid := true
	for _, field := range fm.fields {
		if !field.Validate() {
			allValid = false
		}
	}
	return allValid
}

// GetValues возвращает все значения полей
func (fm *FormManager) GetValues() map[string]string {
	values := make(map[string]string)
	for name, field := range fm.fields {
		values[name] = field.Value()
	}
	return values
}

// RenderForm рендерит всю форму
func (fm *FormManager) RenderForm() string {
	var formContent []string

	// Рендерим только видимые поля
	for _, fieldName := range fm.fieldOrder {
		field := fm.fields[fieldName]
		if field.IsVisible() {
			formContent = append(formContent, field.Render())
			formContent = append(formContent, "") // Пустая строка между полями
		}
	}

	// Инструкции
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted)).
		Italic(styles.TextItalic).
		Render("Tab - след. поле • ↑/↓ - прокрутка • Enter - сохранить • Esc - назад")
	formContent = append(formContent, instructions)

	return lipgloss.JoinVertical(lipgloss.Left, formContent...)
}

// GetFieldOrder возвращает порядок полей
func (fm *FormManager) GetFieldOrder() []string {
	return fm.fieldOrder
}

// getCurrentFieldIndex возвращает индекс текущего поля
func (fm *FormManager) getCurrentFieldIndex() int {
	for i, fieldName := range fm.fieldOrder {
		if fieldName == fm.currentField {
			return i
		}
	}
	return 0
}
