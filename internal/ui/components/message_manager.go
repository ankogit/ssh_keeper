package components

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MessageType определяет тип сообщения
type MessageType int

const (
	MessageTypeSuccess MessageType = iota
	MessageTypeError
	MessageTypeWarning
	MessageTypeInfo
)

// Message представляет сообщение для отображения
type Message struct {
	Type      MessageType
	Text      string
	Timestamp time.Time
	Duration  time.Duration
}

// NewMessage создает новое сообщение
func NewMessage(msgType MessageType, text string) Message {
	return Message{
		Type:      msgType,
		Text:      text,
		Timestamp: time.Now(),
		Duration:  3 * time.Second, // По умолчанию показываем 3 секунды
	}
}

// MessageManager управляет сообщениями
type MessageManager struct {
	messages []Message
	timeout  time.Duration
}

// NewMessageManager создает новый менеджер сообщений
func NewMessageManager() *MessageManager {
	return &MessageManager{
		messages: make([]Message, 0),
		timeout:  3 * time.Second,
	}
}

// AddMessage добавляет сообщение
func (mm *MessageManager) AddMessage(msg Message) {
	mm.messages = append(mm.messages, msg)
}

// AddSuccess добавляет сообщение об успехе
func (mm *MessageManager) AddSuccess(text string) {
	mm.AddMessage(NewMessage(MessageTypeSuccess, text))
}

// AddError добавляет сообщение об ошибке
func (mm *MessageManager) AddError(text string) {
	mm.AddMessage(NewMessage(MessageTypeError, text))
}

// AddWarning добавляет предупреждение
func (mm *MessageManager) AddWarning(text string) {
	mm.AddMessage(NewMessage(MessageTypeWarning, text))
}

// AddInfo добавляет информационное сообщение
func (mm *MessageManager) AddInfo(text string) {
	mm.AddMessage(NewMessage(MessageTypeInfo, text))
}

// GetMessages возвращает все активные сообщения
func (mm *MessageManager) GetMessages() []Message {
	now := time.Now()
	var activeMessages []Message

	for _, msg := range mm.messages {
		if now.Sub(msg.Timestamp) < msg.Duration {
			activeMessages = append(activeMessages, msg)
		}
	}

	// Обновляем список сообщений, удаляя устаревшие
	mm.messages = activeMessages
	return activeMessages
}

// ClearMessages очищает все сообщения
func (mm *MessageManager) ClearMessages() {
	mm.messages = make([]Message, 0)
}

// RenderMessages отображает все активные сообщения
func (mm *MessageManager) RenderMessages(width int) string {
	messages := mm.GetMessages()
	if len(messages) == 0 {
		return ""
	}

	var rendered []string
	for _, msg := range messages {
		rendered = append(rendered, mm.renderMessage(msg, width))
	}

	return fmt.Sprintf("%s\n", fmt.Sprintf("%s", rendered))
}

// renderMessage отображает одно сообщение
func (mm *MessageManager) renderMessage(msg Message, width int) string {
	var style lipgloss.Style
	var icon string

	switch msg.Type {
	case MessageTypeSuccess:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Background(lipgloss.Color("#004400")).
			Padding(0, 1)
		icon = "+"
	case MessageTypeError:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Background(lipgloss.Color("#440000")).
			Padding(0, 1)
		icon = "-"
	case MessageTypeWarning:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffaa00")).
			Background(lipgloss.Color("#442200")).
			Padding(0, 1)
		icon = "!"
	case MessageTypeInfo:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0088ff")).
			Background(lipgloss.Color("#002244")).
			Padding(0, 1)
		icon = "ℹ"
	}

	// Ограничиваем ширину сообщения
	if width > 0 && len(msg.Text) > width-10 {
		msg.Text = msg.Text[:width-10] + "..."
	}

	return style.Render(fmt.Sprintf("%s %s", icon, msg.Text))
}

// MessageCmd команда для обновления сообщений
type MessageCmd struct {
	Message Message
}

// AddMessageCmd создает команду для добавления сообщения
func AddMessageCmd(msg Message) tea.Cmd {
	return func() tea.Msg {
		return MessageCmd{Message: msg}
	}
}

// AddSuccessCmd создает команду для добавления сообщения об успехе
func AddSuccessCmd(text string) tea.Cmd {
	return AddMessageCmd(NewMessage(MessageTypeSuccess, text))
}

// AddErrorCmd создает команду для добавления сообщения об ошибке
func AddErrorCmd(text string) tea.Cmd {
	return AddMessageCmd(NewMessage(MessageTypeError, text))
}
