package widgets

import (
	"github.com/SourcewareLab/Toney/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type HelpSection struct {
	Title string
	Rows  []string // lines of content already formatted
}

type HelpOverlay struct {
	Palette  theme.Palette
	Width    int
	Height   int
	Title    string
	Sections []HelpSection
}

func (h HelpOverlay) View() string {
	pal := h.Palette
	b := theme.Borders()
	// Build section content
	parts := make([]string, 0, len(h.Sections)+1)
	title := lipgloss.NewStyle().Foreground(pal.Fg).Bold(true).Render(h.Title)
	parts = append(parts, title)
	for _, s := range h.Sections {
		st := lipgloss.NewStyle().Foreground(pal.Accent).Bold(true).Render(s.Title)
		body := lipgloss.NewStyle().Foreground(pal.Fg).Render(strings.Join(s.Rows, "\n"))
		parts = append(parts, st, body)
	}
	content := strings.Join(parts, "\n\n")

	box := lipgloss.NewStyle().
		Border(b).
		BorderForeground(pal.Border).
		Foreground(pal.Fg).
		Padding(1, 2)

	// Constrain box size within viewport
	maxW := h.Width - 6
	if maxW < 20 {
		maxW = h.Width - 2
	}
	rendered := box.Width(maxW).Render(content)

	// Center the box in the screen
	return lipgloss.Place(h.Width, h.Height, lipgloss.Center, lipgloss.Center, rendered)
}
