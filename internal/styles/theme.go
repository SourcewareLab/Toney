package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Palette defines a UI color scheme and spacing settings.
type Palette struct {
	Bg         lipgloss.Color
	Fg         lipgloss.Color
	Muted      lipgloss.Color
	Border     lipgloss.Color
	Accent     lipgloss.Color
	AccentAlt  lipgloss.Color
	SelectedBg lipgloss.Color
	SelectedFg lipgloss.Color
	Error      lipgloss.Color
	Warning    lipgloss.Color
	Success    lipgloss.Color
}

// DefaultDarkPalette returns a standard dark palette.
func DefaultDarkPalette() Palette {
	return Palette{
		Bg:         lipgloss.Color("#1a1a1a"), // dark gray
		Fg:         lipgloss.Color("#ffffff"), // white
		Muted:      lipgloss.Color("#808080"), // gray
		Border:     lipgloss.Color("#444444"), // medium gray
		Accent:     lipgloss.Color("#00aaff"), // blue
		AccentAlt:  lipgloss.Color("#ff6600"), // orange
		SelectedBg: lipgloss.Color("#0066cc"), // darker blue
		SelectedFg: lipgloss.Color("#ffffff"), // white
		Error:      lipgloss.Color("#ff4444"),
		Warning:    lipgloss.Color("#ffaa00"),
		Success:    lipgloss.Color("#00cc44"),
	}
}

// Spacing provides consistent paddings and margins.
type Spacing struct {
	PaddingH int
	PaddingV int
	Gap      int
}

func DefaultSpacing() Spacing { return Spacing{PaddingH: 1, PaddingV: 0, Gap: 1} }

// Borders returns common rounded borders.
func Borders() lipgloss.Border { return lipgloss.RoundedBorder() }
