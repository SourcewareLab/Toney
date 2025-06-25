package fileexplorer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SourcewareLab/Toney/internal/enums"
	filetree "github.com/SourcewareLab/Toney/internal/fileTree"
	"github.com/SourcewareLab/Toney/internal/messages"
	filepopup "github.com/SourcewareLab/Toney/internal/models/filePopup"
	"github.com/SourcewareLab/Toney/internal/styles"

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
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "down":
			if m.CurrentIndex >= len(m.VisibleNodes)-1 {
				return m, nil
			}
			m.CurrentIndex += 1
			m.CurrentNode = m.VisibleNodes[m.CurrentIndex]
			return m, m.SelectionChanged(m.CurrentNode)
		case "up":
			if m.CurrentIndex <= 0 {
				return m, nil
			}
			m.CurrentIndex -= 1
			m.CurrentNode = m.VisibleNodes[m.CurrentIndex]
			return m, m.SelectionChanged(m.CurrentNode)
		case "enter":
			// It's possible for currentNode to be nil if the file explorer is empty
			// This check prevents the application from crashing if the user presses enter in that state
			if m.CurrentNode == nil {
				return m, nil
			}

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

		case "c":
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileCreate,
					Curr: m.CurrentNode,
				}
			}
		case "d":
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileDelete,
					Curr: m.CurrentNode,
				}
			}

		case "m":
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileMove,
					Curr: m.CurrentNode,
				}
			}
		case "r":
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
		style = style.BorderForeground(lipgloss.Color("#b4befe"))
	}

	w := (m.Width / 4) - 1
	h := m.Height - 2

	var s string
	// A check is added to endure m.Root is not nil before trying to build thre tree.
	// This prevents a panic if the View is rendered after the initial CreateTree failed
	// or if the tree becomes empty after a refresh
	if m.Root != nil {
		s = filetree.BuildNodeTree(m.Root, "", len(m.Root.Children) == 0, m.CurrentNode)
	}

	return style.Width(w).Height(h).MarginTop(1).Render(s)
}

func (m *FileExplorer) Resize(w int, h int) {
	m.Height = h
	m.Width = w
}

func (m *FileExplorer) SelectionChanged(node *filetree.Node) tea.Cmd {
	// Added a guard clause to if this function is ever called with a nil node.
	// Makes function more robust against unexpected states
	if node == nil {
		return nil
	}

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
	newRoot, err := filetree.CreateTree()

	// From last state this function was a bit of a *respectful mess* lol. So we're going to try to clean it up a little bit at a time.

	// The OG code ignored the error from CreateTree, (e.g due to a permissions issue), NewRoot would be nil. We now handle the error by aborting the refresh, keeping app stable.

	if err != nil {
		// Log or handlme the error appropriately instead of ignorin g it
		// FOr now, we'll just abort the refresh
		return
	}

	// Check ensures we don't try to map expanded states from a nil m.Root, which could happen in rare cases as such. It's small change but added safety.

	if m.Root != nil {
		filepopup.MapExpanded(newRoot, m.Root)
	}

	m.Root = newRoot
	m.VisibleNodes = filetree.FlattenVisibleTree(newRoot)

	// After a refresh (e.g deleting the last file in a dir), the list of visible nodes might be empty. Graceful handling added. Resetting the selection from accessing an empty slice.
	if len(m.VisibleNodes) == 0 {
		m.CurrentIndex = 0
		m.CurrentNode = nil
		return
	}

	idx := -1

	// we only search for the old node if one was actually selected. This avoids nil pointer deference on m.CurrentNode if it was nil before the refresh.
	if m.CurrentNode != nil {
		currentPath := filepopup.GetPath(m.CurrentNode)
		for i, val := range m.VisibleNodes {
			if val.Name == m.CurrentNode.Name && filepopup.GetPath(val) == currentPath {
				idx = i
				// break out of the loop once we find the node
				break
			}
		}
	}

	// Old logic for finding index was buggy and could lead to an index of -1. Added more robustness.

	if idx == -1 {
		// If the previously selected node is gone (e.g. deleted),
		// attempt to select the previous index, but ensure it's within bounds.
		if m.CurrentIndex > 0 && m.CurrentIndex-1 < len(m.VisibleNodes) {
			idx = m.CurrentIndex - 1
		} else {
			// Default to the first item if the old index is invalid.
			idx = 0
		}
	}

	m.CurrentIndex = idx
	m.CurrentNode = m.VisibleNodes[idx]
}
