package filepopup

import (
	"fmt"

	"toney/internal/enums"
	filetree "toney/internal/fileTree"
	"toney/internal/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	popupStyle  = styles.BorderStyle.Align(lipgloss.Left, lipgloss.Top).BorderForeground(lipgloss.Color("#bb9af7"))
	headerStyle = lipgloss.NewStyle().Background(lipgloss.Color("#bb9af7")).Foreground(lipgloss.Color("#1e1e2e"))
)

type FilePopup struct {
	Height    int
	Width     int
	Type      enums.PopupType
	TextInput textinput.Model
	Node      *filetree.Node
}

func NewPopup(typ enums.PopupType, node *filetree.Node) *FilePopup {
	ti := textinput.New()
	ti.Focus()

	return &FilePopup{
		Type:      typ,
		TextInput: ti,
		Node:      node,
	}
}

type ShowPopupMessage struct {
	Type enums.PopupType
	Curr *filetree.Node
}

type HidePopupMessage struct{}

type RefreshFileExplorerMsg struct{}

func (m FilePopup) Init() tea.Cmd {
	return nil
}

func (m *FilePopup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// case PopupMessage:
	// 	return m, func() tea.Msg {
	// 		return msg
	// 	}

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return HidePopupMessage{}
			}
		case "enter":
			return HandleEnter(m)
		}
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	}

	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m *FilePopup) View() string {
	w := m.Width / 3
	h := m.Height / 3

	return popupStyle.Width(w).Height(h).Render(GetText(m.Type, m.TextInput))
}

func GetText(typ enums.PopupType, ti textinput.Model) string {
	header := ""
	switch typ {
	case enums.FileCreate:
		header = "Create a File (names ending with '/' will create directory):- "
	case enums.FileDelete:
		header = "Delete file (y/n)?"
	case enums.FileRename:
		header = "Enter new name for file:- "
	case enums.FileMove:
		header = "Enter new location (relative) for file:- "
	default:
		header = "TBD"
	}

	return fmt.Sprintf("%s\n\n%s", headerStyle.Render(header), ti.View())
}
