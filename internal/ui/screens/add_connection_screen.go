package screens

import (
	"fmt"
	"ssh-keeper/internal/models"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ssh"
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

	errors map[string]string

	// Менеджер формы
	formManager *components.FormManager

	// Менеджер сообщений
	messageManager *components.MessageManager

	// Прокрутка
	viewport viewport.Model
	ready    bool
}

// NewAddConnectionScreen создает новый экран добавления подключения
func NewAddConnectionScreen() *AddConnectionScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Добавить подключение")

	connectionSvc := services.GetGlobalConnectionService()

	// Создаем менеджер формы
	formManager := components.NewFormManager()

	// Создаем менеджер сообщений
	messageManager := components.NewMessageManager()

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
		Label:       "Использовать пароль (←/→)",
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

	// Добавляем кнопки
	formManager.AddField(components.FieldConfig{
		Name:      "save",
		Label:     "Сохранить",
		FieldType: components.FieldTypeButton,
		Required:  false,
		Width:     20,
	})

	formManager.AddField(components.FieldConfig{
		Name:      "test",
		Label:     "Тестировать подключение",
		FieldType: components.FieldTypeButton,
		Required:  false,
		Width:     30,
	})

	// Создаем viewport для прокрутки
	vp := viewport.New(0, 0)

	// Создаем экран
	screen := &AddConnectionScreen{
		BaseScreen:     baseScreen,
		connectionSvc:  connectionSvc,
		errors:         make(map[string]string),
		formManager:    formManager,
		messageManager: messageManager,
		viewport:       vp,
		ready:          false,
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
			// Учитываем все служебные элементы: заголовки, рамки, отступы
			acs.viewport = viewport.New(msg.Width-4, msg.Height-12)
			acs.ready = true
			acs.updateViewportContent()
		} else {
			// Используем одинаковый размер для консистентности
			acs.viewport.Width = msg.Width - 4
			acs.viewport.Height = msg.Height - 12
		}
		return acs, nil

	case ui.NavigateToMsg:
		// Очищаем форму только при навигации к другому экрану (не к самому add_connection)
		if msg.ScreenName != "add_connection" {
			acs.clearForm()
		}
		return acs, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return acs, tea.Quit
		case "esc":
			// Очищаем форму и возвращаемся
			acs.clearForm()
			return acs, ui.GoBackCmd()
		case "tab":
			// Переходим к следующему полю
			acs.formManager.NextField()
		case "shift+tab":
			// Переходим к предыдущему полю
			acs.formManager.PrevField()
		case "enter":
			// Проверяем, какая кнопка была нажата
			currentFieldName := acs.formManager.GetCurrentField()
			currentField := acs.formManager.GetField(currentFieldName)

			if currentField != nil && currentField.IsButton() {
				buttonName := currentField.GetName()
				switch buttonName {
				case "save":
					// Сохраняем подключение
					return acs, acs.saveConnection()
				case "test":
					// Тестируем подключение
					return acs, acs.testConnection()
				}
			} else if acs.formManager.IsLastField() {
				// Если это последнее поле (не кнопка), переходим к первой кнопке
				acs.formManager.NextField()
			} else {
				// Переходим к следующему полю
				acs.formManager.NextField()
			}
		case "ctrl+s":
			// Сохранить подключение
			return acs, acs.saveConnection()
		case "ctrl+t":
			// Тестировать SSH подключение
			return acs, acs.testConnection()
		case "ctrl+p":
			// Переключить показ пароля
			acs.togglePasswordVisibility()
		case "space", "left", "right":
			// Обрабатываем булевое поле через FormManager
			if acs.formManager.GetCurrentField() == components.FieldNameAuth {
				// Обновляем поле в FormManager напрямую
				field := acs.formManager.GetField(components.FieldNameAuth)
				if field != nil {
					_, fieldCmd := field.Update(msg)
					if fieldCmd != nil {
						if teaCmd, ok := fieldCmd.(tea.Cmd); ok {
							cmd = teaCmd
						}
					}
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
			// Добавляем команду для прокрутки после обновления содержимого
			cmd = tea.Sequence(cmd, acs.scrollToCurrentFieldCmd())
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

	// Добавляем сообщения в начало
	messages := acs.messageManager.RenderMessages(80) // Используем фиксированную ширину
	if messages != "" {
		viewportContent = messages + viewportContent
	}

	// Добавляем индикатор прокрутки под viewport если не дошли до конца
	if !acs.viewport.AtBottom() {
		// Создаем более информативный индикатор прокрутки
		scrollIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorMuted)).
			Align(lipgloss.Center).
			Render("↓ Прокрутите вниз для просмотра всех полей ↓")

		// Добавляем индикатор под содержимое viewport
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
func (acs *AddConnectionScreen) saveConnection() tea.Cmd {
	// Валидируем все поля
	if !acs.formManager.ValidateAll() {
		// Показываем сообщение об ошибке валидации
		acs.messageManager.AddError("❌ Пожалуйста, заполните все обязательные поля")
		return nil
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
		UseSSHKey:   !(values[components.FieldNameAuth] == "true"), // Если не пароль, то SSH ключ
		HasPassword: values[components.FieldNameAuth] == "true" && values[components.FieldNamePassword] != "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Добавляем пароль если используется
	if connection.HasPassword {
		connection.Password = values[components.FieldNamePassword]
	}

	// Сохраняем подключение
	err := acs.connectionSvc.AddConnection(connection)
	if err != nil {
		// Показываем ошибку сохранения
		acs.messageManager.AddError(fmt.Sprintf("Ошибка сохранения: %v", err))
		return nil
	}

	// Очищаем ошибки
	acs.errors = make(map[string]string)

	// Добавляем сообщение об успехе
	acs.messageManager.AddSuccess(fmt.Sprintf("Подключение '%s' успешно добавлено!", connection.Name))

	// Очищаем форму
	acs.clearForm()

	// Остаемся на экране добавления
	return nil
}

// testConnection тестирует SSH подключение без сохранения
func (acs *AddConnectionScreen) testConnection() tea.Cmd {
	// Валидируем все поля
	if !acs.formManager.ValidateAll() {
		// Показываем сообщение об ошибке валидации
		acs.messageManager.AddError("❌ Пожалуйста, заполните все обязательные поля для тестирования")
		return nil
	}

	// Получаем значения полей
	values := acs.formManager.GetValues()

	// Создаем подключение для тестирования
	port := 22 // По умолчанию
	if values[components.FieldNamePort] != "" {
		if p, err := strconv.Atoi(values[components.FieldNamePort]); err == nil {
			port = p
		}
	}

	connection := &models.Connection{
		Name:        values[components.FieldNameName],
		Host:        values[components.FieldNameHost],
		Port:        port,
		User:        values[components.FieldNameUser],
		KeyPath:     values[components.FieldNameKey],
		UseSSHKey:   !(values[components.FieldNameAuth] == "true"), // Если не пароль, то SSH ключ
		HasPassword: values[components.FieldNameAuth] == "true" && values[components.FieldNamePassword] != "",
	}

	// Добавляем пароль если используется
	if connection.HasPassword {
		connection.Password = values[components.FieldNamePassword]
	}

	// Тестируем подключение
	acs.messageManager.AddInfo("Тестирование SSH подключения...")

	// Создаем SSH клиент
	clientFactory := ssh.NewClientFactory()
	client := clientFactory.CreateClient(connection)

	// Пытаемся подключиться
	err := client.Connect()
	if err != nil {
		acs.messageManager.AddError(fmt.Sprintf("❌ SSH подключение не удалось: %v", err))
		return nil
	}

	// Показываем успех
	acs.messageManager.AddSuccess(fmt.Sprintf("SSH подключение к %s@%s:%d успешно!", connection.User, connection.Host, connection.Port))

	return nil
}

// clearForm очищает все поля формы
func (acs *AddConnectionScreen) clearForm() {
	// Очищаем все поля через FormManager
	for _, fieldName := range acs.formManager.GetFieldOrder() {
		field := acs.formManager.GetField(fieldName)
		if field != nil {
			field.SetValue("")
		}
	}

	// Очищаем ошибки
	acs.errors = make(map[string]string)

	// Устанавливаем фокус на первое поле
	acs.formManager.SetCurrentField(components.FieldNameName)
	acs.formManager.UpdateFocus()

	// Сбрасываем прокрутку в начало
	if acs.ready {
		acs.viewport.SetYOffset(0)
		acs.updateViewportContent()
	}
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

	// Более точная логика прокрутки
	// Каждое поле занимает примерно 4 строки (лейбл + поле + отступы)
	linesPerField := 4
	targetLine := currentIndex * linesPerField

	// Получаем высоту viewport и общую высоту содержимого
	viewportHeight := acs.viewport.Height
	totalContentHeight := acs.viewport.TotalLineCount()

	// Если содержимое помещается в viewport, не прокручиваем
	if totalContentHeight <= viewportHeight {
		acs.viewport.SetYOffset(0)
		return
	}

	// Рассчитываем оптимальную позицию прокрутки
	// Пытаемся поместить текущее поле в верхнюю треть viewport
	optimalOffset := targetLine - viewportHeight/3

	// Ограничиваем прокрутку границами содержимого
	maxOffset := totalContentHeight - viewportHeight
	if optimalOffset < 0 {
		optimalOffset = 0
	} else if optimalOffset > maxOffset {
		optimalOffset = maxOffset
	}

	// Плавно прокручиваем к оптимальной позиции
	acs.viewport.SetYOffset(optimalOffset)
}

// scrollToCurrentFieldCmd возвращает команду для прокрутки к текущему полю
func (acs *AddConnectionScreen) scrollToCurrentFieldCmd() tea.Cmd {
	return func() tea.Msg {
		acs.scrollToCurrentField()
		return nil
	}
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
