package menu

import (
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	Width   int
	Height  int
	Options map[enums.Page]string
	List    *MenuList
}

func NewMenu(w int, h int) *Menu {
	opts := map[enums.Page]string{
		enums.HomePage:  "Home",
		enums.DailyPage: "Daily Tasks",
		enums.DiaryPage: "Diary",
		enums.Quit:      "Quit",
	}

	list := NewMenuList(w/3, h/2-1, opts)
	return &Menu{
		Width:   w,
		Height:  h,
		Options: opts,
		List:    list,
	}
}

func (m *Menu) Init() tea.Cmd {
	return nil
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	_, cmd = m.List.Update(msg)

	return m, cmd
}

func (m *Menu) View() string {
	logo := styles.GetLogo(m.Width, m.Height/2)
	list := lipgloss.PlaceHorizontal(m.Width/3, lipgloss.Center, m.List.View())

	return lipgloss.JoinVertical(lipgloss.Center, logo, list)
}
