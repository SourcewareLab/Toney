package menu

import (
	"strings"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/messages"
	"github.com/SourcewareLab/Toney/internal/ui/theme"
	"github.com/SourcewareLab/Toney/internal/ui/widgets"
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
	pal := theme.DefaultDarkPalette()

	// Header
	header := widgets.Header{
		Palette: pal,
		Width:   m.Width,
		Title:   "Toney · Menu",
		Right:   "",
	}.View()

	// Body: bordered container for list
	container := lipgloss.NewStyle().
		Border(theme.Borders()).
		BorderForeground(pal.Border).
		Foreground(pal.Fg).
		Width(m.Width).
		Padding(0, 1)

	// Estimate header/footer heights (~3 lines each including borders)
	footer := widgets.Footer{
		Palette: pal,
		Width:   m.Width,
		Hints: []widgets.KeyHint{
			{Key: "↑↓", Desc: "navigate"},
			{Key: "enter", Desc: "select"},
			{Key: "esc", Desc: "back"},
		},
	}.View()

	// Compute body height to avoid clipping
	headerH := lipgloss.Height(header)
	footerH := lipgloss.Height(footer)
	bodyH := m.Height - headerH - footerH
	if bodyH < 3 {
		bodyH = 3
	}

	body := container.Height(bodyH).Render(m.GetTextStyled(pal))

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		footer,
	)
}

func (m *MenuList) GetTextStyled(pal theme.Palette) string {
	text := ""
	base := lipgloss.NewStyle().Width(m.Width-4).Padding(0, 1).Foreground(pal.Fg)
	sel := base.Foreground(pal.SelectedFg).Background(pal.SelectedBg).Bold(true)

	for idx, val := range m.Selections {
		line := m.Options[val]
		if m.Selected == idx {
			text += sel.Render(line) + "\n"
			continue
		}
		text += base.Render(line) + "\n"
	}

	return strings.TrimSuffix(text, "\n")
}
