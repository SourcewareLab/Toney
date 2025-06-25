package fileexplorer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SourcewareLab/Toney/internal/config"
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
	KeyMap        KeyMap
}

type KeyMap struct {
	Up     string
	Down   string
	Select string
	Quit   string
	Create string
	Delete string
	Rename string
	Move   string
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
		KeyMap:       LoadKeyMap(config.AppConfig.KeyBinding),
	}
}

func (m FileExplorer) Init() tea.Cmd {
	return nil
}

func LoadKeyMap(mode string) KeyMap {
	var km KeyMap
	switch mode {
	case "vim":
		km = KeyMap{
			Up:     "k",
			Down:   "j",
			Select: "enter",
			Quit:   "q",
			Create: "a",
			Delete: "d",
			Rename: "r",
			Move:   "m",
		}
	default:
		km = KeyMap{
			Up:     "up",
			Down:   "down",
			Select: "enter",
			Quit:   "q",
			Create: "c",
			Delete: "d",
			Rename: "r",
			Move:   "m",
		}
	}

	// Apply overrides
	// if v, ok := overrides["create"]; ok {
	// 	km.Create = v
	// }
	// if v, ok := overrides["delete"]; ok {
	// 	km.Delete = v
	// }
	// if v, ok := overrides["rename"]; ok {
	// 	km.Rename = v
	// }
	// if v, ok := overrides["move"]; ok {
	// 	km.Move = v
	// }
	// if v, ok := overrides["up"]; ok {
	// 	km.Up = v
	// }
	// if v, ok := overrides["down"]; ok {
	// 	km.Down = v
	// }
	// if v, ok := overrides["select"]; ok {
	// 	km.Select = v
	// }
	// if v, ok := overrides["quit"]; ok {
	// 	km.Quit = v
	// }

	return km
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
		case m.KeyMap.Down:
			if m.CurrentIndex >= len(m.VisibleNodes)-1 {
				return m, nil
			}
			m.CurrentIndex += 1
			m.CurrentNode = m.VisibleNodes[m.CurrentIndex]
			return m, m.SelectionChanged(m.CurrentNode)
		case m.KeyMap.Up:
			if m.CurrentIndex <= 0 {
				return m, nil
			}
			m.CurrentIndex -= 1
			m.CurrentNode = m.VisibleNodes[m.CurrentIndex]
			return m, m.SelectionChanged(m.CurrentNode)
		case m.KeyMap.Select:
			if m.CurrentNode.IsDirectory {
				m.CurrentNode.IsExpanded = !m.CurrentNode.IsExpanded
				m.VisibleNodes = filetree.FlattenVisibleTree(m.Root)
				return m, nil
			}

			relPath := strings.TrimSuffix(filepopup.GetPath(m.CurrentNode), "/")
			homeDir, _ := os.UserHomeDir()
			var fullPath string
			if strings.HasPrefix(relPath, "~") {
				fullPath = filepath.Join(homeDir, relPath[1:])
			} else if filepath.IsAbs(relPath) {
				fullPath = relPath
			} else {
				fullPath = filepath.Join(homeDir, relPath)
			}

			c := exec.Command(config.AppConfig.Editor, fullPath)
			cmd := tea.ExecProcess(c, func(err error) tea.Msg {
				return messages.EditorClose{
					Err: err,
				}
			})
			return m, cmd

		case m.KeyMap.Create:
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileCreate,
					Curr: m.CurrentNode,
				}
			}
		case m.KeyMap.Delete:
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileDelete,
					Curr: m.CurrentNode,
				}
			}

		case m.KeyMap.Move:
			return m, func() tea.Msg {
				return messages.ShowPopupMessage{
					Type: enums.FileMove,
					Curr: m.CurrentNode,
				}
			}
		case m.KeyMap.Rename:
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
