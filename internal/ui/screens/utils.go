package screens

import (
	"ssh-keeper/internal/ui"

	"github.com/charmbracelet/bubbles/list"
)

// convertToListItem конвертирует MenuItem в list.Item
func convertToListItem(menuItems []ui.MenuItem) []list.Item {
	items := make([]list.Item, len(menuItems))
	for i, item := range menuItems {
		items[i] = item
	}
	return items
}
