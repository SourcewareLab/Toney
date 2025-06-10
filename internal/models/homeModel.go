package models

import (
	homemodels "toney/internal/models/homeModels"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	Width        int
	Height       int
	FocusOn      homemodels.Splits
	FileExplorer homemodels.FileExplorer
	Viewer       homemodels.Viewer
}

func NewHome() HomeModel {
	return HomeModel{
		FocusOn:      homemodels.FViewer,
		FileExplorer: *homemodels.NewFileExplorer(),
		Viewer:       homemodels.Viewer{},
	}
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "F":
			m.FocusOn = homemodels.File
			return m, nil
		case "V":
			m.FocusOn = homemodels.FViewer
			return m, nil
		default:
			switch m.FocusOn {
			case homemodels.FViewer:
				return m.Viewer.Update(msg)
			case homemodels.File:
				return m.FileExplorer.Update(msg)
			}
		}

	case tea.WindowSizeMsg:
		m.Height = msg.Height - 2
		m.Width = msg.Width - 2

		m.Viewer.Update(msg)
		m.FileExplorer.Update(msg)
	}

	var cmd tea.Cmd

	if m.FileExplorer.IsFocused {
		_, cmd = m.FileExplorer.Update(msg)
	} else if m.Viewer.IsFocused {
		_, cmd = m.Viewer.Update(msg)
	}

	return m, cmd
}

func (m HomeModel) View() string {
	m.FileExplorer.IsFocused = false
	m.Viewer.IsFocused = false

	if m.FocusOn == homemodels.File {
		m.FileExplorer.IsFocused = true
	} else if m.FocusOn == homemodels.FViewer {
		m.Viewer.IsFocused = true
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, m.FileExplorer.View(), m.Viewer.View())
}
