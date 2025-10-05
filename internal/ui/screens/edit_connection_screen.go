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
)

// EditConnectionScreen представляет экран редактирования подключения
type EditConnectionScreen struct {
	*BaseScreen
	connectionSvc *services.ConnectionService

	// Подключение для редактирования
	connection *models.Connection

	errors map[string]string

	// Менеджер формы
	formManager *components.FormManager

	// Менеджер сообщений
	messageManager *components.MessageManager

	// Прокрутка
	viewport viewport.Model
	ready    bool
}

// NewEditConnectionScreenEmpty создает пустой экран редактирования (для фабрики)
func NewEditConnectionScreenEmpty() *EditConnectionScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Редактировать подключение")

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
		Label:     "Сохранить изменения",
		FieldType: components.FieldTypeButton,
		Required:  false,
		Width:     25,
	})

	formManager.AddField(components.FieldConfig{
		Name:      "delete",
		Label:     "Удалить подключение",
		FieldType: components.FieldTypeButton,
		Required:  false,
		Width:     25,
		Style:     "warning",
	})

	// Создаем viewport для прокрутки
	vp := viewport.New(0, 0)

	// Создаем экран
	screen := &EditConnectionScreen{
		BaseScreen:     baseScreen,
		connectionSvc:  connectionSvc,
		connection:     nil, // Будет установлено через SetData
		errors:         make(map[string]string),
		formManager:    formManager,
		messageManager: messageManager,
		viewport:       vp,
		ready:          false,
	}

	// Устанавливаем фокус на первое поле
	screen.formManager.SetCurrentField(components.FieldNameName)
	screen.formManager.UpdateFocus()

	return screen
}

// SetData устанавливает данные для экрана редактирования
func (ecs *EditConnectionScreen) SetData(data interface{}) {
	if connection, ok := data.(models.Connection); ok {
		ecs.connection = &connection
		// Добавляем отладочную информацию
		ecs.messageManager.AddInfo(fmt.Sprintf("Получены данные для редактирования: %s (ID: %s)", connection.Name, connection.ID))
		// Обновляем заголовок экрана
		ecs.BaseScreen.SetTitle(fmt.Sprintf("SSH Keeper - Редактировать '%s'", connection.Name))

		// Предзаполняем форму данными подключения
		ecs.prefillForm()

		// Добавляем сообщение об успешной загрузке данных
		ecs.messageManager.AddSuccess(fmt.Sprintf("Данные подключения '%s' загружены", connection.Name))
	} else {
		// Добавляем сообщение об ошибке
		ecs.messageManager.AddError("Ошибка: не удалось загрузить данные подключения")
	}
}

// prefillForm предзаполняет форму данными подключения
func (ecs *EditConnectionScreen) prefillForm() {
	if ecs.connection == nil {
		ecs.messageManager.AddError("Ошибка: подключение не инициализировано")
		return
	}

	// Заполняем основные поля напрямую через textinput
	nameField := ecs.formManager.GetField(components.FieldNameName)
	if nameField != nil {
		if textInput, ok := nameField.GetTextInput(); ok {
			textInput.SetValue(ecs.connection.Name)
			nameField.SetTextInput(textInput)
		}
	}

	hostField := ecs.formManager.GetField(components.FieldNameHost)
	if hostField != nil {
		if textInput, ok := hostField.GetTextInput(); ok {
			textInput.SetValue(ecs.connection.Host)
			hostField.SetTextInput(textInput)
		}
	}

	portField := ecs.formManager.GetField(components.FieldNamePort)
	if portField != nil {
		if textInput, ok := portField.GetTextInput(); ok {
			textInput.SetValue(strconv.Itoa(ecs.connection.Port))
			portField.SetTextInput(textInput)
		}
	}

	userField := ecs.formManager.GetField(components.FieldNameUser)
	if userField != nil {
		if textInput, ok := userField.GetTextInput(); ok {
			textInput.SetValue(ecs.connection.User)
			userField.SetTextInput(textInput)
		}
	}

	// Определяем тип аутентификации
	authField := ecs.formManager.GetField(components.FieldNameAuth)
	if authField != nil {
		if ecs.connection.UseSSHKey {
			authField.SetValue("false")
			keyField := ecs.formManager.GetField(components.FieldNameKey)
			if keyField != nil {
				if textInput, ok := keyField.GetTextInput(); ok {
					textInput.SetValue(ecs.connection.KeyPath)
					keyField.SetTextInput(textInput)
				}
			}
		} else {
			authField.SetValue("true")
			if ecs.connection.HasPassword && ecs.connection.Password != "" {
				passwordField := ecs.formManager.GetField(components.FieldNamePassword)
				if passwordField != nil {
					if textInput, ok := passwordField.GetTextInput(); ok {
						textInput.SetValue(ecs.connection.Password)
						passwordField.SetTextInput(textInput)
					}
				}
			}
		}
	}

	// Обновляем видимость полей
	ecs.updateFieldVisibility()

	// Принудительно обновляем фокус полей
	ecs.formManager.UpdateFocus()
}

