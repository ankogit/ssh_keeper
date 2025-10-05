package screens

import (
	"fmt"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"ssh-keeper/internal/ui/styles"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WelcomeScreen представляет экран приветствия и ввода мастер-пароля
type WelcomeScreen struct {
	*BaseScreen
	masterPasswordService *services.MasterPasswordService
	encryptionService     *services.EncryptionService

	// Менеджеры
	formManager    *components.FormManager
	messageManager *components.MessageManager

	// Прокрутка
	viewport viewport.Model
	ready    bool

	// Состояние
	isFirstTime bool
	currentStep int // 0 - ввод пароля, 1 - подтверждение, 2 - готово
}

// NewWelcomeScreen создает новый экран приветствия
func NewWelcomeScreen() *WelcomeScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Добро пожаловать!")

	// Используем глобальные сервисы
	masterPasswordService := services.GetGlobalMasterPasswordService()
	encryptionService := services.GetGlobalEncryptionService()

	// Создаем менеджеры
	formManager := components.NewFormManager()
	messageManager := components.NewMessageManager()

	// Создаем viewport для прокрутки
	vp := viewport.New(0, 0)

	// Определяем, первый ли это запуск - используем глобальную функцию с проверкой подписи
	isFirstTime := !services.IsMasterPasswordInitializedWithSignature()

	screen := &WelcomeScreen{
		BaseScreen:            baseScreen,
		masterPasswordService: masterPasswordService,
		encryptionService:     encryptionService,
		formManager:           formManager,
		messageManager:        messageManager,
		viewport:              vp,
		ready:                 false,
		currentStep:           0,
		isFirstTime:           isFirstTime,
	}

	// Настраиваем форму в зависимости от того, первый ли это запуск
	if isFirstTime {
		screen.setupFirstTimeForm()
	} else {
		screen.setupExistingForm()
	}

	// Устанавливаем фокус на первое поле
	screen.formManager.SetCurrentField("password")
	screen.formManager.UpdateFocus()

	return screen
}

// setupFirstTimeForm настраивает форму для первого запуска
func (ws *WelcomeScreen) setupFirstTimeForm() {
	// Поле ввода мастер-пароля
	ws.formManager.AddField(components.FieldConfig{
		Name:        "password",
		Label:       "Мастер-пароль",
		Placeholder: "Введите мастер-пароль",
		FieldType:   components.FieldTypePassword,
		Required:    true,
		Width:       50,
	})

	// Поле подтверждения пароля
	ws.formManager.AddField(components.FieldConfig{
		Name:        "confirm",
		Label:       "Подтверждение",
		Placeholder: "Подтвердите мастер-пароль",
		FieldType:   components.FieldTypePassword,
		Required:    true,
		Width:       50,
	})

	// Кнопка подтверждения
	ws.formManager.AddField(components.FieldConfig{
		Name:      "submit",
		Label:     "Установить мастер-пароль",
		FieldType: components.FieldTypeButton,
		Required:  false,
		Width:     30,
	})
}

// setupExistingForm настраивает форму для существующего пользователя
func (ws *WelcomeScreen) setupExistingForm() {
	// Поле ввода мастер-пароля
	ws.formManager.AddField(components.FieldConfig{
		Name:        "password",
		Label:       "Мастер-пароль",
		Placeholder: "Введите мастер-пароль",
		FieldType:   components.FieldTypePassword,
		Required:    true,
		Width:       50,
	})

	// Кнопка входа
	ws.formManager.AddField(components.FieldConfig{
		Name:      "submit",
		Label:     "Войти",
		FieldType: components.FieldTypeButton,
		Required:  false,
		Width:     20,
	})
}

// Init инициализирует экран
func (ws *WelcomeScreen) Init() tea.Cmd {
	// Устанавливаем фокус на первое поле через FormManager
	ws.formManager.SetCurrentField("password")
	ws.formManager.UpdateFocus()
	return textinput.Blink
}

