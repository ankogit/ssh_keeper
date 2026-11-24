package screens

import (
	"fmt"
	"os"
	"ssh-keeper/internal/models"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ssh"
	"ssh-keeper/internal/ui"
	"ssh-keeper/internal/ui/components"
	"ssh-keeper/internal/ui/styles"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type promptKind int

const (
	promptNone promptKind = iota
	promptDelete
	promptTrust
	promptReplace
)

type listPrompt struct {
	kind  promptKind
	conn  models.Connection
	check ssh.CheckResult
}

type hostKeyCheckedMsg struct {
	id     int
	conn   models.Connection
	result ssh.CheckResult
}

type sessionFinishedMsg struct {
	name     string
	passFile string
	err      error
}

// ConnectionsScreen представляет экран управления подключениями.
type ConnectionsScreen struct {
	*BaseScreen
	list           list.Model
	searchInput    textinput.Model
	searchActive   bool
	allItems       []list.Item
	messageManager *components.MessageManager
	prompt         listPrompt
	checking       bool
	checkID        int
	width          int
	height         int
}

// NewConnectionsScreen создает новый экран подключений.
func NewConnectionsScreen() *ConnectionsScreen {
	baseScreen := NewBaseScreen("SSH Keeper")

	var listItems []list.Item
	for _, conn := range services.GetConnections() {
		listItems = append(listItems, components.NewConnectionItem(conn))
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.KeyMap.Quit.SetKeys("ctrl+q")
	l.Styles.PaginationStyle = lipgloss.NewStyle().Margin(1, 0, 0, 0)

	searchInput := textinput.New()
	searchInput.Placeholder = "Поиск..."
	searchInput.CharLimit = 25
	searchInput.Width = 40
	searchInput.Blur()

	return &ConnectionsScreen{
		BaseScreen:     baseScreen,
		list:           l,
		searchInput:    searchInput,
		allItems:       listItems,
		messageManager: components.NewMessageManager(),
	}
}

func (cs *ConnectionsScreen) refreshConnections() {
	if err := services.ReloadConnections(); err != nil {
		cs.show(components.MessageTypeError, fmt.Sprintf("Ошибка загрузки подключений: %v", err))
		return
	}
	var listItems []list.Item
	for _, conn := range services.GetConnections() {
		listItems = append(listItems, components.NewConnectionItem(conn))
	}
	cs.allItems = listItems
	cs.filterList()
}

func (cs *ConnectionsScreen) isAddSelected() bool {
	_, ok := cs.list.SelectedItem().(components.AddConnectionItem)
	return ok
}

func (cs *ConnectionsScreen) selected() (models.Connection, bool) {
	item, ok := cs.list.SelectedItem().(components.ConnectionItem)
	if !ok {
		return models.Connection{}, false
	}
	return item.GetConnection(), true
}

func (cs *ConnectionsScreen) editSelectedConnection() tea.Cmd {
	conn, ok := cs.selected()
	if !ok {
		cs.show(components.MessageTypeError, "Нет выбранного подключения")
		return nil
	}
	return ui.NavigateToWithDataCmd("edit_connection", conn)
}

func (cs *ConnectionsScreen) askDelete() {
	conn, ok := cs.selected()
	if !ok {
		cs.show(components.MessageTypeError, "Нет выбранного подключения")
		return
	}
	cs.prompt = listPrompt{kind: promptDelete, conn: conn}
}

func (cs *ConnectionsScreen) beginConnect() tea.Cmd {
	conn, ok := cs.selected()
	if !ok {
		cs.show(components.MessageTypeError, "Нет выбранного подключения")
		return nil
	}
	cs.checking = true
	cs.checkID++
	id := cs.checkID
	cs.show(components.MessageTypeInfo, fmt.Sprintf("Проверяем ключ %s...", conn.Host))
	return func() tea.Msg {
		return hostKeyCheckedMsg{id: id, conn: conn, result: ssh.Check(conn.Host, conn.Port)}
	}
}

func (cs *ConnectionsScreen) handleHostKey(msg hostKeyCheckedMsg) tea.Cmd {
	if msg.id != cs.checkID {
		return nil
	}
	cs.checking = false
	switch msg.result.Status {
	case ssh.StatusTrusted:
		return cs.startSession(msg.conn)
	case ssh.StatusUnknown:
		cs.messageManager.ClearMessages()
		cs.prompt = listPrompt{kind: promptTrust, conn: msg.conn, check: msg.result}
		return nil
	case ssh.StatusChanged:
		cs.messageManager.ClearMessages()
		cs.prompt = listPrompt{kind: promptReplace, conn: msg.conn, check: msg.result}
		return nil
	default:
		text := "сервер недоступен"
		if msg.result.Err != nil {
			text = msg.result.Err.Error()
		}
		cs.show(components.MessageTypeError, "Не удалось проверить ключ: "+text)
		return nil
	}
}

func (cs *ConnectionsScreen) acceptPrompt() tea.Cmd {
	current := cs.prompt
	cs.prompt = listPrompt{}
	switch current.kind {
	case promptDelete:
		if err := services.DeleteConnection(current.conn.ID); err != nil {
			cs.show(components.MessageTypeError, fmt.Sprintf("Ошибка удаления: %v", err))
			return nil
		}
		cs.refreshConnections()
		cs.show(components.MessageTypeSuccess, fmt.Sprintf("Подключение «%s» удалено", current.conn.Name))
		return nil
	case promptTrust, promptReplace:
		if current.check.Presented == nil {
			cs.show(components.MessageTypeError, "Нет ключа сервера для записи")
			return nil
		}
		path, err := ssh.KnownHostsPath()
		if err != nil {
			cs.show(components.MessageTypeError, err.Error())
			return nil
		}
		port := current.conn.Port
		if port == 0 {
			port = 22
		}
		if err := ssh.SaveKey(path, current.conn.Host, port, current.check.Presented); err != nil {
			cs.show(components.MessageTypeError, fmt.Sprintf("Не удалось записать known_hosts: %v", err))
			return nil
		}
		return cs.startSession(current.conn)
	default:
		return nil
	}
}

func (cs *ConnectionsScreen) startSession(conn models.Connection) tea.Cmd {
	cmd, passFile, err := ssh.Prepare(&conn)
	if err != nil {
		if passFile != "" {
			os.Remove(passFile)
		}
		cs.show(components.MessageTypeError, err.Error())
		return nil
	}
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return sessionFinishedMsg{name: conn.Name, passFile: passFile, err: err}
	})
}

