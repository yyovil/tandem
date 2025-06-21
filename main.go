package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyovil/tandem/internal/app"
	"log"
	"os"
)

func main() {
	if _, err := tea.LogToFile("debug.log", "debug"); err != nil {
		os.Exit(1)
	}

	app := app.NewApp()
	if _, err := app.Run(); err != nil {
		log.Println("Error running program:", err.Error())
		os.Exit(1)
	}
}
