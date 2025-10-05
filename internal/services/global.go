package services

import (
	"fmt"
	"ssh-keeper/internal/models"
)

// Global service instances
var (
	globalConnectionService     *ConnectionService
	globalMasterPasswordService *MasterPasswordService
	globalEncryptionService     *EncryptionService
	globalSecurityConfigService *SecurityConfigService
)

// SetGlobalConnectionService sets the global connection service
func SetGlobalConnectionService(service *ConnectionService) {
	globalConnectionService = service
}

// GetGlobalConnectionService returns the global connection service
func GetGlobalConnectionService() *ConnectionService {
	return globalConnectionService
}

// GetConnections returns all connections from the global service
func GetConnections() []models.Connection {
	if globalConnectionService == nil {
		return []models.Connection{}
	}
	return globalConnectionService.GetAllConnections()
}

// AddConnection adds a connection using the global service
func AddConnection(conn *models.Connection) error {
	if globalConnectionService == nil {
		return fmt.Errorf("connection service not initialized")
	}
	return globalConnectionService.AddConnection(conn)
}

// UpdateConnection updates a connection using the global service
func UpdateConnection(conn *models.Connection) error {
	if globalConnectionService == nil {
		return fmt.Errorf("connection service not initialized")
	}
	return globalConnectionService.UpdateConnection(conn)
}

// DeleteConnection deletes a connection using the global service
func DeleteConnection(id string) error {
	if globalConnectionService == nil {
		return fmt.Errorf("connection service not initialized")
	}
	return globalConnectionService.DeleteConnection(id)
}

// GetConnectionByID gets a connection by ID using the global service
func GetConnectionByID(id string) *models.Connection {
	if globalConnectionService == nil {
		return nil
	}
	return globalConnectionService.GetConnectionByID(id)
}

// ReloadConnections reloads connections from file using the global service
func ReloadConnections() error {
	if globalConnectionService == nil {
		return fmt.Errorf("connection service not initialized")
	}
	return globalConnectionService.ReloadConnections()
}

// SetGlobalMasterPasswordService sets the global master password service
func SetGlobalMasterPasswordService(service *MasterPasswordService) {
	globalMasterPasswordService = service
}

// GetGlobalMasterPasswordService returns the global master password service
func GetGlobalMasterPasswordService() *MasterPasswordService {
	return globalMasterPasswordService
}

// SetGlobalEncryptionService sets the global encryption service
func SetGlobalEncryptionService(service *EncryptionService) {
	globalEncryptionService = service
}

// GetGlobalEncryptionService returns the global encryption service
func GetGlobalEncryptionService() *EncryptionService {
	return globalEncryptionService
}

// IsMasterPasswordInitialized checks if master password is initialized globally
func IsMasterPasswordInitialized() bool {
	if globalMasterPasswordService == nil {
		return false
	}
	return globalMasterPasswordService.IsInitialized()
}

// RefreshEncryptionKey refreshes the encryption key globally
func RefreshEncryptionKey() error {
	if globalEncryptionService == nil {
		return fmt.Errorf("encryption service not initialized")
	}
	return globalEncryptionService.RefreshKey()
}

// IsMasterPasswordInitializedWithSignature checks if master password is initialized with signature validation
func IsMasterPasswordInitializedWithSignature() bool {
	if globalMasterPasswordService == nil {
		return false
	}
	return globalMasterPasswordService.IsInitialized()
}

// SetMasterPasswordWithSignature sets master password with signature validation
func SetMasterPasswordWithSignature(password string) error {
	if globalMasterPasswordService == nil {
		return fmt.Errorf("services not initialized")
	}
	return globalMasterPasswordService.SetMasterPassword(password)
}

// GetMasterPasswordWithSignature gets master password with signature validation
func GetMasterPasswordWithSignature() (string, error) {
	if globalMasterPasswordService == nil {
		return "", fmt.Errorf("services not initialized")
	}
	return globalMasterPasswordService.GetMasterPassword()
}

// ClearMasterPasswordWithSignature clears master password with signature validation
func ClearMasterPasswordWithSignature() error {
	if globalMasterPasswordService == nil {
		return fmt.Errorf("services not initialized")
	}
	return globalMasterPasswordService.ClearMasterPassword()
}

// SetRequirePasswordOnStartupWithSignature sets require password on startup setting with signature validation
func SetRequirePasswordOnStartupWithSignature(require bool) error {
	if globalMasterPasswordService == nil {
		return fmt.Errorf("services not initialized")
	}
	return globalMasterPasswordService.SetRequirePasswordOnStartup(require)
}

// GetRequirePasswordOnStartupWithSignature gets require password on startup setting with signature validation
func GetRequirePasswordOnStartupWithSignature() (bool, error) {
	if globalMasterPasswordService == nil {
		return false, fmt.Errorf("services not initialized")
	}
	return globalMasterPasswordService.GetRequirePasswordOnStartup()
}

// ClearRequirePasswordOnStartupWithSignature clears require password on startup setting with signature validation
func ClearRequirePasswordOnStartupWithSignature() error {
	if globalMasterPasswordService == nil {
		return fmt.Errorf("services not initialized")
	}
	return globalMasterPasswordService.ClearRequirePasswordOnStartup()
}

// SetGlobalSecurityConfigService sets the global security config service
func SetGlobalSecurityConfigService(service *SecurityConfigService) {
	globalSecurityConfigService = service
}

// GetGlobalSecurityConfigService returns the global security config service
func GetGlobalSecurityConfigService() *SecurityConfigService {
	return globalSecurityConfigService
}
