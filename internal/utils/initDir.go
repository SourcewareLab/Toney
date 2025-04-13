package utils

import (
	"fmt"
	"os"
)

func CheckNoteDir() (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Home Directory not found.")
	}

	dirPath := userDir + "/.toney"

	if _, err := os.Stat(dirPath); err != nil {
		fmt.Println("No previous note directory found.")
		fmt.Printf("Creating directory .toney in %s.\n", userDir)

		err := os.Mkdir(dirPath, 0o750)
		if err != nil {
			fmt.Println("Error creating directory .toney, exiting.")
			fmt.Printf("Error -> %s", err.Error())
			return "", err
		}

		fmt.Println("Successfully created directory .toney.")
	}
	return dirPath + "/", nil
}
