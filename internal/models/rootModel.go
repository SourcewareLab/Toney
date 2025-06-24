package models

import (
	"fmt"

	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/messages"
	filepopup "github.com/SourcewareLab/Toney/internal/models/filePopup"
	homemodel "github.com/SourcewareLab/Toney/internal/models/homeModel"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	Width         int
	Height        int
	Page          enums.Page
	Home          *homemodel.HomeModel
	ShowPopup     bool
	FilePopupType enums.PopupType
	FilePopup     *filepopup.FilePopup
	isLoading     bool
}

func NewRoot() *RootModel {
	return &RootModel{
		Page:      enums.Home,
		Home:      nil,
		ShowPopup: false,
		isLoading: true,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.ShowLoader:
		m.isLoading = true
		return m, nil
	case messages.HideLoader:
		m.isLoading = false
		return m, nil
	case messages.ShowPopupMessage:
		m.FilePopup = filepopup.NewPopup(msg.Type, msg.Curr)
		m.ShowPopup = true
	case messages.HidePopupMessage:
		m.ShowPopup = false
	case messages.RefreshFileExplorerMsg:
		m.Home.FileExplorer.Update(msg)
		return m, nil
	case messages.EditorClose:
		if msg.Err != nil {
			fmt.Println(msg.Err.Error())
		}
		m.Home.FileExplorer.Update(msg)
		m.Home.Viewer.Update(msg)

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		m.Home = homemodel.NewHome(msg.Width, msg.Height)

		m.isLoading = false

		return m, nil
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

	if m.isLoading {
		return lipgloss.NewStyle().Render("Loading...")
	}

	if m.ShowPopup && m.FilePopup != nil {
		return lipgloss.Place(m.Width+2, m.Height+2, lipgloss.Center, lipgloss.Center, m.FilePopup.View())
	}

	return lipgloss.NewStyle().Background(lipgloss.Color("#1e1e2e")).Render(bg.View())
}
