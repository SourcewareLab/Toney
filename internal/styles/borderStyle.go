package styles

import "github.com/charmbracelet/lipgloss"

var BorderStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Align(lipgloss.Center, lipgloss.Center).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#45475a"))
