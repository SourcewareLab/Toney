package models

import (
	"fmt"

	"toney/internal/enums"
	"toney/internal/messages"
	filepopup "toney/internal/models/filePopup"
	homemodels "toney/internal/models/homeModels"

	"github.com/charmbracelet/bubbles/spinner"
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
	isLoading     bool
	Loader        spinner.Model
	isSized       bool
}

type (
	RefreshFileExplorerMsg struct{}
	NvimDoneMsg            struct {
		Err string
	}
)

func NewRoot() *RootModel {
	loader := spinner.Model{}
	loader.Spinner = spinner.Dot
	return &RootModel{
		Page:      enums.Home,
		Home:      nil,
		ShowPopup: false,
		isLoading: true,
		Loader:    loader,
		isSized:   false,
	}
}

func (m RootModel) Init() tea.Cmd {
	return m.Loader.Tick
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.RendererCreated:
		m.Home.Viewer.Update(msg)
		m.isSized = true
		return m, nil
	case spinner.TickMsg:
		ld, cmd := m.Loader.Update(msg)
		m.Loader = ld
		return m, cmd
	case messages.ShowLoader:
		m.isLoading = true
		return m, nil
	case messages.HideLoader:
		m.isLoading = false
		return m, nil
	case filepopup.ShowPopupMessage:
		m.FilePopup = filepopup.NewPopup(msg.Type, msg.Curr)
		m.ShowPopup = true
	case filepopup.HidePopupMessage:
		fmt.Println("Hiding Popup")
		m.ShowPopup = false
	case filepopup.RefreshFileExplorerMsg:
		m.Home.FileExplorer.Update(msg)
		return m, nil
	case homemodels.EditorClose:
		if msg.Err != nil {
			fmt.Println(msg.Err.Error())
		}

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		if m.isSized {
			return m, nil
		}

		m.Width = msg.Width
		m.Height = msg.Height

		m.Home = NewHome(msg.Width, msg.Height)

		return m, homemodels.InitRenderer(msg.Width)
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

	if !m.isSized {
		content := fmt.Sprintf("%s  %s\n", m.Loader.View(), "Loading...")
		return lipgloss.NewStyle().Render(content)
	}

	if m.ShowPopup && m.FilePopup != nil {
		return lipgloss.Place(m.Width+2, m.Height+2, lipgloss.Center, lipgloss.Center, m.FilePopup.View())
	}

	return bg.View()
}
