package screens

import (
	"fmt"
	"ssh-keeper/internal/models"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"ssh-keeper/internal/ui/styles"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// AddConnectionScreen представляет экран добавления нового подключения
type AddConnectionScreen struct {
	*BaseScreen
	connectionSvc *services.ConnectionService

	// Поля формы
	nameInput     textinput.Model
	hostInput     textinput.Model
	portInput     textinput.Model
	userInput     textinput.Model
	keyPathInput  textinput.Model
	passwordInput textinput.Model
	usePassword   *components.BoolField // Булевое поле для выбора аутентификации
	errors        map[string]string

	// Менеджер формы
	formManager *components.FormManager

	// Прокрутка
	viewport viewport.Model
	ready    bool
}

// NewAddConnectionScreen создает новый экран добавления подключения
func NewAddConnectionScreen() *AddConnectionScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Добавить подключение")

	connectionSvc := services.NewConnectionService()

	// Создаем поля формы
	nameInput := textinput.New()
	nameInput.Placeholder = "Название подключения"
	nameInput.CharLimit = 50
	nameInput.Width = 50

	hostInput := textinput.New()
	hostInput.Placeholder = "example.com"
	hostInput.CharLimit = 100
	hostInput.Width = 50

	portInput := textinput.New()
	portInput.Placeholder = "22"
	portInput.CharLimit = 5
	portInput.Width = 15

	userInput := textinput.New()
	userInput.Placeholder = "username"
	userInput.CharLimit = 50
	userInput.Width = 50

	keyPathInput := textinput.New()
	keyPathInput.Placeholder = "/path/to/private/key"
	keyPathInput.CharLimit = 200
	keyPathInput.Width = 50

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Пароль (если нет ключа)"
	passwordInput.CharLimit = 100
	passwordInput.Width = 50
	passwordInput.EchoMode = textinput.EchoPassword

	// Поля созданы выше

	// Создаем булевое поле
	usePasswordField := components.NewBoolField("Использовать пароль")
	usePasswordField.SetWidth(20)

	// Создаем менеджер формы
	formManager := components.NewFormManager()

	// Добавляем поля в форму
	formManager.AddField(components.FieldConfig{
		Name:        components.FieldNameName,
		Label:       "Название",
		Required:    true,
		Width:       50,
		MaxLength:   50,
		Placeholder: "Название подключения",
		FieldType:   components.FieldTypeText,
	})

	formManager.AddField(components.FieldConfig{
		Name:        components.FieldNameHost,
		Label:       "Хост",
		Required:    true,
		Width:       50,
		MaxLength:   100,
		Placeholder: "IP адрес или домен",
		FieldType:   components.FieldTypeText,
	})

	formManager.AddField(components.FieldConfig{
		Name:        components.FieldNamePort,
		Label:       "Порт",
		Required:    false,
		Width:       15,
		MaxLength:   5,
		Placeholder: "22",
		FieldType:   components.FieldTypePort,
	})

	formManager.AddField(components.FieldConfig{
		Name:        components.FieldNameUser,
		Label:       "Пользователь",
		Required:    true,
		Width:       50,
		MaxLength:   50,
		Placeholder: "Имя пользователя",
		FieldType:   components.FieldTypeText,
	})

	formManager.AddField(components.FieldConfig{
		Name:        components.FieldNameAuth,
		Label:       "Использовать пароль",
		Required:    true,
		Width:       20,
		Placeholder: "",
		FieldType:   components.FieldTypeBool,
	})

	formManager.AddField(components.FieldConfig{
		Name:        components.FieldNamePassword,
		Label:       "Пароль",
		Required:    false,
		Width:       50,
		MaxLength:   100,
		Placeholder: "Пароль (если нет ключа)",
		FieldType:   components.FieldTypePassword,
	})

	formManager.AddField(components.FieldConfig{
		Name:        components.FieldNameKey,
		Label:       "SSH ключ",
		Required:    false,
		Width:       50,
		MaxLength:   200,
		Placeholder: "Путь к SSH ключу",
		FieldType:   components.FieldTypeText,
	})

	// Создаем viewport для прокрутки
	vp := viewport.New(0, 0)

	// Создаем экран
	screen := &AddConnectionScreen{
		BaseScreen:    baseScreen,
		connectionSvc: connectionSvc,
		nameInput:     nameInput,
		hostInput:     hostInput,
		portInput:     portInput,
		userInput:     userInput,
		keyPathInput:  keyPathInput,
		passwordInput: passwordInput,
		usePassword:   usePasswordField,
		errors:        make(map[string]string),
		formManager:   formManager,
		viewport:      vp,
		ready:         false,
	}

	// Устанавливаем фокус на первое поле сразу
	screen.formManager.SetCurrentField(components.FieldNameName)
	screen.formManager.UpdateFocus()

	return screen
}

