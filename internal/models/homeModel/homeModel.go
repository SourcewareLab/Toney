package homemodel

import (
	"os"
	"path/filepath"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	filetree "github.com/SourcewareLab/Toney/internal/fileTree"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	fileexplorer "github.com/SourcewareLab/Toney/internal/models/fileExplorer"
	"github.com/SourcewareLab/Toney/internal/models/fzf"
	viewer "github.com/SourcewareLab/Toney/internal/models/viewer"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	Width        int
	Height       int
	ShowFinder   bool
	FocusOn      enums.Splits
	FileExplorer *fileexplorer.FileExplorer
	Viewer       *viewer.Viewer
	Keymap       keymap.HomeKeyMap
	Finder       fzf.FuzzyFinder
	Help         help.Model
}

func NewHome(w int, h int) *HomeModel {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, config.AppConfig.General.NotesDir)
	files, _ := filetree.ListFilesRel(path)

	return &HomeModel{
		Width:        w,
		Height:       h,
		ShowFinder:   false,
		FocusOn:      enums.File,
		FileExplorer: fileexplorer.NewFileExplorer(w, h),
		Viewer:       viewer.NewViewer(w, h),
		Finder:       fzf.NewFzf(files, w, h),
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
	case messages.FzfSelection:
		if msg.Exited {
			m.ShowFinder = false
			return m, nil
		}

		updated, cmd := m.FileExplorer.Update(msg)
		if fe, ok := updated.(*fileexplorer.FileExplorer); ok { // Type matching, cause I cant assign it straightaway
			m.FileExplorer = fe
		}

		updated, cmd = m.Viewer.Update(cmd())
		if v, ok := updated.(*viewer.Viewer); ok { // Type matching, cause I cant assign it straightaway
			m.Viewer = v
		}

		m.ShowFinder = false
		return m, cmd
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
		updated, _ := m.Viewer.Update(msg)
		if v, ok := updated.(*viewer.Viewer); ok {
			m.Viewer = v
		}

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.FocusExplorer):
			m.FocusOn = enums.File
			return m, nil
		case key.Matches(msg, m.Keymap.FocusViewer):
			m.FocusOn = enums.FViewer
			return m, nil
		case key.Matches(msg, m.Keymap.BackToMenu):
			if m.ShowFinder {
				var cmd tea.Cmd
				m.Finder, cmd = m.Finder.Update(msg)
				return m, cmd
			}
			return m, func() tea.Msg {
				return messages.ChangePage{Page: enums.MenuPage}
			}
		case key.Matches(msg, m.Keymap.Finder):
			m.ShowFinder = true
			return m, nil
		default:
			if m.ShowFinder {
				var cmd tea.Cmd
				m.Finder, cmd = m.Finder.Update(msg)
				return m, cmd
			}
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
	if m.ShowFinder {
		return m.Finder.View()
	}

	m.FileExplorer.IsFocused = false
	m.Viewer.IsFocused = false

	bindings := []key.Binding{m.Keymap.FocusExplorer, m.Keymap.FocusViewer, m.Keymap.BackToMenu, m.Keymap.Finder}

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

	// Ensure proper sizing before joining - leave space for viewer's right border
	explorerView := m.FileExplorer.View()

	// Force viewer to use less width to ensure right border is visible
	m.Viewer.Width = m.Width - (m.Width / 4) - 2
	viewerView := m.Viewer.View()

	main := lipgloss.JoinHorizontal(lipgloss.Top, explorerView, viewerView)

	help := lipgloss.NewStyle().PaddingLeft(2).Render(m.Help.View(keymap.NewDynamic(bindings)))

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}
