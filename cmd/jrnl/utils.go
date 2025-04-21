package jrnl

import (
	"fmt"
	"os"
	"time"

	"toney/internal/utils"
)

func GetTodayJrnl() (*os.File, error) {
	dirPath, err := utils.CheckNoteDir()
	if err != nil {
		return &os.File{}, err
	}

	fileName := dirPath + time.Now().UTC().Format("2006_01_02") + ".toney"

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		fmt.Println("ERROR: error opening journal file.")
		fmt.Printf("err -> %s\n", err.Error())
		return f, err
	}

	return f, nil
}
