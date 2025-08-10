package theme

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

// DefaultDarkPalette returns a tasteful Crush-like dark palette.
func DefaultDarkPalette() Palette {
	return Palette{
		Bg:         lipgloss.Color("#0b0f14"), // near-black
		Fg:         lipgloss.Color("#d3d7de"), // light gray
		Muted:      lipgloss.Color("#7b8794"), // muted gray
		Border:     lipgloss.Color("#2a2f3a"), // subtle border
		Accent:     lipgloss.Color("#4fd1c5"), // teal
		AccentAlt:  lipgloss.Color("#a78bfa"), // purple alt
		SelectedBg: lipgloss.Color("#185b57"), // dark teal bg
		SelectedFg: lipgloss.Color("#e6fffb"), // near-white fg
		Error:      lipgloss.Color("#f87171"),
		Warning:    lipgloss.Color("#fbbf24"),
		Success:    lipgloss.Color("#34d399"),
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
