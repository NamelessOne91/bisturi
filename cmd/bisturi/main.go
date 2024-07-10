package main

import (
	"log"
	"os"

	models "github.com/NamelessOne91/bisturi/tui/models"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Getenv("BISTURI_DEBUG")) > 0 {
		f, err := tea.LogToFile("bisturi_debug.log", "debug")
		if err != nil {
			log.Fatal("Failed to setup logging:", err)
		}
		defer f.Close()
	}

	p := tea.NewProgram(models.NewBisturiModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
