package main

import (
	"github.com/SourcewareLab/Toney/cmd"
	"github.com/SourcewareLab/Toney/internal/config"
)

func main() {
	config.SetConfig()
	cmd.Execute()
}
