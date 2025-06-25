package config

type StylesConfig struct {
	FocusedBorder string `mapstructure:"focused_border"`
	NormalBorder  string `mapstructure:"normal_border"`
	Text          string `mapstructure:"text"`
}

type Config struct {
	NotePath   string       `mapstructure:"note_path"`
	Editor     string       `mapstructure:"editor"`
	KeyBinding string       `mapstructure:"keybinding"`
	Styles     StylesConfig `mapstructure:"styles"`
}
