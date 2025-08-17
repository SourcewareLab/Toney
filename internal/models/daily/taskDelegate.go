package daily

import (
	"fmt"
	"io"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskDelegate struct{}

func (d TaskDelegate) Height() int                               { return 2 }
func (d TaskDelegate) Spacing() int                              { return 1 }
func (d TaskDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TaskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if len(m.Items()) == 0 {
		return
	}

	t, ok := item.(Task)
	if !ok {
		return
	}

	text := ""
	cfg := config.AppConfig.Styles.Icons.TaskIcons

	switch t.Status {
	case enums.Complete:
		text += fmt.Sprintf("%s\n%s",
			styles.CompletedStyle().Title.Render(fmt.Sprintf("%s | %s", cfg.CompletedIcon, Shorten(t.Title(), 35))),
			styles.CompletedStyle().Desc.Render(t.Description()),
		)
	case enums.Pending:
		text += fmt.Sprintf("%s\n%s",
			styles.PendingStyle().Title.Render(fmt.Sprintf("%s | %s", cfg.PendingIcon, Shorten(t.Title(), 35))),
			styles.PendingStyle().Desc.Render(t.Description()),
		)
	case enums.Started:
		text += fmt.Sprintf("%s\n%s",
			styles.StartedStyle().Title.Render(fmt.Sprintf("%s | %s", cfg.StartedIcon, Shorten(t.Title(), 35))),
			styles.StartedStyle().Desc.Render(t.Description()),
		)
	case enums.Abandoned:
		text += fmt.Sprintf("%s\n%s",
			styles.AbandonedStyle().Title.Render(fmt.Sprintf("%s | %s", cfg.AbandonedIcon, Shorten(t.Title(), 35))),
			styles.AbandonedStyle().Desc.Render(t.Description()),
		)
	}

	border := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).
		PaddingLeft(2).
		BorderLeft(true).BorderBottom(false).BorderRight(false).BorderTop(false)

	if index == m.Index() {
		text = border.BorderForeground(colors.ColorPalette().TaskFocusedBar).Render(text)
	} else {
		text = border.BorderForeground(colors.ColorPalette().TaskUnfocusedBar).Render(text)
	}

	io.WriteString(w, text)
}

func Shorten(s string, maxLen int) string {
	if len(s) <= maxLen {
		return lipgloss.NewStyle().Width(maxLen).Render(s)
	}
	if maxLen <= 1 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "…"
}
