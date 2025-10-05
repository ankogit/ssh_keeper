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

// Build-time variables
var (
	version      string
	appSignature string
)

// restoreTerminal восстанавливает терминал после SSH сессий
func restoreTerminal() {
	// Проверяем, что мы в интерактивном терминале
	if !isTerminal(os.Stdin) {
		return
	}

	// Пробуем сначала tput reset, потом reset
	cmd := exec.Command("tput", "reset")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Fallback к обычному reset
		cmd = exec.Command("reset")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			// Игнорируем ошибки в неинтерактивных режимах
			if isTerminal(os.Stdin) {
				fmt.Printf("Warning: Failed to reset terminal: %v\n", err)
			}
		}
	}
}

// isTerminal проверяет, является ли файл терминалом
func isTerminal(file *os.File) bool {
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	// Проверяем что это символьное устройство И что это не pipe/redirect
	return (stat.Mode()&os.ModeCharDevice) != 0 && file.Name() != "/dev/null"
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

	// Terminal will be restored by signal handler if needed
}

// initializeServices initializes all application services
func initializeServices() error {
	// Устанавливаем встроенную подпись если она есть (из CI build)
	// Это нужно сделать ДО загрузки конфигурации
	if appSignature != "" {
		os.Setenv("SECURITY_APP_SIGNATURE", appSignature)
	}

	// Инициализируем конфигурацию (это загрузит .env файл если он есть)
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

	// // Инициализируем сервис автоматических обновлений
	// autoUpdateService := services.NewAutoUpdateService(cfg)
	// services.SetGlobalAutoUpdateService(autoUpdateService)

	// // Выполняем автоматическую проверку обновлений если нужно
	// if autoUpdateService.CheckIfUpdateNeeded() {
	// 	go func() {
	// 		updateInfo, err := autoUpdateService.PerformAutoCheck()
	// 		if err != nil {
	// 			fmt.Printf("Auto-update check failed: %v\n", err)
	// 		} else if updateInfo != nil && updateInfo.IsAvailable {
	// 			// Уведомление об обновлениях временно отключено
	// 			// fmt.Printf("🔄 New version %s is available! Check Settings > Updates to install.\n", updateInfo.Version)
	// 		}
	// 	}()
	// }

	return nil
}

// GetConnectionService returns the global connection service
func GetConnectionService() *services.ConnectionService {
	return connectionService
}
