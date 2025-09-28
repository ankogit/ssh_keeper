package components

import (
	"fmt"
	"ssh-keeper/internal/models"
	"ssh-keeper/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ConnectionItem представляет элемент подключения для списка
type ConnectionItem struct {
	Connection models.Connection
}

// NewConnectionItem создает новый элемент подключения
func NewConnectionItem(conn models.Connection) ConnectionItem {
	return ConnectionItem{
		Connection: conn,
	}
}

// Title возвращает заголовок элемента
func (ci ConnectionItem) Title() string {
	return ci.Connection.Name
}

// Description возвращает описание элемента (компактное)
func (ci ConnectionItem) Description() string {
	// Компактная информация в одну строку
	hostInfo := fmt.Sprintf("%s:%d", ci.Connection.Host, ci.Connection.Port)
	userInfo := ci.Connection.User

	// Тип аутентификации (только иконка)
	var authIcon string
	if ci.Connection.KeyPath != "" {
		authIcon = "🔑"
	} else if ci.Connection.HasPassword {
		authIcon = "🔒"
	} else {
		authIcon = "❓"
	}

	return fmt.Sprintf("%s | %s | %s", hostInfo, userInfo, authIcon)
}

// FilterValue возвращает значение для фильтрации
func (ci ConnectionItem) FilterValue() string {
	// Поиск по названию, хосту и пользователю
	return fmt.Sprintf("%s %s %s",
		ci.Connection.Name,
		ci.Connection.Host,
		ci.Connection.User)
}

// GetConnection возвращает подключение
func (ci ConnectionItem) GetConnection() models.Connection {
	return ci.Connection
}

// RenderCustomItem создает кастомное отображение элемента
func (ci ConnectionItem) RenderCustomItem() string {
	// Стили для различных частей
	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorSecondary)).
		Bold(styles.TextBold)

	hostStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorText))

	userStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorMuted))

	authStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorSecondary))

	// Форматируем части
	name := nameStyle.Render(ci.Connection.Name)
	host := hostStyle.Render(fmt.Sprintf("(%s:%d)", ci.Connection.Host, ci.Connection.Port))
	user := userStyle.Render(fmt.Sprintf("пользователь: %s", ci.Connection.User))

	var auth string
	if ci.Connection.KeyPath != "" {
		auth = authStyle.Render("🔑 ключ")
	} else if ci.Connection.HasPassword {
		auth = authStyle.Render("🔒 пароль")
	} else {
		auth = authStyle.Render("❓ неизвестно")
	}

	return fmt.Sprintf("%s %s | %s | %s", name, host, user, auth)
}
