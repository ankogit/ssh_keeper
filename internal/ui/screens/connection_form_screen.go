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
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type formMode int

const (
	formCreate formMode = iota
	formEdit
)

type formHostCheckMsg struct {
	id     int
	result ssh.CheckResult
}

// ConnectionFormScreen создаёт или редактирует подключение.
type ConnectionFormScreen struct {
	*BaseScreen
	mode           formMode
	connectionSvc  *services.ConnectionService
	connection     *models.Connection
	formManager    *components.FormManager
	messageManager *components.MessageManager
	viewport       viewport.Model
	ready          bool
	confirmDelete  bool
	checking       bool
	checkID        int
}

// NewAddConnectionScreen создает экран добавления подключения.
func NewAddConnectionScreen() *ConnectionFormScreen {
	return newConnectionForm(formCreate)
}

// NewEditConnectionScreenEmpty создает пустой экран редактирования для фабрики.
func NewEditConnectionScreenEmpty() *ConnectionFormScreen {
	return newConnectionForm(formEdit)
}

func newConnectionForm(mode formMode) *ConnectionFormScreen {
	title := "SSH Keeper - Добавить подключение"
	if mode == formEdit {
		title = "SSH Keeper - Редактировать подключение"
	}

	formManager := components.NewFormManager()
	addText := func(name, label, placeholder string, required bool, width, maxLen int, fieldType components.FieldType) {
		formManager.AddField(components.FieldConfig{
			Name:        name,
			Label:       label,
			Required:    required,
			Width:       width,
			MaxLength:   maxLen,
			Placeholder: placeholder,
			FieldType:   fieldType,
		})
	}

	addText(components.FieldNameName, "Название", "Название подключения", true, 50, 50, components.FieldTypeText)
	addText(components.FieldNameHost, "Хост", "IP-адрес или домен", true, 50, 100, components.FieldTypeText)
	addText(components.FieldNamePort, "Порт", "22", false, 15, 5, components.FieldTypePort)
	addText(components.FieldNameUser, "Пользователь", "Имя пользователя", true, 50, 50, components.FieldTypeText)
	formManager.AddField(components.FieldConfig{
		Name:      components.FieldNameAuth,
		Label:     "По паролю",
		Required:  true,
		Width:     20,
		FieldType: components.FieldTypeBool,
	})
	addText(components.FieldNamePassword, "Пароль", "Пусто — спросить при подключении", false, 50, 100, components.FieldTypePassword)
	addText(components.FieldNameKey, "SSH ключ", "Пусто — ключ по умолчанию или агент", false, 50, 200, components.FieldTypeText)

	formManager.AddField(components.FieldConfig{
		Name:      "save",
		Label:     saveButtonLabel(mode),
		FieldType: components.FieldTypeButton,
		Width:     24,
	})
	if mode == formCreate {
		formManager.AddField(components.FieldConfig{
			Name:      "test",
			Label:     "Проверить сервер",
			FieldType: components.FieldTypeButton,
			Width:     24,
		})
	} else {
		formManager.AddField(components.FieldConfig{
			Name:      "delete",
			Label:     "Удалить подключение",
			FieldType: components.FieldTypeButton,
			Width:     24,
			Style:     "warning",
		})
	}

	instructions := "Tab / Shift+Tab — поля • ←/→ или пробел — способ входа • Ctrl+S — сохранить • Ctrl+T — проверить • Ctrl+P — показать пароль • Esc — назад"
	if mode == formEdit {
		instructions = "Tab / Shift+Tab — поля • ←/→ или пробел — способ входа • Ctrl+S — сохранить • Ctrl+T — проверить • Ctrl+D — удалить • Esc — назад"
	}
	formManager.SetInstructions(instructions)

	screen := &ConnectionFormScreen{
		BaseScreen:     NewBaseScreen(title),
		mode:           mode,
		connectionSvc:  services.GetGlobalConnectionService(),
		formManager:    formManager,
		messageManager: components.NewMessageManager(),
		viewport:       viewport.New(0, 0),
	}
	screen.updateFieldVisibility()
	screen.formManager.SetCurrentField(components.FieldNameName)
	screen.formManager.UpdateFocus()
	return screen
}

func saveButtonLabel(mode formMode) string {
	if mode == formEdit {
		return "Сохранить изменения"
	}
	return "Сохранить"
}

// SetData заполняет форму данными подключения.
func (s *ConnectionFormScreen) SetData(data interface{}) {
	connection, ok := data.(models.Connection)
	if !ok {
		s.show(components.MessageTypeError, "Не удалось загрузить подключение")
		return
	}
	s.connection = &connection
	s.SetTitle(fmt.Sprintf("SSH Keeper - Редактировать «%s»", connection.Name))
	s.prefill(connection)
}

