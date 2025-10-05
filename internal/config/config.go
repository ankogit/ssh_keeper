package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	DevEnv  = "development"
	ProdEnv = "production"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	// Основные настройки
	Debug      bool   `envconfig:"DEBUG" default:"false"`
	Env        string `envconfig:"ENV" default:"development"`
	ConfigPath string `envconfig:"CONFIG_PATH" default:"~/.ssh-keeper/config"`

	// Настройки безопасности
	AppSignature string `envconfig:"SECURITY_APP_SIGNATURE"`

	// Настройки SSH
	SSH struct {
		ConfigPath string `envconfig:"SSH_CONFIG_PATH" default:"~/.ssh/config"`
	} `envconfig:"SSH"`

	// Настройки приложения
	App struct {
		Name        string `envconfig:"APP_NAME" default:"ssh-keeper"`
		Version     string `envconfig:"APP_VERSION" default:"1.0.0"`
		Description string `envconfig:"APP_DESCRIPTION" default:"SSH Connection Manager"`
	} `envconfig:"APP"`

	// Настройки логирования
	Logging struct {
		Level  string `envconfig:"LOG_LEVEL" default:"info"`
		Format string `envconfig:"LOG_FORMAT" default:"text"`
	} `envconfig:"LOGGING"`

	// Настройки обновлений
	Updates struct {
		AutoCheck     bool   `envconfig:"AUTO_CHECK_UPDATES" default:"true"`
		CheckInterval int    `envconfig:"UPDATE_CHECK_INTERVAL" default:"24"` // часы
		LastCheck     string `envconfig:"LAST_UPDATE_CHECK" default:""`
	} `envconfig:"UPDATES"`
}

// Init инициализирует конфигурацию из переменных окружения и файлов
func Init() (*Config, error) {
	var cfg Config

	// Загружаем .env файл если он существует
	if err := godotenv.Load(); err != nil {
		// Игнорируем ошибку если файл не найден
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	// Парсим переменные окружения
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	// Валидируем конфигурацию
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Проверяем допустимые значения
	if c.Env != DevEnv && c.Env != ProdEnv {
		return fmt.Errorf("invalid ENV value: %s, must be %s or %s", c.Env, DevEnv, ProdEnv)
	}

	return nil
}

// IsDevelopment проверяет, работает ли приложение в режиме разработки
func (c *Config) IsDevelopment() bool {
	return c.Env == DevEnv
}

// IsProduction проверяет, работает ли приложение в продакшене
func (c *Config) IsProduction() bool {
	return c.Env == ProdEnv
}

// GetAppSignature возвращает подпись приложения
func (c *Config) GetAppSignature() string {
	return c.AppSignature
}

// GetConfigPath возвращает путь к файлу конфигурации приложения
func (c *Config) GetConfigPath() string {
	return c.ConfigPath
}

// GetSSHConfigPath возвращает путь к SSH конфигурации
func (c *Config) GetSSHConfigPath() string {
	return c.SSH.ConfigPath
}

// ValidateAppSignature проверяет подпись приложения
func (c *Config) ValidateAppSignature() error {
	// В development режиме пропускаем проверку подписи
	if c.Env == DevEnv {
		return nil
	}

	if c.AppSignature == "" {
		return fmt.Errorf("подпись приложения не загружена")
	}

	// Простая проверка - любая непустая подпись принимается
	// В будущем здесь будет более строгая проверка
	// Например, проверка хеша приложения, цифровой подписи и т.д.

	return nil
}

// GetUpdatesConfig возвращает настройки обновлений
func (c *Config) GetUpdatesConfig() *struct {
	AutoCheck     bool   `envconfig:"AUTO_CHECK_UPDATES" default:"true"`
	CheckInterval int    `envconfig:"UPDATE_CHECK_INTERVAL" default:"24"` // часы
	LastCheck     string `envconfig:"LAST_UPDATE_CHECK" default:""`
} {
	return &c.Updates
}

// IsAutoUpdateEnabled проверяет, включена ли автоматическая проверка обновлений
func (c *Config) IsAutoUpdateEnabled() bool {
	return c.Updates.AutoCheck
}

// GetUpdateCheckInterval возвращает интервал проверки обновлений в часах
func (c *Config) GetUpdateCheckInterval() int {
	return c.Updates.CheckInterval
}

// GetLastUpdateCheck возвращает время последней проверки обновлений
func (c *Config) GetLastUpdateCheck() string {
	return c.Updates.LastCheck
}
