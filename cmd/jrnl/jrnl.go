package jrnl

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var JrnlCommand = &cobra.Command{
	Use:   "jrnl",
	Short: "interact with your daily journal",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := GetTodayJrnl()
		if err != nil {
			return err
		}

		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		fmt.Println(string(content))

		return nil
	},
}

func init() {
	JrnlCommand.AddCommand(addCommand)
}
