package main

import (
	"fmt"

	"github.com/SourcewareLab/Toney/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(models.NewRoot(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Alas, error")
		fmt.Println(err.Error())
		return
	}
}
