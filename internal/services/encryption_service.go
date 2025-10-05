package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptionService handles encryption/decryption of sensitive data
type EncryptionService struct {
	masterPasswordService *MasterPasswordService
	derivedKey            []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(masterPasswordService *MasterPasswordService) *EncryptionService {
	es := &EncryptionService{
		masterPasswordService: masterPasswordService,
		derivedKey:            nil,
	}

	// Инициализируем ключ если мастер-пароль уже установлен
	if masterPasswordService.IsInitialized() {
		es.refreshKey()
	}

	return es
}

// Encrypt encrypts a plaintext string
func (es *EncryptionService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// Проверяем, что ключ инициализирован
	if es.derivedKey == nil {
		return "", fmt.Errorf("ключ шифрования не инициализирован")
	}

	// Create AES cipher
	block, err := aes.NewCipher(es.derivedKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts an encrypted string
func (es *EncryptionService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Проверяем, что ключ инициализирован
	if es.derivedKey == nil {
		return "", fmt.Errorf("ключ шифрования не инициализирован")
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(es.derivedKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// EncryptPassword encrypts a password for storage
func (es *EncryptionService) EncryptPassword(password string) (string, error) {
	return es.Encrypt(password)
}

// DecryptPassword decrypts a password from storage
func (es *EncryptionService) DecryptPassword(encryptedPassword string) (string, error) {
	return es.Decrypt(encryptedPassword)
}

// RefreshKey обновляет ключ шифрования из мастер-пароля
func (es *EncryptionService) RefreshKey() error {
	return es.refreshKey()
}

// refreshKey внутренний метод для обновления ключа
func (es *EncryptionService) refreshKey() error {
	masterPassword, err := es.masterPasswordService.GetMasterPassword()
	if err != nil {
		return fmt.Errorf("не удалось получить мастер-пароль: %w", err)
	}

	es.derivedKey = es.masterPasswordService.DeriveKey(masterPassword)
	return nil
}

// IsInitialized проверяет, инициализирован ли сервис шифрования
func (es *EncryptionService) IsInitialized() bool {
	return es.derivedKey != nil && es.masterPasswordService.IsInitialized()
}
