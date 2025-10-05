package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ExportScreen представляет экран экспорта конфигурации
type ExportScreen struct {
	*BaseScreen
	formManager       *components.FormManager
	messageManager    *components.MessageManager
	connectionService *services.ConnectionService
	exportPath        string
	isExporting       bool
}

// NewExportScreen создает новый экран экспорта
func NewExportScreen() *ExportScreen {
	baseScreen := NewBaseScreen("Экспорт конфигурации")

	// Создаем менеджер формы
	formManager := components.NewFormManager()

	// Добавляем поля формы
	formManager.AddField(components.FieldConfig{
		Name:        "export_path",
		Label:       "Путь к файлу",
		Required:    true,
		Width:       50,
		Placeholder: "/path/to/exported_config",
		FieldType:   components.FieldTypeText,
	})

	formManager.AddField(components.FieldConfig{
		Name:      "export_button",
		Label:     "Экспорт",
		FieldType: components.FieldTypeButton,
		Style:     "success",
	})

	// Устанавливаем фокус на первое поле
	formManager.SetCurrentField("export_path")
	formManager.UpdateFocus()

	// Создаем менеджер сообщений
	messageManager := components.NewMessageManager()

	// Получаем сервис подключений
	connectionService := services.GetGlobalConnectionService()

	return &ExportScreen{
		BaseScreen:        baseScreen,
		formManager:       formManager,
		messageManager:    messageManager,
		connectionService: connectionService,
		exportPath:        "",
		isExporting:       false,
	}
}

// Update обрабатывает обновления состояния
func (es *ExportScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		es.SetSize(msg.Width, msg.Height)
		return es, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return es, ui.GoBackCmd()

		case "tab":
			es.formManager.NextField()
			es.formManager.UpdateFocus()

		case "shift+tab":
			es.formManager.PrevField()
			es.formManager.UpdateFocus()

		case "enter":
			currentField := es.formManager.GetCurrentField()
			if currentField == "export_button" {
				return es, es.performExport()
			} else {
				// Если это не кнопка, переходим к следующему полю
				es.formManager.NextField()
				es.formManager.UpdateFocus()
			}

		default:
			// Обновляем текущее поле
			if currentField := es.formManager.GetCurrentFieldModel(); currentField != nil {
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

	case ExportResultMsg:
		// Обрабатываем результат экспорта
		switch msg.Type {
		case "error":
			es.messageManager.AddError(msg.Message)
		case "success":
			es.messageManager.AddSuccess(msg.Message)
			// Добавляем дополнительные детали
			for _, detail := range msg.Details {
				es.messageManager.AddInfo(detail)
			}
		}
	}

	// Обновляем менеджер сообщений (MessageManager не имеет метода Update)

	return es, cmd
}

// View возвращает строку для отрисовки
func (es *ExportScreen) View() string {
	es.updateContent()
	return es.BaseScreen.View()
}

// updateContent обновляет содержимое экрана
func (es *ExportScreen) updateContent() {
	// Создаем стили
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Margin(0, 0, 1, 0)

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Margin(0, 0, 2, 0)

	// Создаем заголовок
	header := headerStyle.Render("Экспорт конфигурации SSH")

	// Создаем описание
	description := descriptionStyle.Render("Экспорт создаст файл с конфигурацией SSH в том же формате, но без шифрования паролей.")

	// Рендерим форму
	formContent := es.formManager.RenderForm()

	// Рендерим сообщения
	messages := es.messageManager.RenderMessages(es.width - 20)

	// Объединяем все части
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		description,
		formContent,
		"",
		messages,
	)

	es.SetContent(content)
}

// performExport выполняет экспорт конфигурации
func (es *ExportScreen) performExport() tea.Cmd {
	return func() tea.Msg {
		// Получаем путь из формы и убираем лишние пробелы
		exportPath := strings.TrimSpace(es.formManager.GetField("export_path").Value())

		if exportPath == "" {
			return ExportResultMsg{
				Success: false,
				Message: "Укажите путь к файлу для экспорта",
				Type:    "error",
			}
		}

		// Проверяем, что директория существует
		dir := filepath.Dir(exportPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return ExportResultMsg{
				Success: false,
				Message: fmt.Sprintf("Директория %s не существует", dir),
				Type:    "error",
			}
		}

		// Выполняем экспорт без шифрования паролей
		err := es.connectionService.ExportConfigPlain(exportPath)
		if err != nil {
			return ExportResultMsg{
				Success: false,
				Message: fmt.Sprintf("Ошибка экспорта: %v", err),
				Type:    "error",
			}
		}

		// Получаем количество подключений для отчета
		connections := es.connectionService.GetAllConnections()
		return ExportResultMsg{
			Success: true,
			Message: fmt.Sprintf("Конфигурация успешно экспортирована в %s", exportPath),
			Type:    "success",
			Details: []string{
				fmt.Sprintf("Экспортировано %d подключений", len(connections)),
				"Пароли сохранены в открытом виде - защитите файл!",
			},
		}
	}
}

// ExportResultMsg сообщение с результатом экспорта
type ExportResultMsg struct {
	Success bool
	Message string
	Type    string
	Details []string
}

// Init инициализирует экран
func (es *ExportScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (es *ExportScreen) GetName() string {
	return "export"
}
