package homemodels

import (
	"fmt"

	"toney/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Task struct {
	IsFocused bool
	Height    int
	Width     int
}

func (m Task) Init() tea.Cmd {
	return nil
}

func (m *Task) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Task) View() string {
	style := styles.BorderStyle

	if m.IsFocused {
		style.BorderForeground(lipgloss.Color("#bb9af7"))
	}

	w := (((m.Width) / 4) * 3)
	h := m.Height - 2
	return style.Width(w).Height(h).MarginTop(3).Render(fmt.Sprintf("Tasks Here; Height: %d", h))
}
