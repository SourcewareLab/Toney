package keymap

import "github.com/charmbracelet/bubbles/key"

type ViewerKeyMap struct {
	ScrollUp   key.Binding
	ScrollDown key.Binding
}

func NewViewerKeyMap() ViewerKeyMap {
	return ViewerKeyMap{
		ScrollUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "scroll down"),
		),
	}
}
