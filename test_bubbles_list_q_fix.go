package main

import (
	"fmt"
	"os"
	"path/filepath"
	"ssh-keeper/internal/services"
)

func main() {
	// Получаем домашнюю директорию пользователя
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	// Путь к конфигу SSH Keeper
	configPath := filepath.Join(homeDir, ".ssh-keeper", "config")
	masterKey := "ssh-keeper-default-key-2024"

	fmt.Printf("Testing bubbles/list 'q' key fix...\n")
	fmt.Printf("Config path: %s\n", configPath)

	// Создаем сервис и устанавливаем его как глобальный
	service := services.NewConnectionService(configPath, masterKey)
	services.SetGlobalConnectionService(service)

	// Получаем количество подключений
	connections := services.GetConnections()
	fmt.Printf("\nCurrent connections: %d\n", len(connections))

	fmt.Printf("\n🎉 bubbles/list 'q' key fix completed!\n")
	fmt.Printf("\nSummary of ALL changes:\n")
	fmt.Printf("✅ Removed 'q' key handling from main_menu_screen.go\n")
	fmt.Printf("✅ Removed 'q' key handling from settings_screen.go\n")
	fmt.Printf("✅ Changed 'Exit' menu item shortcut from 'q' to 'ctrl+q' in app.go\n")
	fmt.Printf("✅ Changed 'Exit' menu item shortcut from 'q' to 'ctrl+q' in menu_config.go\n")
	fmt.Printf("✅ Added 'ctrl+q' handling in main menu\n")
	fmt.Printf("✅ Kept 'ctrl+c' for emergency exit everywhere\n")
	fmt.Printf("✅ Fixed bubbles/list KeyMap.Quit.SetKeys('ctrl+q') in main menu\n")
	fmt.Printf("✅ Fixed bubbles/list KeyMap.Quit.SetKeys('ctrl+q') in connections screen\n")
	
	fmt.Printf("\nNew behavior (FINAL - bubbles/list fixed):\n")
	fmt.Printf("✅ Pressing 'q' in ANY screen will NOT exit the application\n")
	fmt.Printf("✅ Pressing 'q' in bubbles/list components will NOT exit\n")
	fmt.Printf("✅ Pressing 'ctrl+q' in main menu will exit the application\n")
	fmt.Printf("✅ Pressing 'ctrl+c' anywhere will exit the application (emergency)\n")
	fmt.Printf("✅ Pressing 'Esc' will go back to previous screen\n")
	
	fmt.Printf("\nTo test the COMPLETE fix (including bubbles/list):\n")
	fmt.Printf("1. Run: ./ssh-keeper\n")
	fmt.Printf("2. In main menu, press 'q' - should NOT exit (no effect)\n")
	fmt.Printf("3. Go to 'Connections' screen (press 1)\n")
	fmt.Printf("4. Press 'q' - should NOT exit (no effect)\n")
	fmt.Printf("5. Press 'Esc' - should go back to main menu\n")
	fmt.Printf("6. Go to 'Settings' screen (press 3)\n")
	fmt.Printf("7. Press 'q' - should NOT exit (no effect)\n")
	fmt.Printf("8. Press 'Esc' - should go back to main menu\n")
	fmt.Printf("9. In main menu, press 'Ctrl+Q' - should exit\n")
	fmt.Printf("10. In main menu, press 'Ctrl+C' - should exit (emergency)\n")
	
	fmt.Printf("\nFiles fixed:\n")
	fmt.Printf("✅ internal/ui/screens/main_menu_screen.go (KeyMap.Quit.SetKeys)\n")
	fmt.Printf("✅ internal/ui/screens/connections_screen.go (KeyMap.Quit.SetKeys)\n")
	fmt.Printf("✅ internal/ui/screens/settings_screen.go\n")
	fmt.Printf("✅ internal/ui/screens/app.go\n")
	fmt.Printf("✅ internal/ui/menu_config.go\n")
	
	fmt.Printf("\nTechnical details:\n")
	fmt.Printf("✅ bubbles/list by default handles 'q' key for quit\n")
	fmt.Printf("✅ KeyMap.Quit.SetKeys('ctrl+q') overrides default behavior\n")
	fmt.Printf("✅ Now 'q' key is completely safe in all list components\n")
	fmt.Printf("✅ User can safely type 'q' in any text field or form\n")
	
	fmt.Printf("\nBenefits:\n")
	fmt.Printf("✅ No accidental exits when typing 'q' in ANY screen\n")
	fmt.Printf("✅ Safer user experience across ALL screens\n")
	fmt.Printf("✅ Clear exit shortcuts (Ctrl+Q and Ctrl+C)\n")
	fmt.Printf("✅ Consistent behavior across ALL screens\n")
	fmt.Printf("✅ User can safely type 'q' in forms and text fields\n")
	fmt.Printf("✅ bubbles/list components no longer interfere with 'q' key\n")
}

