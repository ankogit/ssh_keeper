package models

import (
	"fmt"
	"time"
)

// SSHConfigHost represents a single host configuration in SSH config format
type SSHConfigHost struct {
	// Host patterns (can be wildcards)
	Host []string `yaml:"host"`

	// Connection settings
	Name     string `yaml:"name,omitempty"` // Display name for the connection
	HostName string `yaml:"hostname"`
	Port     int    `yaml:"port,omitempty"`
	User     string `yaml:"user"`

	// Authentication
	IdentityFile string `yaml:"identityfile,omitempty"`
	UseSSHKey    bool   `yaml:"usesshkey,omitempty"` // Whether to use SSH key authentication
	Password     string `yaml:"password,omitempty"`  // Will be encrypted

	// Additional SSH options
	StrictHostKeyChecking string `yaml:"strictHostKeyChecking,omitempty"`
	UserKnownHostsFile    string `yaml:"userKnownHostsFile,omitempty"`
	ServerAliveInterval   int    `yaml:"serverAliveInterval,omitempty"`
	ServerAliveCountMax   int    `yaml:"serverAliveCountMax,omitempty"`

	// SSH Keeper specific metadata
	ID        string    `yaml:"id,omitempty"`
	CreatedAt time.Time `yaml:"created_at,omitempty"`
	UpdatedAt time.Time `yaml:"updated_at,omitempty"`
}

// SSHConfig represents the complete SSH configuration file
type SSHConfig struct {
	// Global settings
	GlobalSettings map[string]string `yaml:"global_settings,omitempty"`

	// Host configurations
	Hosts []SSHConfigHost `yaml:"hosts"`

	// SSH Keeper metadata
	Version   string    `yaml:"version"`
	CreatedAt time.Time `yaml:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at"`
}

// NewSSHConfig creates a new SSH configuration
func NewSSHConfig() *SSHConfig {
	return &SSHConfig{
		GlobalSettings: make(map[string]string),
		Hosts:          make([]SSHConfigHost, 0),
		Version:        "1.0",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// AddHost adds a new host configuration
func (sc *SSHConfig) AddHost(host SSHConfigHost) {
	host.CreatedAt = time.Now()
	host.UpdatedAt = time.Now()
	sc.Hosts = append(sc.Hosts, host)
	sc.UpdatedAt = time.Now()
}

// FindHostByID finds a host by SSH Keeper ID
func (sc *SSHConfig) FindHostByID(id string) *SSHConfigHost {
	for i := range sc.Hosts {
		if sc.Hosts[i].ID == id {
			return &sc.Hosts[i]
		}
	}
	return nil
}

// UpdateHost updates an existing host configuration
func (sc *SSHConfig) UpdateHost(host SSHConfigHost) bool {
	for i := range sc.Hosts {
		if sc.Hosts[i].ID == host.ID {
			host.UpdatedAt = time.Now()
			sc.Hosts[i] = host
			sc.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// RemoveHost removes a host by ID
func (sc *SSHConfig) RemoveHost(id string) bool {
	for i := range sc.Hosts {
		if sc.Hosts[i].ID == id {
			sc.Hosts = append(sc.Hosts[:i], sc.Hosts[i+1:]...)
			sc.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// ConvertToConnection converts SSHConfigHost to Connection model
func (sh *SSHConfigHost) ConvertToConnection() *Connection {
	conn := &Connection{
		ID:          sh.ID,
		Name:        sh.Name,
		Host:        sh.HostName,
		Port:        sh.Port,
		User:        sh.User,
		KeyPath:     sh.IdentityFile,
		UseSSHKey:   sh.UseSSHKey,
		Password:    sh.Password,
		HasPassword: sh.Password != "",
		CreatedAt:   sh.CreatedAt,
		UpdatedAt:   sh.UpdatedAt,
	}

	// Set default port if not specified
	if conn.Port == 0 {
		conn.Port = 22
	}

	// Set fallback name if empty
	if conn.Name == "" {
		conn.Name = fmt.Sprintf("%s@%s", conn.User, conn.Host)
	}

	return conn
}

// ConvertFromConnection converts Connection to SSHConfigHost
func (sh *SSHConfigHost) ConvertFromConnection(conn *Connection) {
	sh.ID = conn.ID
	sh.Name = conn.Name
	sh.HostName = conn.Host
	sh.Port = conn.Port
	sh.User = conn.User
	sh.IdentityFile = conn.KeyPath
	sh.UseSSHKey = conn.UseSSHKey
	sh.Password = conn.Password
	sh.CreatedAt = conn.CreatedAt
	sh.UpdatedAt = conn.UpdatedAt

	// Set default SSH options
	if sh.StrictHostKeyChecking == "" {
		sh.StrictHostKeyChecking = "ask"
	}
	if sh.ServerAliveInterval == 0 {
		sh.ServerAliveInterval = 60
	}
	if sh.ServerAliveCountMax == 0 {
		sh.ServerAliveCountMax = 3
	}

	// Set host pattern based on hostname
	if len(sh.Host) == 0 {
		sh.Host = []string{sh.HostName}
	}
}
