package styles

import (
	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/charmbracelet/lipgloss"
)

var BorderStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Align(lipgloss.Center, lipgloss.Center).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(colors.ColorPalette().Surface1)
