package cmd

import (
	"fmt"

	"toney/cmd/jrnl"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(jrnl.JrnlCommand)

	// Persistent Flag
}

var rootCmd = &cobra.Command{
	Use:   "toney",
	Short: "a note-taking terminal application.",
	Long:  "lorem ipsum",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Toney is a note-taking terminal app.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
