package widgets

import (
	"github.com/SourcewareLab/Toney/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	Palette theme.Palette
	Width   int
	Title   string
	Right   string
}

func (h Header) View() string {
	b := theme.Borders()
	style := lipgloss.NewStyle().
		Border(b).
		BorderForeground(h.Palette.Border).
		Foreground(h.Palette.Fg).
		Padding(0, 1)

	// Title on left, Right info on right
	left := lipgloss.NewStyle().Foreground(h.Palette.Fg).Render(" " + h.Title)
	right := lipgloss.NewStyle().Foreground(h.Palette.Muted).Render(h.Right + " ")
	row := lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		lipgloss.PlaceHorizontal(h.Width-lipgloss.Width(left), lipgloss.Right, right),
	)
	return style.Width(h.Width).Render(row)
}
