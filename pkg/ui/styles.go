package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Styles
	CorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")) // Green
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")) // Red
	UntypedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")) // Grey
	CursorStyle  = lipgloss.NewStyle().Underline(true)                       // Underline
	TitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).MarginBottom(1)
	ResultsStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#888888"))
)
