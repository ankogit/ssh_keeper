package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"ssh-keeper/internal/models"
)

// ConnectionService handles business logic for connections
type ConnectionService struct {
	connections []*models.Connection
}

// NewConnectionService creates a new connection service
func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		connections: make([]*models.Connection, 0),
	}
}

// CreateConnection creates a new connection with business logic
func (s *ConnectionService) CreateConnection(name, host, user string) *models.Connection {
	conn := models.NewConnection(name, host, user)
	conn.ID = s.generateID()
	conn.UpdatedAt = time.Now()

	s.connections = append(s.connections, conn)
	return conn
}

// SetPassword marks that this connection has a password stored
func (s *ConnectionService) SetPassword(conn *models.Connection) {
	conn.HasPassword = true
	conn.UpdatedAt = time.Now()
}

// SetKeyPath sets the SSH key path for this connection
func (s *ConnectionService) SetKeyPath(conn *models.Connection, keyPath string) {
	conn.KeyPath = keyPath
	conn.UpdatedAt = time.Now()
}

// SetPort sets the SSH port for this connection
func (s *ConnectionService) SetPort(conn *models.Connection, port int) {
	conn.Port = port
	conn.UpdatedAt = time.Now()
}

// GetAllConnections returns all connections
func (s *ConnectionService) GetAllConnections() []*models.Connection {
	return s.connections
}

// GetConnectionByID returns a connection by ID
func (s *ConnectionService) GetConnectionByID(id string) *models.Connection {
	for _, conn := range s.connections {
		if conn.ID == id {
			return conn
		}
	}
	return nil
}

// UpdateConnection updates an existing connection
func (s *ConnectionService) UpdateConnection(conn *models.Connection) bool {
	for i, existing := range s.connections {
		if existing.ID == conn.ID {
			conn.UpdatedAt = time.Now()
			s.connections[i] = conn
			return true
		}
	}
	return false
}

// DeleteConnection deletes a connection by ID
func (s *ConnectionService) DeleteConnection(id string) bool {
	for i, conn := range s.connections {
		if conn.ID == id {
			s.connections = append(s.connections[:i], s.connections[i+1:]...)
			return true
		}
	}
	return false
}

// generateID generates a unique ID for the connection
func (s *ConnectionService) generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
