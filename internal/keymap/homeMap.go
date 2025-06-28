package keymap

import "github.com/charmbracelet/bubbles/key"

type HomeKeyMap struct {
	FocusViewer   key.Binding
	FocusExplorer key.Binding
}

func NewHomeKeyMap() HomeKeyMap {
	return HomeKeyMap{
		FocusExplorer: key.NewBinding(
			key.WithKeys("F"),
			key.WithHelp("F", "file explorer"),
		),
		FocusViewer: key.NewBinding(
			key.WithKeys("V"),
			key.WithHelp("V", "viewer"),
		),
	}
}
