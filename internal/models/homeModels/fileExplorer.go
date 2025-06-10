package homemodels

import (
	"fmt"
	"os"

	"toney/internal/enums"
	filetree "toney/internal/fileTree"
	filepopup "toney/internal/models/filePopup"
	"toney/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileExplorer struct {
	path         string
	IsFocused    bool
	Width        int
	Height       int
	Root         *filetree.Node
	CurrentNode  *filetree.Node
	CurrentIndex int
	VisibleNodes []*filetree.Node
}

func NewFileExplorer() *FileExplorer {
	root, err := filetree.CreateTree()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return &FileExplorer{
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
			return m, nil
		case "up":
			if m.CurrentIndex <= 0 {
				return m, nil
			}

			m.CurrentIndex -= 1
			m.CurrentNode = m.VisibleNodes[m.CurrentIndex]
			return m, nil
		case "enter":
			if m.CurrentNode.IsDirectory {
				m.CurrentNode.IsExpanded = !m.CurrentNode.IsExpanded
				m.VisibleNodes = filetree.FlattenVisibleTree(m.Root)
				return m, nil
			}

		case "c":
			return m, func() tea.Msg {
				return filepopup.PopupMessage{
					Type: enums.FileCreate,
					Show: true,
				}
			}

		}
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width

	}

	return m, nil
}

func (m FileExplorer) View() string {
	style := styles.BorderStyle
	style = style.Align(lipgloss.Left, lipgloss.Top)

	if m.IsFocused {
		style = style.BorderForeground(lipgloss.Color("#bb9af7"))
	}

	w := (m.Width / 4) - 1
	h := m.Height - 2

	s := filetree.BuildNodeTree(m.Root, "", len(m.Root.Children) == 0, m.CurrentNode)

	return style.Width(w).Height(h).MarginTop(1).Render(s)
}
