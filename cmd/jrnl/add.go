package jrnl

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"toney/internal/models"

	"github.com/spf13/cobra"

	"github.com/manifoldco/promptui"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "add a task to your journal",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := GetTodayJrnl()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		task := taskPrompt()

		fmt.Println(task)

		err = appendToFile(task, f)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func appendToFile(task models.Task, f *os.File) error {
	n, err := f.WriteString(fmt.Sprintf("[%s]{%s}%s\n", task.Status, task.Title, task.Descr))
	if err != nil {
		return err
	}

	fmt.Println(n)

	return nil
}

func taskPrompt() models.Task {
	statusValidate := func(input string) error {
		input = strings.ToLower(input)
		if input != "x" && input != "." && input != "-" {
			return errors.New("Invalid task status. The supported types are :-\n\t1. 'x' : Completed\n\t2. '.' : Pending\n\t3. '-' : Suspended")
		}

		return nil
	}

	taskStatusPrompt := promptui.Prompt{
		Label:    "Enter current status of task",
		Validate: statusValidate,
	}

	taskTitlePrompt := promptui.Prompt{
		Label: "Enter title of task",
	}

	taskDescrPrompt := promptui.Prompt{
		Label: "Enter task description",
	}

	title, _ := taskTitlePrompt.Run()
	stat, _ := taskStatusPrompt.Run()
	desc, _ := taskDescrPrompt.Run()

	return models.NewTask(title, desc, stat)
}
