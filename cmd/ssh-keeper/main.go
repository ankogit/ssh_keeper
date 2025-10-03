package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Global services
var (
	connectionService *services.ConnectionService
)

func main() {
	// Set up terminal environment
	lipgloss.SetColorProfile(termenv.ColorProfile())

	// Убеждаемся, что терминал восстановится при выходе
	defer restoreTerminal()

	// Set up signal handler to restore terminal on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		restoreTerminal()
		os.Exit(0)
	}()

	// Initialize services
	if err := initializeServices(); err != nil {
		fmt.Printf("Error initializing services: %v\n", err)
		os.Exit(1)
	}

	// Create app with screen manager
	app := screens.NewApp()

	// Create tea program
	p := tea.NewProgram(app, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		restoreTerminal()
		os.Exit(1)
	}
}

// initializeServices initializes all application services
func initializeServices() error {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create SSH Keeper config directory
	configDir := filepath.Join(homeDir, ".ssh-keeper")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set config file path
	configPath := filepath.Join(configDir, "config")

	// For now, use a simple master key
	// In production, this should be derived from user input or stored securely
	masterKey := "ssh-keeper-default-key-2024"

	// Initialize connection service
	connectionService = services.NewConnectionService(configPath, masterKey)

	// Initialize with sample data if no config exists
	if err := connectionService.InitializeWithSampleData(); err != nil {
		return fmt.Errorf("failed to initialize with sample data: %w", err)
	}

	// Set global service
	services.SetGlobalConnectionService(connectionService)

	return nil
}

// GetConnectionService returns the global connection service
func GetConnectionService() *services.ConnectionService {
	return connectionService
}

// restoreTerminal восстанавливает терминал при выходе
func restoreTerminal() {
	// Радикальное восстановление терминала
	fmt.Print("\033[?1049l") // Выход из альтернативного буфера
	fmt.Print("\033[?25h")   // Показать курсор
	fmt.Print("\033[?2004l") // Отключаем bracketed paste mode
	fmt.Print("\033[?1l")    // Отключаем application cursor keys
	fmt.Print("\033[?7h")    // Включаем auto wrap mode
	fmt.Print("\033[?12l")   // Отключаем local echo
	fmt.Print("\033[?1000l") // Отключаем mouse reporting
	fmt.Print("\033[?1001l") // Отключаем mouse reporting
	fmt.Print("\033[?1002l") // Отключаем mouse reporting
	fmt.Print("\033[?1003l") // Отключаем mouse reporting
	fmt.Print("\033[?1005l") // Отключаем mouse reporting
	fmt.Print("\033[?1006l") // Отключаем mouse reporting
	fmt.Print("\033[?1015l") // Отключаем mouse reporting
	fmt.Print("\033[?25h")   // Показать курсор
	fmt.Print("\033[0m")     // Сбрасываем все атрибуты
	fmt.Print("\033c")       // Полный сброс терминала
	fmt.Print("\033[2J")     // Очищаем экран
	fmt.Print("\033[H")      // Перемещаем курсор в начало
	fmt.Print("\033[?25h")   // Показать курсор
	fmt.Print("\033[0m")     // Сбрасываем все атрибуты

	// Радикальная защита через reset и stty
	exec.Command("reset").Run()
	exec.Command("stty", "sane").Run()
	exec.Command("tput", "reset").Run()
}
