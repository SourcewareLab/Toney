package models

import (
	"toney/internal/enums"

	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	Width  int
	Height int
	Page   enums.Page
	Home   HomeModel
}

func NewRoot() *RootModel {
	return &RootModel{
		Page: enums.Home,
		Home: NewHome(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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

	_, cmd = m.Home.Update(msg)

	return m, cmd
}

func (m RootModel) View() string {
	return m.Home.View()
}
