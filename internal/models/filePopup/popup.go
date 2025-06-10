package filepopup

import (
	"toney/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var popupStyle = styles.BorderStyle.Align(lipgloss.Left, lipgloss.Top).BorderForeground(lipgloss.Color("#bb9af7"))

type FilePopup struct {
	Height int
	Width  int
}

func NewPopup() *FilePopup {
	return &FilePopup{}
}

func (m FilePopup) Init() tea.Cmd {
	return nil
}

func (m *FilePopup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	}

	return m, nil
}

func (m FilePopup) View() string {
	s := "This is a popup"

	return popupStyle.Render(s)
}
