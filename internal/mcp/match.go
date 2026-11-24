package mcp

import (
	"fmt"
	"strings"

	"ssh-keeper/internal/models"
)

// MatchConnection находит одно подключение по имени, хосту или пользователю.
// Все слова запроса должны встретиться. Точное имя важнее частичного совпадения.
func MatchConnection(connections []models.Connection, query string) (models.Connection, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return models.Connection{}, fmt.Errorf("не указано имя подключения")
	}

	var exact []models.Connection
	var partial []models.Connection
	for _, conn := range connections {
		if strings.EqualFold(strings.TrimSpace(conn.Name), query) {
			exact = append(exact, conn)
			continue
		}
		if containsAllWords(conn, query) {
			partial = append(partial, conn)
		}
	}

	switch {
	case len(exact) == 1:
		return exact[0], nil
	case len(exact) > 1:
		return models.Connection{}, fmt.Errorf("несколько подключений с именем %q:\n%s", query, formatMatches(exact))
	case len(partial) == 1:
		return partial[0], nil
	case len(partial) > 1:
		return models.Connection{}, fmt.Errorf("запрос %q подходит нескольким серверам:\n%s", query, formatMatches(partial))
	default:
		return models.Connection{}, fmt.Errorf("подключение %q не найдено", query)
	}
}

func containsAllWords(conn models.Connection, query string) bool {
	haystack := strings.ToLower(strings.Join([]string{conn.Name, conn.Host, conn.User}, " "))
	for _, word := range strings.Fields(strings.ToLower(query)) {
		if !strings.Contains(haystack, word) {
			return false
		}
	}
	return true
}

func formatMatches(connections []models.Connection) string {
	lines := make([]string, 0, len(connections))
	for _, conn := range connections {
		lines = append(lines, "- "+describe(conn))
	}
	return strings.Join(lines, "\n")
}

func describe(conn models.Connection) string {
	port := conn.Port
	if port == 0 {
		port = 22
	}
	auth := "пароль"
	if conn.UseSSHKey {
		auth = "ключ"
	}
	name := conn.Name
	if name == "" {
		name = conn.Host
	}
	return fmt.Sprintf("%s | %s@%s:%d | %s", name, conn.User, conn.Host, port, auth)
}

// FormatConnections возвращает список подключений без секретов.
func FormatConnections(connections []models.Connection) string {
	if len(connections) == 0 {
		return "Сохранённых подключений нет."
	}
	lines := make([]string, 0, len(connections))
	for _, conn := range connections {
		lines = append(lines, describe(conn))
	}
	return strings.Join(lines, "\n")
}
