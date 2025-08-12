package utils

import (
	"github.com/SourcewareLab/Toney/internal/messages"
	tea "github.com/charmbracelet/bubbletea"
)

func ReturnError(locn, title string, err error) tea.Cmd {
	return func() tea.Msg {
		return messages.ErrorMsg{
			Locn:  locn,
			Title: title,
			Msg:   err.Error(),
		}
	}
}
