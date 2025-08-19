package fileexplorer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	filetree "github.com/SourcewareLab/Toney/internal/fileTree"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	filepopup "github.com/SourcewareLab/Toney/internal/models/filePopup"
	"github.com/SourcewareLab/Toney/internal/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileExplorer struct {
	path          string
	IsFocused     bool
	Width         int
	Height        int
	Vp            viewport.Model
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
		Vp:           viewport.New(w/4-4, h-6),
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
	case messages.FzfSelection:
		m.FindWithRelativePath(msg.Selection)
		m.Refresh()
		return m, m.SelectionChanged(m.CurrentNode)
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

			// fmt.Println(m.CurrentIndex*5, m.Vp.YOffset+m.Vp.Height-5)
			if m.CurrentIndex > m.Vp.YOffset+m.Vp.Height-5 {
				m.Vp.YOffset += 1
			}

			return m, m.SelectionChanged(m.CurrentNode)
		case key.Matches(msg, m.Keymap.Up):
			if m.CurrentIndex <= 0 {
				return m, nil
			}
			m.CurrentIndex -= 1
			m.CurrentNode = m.VisibleNodes[m.CurrentIndex]

			if m.CurrentIndex < m.Vp.YOffset {
				m.Vp.YOffset -= 1
			}

			return m, m.SelectionChanged(m.CurrentNode)
		case key.Matches(msg, m.Keymap.OpenForEdit):
			if m.CurrentNode.IsDirectory {
				m.CurrentNode.IsExpanded = !m.CurrentNode.IsExpanded
				m.VisibleNodes = filetree.FlattenVisibleTree(m.Root)
				return m, nil
			}

			args := config.AppConfig.General.Editor
			args = append(args, strings.TrimSuffix(filepopup.GetPath(m.CurrentNode), "/"))

			c := exec.Command(args[0], args[1:]...)
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
		default:
			var cmd tea.Cmd
			m.Vp, cmd = m.Vp.Update(msg)
			return m, cmd

		}
	}

	return m, nil
}

func (m FileExplorer) View() string {
	style := styles.BorderStyle()
	style = style.Align(lipgloss.Left, lipgloss.Top).Width(m.Width / 4).Height(m.Height - 3)

	if m.IsFocused {
		style = style.BorderForeground(colors.ColorPalette().FocusedBorder)
	} else {
		style = style.BorderForeground(colors.ColorPalette().Border)
	}

	s := filetree.BuildNodeTree(m.Root, "", len(m.Root.Children) == 0, m.CurrentNode)

	m.Vp.SetContent(s)
	m.Vp.Width = m.Width/4 - 4
	m.Vp.Height = m.Height - 6

	return style.Render(m.Vp.View())
}

func (m *FileExplorer) Resize(w int, h int) {
	m.Height = h
	m.Width = w
	m.Vp.Width = w/4 - 4
	m.Vp.Height = h - 6
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

func (m *FileExplorer) FindWithRelativePath(path string) {
	cleanPath := filepath.ToSlash(path)
	parts := strings.Split(cleanPath, "/")

	i := 0
	curr := m.Root
	for {
		if !curr.IsDirectory {
			break
		}

		for _, v := range curr.Children { // Moving across the tree
			if v.Name == parts[i] {
				curr = v
				if curr.IsDirectory {
					curr.IsExpanded = true
				}
				break
			}
		}

		i++
	}

	m.CurrentNode = curr
	m.VisibleNodes = filetree.FlattenVisibleTree(m.Root)

	currPath := filepopup.GetPath(curr)
	for k, v := range m.VisibleNodes {
		if filepopup.GetPath(v) == currPath {
			m.CurrentIndex = k
		}
	}
}
