package homemodels

import (
	"toney/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileExplorer struct {
	path      string
	IsFocused bool
	Width     int
	Height    int
}

func (m FileExplorer) Init() tea.Cmd {
	return nil
}

func (m *FileExplorer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m FileExplorer) View() string {
	style := styles.BorderStyle

	if m.IsFocused {
		style.BorderForeground(lipgloss.Color("#bb9af7"))
	}

	w := (m.Width / 4) - 1
	h := (m.Height / 2) - 3

	return style.Width(w).Height(h).MarginTop(3).Render("File Explorer Here")
}
