package styles

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type Header struct {
	Palette Palette
	Width   int
	Title   string
	Right   string
}

func (h Header) View() string {
	b := Borders()
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

type KeyHint struct {
	Key  string
	Desc string
}

type Footer struct {
	Palette Palette
	Width   int
	Hints   []KeyHint
}

func (f Footer) View() string {
	b := Borders()
	key := lipgloss.NewStyle().Foreground(f.Palette.SelectedFg).Background(f.Palette.SelectedBg).Bold(true).Padding(0, 1)
	desc := lipgloss.NewStyle().Foreground(f.Palette.Fg)
	segments := make([]string, len(f.Hints))
	for i, h := range f.Hints {
		segments[i] = key.Render(h.Key) + " " + desc.Render(h.Desc)
	}
	content := strings.Join(segments, " • ")

	bar := lipgloss.NewStyle().
		Border(b).
		BorderForeground(f.Palette.Border).
		Foreground(f.Palette.Fg).
		Padding(0, 1).
		Width(f.Width)
	return bar.Render(content)
}

type HelpSection struct {
	Title string
	Rows  []string // lines of content already formatted
}

type HelpOverlay struct {
	Palette  Palette
	Width    int
	Height   int
	Title    string
	Sections []HelpSection
}

func (h HelpOverlay) View() string {
	pal := h.Palette
	b := Borders()
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
