package models

import (
	homemodels "toney/internal/models/homeModels"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	Width        int
	Height       int
	FocusOn      homemodels.Splits
	FileExplorer homemodels.FileExplorer
	Vault        homemodels.Vault
	Task         homemodels.Task
}

func NewHome() HomeModel {
	return HomeModel{
		FocusOn:      homemodels.Tasks,
		FileExplorer: homemodels.FileExplorer{},
		Vault:        homemodels.Vault{},
		Task:         homemodels.Task{},
	}
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.Height = msg.Height - 2
		m.Width = msg.Width - 2

		m.Vault.Update(msg)
		m.Task.Update(msg)
		m.FileExplorer.Update(msg)
	}

	var cmd tea.Cmd

	if m.FileExplorer.IsFocused {
		_, cmd = m.FileExplorer.Update(msg)
	} else if m.Vault.IsFocused {
		_, cmd = m.Vault.Update(msg)
	} else {
		_, cmd = m.Task.Update(msg)
	}

	return m, cmd
}

func (m HomeModel) View() string {
	m.FileExplorer.IsFocused = false
	m.Task.IsFocused = false
	m.Vault.IsFocused = false

	if m.FocusOn == homemodels.File {
		m.FileExplorer.IsFocused = true
	} else if m.FocusOn == homemodels.Vaults {
		m.Vault.IsFocused = true
	} else {
		m.Task.IsFocused = true
	}

	leftCol := lipgloss.JoinVertical(lipgloss.Center, m.FileExplorer.View(), m.Vault.View())
	rightCol := m.Task.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)
}