// Update обрабатывает события списка.
func (cs *ConnectionsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cs.SetSize(msg.Width, msg.Height)
		cs.width = msg.Width
		cs.height = msg.Height
		cs.resizeList()
		return cs, nil
	case ui.NavigateToMsg:
		if msg.ScreenName == "connections" {
			cs.prompt = listPrompt{}
			cs.checking = false
			cs.refreshConnections()
		}
		return cs, nil
	case hostKeyCheckedMsg:
		return cs, cs.handleHostKey(msg)
	case sessionFinishedMsg:
		if msg.passFile != "" {
			os.Remove(msg.passFile)
		}
		if msg.err != nil {
			cs.show(components.MessageTypeError, fmt.Sprintf("Сессия «%s» завершилась с ошибкой: %v", displayName(msg.name), msg.err))
		} else {
			cs.show(components.MessageTypeSuccess, fmt.Sprintf("Сессия «%s» завершена", displayName(msg.name)))
		}
		return cs, nil
	case tea.KeyMsg:
		return cs.handleKey(msg)
	}

	if cs.prompt.kind == promptNone && !cs.searchActive && !cs.checking {
		var cmd tea.Cmd
		cs.list, cmd = cs.list.Update(msg)
		return cs, cmd
	}
	return cs, nil
}

func (cs *ConnectionsScreen) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if cs.prompt.kind != promptNone {
		switch msg.String() {
		case "ctrl+c":
			return cs, tea.Quit
		case "y", "Y", "enter":
			return cs, cs.acceptPrompt()
		case "n", "N", "esc":
			cs.prompt = listPrompt{}
			return cs, nil
		}
		return cs, nil
	}

	if cs.checking {
		switch msg.String() {
		case "ctrl+c":
			return cs, tea.Quit
		case "esc":
			cs.checkID++
			cs.checking = false
			cs.show(components.MessageTypeInfo, "Проверка отменена")
			return cs, nil
		}
		return cs, nil
	}

	if cs.searchActive {
		switch msg.String() {
		case "ctrl+c":
			return cs, tea.Quit
		case "esc":
			cs.searchInput.SetValue("")
			cs.searchInput.Blur()
			cs.searchActive = false
			cs.filterList()
			return cs, nil
		case "enter":
			cs.searchInput.Blur()
			cs.searchActive = false
			return cs, nil
		}
		var cmd tea.Cmd
		cs.searchInput, cmd = cs.searchInput.Update(msg)
		cs.filterList()
		return cs, cmd
	}

	switch msg.String() {
	case "ctrl+c", "esc":
		return cs, tea.Quit
	case "enter":
		if cs.isAddSelected() {
			return cs, ui.NavigateToCmd("add_connection")
		}
		return cs, cs.beginConnect()
	case "a":
		return cs, ui.NavigateToCmd("add_connection")
	case "e":
		if cs.isAddSelected() {
			return cs, ui.NavigateToCmd("add_connection")
		}
		return cs, cs.editSelectedConnection()
	case "d":
		if cs.isAddSelected() {
			cs.show(components.MessageTypeInfo, "Это кнопка добавления, а не сервер")
			return cs, nil
		}
		cs.askDelete()
		return cs, nil
	case "`":
		return cs, ui.NavigateToCmd("settings")
	case "/":
		cs.searchActive = true
		cs.searchInput.Focus()
		return cs, textinput.Blink
	}

	var cmd tea.Cmd
	cs.list, cmd = cs.list.Update(msg)
	return cs, cmd
}

