package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ssh-keeper/internal/config"
)

// AutoUpdateService управляет автоматическими проверками обновлений
type AutoUpdateService struct {
	config        *config.Config
	updateService *UpdateService
}

// NewAutoUpdateService создает новый сервис автоматических обновлений
func NewAutoUpdateService(cfg *config.Config) *AutoUpdateService {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "0.1.0"
	}

	return &AutoUpdateService{
		config:        cfg,
		updateService: NewUpdateService(version),
	}
}

// CheckIfUpdateNeeded проверяет, нужно ли проверить обновления
func (aus *AutoUpdateService) CheckIfUpdateNeeded() bool {
	// Если автоматическая проверка отключена
	if !aus.config.IsAutoUpdateEnabled() {
		return false
	}

	// Получаем время последней проверки
	lastCheckStr := aus.config.GetLastUpdateCheck()
	if lastCheckStr == "" {
		// Если никогда не проверяли, проверяем сейчас
		return true
	}

	// Парсим время последней проверки
	lastCheck, err := time.Parse(time.RFC3339, lastCheckStr)
	if err != nil {
		// Если не удалось распарсить, проверяем сейчас
		return true
	}

	// Проверяем, прошло ли достаточно времени
	interval := time.Duration(aus.config.GetUpdateCheckInterval()) * time.Hour
	return time.Since(lastCheck) >= interval
}

// PerformAutoCheck выполняет автоматическую проверку обновлений
func (aus *AutoUpdateService) PerformAutoCheck() (*UpdateInfo, error) {
	// Проверяем обновления
	updateInfo, err := aus.updateService.CheckForUpdates()
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	// Сохраняем время последней проверки
	if err := aus.saveLastCheckTime(); err != nil {
		// Логируем ошибку, но не прерываем выполнение
		fmt.Printf("Warning: Failed to save last check time: %v\n", err)
	}

	return updateInfo, nil
}

// saveLastCheckTime сохраняет время последней проверки
func (aus *AutoUpdateService) saveLastCheckTime() error {
	// Получаем путь к конфигурации
	configPath := aus.config.GetConfigPath()

	// Разворачиваем ~ в полный путь
	if configPath[0] == '~' {
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

	// Создаем файл с временем последней проверки
	lastCheckFile := filepath.Join(configDir, ".last_update_check")
	now := time.Now().Format(time.RFC3339)

	return os.WriteFile(lastCheckFile, []byte(now), 0644)
}

// GetLastCheckTime возвращает время последней проверки
func (aus *AutoUpdateService) GetLastCheckTime() (time.Time, error) {
	// Получаем путь к конфигурации
	configPath := aus.config.GetConfigPath()

	// Разворачиваем ~ в полный путь
	if configPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, configPath[2:])
	}

	// Читаем файл с временем последней проверки
	configDir := filepath.Dir(configPath)
	lastCheckFile := filepath.Join(configDir, ".last_update_check")

	data, err := os.ReadFile(lastCheckFile)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Time{}, nil // Файл не существует
		}
		return time.Time{}, fmt.Errorf("failed to read last check file: %w", err)
	}

	// Парсим время
	return time.Parse(time.RFC3339, string(data))
}

// GetUpdateService возвращает сервис обновлений
func (aus *AutoUpdateService) GetUpdateService() *UpdateService {
	return aus.updateService
}