// Update обрабатывает обновления состояния
func (acs *AddConnectionScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		acs.SetSize(msg.Width, msg.Height)
		if !acs.ready {
			// Инициализируем viewport при первом изменении размера
			acs.viewport = viewport.New(msg.Width-4, msg.Height-12)
			acs.ready = true
			acs.updateViewportContent()
		} else {
			acs.viewport.Width = msg.Width - 4
			acs.viewport.Height = msg.Height - 8
		}
		return acs, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return acs, tea.Quit
		case "esc":
			return acs, ui.GoBackCmd()
		case "tab":
			// Переходим к следующему полю
			acs.formManager.NextField()
		case "shift+tab":
			// Переходим к предыдущему полю
			acs.formManager.PrevField()
		case "enter":
			if acs.formManager.IsLastField() {
				// Последнее поле - сохраняем
				acs.saveConnection()
			} else {
				// Переходим к следующему полю
				acs.formManager.NextField()
			}
		case "ctrl+s":
			// Сохранить подключение
			acs.saveConnection()
		case "ctrl+p":
			// Переключить показ пароля
			acs.togglePasswordVisibility()
		case "space", "left", "right":
			// Обрабатываем булевое поле
			if acs.formManager.GetCurrentField() == components.FieldNameAuth {
				acs.usePassword, _ = acs.usePassword.Update(msg)
				// Синхронизируем с полем в FormManager
				authField := acs.formManager.GetField(components.FieldNameAuth)
				if acs.usePassword.Value() {
					authField.SetValue("true")
				} else {
					authField.SetValue("false")
				}
			}
		case "up":
			// Прокрутка вверх
			acs.viewport.ScrollUp(1)
		case "down":
			// Прокрутка вниз
			acs.viewport.ScrollDown(1)
		case "pageup":
			// Прокрутка на страницу вверх
			acs.viewport.PageUp()
		case "pagedown":
			// Прокрутка на страницу вниз
			acs.viewport.PageDown()
		}
	}

	// Обновляем все поля через менеджер формы
	for _, fieldName := range acs.formManager.GetFieldOrder() {
		field := acs.formManager.GetField(fieldName)
		_, fieldCmd := field.Update(msg)
		if fieldCmd != nil {
			if teaCmd, ok := fieldCmd.(tea.Cmd); ok {
				cmd = teaCmd
			}
		}
	}

	// Обновляем фокус через менеджер формы
	acs.formManager.UpdateFocus()

	// Обновляем видимость полей на основе выбора аутентификации
	acs.updateFieldVisibility()

	// Обновляем содержимое viewport
	acs.updateViewportContent()

	// Прокручиваем viewport для видимости текущего поля после навигации
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "tab", "shift+tab", "enter":
			acs.scrollToCurrentField()
		}
	}

	// Обновляем базовый экран
	baseScreen, baseCmd := acs.BaseScreen.Update(msg)
	acs.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return acs, cmd
}

// View возвращает строку для отрисовки
func (acs *AddConnectionScreen) View() string {
	if !acs.ready {
		return "Инициализация..."
	}

	// Подготавливаем содержимое viewport
	viewportContent := acs.viewport.View()

	// Добавляем индикатор прокрутки под viewport если не дошли до конца
	if !acs.viewport.AtBottom() {
		// Создаем стрелку вниз
		scrollIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorMuted)).
			Render("↓↓↓")

		// Добавляем стрелку под содержимое viewport
		viewportContent += "\n" + scrollIndicator
	}

	// Устанавливаем содержимое с индикатором
	acs.SetContent(viewportContent)
	return acs.BaseScreen.View()
}

