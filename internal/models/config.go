package models

import (
	"time"
)

// Config represents the application configuration
type Config struct {
	MasterKeyTimeout time.Duration `yaml:"master_key_timeout"`
	SSHPath          string        `yaml:"ssh_path"`
	ExportFormat     string        `yaml:"export_format"`
	DefaultPort      int           `yaml:"default_port"`
	Theme            string        `yaml:"theme"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		MasterKeyTimeout: time.Hour, // 1 hour timeout
		SSHPath:          "ssh",     // Use system ssh
		ExportFormat:     "openssh", // OpenSSH config format
		DefaultPort:      22,        // Default SSH port
		Theme:            "default", // Default theme
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.MasterKeyTimeout <= 0 {
		c.MasterKeyTimeout = time.Hour
	}
	if c.SSHPath == "" {
		c.SSHPath = "ssh"
	}
	if c.ExportFormat == "" {
		c.ExportFormat = "openssh"
	}
	if c.DefaultPort <= 0 {
		c.DefaultPort = 22
	}
	if c.Theme == "" {
		c.Theme = "default"
	}
	return nil
}