// Update обновляет состояние экрана
func (ws *WelcomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ws.SetSize(msg.Width, msg.Height)
		if !ws.ready {
			// Инициализируем viewport при первом изменении размера
			// Учитываем все служебные элементы: заголовки, рамки, отступы
			ws.viewport = viewport.New(msg.Width-4, msg.Height-12)
			ws.ready = true
			ws.updateViewportContent()
		} else {
			// Используем одинаковый размер для консистентности
			ws.viewport.Width = msg.Width - 4
			ws.viewport.Height = msg.Height - 12
		}
		return ws, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return ws, tea.Quit
		case "esc":
			return ws, tea.Quit
		case "ctrl+p":
			// Переключаем видимость пароля
			ws.togglePasswordVisibility()
		case "tab":
			// Переходим к следующему полю
			ws.formManager.NextField()
		case "shift+tab":
			// Переходим к предыдущему полю
			ws.formManager.PrevField()
		case "enter":
			// Проверяем, является ли текущее поле кнопкой
			if currentField := ws.formManager.GetCurrentFieldModel(); currentField != nil && currentField.IsButton() {
				return ws, ws.handleEnter()
			} else {
				// Если это не кнопка, переходим к следующему полю
				ws.formManager.NextField()
			}
		case "up":
			// Прокрутка вверх
			ws.viewport.ScrollUp(1)
		case "down":
			// Прокрутка вниз
			ws.viewport.ScrollDown(1)
		case "pageup":
			// Прокрутка на страницу вверх
			ws.viewport.PageUp()
		case "pagedown":
			// Прокрутка на страницу вниз
			ws.viewport.PageDown()
		}
	}

	// Обновляем все поля через менеджер формы
	for _, fieldName := range ws.formManager.GetFieldOrder() {
		field := ws.formManager.GetField(fieldName)
		_, fieldCmd := field.Update(msg)
		if fieldCmd != nil {
			if teaCmd, ok := fieldCmd.(tea.Cmd); ok {
				cmd = teaCmd
			}
		}
	}

	// Обновляем фокус через менеджер формы
	ws.formManager.UpdateFocus()

	// Обновляем содержимое viewport
	ws.updateViewportContent()

	// Обновляем базовый экран
	baseScreen, baseCmd := ws.BaseScreen.Update(msg)
	ws.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return ws, cmd
}

// View отображает экран
func (ws *WelcomeScreen) View() string {
	if !ws.ready {
		return "Инициализация..."
	}

	// Подготавливаем содержимое viewport
	viewportContent := ws.viewport.View()

	// Добавляем индикатор прокрутки под viewport если не дошли до конца
	if !ws.viewport.AtBottom() {
		// Создаем более информативный индикатор прокрутки
		scrollIndicator := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorMuted)).
			Align(lipgloss.Center).
			Render("↓ Прокрутите вниз для просмотра всех полей ↓")

		// Добавляем индикатор под содержимое viewport
		viewportContent += "\n" + scrollIndicator
	}

	// Добавляем сообщения в конец
	messages := ws.messageManager.RenderMessages(80) // Используем фиксированную ширину
	if messages != "" {
		viewportContent += "\n" + messages
	}

	// Устанавливаем содержимое с индикатором
	ws.SetContent(viewportContent)
	return ws.BaseScreen.View()
}

// updateViewportContent обновляет содержимое viewport
func (ws *WelcomeScreen) updateViewportContent() {
	content := ws.renderForm()
	ws.viewport.SetContent(content)
}

// renderForm рендерит форму в виде вертикального списка
func (ws *WelcomeScreen) renderForm() string {
	if ws.isFirstTime {
		return ws.renderFirstTimeView()
	} else {
		return ws.renderExistingView()
	}
}

// renderFirstTimeView отображает форму для первого запуска
func (ws *WelcomeScreen) renderFirstTimeView() string {
	var content string

	content += styles.TitleStyle.Render("Добро пожаловать в SSH Keeper!") + "\n\n"
	content += styles.SubtitleStyle.Render("Установите мастер-пароль для защиты ваших SSH подключений.") + "\n"

	// Отображаем поля формы
	content += ws.formManager.RenderForm()

	return content
}