func (s *ConnectionFormScreen) prefill(conn models.Connection) {
	set := func(name, value string) {
		if field := s.formManager.GetField(name); field != nil {
			field.SetValue(value)
			field.SetErrorText("")
		}
	}
	set(components.FieldNameName, conn.Name)
	set(components.FieldNameHost, conn.Host)
	port := conn.Port
	if port == 0 {
		port = 22
	}
	set(components.FieldNamePort, strconv.Itoa(port))
	set(components.FieldNameUser, conn.User)
	if conn.UseSSHKey {
		set(components.FieldNameAuth, "false")
	} else {
		set(components.FieldNameAuth, "true")
	}
	set(components.FieldNameKey, conn.KeyPath)
	set(components.FieldNamePassword, conn.Password)
	s.updateFieldVisibility()
	s.formManager.UpdateFocus()
	if s.ready {
		s.updateViewportContent()
	}
}

// Update обрабатывает события формы.
func (s *ConnectionFormScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetSize(msg.Width, msg.Height)
		s.viewport.Width = msg.Width - 4
		s.viewport.Height = msg.Height - 12
		if s.viewport.Height < 5 {
			s.viewport.Height = 5
		}
		s.ready = true
		s.updateViewportContent()
		return s, nil
	case ui.NavigateToMsg:
		if s.mode == formCreate && msg.ScreenName == "add_connection" {
			s.clearForm()
		}
		return s, nil
	case formHostCheckMsg:
		if msg.id != s.checkID {
			return s, nil
		}
		s.checking = false
		s.showHostCheck(msg.result)
		s.updateViewportContent()
		return s, nil
	case tea.KeyMsg:
		return s.handleKey(msg)
	}
	return s, nil
}

func (s *ConnectionFormScreen) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if s.confirmDelete {
		switch msg.String() {
		case "y", "Y", "enter":
			s.confirmDelete = false
			return s, s.deleteConnection()
		case "n", "N", "esc", "ctrl+c":
			if msg.String() == "ctrl+c" {
				return s, tea.Quit
			}
			s.confirmDelete = false
			s.updateViewportContent()
			return s, nil
		}
		return s, nil
	}

	if s.checking {
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		case "esc":
			s.checkID++
			s.checking = false
			s.show(components.MessageTypeInfo, "Проверка отменена")
			s.updateViewportContent()
			return s, nil
		}
		return s, nil
	}

	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c":
		return s, tea.Quit
	case "esc":
		s.clearForm()
		return s, ui.GoBackCmd()
	case "tab":
		s.formManager.NextField()
	case "shift+tab":
		s.formManager.PrevField()
	case "enter":
		return s, s.activateCurrent()
	case "ctrl+s":
		return s, s.saveConnection()
	case "ctrl+t":
		return s, s.testConnection()
	case "ctrl+p":
		s.togglePasswordVisibility()
	case "ctrl+d":
		if s.mode == formEdit {
			s.confirmDelete = true
			s.updateViewportContent()
		}
		return s, nil
	case "up":
		s.viewport.ScrollUp(1)
		return s, nil
	case "down":
		s.viewport.ScrollDown(1)
		return s, nil
	case "pgup":
		s.viewport.PageUp()
		return s, nil
	case "pgdown":
		s.viewport.PageDown()
		return s, nil
	default:
		if field := s.formManager.GetCurrentFieldModel(); field != nil && !field.IsButton() {
			_, fieldCmd := field.Update(msg)
			if teaCmd, ok := fieldCmd.(tea.Cmd); ok {
				cmd = teaCmd
			}
		}
	}

	s.formManager.UpdateFocus()
	s.updateFieldVisibility()
	s.updateViewportContent()
	switch msg.String() {
	case "tab", "shift+tab", "enter":
		s.scrollToCurrentField()
	}
	return s, cmd
}

func (s *ConnectionFormScreen) activateCurrent() tea.Cmd {
	field := s.formManager.GetCurrentFieldModel()
	if field != nil && field.IsButton() {
		switch field.GetName() {
		case "save":
			return s.saveConnection()
		case "test":
			return s.testConnection()
		case "delete":
			s.confirmDelete = true
			s.updateViewportContent()
			return nil
		}
	}
	s.formManager.NextField()
	s.formManager.UpdateFocus()
	s.updateFieldVisibility()
	s.updateViewportContent()
	s.scrollToCurrentField()
	return nil
}

