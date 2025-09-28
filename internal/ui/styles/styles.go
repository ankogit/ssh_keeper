package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

var (
	// Color scheme
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#04B575")
	accentColor    = lipgloss.Color("#F25D94")
	warningColor   = lipgloss.Color("#FFA500")
	errorColor     = lipgloss.Color("#FF6B6B")
	successColor   = lipgloss.Color("#51CF66")

	// Text colors
	textPrimary   = lipgloss.Color("#FFFFFF")
	textSecondary = lipgloss.Color("#B0B0B0")
	textMuted     = lipgloss.Color("#808080")
	textInverse   = lipgloss.Color("#000000")

	// Background colors
	bgPrimary   = lipgloss.Color("#1A1A1A")
	bgSecondary = lipgloss.Color("#2D2D2D")
	bgTertiary  = lipgloss.Color("#404040")
	bgHighlight = lipgloss.Color("#3D3D3D")
)

// Styles for different UI elements
var (
	// Main container
	AppStyle = lipgloss.NewStyle().
			Margin(1, 2).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Background(bgPrimary).
			Foreground(textPrimary)

	// Header
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(bgSecondary).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	// Title
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textPrimary).
			Background(bgPrimary).
			Margin(0, 0, 1, 0)

	// Subtitle
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(textSecondary).
			Background(bgPrimary).
			Margin(0, 0, 2, 0)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(textInverse).
			Background(primaryColor).
			Padding(0, 1).
			Margin(0, 1, 0, 0)

	ButtonHoverStyle = lipgloss.NewStyle().
				Foreground(textInverse).
				Background(secondaryColor).
				Padding(0, 1).
				Margin(0, 1, 0, 0)

	ButtonActiveStyle = lipgloss.NewStyle().
				Foreground(textInverse).
				Background(accentColor).
				Padding(0, 1).
				Margin(0, 1, 0, 0)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Foreground(textPrimary).
			Background(bgSecondary).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(bgSecondary).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(secondaryColor)

	// List styles
	ListStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Background(bgSecondary).
			Padding(0, 1)

	ListItemStyle = lipgloss.NewStyle().
			Foreground(textPrimary).
			Padding(0, 1).
			Margin(0, 0, 0, 0)

	ListItemHoverStyle = lipgloss.NewStyle().
				Foreground(textInverse).
				Background(primaryColor).
				Padding(0, 1).
				Margin(0, 0, 0, 0)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Foreground(textInverse).
				Background(secondaryColor).
				Padding(0, 1).
				Margin(0, 0, 0, 0)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			Italic(true)

	// Search styles
	SearchStyle = lipgloss.NewStyle().
			Foreground(textPrimary).
			Background(bgTertiary).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor)

	SearchFocusedStyle = lipgloss.NewStyle().
				Foreground(textPrimary).
				Background(bgTertiary).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentColor)
)

// GetTerminalProfile returns the current terminal color profile
func GetTerminalProfile() termenv.Profile {
	return termenv.ColorProfile()
}

// IsDarkTheme checks if the terminal is using a dark theme
func IsDarkTheme() bool {
	profile := GetTerminalProfile()
	return profile == termenv.ANSI256 || profile == termenv.TrueColor
}

// GetAdaptiveColor returns a color that adapts to the terminal theme
func GetAdaptiveColor(light, dark lipgloss.Color) lipgloss.Color {
	if IsDarkTheme() {
		return dark
	}
	return light
}

// CreateGradient creates a gradient effect for text
func CreateGradient(text string, start, end colorful.Color) string {
	steps := len(text)
	if steps == 0 {
		return text
	}

	var result string
	for i, char := range text {
		ratio := float64(i) / float64(steps-1)
		color := start.BlendRgb(end, ratio)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color.Hex()))
		result += style.Render(string(char))
	}
	return result
}

// RenderLogo creates a beautiful logo for the application
func RenderLogo() string {
	logo := "SSH Keeper"
	startColor := colorful.Color{R: 0.49, G: 0.34, B: 0.96} // Primary color
	endColor := colorful.Color{R: 0.02, G: 0.71, B: 0.46}   // Secondary color

	return CreateGradient(logo, startColor, endColor)
}

// RenderHeader creates a beautiful header for the application
func RenderHeader(title, subtitle string) string {
	logo := RenderLogo()
	titleStyled := TitleStyle.Render(title)
	subtitleStyled := SubtitleStyle.Render(subtitle)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		titleStyled,
		subtitleStyled,
	)
}

// PrimaryColor returns the primary color
func PrimaryColor() lipgloss.Color {
	return primaryColor
}

// SecondaryColor returns the secondary color
func SecondaryColor() lipgloss.Color {
	return secondaryColor
}
