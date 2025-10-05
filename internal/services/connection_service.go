package services

import (
	"fmt"
	"os"
	"ssh-keeper/internal/models"
	"time"
)

// ConnectionService предоставляет методы для работы с подключениями
type ConnectionService struct {
	connections       []models.Connection
	sshConfigService  *SSHConfigService
	encryptionService *EncryptionService
	configPath        string
}

// NewConnectionService создает новый сервис подключений
func NewConnectionService(configPath string) *ConnectionService {
	sshConfigService := NewSSHConfigService(configPath)
	masterPasswordService := NewMasterPasswordService()
	encryptionService := NewEncryptionService(masterPasswordService)

	cs := &ConnectionService{
		connections:       make([]models.Connection, 0),
		sshConfigService:  sshConfigService,
		encryptionService: encryptionService,
		configPath:        configPath,
	}

	// Try to load connections from config file
	if err := cs.LoadConnectionsFromFile(); err != nil {
		// If loading fails, start with empty connections
		cs.connections = make([]models.Connection, 0)
	}

	return cs
}

// LoadConnectionsFromFile loads connections from SSH config file
func (cs *ConnectionService) LoadConnectionsFromFile() error {
	config, err := cs.sshConfigService.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	connections := cs.sshConfigService.ConvertSSHConfigToConnections(config)

	// Decrypt passwords only if encryption service is initialized
	if cs.encryptionService.IsInitialized() {
		for i := range connections {
			if connections[i].HasPassword && connections[i].Password != "" {
				// Проверяем, является ли пароль зашифрованным (должен быть длинным и содержать base64 символы)
				if len(connections[i].Password) >= 20 && isBase64Like(connections[i].Password) {
					decryptedPassword, err := cs.encryptionService.DecryptPassword(connections[i].Password)
					if err != nil {
						return fmt.Errorf("failed to decrypt password for connection %s: %w", connections[i].ID, err)
					}
					connections[i].Password = decryptedPassword
				}
			}
		}
	}

	cs.connections = connections
	return nil
}

// SaveConnectionsToFile saves connections to SSH config file
func (cs *ConnectionService) SaveConnectionsToFile() error {
	// Create a copy of connections for encryption
	connectionsCopy := make([]models.Connection, len(cs.connections))
	copy(connectionsCopy, cs.connections)

	// Encrypt passwords before saving (only if encryption service is initialized)
	if cs.encryptionService.IsInitialized() {
		for i := range connectionsCopy {
			if connectionsCopy[i].HasPassword && connectionsCopy[i].Password != "" && len(connectionsCopy[i].Password) > 0 {
				encryptedPassword, err := cs.encryptionService.EncryptPassword(connectionsCopy[i].Password)
				if err != nil {
					return fmt.Errorf("failed to encrypt password for connection %s: %w", connectionsCopy[i].ID, err)
				}
				connectionsCopy[i].Password = encryptedPassword
			}
		}
	}

	config := cs.sshConfigService.ConvertConnectionsToSSHConfig(connectionsCopy)
	return cs.sshConfigService.SaveConfig(config)
}

// GetAllConnections возвращает все подключения
func (cs *ConnectionService) GetAllConnections() []models.Connection {
	return cs.connections
}

// GetConnectionByID возвращает подключение по ID
func (cs *ConnectionService) GetConnectionByID(id string) *models.Connection {
	for _, conn := range cs.connections {
		if conn.ID == id {
			return &conn
		}
	}
	return nil
}

// AddConnection добавляет новое подключение
func (cs *ConnectionService) AddConnection(conn *models.Connection) error {
	conn.ID = generateID()
	conn.CreatedAt = time.Now()
	conn.UpdatedAt = time.Now()
	cs.connections = append(cs.connections, *conn)

	// Auto-save to file
	return cs.SaveConnectionsToFile()
}

// UpdateConnection обновляет существующее подключение
func (cs *ConnectionService) UpdateConnection(conn *models.Connection) error {
	for i, existing := range cs.connections {
		if existing.ID == conn.ID {
			conn.UpdatedAt = time.Now()
			cs.connections[i] = *conn

			// Auto-save to file
			return cs.SaveConnectionsToFile()
		}
	}
	return fmt.Errorf("connection with ID %s not found", conn.ID)
}

// DeleteConnection удаляет подключение по ID
func (cs *ConnectionService) DeleteConnection(id string) error {
	for i, conn := range cs.connections {
		if conn.ID == id {
			cs.connections = append(cs.connections[:i], cs.connections[i+1:]...)

			// Auto-save to file
			return cs.SaveConnectionsToFile()
		}
	}
	return fmt.Errorf("connection with ID %s not found", id)
}

