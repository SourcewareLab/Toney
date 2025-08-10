package config

type Config struct {
	General  GeneralConfig  `mapstructure:"general"`
	Styles   StylesConfig   `mapstructure:"styles"`
	Keybinds KeybindsConfig `mapstructure:"keybinds"`
	GitHub   GitHubConfig   `mapstructure:"github"`
}

type GeneralConfig struct {
	Editor      []string `mapstructure:"editor"`
	NotesDir    string   `mapstructure:"notes_dir"`
	StartScript []string `mapstructure:"start_script"`
	StopScript  []string `mapstructure:"stop_script"`
	Script      []string `mapstructure:"script"`
}

type GitHubConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Token   string `mapstructure:"token"`
	Owner   string `mapstructure:"owner"`
	Repo    string `mapstructure:"repo"`
}
