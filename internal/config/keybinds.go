package config

type KeybindsConfig struct {
	Global GlobalKeybinds `mapstructure:"global"`
	Home   HomeKeybinds   `mapstructure:"home"`
	Daily  DailyKeybinds  `mapstructure:"daily"`
	Fuzz   FuzzyKeybinds  `mapstructure:"fuzzy"`
	Diary  DiaryKeybinds  `mapstructure:"diary"`
}

type FuzzyKeybinds struct {
	StartWriting string `mapstructure:"start_writing"`
	Up           string `mapstructure:"up"`
	Down         string `mapstructure:"down"`
	Enter        string `mapstructure:"enter"`
	Exit         string `mapstructure:"exit"`
}

type DiaryKeybinds struct {
	Edit       string `mapstructure:"edit"`
	ScrollUp   string `mapstructure:"scroll_up"`
	ScrollDown string `mapstructure:"scroll_down"`
	Finder     string `mapstructure:"finder"`
	BackToMenu string `mapstructure:"return_to_menu"`
}

type GlobalKeybinds struct {
	Up     string `mapstructure:"up"`
	Down   string `mapstructure:"down"`
	Script string `mapstructure:"script"`
}

type DailyKeybinds struct {
	Create       string `mapstructure:"create"`
	Delete       string `mapstructure:"delete"`
	StatusChange string `mapstructure:"status_change"`
	Edit         string `mapstructure:"edit"`
	BackToMenu   string `mapstructure:"return_to_menu"`
	FormUp       string `mapstructure:"form_up"`
	FormDown     string `mapstructure:"form_down"`
	ExitPopup    string `mapstructure:"exit_popup"`
	Enter        string `mapstructure:"enter"`
}

type HomeKeybinds struct {
	BackToMenu    string `mapstructure:"return_to_menu"`
	FocusViewer   string `mapstructure:"focus_viewer"`
	FocusExplorer string `mapstructure:"focus_explorer"`
	Create        string `mapstructure:"create"`
	Rename        string `mapstructure:"rename"`
	Move          string `mapstructure:"move"`
	Delete        string `mapstructure:"delete"`
	Edit          string `mapstructure:"edit"`
	ScrollUp      string `mapstructure:"scroll_up"`
	ScrollDown    string `mapstructure:"scroll_down"`
	Finder        string `mapstructure:"finder"`
}
