package screens

import (
	"fmt"
	"os"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ImportScreen представляет экран импорта конфигурации
type ImportScreen struct {
	*BaseScreen
	formManager       *components.FormManager
	messageManager    *components.MessageManager
	connectionService *services.ConnectionService
	importPath        string
	isImporting       bool
}

// NewImportScreen создает новый экран импорта
func NewImportScreen() *ImportScreen {
	baseScreen := NewBaseScreen("Импорт конфигурации")

	// Создаем менеджер формы
	formManager := components.NewFormManager()

	// Добавляем поля формы
	formManager.AddField(components.FieldConfig{
		Name:        "import_path",
		Label:       "Путь к файлу",
		Required:    true,
		Width:       50,
		Placeholder: "/path/to/config_file",
		FieldType:   components.FieldTypeText,
	})

	formManager.AddField(components.FieldConfig{
		Name:      "import_button",
		Label:     "Импорт",
		FieldType: components.FieldTypeButton,
		Style:     "success",
	})

	// Устанавливаем фокус на первое поле
	formManager.SetCurrentField("import_path")
	formManager.UpdateFocus()

	// Создаем менеджер сообщений
	messageManager := components.NewMessageManager()

	// Получаем сервис подключений
	connectionService := services.GetGlobalConnectionService()

	return &ImportScreen{
		BaseScreen:        baseScreen,
		formManager:       formManager,
		messageManager:    messageManager,
		connectionService: connectionService,
		importPath:        "",
		isImporting:       false,
	}
}

// Update обрабатывает обновления состояния
func (is *ImportScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		is.SetSize(msg.Width, msg.Height)
		return is, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return is, ui.GoBackCmd()

		case "tab":
			is.formManager.NextField()
			is.formManager.UpdateFocus()

		case "shift+tab":
			is.formManager.PrevField()
			is.formManager.UpdateFocus()

		case "enter":
			currentField := is.formManager.GetCurrentField()
			if currentField == "import_button" {
				return is, is.performImport()
			} else {
				// Если это не кнопка, переходим к следующему полю
				is.formManager.NextField()
				is.formManager.UpdateFocus()
			}

		default:
			// Обновляем текущее поле
			if currentField := is.formManager.GetCurrentFieldModel(); currentField != nil {
				if !currentField.IsButton() {
					_, fieldCmd := currentField.Update(msg)
					if fieldCmd != nil {
						if teaCmd, ok := fieldCmd.(tea.Cmd); ok {
							cmd = teaCmd
						}
					}
				}
			}
		}

	case ImportResultMsg:
		// Обрабатываем результат импорта
		switch msg.Type {
		case "error":
			is.messageManager.AddError(msg.Message)
		case "success":
			is.messageManager.AddSuccess(msg.Message)
			// Добавляем дополнительные детали
			for _, detail := range msg.Details {
				is.messageManager.AddInfo(detail)
			}
		case "warning":
			is.messageManager.AddWarning(msg.Message)
			// Добавляем дополнительные детали
			for _, detail := range msg.Details {
				is.messageManager.AddInfo(detail)
			}
		}
	}

	// Обновляем менеджер сообщений (MessageManager не имеет метода Update)

	return is, cmd
}

// View возвращает строку для отрисовки
func (is *ImportScreen) View() string {
	is.updateContent()
	return is.BaseScreen.View()
}

// updateContent обновляет содержимое экрана
func (is *ImportScreen) updateContent() {
	// Создаем стили
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Margin(0, 0, 1, 0)

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Margin(0, 0, 2, 0)

	// Создаем заголовок
	header := headerStyle.Render("Импорт конфигурации SSH")

	// Создаем описание
	description := descriptionStyle.Render("Импорт загрузит подключения из файла конфигурации SSH. Пароли будут зашифрованы мастер-паролем.")

	// Рендерим форму
	formContent := is.formManager.RenderForm()

	// Рендерим сообщения
	messages := is.messageManager.RenderMessages(is.width - 20)

	// Объединяем все части
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		description,
		formContent,
		"",
		messages,
	)

	is.SetContent(content)
}

// performImport выполняет импорт конфигурации
func (is *ImportScreen) performImport() tea.Cmd {
	return func() tea.Msg {
		// Получаем путь из формы и убираем лишние пробелы
		importPath := strings.TrimSpace(is.formManager.GetField("import_path").Value())

		if importPath == "" {
			return ImportResultMsg{
				Success: false,
				Message: "Укажите путь к файлу для импорта",
				Type:    "error",
			}
		}

		// Проверяем, что файл существует
		if _, err := os.Stat(importPath); os.IsNotExist(err) {
			return ImportResultMsg{
				Success: false,
				Message: fmt.Sprintf("Файл %s не существует", importPath),
				Type:    "error",
			}
		}

		// Получаем количество подключений до импорта
		connectionsBefore := len(is.connectionService.GetAllConnections())

		// Выполняем импорт с шифрованием паролей
		err := is.connectionService.ImportConfigPlain(importPath)
		if err != nil {
			// Проверяем, является ли это ошибкой о дубликатах
			if strings.Contains(err.Error(), "all") && strings.Contains(err.Error(), "connections already exist") {
				return ImportResultMsg{
					Success: true,
					Message: fmt.Sprintf("%v", err),
					Type:    "warning",
					Details: []string{"All connections from file already exist in database"},
				}
			} else {
				return ImportResultMsg{
					Success: false,
					Message: fmt.Sprintf("Ошибка импорта: %v", err),
					Type:    "error",
				}
			}
		}

		// Получаем количество подключений после импорта
		connectionsAfter := len(is.connectionService.GetAllConnections())
		importedCount := connectionsAfter - connectionsBefore

		return ImportResultMsg{
			Success: true,
			Message: fmt.Sprintf("Конфигурация успешно импортирована из %s", importPath),
			Type:    "success",
			Details: []string{
				fmt.Sprintf("Импортировано %d новых подключений", importedCount),
				"All passwords encrypted with master password",
			},
		}
	}
}

// ImportResultMsg сообщение с результатом импорта
type ImportResultMsg struct {
	Success bool
	Message string
	Type    string
	Details []string
}

// Init инициализирует экран
func (is *ImportScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (is *ImportScreen) GetName() string {
	return "import"
}
