package widgets

import (
	"github.com/SourcewareLab/Toney/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type KeyHint struct {
	Key  string
	Desc string
}

type Footer struct {
	Palette theme.Palette
	Width   int
	Hints   []KeyHint
}

func (f Footer) View() string {
	b := theme.Borders()
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