// renderExistingView отображает форму для существующего пользователя
func (ws *WelcomeScreen) renderExistingView() string {
	var content string

	content += styles.TitleStyle.Render("Добро пожаловать обратно!") + "\n\n"
	content += styles.SubtitleStyle.Render("Введите мастер-пароль для доступа к вашим SSH подключениям.") + "\n\n"

	// Отображаем поля формы
	content += ws.formManager.RenderForm()

	content += "\n\n" + styles.HelpStyle.Render("Нажмите Enter для входа, Esc для выхода")

	return content
}

// handleEnter обрабатывает нажатие Enter
func (ws *WelcomeScreen) handleEnter() tea.Cmd {
	if ws.isFirstTime {
		return ws.handleFirstTimeSetup()
	} else {
		return ws.handleExistingSetup()
	}
}

// handleFirstTimeSetup обрабатывает настройку для первого запуска
func (ws *WelcomeScreen) handleFirstTimeSetup() tea.Cmd {
	values := ws.formManager.GetValues()
	password := values["password"]
	confirm := values["confirm"]

	// Проверяем, что пароли совпадают
	if password != confirm {
		ws.messageManager.AddError("Пароли не совпадают")
		return nil
	}

	// Проверяем длину пароля
	if len(password) < 6 {
		ws.messageManager.AddError("Пароль должен содержать минимум 6 символов")
		return nil
	}

	// Сохраняем мастер-пароль с проверкой подписи
	err := services.SetMasterPasswordWithSignature(password)
	if err != nil {
		ws.messageManager.AddError(fmt.Sprintf("Ошибка сохранения: %v", err))
		return nil
	}

	// Обновляем ключ шифрования
	if ws.encryptionService != nil {
		err = ws.encryptionService.RefreshKey()
		if err != nil {
			ws.messageManager.AddError(fmt.Sprintf("Ошибка обновления ключа: %v", err))
			return nil
		}
	}

	ws.messageManager.AddSuccess("Мастер-пароль успешно установлен!")

	// Переходим к главному меню через небольшую задержку
	return tea.Sequence(
		tea.Printf("Мастер-пароль установлен успешно!"),
		func() tea.Msg {
			time.Sleep(1 * time.Second)
			return tea.Msg(nil)
		},
		func() tea.Msg {
			return ui.NavigateToMsg{ScreenName: "main_menu"}
		},
	)
}

// handleExistingSetup обрабатывает вход для существующего пользователя
func (ws *WelcomeScreen) handleExistingSetup() tea.Cmd {
	values := ws.formManager.GetValues()
	password := values["password"]

	// Проверяем пароль с проверкой подписи
	storedPassword, err := services.GetMasterPasswordWithSignature()
	if err != nil {
		ws.messageManager.AddError(fmt.Sprintf("Ошибка получения пароля: %v", err))
		return nil
	}

	if password != storedPassword {
		ws.messageManager.AddError("Неверный мастер-пароль")
		return nil
	}

	// Обновляем ключ шифрования
	if ws.encryptionService != nil {
		err = ws.encryptionService.RefreshKey()
		if err != nil {
			ws.messageManager.AddError(fmt.Sprintf("Ошибка обновления ключа: %v", err))
			return nil
		}
	}

	ws.messageManager.AddSuccess("Добро пожаловать!")

	// Переходим к главному меню через небольшую задержку
	return tea.Sequence(
		tea.Printf("Вход выполнен успешно!"),
		func() tea.Msg {
			time.Sleep(1 * time.Second)
			return tea.Msg(nil)
		},
		func() tea.Msg {
			return ui.NavigateToMsg{ScreenName: "main_menu"}
		},
	)
}

// togglePasswordVisibility переключает видимость пароля
func (ws *WelcomeScreen) togglePasswordVisibility() {
	currentFieldName := ws.formManager.GetCurrentField()
	currentField := ws.formManager.GetField(currentFieldName)
	if currentField == nil {
		return
	}

	// Получаем textinput.Model из FormField
	if textInput, ok := currentField.GetTextInput(); ok {
		if textInput.EchoMode == textinput.EchoPassword {
			textInput.EchoMode = textinput.EchoNormal
		} else {
			textInput.EchoMode = textinput.EchoPassword
		}
		// Обновляем поле в FormField
		currentField.SetTextInput(textInput)
	}
}