// ExportConfig exports connections to SSH config file
func (cs *ConnectionService) ExportConfig(exportPath string) error {
	exportService := NewSSHConfigService(exportPath)
	config := cs.sshConfigService.ConvertConnectionsToSSHConfig(cs.connections)
	return exportService.SaveConfig(config)
}

// ImportConfig imports connections from SSH config file
func (cs *ConnectionService) ImportConfig(importPath string) error {
	importService := NewSSHConfigService(importPath)
	config, err := importService.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config from %s: %w", importPath, err)
	}

	importedConnections := importService.ConvertSSHConfigToConnections(config)

	// Decrypt passwords if they are encrypted
	for i := range importedConnections {
		if importedConnections[i].HasPassword && importedConnections[i].Password != "" {
			// Try to decrypt - if it fails, assume it's already plaintext
			if decryptedPassword, err := cs.encryptionService.DecryptPassword(importedConnections[i].Password); err == nil {
				importedConnections[i].Password = decryptedPassword
			}
		}
	}

	// Add imported connections to existing ones
	for _, conn := range importedConnections {
		// Generate new ID to avoid conflicts
		conn.ID = generateID()
		conn.CreatedAt = time.Now()
		conn.UpdatedAt = time.Now()
		cs.connections = append(cs.connections, conn)
	}

	// Save all connections
	return cs.SaveConnectionsToFile()
}

// ExportConfigPlain exports connections to SSH config file without password encryption
func (cs *ConnectionService) ExportConfigPlain(exportPath string) error {
	exportService := NewSSHConfigService(exportPath)

	// Create a copy of connections without encryption
	connectionsCopy := make([]models.Connection, len(cs.connections))
	copy(connectionsCopy, cs.connections)

	// Decrypt passwords before export (only if encryption service is initialized)
	if cs.encryptionService.IsInitialized() {
		for i := range connectionsCopy {
			if connectionsCopy[i].HasPassword && connectionsCopy[i].Password != "" && len(connectionsCopy[i].Password) > 0 {
				// Check if password is encrypted (base64-like)
				if len(connectionsCopy[i].Password) >= 20 && isBase64Like(connectionsCopy[i].Password) {
					decryptedPassword, err := cs.encryptionService.DecryptPassword(connectionsCopy[i].Password)
					if err == nil {
						connectionsCopy[i].Password = decryptedPassword
					}
					// If decryption fails, keep the password as is (might be plaintext already)
				}
			}
		}
	}

	config := cs.sshConfigService.ConvertConnectionsToSSHConfig(connectionsCopy)
	return exportService.SaveConfig(config)
}

// ImportConfigPlain imports connections from SSH config file and encrypts passwords
func (cs *ConnectionService) ImportConfigPlain(importPath string) error {
	importService := NewSSHConfigService(importPath)
	config, err := importService.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config from %s: %w", importPath, err)
	}

	importedConnections := importService.ConvertSSHConfigToConnections(config)

	// Encrypt passwords if encryption service is initialized
	if cs.encryptionService.IsInitialized() {
		for i := range importedConnections {
			if importedConnections[i].HasPassword && importedConnections[i].Password != "" {
				// Check if password is already encrypted
				if len(importedConnections[i].Password) < 20 || !isBase64Like(importedConnections[i].Password) {
					// Password is plaintext, encrypt it
					encryptedPassword, err := cs.encryptionService.EncryptPassword(importedConnections[i].Password)
					if err != nil {
						return fmt.Errorf("failed to encrypt password for connection %s: %w", importedConnections[i].Name, err)
					}
					importedConnections[i].Password = encryptedPassword
				}
			}
		}
	}

	// Check for duplicate connections and add only new ones
	var addedCount int
	var skippedCount int

	for _, conn := range importedConnections {
		// Check if connection already exists
		isDuplicate := false

		// If connection has ID, check by ID
		if conn.ID != "" {
			for _, existingConn := range cs.connections {
				if existingConn.ID == conn.ID {
					isDuplicate = true
					skippedCount++
					break
				}
			}
		} else {
			// If no ID, check by combination of Host+Port+User+Name
			for _, existingConn := range cs.connections {
				if existingConn.Host == conn.Host &&
					existingConn.Port == conn.Port &&
					existingConn.User == conn.User &&
					existingConn.Name == conn.Name {
					isDuplicate = true
					skippedCount++
					break
				}
			}
		}

		if !isDuplicate {
			// Generate new ID only if it's empty
			if conn.ID == "" {
				conn.ID = generateID()
			}
			conn.CreatedAt = time.Now()
			conn.UpdatedAt = time.Now()
			cs.connections = append(cs.connections, conn)
			addedCount++
		}
	}

	// Save all connections only if we added new ones
	if addedCount > 0 {
		return cs.SaveConnectionsToFile()
	}

	// If no new connections were added, return a specific error
	if skippedCount > 0 {
		return fmt.Errorf("all %d connections already exist (duplicates skipped)", skippedCount)
	}

	return nil
}

