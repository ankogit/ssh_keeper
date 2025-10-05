package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"ssh-keeper/internal/config"
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

// restoreTerminal восстанавливает терминал после SSH сессий
func restoreTerminal() {
	cmd := exec.Command("reset")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Failed to reset terminal: %v\n", err)
	}
}

func main() {
	// ПРИНУДИТЕЛЬНО восстанавливаем терминал через reset в самом начале
	restoreTerminal()

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

	// Восстанавливаем терминал после нормального завершения
	restoreTerminal()

	// Set up signal handler to restore terminal on exit
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-c
	// 	fmt.Println("QQQQQQQQQQQ: Restoring terminal...")
	// 	exec.Command("reset").Run()
	// 	os.Exit(0)
	// }()
	exec.Command("reset").Run()
}

// initializeServices initializes all application services
func initializeServices() error {
	// Инициализируем конфигурацию
	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Получаем путь к конфигурации из настроек
	configPath := cfg.GetConfigPath()

	// Разворачиваем ~ в полный путь
	if strings.HasPrefix(configPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, configPath[2:])
	}

	// Создаем директорию конфигурации если она не существует
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Initialize master password service
	masterPasswordService := services.NewMasterPasswordService()
	services.SetGlobalMasterPasswordService(masterPasswordService)

	// Initialize encryption service
	encryptionService := services.NewEncryptionService(masterPasswordService)
	services.SetGlobalEncryptionService(encryptionService)

	// If master password is already initialized with signature, refresh the encryption key
	if services.IsMasterPasswordInitializedWithSignature() {
		if err := encryptionService.RefreshKey(); err != nil {
			fmt.Printf("Warning: Failed to refresh encryption key: %v\n", err)
		}
	}

	// Initialize connection service
	connectionService = services.NewConnectionService(configPath)

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
