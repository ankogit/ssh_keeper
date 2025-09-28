package screens

import (
	"fmt"
	"ssh-keeper/internal/models"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"ssh-keeper/internal/ui/styles"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// Состояние формы
	currentField int
	fields       []textinput.Model

	// Условные поля
	showPasswordField bool
	showKeyField      bool

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

	// Создаем массив полей для удобного переключения
	fields := []textinput.Model{
		nameInput,
		hostInput,
		portInput,
		userInput,
		keyPathInput,
		passwordInput,
	}

	// Создаем булевое поле
	usePasswordField := components.NewBoolField("Использовать пароль")
	usePasswordField.SetWidth(20)

	// Создаем viewport для прокрутки
	vp := viewport.New(0, 0)

	return &AddConnectionScreen{
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
		currentField:  0,
		fields:        fields,
		viewport:      vp,
		ready:         false,
	}
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
			acs.nextField()
		case "shift+tab":
			// Переходим к предыдущему полю
			acs.prevField()
		case "enter":
			if acs.currentField == 5 {
				// Последнее поле - сохраняем
				acs.saveConnection()
			} else {
				// Переходим к следующему полю
				acs.nextField()
			}
		case "ctrl+s":
			// Сохранить подключение
			acs.saveConnection()
		case "ctrl+p":
			// Переключить показ пароля
			acs.togglePasswordVisibility()
		case "space", "left", "right":
			// Обрабатываем булевое поле
			if acs.currentField == 4 {
				acs.usePassword, _ = acs.usePassword.Update(msg)
			}
		case "up", "k":
			// Прокрутка вверх
			acs.viewport.ScrollUp(1)
		case "down", "j":
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

	// Обновляем поля ПОСЛЕ обработки KeyMsg (как работало)
	acs.nameInput, cmd = acs.nameInput.Update(msg)
	acs.hostInput, _ = acs.hostInput.Update(msg)
	acs.portInput, _ = acs.portInput.Update(msg)
	acs.userInput, _ = acs.userInput.Update(msg)

	// Обновляем условные поля только если они видимы
	if acs.showPasswordField {
		acs.passwordInput, _ = acs.passwordInput.Update(msg)
	}
	if acs.showKeyField {
		acs.keyPathInput, _ = acs.keyPathInput.Update(msg)
	}

	// Убеждаемся, что текущее поле имеет фокус
	switch acs.currentField {
	case 0:
		acs.nameInput.Focus()
		acs.hostInput.Blur()
		acs.portInput.Blur()
		acs.userInput.Blur()
		acs.usePassword.Blur()
		acs.passwordInput.Blur()
		acs.keyPathInput.Blur()
	case 1:
		acs.nameInput.Blur()
		acs.hostInput.Focus()
		acs.portInput.Blur()
		acs.userInput.Blur()
		acs.usePassword.Blur()
		acs.passwordInput.Blur()
		acs.keyPathInput.Blur()
	case 2:
		acs.nameInput.Blur()
		acs.hostInput.Blur()
		acs.portInput.Focus()
		acs.userInput.Blur()
		acs.usePassword.Blur()
		acs.passwordInput.Blur()
		acs.keyPathInput.Blur()
	case 3:
		acs.nameInput.Blur()
		acs.hostInput.Blur()
		acs.portInput.Blur()
		acs.userInput.Focus()
		acs.usePassword.Blur()
		acs.passwordInput.Blur()
		acs.keyPathInput.Blur()
	case 4:
		acs.nameInput.Blur()
		acs.hostInput.Blur()
		acs.portInput.Blur()
		acs.userInput.Blur()
		acs.usePassword.Focus()
		acs.passwordInput.Blur()
		acs.keyPathInput.Blur()
	case 5:
		acs.nameInput.Blur()
		acs.hostInput.Blur()
		acs.portInput.Blur()
		acs.userInput.Blur()
		acs.usePassword.Blur()
		if acs.showPasswordField {
			acs.passwordInput.Focus()
		} else if acs.showKeyField {
			acs.keyPathInput.Focus()
		}
	}

	// Обновляем видимость полей на основе выбора аутентификации
	acs.updateFieldVisibility()

	// Обновляем содержимое viewport
	acs.updateViewportContent()

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
	usePassword := acs.usePassword.Value()
	acs.showPasswordField = usePassword
	acs.showKeyField = !usePassword
}

// updateViewportContent обновляет содержимое viewport
func (acs *AddConnectionScreen) updateViewportContent() {
	content := acs.renderForm()
	acs.viewport.SetContent(content)
}

