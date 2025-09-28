package services

import (
	"ssh-keeper/internal/models"
	"time"
)

// ConnectionService предоставляет методы для работы с подключениями
type ConnectionService struct {
	connections []models.Connection
}

// NewConnectionService создает новый сервис подключений
func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		connections: getSampleConnections(),
	}
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
func (cs *ConnectionService) AddConnection(conn *models.Connection) {
	conn.ID = generateID()
	conn.CreatedAt = time.Now()
	conn.UpdatedAt = time.Now()
	cs.connections = append(cs.connections, *conn)
}

// UpdateConnection обновляет существующее подключение
func (cs *ConnectionService) UpdateConnection(conn *models.Connection) bool {
	for i, existing := range cs.connections {
		if existing.ID == conn.ID {
			conn.UpdatedAt = time.Now()
			cs.connections[i] = *conn
			return true
		}
	}
	return false
}

// DeleteConnection удаляет подключение по ID
func (cs *ConnectionService) DeleteConnection(id string) bool {
	for i, conn := range cs.connections {
		if conn.ID == id {
			cs.connections = append(cs.connections[:i], cs.connections[i+1:]...)
			return true
		}
	}
	return false
}

// getSampleConnections возвращает примеры подключений
func getSampleConnections() []models.Connection {
	return []models.Connection{
		{
			ID:          "1",
			Name:        "Production Server",
			Host:        "prod.example.com",
			Port:        22,
			User:        "admin",
			KeyPath:     "~/.ssh/id_rsa",
			HasPassword: false,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-7 * 24 * time.Hour),
		},
		{
			ID:          "2",
			Name:        "Development Server",
			Host:        "dev.example.com",
			Port:        2222,
			User:        "developer",
			KeyPath:     "",
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
			HasPassword: true,
			CreatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-1 * time.Hour),
		},
	}
}

// generateID генерирует простой ID (в реальном приложении используйте uuid)
func generateID() string {
	return time.Now().Format("20060102150405")
}
