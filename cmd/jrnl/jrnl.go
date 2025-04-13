package jrnl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var JrnlCommand = &cobra.Command{
	Use:   "jrnl",
	Short: "interact with your daily journal",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("jrnl subcommand")
		return nil
	},
}

func init() {
	JrnlCommand.AddCommand(addCommand)
}
