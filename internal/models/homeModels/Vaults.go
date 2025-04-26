package homemodels

import (
	"toney/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Vault struct {
	IsFocused bool
	Height    int
	Width     int
}

func (m Vault) Init() tea.Cmd {
	return nil
}

func (m *Vault) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Vault) View() string {
	style := styles.BorderStyle

	if m.IsFocused {
		style.BorderForeground(lipgloss.Color("#bb9af7"))
	}

	w := (m.Width / 4) - 1
	h := (m.Height / 2) - 1

	return style.Width(w).Height(h).Render("Vaults Here")
}
