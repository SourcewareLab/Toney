package keymap

import "github.com/charmbracelet/bubbles/key"

type DynamicMap struct {
	keys []key.Binding
}

func NewDynamic(keys []key.Binding) *DynamicMap {
	return &DynamicMap{
		keys: keys,
	}
}

func (d DynamicMap) ShortHelp() []key.Binding {
	return d.keys
}

func (d DynamicMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{d.keys}
}
