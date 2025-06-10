package homemodels

import (
	"fmt"

	"toney/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Viewer struct {
	IsFocused bool
	Height    int
	Width     int
}

func (m Viewer) Init() tea.Cmd {
	return nil
}

func (m *Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	}

	return m, nil
}

func (m Viewer) View() string {
	style := styles.BorderStyle

	if m.IsFocused {
		style = style.BorderForeground(lipgloss.Color("#bb9af7"))
	}

	w := (((m.Width)/4)*3 - 1)
	h := m.Height - 2
	return style.Width(w).Height(h).MarginTop(1).Render(fmt.Sprintf("Viewers Here; Height: %d", h))
}
