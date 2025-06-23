package homemodel

import (
	"github.com/NucleoFusion/Toney/internal/enums"
	"github.com/NucleoFusion/Toney/internal/messages"
	viewer "github.com/NucleoFusion/Toney/internal/models/Viewer"
	fileexplorer "github.com/NucleoFusion/Toney/internal/models/fileExplorer"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	Width        int
	Height       int
	FocusOn      enums.Splits
	FileExplorer *fileexplorer.FileExplorer
	Viewer       *viewer.Viewer
}

func NewHome(w int, h int) *HomeModel {
	return &HomeModel{
		Width:        w,
		Height:       h,
		FocusOn:      enums.File,
		FileExplorer: fileexplorer.NewFileExplorer(w, h),
		Viewer:       viewer.NewViewer(w, h),
	}
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.RefreshView:
		return m.Viewer.Update(msg)
	case messages.ShowPopupMessage:
		return m, func() tea.Msg {
			return msg
		}
	case messages.ChangeFileMessage:
		return m.Viewer.Update(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "F":
			m.FocusOn = enums.File
			return m, nil
		case "V":
			m.FocusOn = enums.FViewer
			return m, nil
		default:
			switch m.FocusOn {
			case enums.FViewer:
				return m.Viewer.Update(msg)
			case enums.File:
				return m.FileExplorer.Update(msg)
			}
		}
	}

	return m, nil
}

func (m HomeModel) View() string {
	m.FileExplorer.IsFocused = false
	m.Viewer.IsFocused = false

	if m.FocusOn == enums.File {
		m.FileExplorer.IsFocused = true
	} else if m.FocusOn == enums.FViewer {
		m.Viewer.IsFocused = true
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, m.FileExplorer.View(), m.Viewer.View())
}