// renderForm рендерит форму в виде вертикального списка
func (acs *AddConnectionScreen) renderForm() string {
	// Стили
	fieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(styles.ColorPrimary)).
		Width(50)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorSecondary)).
		Bold(styles.TextBold).
		Width(15)

	instructionsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted)).
		Italic(styles.TextItalic)

	// Создаем два поля
	var formContent []string

	// Название
	nameLabel := labelStyle.Render("Название:")
	nameField := fieldStyle.Render(acs.nameInput.View())
	if _, exists := acs.errors["name"]; exists {
		// Поле с ошибкой - красная рамка
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorError)).
			Width(50).
			Render(acs.nameInput.View())
	} else if acs.currentField == 0 {
		// Активное поле - оранжевая рамка
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorWarning)).
			Width(50).
			Render(acs.nameInput.View())
	}
	formContent = append(formContent, lipgloss.JoinHorizontal(lipgloss.Center, nameLabel, nameField))
	formContent = append(formContent, "") // Пустая строка

	// Хост
	hostLabel := labelStyle.Render("Хост:")
	hostField := fieldStyle.Render(acs.hostInput.View())
	if _, exists := acs.errors["host"]; exists {
		// Поле с ошибкой - красная рамка
		hostField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorError)).
			Width(50).
			Render(acs.hostInput.View())
	} else if acs.currentField == 1 {
		// Активное поле - оранжевая рамка
		hostField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorWarning)).
			Width(50).
			Render(acs.hostInput.View())
	}
	formContent = append(formContent, lipgloss.JoinHorizontal(lipgloss.Center, hostLabel, hostField))
	formContent = append(formContent, "") // Пустая строка

	// Порт
	portLabel := labelStyle.Render("Порт:")
	portFieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(styles.ColorPrimary)).
		Padding(0, 1).
		Width(15)

	portField := portFieldStyle.Render(acs.portInput.View())
	if _, exists := acs.errors["port"]; exists {
		// Поле с ошибкой - красная рамка
		portField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorError)).
			Padding(0, 1).
			Width(15).
			Render(acs.portInput.View())
	} else if acs.currentField == 2 {
		// Активное поле - оранжевая рамка
		portField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorWarning)).
			Padding(0, 1).
			Width(15).
			Render(acs.portInput.View())
	}
	formContent = append(formContent, lipgloss.JoinHorizontal(lipgloss.Center, portLabel, portField))
	formContent = append(formContent, "") // Пустая строка

	// Пользователь
	userLabel := labelStyle.Render("Пользователь:")
	userField := fieldStyle.Render(acs.userInput.View())
	if _, exists := acs.errors["user"]; exists {
		// Поле с ошибкой - красная рамка
		userField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorError)).
			Padding(0, 1).
			Width(50).
			Render(acs.userInput.View())
	} else if acs.currentField == 3 {
		// Активное поле - оранжевая рамка
		userField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.ColorWarning)).
			Padding(0, 1).
			Width(50).
			Render(acs.userInput.View())
	}
	formContent = append(formContent, lipgloss.JoinHorizontal(lipgloss.Center, userLabel, userField))
	formContent = append(formContent, "") // Пустая строка

	// Булевое поле - использование пароля
	authLabel := labelStyle.Render("Использовать пароль: (←/→)")
	authField := acs.usePassword.View()
	formContent = append(formContent, lipgloss.JoinHorizontal(lipgloss.Center, authLabel, authField))
	formContent = append(formContent, "") // Пустая строка

	// Условные поля в зависимости от выбора аутентификации
	if acs.showPasswordField {
		// Поле пароля
		passwordLabel := labelStyle.Render("Пароль:")
		passwordField := fieldStyle.Render(acs.passwordInput.View())
		if _, exists := acs.errors["password"]; exists {
			// Поле с ошибкой - красная рамка
			passwordField = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(styles.ColorError)).
				Padding(0, 1).
				Width(50).
				Render(acs.passwordInput.View())
		} else if acs.currentField == 5 {
			// Активное поле - оранжевая рамка
			passwordField = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(styles.ColorWarning)).
				Padding(0, 1).
				Width(50).
				Render(acs.passwordInput.View())
		}
		formContent = append(formContent, lipgloss.JoinHorizontal(lipgloss.Center, passwordLabel, passwordField))
		formContent = append(formContent, "") // Пустая строка
	}

	if acs.showKeyField {
		// Поле SSH ключа
		keyLabel := labelStyle.Render("SSH ключ:")
		keyField := fieldStyle.Render(acs.keyPathInput.View())
		if _, exists := acs.errors["keyPath"]; exists {
			// Поле с ошибкой - красная рамка
			keyField = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(styles.ColorError)).
				Padding(0, 1).
				Width(50).
				Render(acs.keyPathInput.View())
		} else if acs.currentField == 5 {
			// Активное поле - оранжевая рамка
			keyField = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(styles.ColorWarning)).
				Padding(0, 1).
				Width(50).
				Render(acs.keyPathInput.View())
		}
		formContent = append(formContent, lipgloss.JoinHorizontal(lipgloss.Center, keyLabel, keyField))
		formContent = append(formContent, "") // Пустая строка
	}

	// Инструкции
	instructions := instructionsStyle.Render("Tab - след. поле • ↑/↓ - прокрутка • Enter - сохранить • Esc - назад")
	formContent = append(formContent, instructions)

	// Объединяем все в вертикальный список
	content := lipgloss.JoinVertical(lipgloss.Left, formContent...)

	return content
}

