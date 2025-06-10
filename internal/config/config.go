package config

import "github.com/spf13/viper"

func SetConfig() {
	viper.SetConfigFile("config")
	viper.SetConfigFile("yaml")
	viper.AddConfigPath("$HOME/.config/toney")
	viper.AddConfigPath("$HOME/.toney")
}
