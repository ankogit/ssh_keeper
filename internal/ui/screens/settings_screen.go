package screens

import (
	"fmt"
	"strings"

	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// addMessageMsg сообщение для добавления сообщения в MessageManager
type addMessageMsg struct {
	messageType string // "error", "success", "info", "warning"
	text        string
}

// updateListMsg сообщение для обновления списка
type updateListMsg struct{}

// SettingsScreen представляет экран настроек
type SettingsScreen struct {
	*BaseScreen
	list           list.Model
	menuItems      []ui.MenuItem
	config         ui.MenuConfig
	messageManager *components.MessageManager
}

// getPasswordStartupDescription возвращает динамическое описание для настройки пароля
func getPasswordStartupDescription() string {
	currentSetting, err := services.GetRequirePasswordOnStartupWithSignature()
	if err != nil {
		return "Переключить требование ввода мастер-пароля при каждом запуске"
	}

	if currentSetting {
		return "Сейчас запрашиваем пароль при запуске (нажмите для отключения)"
	} else {
		return "Сейчас НЕ запрашиваем пароль при запуске (нажмите для включения)"
	}
}

// NewSettingsScreen создает новый экран настроек
func NewSettingsScreen() *SettingsScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Настройки")
	messageManager := components.NewMessageManager()

	// Создаем конфигурацию меню с действиями
	config := ui.MenuConfig{
		Title: "Настройки",
		Items: []ui.MenuItemConfig{
			{
				Title:       "Запрашивать пароль при запуске",
				Description: getPasswordStartupDescription(),
				Shortcut:    "1",
				Action: func() tea.Cmd {
					// Получаем текущую настройку
					currentSetting, err := services.GetRequirePasswordOnStartupWithSignature()
					if err != nil {
						return func() tea.Msg {
							return addMessageMsg{
								messageType: "error",
								text:        fmt.Sprintf("Ошибка получения настройки: %v", err),
							}
						}
					}

					// Переключаем настройку
					newSetting := !currentSetting
					err = services.SetRequirePasswordOnStartupWithSignature(newSetting)
					if err != nil {
						return func() tea.Msg {
							return addMessageMsg{
								messageType: "error",
								text:        fmt.Sprintf("Ошибка изменения настройки: %v", err),
							}
						}
					}

					if newSetting {
						return tea.Sequence(
							func() tea.Msg {
								return updateListMsg{}
							},
						)
					} else {
						return tea.Sequence(
							func() tea.Msg {
								return updateListMsg{}
							},
						)
					}
				},
			},
			{
				Title:       "Сбросить мастер-пароль",
				Description: "Удалить мастер-пароль и выйти из приложения",
				Shortcut:    "2",
				Action: func() tea.Cmd {
					// Сбрасываем мастер-пароль с проверкой подписи
					err := services.ClearMasterPasswordWithSignature()
					if err != nil {
						return func() tea.Msg {
							return addMessageMsg{
								messageType: "error",
								text:        fmt.Sprintf("Ошибка сброса мастер-пароля: %v", err),
							}
						}
					}
					return func() tea.Msg {
						return addMessageMsg{
							messageType: "success",
							text:        "✅ Мастер-пароль сброшен. Приложение будет закрыто.",
						}
					}
				},
			},
			{
				Title:       "Экспорт подключений",
				Description: "Экспортировать все подключения в файл",
				Shortcut:    "3",
				Action: func() tea.Cmd {
					return ui.NavigateToCmd("export")
				},
			},
			{
				Title:       "Импорт подключений",
				Description: "Импортировать подключения из файла",
				Shortcut:    "4",
				Action: func() tea.Cmd {
					return ui.NavigateToCmd("import")
				},
			},
			{
				Title:       "Назад",
				Description: "Вернуться к главному меню",
				Shortcut:    "esc",
				Action: func() tea.Cmd {
					return ui.GoBackCmd()
				},
			},
		},
		ShowBack: true,
		ShowQuit: true,
	}

	// Создаем элементы меню с действиями
	menuItems := make([]ui.MenuItem, len(config.Items))
	for i, itemConfig := range config.Items {
		menuItems[i] = ui.NewMenuItem(itemConfig)
	}

	// Создаем список
	l := list.New(convertToListItem(menuItems), list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	// Отключаем обработку клавиши 'q' в списке
	l.KeyMap.Quit.SetKeys("ctrl+q")

	l.Styles.PaginationStyle = lipgloss.NewStyle().
		Margin(styles.ListPaginationMargin, 0, 0, 0)

	return &SettingsScreen{
		BaseScreen:     baseScreen,
		list:           l,
		config:         config,
		menuItems:      menuItems,
		messageManager: messageManager,
	}
}

// updateMenuItems обновляет элементы меню с новыми описаниями
func (ss *SettingsScreen) updateMenuItems() {
	// Обновляем описание для настройки пароля
	ss.config.Items[0].Description = getPasswordStartupDescription()

	// Пересоздаем элементы меню
	ss.menuItems = make([]ui.MenuItem, len(ss.config.Items))
	for i, itemConfig := range ss.config.Items {
		ss.menuItems[i] = ui.NewMenuItem(itemConfig)
	}

	// Обновляем список
	items := make([]list.Item, len(ss.menuItems))
	for i, item := range ss.menuItems {
		items[i] = item
	}
	ss.list.SetItems(items)
}

// Update обрабатывает обновления состояния
func (ss *SettingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ss.SetSize(msg.Width, msg.Height)
		// Учитываем все служебные элементы: заголовки, рамки, отступы
		ss.list.SetSize(msg.Width-4, msg.Height-12)
		return ss, nil

	case updateListMsg:
		// Обновляем список с новыми описаниями
		ss.updateMenuItems()
		return ss, nil

	case addMessageMsg:
		// Добавляем сообщение в MessageManager
		switch msg.messageType {
		case "error":
			ss.messageManager.AddError(msg.text)
		case "success":
			ss.messageManager.AddSuccess(msg.text)
		case "info":
			ss.messageManager.AddInfo(msg.text)
		case "warning":
			ss.messageManager.AddWarning(msg.text)
		}

		// Если это сообщение об успешном сбросе пароля, выходим из приложения
		if msg.messageType == "success" && strings.Contains(msg.text, "Мастер-пароль сброшен") {
			return ss, tea.Quit
		}

		return ss, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return ss, tea.Quit
		case "esc":
			// Возвращаемся к главному меню
			return ss, func() tea.Msg {
				return ui.NavigateToMsg{ScreenName: "main_menu"}
			}
		case "enter":
			// Обрабатываем выбор элемента
			selectedItem := ss.list.SelectedItem()
			if item, ok := selectedItem.(ui.MenuItem); ok {
				// Выполняем действие элемента меню
				actionCmd := item.Execute()
				if actionCmd != nil {
					cmd = actionCmd
				}
			}
		default:
			// Проверяем горячие клавиши
			for _, menuItem := range ss.menuItems {
				if menuItem.GetShortcut() == msg.String() {
					// Выполняем действие элемента меню
					actionCmd := menuItem.Execute()
					if actionCmd != nil {
						cmd = actionCmd
					}
					break
				}
			}
		}
	}

	// Обновляем список
	var listCmd tea.Cmd
	ss.list, listCmd = ss.list.Update(msg)
	if listCmd != nil {
		cmd = listCmd
	}

	// Обновляем базовый экран
	baseScreen, baseCmd := ss.BaseScreen.Update(msg)
	ss.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return ss, cmd
}

// View возвращает строку для отрисовки
func (ss *SettingsScreen) View() string {
	ss.updateContent()
	return ss.BaseScreen.View()
}

// updateContent обновляет содержимое экрана
func (ss *SettingsScreen) updateContent() {
	// Рендерим список
	listContent := ss.list.View()

	// Создаем стили
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorPrimary)).
		Bold(true).
		Margin(0, 0, 1, 0)

	// Создаем заголовок
	header := headerStyle.Render("Выберите действие:")

	// Рендерим сообщения
	messagesContent := ss.messageManager.RenderMessages(80)

	// Объединяем все элементы
	var content string
	if messagesContent != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, header, messagesContent, listContent)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left, header, listContent)
	}

	ss.SetContent(content)
}

// Init инициализирует экран
func (ss *SettingsScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (ss *SettingsScreen) GetName() string {
	return "settings"
}