// updateFieldVisibility обновляет видимость полей на основе выбора аутентификации
func (ecs *EditConnectionScreen) updateFieldVisibility() {
	// Получаем значение из FormManager
	authField := ecs.formManager.GetField(components.FieldNameAuth)
	usePassword := authField.Value() == "true"

	// Обновляем видимость в менеджере полей
	ecs.formManager.GetField(components.FieldNamePassword).SetVisible(usePassword)
	ecs.formManager.GetField(components.FieldNameKey).SetVisible(!usePassword)
}

// Update обрабатывает обновления состояния
func (ecs *EditConnectionScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ecs.SetSize(msg.Width, msg.Height)
		if !ecs.ready {
			// Инициализируем viewport при первом изменении размера
			ecs.viewport = viewport.New(msg.Width-4, msg.Height-12)
			ecs.ready = true
			ecs.updateViewportContent()
		} else {
			// Используем одинаковый размер для консистентности
			ecs.viewport.Width = msg.Width - 4
			ecs.viewport.Height = msg.Height - 12
		}
		return ecs, nil

	case ui.NavigateToMsg:
		// Очищаем форму при навигации к другому экрану
		if msg.ScreenName != "edit_connection" {
			ecs.clearForm()
		}
		return ecs, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return ecs, tea.Quit
		case "esc":
			// Очищаем форму и возвращаемся
			ecs.clearForm()
			return ecs, ui.GoBackCmd()
		case "tab":
			// Переходим к следующему полю
			ecs.formManager.NextField()
		case "shift+tab":
			// Переходим к предыдущему полю
			ecs.formManager.PrevField()
		case "enter":
			// Проверяем, какая кнопка была нажата
			currentFieldName := ecs.formManager.GetCurrentField()
			currentField := ecs.formManager.GetField(currentFieldName)

			if currentField != nil && currentField.IsButton() {
				buttonName := currentField.GetName()
				switch buttonName {
				case "save":
					// Сохраняем изменения
					return ecs, ecs.saveConnection()
				case "delete":
					// Удаляем подключение
					return ecs, ecs.deleteConnection()
				}
			} else if ecs.formManager.IsLastField() {
				// Если это последнее поле (не кнопка), переходим к первой кнопке
				ecs.formManager.NextField()
			} else {
				// Переходим к следующему полю
				ecs.formManager.NextField()
			}
		case "ctrl+s":
			// Сохранить изменения
			return ecs, ecs.saveConnection()
		case "ctrl+d":
			// Удалить подключение
			return ecs, ecs.deleteConnection()
		case "space", "left", "right":
			// Обрабатываем булевое поле через FormManager
			if ecs.formManager.GetCurrentField() == components.FieldNameAuth {
				// Обновляем поле в FormManager напрямую
				field := ecs.formManager.GetField(components.FieldNameAuth)
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
			ecs.viewport.ScrollUp(1)
		case "down":
			// Прокрутка вниз
			ecs.viewport.ScrollDown(1)
		case "pageup":
			// Прокрутка на страницу вверх
			ecs.viewport.PageUp()
		case "pagedown":
			// Прокрутка на страницу вниз
			ecs.viewport.PageDown()
		}
	}

	// Обновляем все поля через менеджер формы
	for _, fieldName := range ecs.formManager.GetFieldOrder() {
		field := ecs.formManager.GetField(fieldName)
		_, fieldCmd := field.Update(msg)
		if fieldCmd != nil {
			if teaCmd, ok := fieldCmd.(tea.Cmd); ok {
				cmd = teaCmd
			}
		}
	}

	// Обновляем фокус через менеджер формы
	ecs.formManager.UpdateFocus()

	// Обновляем видимость полей на основе выбора аутентификации
	ecs.updateFieldVisibility()

	// Обновляем содержимое viewport
	ecs.updateViewportContent()

	// Прокручиваем viewport для видимости текущего поля после навигации
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "tab", "shift+tab", "enter":
			// Добавляем команду для прокрутки после обновления содержимого
			cmd = tea.Sequence(cmd, ecs.scrollToCurrentFieldCmd())
		}
	}

	// Обновляем базовый экран
	baseScreen, baseCmd := ecs.BaseScreen.Update(msg)
	ecs.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return ecs, cmd
}

