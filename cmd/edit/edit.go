package edit

import (
	"os"
	"os/exec"

	"toney/internal/utils"

	"github.com/spf13/cobra"
)

var EditCommand = &cobra.Command{
	Use:   "edit",
	Short: "interact with your toney notes using an editor",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		dirPath, err := utils.CheckNoteDir()
		if err != nil {
			return err
		}

		editorCommand := exec.Command("nvim", dirPath)
		editorCommand.Stdout = os.Stdout
		editorCommand.Stderr = os.Stderr
		editorCommand.Stdin = os.Stdin

		if err := editorCommand.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
}
