package keymap

import (
	"reflect"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/charmbracelet/bubbles/key"
)

type DailyTaskMap struct {
	CreateUnique    key.Binding
	CreateRecurring key.Binding
	EditTask        key.Binding
	ChangeStatus    key.Binding
	DeleteTask      key.Binding
	BackToMenu      key.Binding
	TabRight        key.Binding
	TabLeft         key.Binding
}

func NewDailyTaskMap() DailyTaskMap {
	cfg := config.AppConfig.Keybinds.Daily
	return DailyTaskMap{
		CreateUnique: key.NewBinding(
			key.WithKeys(cfg.Create),
			key.WithHelp(cfg.Create, "create"),
		),
		CreateRecurring: key.NewBinding(
			key.WithKeys(cfg.CreateRecurring),
			key.WithHelp(cfg.CreateRecurring, "new recurring"),
		),
		EditTask: key.NewBinding(
			key.WithKeys(cfg.Edit),
			key.WithHelp(cfg.Edit, "edit"),
		),
		ChangeStatus: key.NewBinding(
			key.WithKeys(cfg.StatusChange),
			key.WithHelp(cfg.StatusChange, "status change"),
		),
		DeleteTask: key.NewBinding(
			key.WithKeys(cfg.Delete),
			key.WithHelp(cfg.Delete, "delete"),
		),
		BackToMenu: key.NewBinding(
			key.WithKeys(cfg.BackToMenu),
			key.WithHelp(cfg.BackToMenu, "return to menu"),
		),
		TabRight: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("right", "right tab"),
		),
		TabLeft: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("left", "left tab"),
		),
	}
}

func (m DailyTaskMap) Bindings() []key.Binding {
	var bindings []key.Binding

	v := reflect.ValueOf(m)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if binding, ok := field.Interface().(key.Binding); ok {
			bindings = append(bindings, binding)
		}
	}

	return bindings
}
