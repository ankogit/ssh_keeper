package main

import (
	"fmt"
	"os"

	"ssh-keeper/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func main() {
	// Set up terminal environment
	lipgloss.SetColorProfile(termenv.ColorProfile())

	// Create main menu screen
	screen := screens.NewMainMenuScreen()

	// Create tea program
	p := tea.NewProgram(screen, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
