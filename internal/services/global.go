package services

import (
	"fmt"
	"ssh-keeper/internal/models"
)

// Global service instances
var (
	globalConnectionService *ConnectionService
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
