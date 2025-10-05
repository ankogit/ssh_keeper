package services

import (
	"crypto/sha256"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	// ServiceName имя сервиса для go-keyring
	ServiceName = "ssh-keeper"
	// MasterPasswordKey ключ для хранения мастер-пароля
	MasterPasswordKey = "master-password"
	// SignatureTokenKey ключ для хранения токена подписи
	SignatureTokenKey = "signature-token"
	// RequirePasswordOnStartupKey ключ для настройки запроса пароля при запуске
	RequirePasswordOnStartupKey = "require-password-on-startup"
)

// MasterPasswordService управляет мастер-паролем через go-keyring
type MasterPasswordService struct {
	securityConfigService *SecurityConfigService
}

// NewMasterPasswordService создает новый сервис для работы с мастер-паролем
func NewMasterPasswordService() *MasterPasswordService {
	securityConfigService := NewSecurityConfigService()
	return &MasterPasswordService{
		securityConfigService: securityConfigService,
	}
}

// SetMasterPassword сохраняет мастер-пароль в системном хранилище
func (mps *MasterPasswordService) SetMasterPassword(password string) error {
	if password == "" {
		return fmt.Errorf("мастер-пароль не может быть пустым")
	}

	return keyring.Set(ServiceName, MasterPasswordKey, password)
}

// GetMasterPassword получает мастер-пароль из системного хранилища с проверкой подписи
func (mps *MasterPasswordService) GetMasterPassword() (string, error) {
	// ВСЕГДА проверяем подпись при получении мастер-пароля
	if err := mps.securityConfigService.ValidateSignature(); err != nil {
		return "", fmt.Errorf("приложение не прошло проверку подписи: %w", err)
	}

	password, err := keyring.Get(ServiceName, MasterPasswordKey)
	if err != nil {
		return "", fmt.Errorf("не удалось получить мастер-пароль: %w", err)
	}
	return password, nil
}

// IsInitialized проверяет, инициализирован ли мастер-пароль с проверкой подписи
func (mps *MasterPasswordService) IsInitialized() bool {
	// ВСЕГДА проверяем подпись при проверке инициализации
	if err := mps.securityConfigService.ValidateSignature(); err != nil {
		return false
	}

	_, err := keyring.Get(ServiceName, MasterPasswordKey)
	return err == nil
}

// ClearMasterPassword удаляет мастер-пароль из системного хранилища с проверкой подписи
func (mps *MasterPasswordService) ClearMasterPassword() error {
	// ВСЕГДА проверяем подпись при удалении мастер-пароля
	if err := mps.securityConfigService.ValidateSignature(); err != nil {
		return fmt.Errorf("приложение не прошло проверку подписи: %w", err)
	}

	return keyring.Delete(ServiceName, MasterPasswordKey)
}

// DeriveKey создает ключ шифрования из мастер-пароля
func (mps *MasterPasswordService) DeriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// ValidateMasterPassword проверяет корректность мастер-пароля
func (mps *MasterPasswordService) ValidateMasterPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("мастер-пароль должен содержать минимум 8 символов")
	}

	// Дополнительные проверки можно добавить здесь
	// Например, проверка на наличие цифр, специальных символов и т.д.

	return nil
}

// ChangeMasterPassword меняет мастер-пароль
func (mps *MasterPasswordService) ChangeMasterPassword(oldPassword, newPassword string) error {
	// Проверяем старый пароль
	currentPassword, err := mps.GetMasterPassword()
	if err != nil {
		return fmt.Errorf("не удалось получить текущий мастер-пароль: %w", err)
	}

	if currentPassword != oldPassword {
		return fmt.Errorf("неверный текущий мастер-пароль")
	}

	// Валидируем новый пароль
	if err := mps.ValidateMasterPassword(newPassword); err != nil {
		return fmt.Errorf("неверный новый мастер-пароль: %w", err)
	}

	// Сохраняем новый пароль
	return mps.SetMasterPassword(newPassword)
}

// SetRequirePasswordOnStartup устанавливает настройку запроса пароля при запуске
func (mps *MasterPasswordService) SetRequirePasswordOnStartup(require bool) error {
	// ВСЕГДА проверяем подпись при изменении настроек
	if err := mps.securityConfigService.ValidateSignature(); err != nil {
		return fmt.Errorf("приложение не прошло проверку подписи: %w", err)
	}

	value := "false"
	if require {
		value = "true"
	}

	return keyring.Set(ServiceName, RequirePasswordOnStartupKey, value)
}

// GetRequirePasswordOnStartup получает настройку запроса пароля при запуске
func (mps *MasterPasswordService) GetRequirePasswordOnStartup() (bool, error) {
	// ВСЕГДА проверяем подпись при получении настроек
	if err := mps.securityConfigService.ValidateSignature(); err != nil {
		return false, fmt.Errorf("приложение не прошло проверку подписи: %w", err)
	}

	value, err := keyring.Get(ServiceName, RequirePasswordOnStartupKey)
	if err != nil {
		// Если настройка не найдена, возвращаем значение по умолчанию (true)
		return true, nil
	}

	return value == "true", nil
}

// ClearRequirePasswordOnStartup удаляет настройку запроса пароля при запуске
func (mps *MasterPasswordService) ClearRequirePasswordOnStartup() error {
	// ВСЕГДА проверяем подпись при удалении настроек
	if err := mps.securityConfigService.ValidateSignature(); err != nil {
		return fmt.Errorf("приложение не прошло проверку подписи: %w", err)
	}

	return keyring.Delete(ServiceName, RequirePasswordOnStartupKey)
}