// GetConfigPath returns the current config file path
func (cs *ConnectionService) GetConfigPath() string {
	return cs.configPath
}

// ReloadConnections reloads connections from file
func (cs *ConnectionService) ReloadConnections() error {
	return cs.LoadConnectionsFromFile()
}

// InitializeWithSampleData initializes the service with sample connections if no config exists
func (cs *ConnectionService) InitializeWithSampleData() error {
	// Check if config file exists
	if _, err := os.Stat(cs.configPath); os.IsNotExist(err) {
		// Create sample connections
		sampleConnections := cs.getSampleConnections()

		// Add all sample connections
		for _, conn := range sampleConnections {
			if err := cs.AddConnection(&conn); err != nil {
				return fmt.Errorf("failed to add sample connection %s: %w", conn.Name, err)
			}
		}
	}
	return nil
}

// getSampleConnections returns sample connections for initial setup
func (cs *ConnectionService) getSampleConnections() []models.Connection {
	return []models.Connection{
		{
			ID:          "2",
			Name:        "Development Server",
			Host:        "dev.example.com",
			Port:        2222,
			User:        "developer",
			KeyPath:     "",
			UseSSHKey:   false, // Только пароль
			HasPassword: true,
			CreatedAt:   time.Now().Add(-15 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-3 * 24 * time.Hour),
		},
		{
			ID:          "3",
			Name:        "Staging Environment",
			Host:        "staging.example.com",
			Port:        22,
			User:        "deploy",
			KeyPath:     "~/.ssh/staging_key",
			UseSSHKey:   true, // Использовать SSH ключ
			HasPassword: false,
			CreatedAt:   time.Now().Add(-10 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-1 * 24 * time.Hour),
		},
		{
			ID:          "4",
			Name:        "Database Server",
			Host:        "db.internal.com",
			Port:        22,
			User:        "dbadmin",
			KeyPath:     "",
			UseSSHKey:   false, // Только пароль
			HasPassword: true,
			CreatedAt:   time.Now().Add(-20 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID:          "5",
			Name:        "Web Server",
			Host:        "web.example.com",
			Port:        22,
			User:        "www-data",
			KeyPath:     "~/.ssh/web_key",
			UseSSHKey:   true, // Использовать SSH ключ
			HasPassword: false,
			CreatedAt:   time.Now().Add(-25 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-2 * 24 * time.Hour),
		},
		{
			ID:          "6",
			Name:        "Backup Server",
			Host:        "backup.example.com",
			Port:        22,
			User:        "backup",
			KeyPath:     "",
			UseSSHKey:   false, // Только пароль
			HasPassword: true,
			CreatedAt:   time.Now().Add(-12 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-1 * 24 * time.Hour),
		},
		{
			ID:          "7",
			Name:        "Monitoring Server",
			Host:        "monitor.internal.com",
			Port:        22,
			User:        "monitor",
			KeyPath:     "~/.ssh/monitor_key",
			UseSSHKey:   true, // Использовать SSH ключ
			HasPassword: false,
			CreatedAt:   time.Now().Add(-8 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-6 * time.Hour),
		},
		{
			ID:          "8",
			Name:        "Local Development",
			Host:        "localhost",
			Port:        22,
			User:        "local",
			KeyPath:     "",
			UseSSHKey:   true, // Использовать дефолтные ключи
			HasPassword: false,
			CreatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-2 * time.Hour),
		},
	}
}

// generateID генерирует простой ID (в реальном приложении используйте uuid)
func generateID() string {
	return time.Now().Format("20060102150405")
}

// isBase64Like проверяет, похож ли пароль на base64 (зашифрованный)
func isBase64Like(password string) bool {
	// Проверяем, что строка содержит только base64 символы и знаки равенства
	for _, char := range password {
		if !((char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') ||
			char == '+' || char == '/' || char == '=') {
			return false
		}
	}
	return true
}
