package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// MenuAction представляет функцию, которая выполняется при выборе пункта меню
type MenuAction func() tea.Cmd

// MenuItemConfig представляет конфигурацию элемента меню
type MenuItemConfig struct {
	Title       string
	Description string
	Action      MenuAction
	Shortcut    string // Горячая клавиша (опционально)
}

// MenuItem представляет элемент меню с конфигурацией
type MenuItem struct {
	config MenuItemConfig
}

// NewMenuItem создает новый элемент меню
func NewMenuItem(config MenuItemConfig) MenuItem {
	return MenuItem{
		config: config,
	}
}

// Title возвращает заголовок элемента меню
func (m MenuItem) Title() string {
	return m.config.Title
}

// Description возвращает описание элемента меню
func (m MenuItem) Description() string {
	return m.config.Description
}

// FilterValue возвращает значение для фильтрации
func (m MenuItem) FilterValue() string {
	return m.config.Title
}

// Execute выполняет действие элемента меню
func (m MenuItem) Execute() tea.Cmd {
	if m.config.Action != nil {
		return m.config.Action()
	}
	return nil
}

// GetShortcut возвращает горячую клавишу
func (m MenuItem) GetShortcut() string {
	return m.config.Shortcut
}

// MenuConfig представляет конфигурацию меню
type MenuConfig struct {
	Title    string
	Items    []MenuItemConfig
	ShowBack bool // Показывать кнопку "Назад"
	ShowQuit bool // Показывать кнопку "Выход"
}
