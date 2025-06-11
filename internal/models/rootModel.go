package models

import (
	"fmt"

	"toney/internal/enums"
	filepopup "toney/internal/models/filePopup"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	Width         int
	Height        int
	Page          enums.Page
	Home          *HomeModel
	ShowPopup     bool
	FilePopupType enums.PopupType
	FilePopup     *filepopup.FilePopup
}

type RefreshFileExplorerMsg struct{}

func NewRoot() *RootModel {
	return &RootModel{
		Page:      enums.Home,
		Home:      NewHome(),
		ShowPopup: false,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case filepopup.ShowPopupMessage:
		fmt.Println("Received PopupMessage", msg.Type)
		m.FilePopup = filepopup.NewPopup(msg.Type, msg.Curr)
		m.ShowPopup = true
	case filepopup.HidePopupMessage:
		fmt.Println("Hiding Popup")
		m.ShowPopup = false
	case filepopup.RefreshFileExplorerMsg:
		m.Home.FileExplorer.Update(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width - 2
		m.Height = msg.Height - 2

		m.Home.Update(msg)
	}

	var cmd tea.Cmd

	if m.ShowPopup {
		_, cmd = m.FilePopup.Update(msg)
	} else {
		_, cmd = m.Home.Update(msg)
	}

	return m, cmd
}

func (m *RootModel) View() string {
	bg := m.Home

	if m.ShowPopup && m.FilePopup != nil {
		return lipgloss.Place(m.Width+2, m.Height+2, lipgloss.Center, lipgloss.Center, m.FilePopup.View())
	}

	return bg.View()
}
