package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ssh-keeper/internal/config"
	"ssh-keeper/internal/mcp"
	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ssh"
	"ssh-keeper/internal/ui/screens"
	"ssh-keeper/internal/version"

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
	appVersion   string = "dev" // Default version for development
	appBuildTime string = "dev" // Default build time for development
	appSignature string
)

func main() {
	if os.Getenv(ssh.AskpassEnv) == "1" {
		runAskpass()
		return
	}

	// Проверяем флаги командной строки
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("SSH Keeper version %s\n", version.GetVersion())
			return
		case "--help", "-h":
			fmt.Printf("SSH Keeper - Secure SSH Connection Manager\n")
			fmt.Printf("Version: %s\n", version.GetVersion())
			fmt.Printf("\nUsage: ssh-keeper [options]\n")
			fmt.Printf("\nOptions:\n")
			fmt.Printf("  --version, -v    Show version information\n")
			fmt.Printf("  --help, -h       Show this help message\n")
			fmt.Printf("  mcp              Run the local MCP server on stdin/stdout\n")
			return
		case "mcp":
			if err := initializeServices(); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing services: %v\n", err)
				os.Exit(1)
			}
			if err := mcp.Serve(); err != nil {
				fmt.Fprintf(os.Stderr, "MCP server: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	// Set up terminal environment
	lipgloss.SetColorProfile(termenv.ColorProfile())

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
		os.Exit(1)
	}
}

// runAskpass печатает сохранённый пароль для OpenSSH и сразу завершается.
func runAskpass() {
	data, err := os.ReadFile(os.Getenv(ssh.PassFileEnv))
	if err != nil {
		fmt.Fprintf(os.Stderr, "askpass: %v\n", err)
		os.Exit(1)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	if _, err := os.Stdout.Write(data); err != nil {
		os.Exit(1)
	}
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
			fmt.Fprintf(os.Stderr, "Warning: Failed to refresh encryption key: %v\n", err)
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

// GetVersion возвращает версию приложения
func GetVersion() string {
	return version.GetVersion()
}
