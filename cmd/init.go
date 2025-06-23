package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Toney configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing Toney...")
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Could not find $HOME directory", err.Error())
			return
		}

		err = os.MkdirAll(home+"/.toney", 0o755)
		if err != nil {
			fmt.Println("Could not create .toney directory", err.Error())
			return
		}

		fmt.Println("Toney setup was succesfull!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
