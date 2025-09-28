package components

import (
	"github.com/charmbracelet/bubbles/textinput"
)

// FieldNavigator управляет навигацией между полями формы
type FieldNavigator struct {
	currentField string
	fields       map[string]FieldInfo
	fieldOrder   []string
}

// FieldInfo содержит информацию о поле
type FieldInfo struct {
	Name     string
	Required bool
	Visible  bool
}

// NewFieldNavigator создает новый навигатор полей
func NewFieldNavigator() *FieldNavigator {
	return &FieldNavigator{
		currentField: "",
		fields:       make(map[string]FieldInfo),
		fieldOrder:   make([]string, 0),
	}
}

// AddField добавляет поле в навигатор
func (fn *FieldNavigator) AddField(name string, required bool) {
	fn.fields[name] = FieldInfo{
		Name:     name,
		Required: required,
		Visible:  true,
	}
	fn.fieldOrder = append(fn.fieldOrder, name)

	// Устанавливаем первое поле как текущее
	if fn.currentField == "" {
		fn.currentField = name
	}
}

// SetFieldVisibility устанавливает видимость поля
func (fn *FieldNavigator) SetFieldVisibility(fieldName string, visible bool) {
	if field, exists := fn.fields[fieldName]; exists {
		field.Visible = visible
		fn.fields[fieldName] = field
	}
}

// GetVisibleFields возвращает количество видимых полей
func (fn *FieldNavigator) GetVisibleFields() int {
	count := 0
	for _, field := range fn.fields {
		if field.Visible {
			count++
		}
	}
	return count
}

// GetTotalFields возвращает общее количество полей
func (fn *FieldNavigator) GetTotalFields() int {
	return len(fn.fields)
}

// GetCurrentField возвращает текущее поле
func (fn *FieldNavigator) GetCurrentField() string {
	return fn.currentField
}

// SetCurrentField устанавливает текущее поле
func (fn *FieldNavigator) SetCurrentField(fieldName string) {
	if _, exists := fn.fields[fieldName]; exists {
		fn.currentField = fieldName
	}
}

// NextField переходит к следующему видимому полю
func (fn *FieldNavigator) NextField() {
	currentIndex := fn.getCurrentFieldIndex()
	for i := 0; i < len(fn.fieldOrder); i++ {
		nextIndex := (currentIndex + 1 + i) % len(fn.fieldOrder)
		fieldName := fn.fieldOrder[nextIndex]
		if fn.fields[fieldName].Visible {
			fn.currentField = fieldName
			return
		}
	}
}

// PrevField переходит к предыдущему видимому полю
func (fn *FieldNavigator) PrevField() {
	currentIndex := fn.getCurrentFieldIndex()
	for i := 0; i < len(fn.fieldOrder); i++ {
		prevIndex := (currentIndex - 1 - i + len(fn.fieldOrder)) % len(fn.fieldOrder)
		fieldName := fn.fieldOrder[prevIndex]
		if fn.fields[fieldName].Visible {
			fn.currentField = fieldName
			return
		}
	}
}

// IsLastField проверяет, является ли текущее поле последним видимым
func (fn *FieldNavigator) IsLastField() bool {
	// Находим последнее видимое поле
	for i := len(fn.fieldOrder) - 1; i >= 0; i-- {
		fieldName := fn.fieldOrder[i]
		if fn.fields[fieldName].Visible {
			return fn.currentField == fieldName
		}
	}
	return false
}

// GetCurrentFieldInfo возвращает информацию о текущем поле
func (fn *FieldNavigator) GetCurrentFieldInfo() *FieldInfo {
	if field, exists := fn.fields[fn.currentField]; exists {
		return &field
	}
	return nil
}

// getCurrentFieldIndex возвращает индекс текущего поля в порядке
func (fn *FieldNavigator) getCurrentFieldIndex() int {
	for i, fieldName := range fn.fieldOrder {
		if fieldName == fn.currentField {
			return i
		}
	}
	return 0
}

