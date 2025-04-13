package jrnl

import (
	"os"
	"time"

	"toney/internal/utils"

	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "add a task to your journal",
	Run: func(cmd *cobra.Command, args []string) {
		dirPath, err := utils.CheckNoteDir()
		if err != nil {
			return
		}

		fileName := dirPath + time.Now().UTC().Format("2006_01_02") + ".toney"

		_, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0o644)
	},
}
