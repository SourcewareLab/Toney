package errorpopup

import (
	"fmt"
	"strings"

	"github.com/SourcewareLab/Toney/internal/colors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrorPopup struct {
	Width    int
	Height   int
	Message  string
	Title    string
	Location string
}

func NewErrorPopup(w, h int, msg, title, location string) *ErrorPopup {
	return &ErrorPopup{
		Width:    w,
		Height:   h,
		Message:  WrapAndLimit(msg, 20, 3),
		Title:    title,
		Location: location,
	}
}

func (s *ErrorPopup) Init() tea.Cmd { return nil }

func (s *ErrorPopup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s *ErrorPopup) View() string {
	style := lipgloss.NewStyle()

	text := fmt.Sprintf("%s\n\n%s", style.Foreground(colors.ColorPalette().ErrorText).Render(s.Title+" | "+s.Location),

		style.Foreground(colors.ColorPalette().Text).Render(s.Message))

	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colors.ColorPalette().FocusedBorder).
		Background(colors.ColorPalette().ErrorBg).Padding(1, 3).Render(text)
}

// Wraps the text to the given length and also limits no. of lines, adds a "..." line if exceeding.
func WrapAndLimit(s string, maxLen, maxLines int) string {
	var lines []string

	for i := 0; i < len(s); i += maxLen {
		end := i + maxLen
		if end > len(s) {
			end = len(s)
		}
		lines = append(lines, s[i:end])
	}

	if len(lines) > maxLines {
		lines = append(lines[:maxLines], "...")
	}

	return strings.Join(lines, "\n")
}
