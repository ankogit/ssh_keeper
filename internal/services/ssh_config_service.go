package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ssh-keeper/internal/models"
)

// SSHConfigService handles SSH configuration file operations
type SSHConfigService struct {
	configPath string
}

// NewSSHConfigService creates a new SSH config service
func NewSSHConfigService(configPath string) *SSHConfigService {
	return &SSHConfigService{
		configPath: configPath,
	}
}

// LoadConfig loads SSH configuration from file
func (scs *SSHConfigService) LoadConfig() (*models.SSHConfig, error) {
	if _, err := os.Stat(scs.configPath); os.IsNotExist(err) {
		// Create default config if file doesn't exist
		return models.NewSSHConfig(), nil
	}

	file, err := os.Open(scs.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := models.NewSSHConfig()
	scanner := bufio.NewScanner(file)

	var currentHost *models.SSHConfigHost
	var inHostBlock bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Check for Host directive
		if strings.HasPrefix(strings.ToLower(line), "host ") {
			// Save previous host if exists
			if currentHost != nil && inHostBlock {
				config.AddHost(*currentHost)
			}

			// Parse host patterns
			hostPatterns := strings.Fields(line)[1:] // Skip "Host" keyword
			currentHost = &models.SSHConfigHost{
				Host: hostPatterns,
			}
			inHostBlock = true
			continue
		}

		// If we're not in a host block, treat as global setting
		if !inHostBlock {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				config.GlobalSettings[strings.ToLower(parts[0])] = parts[1]
			}
			continue
		}

		// Parse host-specific settings
		if currentHost != nil {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.ToLower(parts[0])
			value := parts[1]

			switch key {
			case "name":
				currentHost.Name = value
			case "hostname":
				currentHost.HostName = value
			case "port":
				if port, err := strconv.Atoi(value); err == nil {
					currentHost.Port = port
				}
			case "user":
				currentHost.User = value
			case "identityfile":
				currentHost.IdentityFile = value
			case "usesshkey":
				currentHost.UseSSHKey = strings.ToLower(value) == "true" || value == "yes" || value == "1"
			case "password":
				currentHost.Password = value
			case "stricthostkeychecking":
				currentHost.StrictHostKeyChecking = value
			case "userknownhostsfile":
				currentHost.UserKnownHostsFile = value
			case "serveraliveinterval":
				if interval, err := strconv.Atoi(value); err == nil {
					currentHost.ServerAliveInterval = interval
				}
			case "serveralivecountmax":
				if count, err := strconv.Atoi(value); err == nil {
					currentHost.ServerAliveCountMax = count
				}
			case "ID":
				currentHost.ID = value
			case "CreatedAt":
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					currentHost.CreatedAt = t
				}
			case "UpdatedAt":
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					currentHost.UpdatedAt = t
				}
			}
		}
	}

	// Don't forget the last host
	if currentHost != nil && inHostBlock {
		config.AddHost(*currentHost)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves SSH configuration to file
func (scs *SSHConfigService) SaveConfig(config *models.SSHConfig) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(scs.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	file, err := os.Create(scs.configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write header comment
	fmt.Fprintf(writer, "# SSH Keeper Configuration File\n")
	fmt.Fprintf(writer, "# Generated on %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(writer, "# Version: %s\n\n", config.Version)

	// Write global settings
	if len(config.GlobalSettings) > 0 {
		fmt.Fprintf(writer, "# Global Settings\n")
		for key, value := range config.GlobalSettings {
			fmt.Fprintf(writer, "%s %s\n", key, value)
		}
		fmt.Fprintf(writer, "\n")
	}

	// Write host configurations
	for _, host := range config.Hosts {
		fmt.Fprintf(writer, "Host %s\n", strings.Join(host.Host, " "))

		// Write name if available
		if host.Name != "" {
			fmt.Fprintf(writer, "    Name %s\n", host.Name)
		}

		if host.HostName != "" {
			fmt.Fprintf(writer, "    HostName %s\n", host.HostName)
		}
		if host.Port != 0 && host.Port != 22 {
			fmt.Fprintf(writer, "    Port %d\n", host.Port)
		}
		if host.User != "" {
			fmt.Fprintf(writer, "    User %s\n", host.User)
		}
		if host.IdentityFile != "" {
			fmt.Fprintf(writer, "    IdentityFile %s\n", host.IdentityFile)
		}
		if host.UseSSHKey {
			fmt.Fprintf(writer, "    UseSSHKey true\n")
		}
		if host.Password != "" {
			fmt.Fprintf(writer, "    Password %s\n", host.Password)
		}
		if host.StrictHostKeyChecking != "" {
			fmt.Fprintf(writer, "    StrictHostKeyChecking %s\n", host.StrictHostKeyChecking)
		}
		if host.UserKnownHostsFile != "" {
			fmt.Fprintf(writer, "    UserKnownHostsFile %s\n", host.UserKnownHostsFile)
		}
		if host.ServerAliveInterval != 0 {
			fmt.Fprintf(writer, "    ServerAliveInterval %d\n", host.ServerAliveInterval)
		}
		if host.ServerAliveCountMax != 0 {
			fmt.Fprintf(writer, "    ServerAliveCountMax %d\n", host.ServerAliveCountMax)
		}

		// SSH Keeper metadata
		if host.ID != "" {
			fmt.Fprintf(writer, "    ID %s\n", host.ID)
		}
		if host.CreatedAt != (time.Time{}) {
			fmt.Fprintf(writer, "    CreatedAt %s\n", host.CreatedAt.Format(time.RFC3339))
		}
		if host.UpdatedAt != (time.Time{}) {
			fmt.Fprintf(writer, "    UpdatedAt %s\n", host.UpdatedAt.Format(time.RFC3339))
		}

		fmt.Fprintf(writer, "\n")
	}

	return nil
}

// ConvertConnectionsToSSHConfig converts ConnectionService connections to SSH config
func (scs *SSHConfigService) ConvertConnectionsToSSHConfig(connections []models.Connection) *models.SSHConfig {
	config := models.NewSSHConfig()

	for _, conn := range connections {
		host := &models.SSHConfigHost{}
		host.ConvertFromConnection(&conn)
		config.AddHost(*host)
	}

	return config
}

// ConvertSSHConfigToConnections converts SSH config to ConnectionService connections
func (scs *SSHConfigService) ConvertSSHConfigToConnections(config *models.SSHConfig) []models.Connection {
	connections := make([]models.Connection, 0, len(config.Hosts))

	for _, host := range config.Hosts {
		conn := host.ConvertToConnection()
		connections = append(connections, *conn)
	}

	return connections
}
