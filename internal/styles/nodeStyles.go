package styles

import (
	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/charmbracelet/lipgloss"
)

var CurrentNodeStyle = lipgloss.NewStyle().Background(colors.ColorPalette().Lavender).Foreground(colors.ColorPalette().Base)