// View возвращает строку для отрисовки
func (ecs *EditConnectionScreen) View() string {
	if !ecs.ready {
		return "Инициализация..."
	}

	// Подготавливаем содержимое viewport
	viewportContent := ecs.viewport.View()

	// Добавляем сообщения в начало
	messages := ecs.messageManager.RenderMessages(80) // Используем фиксированную ширину
	if messages != "" {
		viewportContent = messages + viewportContent
	}

	// Добавляем индикатор прокрутки под viewport если не дошли до конца
	if !ecs.viewport.AtBottom() {
		// Создаем более информативный индикатор прокрутки
		scrollIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorMuted)).
			Align(lipgloss.Center).
			Render("↓ Прокрутите вниз для просмотра всех полей ↓")

		// Добавляем индикатор под содержимое viewport
		viewportContent += "\n" + scrollIndicator
	}

	// Устанавливаем содержимое с индикатором
	ecs.SetContent(viewportContent)
	return ecs.BaseScreen.View()
}

// updateViewportContent обновляет содержимое viewport
func (ecs *EditConnectionScreen) updateViewportContent() {
	if !ecs.ready {
		// Если viewport еще не готов, инициализируем его с минимальными размерами
		ecs.viewport = viewport.New(80, 20) // Минимальные размеры
		ecs.ready = true
	}
	content := ecs.renderForm()
	ecs.viewport.SetContent(content)
}

// renderForm рендерит форму в виде вертикального списка
func (ecs *EditConnectionScreen) renderForm() string {
	// Используем менеджер формы для рендеринга
	return ecs.formManager.RenderForm()
}

// saveConnection сохраняет изменения подключения
func (ecs *EditConnectionScreen) saveConnection() tea.Cmd {
	// Валидируем все поля
	if !ecs.formManager.ValidateAll() {
		// Показываем сообщение об ошибке валидации
		ecs.messageManager.AddError("❌ Пожалуйста, заполните все обязательные поля")
		return nil
	}

	// Получаем значения полей
	values := ecs.formManager.GetValues()

	// Обновляем подключение
	port := 22 // По умолчанию
	if values[components.FieldNamePort] != "" {
		if p, err := strconv.Atoi(values[components.FieldNamePort]); err == nil {
			port = p
		}
	}

	// Обновляем поля подключения
	ecs.connection.Name = values[components.FieldNameName]
	ecs.connection.Host = values[components.FieldNameHost]
	ecs.connection.Port = port
	ecs.connection.User = values[components.FieldNameUser]
	ecs.connection.KeyPath = values[components.FieldNameKey]
	ecs.connection.UseSSHKey = !(values[components.FieldNameAuth] == "true") // Если не пароль, то SSH ключ
	ecs.connection.HasPassword = values[components.FieldNameAuth] == "true" && values[components.FieldNamePassword] != ""
	ecs.connection.UpdatedAt = time.Now()

	// Обновляем пароль если используется
	if ecs.connection.HasPassword {
		ecs.connection.Password = values[components.FieldNamePassword]
	} else {
		ecs.connection.Password = ""
	}

	// Сохраняем изменения
	ecs.messageManager.AddInfo(fmt.Sprintf("Сохраняем подключение: %s (ID: %s)", ecs.connection.Name, ecs.connection.ID))
	err := ecs.connectionSvc.UpdateConnection(ecs.connection)
	if err != nil {
		// Показываем ошибку сохранения
		ecs.messageManager.AddError(fmt.Sprintf("Ошибка сохранения: %v", err))
		return nil
	}

	// Очищаем ошибки
	ecs.errors = make(map[string]string)

	// Добавляем сообщение об успехе
	ecs.messageManager.AddSuccess(fmt.Sprintf("Подключение '%s' успешно обновлено!", ecs.connection.Name))

	// Возвращаемся к списку подключений
	return ui.NavigateToCmd("connections")
}