// View рисует форму.
func (s *ConnectionFormScreen) View() string {
	if !s.ready {
		return "Инициализация..."
	}
	content := s.viewport.View()
	if s.confirmDelete {
		name := "это подключение"
		if s.connection != nil && s.connection.Name != "" {
			name = fmt.Sprintf("«%s»", s.connection.Name)
		}
		content += "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorWarning)).
			Render(fmt.Sprintf("Удалить %s? y — да, n — нет", name))
	}
	if !s.viewport.AtBottom() {
		content += "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ColorMuted)).
			Render("↓ Ниже есть ещё поля")
	}
	if messages := s.messageManager.RenderMessages(80); messages != "" {
		content += "\n" + messages
	}
	s.SetContent(content)
	return s.BaseScreen.View()
}

func (s *ConnectionFormScreen) updateFieldVisibility() {
	authField := s.formManager.GetField(components.FieldNameAuth)
	if authField == nil {
		return
	}
	usePassword := authField.Value() == "true"
	if field := s.formManager.GetField(components.FieldNamePassword); field != nil {
		field.SetVisible(usePassword)
	}
	if field := s.formManager.GetField(components.FieldNameKey); field != nil {
		field.SetVisible(!usePassword)
	}
}

func (s *ConnectionFormScreen) updateViewportContent() {
	if !s.ready {
		return
	}
	s.viewport.SetContent(s.formManager.RenderForm())
}

func (s *ConnectionFormScreen) saveConnection() tea.Cmd {
	conn, ok := s.readForm()
	if !ok {
		return nil
	}
	if s.connectionSvc == nil {
		s.show(components.MessageTypeError, "Сервис подключений недоступен")
		s.updateViewportContent()
		return nil
	}

	var err error
	if s.mode == formCreate {
		err = s.connectionSvc.AddConnection(conn)
	} else {
		err = s.connectionSvc.UpdateConnection(conn)
	}
	if err != nil {
		s.show(components.MessageTypeError, fmt.Sprintf("Ошибка сохранения: %v", err))
		s.updateViewportContent()
		return nil
	}

	if s.mode == formCreate {
		s.show(components.MessageTypeSuccess, fmt.Sprintf("Подключение «%s» добавлено", conn.Name))
		s.clearForm()
		return nil
	}

	s.show(components.MessageTypeSuccess, fmt.Sprintf("Подключение «%s» обновлено", conn.Name))
	s.updateViewportContent()
	return ui.NavigateToCmd("connections")
}

func (s *ConnectionFormScreen) testConnection() tea.Cmd {
	conn, ok := s.readForm()
	if !ok {
		return nil
	}
	s.checking = true
	s.checkID++
	id := s.checkID
	s.show(components.MessageTypeInfo, fmt.Sprintf("Проверяем %s:%d...", conn.Host, conn.Port))
	s.updateViewportContent()
	return func() tea.Msg {
		return formHostCheckMsg{id: id, result: ssh.Check(conn.Host, conn.Port)}
	}
}

func (s *ConnectionFormScreen) showHostCheck(result ssh.CheckResult) {
	switch result.Status {
	case ssh.StatusTrusted:
		s.show(components.MessageTypeSuccess, fmt.Sprintf("Сервер доступен, ключ совпадает (%s)", result.PresentedFingerprint))
	case ssh.StatusUnknown:
		s.show(components.MessageTypeInfo, fmt.Sprintf("Сервер доступен. Новый ключ %s. Он сохранится при подключении.", result.PresentedFingerprint))
	case ssh.StatusChanged:
		s.show(components.MessageTypeError, fmt.Sprintf("Ключ изменился. Было %s, стало %s. При подключении его можно перезаписать.", result.StoredFingerprint, result.PresentedFingerprint))
	default:
		text := "сервер недоступен"
		if result.Err != nil {
			text = result.Err.Error()
		}
		s.show(components.MessageTypeError, "Проверка не удалась: "+text)
	}
}

func (s *ConnectionFormScreen) deleteConnection() tea.Cmd {
	if s.connection == nil || s.connectionSvc == nil {
		s.show(components.MessageTypeError, "Подключение не загружено")
		s.updateViewportContent()
		return nil
	}
	if err := s.connectionSvc.DeleteConnection(s.connection.ID); err != nil {
		s.show(components.MessageTypeError, fmt.Sprintf("Ошибка удаления: %v", err))
		s.updateViewportContent()
		return nil
	}
	return ui.NavigateToCmd("connections")
}