// FieldManager управляет полями формы
type FieldManager struct {
	navigator     *FieldNavigator
	nameInput     textinput.Model
	hostInput     textinput.Model
	portInput     textinput.Model
	userInput     textinput.Model
	keyPathInput  textinput.Model
	passwordInput textinput.Model
	usePassword   *BoolField
}

// NewFieldManager создает новый менеджер полей
func NewFieldManager(
	nameInput textinput.Model,
	hostInput textinput.Model,
	portInput textinput.Model,
	userInput textinput.Model,
	keyPathInput textinput.Model,
	passwordInput textinput.Model,
	usePassword *BoolField,
) *FieldManager {
	navigator := NewFieldNavigator()

	// Добавляем поля в навигатор
	navigator.AddField(FieldNameName, true)      // FieldIndexName
	navigator.AddField(FieldNameHost, true)      // FieldIndexHost
	navigator.AddField(FieldNamePort, false)     // FieldIndexPort
	navigator.AddField(FieldNameUser, true)      // FieldIndexUser
	navigator.AddField(FieldNameAuth, true)      // FieldIndexAuth
	navigator.AddField(FieldNamePassword, false) // FieldIndexPassword (условное)
	navigator.AddField(FieldNameKey, false)      // FieldIndexKey (условное)

	return &FieldManager{
		navigator:     navigator,
		nameInput:     nameInput,
		hostInput:     hostInput,
		portInput:     portInput,
		userInput:     userInput,
		keyPathInput:  keyPathInput,
		passwordInput: passwordInput,
		usePassword:   usePassword,
	}
}

// GetCurrentField возвращает текущее поле
func (fm *FieldManager) GetCurrentField() string {
	return fm.navigator.GetCurrentField()
}

// NextField переходит к следующему полю
func (fm *FieldManager) NextField() {
	fm.navigator.NextField()
	fm.updateFocus()
}

// PrevField переходит к предыдущему полю
func (fm *FieldManager) PrevField() {
	fm.navigator.PrevField()
	fm.updateFocus()
}

// IsLastField проверяет, является ли текущее поле последним
func (fm *FieldManager) IsLastField() bool {
	return fm.navigator.IsLastField()
}

// updateFocus обновляет фокус полей
func (fm *FieldManager) updateFocus() {
	// Убираем фокус со всех полей
	fm.nameInput.Blur()
	fm.hostInput.Blur()
	fm.portInput.Blur()
	fm.userInput.Blur()
	fm.usePassword.Blur()
	fm.passwordInput.Blur()
	fm.keyPathInput.Blur()

	// Устанавливаем фокус на текущее поле
	switch fm.navigator.GetCurrentField() {
	case FieldNameName:
		fm.nameInput.Focus()
	case FieldNameHost:
		fm.hostInput.Focus()
	case FieldNamePort:
		fm.portInput.Focus()
	case FieldNameUser:
		fm.userInput.Focus()
	case FieldNameAuth:
		fm.usePassword.Focus()
	case FieldNamePassword:
		fm.passwordInput.Focus()
	case FieldNameKey:
		fm.keyPathInput.Focus()
	}
}

// UpdateFieldVisibility обновляет видимость условных полей
func (fm *FieldManager) UpdateFieldVisibility(usePassword bool) {
	// Поле пароля видимо только если usePassword = true
	fm.navigator.SetFieldVisibility(FieldNamePassword, usePassword)
	// Поле ключа видимо только если usePassword = false
	fm.navigator.SetFieldVisibility(FieldNameKey, !usePassword)
}

// GetActiveField возвращает активное поле для условной логики
func (fm *FieldManager) GetActiveField() (textinput.Model, *BoolField) {
	switch fm.navigator.GetCurrentField() {
	case FieldNameName:
		return fm.nameInput, nil
	case FieldNameHost:
		return fm.hostInput, nil
	case FieldNamePort:
		return fm.portInput, nil
	case FieldNameUser:
		return fm.userInput, nil
	case FieldNameAuth:
		return textinput.Model{}, fm.usePassword
	case FieldNamePassword:
		return fm.passwordInput, nil
	case FieldNameKey:
		return fm.keyPathInput, nil
	default:
		return textinput.Model{}, nil
	}
}
