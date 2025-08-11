package menu

import (
	"strings"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuList struct {
	Width      int
	Height     int
	Options    map[enums.Page]string
	Selections []enums.Page
	Selected   int
}

func NewMenuList(w int, h int, opts map[enums.Page]string) *MenuList {
	selections := []enums.Page{enums.HomePage, enums.DailyPage, enums.DiaryPage}

	// Add GitHub option if it exists in the options map
	if _, hasGitHub := opts[enums.GitHubPage]; hasGitHub {
		selections = append(selections, enums.GitHubPage)
	}

	// Always add Quit at the end
	selections = append(selections, enums.Quit)

	return &MenuList{
		Width:      w,
		Height:     h,
		Options:    opts,
		Selections: selections,
		Selected:   0,
	}
}

func (m *MenuList) Init() tea.Cmd {
	return nil
}

func (m *MenuList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case config.AppConfig.Keybinds.Global.Down:
			if m.Selected < len(m.Selections)-1 {
				m.Selected += 1
			}
			return m, nil
		case config.AppConfig.Keybinds.Global.Up:
			if m.Selected > 0 {
				m.Selected -= 1
			}
			return m, nil
		case "enter":
			return m, func() tea.Msg {
				return messages.ChangePage{
					Page: m.Selections[m.Selected],
				}
			}
		}
	}

	return m, nil
}

func (m *MenuList) View() string {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(colors.ColorPalette().Border).Render(m.GetText())
}

func (m *MenuList) GetText() string {
	text := ""
	style := lipgloss.NewStyle().Width(m.Width).Padding(0, 2).Foreground(colors.ColorPalette().Text)

	for idx, val := range m.Selections {
		line := m.Options[val]
		if m.Selected == idx {
			text += style.Background(colors.ColorPalette().MenuSelectedBg).
				Foreground(colors.ColorPalette().MenuSelectedText).
				Render(line) + "\n"

			continue
		}

		text += style.Render(line) + "\n"
	}

	return strings.TrimSuffix(text, "\n")
}