// deleteConnection удаляет подключение
func (ecs *EditConnectionScreen) deleteConnection() tea.Cmd {
	if ecs.connection == nil {
		ecs.messageManager.AddError("Ошибка: подключение не инициализировано")
		return nil
	}

	// Удаляем подключение через сервис
	err := ecs.connectionSvc.DeleteConnection(ecs.connection.ID)
	if err != nil {
		// Показываем ошибку удаления
		ecs.messageManager.AddError(fmt.Sprintf("Ошибка удаления: %v", err))
		return nil
	}

	// Добавляем сообщение об успехе
	ecs.messageManager.AddSuccess(fmt.Sprintf("Подключение '%s' успешно удалено!", ecs.connection.Name))

	// Возвращаемся к списку подключений
	return ui.NavigateToCmd("connections")
}

// clearForm очищает все поля формы
func (ecs *EditConnectionScreen) clearForm() {
	// Очищаем все поля через FormManager
	for _, fieldName := range ecs.formManager.GetFieldOrder() {
		field := ecs.formManager.GetField(fieldName)
		if field != nil {
			field.SetValue("")
		}
	}

	// Очищаем ошибки
	ecs.errors = make(map[string]string)

	// Устанавливаем фокус на первое поле
	ecs.formManager.SetCurrentField(components.FieldNameName)
	ecs.formManager.UpdateFocus()

	// Сбрасываем прокрутку в начало
	if ecs.ready {
		ecs.viewport.SetYOffset(0)
		ecs.updateViewportContent()
	}
}

// scrollToCurrentField прокручивает viewport для видимости текущего поля
func (ecs *EditConnectionScreen) scrollToCurrentField() {
	if !ecs.ready {
		return
	}

	// Получаем текущее поле
	currentField := ecs.formManager.GetCurrentField()
	fieldOrder := ecs.formManager.GetFieldOrder()

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
	viewportHeight := ecs.viewport.Height
	totalContentHeight := ecs.viewport.TotalLineCount()

	// Если содержимое помещается в viewport, не прокручиваем
	if totalContentHeight <= viewportHeight {
		ecs.viewport.SetYOffset(0)
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
	ecs.viewport.SetYOffset(optimalOffset)
}

// scrollToCurrentFieldCmd возвращает команду для прокрутки к текущему полю
func (ecs *EditConnectionScreen) scrollToCurrentFieldCmd() tea.Cmd {
	return func() tea.Msg {
		ecs.scrollToCurrentField()
		return nil
	}
}

// Init инициализирует экран
func (ecs *EditConnectionScreen) Init() tea.Cmd {
	// Устанавливаем фокус на первое поле через FormManager
	ecs.formManager.SetCurrentField(components.FieldNameName)
	ecs.formManager.UpdateFocus()
	return textinput.Blink
}

// GetName возвращает имя экрана
func (ecs *EditConnectionScreen) GetName() string {
	return "edit_connection"
}
