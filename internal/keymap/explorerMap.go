package keymap

import "github.com/charmbracelet/bubbles/key"

type ExplorerKeyMap struct {
	CreateFile  key.Binding
	MoveFile    key.Binding
	RenameFile  key.Binding
	DeleteFile  key.Binding
	OpenForEdit key.Binding
	Up          key.Binding
	Down        key.Binding
}

func NewExplorerKeyMap() ExplorerKeyMap {
	return ExplorerKeyMap{
		CreateFile: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create"),
		),
		MoveFile: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "move"),
		),
		RenameFile: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename"),
		),
		DeleteFile: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		OpenForEdit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "edit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
	}
}