// View возвращает строку для отрисовки.
func (cs *ConnectionsScreen) View() string {
	cs.updateContent()
	return cs.BaseScreen.View()
}

func (cs *ConnectionsScreen) updateContent() {
	cs.resizeList()
	instructionsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted)).
		Italic(styles.TextItalic)

	parts := []string{cs.list.View()}
	if cs.searchActive || cs.searchInput.Value() != "" {
		border := styles.ColorMuted
		if cs.searchActive {
			border = styles.ColorWarning
		}
		searchView := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(border)).
			Padding(0, 1).
			Render(cs.searchInput.View())
		parts = []string{searchView, "", cs.list.View()}
	}

	if promptText := cs.promptText(); promptText != "" {
		parts = append(parts, lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorWarning)).
			Render(promptText))
	}

	parts = append(parts, instructionsStyle.Render(cs.instructions()))
	if messages := cs.messageManager.RenderMessages(cs.messageWidth()); messages != "" {
		parts = append(parts, messages)
	}
	cs.SetContent(lipgloss.JoinVertical(lipgloss.Left, parts...))
}

func (cs *ConnectionsScreen) instructions() string {
	switch {
	case cs.prompt.kind != promptNone:
		return "y — подтвердить • n или Esc — отмена"
	case cs.checking:
		return "Проверяем ключ сервера... • Esc — отмена"
	case cs.searchActive:
		return "Введите запрос • Enter — закрыть поиск • Esc — сбросить"
	default:
		return "/ поиск • ↑/↓ выбор • Enter открыть • e править • d удалить • ` настройки • Esc выход"
	}
}

func (cs *ConnectionsScreen) promptText() string {
	switch cs.prompt.kind {
	case promptDelete:
		return fmt.Sprintf("Удалить «%s»? y — да, n — нет", displayName(cs.prompt.conn.Name))
	case promptTrust:
		return fmt.Sprintf("Новый ключ %s (%s). Доверять и подключиться? y — да, n — нет",
			cs.prompt.check.PresentedFingerprint, endpoint(cs.prompt.conn))
	case promptReplace:
		return fmt.Sprintf("Ключ %s изменился. Было %s, стало %s. Перезаписать known_hosts и подключиться? y — да, n — нет",
			endpoint(cs.prompt.conn), cs.prompt.check.StoredFingerprint, cs.prompt.check.PresentedFingerprint)
	default:
		return ""
	}
}

func (cs *ConnectionsScreen) resizeList() {
	height := cs.height - 12
	if cs.searchActive || cs.searchInput.Value() != "" {
		height -= 3
	}
	if cs.prompt.kind != promptNone {
		height -= 3
	}
	if height < 5 {
		height = 5
	}
	width := cs.width - 4
	if width < 10 {
		width = 10
	}
	cs.list.SetSize(width, height)
}

func (cs *ConnectionsScreen) messageWidth() int {
	if cs.width > 20 {
		return cs.width - 10
	}
	return 70
}

// Init перезагружает подключения.
func (cs *ConnectionsScreen) Init() tea.Cmd {
	cs.refreshConnections()
	return nil
}

// GetName возвращает имя экрана.
func (cs *ConnectionsScreen) GetName() string {
	return "connections"
}

func (cs *ConnectionsScreen) filterList() {
	query := strings.ToLower(cs.searchInput.Value())
	if query == "" {
		cs.list.SetItems(append(append([]list.Item{}, cs.allItems...), components.AddConnectionItem{}))
		return
	}
	var filtered []list.Item
	for _, item := range cs.allItems {
		connItem, ok := item.(components.ConnectionItem)
		if !ok {
			continue
		}
		haystack := strings.ToLower(connItem.Title() + " " + connItem.Description() + " " + connItem.FilterValue())
		if strings.Contains(haystack, query) {
			filtered = append(filtered, item)
		}
	}
	filtered = append(filtered, components.AddConnectionItem{})
	cs.list.SetItems(filtered)
}

func (cs *ConnectionsScreen) show(kind components.MessageType, text string) {
	cs.messageManager.ClearMessages()
	cs.messageManager.AddMessage(components.Message{
		Type:      kind,
		Text:      text,
		Timestamp: time.Now(),
		Duration:  8 * time.Second,
	})
}

func displayName(name string) string {
	if name == "" {
		return "подключение"
	}
	return name
}

func endpoint(conn models.Connection) string {
	port := conn.Port
	if port == 0 {
		port = 22
	}
	return fmt.Sprintf("%s:%d", conn.Host, port)
}
