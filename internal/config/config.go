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
	Security struct {
		AppSignature string `envconfig:"APP_SIGNATURE"`
	} `envconfig:"SECURITY"`

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
	// Проверяем обязательные поля
	if c.Security.AppSignature == "" {
		return fmt.Errorf("APP_SIGNATURE is required")
	}

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
	return c.Security.AppSignature
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
	if c.Security.AppSignature == "" {
		return fmt.Errorf("подпись приложения не загружена")
	}

	// В режиме разработки принимаем подпись ssh-keeper-sig-dev
	if c.IsDevelopment() && c.Security.AppSignature == "ssh-keeper-sig-dev" {
		return nil
	}

	// В продакшене здесь будет более строгая проверка
	// Например, проверка хеша приложения, цифровой подписи и т.д.

	return nil
}
