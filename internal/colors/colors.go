package colors

import (
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/charmbracelet/lipgloss"
)

type Colors struct {
	Text             lipgloss.Color
	Background       lipgloss.Color
	Border           lipgloss.Color
	FocusedBorder    lipgloss.Color
	MenuSelectedBg   lipgloss.Color
	MenuSelectedText lipgloss.Color
	TaskFocusedBar   lipgloss.Color
	TaskUnfocusedBar lipgloss.Color
	ErrorBg          lipgloss.Color
	ErrorText        lipgloss.Color
	CompletedTask    TaskColors
	AbandonedTask    TaskColors
	PendingTask      TaskColors
	StartedTask      TaskColors
}

type TaskColors struct {
	TaskTitle lipgloss.Color
	TaskDesc  lipgloss.Color
}

func ColorPalette() Colors {
	cfg := config.AppConfig.Styles
	return Colors{
		Text:             lipgloss.Color(cfg.Text),
		Background:       lipgloss.Color(cfg.Background),
		Border:           lipgloss.Color(cfg.Border),
		FocusedBorder:    lipgloss.Color(cfg.FocusedBorder),
		TaskFocusedBar:   lipgloss.Color(cfg.TaskStyles.FocusedBar),
		TaskUnfocusedBar: lipgloss.Color(cfg.TaskStyles.UnfocusedBar),
		MenuSelectedBg:   lipgloss.Color(cfg.MenuSelectedBg),
		MenuSelectedText: lipgloss.Color(cfg.MenuSelectedText),
		ErrorBg:          lipgloss.Color(cfg.ErrorBg),
		ErrorText:        lipgloss.Color(cfg.ErrorText),
		CompletedTask: TaskColors{
			TaskTitle: lipgloss.Color(cfg.TaskStyles.CompletedStyle.Title),
			TaskDesc:  lipgloss.Color(cfg.TaskStyles.CompletedStyle.Desc),
		},
		AbandonedTask: TaskColors{
			TaskTitle: lipgloss.Color(cfg.TaskStyles.AbandonedStyle.Title),
			TaskDesc:  lipgloss.Color(cfg.TaskStyles.AbandonedStyle.Desc),
		},
		PendingTask: TaskColors{
			TaskTitle: lipgloss.Color(cfg.TaskStyles.PendingStyle.Title),
			TaskDesc:  lipgloss.Color(cfg.TaskStyles.PendingStyle.Desc),
		},
		StartedTask: TaskColors{
			TaskTitle: lipgloss.Color(cfg.TaskStyles.StartedStyle.Title),
			TaskDesc:  lipgloss.Color(cfg.TaskStyles.StartedStyle.Desc),
		},
	}
}
