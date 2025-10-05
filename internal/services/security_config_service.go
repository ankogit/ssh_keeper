package services

import (
	"fmt"
	"ssh-keeper/internal/config"
)

// SecurityConfigService управляет конфигурацией безопасности
type SecurityConfigService struct {
	config *config.Config
}

// NewSecurityConfigService создает новый сервис конфигурации безопасности
func NewSecurityConfigService() *SecurityConfigService {
	// Инициализируем конфигурацию
	cfg, err := config.Init()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
		// Создаем минимальную конфигурацию для fallback
		cfg = &config.Config{}
	}

	return &SecurityConfigService{
		config: cfg,
	}
}

// GetAppSignature возвращает подпись приложения из конфигурации
func (scs *SecurityConfigService) GetAppSignature() string {
	return scs.config.GetAppSignature()
}

// ValidateSignature проверяет подпись приложения против конфигурации
func (scs *SecurityConfigService) ValidateSignature() error {
	return scs.config.ValidateAppSignature()
}