// nextField переходит к следующему полю
func (acs *AddConnectionScreen) nextField() {
	switch acs.currentField {
	case 0:
		acs.nameInput.Blur()
		acs.currentField = 1
		acs.hostInput.Focus()
	case 1:
		acs.hostInput.Blur()
		acs.currentField = 2
		acs.portInput.Focus()
	case 2:
		acs.portInput.Blur()
		acs.currentField = 3
		acs.userInput.Focus()
	case 3:
		acs.userInput.Blur()
		acs.currentField = 4
		acs.usePassword.Focus()
	case 4:
		acs.usePassword.Blur()
		acs.currentField = 5
		if acs.showPasswordField {
			acs.passwordInput.Focus()
		} else if acs.showKeyField {
			acs.keyPathInput.Focus()
		}
	case 5:
		acs.passwordInput.Blur()
		acs.keyPathInput.Blur()
		acs.currentField = 0
		acs.nameInput.Focus()
	}
}

// prevField переходит к предыдущему полю
func (acs *AddConnectionScreen) prevField() {
	switch acs.currentField {
	case 0:
		acs.nameInput.Blur()
		acs.currentField = 5
		if acs.showPasswordField {
			acs.passwordInput.Focus()
		} else if acs.showKeyField {
			acs.keyPathInput.Focus()
		}
	case 1:
		acs.hostInput.Blur()
		acs.currentField = 0
		acs.nameInput.Focus()
	case 2:
		acs.portInput.Blur()
		acs.currentField = 1
		acs.hostInput.Focus()
	case 3:
		acs.userInput.Blur()
		acs.currentField = 2
		acs.portInput.Focus()
	case 4:
		acs.usePassword.Blur()
		acs.currentField = 3
		acs.userInput.Focus()
	case 5:
		acs.passwordInput.Blur()
		acs.keyPathInput.Blur()
		acs.currentField = 4
		acs.usePassword.Focus()
	}
}

// togglePasswordVisibility переключает видимость пароля
func (acs *AddConnectionScreen) togglePasswordVisibility() {
	if acs.passwordInput.EchoMode == textinput.EchoPassword {
		acs.passwordInput.EchoMode = textinput.EchoNormal
	} else {
		acs.passwordInput.EchoMode = textinput.EchoPassword
	}
}

// saveConnection сохраняет подключение
func (acs *AddConnectionScreen) saveConnection() {
	// Очищаем предыдущие ошибки
	acs.errors = make(map[string]string)

	// Валидация обязательных полей
	if acs.nameInput.Value() == "" {
		acs.errors["name"] = "Название обязательно"
	}
	if acs.hostInput.Value() == "" {
		acs.errors["host"] = "Хост обязателен"
	}
	if acs.userInput.Value() == "" {
		acs.errors["user"] = "Пользователь обязателен"
	}

	// Валидация условных полей
	if acs.showPasswordField && acs.passwordInput.Value() == "" {
		acs.errors["password"] = "Пароль обязателен"
	}
	// SSH ключ не обязателен - убираем валидацию

	// Если есть ошибки, не сохраняем
	if len(acs.errors) > 0 {
		return
	}

	// Выводим все данные формы в консоль
	fmt.Println("\n=== ДАННЫЕ ФОРМЫ ===")
	fmt.Printf("Название: %s\n", acs.nameInput.Value())
	fmt.Printf("Хост: %s\n", acs.hostInput.Value())
	fmt.Printf("Порт: %s\n", acs.portInput.Value())
	fmt.Printf("Пользователь: %s\n", acs.userInput.Value())
	fmt.Printf("Использовать пароль: %t\n", acs.usePassword.Value())

	if acs.showPasswordField {
		fmt.Printf("Пароль: %s\n", acs.passwordInput.Value())
	}
	if acs.showKeyField {
		fmt.Printf("SSH ключ: %s\n", acs.keyPathInput.Value())
	}
	fmt.Println("==================\n")

	// Создаем подключение
	port := 22 // Значение по умолчанию
	if acs.portInput.Value() != "" {
		if parsedPort, err := strconv.Atoi(acs.portInput.Value()); err == nil {
			port = parsedPort
		}
	}
	_ = models.Connection{
		Name:        acs.nameInput.Value(),
		Host:        acs.hostInput.Value(),
		Port:        port,
		User:        acs.userInput.Value(),
		KeyPath:     acs.keyPathInput.Value(),
		HasPassword: acs.usePassword.Value(),
	}

	// TODO: Сохранить в сервисе
	// Возвращаемся к списку подключений
	// В реальном приложении здесь будет команда навигации
}

// Init инициализирует экран
func (acs *AddConnectionScreen) Init() tea.Cmd {
	// Устанавливаем фокус на поле названия
	acs.nameInput.Focus()
	return textinput.Blink
}

// GetName возвращает имя экрана
func (acs *AddConnectionScreen) GetName() string {
	return "add_connection"
}
