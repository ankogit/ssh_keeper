package models

import (
	"time"
)

// Connection represents an SSH connection configuration
type Connection struct {
	ID          string    `yaml:"id"`
	Name        string    `yaml:"name"`
	Host        string    `yaml:"host"`
	Port        int       `yaml:"port,omitempty"`
	User        string    `yaml:"user"`
	KeyPath     string    `yaml:"key_path,omitempty"`
	HasPassword bool      `yaml:"has_password"`
	CreatedAt   time.Time `yaml:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at"`
}

// NewConnection creates a new connection with default values
func NewConnection(name, host, user string) *Connection {
	return &Connection{
		ID:        "",
		Name:      name,
		Host:      host,
		User:      user,
		Port:      22, // Default SSH port
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
