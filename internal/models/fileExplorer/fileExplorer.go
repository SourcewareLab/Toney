package fileexplorer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/enums"
	filetree "github.com/SourcewareLab/Toney/internal/fileTree"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	filepopup "github.com/SourcewareLab/Toney/internal/models/filePopup"
	"github.com/SourcewareLab/Toney/internal/styles"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileExplorer struct {
	path          string
	IsFocused     bool
	Width         int
	Height        int
	Root          *filetree.Node
	CurrentNode   *filetree.Node
	CurrentIndex  int
	VisibleNodes  []*filetree.Node
	LastSelection string
	Keymap        keymap.ExplorerKeyMap
}

func NewFileExplorer(w int, h int) *FileExplorer {
	root, err := filetree.CreateTree()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return &FileExplorer{
		Width:        w,
		Height:       h,
		Root:         root,
		CurrentNode:  root,
		CurrentIndex: 0,
		VisibleNodes: filetree.FlattenVisibleTree(root),
		Keymap:       keymap.NewExplorerKeyMap(),
	}
}

func (m FileExplorer) Init() tea.Cmd {
	return nil
}

func (m *FileExplorer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.EditorClose:
		m.Refresh()
		return m, m.SelectionChanged(m.CurrentNode)
	case messages.RefreshFileExplorerMsg:
		m.Refresh()
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Down):
			if m.CurrentIndex >= len(m.VisibleNodes)-1 {
				return m, nil
			}
			m.CurrentIndex += 1
			m.CurrentNode = m.VisibleNodes[m.CurrentIndex]
			return m, m.SelectionChanged(m.CurrentNode)
		case key.Matches(msg, m.Keymap.Up):
			if m.CurrentIndex <= 0 {
				return m, nil
			}
			m.CurrentIndex -= 1
			m.CurrentNode = m.VisibleNodes[m.CurrentIndex]
			return m, m.SelectionChanged(m.CurrentNode)
		case key.Matches(msg, m.Keymap.OpenForEdit):
			if m.CurrentNode.IsDirectory {
				m.CurrentNode.IsExpanded = !m.CurrentNode.IsExpanded
				m.VisibleNodes = filetree.FlattenVisibleTree(m.Root)
				return m, nil
			}

			c := exec.Command("nvim", strings.TrimSuffix(filepopup.GetPath(m.CurrentNode), "/"))
			cmd := tea.ExecProcess(c, func(err error) tea.Msg {
				return messages.EditorClose{
					Err: err,
				}
			})
			return m, cmd

		case key.Matches(msg, m.Keymap.CreateFile):
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileCreate,
					Curr: m.CurrentNode,
				}
			}
		case key.Matches(msg, m.Keymap.DeleteFile):
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileDelete,
					Curr: m.CurrentNode,
				}
			}
		case key.Matches(msg, m.Keymap.MoveFile):
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileMove,
					Curr: m.CurrentNode,
				}
			}
		case key.Matches(msg, m.Keymap.RenameFile):
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileRename,
					Curr: m.CurrentNode,
				}
			}

		}
	}

	return m, nil
}

func (m FileExplorer) View() string {
	style := styles.BorderStyle
	style = style.Align(lipgloss.Left, lipgloss.Top)

	if m.IsFocused {
		style = style.BorderForeground(colors.ColorPalette().Lavender)
	}

	w := (m.Width / 4) - 1
	h := m.Height - 3

	s := filetree.BuildNodeTree(m.Root, "", len(m.Root.Children) == 0, m.CurrentNode)

	return style.Width(w).Height(h).MarginTop(1).Render(s)
}

func (m *FileExplorer) Resize(w int, h int) {
	m.Height = h
	m.Width = w
}

func (m *FileExplorer) SelectionChanged(node *filetree.Node) tea.Cmd {
	path := filepopup.GetPath(node)
	if node.IsDirectory || m.LastSelection == path {
		return nil
	}

	m.LastSelection = path

	return func() tea.Msg {
		return messages.ChangeFileMessage{
			Path: path,
		}
	}
}

func (m *FileExplorer) Refresh() {
	newRoot, _ := filetree.CreateTree()

	filepopup.MapExpanded(newRoot, m.Root)

	m.Root = newRoot
	m.VisibleNodes = filetree.FlattenVisibleTree(newRoot)

	idx := -1

	for i, val := range m.VisibleNodes {
		if val.Name == m.CurrentNode.Name && filepopup.GetPath(val) == filepopup.GetPath(m.CurrentNode) {
			idx = i
		}
	}

	if idx == -1 {
		if m.CurrentIndex != 0 {
			idx = m.CurrentIndex - 1
		}
	}

	m.CurrentIndex = idx
	m.CurrentNode = m.VisibleNodes[idx]
}
