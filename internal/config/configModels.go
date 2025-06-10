package config

type StylesConfig struct {
	FocusedBorder string `mapstructure:"focused_border"`
	NormalBorder  string `mapstructure:"normal_border"`
	Text          string `mapstructure:"text"`
}

type Config struct {
	Styles StylesConfig `mapstructure:"styles"`
}