// updateFieldVisibility обновляет видимость полей на основе выбора аутентификации
func (acs *AddConnectionScreen) updateFieldVisibility() {
	// Получаем значение из FormManager
	authField := acs.formManager.GetField(components.FieldNameAuth)
	usePassword := authField.Value() == "true"

	// Обновляем видимость в менеджере полей
	acs.formManager.GetField(components.FieldNamePassword).SetVisible(usePassword)
	acs.formManager.GetField(components.FieldNameKey).SetVisible(!usePassword)
}

// updateViewportContent обновляет содержимое viewport
func (acs *AddConnectionScreen) updateViewportContent() {
	content := acs.renderForm()
	acs.viewport.SetContent(content)
}

// renderForm рендерит форму в виде вертикального списка
func (acs *AddConnectionScreen) renderForm() string {
	// Используем менеджер формы для рендеринга
	return acs.formManager.RenderForm()
}

// Методы nextField и prevField удалены - навигация теперь обрабатывается через FieldManager

// togglePasswordVisibility переключает видимость пароля
func (acs *AddConnectionScreen) togglePasswordVisibility() {
	passwordField := acs.formManager.GetField(components.FieldNamePassword)
	if passwordField == nil {
		return
	}

	// Получаем textinput.Model из FormField
	if textInput, ok := passwordField.GetTextInput(); ok {
		if textInput.EchoMode == textinput.EchoPassword {
			textInput.EchoMode = textinput.EchoNormal
		} else {
			textInput.EchoMode = textinput.EchoPassword
		}
		// Обновляем поле в FormField
		passwordField.SetTextInput(textInput)
	}
}

// saveConnection сохраняет подключение
func (acs *AddConnectionScreen) saveConnection() {
	// Валидируем все поля
	if !acs.formManager.ValidateAll() {
		return
	}

	// Получаем значения полей
	values := acs.formManager.GetValues()

	// Создаем подключение
	port := 22 // По умолчанию
	if values[components.FieldNamePort] != "" {
		if p, err := strconv.Atoi(values[components.FieldNamePort]); err == nil {
			port = p
		}
	}

	connection := &models.Connection{
		ID:          uuid.New().String(),
		Name:        values[components.FieldNameName],
		Host:        values[components.FieldNameHost],
		Port:        port,
		User:        values[components.FieldNameUser],
		KeyPath:     values[components.FieldNameKey],
		HasPassword: values[components.FieldNameAuth] == "true",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Сохраняем подключение
	acs.connectionSvc.AddConnection(connection)

	// Выводим данные в консоль
	fmt.Printf("Подключение сохранено:\n")
	fmt.Printf("  Название: %s\n", connection.Name)
	fmt.Printf("  Хост: %s\n", connection.Host)
	fmt.Printf("  Порт: %d\n", connection.Port)
	fmt.Printf("  Пользователь: %s\n", connection.User)
	fmt.Printf("  Использовать пароль: %t\n", connection.HasPassword)
	if connection.HasPassword {
		fmt.Printf("  Пароль: %s\n", values[components.FieldNamePassword])
	} else {
		fmt.Printf("  SSH ключ: %s\n", values[components.FieldNameKey])
	}

	// Возвращаемся к предыдущему экрану
	// В реальной реализации здесь будет команда навигации
}

// scrollToCurrentField прокручивает viewport для видимости текущего поля
func (acs *AddConnectionScreen) scrollToCurrentField() {
	if !acs.ready {
		return
	}

	// Получаем текущее поле
	currentField := acs.formManager.GetCurrentField()
	fieldOrder := acs.formManager.GetFieldOrder()

	// Находим индекс текущего поля
	currentIndex := -1
	for i, fieldName := range fieldOrder {
		if fieldName == currentField {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return
	}

	// Упрощенная логика: просто прокручиваем на несколько строк вниз
	// Каждое поле занимает примерно 2 строки (лейбл + поле)
	linesPerField := 2
	targetLine := currentIndex * linesPerField

	// Простая прокрутка: устанавливаем позицию на текущее поле
	acs.viewport.SetYOffset(targetLine)
}

// Init инициализирует экран
func (acs *AddConnectionScreen) Init() tea.Cmd {
	// Устанавливаем фокус на первое поле через FormManager
	acs.formManager.SetCurrentField(components.FieldNameName)
	acs.formManager.UpdateFocus()
	return textinput.Blink
}

// GetName возвращает имя экрана
func (acs *AddConnectionScreen) GetName() string {
	return "add_connection"
}