func (s *ConnectionFormScreen) readForm() (*models.Connection, bool) {
	if !s.formManager.ValidateAll() {
		s.updateViewportContent()
		s.focusFirstError()
		return nil, false
	}

	values := s.formManager.GetValues()
	port := 22
	if strings.TrimSpace(values[components.FieldNamePort]) != "" {
		parsed, err := strconv.Atoi(strings.TrimSpace(values[components.FieldNamePort]))
		if err != nil {
			return nil, false
		}
		port = parsed
	}

	usePassword := values[components.FieldNameAuth] == "true"
	now := time.Now()
	conn := &models.Connection{
		Name:      strings.TrimSpace(values[components.FieldNameName]),
		Host:      strings.TrimSpace(values[components.FieldNameHost]),
		Port:      port,
		User:      strings.TrimSpace(values[components.FieldNameUser]),
		UseSSHKey: !usePassword,
		UpdatedAt: now,
	}
	if usePassword {
		conn.Password = values[components.FieldNamePassword]
		conn.HasPassword = conn.Password != ""
	} else {
		conn.KeyPath = strings.TrimSpace(values[components.FieldNameKey])
	}

	if s.mode == formEdit {
		if s.connection == nil {
			s.show(components.MessageTypeError, "Подключение не загружено")
			s.updateViewportContent()
			return nil, false
		}
		conn.ID = s.connection.ID
		conn.CreatedAt = s.connection.CreatedAt
		s.connection = conn
	} else {
		conn.ID = uuid.New().String()
		conn.CreatedAt = now
	}
	return conn, true
}

func (s *ConnectionFormScreen) focusFirstError() {
	for _, name := range s.formManager.GetFieldOrder() {
		field := s.formManager.GetField(name)
		if field != nil && field.IsVisible() && field.HasError() {
			s.formManager.SetCurrentField(name)
			s.formManager.UpdateFocus()
			s.scrollToCurrentField()
			return
		}
	}
}

func (s *ConnectionFormScreen) togglePasswordVisibility() {
	field := s.formManager.GetField(components.FieldNamePassword)
	if field == nil {
		return
	}
	input, ok := field.GetTextInput()
	if !ok {
		return
	}
	if input.EchoMode == textinput.EchoPassword {
		input.EchoMode = textinput.EchoNormal
	} else {
		input.EchoMode = textinput.EchoPassword
	}
	field.SetTextInput(input)
	s.updateViewportContent()
}

func (s *ConnectionFormScreen) clearForm() {
	for _, name := range s.formManager.GetFieldOrder() {
		field := s.formManager.GetField(name)
		if field == nil || field.IsButton() {
			continue
		}
		field.SetValue("")
		field.SetErrorText("")
	}
	s.confirmDelete = false
	s.formManager.SetCurrentField(components.FieldNameName)
	s.formManager.UpdateFocus()
	s.updateFieldVisibility()
	if s.ready {
		s.viewport.SetYOffset(0)
		s.updateViewportContent()
	}
}

func (s *ConnectionFormScreen) scrollToCurrentField() {
	if !s.ready {
		return
	}
	current := s.formManager.GetCurrentField()
	visibleIndex := 0
	found := -1
	for _, name := range s.formManager.GetFieldOrder() {
		field := s.formManager.GetField(name)
		if field == nil || !field.IsVisible() {
			continue
		}
		if name == current {
			found = visibleIndex
			break
		}
		visibleIndex++
	}
	if found < 0 {
		return
	}

	targetLine := found * 4
	viewportHeight := s.viewport.Height
	total := s.viewport.TotalLineCount()
	if total <= viewportHeight {
		s.viewport.SetYOffset(0)
		return
	}
	offset := targetLine - viewportHeight/3
	maxOffset := total - viewportHeight
	if offset < 0 {
		offset = 0
	} else if offset > maxOffset {
		offset = maxOffset
	}
	s.viewport.SetYOffset(offset)
}

func (s *ConnectionFormScreen) show(kind components.MessageType, text string) {
	s.messageManager.ClearMessages()
	s.messageManager.AddMessage(components.Message{
		Type:      kind,
		Text:      text,
		Timestamp: time.Now(),
		Duration:  8 * time.Second,
	})
}

// Init включает курсор в первом поле.
func (s *ConnectionFormScreen) Init() tea.Cmd {
	if s.formManager.GetCurrentField() == "" {
		s.formManager.SetCurrentField(components.FieldNameName)
	}
	s.formManager.UpdateFocus()
	return textinput.Blink
}

// GetName возвращает имя экрана.
func (s *ConnectionFormScreen) GetName() string {
	if s.mode == formEdit {
		return "edit_connection"
	}
	return "add_connection"
}
