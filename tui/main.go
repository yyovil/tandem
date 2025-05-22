package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	components "github.com/yyovil/tui/internal/components"
)

func main() {
	if _, err := tea.LogToFile("debug.log", "debug"); err != nil {
		os.Exit(1)
	}

	fpc := components.NewInput()
	if _, err := tea.NewProgram(&fpc, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		log.Println("Error running program:", err.Error())
		os.Exit(1)
	}
}
