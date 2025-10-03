package main

import (
	"fmt"
	"os"
	"path/filepath"
	"ssh-keeper/internal/models"
	"ssh-keeper/internal/services"
	"time"
)

func main() {
	// Создаем временную директорию для тестов
	tempDir := "/tmp/ssh-keeper-test"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test_config")
	masterKey := "test-master-key-2024"

	// Создаем сервис
	service := services.NewConnectionService(configPath, masterKey)

	// Загружаем примеры подключений
	fmt.Println("Loading sample connections...")
	sampleConnections := getSampleConnections()

	// Добавляем все примеры подключений
	for _, conn := range sampleConnections {
		err := service.AddConnection(&conn)
		if err != nil {
			fmt.Printf("Error adding connection %s: %v\n", conn.Name, err)
			return
		}
	}

	// Проверяем, что подключение добавлено
	connections := service.GetAllConnections()
	fmt.Printf("Total connections: %d\n", len(connections))

	// Проверяем содержимое файла конфига
	fmt.Println("\nConfig file contents:")
	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}
	fmt.Println(string(content))

	// Тестируем импорт
	fmt.Println("\nTesting import...")
	importPath := filepath.Join(tempDir, "import_config")
	err = os.WriteFile(importPath, []byte(`# Test import config
Host imported-server
    HostName 10.0.0.1
    Port 22
    User admin
    Password testpass
`), 0644)
	if err != nil {
		fmt.Printf("Error creating import file: %v\n", err)
		return
	}

	err = service.ImportConfig(importPath)
	if err != nil {
		fmt.Printf("Error importing config: %v\n", err)
		return
	}

	connections = service.GetAllConnections()
	fmt.Printf("Total connections after import: %d\n", len(connections))

	fmt.Println("\nTest completed successfully!")
}

// getSampleConnections возвращает примеры подключений
func getSampleConnections() []models.Connection {
	return []models.Connection{
		{
			ID:          "1",
			Name:        "Test Server (Password Auth)",
			Host:        "213.165.35.209",
			Port:        22,
			User:        "root",
			KeyPath:     "",
			HasPassword: true,
			Password:    "yrw6hs1IsLxt",
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-7 * 24 * time.Hour),
		},
		{
			ID:          "1b",
			Name:        "Test Server 2 (SSH Key)",
			Host:        "95.216.195.202",
			Port:        22,
			User:        "root",
			KeyPath:     "", // Пустой путь означает использование дефолтных ключей
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
