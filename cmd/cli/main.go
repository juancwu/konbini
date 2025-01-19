package main

import (
	"log"
	"os"

	command "konbini/cli/commands"
	"konbini/cli/models"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) == 1 {
		m := models.NewModel()
		p := tea.NewProgram(m, tea.WithAltScreen())
		_, err := p.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := command.Execute()
		if err != nil {
			log.Fatal(err)
		}
	}
}
