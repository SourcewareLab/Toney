package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var AppConfig Config
var NoteRootPath string

func SetConfig() {
	home, _ := os.UserHomeDir()

	viper.AddConfigPath(filepath.Join(home, ".config", "toney"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, using defaults")
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}

	if AppConfig.KeyBinding == "" {
		AppConfig.KeyBinding = "normal"
	}

	NoteRootPath = filepath.Join(home, AppConfig.NotePath)
}
