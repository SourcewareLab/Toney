package homemodel

import (
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	viewer "github.com/SourcewareLab/Toney/internal/models/Viewer"
	fileexplorer "github.com/SourcewareLab/Toney/internal/models/fileExplorer"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	Width        int
	Height       int
	FocusOn      enums.Splits
	FileExplorer *fileexplorer.FileExplorer
	Viewer       *viewer.Viewer
	Keymap       keymap.HomeKeyMap
	Help         help.Model
}

func NewHome(w int, h int) *HomeModel {
	return &HomeModel{
		Width:        w,
		Height:       h,
		FocusOn:      enums.File,
		FileExplorer: fileexplorer.NewFileExplorer(w, h),
		Viewer:       viewer.NewViewer(w, h),
		Keymap:       keymap.NewHomeKeyMap(),
		Help:         help.New(),
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
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		m.FileExplorer.Resize(msg.Width, msg.Height)
		m.Viewer = viewer.NewViewer(msg.Width, m.Height)

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.FocusExplorer):
			m.FocusOn = enums.File
			return m, nil
		case key.Matches(msg, m.Keymap.FocusViewer):
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

	bindings := []key.Binding{m.Keymap.FocusExplorer, m.Keymap.FocusViewer}

	if m.FocusOn == enums.File {
		m.FileExplorer.IsFocused = true
		bindings = append(bindings, m.FileExplorer.Keymap.CreateFile,
			m.FileExplorer.Keymap.RenameFile,
			m.FileExplorer.Keymap.MoveFile,
			m.FileExplorer.Keymap.DeleteFile,
			m.FileExplorer.Keymap.OpenForEdit,
		)
	} else if m.FocusOn == enums.FViewer {
		m.Viewer.IsFocused = true
		bindings = append(bindings,
			m.Viewer.Keymap.ScrollUp,
			m.Viewer.Keymap.ScrollDown,
		)
	}

	main := lipgloss.JoinHorizontal(lipgloss.Top, m.FileExplorer.View(), m.Viewer.View())

	help := lipgloss.NewStyle().PaddingLeft(2).Render(m.Help.View(keymap.NewDynamic(bindings)))

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}
