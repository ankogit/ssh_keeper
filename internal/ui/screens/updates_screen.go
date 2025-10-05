package screens

import (
	"fmt"
	"os"
	"time"

	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// updateCheckMsg сообщение о результате проверки обновлений
type updateCheckMsg struct {
	updateInfo *services.UpdateInfo
	err        error
}

// updateDownloadMsg сообщение о результате загрузки обновления
type updateDownloadMsg struct {
	err error
}

// UpdatesScreen представляет экран управления обновлениями
type UpdatesScreen struct {
	*BaseScreen
	list           list.Model
	menuItems      []ui.MenuItem
	config         ui.MenuConfig
	messageManager *components.MessageManager
	updateService  *services.UpdateService
	updateInfo     *services.UpdateInfo
	isChecking     bool
	isDownloading  bool
}

// NewUpdatesScreen создает новый экран обновлений
func NewUpdatesScreen() *UpdatesScreen {
	baseScreen := NewBaseScreen("SSH Keeper - Обновления")
	messageManager := components.NewMessageManager()

	// Получаем глобальный сервис автоматических обновлений
	autoUpdateService := services.GetGlobalAutoUpdateService()
	var updateService *services.UpdateService
	if autoUpdateService != nil {
		updateService = autoUpdateService.GetUpdateService()
	} else {
		// Fallback если сервис не инициализирован
		version := os.Getenv("APP_VERSION")
		if version == "" {
			version = "0.1.0"
		}
		updateService = services.NewUpdateService(version)
	}

	// Создаем конфигурацию меню с действиями
	config := ui.MenuConfig{
		Title: "Управление обновлениями",
		Items: []ui.MenuItemConfig{
			{
				Title:       "Проверить обновления",
				Description: "Проверить наличие новых версий",
				Shortcut:    "1",
				Action: func() tea.Cmd {
					return func() tea.Msg {
						updateInfo, err := updateService.CheckForUpdates()
						return updateCheckMsg{updateInfo: updateInfo, err: err}
					}
				},
			},
			{
				Title:       "Загрузить и установить обновление",
				Description: "Загрузить и установить доступное обновление",
				Shortcut:    "2",
				Action: func() tea.Cmd {
					// Этот метод будет переопределен в updateMenuItems
					return nil
				},
			},
			{
				Title:       "Настройки обновлений",
				Description: "Настроить автоматическую проверку обновлений",
				Shortcut:    "3",
				Action: func() tea.Cmd {
					// TODO: Реализовать экран настроек обновлений
					return components.AddMessageCmd(components.NewMessage(components.MessageTypeInfo, "Настройки обновлений будут добавлены в следующей версии"))
				},
			},
			{
				Title:       "Назад",
				Description: "Вернуться к настройкам",
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

	return &UpdatesScreen{
		BaseScreen:     baseScreen,
		list:           l,
		config:         config,
		menuItems:      menuItems,
		messageManager: messageManager,
		updateService:  updateService,
	}
}

// Update обрабатывает обновления состояния
func (us *UpdatesScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		us.SetSize(msg.Width, msg.Height)
		// Учитываем все служебные элементы: заголовки, рамки, отступы
		us.list.SetSize(msg.Width-4, msg.Height-12)
		return us, nil

	case updateCheckMsg:
		us.isChecking = false
		if msg.err != nil {
			us.messageManager.AddError(fmt.Sprintf("Ошибка проверки обновлений: %v", msg.err))
		} else {
			us.updateInfo = msg.updateInfo
			if msg.updateInfo.IsAvailable {
				us.messageManager.AddSuccess(fmt.Sprintf("Доступно обновление до версии %s", msg.updateInfo.Version))
				us.updateMenuItems()
			} else {
				us.messageManager.AddInfo("У вас установлена последняя версия")
			}
		}
		return us, nil

	case updateDownloadMsg:
		us.isDownloading = false
		if msg.err != nil {
			us.messageManager.AddError(fmt.Sprintf("Ошибка установки обновления: %v", msg.err))
		} else {
			us.messageManager.AddSuccess("Обновление успешно установлено! Приложение будет перезапущено.")
			// Перезапускаем приложение через несколько секунд
			return us, tea.Sequence(
				func() tea.Msg {
					time.Sleep(2 * time.Second)
					return tea.Quit
				},
			)
		}
		return us, nil

	case components.MessageCmd:
		// Добавляем сообщение в MessageManager
		us.messageManager.AddMessage(msg.Message)
		return us, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return us, tea.Quit
		case "esc":
			// Возвращаемся к настройкам
			return us, func() tea.Msg {
				return ui.NavigateToMsg{ScreenName: "settings"}
			}
		case "enter":
			// Обрабатываем выбор элемента
			selectedItem := us.list.SelectedItem()
			if item, ok := selectedItem.(ui.MenuItem); ok {
				// Выполняем действие элемента меню
				actionCmd := item.Execute()
				if actionCmd != nil {
					cmd = actionCmd
				}
			}
		default:
			// Проверяем горячие клавиши
			for _, menuItem := range us.menuItems {
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
	us.list, listCmd = us.list.Update(msg)
	if listCmd != nil {
		cmd = listCmd
	}

	// Обновляем базовый экран
	baseScreen, baseCmd := us.BaseScreen.Update(msg)
	us.BaseScreen = baseScreen.(*BaseScreen)
	if baseCmd != nil {
		cmd = baseCmd
	}

	return us, cmd
}

// updateMenuItems обновляет элементы меню в зависимости от состояния обновлений
func (us *UpdatesScreen) updateMenuItems() {
	// Обновляем описание и действие для загрузки обновления
	if us.updateInfo != nil && us.updateInfo.IsAvailable {
		us.config.Items[1].Description = fmt.Sprintf("Загрузить версию %s (%.1f MB)",
			us.updateInfo.Version, float64(us.updateInfo.Size)/(1024*1024))

		// Устанавливаем действие для загрузки обновления
		us.config.Items[1].Action = func() tea.Cmd {
			return func() tea.Msg {
				err := us.updateService.DownloadAndInstallUpdate(us.updateInfo)
				return updateDownloadMsg{err: err}
			}
		}
	} else {
		us.config.Items[1].Description = "Загрузить и установить доступное обновление"
		us.config.Items[1].Action = func() tea.Cmd {
			return components.AddMessageCmd(components.NewMessage(components.MessageTypeWarning, "Нет доступных обновлений для загрузки"))
		}
	}

	// Пересоздаем элементы меню
	us.menuItems = make([]ui.MenuItem, len(us.config.Items))
	for i, itemConfig := range us.config.Items {
		us.menuItems[i] = ui.NewMenuItem(itemConfig)
	}

	// Обновляем список
	items := make([]list.Item, len(us.menuItems))
	for i, item := range us.menuItems {
		items[i] = item
	}
	us.list.SetItems(items)
}

// View возвращает строку для отрисовки
func (us *UpdatesScreen) View() string {
	us.updateContent()
	return us.BaseScreen.View()
}

// updateContent обновляет содержимое экрана
func (us *UpdatesScreen) updateContent() {
	// Рендерим список
	listContent := us.list.View()

	// Создаем стили
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorPrimary)).
		Bold(true).
		Margin(0, 0, 1, 0)

	// Создаем заголовок
	header := headerStyle.Render("Выберите действие:")

	// Добавляем информацию об обновлениях
	var updateInfo string
	if us.updateInfo != nil && us.updateInfo.IsAvailable {
		updateInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorSuccess)).
			Bold(true).
			Margin(0, 0, 1, 0).
			Render(fmt.Sprintf("🔄 Доступно обновление до версии %s", us.updateInfo.Version))
	} else if us.updateInfo != nil && !us.updateInfo.IsAvailable {
		updateInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorSuccess)).
			Margin(0, 0, 1, 0).
			Render("✅ У вас установлена последняя версия")
	}

	// Рендерим сообщения
	messagesContent := us.messageManager.RenderMessages(80)

	// Объединяем все элементы
	var content string
	if updateInfo != "" {
		if messagesContent != "" {
			content = lipgloss.JoinVertical(lipgloss.Left, header, updateInfo, messagesContent, listContent)
		} else {
			content = lipgloss.JoinVertical(lipgloss.Left, header, updateInfo, listContent)
		}
	} else {
		if messagesContent != "" {
			content = lipgloss.JoinVertical(lipgloss.Left, header, messagesContent, listContent)
		} else {
			content = lipgloss.JoinVertical(lipgloss.Left, header, listContent)
		}
	}

	us.SetContent(content)
}

// Init инициализирует экран
func (us *UpdatesScreen) Init() tea.Cmd {
	return nil
}

// GetName возвращает имя экрана
func (us *UpdatesScreen) GetName() string {
	return "updates"
}
