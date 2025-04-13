package jrnl

import (
	"toney/internal/utils"

	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "add a task to your journal",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := utils.CheckNoteDir()
		if err != nil {
			return
		}
	},
}
